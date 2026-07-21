# Guion de clase — ConjAnalyzer

> Esta es tu nota de trabajo para DAR la clase, no material para estudiantes. Las diapositivas HTML (`etapaN.html`) están escritas en 3ª persona para la audiencia; esto en cambio son consejos directos para vos: qué decir, qué mostrar en vivo, cómo responder si preguntan.

## Antes de empezar

**Cómo abrir la presentación:** doble clic en `C:\Users\72358\Desktop\Hades\presentacion-conjanalyzer\index.html`. Es autocontenida (sin fetch, sin servidor) — abre directo en el navegador. Desde el índice vas a `etapa0.html` y de ahí navegás con las flechas ← → o los botones de abajo; el botón ◐ cambia tema claro/oscuro si la sala tiene mala luz.

**Duración total estimada: 90–100 minutos** (7 etapas + la profundización ★ + cierre). Repartido así:

| Bloque | Tiempo |
|---|---|
| Etapa 0 | 8 min |
| Etapa 1 | 12 min |
| Etapa 2 | 12 min |
| Etapa 3 | 15 min |
| Etapa 4 | 15 min |
| Etapa 5 | 10 min |
| Etapa 6 | 10 min |
| ★ Simplificador (profundización) | 8 min |
| Cierre | 5 min |

Si vas corto de tiempo, la ★ es la primera que podés recortar o dejar como "para repasar en casa" — no rompe la continuidad porque etapa4/etapa6 ya cubren lo esencial.

**Qué preparar antes de arrancar:**

1. Tené **IntelliJ IDEA abierto** con el proyecto `ConjAnalyzer` (el que tiene el `pom.xml`) ya cargado — así no perdés tiempo de clase esperando que Maven indexe.
2. Confirmá que compila una vez antes de la clase (ventana Maven → Lifecycle → `compile`), para no descubrir un problema de JFlex/CUP en vivo.
3. Si vas a usar la terminal en lugar del botón ▶ Run de IntelliJ, el comando exacto con el Maven embebido del IDE es:

```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemploN.ca"
```

   Reemplazá `ejemploN.ca` por el archivo que toque en cada etapa (ver abajo). El `mainClass` por defecto del `exec-maven-plugin` en el `pom.xml` ya es `conjanalyzer.TestInterprete`, así que no hace falta agregar `-Dexec.mainClass`.
4. Para la demo de la GUI (Etapa 5), el comando análogo es:

```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" clean javafx:run
```

   O simplemente ▶ Run sobre la clase `Lanzador` desde IDEA (toolbar: Nuevo / Abrir / Guardar / ▶ Ejecutar / Reportes) — nunca sobre `EditorApp` directamente: el launcher de Java revisa si la clase con `main` extiende `Application` y, si JavaFX no está en el module-path, aborta con "JavaFX runtime components are missing". `Lanzador` existe justo para esquivar ese chequeo.
5. Tené a mano, en una pestaña del explorador de archivos, la carpeta `ConjAnalyzer/reportes/` — después de cada corrida de consola ahí aparecen `tokens.html`, `errores.html`, `operaciones.html` y `simplificacion.json`, y los vas a querer abrir en vivo más de una vez.

---

## Etapa 0 — Conceptos base

**Objetivo pedagógico:** que el grupo entienda el dominio (conjuntos, no aritmética) y las 3 decisiones que se arrastran por todo el proyecto, antes de ver una línea de código.

**Tiempo sugerido:** 8 minutos.

**Puntos clave a enfatizar oralmente (no releas la diapositiva):**
- Empezá comparando con DataForge: mismo esqueleto de 3 capas, dominio distinto. Esto les da un ancla si ya vieron DataForge.
- El universo ASCII 33–126 es la pieza que hace posible calcular un complemento — remarcá el porqué, no solo el rango. Sin universo concreto, "todo lo que no está en A" no se puede enumerar.
- Case sensitive real es la diferencia de diseño más visible con DataForge — anunciá que va a reaparecer en el lexer (Etapa 1) y en el `Entorno` (Etapa 3), para que no lo vean como un detalle aislado.
- La notación prefija es solo un anticipo acá — no te enredés explicando el porqué todavía, eso es la Etapa 2. Con mostrar 2-3 ejemplos de la tabla alcanza.

