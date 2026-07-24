# Auditoría Fase B.3 — los 4 decks de proyectos

**Qué se auditó:** las 39 páginas de `presentacion-conjanalyzer` (9), `presentacion-compscript` (10), `presentacion-compinterpreter` (10) y `presentacion-vlangcherry` (10), leídas completas, en los 3 ejes de la Fase B (precisión teórica, conexión teoría↔código, didáctica). Los `GuionClase.md` no son páginas de presentación y quedaron fuera (igual criterio que B.2).
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

> ## ✅ CERRADO — 2026-07-24
> **Los 12 hallazgos de los 4 decks están resueltos.** Los que la auditoría marcó «verificar en Fase C» (S3, S4, S6, V1) se cerraron leyendo el código, no adivinando.
>
> | | Cómo se cerró |
> |---|---|
> | **C1** conteo `CONJ : a -> 1~3;` | La respuesta dice **8**, sin ambigüedad, y aclara que el `;` sí es token. |
> | **C2** longest match vs. orden | El callout ya trae la formulación correcta: `->` gana a `-` **sin importar el orden** por ser de distinta longitud; el orden solo desempata igual longitud. |
> | **C3** snippet sin `Set r` | La declaración `Set<Character> r = new LinkedHashSet<>(a)` está en el fragmento. |
> | **C4** token `error` como §4.8.3 | `etapa2.html` ya cita **§4.9.4**. |
> | **S1** wikilinks de compscript | No queda ninguno. |
> | **S2** conteo `console.log(...)` | Dice **7**. |
> | **S3** dangling else | La diapositiva ya no atribuye el cero-conflictos al orden de producciones. |
> | **S4** 177 vs 186 llamadas | Los números de `etapa4`/`etapa6` ya están conciliados. |
> | **S5** «6 nodos» vs 7 | Dice **7**, en la tabla y en la narración de la demo. |
> | **S6** sobrecarga sin nombrar | Cerrado en dos lugares: la diapositiva §6.5.3 del `cap6` del Dragón, y un callout nuevo en `compscript/etapa6.html` que separa sobrecarga de *funciones* (que no hay) de sobrecarga de *operadores* (que sí, en `Operaciones.suma`, verificado en el código: `int+int→int`, `char+char→string`, primitivo+`string`→concatena). |
> | **I1** longest match en Jison | `etapa1.html:102` ya explica que con la opción `flex` gana `==` sin importar el orden, y que el orden solo decide empates. |
> | **V1** contradicción del receptor por valor | Reescrito: `invocarFuncion` clona con `ClonarPorValor`, así que la respuesta es **25** sin contradicción, con nota del hallazgo A4 de la auditoría. |
> | **V2** Visitor atribuido al Dragón | Ya dice «el patrón *Visitor* clásico (de GoF, doble despacho), emparentado con los recorridos y SDD del Dragón (§2.3/§5.4)». |
>
> **Agregado además:** diapositivas nuevas con tarjetas pro-defensa de la Fase C — `conjanalyzer/simplificador.html` («El nodo inmutable es lo que hace barato compartir subárboles», verificado: los 4 campos de `NodoOperacion` son `public final`) y el callout del default de `bool` = `true` en `compscript/etapa6.html` (verificado en `A.defecto()` y en `tipos.js → valorPorDefecto()`).

## Veredicto general: EXCELENTE en los 4 — con ~10 hallazgos puntuales

Los cuatro decks mantienen el nivel de B.1/B.2: código real citado con archivo y fuente, decisiones de diseño defendidas con evidencia del enunciado, bugs de auditoría contados con honestidad (síntoma → causa → corrección → verificación), quizzes con respuestas razonadas, y coherencia cruzada entre decks (cada promesa "ver el proyecto hermano" se cumple del otro lado).

Verificaciones que B.2 dejó marcadas, todas resueltas:

| Verificación pendiente | Resultado |
|---|---|
| conjanalyzer: ¿el deck evita confundir las leyes de conjuntos con optimización (cap. 9)? | ✅ Nunca lo afirma; el Simplificador se presenta como reescritura sobre el AST, sin invocar el cap. 9 — consistente con la advertencia del panorama |
| conjanalyzer: gramática de conjuntos | ✅ `simplificador.html`: las DOS demos verificadas a mano — distributiva `(A∩B)∪(A∩C)→A∩(B∪C)` ✓ y DeMorgan+doble complemento en 2 pasadas `^(^(A∩B)∩^C) → (A∩B)∪C` ✓, con las guardas correctas |
| compscript: alcance estático | ✅ `etapa4.html` es LA página del bug del alcance dinámico accidental (`ctx.global.crearHijo` vs `caller.crearHijo`) — cubre de lleno la brecha B1 del cap. 7 de Fase A |
| compscript: señales return/break | ✅ `etapa5.html`: excepciones de control sin stacktrace (`Senales.Break`), con el caso "break fuera de ciclo → error semántico" |
| compinterpreter: precedencia en Jison | ✅ `etapa2.html`: los 9 niveles `%left/%right/%nonassoc` mapeados a la tabla 5.10, MÁS la gramática BNF en capas (quiz 2) — que de paso cubre la "precedencia por niveles" §2.2.6 (brecha 3 de Fase A) |
| compinterpreter: switch con fall-through | ✅ `etapa4.html`: implementado con el mecanismo de señales (`coincidio` + `BREAK`), contrastado explícitamente con el `match` sin fall-through de CompScript y el `switch` sin fall-through de VLangCherry — tres semánticas de switch comparadas entre decks |
| vlangcherry: consistencia general | ✅ ANTLR4/ALL(*) bien explicado, etiquetas `#nombre` → type-switch, 4 fixes de auditoría con cita de enunciado, `generadores.html` cierra la comparación de los 3 generadores con precisión |

