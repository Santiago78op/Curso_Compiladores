---
tags: [concepto, compiladores]
aliases: ["modo pánico", "recuperación de errores", "pánico léxico", "pánico sintáctico", "token error", "producción de error", "yyerrok", "reporte de errores"]
fuente: "Libro del Dragón §3.1.4 (léxico), §4.1.4 y §4.9.4 (sintáctico/token error), cap. 6 (semántico)"
fecha: 2026-07-24
---

# Manejo de errores (léxicos, sintácticos, semánticos)

- **Léxico:** un carácter/lexema no reconocido (ej. `$` fuera del lenguaje). Ojo: el lexer casi no detecta errores por sí solo — `fi (a==f(x))` devuelve `id` (no sabe si `fi` es un `if` mal escrito), y es otra fase la que se queja (§3.1.4).
- **Sintáctico:** la secuencia de tokens no respeta la gramática (ej. "se esperaba Expresión").
- **Semántico:** viola reglas de significado (ej. "no se puede sumar CADENA y CADENA", usar variable no declarada, reasignar `const`) — se chequea consultando la [[Tabla de símbolos]] y los [[Entornos y alcance|entornos]].

## Recuperación: dos "modos pánico" DISTINTOS

- **Pánico léxico (§3.1.4).** Cuando ningún patrón calza un prefijo de la entrada, se **eliminan caracteres sucesivos** hasta poder formar un token bien formado al principio de lo que queda. Reparaciones de un solo carácter: **borrar / insertar / sustituir / transponer**. En los proyectos: *"léxico descarta el carácter"*.
- **Pánico sintáctico (§4.1.4).** El parser **descarta tokens** hasta un **token de sincronización** (típicamente del [[FIRST y FOLLOW|FOLLOW]] del no terminal en curso, p. ej. `;`).

## El token `error` de CUP/Yacc (§4.9.4) — la teoría del `error PUNTO_COMA`

Es la implementación concreta del pánico sintáctico en un parser LR, y es exactamente lo que hace la línea `instruccion ::= error PUNTO_COMA` del `parser.cup` de DataForge:

1. El usuario agrega **producciones de error** `A → error α` a los no terminales "importantes" (instrucción, expresión, bloque…); `error` es un token reservado.
2. Al fallar, el parser **saca estados de la pila** hasta encontrar uno cuyo conjunto de [[Elementos LR(0) y la tabla SLR|elementos LR(0)]] incluya `A → · error α`, y **desplaza** el token ficticio `error`.
3. Luego **descarta símbolos de entrada** hasta poder reducir `error α` a `A` (si `α` son terminales como `;`, avanza hasta encontrarlos) y continúa el análisis normal.
4. `yyerrok` (o el `RESULT` de la acción en CUP) le indica que ya salió del modo de error.

El efecto: **un error por instrucción rota, y el análisis sigue** — justo lo que exige el reporte de errores del curso. (Las otras 3 estrategias del §4.1.4: nivel de frase, producciones de error, corrección global.)

## Los semánticos no tienen modo pánico

Los errores **semánticos** se reportan donde se detectan, y qué hacer después **no lo define el libro sino el enunciado**: abortar, propagar un valor nulo, o acumular y seguir. Los cinco proyectos del curso eligieron tres caminos distintos y los tres se defienden — la comparación completa está en [[Políticas de error semántico]].

Los 4 proyectos generan un **reporte de errores** con tipo, descripción, línea y columna.

> **Detalle de implementación que se paga caro:** si el pipeline es *best-effort* (sigue traduciendo aunque el parseo haya fallado, para juntar todos los errores), hay que evitar que el panic esperable de recorrer un árbol roto se presente al usuario como un "error interno". Junto a un error sintáctico real, ese ruido hace parecer que el intérprete se rompió cuando en realidad detectó bien el problema. Bug real corregido en [[VLangCherry]].

## Relacionadas
- [[Políticas de error semántico]]
- [[Comprobación de tipos]]
- [[FIRST y FOLLOW]]
- [[Elementos LR(0) y la tabla SLR]]
- [[CUP]]
- [[Fases de un compilador]]
