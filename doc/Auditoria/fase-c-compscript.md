# Fase C.2 — Code review de CompScript (Java + JFlex + CUP + AST)

**Qué se revisó:** ~2,600 renglones leídos completos — `interprete/` (Senales, Entorno, Contexto, Tipo, Valor, Simbolo, Operaciones, Interprete), `ast/A.java` (892), `parser.cup`, `Lexer.flex`, `reportes/Reportes.java` — más el enunciado (`OLC1_PT1_ VD2024.clean.md`) en los puntos dudosos. `EditorApp`/`TestInterprete`/`Lanzador` solo por patrón (idénticos al esquema DataForge ya verificado).
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto

**El proyecto más sólido de los revisados hasta ahora.** Los dos puntos donde VLangCherry falla, CompScript los resuelve bien: `return;` en función no-void → error; `return valor` en void → error; cortocircuito implementado (con el comentario del bug corregido en el código); `StackOverflowError` capturado y reportado como "desbordamiento de pila" (Java sí puede — Go no); señales que escapan capturadas con mensaje correcto; `const` validado en los 4 caminos de mutación. No hay hallazgos ALTOS. Quedan 1 medio, 2 medio-bajos con decisión de diseño pendiente, y varios bajos.

### Verificaciones que pasaron (escudo pro-defensa)

| Punto | Resultado |
|---|---|
| Default de `bool` | `true` — y la tabla 5.3 del enunciado dice exactamente eso. **Trampa de examen implementada bien** (cualquiera esperaría `false`) |
| `let p2: Persona = p1;` (sospecha de init ignorado) | **No existe**: la gramática solo admite struct sin init o con literal (`parser.cup:131-134`) — es error sintáctico, no silencio |
| Cortocircuito `&&`/`||` | ✅ `A.java:165-181`, con el test `x != 0 && 10/x > 1` citado en el comentario |
| `return` en void / sin valor en no-void | ✅ ambos errores reportados (`invocar`, A.java:869-878) |
| Guarda de `validarVector` (bug de la auditoría previa) | ✅ presente (A.java:461-469) |
| `parser.raiz` | inicializada en la declaración — no hay NPE aunque el parse aborte |
| Tabla de símbolos | registro histórico (una fila por declaración) **y** con valor final (Simbolo.valor mutable) — mejor que VLangCherry en ambos ejes |
| Reservadas/longest match en `.flex` | orden correcto, comentario correcto ("longest match resuelve == vs =") |

## Hallazgos

### M1 (MEDIA) — `let l: list<int> = expr;` descarta el init en silencio
La producción de declaración (`parser.cup:123`) acepta `opt_init` para cualquier `tipo_dato`, incluida `list<T>`; pero `Declaracion.ejecutar` (A.java:426) hace `case LIST: val = new Valor(tipo, new ArrayList<>())` **ignorando `init` por completo**. `let l: list<int> = [1,2];` produce una lista vacía sin ningún error.
**Fix (elegir uno):** (a) rechazarlo — error semántico "una lista dinámica se inicializa vacía; usá push()" (mínimo, consistente con el diseño actual); o (b) soportar init desde un LiteralVector. Documentar la decisión en la gramática BNF entregable.

### M2 (MEDIA-BAJA, decisión de diseño sin documentar) — aliasing inconsistente en vectores
- Declaración con init **copia**: `validarVector` construye una lista nueva (`out`).
- Asignación posterior **comparte**: `AsignacionVariable` hace `s.valor = v` — `b = a;` deja a `b` y `a` apuntando a la MISMA lista Java; mutar `b[0]` cambia `a[0]`.
Dos semánticas distintas para "poner un vector en una variable". No es incorrecta per se (el enunciado no fija semántica de referencia), pero es inconsistente y sin documentar — pregunta de defensa incómoda. **Fix:** unificar (copiar también en la asignación, o compartir también en la declaración) y documentar en el Manual Técnico. Structs: la asignación de variable completa también comparte el mapa (mismo tema).

### M3 (MEDIA-BAJA, verificar contra §5.20) — campos anidados inaccesibles
La gramática de acceso/asignación de campo es de UN nivel (`ID PUNTO ID`, parser.cup:169/335). Los structs anidados **se pueden declarar** (`campo_struct` admite tipo `ID`), y `p.domicilio` se puede leer/reemplazar completo — pero `p.domicilio.ciudad` es **error sintáctico**, tanto para leer como para asignar. Si el enunciado §5.20 exige structs anidados operables, es un gap funcional; si solo exige declararlos, conviene documentar el límite. (El deck de tabla-simbolos de la presentación muestra structs anidados como upgrade — revisar qué promete.)

### Bajos
- **B1** Orden de evaluación de argumentos: `invocar` evalúa en orden de **parámetros** (itera `f.params`), no en el orden textual de la llamada — observable si dos argumentos tienen efectos secundarios (llamadas con push/pop). Documentar o iterar `args` primero.
- **B2** `const` sobre lista dinámica se rechaza desde la ACCIÓN del parser con tipo "Sintactico" (parser.cup:125) — es un chequeo semántico; categorizarlo como tal (o moverlo a `Declaracion.ejecutar`).
- **B3** El comentario de parser.cup:258 ("el else colgante se resuelve por defecto con SHIFT") repite la atribución incorrecta del hallazgo S3 de C.0: con `bloque ::= LLAVE_IZQ ... LLAVE_DER` el conflicto **no puede existir**. Corregir el comentario junto con la diapositiva.
- **B4** `astDot` escapa `\` y `"` pero no saltos de línea en las etiquetas: un `Literal` de cadena con `\n` real (ya procesado por el lexer) rompe el `.dot`. Escapar `\n` → `\\n`.
- **B5** `round()` acepta `char` (esNumerico incluye CHAR): `round('a')` = 97. Consistente con las tablas si char cuenta como numérico — verificar §5.25; si no, restringir.

### Observación (no es bug — es el diseño exigido)
Un error semántico **aborta** la ejecución (`Entorno.errorSemantico` lanza; `Interprete` captura) — incluido un `match` cuyo caso tiene tipo incomparable. Es lo que pide el enunciado 4.3 ("poder terminar la ejecución ante un error semántico") y la diferencia deliberada con DataForge (propagación por null) y VLangCherry (acumula y sigue). Tres proyectos, tres políticas — las tres correctas respecto a su enunciado. Excelente material de defensa comparada.

## Insumos para Fase D

1. Fix M1 (una decisión + ~5 líneas) y B2/B3/B4 (mecánicos).
2. Decidir y documentar M2 (aliasing) y M3 (anidados) contra el enunciado §5.18-5.20 antes de tocar código.
3. Presentaciones: agregar a presentacion-compscript la tarjeta de **sobrecarga del `+`** (S6 de C.0, `Operaciones.suma`: string+primitivo, char+char→cadena, §6.5.3) y el dato pro-defensa del **default true de bool** (tabla 5.3) — ambos son preguntas de examen probables.
4. Los números corregidos del reporte de símbolos (S4 de C.0: 186 = 182 `n` + 4 `base`/`exp`) para etapa4/etapa6.
