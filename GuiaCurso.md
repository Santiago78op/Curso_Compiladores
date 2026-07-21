# Guía maestra del curso — Hades

> Este documento es tuyo, para vos que das la clase. No repite lo que ya está escrito en cada `GuionClase.md` / `GuiaProgramacion.md` / `README.md` — es la capa de **navegación y secuencia** por encima de todo eso: qué sesión toca, qué abrís en cada momento, y cuándo cambia la dinámica de clase. Cuando dice "leé acá", andá al documento real — esta guía no duplica su contenido.

---

## 1. Panorama general

Hades sostiene **dos cursos distintos**:

- **OLC1** (Organización de Lenguajes y Compiladores 1): 4 proyectos de intérprete (DataForge, ConjAnalyzer, CompScript, CompInterpreter) + una presentación teórica completa del Libro del Dragón (caps. 1–7 + panorama 8–12).
- **OLC2** (Organización de Lenguajes y Compiladores 2): **VLangCherry**, un único proyecto grupal (3 integrantes) con defensa oral — es una materia y una dinámica completamente distintas, no lo trates como "el quinto proyecto de la secuencia de OLC1".

### Mapa de material por proyecto

| Proyecto | Curso | Presentación (diapositivas) | GuionClase (tu chuleta) | GuiaProgramacion (temario) | Modo de clase |
|---|---|---|---|---|---|
| Libro del Dragón (teoría) | OLC1 | [`presentacion-libro-dragon/index.html`](presentacion-libro-dragon/index.html) | *(no existe — ver nota abajo)* | — | Exposición teórica |
| DataForge | OLC1 | [`presentacion-dataforge/index.html`](presentacion-dataforge/index.html) | [`presentacion-dataforge/GuionClase.md`](presentacion-dataforge/GuionClase.md) | [`DataForge/docs/GuiaProgramacion.md`](DataForge/docs/GuiaProgramacion.md) | **Demostración guiada** (vos programás en vivo) |
| ConjAnalyzer | OLC1 | [`presentacion-conjanalyzer/index.html`](presentacion-conjanalyzer/index.html) | [`presentacion-conjanalyzer/GuionClase.md`](presentacion-conjanalyzer/GuionClase.md) | [`ConjAnalyzer/docs/GuiaProgramacion.md`](ConjAnalyzer/docs/GuiaProgramacion.md) | **Tutor invertido** (el estudiante programa, vos revisás) |
| CompScript | OLC1 | [`presentacion-compscript/index.html`](presentacion-compscript/index.html) | [`presentacion-compscript/GuionClase.md`](presentacion-compscript/GuionClase.md) | [`CompScript/docs/GuiaProgramacion.md`](CompScript/docs/GuiaProgramacion.md) | **Tutor invertido** |
| CompInterpreter | OLC1 | [`presentacion-compinterpreter/index.html`](presentacion-compinterpreter/index.html) | [`presentacion-compinterpreter/GuionClase.md`](presentacion-compinterpreter/GuionClase.md) | [`CompInterpreter/docs/GuiaProgramacion.md`](CompInterpreter/docs/GuiaProgramacion.md) | **Tutor invertido** |
| VLangCherry | **OLC2** (aparte) | [`presentacion-vlangcherry/index.html`](presentacion-vlangcherry/index.html) | [`presentacion-vlangcherry/GuionClase.md`](presentacion-vlangcherry/GuionClase.md) | [`VLangCherry/docs/GuiaProgramacion.md`](VLangCherry/docs/GuiaProgramacion.md) | Grupal, con **defensa oral** en vivo |

**Nota sobre el Libro del Dragón**: a diferencia de los 5 proyectos, `presentacion-libro-dragon/` no tiene un `GuionClase.md` propio — es una presentación de teoría pura (8 páginas: `cap1.html`…`cap7.html` + `panorama.html`), sin demos de código para correr. Los tiempos que le asigno en el plan de sesiones más abajo son **estimaciones mías** (no hay cifra documentada por capítulo, a diferencia de los proyectos donde el tiempo sí sale literal de cada `GuionClase.md`) — ajustalos con margen la primera vez que la des.