**Demo en vivo:** no aplica en esta etapa — es puramente conceptual.

**Preguntas probables y cómo responderlas:**
- *"¿Por qué no usar Unicode completo como universo, ya que Java lo soporta?"* — Porque el enunciado fija el rango ASCII imprimible (33–126) como universo computable; ampliarlo es válido en teoría pero no es lo que pide la sección 4.4, y complicaría la definición de "elemento" sin ganar nada para este proyecto.
- *"¿`conjA` y `conjIA` cuentan como el mismo identificador si difieren solo en mayúscula?"* — No. Es la pregunta 1 del quiz real: el `Entorno` no normaliza claves, a diferencia de DataForge. Aprovechá para remarcar que esta decisión es la sección 4.1 del enunciado, no un capricho de implementación.
- *"¿Por qué el operador va antes en vez de en medio, si a nadie le sale natural escribirlo así?"* — Adelantá que la ventaja real (nunca hace falta paréntesis) se ve recién en la Etapa 2, cuando armás la gramática; acá alcanza con que sepan leerlo de adentro hacia afuera.

**Transición a Etapa 1:** "Ya sabemos qué problema resuelve y con qué reglas de fondo. Ahora sí, ¿qué vocabulario reconoce el archivo `.ca`? Esa es la Etapa 1."

---

## Etapa 1 — Análisis léxico (JFlex)

**Objetivo pedagógico:** que vean la tabla de tokens completa y entiendan por qué el orden de las reglas en el `.flex` importa (match más largo + orden de declaración).

**Tiempo sugerido:** 12 minutos.

**Puntos clave a enfatizar oralmente:**
- Las 3 reservadas (`CONJ`, `OPERA`, `EVALUAR`) van siempre en mayúsculas exactas — conectalo con el case sensitive de la Etapa 0.
- Reservar `U`, `&`, `-`, `^` como tokens fijos (no IDs) tiene una consecuencia concreta: ningún conjunto puede llamarse literalmente `U`. Vale la pena decirlo en voz alta porque es un error común al escribir casos de prueba.
- El caso `->` vs `-` es el ejemplo de libro de "regla del match más largo" — pero aclará el matiz real: JFlex prefiere el lexema más largo, y ESO solo es posible si la regla larga (`->`) está declarada ANTES que la corta (`-`) en el archivo. No es automático por longitud del lexema, es por orden de declaración cuando dos reglas podrían calzar el mismo prefijo.
- `SIMBOLO` es el comodín que cubre todo el ASCII imprimible no reservado — es lo que permite conjuntos de signos de puntuación (`ejemplo4_nosimplificable.ca` lo prueba con `!,?,@,$,%`).
- Todo lo que cae FUERA de 33–126 (como una `ñ`) es error léxico: se descarta el carácter y el análisis sigue. Este es un buen momento para sembrar que "ningún error aborta la ejecución" — lo van a ver reforzado en cada etapa siguiente.

**Demo en vivo:**

```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo2_errores.ca"
```

Corré `ejemplo2_errores.ca` y señalá en la consola el error léxico de la `ñ` (`conjuñ` → se reporta, el nombre queda truncado a `conju`, y ese conjunto SÍ se termina definiendo). Después abrí `reportes/tokens.html` y mostrá la tabla real: lexema + nombre de token + línea/columna para cada token reconocido — conectalo con que el nombre sale por reflexión sobre `sym` (esto se profundiza recién en Etapa 6, no te desvíes ahora, solo mencionalo de pasada).

**Preguntas probables y cómo responderlas:**
- *"¿Cuántos tokens produce `CONJ : a -> 1~3;`?"* — 8 contando el `;` final: `CONJ` · `:` · `a` (ID) · `->` · `1` (NUMERO) · `~` · `3` (NUMERO) · `;`. Los espacios no cuentan.
- *"¿Qué token produce un `@` suelto?"* — `SIMBOLO`: está en el universo ASCII pero no calza con ninguna regla de reservada u operador, así que cae en el comodín y queda disponible como elemento de un conjunto.
- *"¿Por qué no usar `%ignorecase` como en DataForge y ya?"* — Porque el enunciado (4.1) exige explícitamente sensibilidad a mayúsculas para ConjAnalyzer; es una decisión de diseño documentada con un comentario propio en el `.flex`, no un olvido.

