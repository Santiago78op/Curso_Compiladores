---
tags: [concepto, semantico]
fuente: "Libro del Dragón, cap. 6.5.2"
fecha: 2026-07-10
---

# Conversión de tipos (coerción y cast)

- **Ampliación (widening):** preserva información (`char → int → float`).
- **Reducción (narrowing):** puede perder información (sentido inverso).
- **Coerción (implícita):** el compilador la inserta solo (ej. `(float)2` en `2 * 3.14`).
- **Cast (explícita):** el programador la escribe.

> El `CAST(exp AS tipo)` de CompScript/CompInterpreter es la conversión **explícita**; las promociones int→double en operaciones son **coerciones** implícitas.

## Relacionadas
- [[Comprobación de tipos]]
- [[CompScript]]
- [[CompInterpreter]]
