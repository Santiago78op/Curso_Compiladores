# Fase C.4 — Code review de ConjAnalyzer (Java + JFlex + CUP, árbol acotado)

**Qué se revisó:** ~1,600 renglones leídos completos — `Simplificador` (verificado ley por ley), `NodoOperacion`, `Entorno`, `parser.cup`, `Lexer.flex`, `DiagramaVenn`, `JsonSalida`, `Operacion`/`Conjunto`/`Interprete` — y spot-check de `Reportes`/`EditorApp`/`TestInterprete` (escapado HTML, entorno fresco, patrón DataForge ya verificado).
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto: el más limpio de los 5 proyectos

**Cero hallazgos altos y cero medios.** El corazón del proyecto — el Simplificador — resiste la verificación manual completa:

- **Doble complemento** `^^X → X` ✓; **DeMorgan guardado** (solo si un operando ya es `^`, generando el `^^` que la siguiente pasada cancela) ✓; **idempotentes** vía `equivalentes()` ✓; **absorción** con sus variantes conmutadas ✓; **distributiva** solo en el sentido que factoriza, con las 4 combinaciones de factor común ✓.
- `canon()` es correcto en el detalle difícil: aplana y ordena **solo** las cadenas de `U`/`&` (conmutativos/asociativos) y trata `-` aparte como no conmutativo (`-(izq,der)` posicional) — el error clásico habría sido aplanar también la diferencia. La deduplicación dentro de `canon` está matemáticamente justificada (idempotencia en la forma canónica).
- El **orden de reglas por nodo** (idempotencia → absorción → distributiva) evita solapamientos: el caso `(X&Y) U (X&Y)` lo captura idempotencia antes de que distributiva lo toque.
- Punto fijo con `MAX_ITER=100` ✓; el Simplificador trabaja sobre `copia()` ✓; y como `NodoOperacion` es **inmutable** (campos `final`), el sharing de subárboles que hace `paso()` (devolver `hijo.izq`, `a`, `b` sin copiar) es seguro por construcción — un detalle de diseño fino que merece mención en defensa.
- El reporte "honesto" de conmutativa/asociativa (solo cuando la comparación no trivial las usó) implementado con heurísticas razonables (`registrarCommAssoc`/`registrarDistributiva`/`registrarAbsorcion`).

El resto acompaña: caso 4.8 del enunciado (la salida matemáticamente incorrecta) documentado en el javadoc de `Entorno.evaluar` con la decisión tomada; validación de referencias ANTES de evaluar; `charUnico` valida longitud Y universo; Venn píxel-exacto reutilizando `pertenenciaRegion` con el complemento sombreando "afuera" en tono claro; JSON con el string literal exacto y `disableHtmlEscaping` justificado; `esc()` aplicado en todas las celdas de los reportes; entorno fresco por ejecución en GUI y consola.

## Hallazgos (todos BAJOS)

- **B1** — `definirConjuntoLista` **aborta el conjunto completo** ante un elemento inválido (`return` en el loop, Entorno.java:83), mientras que `evaluar()` usa `continue` por elemento. Dos criterios distintos para el mismo tipo de falla. Ambos defendibles; elegir uno y documentarlo (sugerencia: definir el conjunto con los elementos válidos + error por cada inválido, como hace EVALUAR).
- **B2** — Un `CONJ` y una `OPERA` **pueden compartir nombre** (mapas separados, sin chequeo cruzado). Probablemente válido (espacios de nombres distintos por diseño del lenguaje: `{X}` siempre refiere a conjuntos, `EVALUAR(..., X)` siempre a operaciones) — documentarlo en `docs/gramatica.txt` como decisión.
- **B3** — El comentario de `Lexer.flex:55-57` justifica el orden `"->"` antes de `"-"` con la regla de desempate ("a igual longitud gana la primera") cuando en realidad ahí decide el *longest match* — como correctamente dice el propio archivo 20 líneas después (:78). Misma confusión que los hallazgos C2/I1 de B.3: corregir en el mismo lote (2 presentaciones + este comentario).
- **B4** — `syntax_error` imprime a `System.err` además de registrar en el entorno (también en CompScript) — ruido inofensivo; opcional silenciarlo.

## Observaciones (sin acción)

- `NodoOperacion.evaluar` defensivo contra referencias inexistentes (null) aunque `definirOperacion` ya las valida antes — redundancia sana.
- `EVALUAR` con datos duplicados imprime una línea por aparición — fiel al formato por-dato de 4.8.
- El bucle del Venn reutiliza el mismo `HashMap` de región por píxel (claves sobrescritas) — correcto y eficiente.

## Insumos para Fase D

Nada urgente de código. B1/B2 son decisiones a documentar (gramatica.txt + Manual); B3 entra en el lote transversal "longest match vs orden" (con conjanalyzer/etapa1.html y compinterpreter/etapa1.html); B4 opcional. ConjAnalyzer queda como **referencia de calidad** del grupo — y su Simplificador verificado a mano es un activo de defensa: las dos demos de la presentación (profundización ★) coinciden exactamente con el código.