**Transición a Etapa 2:** "Ya reconocemos las piezas sueltas del vocabulario. Ahora falta decir en qué ORDEN pueden aparecer para formar una sentencia válida — eso es la gramática."

---

## Etapa 2 — Análisis sintáctico (CUP)

**Objetivo pedagógico:** que entiendan cómo la notación prefija se traduce en una gramática recursiva sin ambigüedad, y por qué acá el no terminal `operacion` sintetiza un objeto (`NodoOperacion`) en vez de solo validar forma.

**Tiempo sugerido:** 12 minutos.

**Puntos clave a enfatizar oralmente:**
- La estructura general (`{ sentencias }` con recursión izquierda) es prácticamente igual a la de DataForge — no te detengas mucho ahí, es puro reconocimiento.
- La regla central de la etapa es la de `<operacion>`: cada operador consume una cantidad FIJA de operandos (2 binarios, 1 el complemento, 0 la hoja `{ID}`). Remarcá la idea "reglas finitas, lenguaje infinito" — es la misma idea de recursión que ya vieron en DataForge, aplicada a otro dominio.
- Hacé el ejercicio de derivar `U U {A} {B} {C}` en el pizarrón o mostrando el slide — leerlo de afuera hacia adentro es lo que más cuesta la primera vez.
- La diferencia grande con DataForge: acá `operacion` tiene TIPO (`NodoOperacion`), no `void`. Es la primera pista de por qué va a hacer falta un árbol — no resuelvas el "por qué" todavía, eso es toda la Etapa 3, solo dejalo picando.
- El modo pánico (`error PUNTO_COMA`) es la misma técnica de DataForge (Dragón §4.8.3) aplicada acá.

**Demo en vivo:** reutilizá la misma corrida de `ejemplo2_errores.ca` de la Etapa 1 (si ya la corriste, no hace falta repetir el comando) pero esta vez señalá los 2 errores SINTÁCTICOS en la consola: `CONJ : -> a~z;` (falta el ID, reporta "no se esperaba '->'") y `OPERA : mala -> U {vocales};` (falta el segundo operando, reporta "no se esperaba ';'"). Abrí `reportes/errores.html` y mostrá que la sentencia `buena` que viene después de ambos errores se define y evalúa sin problema — es la prueba visual de que el modo pánico no aborta el análisis.

**Preguntas probables y cómo responderlas:**
- *"¿Cuántas `<operacion>` consume `^` y cuántas `U`?"* — 1 y 2 respectivamente. Por eso `^` es el único operador unario del lenguaje.
- *"¿Por qué `operacion` tiene tipo `NodoOperacion` acá y no `void` como en DataForge?"* — Porque una operación necesita recorrerse MÁS de una vez (evaluar, simplificar, servir al diagrama de Venn) — eso solo es posible si queda representada como estructura de datos. Se profundiza en la Etapa 3, así que si insisten podés decir "en 10 minutos lo vemos con detalle".
- *"`& U {C} {A} {B}` ¿qué operación matemática es?"* — `(C ∪ A) ∩ B`. Es el ejemplo textual del enunciado (4.6): el `&` exterior toma como primer operando todo `U {C} {A}` y como segundo operando `{B}`.

**Transición a Etapa 3:** "Ya tenemos un árbol construyéndose en cada reducción del parser. Ahora la pregunta que decide toda la arquitectura: ¿por qué hace falta guardarlo, si DataForge nunca lo necesitó?"

---

## Etapa 3 — Ejecución: Entorno, Conjunto, NodoOperacion

**Objetivo pedagógico:** que entiendan la razón arquitectónica de por qué ConjAnalyzer SÍ necesita un árbol (aunque acotado) cuando DataForge no lo necesitaba, y que vean el caso de estudio de la inconsistencia real del enunciado (4.8).

**Tiempo sugerido:** 15 minutos — es la etapa más densa conceptualmente, no la apures.

