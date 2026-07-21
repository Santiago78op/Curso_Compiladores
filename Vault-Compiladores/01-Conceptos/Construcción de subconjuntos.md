---
tags: [concepto, lexico, algoritmo]
fuente: "Libro del Dragón, cap. 3.7.1"
fecha: 2026-07-10
---

# Construcción de subconjuntos

Convierte un **AFN → AFD**. Cada estado del AFD es un **conjunto** de estados del AFN. Usa dos operaciones:
- **ε-cerradura(T):** estados alcanzables desde `T` solo con ε.
- **mover(T, a):** estados alcanzables desde `T` con el símbolo `a`.

Regla: `Dtran[A, a] = ε-cerradura(mover(A, a))`.

Ejemplo del libro para `(a|b)*abb`:

```mermaid
stateDiagram-v2
    direction LR
    [*] --> A
    A --> B: a
    A --> C: b
    B --> B: a
    B --> D: b
    C --> B: a
    C --> C: b
    D --> B: a
    D --> E: b
    E --> B: a
    E --> C: b
    E --> [*]
```

Luego se aplica [[Minimización de AFD]] para reducir estados.

## Relacionadas
- [[Construcción de Thompson]]
- [[Minimización de AFD]]
- [[Autómata finito (AFN y AFD)]]