### El cambio de dinámica más importante del curso

DataForge la construyó Claude como **demostración guiada**: vos programás en vivo mientras explicás, el grupo mira y pregunta. A partir de **ConjAnalyzer**, los roles se invierten — **el estudiante escribe el código, vos tutoreás y revisás**. Esto está documentado en `CLAUDE.md` (regla de trabajo) y reforzado en el cierre de cada `GuionClase.md`:

- Cierre de [`presentacion-dataforge/GuionClase.md`](presentacion-dataforge/GuionClase.md) (sección "Qué conecta con el próximo proyecto"): anticipá el cambio ANTES de que termine la sesión de DataForge, no al empezar ConjAnalyzer — que no esperen la misma dinámica.
- Apertura de [`presentacion-conjanalyzer/GuionClase.md`](presentacion-conjanalyzer/GuionClase.md): arranca comparando con DataForge para darles un ancla, pero la mecánica de la clase ya es otra.

Qué decirle al grupo en ese momento exacto (versión corta para el aula): *"Hasta acá yo armé el código en vivo y ustedes vieron cómo se construye un intérprete de punta a punta. Desde ConjAnalyzer, el código lo escriben ustedes — yo reviso, hago preguntas, y solo destrabo cuando hace falta. Usen la `GuiaProgramacion.md` de cada proyecto como su tutorial paso a paso; yo la sigo en paralelo para saber dónde están y qué preguntarles."* Qué esperar del grupo: van a ir más lento que en DataForge (es normal, están escribiendo, no mirando), van a trabar en los mismos puntos que ya señalan los "Errores comunes" / "Dónde se rompe si..." de cada `GuiaProgramacion.md` — usalos como tu lista de síntomas para diagnosticar rápido en vivo.

---

## 2. Plan de sesiones — OLC1

Diez sesiones, intercalando los capítulos del Dragón con las etapas de cada proyecto siguiendo las conexiones que ya documentó la tabla "Contenido" de [`presentacion-libro-dragon/README.md`](presentacion-libro-dragon/README.md) (cap3 ↔ `automatas.html`, cap4 ↔ `gramaticas.html`, cap5 ↔ el atributo-AST DataForge-vs-CompScript, etc.). Los tiempos de cada etapa de proyecto son los que ya calculó cada `GuionClase.md`; los de los capítulos del Dragón son estimados (ver nota de la sección 1).

### Sesión 1 — Apertura teórica (sin código todavía)
- **Dragón cap. 1** — Introducción: compilador vs. intérprete, fases, tabla de símbolos, estático vs. dinámico. Abrí con esto: es el mapa mental que necesitan antes de ver una línea de DataForge. (~25-30 min, estimado)
- **Dragón cap. 2** — Un traductor sencillo: BNF, derivaciones, acciones semánticas y postfijo, descendente vs. LALR. Es el "primer compilador de juguete" del libro — preparación directa para ver DataForge Etapa 1-2. (~30-35 min, estimado)
- Sin demo en vivo — es la única sesión 100% teórica del curso, apoyate solo en `presentacion-libro-dragon/index.html`.

### Sesión 2 — DataForge completo (demostración guiada)
- [`presentacion-dataforge/GuionClase.md`](presentacion-dataforge/GuionClase.md), Etapas 0 a 6 (roadmap principal completo) — **135 min (~2h15)**, tal como el propio guion recomienda darlo de una sola vez, sin partirlo.
- Vos programás en vivo; seguí el guion al pie de la letra (tiene los comandos exactos, qué archivo `.df` correr en cada etapa, y las preguntas típicas del quiz).
- Al cerrar, decí explícitamente el mensaje de cambio de dinámica (sección 1 de esta guía) — el guion ya trae el texto sugerido en su sección de cierre.

