---
tags: [concepto, sintactico]
fuente: "Libro del Dragón, cap. 4.3"
fecha: 2026-07-10
---

# Recursividad por la izquierda y factorización

Dos transformaciones para preparar una gramática para el análisis descendente:

- **Eliminar recursividad por la izquierda:** `A → Aα | β` (rompe a los parsers descendentes) se reescribe:
  ```
  A  → β A'
  A' → α A' | ε
  ```
- **Factorización por la izquierda:** cuando dos producciones comparten prefijo, se factoriza para decidir con un solo token de anticipación.

Ninguna gramática recursiva por la izquierda o ambigua puede ser [[Análisis sintáctico descendente LL(1)|LL(1)]].

## Relacionadas
- [[Gramática libre de contexto (BNF)]]
- [[FIRST y FOLLOW]]
- [[Análisis sintáctico descendente LL(1)]]