**Puntos clave a enfatizar oralmente:**
- Arrancá con la pregunta central tal cual la plantea la diapositiva: "¿cuántas veces se necesita recorrer una operación?". Una operación sirve para 3 cosas distintas (evaluarse, simplificarse, servir al Venn) — eso es lo que la distingue de un `CONJ`, que se calcula una sola vez y se guarda.
- La regla general que querés que se lleven: si algo se recorre una sola vez, ejecutarlo directo en la acción del `.cup` alcanza (como hace DataForge con todo). Si se recorre más de una vez, de formas DISTINTAS, hace falta materializarlo como árbol. Esta es la idea más transferible de toda la clase — va a reaparecer en CompScript con un AST completo.
- `NodoOperacion` tiene solo 3 formas (hoja / unario / binario) — no es un AST del programa completo, es acotado solo a `OPERA`. `CONJ` y `EVALUAR` se siguen ejecutando directo, igual que en DataForge.
- Mostrá los dos recorridos (`evaluar()` en postorden, `pertenenciaRegion()` como función booleana) como dos formas DISTINTAS de recorrer el MISMO árbol — es la prueba de por qué hacía falta guardarlo como estructura y no como resultado ya calculado.

**El caso de estudio (sección 4.8) — dale tiempo, es de lo mejor de la clase:**

El enunciado define `conjuntoA -> 1,2,3,a,b` y `conjuntoB -> a~z`, con `operacion1 = conjuntoA & conjuntoB`. Para `EVALUAR({1, b}, operacion1)`, el PDF del enunciado muestra:

```
1 -> exitoso
b -> exitoso
```

Pero la intersección real es `{a, b}`: `conjuntoA` tiene el carácter `1`, pero `conjuntoB` (rango `a~z`) NO lo tiene — `1` no pertenece a la intersección. Lo correcto (y lo que produce ConjAnalyzer) es:

```
1 -> fallo
b -> exitoso
```

Presentalo así, casi como un cliffhanger: "el enunciado tiene un error, y el código NO lo reproduce — vamos a ver por qué eso es lo correcto, no un bug". El punto pedagógico fuerte acá es: **leer el enunciado con sentido crítico**. Copiar ciegamente la salida del PDF para "coincidir" habría introducido un bug real solo por parecerse al enunciado. La decisión correcta es priorizar la semántica matemática y DOCUMENTAR la discrepancia en el código (`Entorno.evaluar`) — eso es defendible en la revisión oral con la evidencia matemática en la mano, y es exactamente el tipo de criterio que un evaluador de OLC1 valora.

**Demo en vivo:**

```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo1.ca"
```

Corré `ejemplo1.ca` (es justo el caso 4.8) y mostrá en la consola las dos líneas de `EVALUAR({1, b}, operacion1)`: `1 -> fallo` / `b -> exitoso`. Si alguien duda, hacé la cuenta en el pizarrón en vivo: `conjuntoA = {1,2,3,a,b}`, `conjuntoB = {a..z}`, intersección = `{a,b}` — `1` claramente no está.

**Preguntas probables y cómo responderlas:**
- *"Si ConjAnalyzer no tuviera que simplificar ni dibujar el Venn, ¿haría falta igual el árbol?"* — No necesariamente: con un solo recorrido (evaluar y listo) se podría ejecutar directo en las acciones del `.cup`, exactamente como DataForge. El árbol se vuelve necesario justo cuando aparece un SEGUNDO uso de la misma estructura que no puede reconstruirse desde el resultado ya evaluado.
- *"¿Por qué `evaluar()` recorre en postorden y no al revés?"* — Porque el resultado de un nodo depende del resultado de sus hijos: no se puede calcular `A ∪ B` sin tener antes los conjuntos evaluados de A y B. Postorden = resolver abajo antes de combinar arriba.
- *"¿Por qué el código no reproduce la salida tal cual la pone el PDF?"* — Porque esa salida es matemáticamente incorrecta dado cómo el propio enunciado definió los conjuntos. Priorizar "coincidir con el PDF" habría metido un bug real a propósito.

**Transición a Etapa 4:** "Ya vimos el primer uso real del árbol: evaluarlo. El segundo uso — reescribirlo aplicando las leyes de la teoría de conjuntos — es la Etapa 4, y también trae una historia real de auditoría."

---

## Etapa 4 — El Simplificador (sección 7)

