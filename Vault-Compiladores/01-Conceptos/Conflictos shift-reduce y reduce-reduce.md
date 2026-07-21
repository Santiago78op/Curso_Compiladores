---
tags: [concepto, sintactico]
fuente: "Libro del Dragón, caps. 4.5.4 y 4.8"
fecha: 2026-07-10
---

# Conflictos shift-reduce y reduce-reduce

En el [[Análisis sintáctico ascendente LR|análisis LR]] pueden surgir conflictos:
- **Shift-reduce:** el parser no sabe si desplazar o reducir (p. ej. el ***else* colgante**). Suele resolverse con **precedencia/asociatividad** o desplazando por defecto.
- **Reduce-reduce:** dos producciones podrían reducir en el mismo punto (síntoma de gramática mal diseñada).

Los generadores ([[CUP]], [[Jison]]) reportan estos conflictos al construir las tablas → hay que depurar la gramática o declarar precedencias.

## Relacionadas
- [[Ambigüedad, precedencia y asociatividad]]
- [[Análisis sintáctico ascendente LR]]
- [[CUP]]
