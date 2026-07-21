---
tags: [proyecto, lexico, sintactico, semantico]
fuente: "Enunciado OLC1_PT1_VD2024 (Proyecto 1, Vacaciones Dic 2024)"
fecha: 2026-07-10
---

# CompScript

**Objetivo:** **intérprete** completo con análisis léxico, sintáctico y **semántico**, y **AST**. Extensión `.cs`, *case insensitive*. Se entrega en 2 fases.

**Etapa del compilador:** front-end completo + intérprete sobre AST.

**Tecnologías:** [[JFlex]] + [[CUP]] (Java), GUI [[JavaFX y Scene Builder]], [[Maven]], AST con [[Graphviz]] (y [[Mermaid]] en el vault).

**Rasgos:** tipos `int, double, bool, char, string`; tablas de compatibilidad de tipos; casteos; `if/match`; `while/for/do-while`; `break/continue/return`; vectores, listas dinámicas, structs, funciones, métodos, `RUN_MAIN`. Reportes: tokens, errores (léx/sint/sem), AST, [[Tabla de símbolos|tabla de símbolos]] por entorno.

**Conceptos aplicados:** [[Árbol de sintaxis abstracta (AST)]] · [[Comprobación de tipos]] · [[Conversión de tipos (coerción y cast)]] · [[Flujo de control y switch]] · [[Entornos y alcance]] · [[Registro de activación y pila de control]]

## Guía de elaboración
- [[Guía CompScript]]

## Relacionadas
- [[CompInterpreter]]
- [[Cap 5 - Traducción dirigida por la sintaxis]]
- [[Cap 7 - Entornos en tiempo de ejecución]]