**Objetivo pedagógico:** que entiendan el motor de reescritura hasta punto fijo con guardas, y que vean los 2 hallazgos reales de la auditoría de calidad como ejemplo de qué es una auditoría de código seria (no solo "correr un linter").

**Tiempo sugerido:** 15 minutos — la otra etapa densa, dale el tiempo que necesite.

**Puntos clave a enfatizar oralmente:**
- Las 13 propiedades de la sección 7 se agrupan en 7 familias (doble complemento, DeMorgan, conmutativa, asociativa, distributiva, idempotente, absorción), pero solo 5 tienen transformación REAL de reescritura (doble complemento, DeMorgan, idempotentes, absorción, distributivas). Conmutativa y asociativa NO reescriben nada por sí solas — son solo ayuda de comparación vía `canon()`. Esta distinción es la primera decisión de diseño de la etapa, remarcala antes de entrar al motor.
- El motor: postorden (a lo sumo una regla por nodo) repetido hasta punto fijo (`toPrefijo()` deja de cambiar), con tope de seguridad de 100 iteraciones. La pregunta "¿por qué repetir pasadas completas?" tiene una respuesta concreta: aplicar una ley puede exponer la oportunidad de aplicar otra en el nivel de arriba — se ve clarísimo en la demo de DeMorgan de la profundización ★.
- Las guardas son la parte más fina: DeMorgan solo se aplica si YA hay un `^` adentro (para que genere un `^^` que se cancele después); distributiva solo factoriza, nunca expande. Sin esas guardas, el motor podría "inflar" el árbol o incluso oscilar sin converger.

**Los 2 hallazgos de la auditoría — este es el otro momento fuerte de la clase, dale tiempo:**

1. **Ley distributiva faltante.** La sección 7 del enunciado lista las propiedades distributivas explícitamente, y la sección 5.4 (JSON de simplificación) pide aplicarlas. El código construido originalmente (20 de julio) implementaba 4 de las 5 leyes con transformación real y dejaba la distributiva sin implementar — ni siquiera como detección de "no simplificable". Marcá bien la diferencia con conmutativa/asociativa: esas SÍ estaban documentadas a propósito como "solo ayuda de comparación", con razón matemática válida (reordenar no reduce, podría ciclar). La distributiva en cambio SÍ tiene un sentido que reduce el árbol (factorizar, análogo a absorción) — era un gap real y especificado, no una decisión de diseño. Se encontró cruzando el grafo de llamadas del proyecto (`codebase-memory-mcp`) contra el enunciado y el Manual Técnico. El arreglo: se agregó `distributiva(char op, NodoOperacion a, NodoOperacion b)`, activado desde `paso(...)` justo después de absorción, con la misma guarda de "solo factorizar".

2. **Código muerto eliminado: `Entorno.getUniverso()`.** El mismo grafo de llamadas mostró que ese getter no tenía ningún llamador real en todo el proyecto — el universo se pasa siempre como parámetro explícito donde hace falta. Se eliminó. Mencioná también el matiz honesto: la MISMA auditoría descartó falsos positivos (`definirOperacion`, `registrarToken` "parecían" sin llamadores porque sus únicos llamadores están en código GENERADO por CUP/JFlex, fuera del índice de análisis) — es un buen ejemplo de que una auditoría de calidad real no es "borrar todo lo que el grafo marca en rojo", sino distinguir señal de ruido.

Usalo como gancho de cierre de este bloque: "esto es lo que se espera de ustedes en la revisión de su propio código antes de entregar: no solo que compile, sino poder decir con evidencia por qué algo se sacó o se agregó."

**Demo en vivo:**

```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo3_simplificacion.ca"
```

Corré `ejemplo3_simplificacion.ca` (tiene los 5 casos: demorgan, doble, idempotente, absorcion, distributiva) y abrí `reportes/simplificacion.json`. Señalá puntualmente la entrada `"distributiva"` con `"leyes": ["Propiedades distributivas"]` y `"conjunto simplificado": "& {conjA} U {conjB} {conjC}"` — es el caso agregado en la auditoría, verificado de punta a punta.