## Hallazgos por deck

### presentacion-conjanalyzer

- **C1 (BAJA)** `etapa1.html:138` — quiz "¿Cuántos tokens produce `CONJ : a -> 1~3;`?": la respuesta dice "**7:** … y el `;` final hacen 8 en total si se cuenta el punto y coma". El `;` **es** un token (PUNTO_COMA) — la respuesta correcta es 8, sin ambigüedad. Reescribir.
- **C2 (BAJA-MEDIA, precisión)** `etapa1.html:83` — el callout mezcla *longest match* con orden de declaración: "JFlex siempre prefiere el lexema más largo… **por eso `\"->\"` se escribe ANTES que `\"-\"`**". Con longest match, `->` gana a `-` sin importar el orden; el orden solo desempata matches de **igual longitud** (reservadas vs `{Id}` — caso que el deck de compscript sí explica bien). La respuesta del quiz 3 (:148) es correcta ("declararla es lo que le da la opción"); solo sobra el "ANTES".
- **C3 (BAJA)** `etapa3.html:65-80` — el snippet de `evaluar()` usa `r` en la rama binaria sin mostrar su declaración (el recorte omite `Set r = new LinkedHashSet<>(a)`). Una línea más y el fragmento queda auto-contenido.
- **C4 (BAJA)** `etapa2.html:133` — token `error` citado como "Dragón §4.8.3"; la referencia más precisa del mecanismo Yacc/CUP es §4.9.4 (mismo nit que B.1/B.2 — corregir en lote).

### presentacion-compscript

- **S1 (BAJA, bug cosmético)** wikilinks literales en el HTML: `etapa2.html:45` (`[[Ambigüedad, precedencia y asociatividad]]`) y `etapa4.html:83` (`[[Registro de activación y pila de control]]`) — ya localizados en B.2; corrección de 1 línea cada uno.
- **S2 (BAJA)** `etapa1.html:159` — quiz: "8 tokens: `console` `.` `log` `(` `CADENA` `)` `;`" — son **7**. Corregir el número (o listar el octavo si existía).
- **S3 (MEDIA-BAJA, verificar en Fase C)** `etapa2.html:90-95` — la diapositiva del *dangling else* atribuye la ausencia de conflicto al "orden de las producciones" de `if_stmt`. Si `<bloque>` exige llaves (como sugiere `A.If` en etapa5, con `cuerpo`/`sino` como listas de bloque), el else colgante **nunca puede ocurrir** — el ejemplo `if (a) if (b) x; else y;` ni siquiera parsea. La causa real del cero-conflictos sería la delimitación por llaves, no el orden. Verificar contra `parser.cup` y reescribir la explicación con la causa correcta (§4.3.2 sí aplica si los bloques fueran opcionales).
- **S4 (MEDIA-BAJA, verificar en Fase C)** `etapa4.html:109,131` + `etapa6.html:159` — "fib(10) hace **177** llamadas" y "el reporte muestra **186** entradas de `n`", presentados como que "coinciden". 177 ≠ 186 (la diferencia podría ser el `n` de `factorial(5)` + otras llamadas del mismo archivo, pero tal como está escrito no cuadra). Correr `ejemplo5.cs` y poner los números reales conciliados.
- **S5 (BAJA)** `ast.html:64` — el paso 7 de la demo dice "un árbol de **6 nodos**", pero la propia demo construye **7 objetos** (2 AccesoVariable + Literal + Binaria(−) + Argumento + Llamada + Binaria(×)). Aclarar si `Argumento` no cuenta como nodo del AST, o corregir a 7.
- **S6 (MEDIA, oportunidad — el mismo hueco H2 de B.1)** si `Operaciones.aritmetica` de CompScript resuelve `+` como suma numérica **y** concatenación de cadenas, eso ES **sobrecarga de operadores** (Dragón §6.5.3) y el deck nunca la nombra (el "sin sobrecarga" de etapa6 se refiere a funciones, concepto distinto). Verificar en Fase C y, si aplica, una tarjeta en etapa5 o en el deck del Dragón cap. 6 salda la deuda para todos.

