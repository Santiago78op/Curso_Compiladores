# Auditoría Fase B.2 — presentacion-dataforge

**Qué se auditó:** las 14 páginas (index, etapas 0–6, 6 profundizaciones ★) contra el libro (Fase A) y las promesas del deck del Dragón (B.1).
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

> ## ✅ CERRADO — 2026-07-24
> **Los 3 hallazgos accionables están resueltos** (H4 no pedía cambios). Verificado contra el HTML actual.
>
> | Hallazgo | Cómo se cerró |
> |---|---|
> | H1 · wikilinks literales | No queda **ningún** `[[…]]` en las 62 páginas de las 6 presentaciones (`grep -rn '\[\['` vacío). Incluye los 2 de compscript que este informe adelantaba. |
> | H2 · citas menores | `gramaticas.html` ya cita **§4.6.1** para la jerarquía LL⊂LR; `etapa6.html` ya cita §4.8.3 (pánico) **y** §4.9.4 (token `error`). |
> | H3 · tabla ACTION/GOTO | La diapositiva con la matriz existe en `cap4` del deck del Dragón (commit `01b0987`), que es donde este informe proponía saldarla para los dos decks. |
>
> **Agregado además:** diapositiva nueva en `ast.html` — «Sin AST también significa *sin esas clases de bugs*», la tarjeta pro-defensa nº 2 de la Fase C, con las 4 familias de bugs que la ausencia de control de flujo vuelve imposibles y dónde sí aparecieron.

## Veredicto general: EXCELENTE — todas las promesas cumplidas

Las 5 verificaciones que el deck del Dragón delegaba aquí, **cumplidas y fieles al libro**:

| Promesa | Verificación |
|---|---|
| `automatas.html`: Thompson pieza a pieza + subconjuntos + simulación `ababb` | ✅ Fig. 3.34 fiel (estados 0–10); tabla de subconjuntos **correcta** — incluso calcula bien E={1,2,4,5,6,7,10} donde mi conversión OCR del libro traía errata; simulación del AFD (fig. 3.36) paso a paso; **bonus**: longest match con retroceso al último ✓ (cubre la brecha B2 del vault cap. 3) y minimización A≡C → 4 estados (fig. 3.65) |
| `gramaticas.html`: tabla M, shift-reduce, ambigüedad/precedence | ✅ Tabla M completa y correcta (fig. 4.17) + **traza LL con pila** (estilo fig. 4.21 — cubre la brecha B4 del vault cap. 4) + demo shift-reduce con `SUM(1,2)` + warning real de CUP + precedence con `left/right/nonassoc` |
| `tabla-simbolos.html`: ámbitos encadenados | ✅ Demo animada de la cadena (get que sube, shadowing, descarte al cerrar bloque), cita correcta a la clase `Env` fig. 2.37, comparación honesta ("cadena de un eslabón, correcta por diseño") y el upgrade a CompScript con el código del `padre` |
| `ast.html`: traza comparada directa-vs-AST | ✅ + honestidad ejemplar: "los nodos del Dragón hacen `gen()`, los nuestros `evaluar()` — esta variante está más cerca de Appel" |
| `semantica.html`: sistema de tipos formalizado | ✅ Tabla completa de reglas de tipos del proyecto + el argumento original (y correcto) "chequeo dinámico ≡ estático porque todo corre exactamente una vez" + las 3 familias de errores consolidadas con sus 3 recuperaciones |

Las etapas 0–6 son precisas, con conexión enunciado↔teoría↔código constante (ej.: la prueba de que `titulo` no puede ser reservada porque el propio enunciado la usa como variable; el typo `graphBar`/`grapBar` del enunciado documentado como decisión).

## Hallazgos

### H1. Wikilinks literales en el HTML — **BAJA (bug cosmético, corrección de 1 línea c/u)**
Sintaxis de Obsidian pegada sin adaptar; en el navegador se ven los corchetes:
- `presentacion-dataforge\etapa1.html:189` → `nota [[JFlex]] del vault`
- `presentacion-dataforge\gramaticas.html:584` → `la sección %left de [[Jison]]`
- (Para la B.3: también `presentacion-compscript\etapa2.html:45` y `etapa4.html:83`.)

### H2. Citas menores — **BAJA**
- `gramaticas.html` (comparación LL/LR): "la jerarquía del libro (§4.5)" — la afirmación LL⊂LR está en **§4.6.1**.
- `etapa6.html`: el terminal `error` se cita como §4.8.3 (aceptable: recuperación LR); la referencia más precisa del mecanismo Yacc/CUP es **§4.9.4**.

### H3. El paso "estados → tabla ACTION/GOTO" implícito — **MEDIA-BAJA (compartido con B.1)**
La demo shift-reduce cierra con "¿cómo sabe CUP cuándo shift y cuándo reduce? Con una tabla ACTION/GOTO que genera" — sin mostrar la matriz. Misma observación que H3 de B.1: **una sola diapositiva nueva** (en el cap. 4 del deck del Dragón, con la tabla del ejemplo `E→E+n|n`) salda la deuda para ambos decks; aquí bastaría enlazar.

### H4. Nada que corregir en contenido — verificaciones puntuales que pasaron
- `57.end` → NUMERO 57 + `.` como error léxico: consistente con el patrón `{Digito}+("."{Digito}+)?` y la regla `[^]`.
- Conteos de tokens de los quizzes (9 y 9): correctos.
- La afirmación "DataForge sin `precedence` porque las operaciones son funciones con paréntesis": consistente con la gramática del proyecto.
- Derivaciones, cajas y trazas: fieles a las figuras del libro que citan.

## Fortalezas a preservar

- El **"pack de parcial"**: recetas LL, tabla M con regla de llenado, traza con pila, y el test formal "celda doble = no LL(1)".
- Las **decisiones de diseño documentadas con evidencia del enunciado** (atributos de gráfica como ID, typo graphBar) — oro para la defensa.
- La **coherencia inter-deck**: cada promesa cruzada entre presentaciones se cumple; los conceptos se retoman con referencia explícita ("Demo 6 de Gramáticas").
- `semantica.html` como plantilla de "sistema de tipos formalizado" — replicable en los decks de los otros proyectos.

## Para la etapa B.3 (los 4 decks de proyectos restantes)

1. Corregir los 2 wikilinks literales de compscript ya localizados.
2. Verificar en `presentacion-compscript`: el alcance estático (etapa4 ya asoma "el padre es siempre el global — nunca el del llamador" ✓ por el grep), señales return/break, y si el deck cubre **sobrecarga** en su sistema de tipos (allí sí aplica: string+number).
3. Verificar en `presentacion-compinterpreter`: precedencia `%left` en Jison (la promesa de gramaticas.html), switch con fall-through.
4. `presentacion-conjanalyzer`: leyes de conjuntos ≠ cap. 9 (la advertencia del panorama), gramática de conjuntos.
5. `presentacion-vlangcherry`: alcance OLC2 — auditar consistencia general.