**Preguntas probables y cómo responderlas:**
- *"¿Por qué conmutativa/asociativa no son reglas de reescritura propias?"* — Porque reordenar o reagrupar un árbol no lo simplifica por sí solo, y podría incluso ciclar. Se usan como criterio de comparación (`canon()`) para que otras leyes reconozcan operandos equivalentes en cualquier orden.
- *"¿Qué distingue el hallazgo de la distributiva de la decisión de no aplicar conmutativa?"* — La distributiva SÍ tiene un sentido que reduce el árbol (factorizar), tal como pide la sección 7; omitirla era un gap real. Conmutativa/asociativa no reducen nada por sí mismas: no aplicarlas es una decisión correcta y documentada, no un bug.
- *"Si se quitara la guarda de la distributiva, ¿qué pasaría?"* — Con `& {A} U {B} {C}` se reescribiría de vuelta a `U & {A} {B} & {A} {C}` — el árbol original, más grande. El motor podría oscilar entre ambas formas sin converger nunca a un punto fijo.

**Transición a Etapa 5:** "Ya tenemos el árbol evaluado y simplificado. Falta la parte visual que pide el enunciado: el diagrama de Venn, que reutiliza uno de los recorridos que ya vimos en la Etapa 3."

---

## Etapa 5 — Diagrama de Venn (Canvas)

**Objetivo pedagógico:** que entiendan que el sombreado es EXACTO (calculado píxel a píxel con la misma función booleana del intérprete) y no un dibujo aproximado hecho aparte.

**Tiempo sugerido:** 10 minutos.

**Puntos clave a enfatizar oralmente:**
- La idea central en una frase: por cada píxel se pregunta "¿en qué círculos cae?" y se evalúa `pertenenciaRegion` (la MISMA función de la Etapa 3) sobre esa combinación de booleanos — si da `true`, se pinta. No hay lógica de dibujo separada de la lógica de evaluación.
- El límite geométrico honesto: de 1 a 3 conjuntos se dibuja con círculos; con 4 o más, el proyecto reconoce que no hay disposición geométrica razonable y degrada a solo texto. Vale la pena decir que esto es una decisión de diseño defendible, no una limitación oculta.
- El complemento sombrea también "lo de afuera de todos los círculos" — dos tonos del mismo color (`DENTRO`/`FUERA`), mismo `if`, solo cambia el tono según si el píxel cae dentro de algún círculo.
- El dato que suele sorprender: el diagrama usa el árbol ORIGINAL (`op.arbol`), no el simplificado — el panel solo MUESTRA el simplificado como texto aparte. Si el Simplificador reescribe el árbol, no hay que tocar una línea de `DiagramaVenn` para que el dibujo siga siendo correcto.

**Demo en vivo:**

```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" clean javafx:run
```

(o ▶ Run sobre `Lanzador` desde IDEA — nunca sobre `EditorApp` directamente, ver la nota en "Qué preparar antes de arrancar"). En el editor que abre, usá el botón **Abrir** para cargar `entradas/ejemplo1.ca` (o `ejemplo3_simplificacion.ca` si querés mostrar una operación con 3 conjuntos base para el trébol clásico), tocá **▶ Ejecutar** y después **Reportes** — o el panel de Venn que corresponda según cómo esté cableada la navegación en tu build — para mostrar el Canvas sombreado en vivo, con la expresión prefija, el resultado y la simplificación debajo.

**Preguntas probables y cómo responderlas:**
- *"¿Qué pasa con una operación de 5 conjuntos base?"* — Solo texto: "Venn no disponible para 5 conjuntos base". No existe disposición geométrica razonable con círculos simples para más de 3, y el proyecto lo reconoce en vez de dibujar algo engañoso.
- *"¿Por qué no hace falta recorrer el `Set<Character>` completo del universo para sombrear?"* — Porque en vez de preguntar "¿qué caracteres pertenecen al resultado?", se pregunta "¿este punto del canvas satisface la función booleana de la operación?" — `pertenenciaRegion` resuelve eso sin conocer ningún carácter concreto.
- *"Si simplifico la operación, ¿cambia el dibujo?"* — No. El diagrama usa el árbol original evaluado, no el simplificado — son independientes; el simplificado solo se muestra como texto informativo.

**Transición a Etapa 6:** "Ya tenemos evaluación, simplificación y visual. Falta el último eslabón: dejar todo eso documentado en archivos — los reportes y el JSON que pide el enunciado."

---

## Etapa 6 — Reportes HTML + JSON (Gson)

