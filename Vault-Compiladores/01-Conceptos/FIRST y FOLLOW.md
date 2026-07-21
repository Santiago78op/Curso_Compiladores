---
tags: [concepto, sintactico, algoritmo]
fuente: "Libro del Dragón, cap. 4.4.2"
fecha: 2026-07-10
---

# FIRST y FOLLOW

Dos funciones que guían la construcción de analizadores sintácticos.

- **FIRST(α):** terminales con los que pueden **empezar** las cadenas derivadas de α (incluye ε si α ⇒ ε).
- **FOLLOW(A):** terminales que pueden aparecer **inmediatamente a la derecha** de A; incluye `$` si A puede ser el último símbolo.

Ejemplo del libro (gramática de expresiones):

| No terminal | FIRST | FOLLOW |
|---|---|---|
| E | `(, id` | `), $` |
| E' | `+, ε` | `), $` |
| T | `(, id` | `+, ), $` |
| T' | `*, ε` | `+, ), $` |
| F | `(, id` | `+, *, ), $` |

Con estos conjuntos se llena la tabla de [[Análisis sintáctico descendente LL(1)|LL(1)]].

## Relacionadas
- [[Análisis sintáctico descendente LL(1)]]
- [[Recursividad por la izquierda y factorización]]
- [[Cap 4 - Análisis sintáctico]]