### Sesión 3 — Léxico y sintáctico a fondo
- **Dragón cap. 3** — Análisis léxico: rol del lexer, ER, definiciones regulares, pipeline ER→AFN→AFD. (~30 min, estimado)
- [`presentacion-dataforge/GuionClase.md`](presentacion-dataforge/GuionClase.md), profundización `automatas.html` — 25 min. Es la pareja documentada de cap. 3 (el propio README de libro-dragon dice que cap3 "enlaza a las demos de `automatas.html`"): dala inmediatamente después del capítulo, sobre el mismo hilo (el token `NUMERO` de DataForge).
- **Dragón cap. 4** — Análisis sintáctico: recuperación de errores, *dangling else*, LL, FIRST/FOLLOW, autómata LR(0) con ítems, escalera SLR→LALR. (~35 min, estimado)
- Profundización `gramaticas.html` — 35 min. Pareja de cap. 4 (README: "complementa `gramaticas.html`") — mismo criterio, dala justo después.
- Total aproximado: ~2h05.

### Sesión 4 — Semántica, AST y por qué DataForge no lo necesita
- **Dragón cap. 5** — Traducción dirigida por la sintaxis: SDD, sintetizados/heredados, la pila `RESULT` de CUP, S/L-atribuidas, **el atributo-AST comparando DataForge vs. CompScript** — esta última parte es el puente real hacia el próximo proyecto, citalo textual.
- Profundización `ast.html` de DataForge (15 min) — por qué DataForge no construye AST; es literalmente el mismo argumento que acaba de dar cap. 5.
- **Dragón cap. 6** — Código intermedio con lente de intérprete: chequeo de tipos con demo de coerción, corto circuito, *backpatching* (por qué no aplica acá), switch. El corto circuito es un anticipo útil: es exactamente el bug real que van a auditar en CompScript Etapa 5.
- Profundización `semantica.html` (15 min) + `fases.html` (10 min) de DataForge para cerrar el bloque.
- Total aproximado: ~1h45-2h.

### Sesión 5 — ConjAnalyzer completo — CAMBIO DE DINÁMICA
- [`presentacion-conjanalyzer/GuionClase.md`](presentacion-conjanalyzer/GuionClase.md) completo: Etapas 0–6 + profundización ★ `simplificador.html` + cierre — **90-100 min**.
- Esta es la **primera sesión en modo tutor invertido**: decí el mensaje de cambio de dinámica ANTES de arrancar (ver sección 1) y tené a mano `ConjAnalyzer/docs/GuiaProgramacion.md` — es lo que el estudiante va siguiendo mientras escribe.
- Mencioná de pasada (no hace falta profundizar todavía) que el caso de estudio de la sección 4.8 del enunciado —el enunciado tiene un error real que el código no reproduce— es un buen anticipo de la actitud crítica que se espera en la revisión de código de todos los proyectos que siguen.

### Sesión 6 — CompScript, primera mitad
- **Dragón cap. 7** — Entornos en tiempo de ejecución: zonas de memoria, árbol de activación, registro ≈ `Entorno`, demo de pila con `fact(3)`, alcance estático vs. dinámico. (~30-35 min, estimado)
- Repaso rápido de `tabla-simbolos.html` de DataForge (10 min) — recordá el caso de un solo ámbito antes de mostrar la cadena de ámbitos real de CompScript.
- [`presentacion-compscript/GuionClase.md`](presentacion-compscript/GuionClase.md), Etapas 0 a 4 (por qué hace falta AST, léxico, sintáctico, el AST real con el bug de `validarVector`, tabla de símbolos con ámbitos anidados) — ~62 min documentados.
- Total aproximado: ~1h35-1h45. Seguís en modo tutor invertido.

