---
tags: [concepto, sintactico]
aliases: [GLC, CFG, BNF, gramática, Backus-Naur, producciones]
fuente: "Libro del Dragón, caps. 2 y 4"
fecha: 2026-07-10
---

# Gramática libre de contexto (BNF)

Notación para describir la sintaxis de un lenguaje. Tiene **4 componentes**: terminales (tokens), no terminales (variables), símbolo inicial y **producciones** `cabeza → cuerpo`.

Ejemplo: `instr → if ( expr ) instr else instr`.

**BNF** (Backus-Naur Form) es el formato exacto del **archivo de gramática** que se entrega en los proyectos (limpio, no copiado del archivo de CUP/Jison).

## Relacionadas
- [[Derivaciones y árbol de análisis sintáctico]]
- [[Ambigüedad, precedencia y asociatividad]]
- [[Recursividad por la izquierda y factorización]]
