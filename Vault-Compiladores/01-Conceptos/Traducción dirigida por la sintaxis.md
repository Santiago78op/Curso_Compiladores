---
tags: [concepto, semantico]
aliases: [SDT, SDD, TDS, syntax-directed translation, reglas semánticas, acciones semánticas]
fuente: "Libro del Dragón, cap. 5"
fecha: 2026-07-10
---

# Traducción dirigida por la sintaxis

Asocia a cada producción de la gramática **reglas semánticas** (SDD) o **acciones** dentro de la producción (SDT) que calculan [[Atributos sintetizados y heredados|atributos]].

Una SDD con **solo atributos sintetizados** (S-atribuida) se implementa directo con un parser LR → por eso las acciones `{: RESULT = … :}` de [[CUP]] y `{ $$ = … }` de [[Jison]] son reglas semánticas.

Aplicación típica: construir el [[Árbol de sintaxis abstracta (AST)]] o **evaluar** expresiones.

## Relacionadas
- [[Atributos sintetizados y heredados]]
- [[Árbol de sintaxis abstracta (AST)]]
- [[Cap 5 - Traducción dirigida por la sintaxis]]