### Sesión 7 — CompScript, segunda mitad
- [`presentacion-compscript/GuionClase.md`](presentacion-compscript/GuionClase.md), Etapas 5 a 7 (flujo de control con el bug real de cortocircuito en `&&`/`||`, funciones y recursión, editor + 5 reportes) + profundización ★ obligatoria `ast.html` — ~67 min documentados, no la recortes (el propio guion la marca como obligatoria).
- Es el cierre natural de cap. 5/6 del Dragón que diste en la Sesión 4 — el bug de cortocircuito es la aplicación real de "por qué corto circuito importa" de cap. 6.

### Sesión 8 — CompInterpreter, primera mitad
- [`presentacion-compinterpreter/GuionClase.md`](presentacion-compinterpreter/GuionClase.md), Etapas 0 a 4 (cambio de stack a JS/Jison, léxico, parser + fábrica de nodos del AST, intérprete de dos pasadas) — no hace falta repetir teoría de autómatas o AST, ya la vieron dos veces; andá directo a lo específico de Jison.
- Recordá el mensaje de la Etapa 0 de este guion: la teoría NO cambia, solo la herramienta (Jison en vez de JFlex+CUP) — es la continuación directa del argumento de cap. 2/3/4 del Dragón, ya aplicado dos veces antes.

### Sesión 9 — CompInterpreter, segunda mitad
- [`presentacion-compinterpreter/GuionClase.md`](presentacion-compinterpreter/GuionClase.md), Etapa 4 (señales de control + los 2 bugs reales de auditoría, si no llegaste en la Sesión 8) + Etapas 5 a 7 (servidor Express, cliente React, suite Playwright) + profundización ★ `arquitectura-rest.html` (arquitectura cliente-servidor, la pieza que distingue a este proyecto de los otros 3).
- Cerrá mencionando el puente documentado en el propio guion hacia VLangCherry (OLC2): "el próximo nivel es generación de código real, fuera del alcance de OLC1" — sin entrar en detalle si esta sesión no continúa directo a OLC2.

### Sesión 10 — Cierre teórico de OLC1
- **Dragón panorama.html (caps. 8–12)** — la frontera front-end/back-end (con la frase para la defensa), generación de código, optimización (**y por qué ConjAnalyzer NO es el cap. 9** — aclaralo explícitamente, es una confusión real que el propio material anticipa), paralelismo, Apéndice A, y el **mapa final capítulo → proyecto**.
- Esta es la sesión de síntesis: recién ahora tiene sentido mostrar la tabla completa capítulo↔proyecto, porque el grupo ya vio los 4 proyectos y puede reconocer cada pieza.
- Cierre de todo OLC1: repasá en una frase cada proyecto (DataForge sin AST, ConjAnalyzer con árbol acotado, CompScript con AST completo, CompInterpreter con arquitectura cliente-servidor) — es literalmente el resumen que ya arma el cierre de [`presentacion-compinterpreter/GuionClase.md`](presentacion-compinterpreter/GuionClase.md), reusalo.

**Total: 10 sesiones para OLC1** (asumiendo bloques de ~2 horas; las Sesiones 2, 5 y 7 corren un poco más largas por la densidad de sus proyectos, dejá margen de 15-20 min extra en esas tres).

---

## 3. Plan de sesiones — OLC2 (VLangCherry)

Curso y dinámica aparte: **trabajo grupal de 3 integrantes con defensa oral en vivo** (documentado en [`presentacion-vlangcherry/README.md`](presentacion-vlangcherry/README.md) y desarrollado a fondo en [`presentacion-vlangcherry/GuionClase.md`](presentacion-vlangcherry/GuionClase.md)). No repito ese contenido acá — solo la secuencia.

El propio guion ya define dos versiones de tiempo (completa ~110 min para estudio, defensa acotada ~35-40 min): usalas como estén documentadas, no las recalcules.

