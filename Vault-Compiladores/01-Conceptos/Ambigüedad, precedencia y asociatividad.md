---
tags: [concepto, sintactico]
aliases: ["precedencia por niveles", "gramática no ambigua", "n+1 no terminales", "expr term factor", "gramática ambigua más precedence", "reglas por defecto de Yacc"]
fuente: "Libro del Dragón §2.2.5–2.2.6 (precedencia estructural) y §4.8.1 (ambigua + precedence)"
fecha: 2026-07-24
---

# Ambigüedad, precedencia y asociatividad

Una gramática es **ambigua** si una cadena tiene ≥2 árboles de análisis (p. ej. `E → E+E | E*E`, o `9+5*2` que sin reglas admite `(9+5)*2` y `9+(5*2)`). Se resuelve fijando:
- **Precedencia:** qué operador agrupa primero (`*` recibe sus operandos antes que `+`).
- **Asociatividad:** izquierda / derecha / no asociativa.

Hay **dos formas** de imponerlas, y conviene no confundirlas:

## 1. Por declaraciones (en el generador)

`precedence left/right` en [[CUP]] y `%left/%right/%nonassoc` en [[Jison]] (y `%prec` para el menos unario). El generador usa esas declaraciones para resolver los [[Conflictos shift-reduce y reduce-reduce|conflictos shift-reduce]] de una gramática ambigua, sin cambiar la gramática. El *else* colgante se resuelve así, asociando el `else` con el `if` abierto más cercano (SHIFT).

## 2. Por la estructura de la gramática (§2.2.6 + recuadro "Generalización", pp. 49–50)

Se codifica la precedencia **en la gramática misma**, sin declaraciones. Un `factor` es "una expresión que **no puede separarse** por ningún operador" (un operando o una `( expr )` protegida); un `term` solo se separa por `* /`; una `expr` por cualquiera:

```
expr   → expr + term | expr - term | term
term   → term * factor | term / factor | factor
factor → dígito | ( expr )
```

**Regla general:** para *n* niveles de precedencia se usan **n+1 no terminales**; el primero (`factor`) nunca se separa, y cada nivel apunta al de arriba. La recursión **izquierda** (`expr → expr + term`) da asociación por la izquierda.

## ¿Cuándo cada una? (§4.8.1, pp. 278–281)

- La **BNF entregable** debe ser **no ambigua por sí sola** (no lleva tablas adjuntas) → se escribe con la construcción **por niveles** (forma 2).
- El **`.cup`/`.jison` real** prefiere la gramática **ambigua + `precedence`** (forma 1) por **ingeniería**: el parser tiene **menos estados** y **no pierde tiempo reduciendo producciones simples** (`E→T`, `T→F`) cuya única función es codificar precedencia. Misma precedencia, autómata más chico.

### Reglas por defecto de Yacc/CUP ante conflictos

Aunque no se declare nada, el generador resuelve:
- **shift/reduce → desplazar** (SHIFT). Esto resuelve el *else* colgante gratis (el `else` se pega al `if` más cercano).
- **reduce/reduce → la primera producción** listada en el archivo.

Conviene no depender de estos defaults salvo el del *else* colgante; los demás conflictos suelen indicar una gramática mal escrita (ver [[Conflictos shift-reduce y reduce-reduce]]).

## Relacionadas
- [[Gramática libre de contexto (BNF)]]
- [[Conflictos shift-reduce y reduce-reduce]]
- [[Recursividad por la izquierda y factorización]]
- [[Elementos LR(0) y la tabla SLR]]
- [[CUP]]
