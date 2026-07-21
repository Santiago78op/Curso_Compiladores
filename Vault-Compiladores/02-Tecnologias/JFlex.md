---
tags: [tecnologia, lexico, java]
fuente: "https://www.jflex.de/manual.html"
fecha: 2026-07-10
---

# JFlex

Generador de **analizadores léxicos** para Java. Lee un `.flex` (ER + acciones Java) y genera una clase con `yylex()` que devuelve tokens. Internamente construye un [[Autómata finito (AFN y AFD)|AFD]].

## Estructura del `.flex` (tres secciones separadas por `%%`)
```jflex
import java_cup.runtime.*;
%%
%class Lexer
%unicode
%cup        // produce Symbol de CUP
%line
%column
%{
  // OJO: symbol() NO lo genera JFlex — es un helper que definís vos
  // en este bloque %{...%} (código copiado tal cual dentro de la clase).
  private Symbol symbol(int type) {
    return new Symbol(type, yyline, yycolumn);
  }
  private Symbol symbol(int type, Object value) {
    return new Symbol(type, yyline, yycolumn, value);
  }
%}
Digit = [0-9]
Id    = [a-zA-Z_][a-zA-Z0-9_]*
%state STRING
%%
<YYINITIAL> {
  "if"     { return symbol(sym.IF); }
  {Digit}+ { return symbol(sym.NUMBER, yytext()); }
  {Id}     { return symbol(sym.ID, yytext()); }
  [^]      { /* error léxico */ }
}
```

**API:** `yytext()`, `yyline`, `yycolumn`, `yybegin(ESTADO)`, `yypushback(n)`.
**Reglas clave:** *longest match* + prioridad de la regla escrita primero (palabras reservadas antes que `Id`).
**Comando:** `jflex Lexer.flex`

Equivale al **Lex** del libro y a la sección `%lex` de [[Jison]]. Se integra con [[CUP]] vía `%cup`.

## Usado en
[[DataForge]], [[ConjAnalyzer]], [[CompScript]]

## Relacionadas
- [[CUP]]
- [[Definiciones regulares]]
- [[Cap 3 - Análisis léxico]]
- [[Codebase-Memory-MCP]] — código fuente real de jflex indexado como grafo consultable
