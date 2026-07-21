---
tags: [concepto, lexico, algoritmo]
aliases: [AFD mínimo, estados equivalentes, refinamiento de particiones]
fuente: "Libro del Dragón, cap. 3.9.6"
fecha: 2026-07-10
---

# Minimización de AFD

Reduce el número de estados de un AFD fusionando **estados equivalentes**, por **refinamiento de particiones**: se parte de dos grupos (aceptación / no aceptación) y se subdividen mientras dos estados del mismo grupo lleven a grupos distintos con algún símbolo.

**Ejemplo trabajado** — AFD de `(a|b)*abb` (el de [[Construcción de subconjuntos]], estados A–E, acepta E):

| Ronda | Partición | ¿Por qué se divide? |
|-------|-----------|---------------------|
| 0 | `{A,B,C,D}` `{E}` | Inicial: no-aceptación / aceptación |
| 1 | `{A,B,C}` `{D}` `{E}` | Con `b`, D va a E (grupo aceptación); A, B y C no |
| 2 | `{A,C}` `{B}` `{D}` `{E}` | Con `b`, B va a D (grupo propio); A y C van a C |

Ninguna división más es posible → **A ≡ C** se fusionan y el AFD queda en **4 estados**. Antes/después:

```mermaid
stateDiagram-v2
    direction LR
    state "AFD original (5)" as orig {
        A --> B : a
        A --> C : b
        B --> B : a
        B --> D : b
        C --> B : a
        C --> C : b
        D --> B : a
        D --> E : b
        E --> B : a
        E --> C : b
    }
```

```mermaid
stateDiagram-v2
    direction LR
    state "AFD mínimo (4)" as min {
        AC --> B : a
        AC --> AC : b
        B --> B : a
        B --> D : b
        D --> B : a
        D --> E2 : b
        E2 --> B : a
        E2 --> AC : b
    }
```

## Relacionadas
- [[Construcción de subconjuntos]]
- [[Autómata finito (AFN y AFD)]]
- [[Cap 3 - Análisis léxico]]
