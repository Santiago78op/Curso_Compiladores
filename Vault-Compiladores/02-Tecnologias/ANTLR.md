---
tags: [tecnologia, lexico, sintactico, java, go]
fuente: "https://github.com/antlr/antlr4 (documentación oficial) + uso real en VLangCherry"
fecha: 2026-07-21
---

# ANTLR

Generador de **lexers y parsers** (ANother Tool for Language Recognition, v4). A diferencia de [[JFlex]]+[[CUP]] (dos herramientas separadas, LALR) o [[Jison]] (una sola herramienta, LALR/SLR/LR), ANTLR describe léxico y sintaxis en un **único archivo `.g4`** y genera un parser **LL(\*) adaptativo** (algoritmo ALL(\*)): decide la alternativa correcta mirando cuanta entrada haga falta, en vez de tablas LR precomputadas — evita a mano buena parte de los conflictos de [[Análisis sintáctico ascendente LR|LR]] y de [[Conflictos shift-reduce y reduce-reduce]].

## Estructura del `.g4`

Un archivo declara reglas de parser (minúscula) y reglas de lexer (MAYÚSCULA) en la misma gramática:

```antlr
grammar VLangCherry;

programa
    : declaracionGlobal* EOF
    ;

expr
    : expr '(' listaArgumentos? ')'      # exprLlamada
    | expr op=('*'|'/'|'%') expr         # exprMultiplicativa
    | expr op=('+'|'-') expr             # exprAditiva
    | literal                            # exprLiteral
    | ID                                  # exprIdentificador
    ;

ID      : [a-zA-Z_][a-zA-Z0-9_]* ;
WS      : [ \t\r\n]+ -> skip ;
```

- El **orden de las alternativas en una regla recursiva por la izquierda define la precedencia** (mayor a menor, de arriba hacia abajo): ANTLR reescribe internamente la recursión izquierda directa, algo que en una gramática LL manual habría que resolver aplicando [[Recursividad por la izquierda y factorización]].
- Las etiquetas `# nombreAlternativa` generan, por cada alternativa de una regla, un **tipo de contexto Go/Java concreto** (`ExprMultiplicativaContext`, `ExprAditivaContext`, ...) en vez de un único contexto genérico — esto es lo que permite recorrer el parse tree con un `switch` por tipo en vez de implementar el patrón *Visitor* completo (ver "Usado en").
- `-> skip` descarta el token (equivalente a no emitir acción en [[JFlex]]); `fragment` declara una subregla de lexer que no genera token propio (igual que en JFlex).

## Comando de generación

```
java -jar antlr.jar -Dlanguage=Go -visitor -o internal/parser -package parser grammar/VLangCherry.g4
```

Genera, en el paquete de salida: `<Gramatica>Lexer.go`, `<Gramatica>Parser.go`, `<Gramatica>BaseVisitor.go`, `<Gramatica>Visitor.go` — todos **código generado**, igual regla que con [[JFlex]]/[[CUP]]: nunca se edita a mano, cualquier cambio de léxico o sintaxis va en el `.g4` y se regenera.

## Runtime necesario

El código generado depende de una librería runtime por lenguaje objetivo (`github.com/antlr4-go/antlr/v4` en Go, `org.antlr:antlr4-runtime` en Java, `antlr4-runtime` en JS/TS/Python).

## Usado en

[[VLangCherry]] — parser generado para Go (`-Dlanguage=Go`); en vez de implementar el patrón Visitor completo (las ~60 firmas de `BaseVLangCherryVisitor`), el traductor propio (`internal/traductor/traductor.go`) hace un **type-switch directo** sobre los contextos etiquetados que ANTLR ya genera por cada alternativa, y produce un AST propio (`internal/ast/`) en vez de operar sobre el parse tree de ANTLR.

## Relacionadas
- [[JFlex]] · [[CUP]] · [[Jison]]
- [[Análisis sintáctico descendente LL(1)]]
- [[Análisis sintáctico ascendente LR]]
- [[Recursividad por la izquierda y factorización]]
- [[Árbol de sintaxis abstracta (AST)]]
- [[Codebase-Memory-MCP]] — código fuente real de antlr4 indexado como grafo consultable
