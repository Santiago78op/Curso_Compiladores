---
tags: [tecnologia, lexico, sintactico, javascript]
fuente: "https://gerhobbelt.github.io/jison/docs/"
fecha: 2026-07-10
---

# Jison

Generador de parsers para **JavaScript** (equivalente a Bison/Yacc). Un solo archivo `.jison` incluye **lexer** (`%lex`) **y** gramática. Soporta LALR(1) (por defecto), SLR, LR(1), y **EBNF** (`* + ? ()`).

## Estructura
```jison
%lex
%%
\s+          /* ignora */
[0-9]+       return 'NUMBER';
"+"          return '+';
/lex

%left '+'
%left '*'

%%
e : e '+' e   { $$ = $1 + $3; }
  | NUMBER    { $$ = Number(yytext); }
  ;
```

- Acciones en JS: `$$` (resultado), `$1/$2` (por posición), `@1` (ubicación), objeto `yy` compartido.
- **Uso:** `jison calc.jison` → `calc.js`; en Node `require('./calc').parse('1+2')`; también en navegador.

## Usado en
[[CompInterpreter]]

## Relacionadas
- [[JFlex]]
- [[CUP]]
- [[Análisis sintáctico ascendente LR]]
- [[Codebase-Memory-MCP]] — código fuente real de jison indexado como grafo consultable
