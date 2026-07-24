# Fase C — Resumen consolidado del code review de los 5 proyectos

**Qué se hizo:** lectura completa del código fuente propio de los 5 proyectos (~10,000 renglones: VLangCherry 2,760 · CompScript 2,600 · CompInterpreter 1,880 · ConjAnalyzer 1,600 · DataForge 1,370), verificando semántica contra los enunciados reales donde hubo duda. Informes por proyecto: `fase-c-verificaciones.md` (C.0) + `fase-c-{vlangcherry,compscript,compinterpreter,conjanalyzer,dataforge}.md`.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## El mapa de un vistazo

| Proyecto | Altos | Medios | Bajos | Estado |
|---|---|---|---|---|
| VLangCherry (C.1) | **5** 🐞 | 4 | 4 | El que más necesita Fase D |
| CompInterpreter (C.3) | **1** 🐞 | 3 | 5 | Sólido con un bug transversal |
| CompScript (C.2) | 0 | 1 | 5 | Muy sólido; decisiones a documentar |
| ConjAnalyzer (C.4) | 0 | 0 | 4 | Referencia de calidad |
| DataForge (C.5) | 0 | 0 | 4 | Limpio; su arquitectura elimina clases de bug |

## Los 3 hallazgos transversales (una decisión, varios proyectos)

1. **Cortocircuito de `&&`/`||`** — el mismo bug en VLangCherry y CompInterpreter; CompScript lo encontró en su auditoría y lo corrigió con test. Fix ~8 líneas por proyecto + caso de prueba (`x != 0 && 10/x > 1`). **Regla para Fase D: todo fix confirmado en un proyecto se verifica en los otros cuatro.**
2. **`return` mal validado** — `return;` en función tipada devuelve 0 (VLangCherry) y `return valor` en void se descarta en silencio (VLangCherry + CompInterpreter). CompScript valida ambos: usar su `invocar` como referencia.
3. **"Longest match vs orden de reglas"** — la confusión aparece en 2 presentaciones (B.3: C2, I1) y 2 comentarios de código (ConjAnalyzer Lexer.flex:55, y el del dangling-else en CompScript parser.cup:258 que es tema aparte pero mismo lote). La formulación correcta: *longest match decide entre longitudes distintas; el orden de declaración desempata igual longitud (reservadas vs Id)*. DataForge y CompScript ya lo dicen bien.

## Cola de fixes para Fase D, por prioridad

