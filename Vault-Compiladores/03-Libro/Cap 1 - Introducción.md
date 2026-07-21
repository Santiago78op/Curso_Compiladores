---
tags: [libro, compiladores]
fuente: "Libro del Dragón (Aho, Lam, Sethi, Ullman), 2ª ed., cap. 1"
fecha: 2026-07-10
---

# Cap 1 - Introducción

Panorama de la compilación. Un compilador traduce fuente→destino; un intérprete ejecuta directamente ([[Compilador vs intérprete]]). Presenta las **[[Fases de un compilador|fases]]**: análisis (front-end) y síntesis (back-end).

**Ejemplo guía:** `posicion = inicial + velocidad * 60` recorre léxico → sintáctico → semántico (inserta coerción `inttofloat`) → intermedio → optimización → código.

Conceptos base: [[Token, lexema y patrón]], [[Tabla de símbolos]], estático vs dinámico.

> Aplica a los 4 proyectos como marco general.

## Relacionadas
- [[Fases de un compilador]]
- [[Compilador vs intérprete]]
- [[Cap 2 - Traductor simple orientado a la sintaxis]]