**Objetivo pedagógico:** que entiendan cómo se generan los 4 artefactos de salida (3 HTML + 1 JSON) y por qué el `Entorno` fresco por ejecución es un requisito que atraviesa todo el proyecto, no un detalle de esta etapa.

**Tiempo sugerido:** 10 minutos.

**Puntos clave a enfatizar oralmente:**
- Los 3 HTML (`tokens.html`, `errores.html`, `operaciones.html`) son autocontenidos (CSS embebido, se abren con doble clic) — mismo criterio que las presentaciones del curso.
- El truco de la reflexión sobre `sym`: el nombre de cada token se obtiene recorriendo los campos de la clase generada por CUP, en vez de un `switch` manual que se desincronizaría cada vez que alguien agregue un terminal. Mismo truco que usa DataForge — si ya dieron esa clase, mencionalo como refuerzo, no como novedad.
- El JSON tiene un formato exacto que hay que respetar al pie de la letra: para operaciones que SÍ simplifican, un objeto con `"leyes"` y `"conjunto simplificado"`; para las que NO, el STRING LITERAL `"No se puede simplificar la operacion"` — no un objeto con arreglo vacío. Y `disableHtmlEscaping()` es obligatorio porque si no, Gson escapa el `&` que la notación prefija usa todo el tiempo.
- El `Entorno` fresco por ejecución: `Interprete.ejecutar(codigo)` crea una instancia nueva cada vez — nunca acumula estado de una corrida a la siguiente. Esto conecta directo con la sección 5 del enunciado ("no deberá mostrarse reportes de análisis anteriores").

**Demo en vivo:** no hace falta una corrida nueva — reutilizá el `reportes/simplificacion.json` que ya abriste en la Etapa 4 y mostrá también `tokens.html` (de la Etapa 1) y `errores.html` (de la Etapa 2) para cerrar el círculo: "estos 4 archivos son el resultado de la MISMA corrida, del MISMO `Entorno`". Si querés un cambio de ejemplo, corré `ejemplo4_nosimplificable.ca` y mostrá en el JSON que las 3 operaciones (`union1`, `combo`, `resta`) aparecen con el string literal de "no se puede simplificar" — es la contraparte exacta de lo que vieron en Etapa 4.

**Preguntas probables y cómo responderlas:**
- *"¿Qué pasaría con el reporte de tokens si mañana agrego un token `XOR` al `.cup`?"* — Nada que arreglar en `Reportes.java`: como el nombre sale por reflexión sobre `sym`, el nuevo token aparecería con su nombre correcto sin tocar el reporte.
- *"¿Por qué las operaciones que no simplifican no tienen un objeto `{"leyes": [], ...}` vacío?"* — Porque el formato del enunciado (5.4) pide específicamente el string literal como valor directo de la clave, no un objeto con arreglo vacío.
- *"¿Qué garantiza que dos ejecuciones seguidas del editor no mezclen reportes?"* — Que `Interprete.ejecutar(...)` crea un `Entorno` nuevo en cada llamada, y que `Reportes`/`JsonSalida`/`DiagramaVenn` siempre reciben ESE `Entorno` como parámetro — nunca leen estado global acumulado.

**Transición a la profundización ★:** "Con esto el pipeline completo ya está armado y las 6 etapas quedan cerradas. Si tienen tiempo y ganas, hay una profundización que vale la pena mostrar: ver el Simplificador reescribiendo un árbol paso a paso, en vivo, con las mismas reglas que acaban de ver."

---

## ★ Profundización — El Simplificador, en movimiento

**Objetivo pedagógico:** consolidar el motor de reescritura de la Etapa 4 viéndolo aplicarse paso a paso sobre dos casos reales, con foco en distinguir "una sola pasada" de "dos pasadas hasta punto fijo".

**Tiempo sugerido:** 8 minutos (es opcional/recortable si vas corto de tiempo — avisalo así al grupo para que no sientan que se están saltando algo obligatorio).

