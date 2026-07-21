---
tags: [concepto, semantico]
fuente: "Libro del Dragón, caps. 6.6 y 6.8"
fecha: 2026-07-10
---

# Flujo de control y switch

- **`if`, `while`, `for`:** se evalúa una **expresión booleana** y se elige la rama. Con **corto circuito**, `&&`/`||` solo evalúan lo necesario.
- **`switch`/`match`:** **bifurcación de n vías**: evaluar E, buscar el caso igual (o `default`), ejecutar su instrucción.

> En un [[Compilador vs intérprete|intérprete]] no se generan etiquetas/gotos: se recorre el [[Árbol de sintaxis abstracta (AST)|AST]] y se ejecuta la rama directamente. Base del `if/else if`, ciclos, `match` (CompScript) y `switch case` con *fall-through* (CompInterpreter).

## Relacionadas
- [[Árbol de sintaxis abstracta (AST)]]
- [[Comprobación de tipos]]
- [[CompInterpreter]]
