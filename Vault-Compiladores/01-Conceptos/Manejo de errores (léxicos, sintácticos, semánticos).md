---
tags: [concepto, compiladores]
aliases: ["modo pánico", "recuperación de errores", "pánico léxico", "pánico sintáctico", "reporte de errores"]
fuente: "Libro del Dragón §3.1.4 (léxico), §4.1.3–4.1.4 (sintáctico), cap. 6 (semántico)"
fecha: 2026-07-24
---

# Manejo de errores (léxicos, sintácticos, semánticos)

- **Léxico:** un carácter/lexema no reconocido (ej. `$` fuera del lenguaje). Ojo: el lexer casi no detecta errores por sí solo — `fi (a==f(x))` devuelve `id` (no sabe si `fi` es un `if` mal escrito), y es otra fase la que se queja (§3.1.4).
- **Sintáctico:** la secuencia de tokens no respeta la gramática (ej. "se esperaba Expresión").
- **Semántico:** viola reglas de significado (ej. "no se puede sumar CADENA y CADENA", usar variable no declarada, reasignar `const`) — estos chequeos se hacen consultando la [[Tabla de símbolos]] y los [[Entornos y alcance|entornos]].

## Recuperación: dos "modos pánico" DISTINTOS

Conviene no confundirlos —operan sobre unidades diferentes—:

- **Pánico léxico (§3.1.4).** Cuando ningún patrón calza un prefijo de la entrada, se **eliminan caracteres sucesivos** hasta poder formar un token bien formado al principio de lo que queda. Alternativas de reparación de un solo carácter: **borrar / insertar / sustituir / transponer** (la mayoría de los errores léxicos son de un carácter). En los proyectos: *"léxico descarta el carácter"*.
- **Pánico sintáctico (§4.1.4).** Cuando el parser no puede continuar, **descarta tokens** hasta encontrar un **token de sincronización** (típicamente del conjunto [[FIRST y FOLLOW|FOLLOW]] del no terminal en curso, p. ej. `;`). En CUP/Jison esto se escribe con el token especial `error`: `instruccion ::= error PUNTO_COMA`.

Los errores **semánticos** no usan modo pánico: se reportan en el punto donde se detectan y (según el proyecto) se aborta o se acumula y sigue.

Los 4 proyectos generan un **reporte de errores** con tipo, descripción, línea y columna.

## Relacionadas
- [[Comprobación de tipos]]
- [[FIRST y FOLLOW]]
- [[CUP]]
- [[Fases de un compilador]]