### VLangCherry (los 5 altos — ~60 líneas en total + pruebas)
1. **A5** `ValorPorDefecto` de slice/struct pierde el tipo declarado → toda asignación posterior rechazada. Fix: `Valor{Tipo: t}` con puntero nil (la forma que `EsNil` ya anticipa). Junto con **M3-nativas** (append/join panic sobre slice nil).
2. **A2** cortocircuito (transversal #1).
3. **A3** sin guardas: `for true {}` cuelga la petición; recursión infinita = **stack overflow fatal de Go que `recover()` NO atrapa → muere el proceso**. Fix: MAX_DEPTH (~2000) en `invocarFuncion` + MAX_ITER (~1M) en `ejecutarFor` (copiar el patrón de CompInterpreter).
4. **A1** `mut` decorativo: la inmutabilidad jamás se valida (el propio ManualUsuario la documenta). Fix transversal: AST + traductor + Entorno + 3 validaciones.
5. **A4** receptor por valor comparte el `StructVal` (muta al llamador); `ReceptorPuntero` es código muerto. Fix: clonar en `invocarFuncion` si `!ReceptorPuntero`.
6. Medios: M1/M2 (return, transversal #2), M4 (tabla de símbolos ni-historia-ni-final), B1 (errores silenciosos en defaults de switch).
7. Nuevo `entradas/ejemplo7_semantica.vch` que ejercite TODO (hoy ningún ejemplo toca estos caminos).

### CompInterpreter
1. **A1** cortocircuito (transversal #1).
2. **M1** índices sin validar tipo (NaN pasa los chequeos de rango) — 4 puntos.
3. **M2** método void como expresión → null silencioso; **M3** `return valor` en método (transversal #2).
4. **M4** nativas sobre vector `new` (elementos null) → NaN silencioso.
5. Bajos: coerción char asimétrica, default de char vs comentario, documentar políticas.

### CompScript
1. **M1** `let l: list<int> = expr;` descarta el init en silencio → decidir (error o soporte) y documentar.
2. Documentar M2 (aliasing declaración-copia vs asignación-comparte) y M3 (campos anidados de un nivel) contra el enunciado §5.18-5.20.
3. Bajos mecánicos: comentario dangling-else, `.dot` sin escapar `\n`, `const list` como "Sintactico", conteo B.3.

### ConjAnalyzer y DataForge (solo bajos)
- ConjAnalyzer: unificar criterio elemento-inválido (definir vs evaluar), documentar CONJ/OPERA homónimos, comentario del flex (transversal #3).
- DataForge: `Cadena` sin `\r\n` (comilla sin cerrar se traga líneas), documentar tipo-de-atributo descartado y EXEC sin validar tipo, mensaje de lista vacía.

## Hallazgos pro-defensa que dejó la Fase C (para presentaciones/vault en Fase D)

> **Estado 2026-07-24 — sembrados en las presentaciones:** ①→ diapositiva nueva en `libro-dragon/cap6.html`; ②→ diapositiva nueva en `dataforge/ast.html`; ③→ callout en `compscript/etapa6.html` (verificado en `A.defecto()` y `tipos.js`); ④→ diapositiva §6.5.3 en `cap6` + callout en `compscript/etapa6.html`; ⑤→ diapositiva nueva en `conjanalyzer/simplificador.html`; ⑥→ callout en `libro-dragon/cap7.html`. Falta solo ⑦ y el volcado al vault.

1. **"Tres políticas de error semántico, las tres correctas"**: DataForge propaga null y sigue (su enunciado quiere reportes completos); CompScript aborta (su enunciado 4.3 lo exige); CompInterpreter/VLangCherry acumulan y siguen. Comparativa lista para una tarjeta.
2. **"Sin AST también significa sin esas clases de bugs"** (DataForge): sin lógicos→sin cortocircuito, sin índices→sin NaN, sin ciclos→sin guardas. La decisión arquitectónica como decisión de *superficie de error*.
3. **El default de `bool` es `true` por la tabla 5.3 del enunciado** (CompScript y CompInterpreter lo implementan bien) — trampa de examen.
4. **`+` sobrecargado** (CompScript `Operaciones.suma`; CompInterpreter tabla T_SUMA): §6.5.3 implementado y sin nombrar — la tarjeta pendiente de B.1-H2.
5. **NodoOperacion inmutable → sharing seguro en el Simplificador** (ConjAnalyzer): diseño fino verificado.
6. **Java atrapa `StackOverflowError`, Go no puede** (CompScript vs VLangCherry): comparativa de runtimes real, medida en este mismo workspace.
7. Los proyectos nuevos repitieron bugs que el viejo ya había corregido → "los fixes se propagan, no se quedan donde se encontraron".

## Estado del plan tras la Fase C

- **Fase A** ✅ libro (12 brechas ALTAS, backlog por capítulo). **Fase B** ✅ 7 presentaciones (~20 hallazgos puntuales; decks por delante del vault). **Fase C** ✅ 5 proyectos (6 altos, 8 medios, 22 bajos; 3 transversales).
- **Fase D (Opus 4.8 + agentes)** tiene ahora 3 frentes con spec ejecutable:
  1. **Vault** — un agente por capítulo con su `cap0N-brechas.md` (sembrando desde los HTML de las presentaciones).
  2. **Presentaciones** — los hallazgos de `fase-b-*.md` (citas, wikilinks, conteos, longest-match, S3/S4/V1 ya verificados en C).
  3. **Código** — esta cola de fixes, en el orden de arriba, con la regla del barrido cruzado y verificación de Fable contra estos informes.