**Puntos clave a enfatizar oralmente:**
- Es la MISMA distributiva que mostraste en el JSON de la Etapa 4 (`(A∩B) ∪ (A∩C) → A∩(B∪C)`), pero ahora paso a paso: llega al punto fijo en UNA sola pasada porque el factor común (`{A}`) ya estaba explícito de entrada.
- La segunda demo (DeMorgan + doble complemento) es la que mejor explica por qué el motor repite pasadas completas: la ley de DeMorgan, aplicada en la raíz, genera DOS `^^` nuevos en un nivel del árbol que `paso()` ya había recorrido en esa misma pasada (postorden = hijos antes que el nodo) — hace falta una pasada nueva para cancelarlos.
- Usá el botón «▶ Paso» de la demo en vivo durante la clase — está pensado para ir a tu ritmo, no es una animación automática. Dejá que el grupo prediga el siguiente paso ANTES de tocar el botón.

**Demo en vivo:** es la propia diapositiva HTML (`simplificador.html`), no hace falta terminal — navegá a esa página y usá los botones «▶ Paso» / «⟲ Reiniciar» de cada una de las 2 demos.

**Preguntas probables y cómo responderlas:**
- *"En la Demo 1, ¿por qué se factoriza con `{A}` y no con `{B}` o `{C}`?"* — Porque `distributiva(op, a, b)` prueba las 4 combinaciones posibles entre los operandos de ambas intersecciones, y `{A}` es el único que aparece en AMBAS — el único candidato a factor común.
- *"Si el árbol de la Demo 2 no tuviera el `^` interno antes de `&{A}{B}`, ¿se aplicaría DeMorgan igual?"* — Depende de si CUALQUIERA de los dos operandos de la intersección es un complemento — la guarda exige `hijo.izq.esUnario() || hijo.der.esUnario()`. Si ningún operando lo fuera, DeMorgan no dispararía, porque no generaría ninguna cancelación posterior.
- *"¿Cuántas pasadas como máximo permite el motor antes de rendirse?"* — 100 (`MAX_ITER`), una guarda de seguridad contra un bucle infinito si alguna combinación de reglas oscilara sin converger — algo que no debería pasar con las guardas actuales, pero se protege igual.

**Transición al cierre:** "Con esto termina el recorrido completo del proyecto. Cerremos con la síntesis y hacia dónde va lo que viene."

---

## Cierre

**Síntesis final (2-3 minutos, decila casi textual, es el resumen que se van a llevar):**

"ConjAnalyzer toma el mismo esqueleto de 3 capas de DataForge — léxico, sintáctico, ejecución — y lo aplica a un dominio distinto: conjuntos en vez de aritmética. La diferencia arquitectónica más importante es que acá SÍ hace falta un árbol, aunque acotado solo a las operaciones (`NodoOperacion`), porque una operación se recorre de tres formas distintas: se evalúa, se simplifica con las leyes de la sección 7, y sirve de función booleana para el diagrama de Venn. Vimos también dos lecciones que van más allá del código: que hay que leer el enunciado con sentido crítico — la sección 4.8 tiene un error real y el proyecto no lo copia ciegamente — y que una auditoría de calidad real no es pasar un linter, es completar un gap real (la distributiva faltante) y distinguir código muerto real de falsos positivos (el código generado por CUP/JFlex que el análisis no ve)."

**Conexión con el próximo proyecto (CompScript):**

Cerrá señalando hacia adelante: "ConjAnalyzer necesitó un árbol acotado — solo para las operaciones, porque solo ellas se recorren más de una vez. El próximo proyecto, CompScript, va un paso más allá: TODO el programa se representa como AST, porque va a tener control de flujo real (`if`, `while`, `for`), funciones y recursión — cosas que DataForge y ConjAnalyzer nunca necesitaron, porque ninguno de los dos tiene una instrucción que dependa del resultado de otra ejecutada condicionalmente. La pregunta que usamos hoy — '¿cuántas veces hace falta recorrer esto, y de cuántas formas distintas?' — es exactamente la misma pregunta que va a justificar el AST completo de CompScript."

Si tenés tiempo, invitá a abrir `../presentacion-compscript/index.html` como adelanto visual antes de cerrar la sesión.

---

*Guion elaborado a partir del contenido real de `index.html`, `etapa0.html`–`etapa6.html`, `simplificador.html` y `README.md` de esta carpeta, verificado contra `ConjAnalyzer/entradas/*.ca` y `ConjAnalyzer/pom.xml`.*
