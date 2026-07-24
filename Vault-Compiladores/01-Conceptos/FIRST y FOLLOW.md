---
tags: [concepto, sintactico, algoritmo]
aliases: ["PRIMERO", "SIGUIENTE", "reglas de FIRST", "reglas de FOLLOW", "calcular primero y siguiente"]
fuente: "Libro del Dragón §4.4.2, pp. 220–222"
fecha: 2026-07-24
---

# FIRST y FOLLOW

Dos funciones que guían la construcción de analizadores sintácticos.

- **FIRST(α):** terminales con los que pueden **empezar** las cadenas derivadas de α (incluye ε si α ⇒ ε).
- **FOLLOW(A):** terminales que pueden aparecer **inmediatamente a la derecha** de A; incluye `$` si A puede ser el último símbolo.

## Reglas de cálculo (a mano, hasta punto fijo)

**FIRST(X)** — repetir hasta que no se agregue nada:
1. Si `X` es **terminal**, `FIRST(X) = {X}`.
2. Si `X → Y₁ Y₂ … Yₖ`: agregar a `FIRST(X)` todo lo de `FIRST(Yᵢ)` (menos ε) **si** `Y₁ … Yᵢ₋₁ ⇒ ε` (es decir, ε ∈ FIRST de todos los anteriores). Si **todos** los `Yⱼ` derivan ε, agregar **ε** a `FIRST(X)`.
3. Si `X → ε` es producción, agregar **ε** a `FIRST(X)`.

Para una cadena `X₁X₂…Xₙ`: se acumula `FIRST(X₁)` sin ε; se sigue con `FIRST(X₂)` solo si ε ∈ FIRST(X₁); y así. ε entra solo si todos derivan ε.

**FOLLOW(A)** — repetir hasta que no se agregue nada:
1. `$` ∈ `FOLLOW(S)` (S = símbolo inicial).
2. Si `A → α B β`: todo `FIRST(β)` **excepto ε** ⊆ `FOLLOW(B)`.
3. Si `A → α B`, o `A → α B β` con **ε ∈ FIRST(β)**: `FOLLOW(A)` ⊆ `FOLLOW(B)`.

> ε **nunca** entra en un conjunto FOLLOW (solo terminales y `$`). La regla 3 es la que se olvida en examen: lo que sigue a `A` también sigue a `B` cuando `B` está al final (o lo que le sigue puede desaparecer).

## Ejemplo del libro (gramática de expresiones no recursiva izq., ej. 4.30)

| No terminal | FIRST | FOLLOW |
|---|---|---|
| E | `(, id` | `), $` |
| E' | `+, ε` | `), $` |
| T | `(, id` | `+, ), $` |
| T' | `*, ε` | `+, ), $` |
| F | `(, id` | `+, *, ), $` |

Con estos conjuntos se llena la tabla de [[Análisis sintáctico descendente LL(1)|LL(1)]] (por FIRST del cuerpo, y por FOLLOW cuando el cuerpo deriva ε), y también deciden los *reduce* de la [[Elementos LR(0) y la tabla SLR|tabla SLR]].

## Relacionadas
- [[Análisis sintáctico descendente LL(1)]]
- [[Elementos LR(0) y la tabla SLR]]
- [[Recursividad por la izquierda y factorización]]
- [[Cap 4 - Análisis sintáctico]]
