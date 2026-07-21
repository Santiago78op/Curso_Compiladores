---
tags: [proyecto, lexico, sintactico, semantico, web]
fuente: "Enunciado OLC1_PT2_2S2024 (Proyecto 2, 2do Sem 2024)"
fecha: 2026-07-10
---

# CompInterpreter

**Objetivo:** **intérprete** léxico/sintáctico/semántico con **arquitectura cliente-servidor** (web). Extensión `.ci`, *case insensitive*.

**Etapa del compilador:** front-end completo + intérprete sobre AST, en la web.

**Tecnologías:** [[Jison]] (JS/TS), cliente en **React** (o Angular/Vue), servidor **Express (Node)** vía REST, AST interactivo con [[vis-network]] (y [[Mermaid]] en el vault).

**Rasgos:** tipos `int, double, bool, char, string, null`; operador **ternario**; `switch case` con *fall-through*; `while/for/do-until/loop`; funciones/métodos; `echo`; funciones nativas (`lower, upper, round, length, truncate, is, toString, toCharArray, reverse, max, min, sum, average`); sentencia `ejecutar`. Reportes: errores, AST, [[Tabla de símbolos|tabla de símbolos]], consola.

**Conceptos aplicados:** [[Árbol de sintaxis abstracta (AST)]] · [[Comprobación de tipos]] · [[Flujo de control y switch]] · [[Entornos y alcance]] · [[Análisis sintáctico ascendente LR]]

## Guía de elaboración
- [[Guía CompInterpreter]]

## Relacionadas
- [[CompScript]]
- [[Jison]]
- [[Cap 6 - Generación de código intermedio]]
