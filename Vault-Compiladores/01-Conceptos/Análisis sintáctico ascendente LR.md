---
tags: [concepto, sintactico]
aliases: [LR, SLR, LALR, "LALR(1)", bottom-up, shift-reduce, ascendente]
fuente: "Libro del Dragón, caps. 4.5–4.7"
fecha: 2026-07-10
---

# Análisis sintáctico ascendente LR

Construye el árbol de las **hojas a la raíz** con una pila y dos acciones: **desplazar** (shift) y **reducir** (reduce, cuando la cima es el cuerpo de una producción = *mango*). Es lo que hacen [[CUP]] y [[Jison]].

Familia, de menor a mayor poder:
```mermaid
flowchart LR
    LR0["LR(0)"] --> SLR["SLR (usa FOLLOW)"]
    SLR --> LALR["LALR(1) — CUP y Jison"]
    LALR --> LR1["LR(1) canónico"]
```

**LALR(1)** es el algoritmo por defecto de CUP y Jison: tablas compactas con casi el poder de LR(1).

## Relacionadas
- [[Conflictos shift-reduce y reduce-reduce]]
- [[Análisis sintáctico descendente LL(1)]]
- [[Ambigüedad, precedencia y asociatividad]]
