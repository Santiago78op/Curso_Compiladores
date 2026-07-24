# Auditoría Fase B.1 — presentacion-libro-dragon

**Qué se auditó:** las 9 páginas (`index`, `cap1`–`cap7`, `panorama`) contra el libro real leído en Fase A, en 3 ejes: precisión teórica, conexión teoría↔proyectos, didáctica.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

> ## ✅ CERRADO — 2026-07-24
> **Los 5 hallazgos de este informe están resueltos y verificados contra el HTML actual.** No queda nada por hacer acá; lo que sigue se conserva como registro de la auditoría.
>
> | Hallazgo | Cómo se cerró |
> |---|---|
> | H1 · 3 citas menores | `cap5` ya cita **§5.2.3–5.2.4** (definiciones) y **§5.4.3** (acciones internas); `cap7` ya cita **fig. 7.1**. Verificado línea por línea. |
> | H2 · sobrecarga ausente | Diapositiva nueva en `cap6` (§6.5.3): «el `+` que suma *o* concatena», con la resolución por firma y la advertencia de no confundirla con sobrecarga de *funciones*. Commit `01b0987`. |
> | H3 · tabla ACCION/ir_A | Diapositiva nueva en `cap4` con la matriz del ejemplo `E→E+n\|n`, consistente con los estados I0–I4 del autómata que el deck ya animaba. Commit `01b0987`. |
> | H4 · paso de parámetros | Diapositiva nueva en `cap7` (§1.6.6–7): valor/referencia/nombre, el matiz de Java y el aliasing. Commit `01b0987`. |
> | H5 · redacción del token `error` | El callout de `cap4` ya dice «la ③ implementada con ①». |
>
> **Agregado además (no estaba en este informe):** en `cap4`, las dos citas del token `error` que decían §4.8.3 ahora distinguen §4.8.3 (modo pánico) de **§4.9.4** (mecanismo Yacc/CUP). Y dos tarjetas pro-defensa de la Fase C: «Tres políticas de error semántico» (`cap6`) y «Java atrapa `StackOverflowError`, Go no puede» (`cap7`).

## Veredicto general: EXCELENTE — va por delante del vault

La presentación cubre **la mayoría de las brechas ALTAS que la Fase A detectó en el vault**. Verificaciones pendientes de Fase A, resueltas:

| Brecha del vault (Fase A) | ¿Está en la presentación? |
|---|---|
| Ejemplo guía fig. 1.7 animado | ✅ cap1: demo "el viaje de una instrucción" (7 pasos, fiel al libro) |
| Precedencia por niveles (n+1 no terminales) | ✅ cap2: tabla + "cuanto más abajo vive, más fuerte liga" + contraste con `precedence` de CUP |
| Reglas de cálculo de FIRST/FOLLOW | ✅ cap4: las 6 reglas + resultados del ej. 4.30 |
| Token `error` de CUP | ✅ cap4: mecánica completa (`instruccion ::= error PUNTO_COMA`) |
| Motor LR (ítems, estados) | ✅ cap4: demo del autómata LR(0) de `E→E+n\|n` con CERRADURA y REDUCE (ver "faltantes" abajo) |
| Escalera SLR→LR(1)→LALR con porqués | ✅ cap4: tabla comparativa correcta |
| Pila de valores de `RESULT` (§5.4.2) | ✅ cap5: mecánica exacta de la reducción en 3 pasos |
| Grafo de dependencias + ciclos | ✅ cap5 |
| Marcadores M→ε y acciones a mitad de regla | ✅ cap5, con la advertencia práctica de conflictos |
| GDA / subexpresiones comunes | ✅ cap6, con el ejemplo del libro |
| Coerción (demo `inttofloat`) + widening/narrowing | ✅ cap6 (cierra el círculo con cap1 — gran decisión didáctica) |
| Backpatching y por qué no aplica a intérpretes | ✅ cap6 |
| Enlace de acceso ≠ control + bug del alcance dinámico | ✅ cap7: tabla frame↔Entorno con "ojo: NO es el llamador" + demo del test `imprime 1, no 99` |
| Árbol de activación + pila con recursión | ✅ cap7: demo `fact(3)` completa |
| break/continue/return como señales | ✅ cap7: checklist + quiz 3 |
| GC como cultura + conexión JVM | ✅ cap7 |
| Frontera front/back-end + frase de defensa | ✅ panorama, incluida la advertencia "ConjAnalyzer NO es el cap. 9" |

