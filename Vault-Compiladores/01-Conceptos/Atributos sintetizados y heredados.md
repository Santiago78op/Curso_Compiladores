---
tags: [concepto, semantico]
fuente: "Libro del Dragón, cap. 5.1"
fecha: 2026-07-10
---

# Atributos sintetizados y heredados

- **Sintetizado:** se calcula desde los valores de los **hijos** (de abajo hacia arriba). Ej.: `E.val = E1.val + T.val`.
- **Heredado:** se define desde el **padre**, **hermanos** o el propio nodo (de arriba/lados). Ej.: propagar el tipo a una lista de variables.

Los terminales solo tienen atributos sintetizados (su valor lo da el léxico). Una SDD **S-atribuida** (solo sintetizados) encaja con parsing LR; una **L-atribuida** con parsing LL.

## Relacionadas
- [[Traducción dirigida por la sintaxis]]
- [[Árbol de sintaxis abstracta (AST)]]
- [[Comprobación de tipos]]
