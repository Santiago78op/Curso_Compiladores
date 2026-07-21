---
tags: [concepto, lexico]
fuente: "Libro del Dragón, cap. 3"
fecha: 2026-07-10
---

# Token, lexema y patrón

- **Token:** par `⟨nombre-token, atributo⟩`; la **categoría** léxica (ID, NUMBER, IF…). Es lo que consume el analizador sintáctico.
- **Lexema:** la **secuencia concreta** de caracteres que casó (p. ej. el texto `posicion`).
- **Patrón:** la regla ([[Expresiones regulares|expresión regular]]) que describe qué lexemas forman un token.

Ejemplo: en `posicion = inicial + 60`, `posicion` es un lexema del token `⟨id⟩`, `60` del token `⟨número⟩`.

> El **reporte de tokens** de los proyectos (# / lexema / tipo / línea / columna) sale de aquí.

## Relacionadas
- [[Expresiones regulares]]
- [[Tabla de símbolos]]
- [[JFlex]]
