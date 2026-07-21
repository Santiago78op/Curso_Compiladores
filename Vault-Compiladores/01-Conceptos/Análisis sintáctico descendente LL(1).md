---
tags: [concepto, sintactico]
aliases: [LL(1), top-down, descenso recursivo, análisis predictivo, parser predictivo]
fuente: "Libro del Dragón, cap. 4.4"
fecha: 2026-07-10
---

# Análisis sintáctico descendente LL(1)

Construye el árbol de la **raíz a las hojas**. **LL(1)** = leer de izquierda a derecha, derivación por la izquierda, 1 token de anticipación.

Una gramática es LL(1) si para toda `A → α | β`: [[FIRST y FOLLOW|FIRST]](α) y FIRST(β) son disjuntos, a lo sumo una deriva ε, etc. Con FIRST/FOLLOW se llena la **tabla predictiva** `M[A,a]`. El **descenso recursivo** es su implementación "a mano" (un método por no terminal).

Orden de construcción **top-down** para `id + id` con `E → T E'`, `E' → + T E' | ε`, `T → id` (los números ①–⑥ marcan en qué paso se expande cada nodo, siempre el no terminal **más a la izquierda**, mirando 1 token):

```mermaid
flowchart TD
    E["① E"] --> T["② T"]
    E --> Ep["④ E'"]
    T --> id1["③ id"]
    Ep --> mas["+"]
    Ep --> T2["⑤ T"]
    Ep --> Ep2["⑥ E' → ε"]
    T2 --> id2["id"]
```

> No permite gramáticas recursivas por la izquierda ni ambiguas ([[Recursividad por la izquierda y factorización]]).

## Relacionadas
- [[FIRST y FOLLOW]]
- [[Análisis sintáctico ascendente LR]]
- [[Derivaciones y árbol de análisis sintáctico]]