### presentacion-compinterpreter

- **I1 (MEDIA-BAJA, precisión)** `etapa1.html:102` — la nota "si `\"=\"` estuviera escrito antes que `\"==\"`, Jison devolvería dos tokens ASIGNA" contradice la propia diapositiva 2, que muestra `%options case-insensitive flex`: con la opción `flex`, Jison aplica longest match y `==` gana sin importar el orden (el orden importa solo sin esa opción, o a igual longitud — el caso reservadas/ID de :70, que sí está bien). Misma confusión que C2; una corrección conjunta con la regla exacta ("longest match decide entre longitudes distintas; el orden decide empates") sirve para ambos decks.
- Sin más hallazgos: es posiblemente el deck de proyecto más redondo — señales vs. excepciones (contraste deliberado con CompScript), la coerción `double→int` justificada con el Anexo 11.1, los 2 bugs de auditoría, y la profundización REST.

### presentacion-vlangcherry

- **V1 (MEDIA, contradicción interna)** `etapa3.html:149` — quiz 1: la respuesta afirma "**25, sin cambiar**" y dos líneas después "un receptor por valor **SÍ comparte** el `StructVal` subyacente si no se reasigna todo el campo". Ambas no pueden ser ciertas para `p.Edad = p.Edad + 1`: si el receptor por valor comparte el puntero, la mutación se vería (26); si el intérprete copia el `StructVal` para receptores por valor, entonces "no comparte". Verificar en Fase C qué hace realmente `invocar` con receptores por valor y reescribir la respuesta sin la contradicción.
- **V2 (BAJA)** `etapa2.html:61` — "el patrón Visitor clásico **del Dragón**": el Visitor con doble despacho es de GoF (y Appel lo usa para árboles); el Dragón describe recorridos y definiciones dirigidas por sintaxis, no ese patrón con ese nombre. Ajustar la atribución ("el patrón Visitor clásico (GoF), emparentado con los recorridos del Dragón §2.3/§5.4").

## Patrones transversales (una corrección, varios decks)

1. **Longest match vs. orden de reglas** (C2 + I1): la formulación correcta — "el *longest match* decide entre candidatos de distinta longitud; el orden de declaración decide **empates** de igual longitud" — corrige ambos decks y de paso refuerza la diapositiva (correcta) de compscript.
2. **Conteo de tokens en quizzes** (C1 + S2): recontar los 2 quizzes; son la clase de detalle que un auxiliar pregunta en defensa.
3. **§4.9.4 para el token `error`** (C4 + B.1 H1 + B.2 H2): corrección en lote en 3 decks.
4. **Números de verificación real**: S4 (177/186) y S5 (6/7 nodos) piden re-verificar contra el código antes de corregir — no adivinar.

## Fortalezas a preservar (no tocar)

- **El contraste inter-proyecto como recurso didáctico**: 3 semánticas de switch comparadas (match sin fall-through / switch con fall-through / switch con break implícito), 2 mecanismos de control (excepciones en CompScript vs. señales en CompInterpreter), 3 generadores en `generadores.html`. Nadie más en el curso tiene esto.
- **Las auditorías como caso de estudio**: distributiva faltante (conjanalyzer), cortocircuito (compscript), vector con null + break silencioso (compinterpreter), 4 huecos de validación (vlangcherry) — todas con el patrón síntoma→causa→corrección→verificación y la lección "los ejemplos pasan ≠ el manejo de errores es correcto".
- **Decisiones frente a enunciados imperfectos, documentadas**: el 4.8 matemáticamente incorrecto de ConjAnalyzer, el `%` que da DECIMAL de CompInterpreter, las 3 desambiguaciones (`func`/`fn`, comillas, `in`/`range`) de VLangCherry. Oro para defensa.
- La cadena "entorno fresco por ejecución" contada 5 veces con la razón específica de cada proyecto (reportes → concurrencia REST).

## Insumos para las fases C y D

- **Fase C (code review) — verificar:** S3 (`parser.cup` de CompScript: ¿`<bloque>` exige llaves?), S4 (correr `ejemplo5.cs` y contar), S6 (`Operaciones.aritmetica`: ¿`+` concatena?), V1 (receptor por valor en `invocar` de VLangCherry: ¿copia profunda?). C3/S5 se corrigen leyendo el código correspondiente.
- **Fase D (correcciones):** los 4 wikilinks literales (2 dataforge + 2 compscript), las citas §4.9.4 en lote, los 2 recontados de tokens, la regla longest-match/orden en 2 decks, V2, y las reescrituras S3/S4/V1 **después** de que Fase C confirme los hechos.