### Sesión OLC2-1 — VLangCherry, primera mitad
- Etapas 0 a 4: por qué ANTLR4 + Go (requisito del enunciado, no elección), la gramática ANTLR4, del *parse tree* al AST propio (decisión de `type-switch` en vez de Visitor generado), tipos estáticos/structs/métodos "por referencia" gratis vía punteros de Go, slices 1D/2D.

### Sesión OLC2-2 — VLangCherry, segunda mitad
- Etapas 5 a 7: control de flujo + recursión + **los 4 hallazgos de la auditoría semántica (2026-07-21)** — es la etapa que el propio guion marca como la más importante para la defensa, no la recortes — servidor REST con `net/http`, cliente React reusando componentes de CompInterpreter.
- Profundización ★ `generadores.html`: ANTLR4 vs. JFlex+CUP vs. Jison — buen cierre comparativo de TODO el curso (los 3 generadores usados en los 5 proyectos), aunque es la primera candidata a recortar si el tiempo aprieta.

### Ensayo de defensa (sesión aparte, cerca de la fecha real)
- Repasá con el equipo la sección **"Para la defensa oral"** del propio [`GuionClase.md`](presentacion-vlangcherry/GuionClase.md#para-la-defensa-oral--modificaciones-que-el-equipo-debería-poder-hacer-en-vivo-con-confianza): los 5 movimientos que el tribunal probablemente les va a pedir en vivo (función nativa nueva, validación semántica nueva, explicar por qué structs son "por referencia" gratis, agregar un operador completo a la gramática, explicar el intérprete de dos pasadas).
- Usá la "versión defensa" (~35-40 min: Etapa 0 → Etapa 3 → Etapa 5 → Etapa 6 → cierre) como el recorrido de práctica cronometrado, tal cual lo sugiere el guion.

**Total: 2 sesiones + 1 ensayo cronometrado para OLC2.**

---

## 4. Cómo usar cada documento en el momento

No mires los tres documentos a la vez sin saber cuál te toca — esta es la regla rápida:

| Momento de la clase | Documento que abrís | Por qué |
|---|---|---|
| Antes de la clase, preparando qué vas a decir | `GuionClase.md` del proyecto de hoy | Tiene el guion completo: tiempos, puntos clave a decir en voz alta, preguntas probables con respuesta lista, y el comando EXACTO de terminal para cada demo |
| Proyectando en el salón, para que la audiencia vea | `index.html` de la `presentacion-X/` correspondiente (doble clic, `file://`, sin servidor) | Es lo único que ve el estudiante — texto en 3ª persona, sin tus notas de director |
| El estudiante ya está escribiendo código (desde ConjAnalyzer en adelante) | `docs/GuiaProgramacion.md` del proyecto en `ConjAnalyzer/`, `CompScript/`, `CompInterpreter/` o `VLangCherry/` | Es el tutorial paso a paso que el estudiante sigue; vos lo tenés abierto en paralelo para saber en qué sección está y qué preguntarle. Las secciones "Errores comunes" / "Dónde se rompe si..." son tu lista de síntomas para diagnosticar rápido |
| Necesitás confirmar un dato técnico del proyecto (no de la clase) | `docs/ManualTecnico.md` del proyecto (si existe) o el código real en `DataForge/`, `ConjAnalyzer/`, etc. | Es la fuente de verdad técnica — nunca inventes un dato que no verificaste ahí |
| Querés repasar la secuencia completa del curso, o qué sigue después de hoy | Este documento (`GuiaCurso.md`) | Es el único que amarra el orden entre proyectos y capítulos del Dragón |

Regla de oro: **el `GuionClase.md` es tuyo, la `index.html` es de ellos, la `GuiaProgramacion.md` es de quien escribe código.** Si te confundís cuál abrir, preguntate primero "¿quién está mirando la pantalla en este momento — yo, el grupo, o el estudiante que programa?".

---

## 5. Checklist de preparación antes de arrancar el curso completo

No repito los comandos exactos (cada `GuionClase.md` los trae completos) — solo qué tener listo, resumido:

### Para OLC1 (DataForge, ConjAnalyzer, CompScript, CompInterpreter)

- [ ] **IntelliJ IDEA** instalado y con el proyecto correspondiente ya abierto e indexado (Maven resuelto) antes de que empiece cada sesión de DataForge/ConjAnalyzer/CompScript — evita el "downloading dependencies" en vivo.
- [ ] **JDK 25 gestionado por IDEA** disponible en `C:\Users\72358\.jdks\openjdk-25.0.1` — **nunca uses el `java` del PATH (es 1.8)**. Si vas a usar terminal en vez del botón ▶ Run de IDEA, el comando necesita `JAVA_HOME` apuntando a ese JDK y el Maven embebido del propio IDEA (`.../plugins/maven/lib/maven3/bin/mvn`) — el patrón exacto está en cada `GuionClase.md`.
- [ ] Compilar una vez cada proyecto ANTES de la clase (Maven → Lifecycle → `compile`) para no descubrir un problema de JFlex/CUP en vivo.
- [ ] Para las demos con GUI JavaFX (DataForge, ConjAnalyzer, CompScript): correr siempre `Lanzador`, **nunca `EditorApp` directamente** — el error "JavaFX runtime components are missing" es un caso real de estudio en la Etapa 4 de DataForge, así que si te pasa en vivo por accidente, aprovechalo.
- [ ] Para CompInterpreter: `Node.js`/`npm` instalados, y **dos terminales** listas — una para `server` (`npm start`, puerto 4000) y otra para `client` (`npm run dev`, Vite puerto 5173) — dejalas corriendo toda la clase. `npx playwright test` para la Etapa 7 (reutiliza los servidores si ya están arriba).
- [ ] Tener a mano en el explorador de archivos las carpetas de `entradas/` y `reportes/` (o equivalente) de cada proyecto — vas a abrir los `.html`/`.json` generados varias veces en vivo.
- [ ] Repasar `presentacion-libro-dragon/index.html` antes de la Sesión 1 y 10 — no tiene código para correr, pero sí conviene tener listo el pizarrón/proyector para dibujar diagramas (fases, árbol de activación, etc.) mencionados en varias secciones de esta guía.
- [ ] Probar el botón ◐ (tema claro/oscuro) de cada presentación antes de la clase, según la luz del salón/proyector — es el mismo motor (`assets/estilo.css` + `assets/deck.js`) en las 6 presentaciones.

### Para OLC2 (VLangCherry)

- [ ] **Go** instalado y funcionando; confirmar que `go run ./cmd/servidor` levanta en `localhost:4100` (`http://localhost:4100/salud` debe responder `{"estado":"ok"}`).
- [ ] **Node/npm** para el cliente React (`npm run dev`, Vite).
- [ ] Dos terminales corriendo toda la sesión (servidor Go + cliente React), igual que CompInterpreter.
- [ ] Alternativa liviana sin cliente gráfico: `go run ./cmd/cli entradas/ejemploN.vch` — útil para las etapas tempranas donde todavía no hace falta mostrar la interfaz.
- [ ] Si vas a practicar las modificaciones en vivo de la sección "Para la defensa oral": tener abiertos de antemano `internal/runtime/nativas.go`, `internal/runtime/interprete.go` y `grammar/VLangCherry.g4`, y el `.jar` de ANTLR4 (`tools/antlr.jar`) accesible para regenerar el parser si hace falta.
- [ ] Correr `npm run e2e` (suite Playwright) una vez antes de la defensa real, para confirmar que el pipeline completo sigue en verde.

### Publicación (opcional, para las 6 presentaciones)

- [ ] Si vas a activar GitHub Pages para alguna presentación: Settings → Pages → Deploy from branch → `main` / root (mismo procedimiento en las 6, documentado en cada `README.md`).