**Consecuencia para la Fase D:** varias notas del vault pueden **sembrarse desde estas páginas** (adaptando la prosa que ya existe y agregando la cita del libro) en vez de redactarse desde cero. Los agentes de Opus deben recibir ambas fuentes: `dragon-md\` y estos HTML.

## Hallazgos (lo que SÍ hay que corregir/agregar)

### H1. Citas menores incorrectas — **BAJA (corrección puntual)**
- `cap5` diapositiva "S-atribuidas y L-atribuidas": cita *"Dragón §5.3 (definiciones)"* — las definiciones están en **§5.2.3–5.2.4** (§5.3 es "Aplicaciones").
- `cap5` diapositiva de SDT/marcadores: cita *"§5.4.4 (SDT con acciones internas)"* — las acciones dentro de producciones son **§5.4.3** (§5.4.4 es eliminación de recursividad izquierda).
- `cap7` diapositiva de zonas de memoria: cita *"fig. 7.2"* — la subdivisión de memoria es la **fig. 7.1** (la 7.2 es el programa quicksort).

### H2. Sobrecarga de operadores ausente — **MEDIA (el mismo hueco del vault)**
`cap6` cubre coerción y cast pero **nunca usa el término "sobrecarga"** (§6.5.3) ni el caso emblema: `+` = suma entre números / concatenación entre cadenas — que es exactamente lo que implementa `Operaciones` y lo que preguntan en defensa. Encaja natural: una tarjeta o mini-sección en la diapositiva del chequeo de tipos de cap6, con la regla de resolución por firma.

### H3. El paso "de los estados a la tabla SLR" queda implícito — **MEDIA-BAJA**
La demo LR(0) de `cap4` es correcta y dice "los 5 estados SON la tabla", pero no muestra la **tabla ACCION/ir_A** como matriz (filas=estados, columnas=terminales, celdas s5/r2/acc) ni una traza con **pila de estados**. Para el ejercicio de examen "construya la tabla SLR" falta ese último paso visual. Opción mínima: una diapositiva con la tabla del propio ejemplo `E→E+n|n` (5 estados × 3 símbolos). El resto del andamiaje ya existe.

### H4. Paso de parámetros (§1.6.6) sin cubrir — **MEDIA-BAJA**
Ni cap1 ni cap7 tocan valor vs. referencia (y el matiz "Java pasa referencias por valor"). Es pregunta de defensa para CompScript/CompInterpreter ("¿cómo pasa parámetros tu lenguaje?"). Encaja en cap7 junto a la tabla frame↔Entorno (los "parámetros reales" ya están en esa tabla — falta el *cómo* llegan).

### H5. Redacción confusa puntual — **BAJA**
`cap4`, callout del token `error`: *"es la ③ hecha con ②①"* — el mecanismo es producciones de error (③) implementadas con modo pánico (①); incluir ② (nivel de frase) confunde la taxonomía que la propia diapositiva acaba de dar. Sugerencia: "la ③ implementada con ①".

## Fortalezas a preservar (no tocar)

- **Los quizzes** son de calidad excepcional — en particular: la cadena voraz `\".*\"` (cap3), el corto circuito como semántica observable con el test que "destruye" la implementación ingenua (cap6), y el bug 99-vs-1 del alcance dinámico (cap7).
- **Los checklists por capítulo** (léxico, sintáctico, semántico, runtime) convierten la teoría en decisiones de proyecto — únicos en su tipo.
- El patrón "cerrar círculos" entre capítulos (el `inttofloat` del cap1 explicado en cap6; la pregunta del cap. 2 sobre recursión izquierda respondida en cap4).
- La delegación explícita a las demos de `presentacion-dataforge` (automatas.html, gramaticas.html) sin duplicar — verificar en B.2 que esas demos cumplen lo prometido: pipeline Thompson→subconjuntos→AFD con simulación, tabla M de LL, shift/reduce con `SUM(1,2)`, cura con `precedence`.

## Para la etapa B.2 (presentacion-dataforge)

Promesas hechas por este deck que hay que verificar allá: (1) `automatas.html` construye Thompson pieza a pieza + subconjuntos + simulación con `ababb`; (2) `gramaticas.html` tiene demo de tabla M, demo shift/reduce y demo de ambigüedad/precedence; (3) `tabla-simbolos.html` demuestra ámbitos encadenados de bloques; (4) `ast.html` tiene la traza comparada DataForge-vs-CompScript; (5) `semantica.html` formaliza el sistema de tipos del proyecto.
