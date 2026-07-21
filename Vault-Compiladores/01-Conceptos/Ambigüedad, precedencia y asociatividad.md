---
tags: [concepto, sintactico]
fuente: "Libro del Dragón, cap. 4.8"
fecha: 2026-07-10
---

# Ambigüedad, precedencia y asociatividad

Una gramática es **ambigua** si una cadena tiene ≥2 árboles de análisis (p. ej. `E → E+E | E*E`). Se resuelve declarando:
- **Precedencia:** qué operador se agrupa primero (`*` antes que `+`).
- **Asociatividad:** izquierda / derecha / no asociativa.

Ejemplo `id + id * id`: si `*` tiene más precedencia → se desplaza `*`. El ***else* colgante** se resuelve asociando el `else` con el `then` más cercano.

> Esto es literal en las tablas de precedencia de los proyectos: `precedence left/right` en [[CUP]] y `%left/%right` en [[Jison]].

## Relacionadas
- [[Conflictos shift-reduce y reduce-reduce]]
- [[Análisis sintáctico ascendente LR]]
- [[CompScript]]
