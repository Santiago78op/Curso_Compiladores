---
tags: [concepto, compiladores]
fuente: "Libro del Dragón, caps. 3.1.4, 4.1.3, 6.5"
fecha: 2026-07-10
---

# Manejo de errores (léxicos, sintácticos, semánticos)

- **Léxico:** un carácter/lexema no reconocido (ej. `$` fuera del lenguaje).
- **Sintáctico:** la secuencia de tokens no respeta la gramática (ej. "se esperaba Expresión").
- **Semántico:** viola reglas de significado (ej. "no se puede sumar CADENA y CADENA", usar variable no declarada, reasignar `const`) — estos chequeos se hacen consultando la [[Tabla de símbolos]] y los [[Entornos y alcance|entornos]].

**Recuperación** típica: **modo pánico** (descartar hasta un token de sincronización de FOLLOW). Los 4 proyectos generan un **reporte de errores** con tipo, descripción, línea y columna.

## Relacionadas
- [[Comprobación de tipos]]
- [[FIRST y FOLLOW]]
- [[Fases de un compilador]]
