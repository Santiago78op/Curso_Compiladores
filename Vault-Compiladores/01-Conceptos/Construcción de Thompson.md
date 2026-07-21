---
tags: [concepto, lexico, algoritmo]
aliases: [Thompson, McNaughton-Yamada-Thompson, ER a AFN, regex a NFA]
fuente: "Libro del Dragón, cap. 3.7.4"
fecha: 2026-07-10
---

# Construcción de Thompson

Algoritmo (McNaughton-Yamada-Thompson) que convierte una [[Expresiones regulares|expresión regular]] en un **AFN**, componiendo fragmentos: un fragmento para ε, uno por símbolo, y reglas para `|` (unión), concatenación y `*` (cerradura). Cada fragmento tiene un estado inicial y uno de aceptación conectados por transiciones‑ε.

**Fragmento para la unión `a|b`** — un inicio y un fin nuevos, conectados por ε a los fragmentos de `a` y de `b`:

```mermaid
stateDiagram-v2
    direction LR
    [*] --> i
    i --> p1 : ε
    i --> q1 : ε
    p1 --> p2 : a
    q1 --> q2 : b
    p2 --> f : ε
    q2 --> f : ε
    f --> [*]
```

**Fragmento para la cerradura `a*`** — permite cero repeticiones (ε directo de `i` a `f`) o volver a empezar (ε de vuelta):

```mermaid
stateDiagram-v2
    direction LR
    [*] --> i
    i --> p : ε
    p --> q : a
    q --> p : ε
    q --> f : ε
    i --> f : ε
    f --> [*]
```

La **concatenación `ab`** simplemente fusiona el estado final del fragmento de `a` con el inicial del de `b` (sin ε extra).

Es el primer paso del pipeline que luego pasa a AFD con la [[Construcción de subconjuntos]].

## Relacionadas
- [[Autómata finito (AFN y AFD)]]
- [[Construcción de subconjuntos]]
- [[Expresiones regulares]]
