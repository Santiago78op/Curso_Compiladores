# Guion de clase / defensa oral — VLangCherry

> Esta es tu chuleta para presentar y defender VLangCherry, no una diapositiva más. Está en 2ª persona porque es para vos (o para cualquiera del equipo que la use el día de la defensa): consejos directos, qué decir, qué mostrar en vivo y qué responder si preguntan. Las diapositivas HTML (`etapa0.html`…`etapa7.html`, `generadores.html`) son el material neutro que vas a proyectar; este documento es la voz en tu cabeza mientras las proyectás.

Fuente de todo lo que sigue: las 9 páginas de `presentacion-vlangcherry/` (índice + etapas 0-7 + profundización ★), el `README.md` del repo de la presentación, `VLangCherry/docs/ManualTecnico.md`, los 6 `.vch` de `VLangCherry/entradas/` y `VLangCherry/client/e2e/vlangcherry.spec.js`. Si algo no está citado ahí, no lo inventes en la defensa — decí "no lo tengo memorizado, lo reviso en el manual" antes que arriesgar un dato falso frente al tribunal.

---

## Antes de empezar

**Cómo abrís la presentación:** doble clic en `presentacion-vlangcherry\index.html`. Es autocontenida (CSS y JS propios, sin fetch a internet), así que funciona sin conexión y sin servidor — perfecta para el salón de la defensa aunque falle el wifi. Dentro de cada etapa navegás con las flechas del teclado o los botones ← Anterior / Siguiente →; el botón ◐ cambia el tema claro/oscuro si el proyector se ve mal con uno de los dos.

**Duración total estimada: 100-110 minutos** si das las 8 etapas completas + la profundización ★, a ritmo de clase con preguntas. Si es defensa oral con tiempo acotado, recortá así:
- Versión completa (~110 min): todas las etapas + generadores, para una sesión de estudio del equipo antes de la defensa.
- Versión defensa (~35-40 min): Etapa 0 (rápida) → Etapa 3 (referencia) → Etapa 5 (control de flujo + auditoría, el plato fuerte) → Etapa 6 (servidor) → cierre con la sección "Para la defensa oral" de este guion. Etapas 1, 2, 4 y 7 quedan como respaldo si preguntan puntualmente por gramática, AST, slices o cliente.

**Qué preparar ANTES de entrar a la sala** (los comandos son los reales del proyecto, verificados en `ManualTecnico.md` §13 y `client/playwright.config.js`):

1. Terminal 1 — servidor Go, desde `VLangCherry/server`:
   ```
   go run ./cmd/servidor
   ```
   Levanta en `localhost:4100` (variable de entorno `PORT`, default `"4100"` si no se define — está hardcodeado en `cmd/servidor/main.go`). Verificá que responde antes de arrancar: `http://localhost:4100/salud` debe devolver `{"estado":"ok"}`.

2. Terminal 2 — cliente React, desde `VLangCherry/client`:
   ```
   npm run dev
   ```
   Levanta con Vite (puerto por defecto de Vite, típicamente 5173). El cliente apunta a `http://localhost:4100` por defecto (`VITE_API_URL` en `src/api.js`), así que con el servidor ya arriba no hace falta configurar nada más.

3. Si querés demostrar el pipeline SIN cliente gráfico (más rápido, menos que pueda fallar en vivo), desde `VLangCherry/server`:
   ```
   go run ./cmd/cli entradas/ejemplo1_basico.vch
   ```
   Imprime `=== ERRORES ===`, `=== CONSOLA ===`, `=== SIMBOLOS ===` y `=== AST ===` directo en la terminal. Tenés los 6 ejemplos listos en `entradas/`: `ejemplo1_basico.vch` (tipos y operadores), `ejemplo2_structs.vch` (structs y métodos), `ejemplo3_slices.vch` (slices 1D/2D), `ejemplo4_control.vch` (if/switch/for), `ejemplo5_funciones.vch` (recursión), `ejemplo6_errores.vch` (los 3 tipos de error a propósito).

4. Si te van a pedir modificar código en vivo (ver sección de defensa), tené abierto en el editor, de antemano: `internal/runtime/nativas.go`, `internal/runtime/interprete.go` y `grammar/VLangCherry.g4` — son los tres archivos donde vas a tocar código con mayor probabilidad.

5. Opcional pero recomendable: corré `npm run e2e` una vez antes de la defensa (desde `client/`) para confirmar que la suite Playwright sigue en verde — te da la tranquilidad de que el pipeline completo (servidor + cliente) funciona de punta a punta justo antes de presentar.

---

## Etapa 0 — Conceptos base: por qué ANTLR4 + Go

**Objetivo pedagógico:** que quede claro que Go y ANTLR4 no son gustos personales — Go es requisito del enunciado (2.3) y ANTLR4 es, de los tres generadores del curso, el único que emite código Go.

**Tiempo sugerido:** 6-8 min. Es la etapa más corta y no tiene demo — arrancá fuerte con esto para no perder tiempo temprano.

**Puntos clave para enfatizar oralmente:**
- Decí explícito: "el enunciado no nos dejó elegir el lenguaje de implementación — dice literalmente que toda la lógica del compilador debe estar en Go". Eso corta de raíz cualquier pregunta de "¿por qué no usaron JFlex+CUP como en OLC1?".
- El pipeline tiene una fase que DataForge (OLC1) no tiene: el **traductor** (parse tree → AST propio). Marcá que esto es porque VLangCherry SÍ tiene control de flujo y recursión — DataForge no, así que a DataForge le alcanzaba con ejecutar directo en las acciones del `.cup`.
- Mencioná las 4 decisiones macro (AST propio, traductor sin Visitor, referencia vía punteros, dos pasadas) como el "mapa" de lo que viene — sirve de ancla mental para el resto de la charla.

**Demo en vivo:** ninguna. Es la única etapa puramente conceptual.

**Preguntas probables (del quiz real + dudas típicas):**
- *"¿Por qué JFlex+CUP o Jison no son opciones válidas acá?"* → Porque ninguno de los dos genera código Go (generan Java y JavaScript respectivamente), y el enunciado exige Go.
- *"¿Qué diferencia de fondo hay entre el algoritmo de ANTLR4 y el de JFlex+CUP/Jison?"* → ANTLR4 genera un parser LL adaptativo (ALL(*)): decide la alternativa en tiempo de análisis. Los otros dos generan LALR/SLR/LR con tablas precomputadas en la generación.
- Duda típica que no está en el quiz pero cae seguido: *"¿Por qué no usaron el mismo stack que los otros 4 proyectos del curso, para ser consistentes?"* → Porque VLangCherry es de OLC2 (curso distinto), y el enunciado de OLC2 impone Go — no es inconsistencia, es un requisito distinto de un curso distinto.

**Transición:** "Ya sabemos por qué Go y por qué ANTLR4. Ahora veamos cómo se ve esa gramática ANTLR4 en la práctica."

---

## Etapa 1 — La gramática ANTLR4

**Objetivo pedagógico:** mostrar que léxico y sintaxis conviven en un solo `.g4`, y que las etiquetas `#nombre` son la pieza que habilita todo el diseño posterior (traductor sin Visitor).

**Tiempo sugerido:** 10-12 min.

**Puntos clave para enfatizar oralmente:**
- La regla `expr` es LA estrella de esta etapa: el orden de las alternativas ES la precedencia (multiplicativa antes que aditiva, etc., según la tabla 4.6 del enunciado). Remarcá que ANTLR4 resuelve la recursión izquierda por vos — en un parser LL manual eso obligaría a factorizar en una cascada `expr → term → factor`.
- Las etiquetas `#exprMultiplicativa`, `#exprLlamada`, etc. generan un **tipo de contexto Go concreto por alternativa**. Esta frase la vas a repetir en la Etapa 2 — es el hilo conductor de casi toda la arquitectura.
- Las 3 desambiguaciones documentadas (`func` no `fn`, comillas dobles sin interpolación, `in` no `range`) son munición perfecta para la defensa: muestran que leyeron el enunciado con ojo crítico y tomaron decisiones, no que improvisaron.

**Demo en vivo (opcional si hay tiempo):** abrí `grammar/VLangCherry.g4` real y mostrá en vivo la regla `expr` con sus etiquetas, y las reglas léxicas de `CADENA`/`RUNE`/`fragment ESC`. No hace falta ejecutar nada — es lectura de código guiada.

**Preguntas probables:**
- *"¿Qué pasa si agregás una alternativa nueva a `expr` sin etiquetarla?"* → ANTLR4 no genera un tipo de contexto propio; el `type-switch` del traductor no la puede distinguir de las demás.
- *"¿Por qué `exprMultiplicativa` va antes que `exprAditiva`?"* → El orden codifica precedencia: `*`/`/`/`%` deben resolverse antes que `+`/`-`.
- *"¿Cuántos archivos hay que tocar a mano para agregar un operador nuevo?"* → Uno: `grammar/VLangCherry.g4`. Los 4 archivos de `internal/parser/` se regeneran completos.

**Transición:** "La gramática ya nos da un parse tree completo. La pregunta de la próxima etapa es: ¿por qué no interpretamos directamente sobre ese árbol?"

---

## Etapa 2 — Del parse tree al AST propio

**Objetivo pedagógico:** justificar el AST propio (desacoplamiento + control del reporte) y explicar la decisión más defendible del proyecto: traductor por `type-switch` en vez del Visitor generado.

**Tiempo sugerido:** 10-12 min.

**Puntos clave para enfatizar oralmente:**
- Frase ancla: *"no implementamos el patrón Visitor clásico del Dragón — usamos un `type-switch` porque cada alternativa etiquetada de la Etapa 1 YA genera su propio tipo Go concreto"*. Si te preguntan "¿conocen el patrón Visitor?", esta es la respuesta que demuestra que sí lo conocen y que decidieron no usarlo por una razón concreta, no por desconocimiento.
- El detalle de `traducirLugar` vs `traducirExpr` (mismos nodos para lectura y escritura) es sutil pero importa: el AST no distingue "leer" de "mutar" — esa decisión vive en el intérprete (Etapa 3), no en el AST. Es una buena respuesta si preguntan "¿por qué el AST es tan simple?".
- El recorrido por `reflect` en `grafo.go` es mantenimiento gratis: agregás un tipo de nodo nuevo y el reporte de AST lo incorpora solo, sin tocar `grafo.go`.

**Demo en vivo:** corré `go run ./cmd/cli entradas/ejemplo1_basico.vch` y mostrá la sección `=== AST ===` de la salida (conteo de nodos y aristas) — o, si tenés el cliente arriba, ejecutá cualquier ejemplo y andá a la pestaña "AST" para mostrar el grafo con vis-network en vivo.

**Preguntas probables:**
- *"¿Qué habría que implementar si el proyecto SÍ hubiera usado el patrón Visitor generado?"* → Las ~60 firmas de `BaseVLangCherryVisitor`.
- *"Si agregan un tipo de nodo nuevo al AST, ¿hay que modificar `grafo.go`?"* → No, si sus campos son exportados y de tipo reconocido. Sí hay que agregar el `case` en el traductor y, si corresponde, la etiqueta `#nombre` en el `.g4`.
- *"¿Por qué `numeros[i] = x` y leer `numeros[i]` usan el mismo nodo `ast.ExprIndice`?"* → Porque el traductor no distingue lectura de escritura; el intérprete decide según el contexto (Etapa 3, `resolverLugar`).

**Transición:** "Con el AST armado, veamos cómo el intérprete logra que structs y slices se pasen 'por referencia' sin ningún mecanismo especial."

---

## Etapa 3 — Tipos estáticos, structs y métodos "por referencia" gratis

**Objetivo pedagógico:** demostrar que `Valor.Slice`/`Valor.Struct` siendo punteros de Go resuelve el requisito de paso por referencia (7.1 del enunciado) sin tabla de referencias ni "boxing" explícito.

**Tiempo sugerido:** 10-12 min. Ideal candidato para pedir tiempo extra si el tribunal se engancha — es uno de los puntos de diseño más elegantes del proyecto.

**Puntos clave para enfatizar oralmente:**
- El truco completo en una frase: *"Go copia structs por valor automáticamente, pero si `Valor` tiene un campo `*StructVal` (puntero), copiar el `Valor` copia el puntero, no lo apuntado — dos variables VLangCherry que comparten un struct terminan apuntando al mismo `StructVal` real"*. Practicá decir esto de memoria, es la pregunta más probable de toda la defensa.
- Contrastá receptor por valor (`func (p Persona) Saludar()`) vs receptor por puntero (`func (p *Persona) Cumplir()`) con el ejemplo real de `entradas/ejemplo2_structs.vch`: `Saludar` no muta nada, `Cumplir` sí incrementa `Edad` y el cambio persiste afuera.
- `append` es la ÚNICA excepción a la regla de referencia compartida: siempre devuelve un `SliceVal` nuevo (como el Go real), por eso hace falta `numeros = append(numeros, 60)` con reasignación explícita. Si no marcás esta excepción, alguien te puede acorralar preguntando "¿entonces por qué necesito reasignar si todo es por referencia?".
- La relajación int↔float64 (`tipoCompatible`/`coercionar`) es una decisión documentada frente a una inconsistencia real del enunciado (3.6 exige tipo exacto, pero 4.3 muestra promoción automática) — no es un bug, es una decisión consciente y documentada en `docs/gramatica.txt`.

**Demo en vivo:** corré `entradas/ejemplo2_structs.vch` (`go run ./cmd/cli entradas/ejemplo2_structs.vch` o desde el cliente) y mostrá cómo `persona.Edad` pasa de 25 a 26 después de `persona.Cumplir()`. Si te sentís cómodo, cambiá en vivo el receptor de `Cumplir` de `(p *Persona)` a `(p Persona)` y volvé a correr — mostrá que en la práctica el cambio de `Edad` SIGUE reflejándose porque `Struct` ya es puntero (matiz real documentado en la respuesta del quiz de esta etapa: la diferencia entre valor/puntero se nota sobre todo cuando se reasigna el struct entero o un campo anidado, no un escalar simple).

**Preguntas probables:**
- *"Si `Cumplir` tuviera receptor `(p Persona)`, ¿qué imprimiría `persona.Edad` después?"* → En la práctica, sigue reflejando el cambio porque `Valor.Struct` ya es puntero — la distinción real entre valor y puntero se nota con reasignaciones de campos anidados o del struct completo (7.1 del enunciado).
- *"¿Por qué `append` necesita reasignación si structs y slices ya se comparten por puntero?"* → Porque `append` es la excepción intencional: siempre devuelve un `SliceVal` nuevo, igual que el Go real.
- *"¿Es válido asignar `"hola"` a una variable `int`?"* → No, error semántico. La única mezcla relajada es int↔float64.

**Transición:** "Los slices ya aparecieron en el ejemplo de `append`. Veamos ahora cómo funcionan las dos dimensiones."

---

## Etapa 4 — Slices 1D/2D

**Objetivo pedagógico:** mostrar que no existe un tipo especial para "2D" — un `[][]int` es un `SliceVal` cuyo `TipoElem` es a su vez un tipo slice, pura recursión de la misma estructura.

**Tiempo sugerido:** 8-10 min. Es la etapa más liviana después de la 0 — no te extiendas de más.

**Puntos clave para enfatizar oralmente:**
- Repetí la idea de "sin caso especial": ni en el runtime (`SliceVal`) ni en la gramática (`tipoSlice: '[' ']' tipo`, definida en términos de sí misma) hay una rama separada para 2D. Es la misma regla aplicada dos veces.
- Las 4 funciones nativas sobre slices (`len`, `append`, `indexOf`, `join`) son sección 7.2 del enunciado — tenelas frescas porque son terreno fértil para "agregá una función nativa nueva" en la defensa (ver sección final de este guion).
- El literal 2D usa **filas** (`{1,2,3}`, `filaSlice` en la gramática), no una lista plana — mostralo con el ejemplo de la matriz.

**Demo en vivo:** `entradas/ejemplo3_slices.vch` — mostrá el `for f, fila in mtx { for c, val in fila { ... } }` anidado y la salida `mtx[0][0] = 1`, etc.

**Preguntas probables:**
- *"¿Existe un tipo `Slice2D` separado de `SliceVal`?"* → No. Es el mismo `SliceVal`, con `TipoElem` siendo a su vez un tipo slice.
- *"¿Por qué no alcanza con `append(numeros, 60)` sin reasignar?"* → Porque `append` siempre produce un `SliceVal` nuevo (Etapa 3).
- *"En el `for` anidado, ¿qué tipo tiene `fila` en cada iteración externa?"* → `[]int` — cada elemento de un `[][]int` es a su vez `[]int`.

**Transición:** "Con tipos, structs y slices resueltos, falta la pieza más grande: todo el control de flujo — y ahí encontramos algo interesante durante una auditoría del código."

---

## Etapa 5 — Control de flujo, recursión y la auditoría de validación semántica ★★★

**Esta es la etapa a la que más tiempo y energía le tenés que dedicar en la defensa.** Los 4 hallazgos de la auditoría son la evidencia más fuerte de que el equipo entiende la diferencia entre "el programa corre" y "el programa valida lo que el enunciado exige" — que es, en el fondo, lo que un jurado de compiladores quiere ver.

**Objetivo pedagógico:** cubrir `if/else if/else`, `switch` sin fall-through, las 3 formas de `for`, recursión con llamada hacia adelante — y presentar los 4 hallazgos de la auditoría del 2026-07-21 como evidencia de rigor semántico.

**Tiempo sugerido:** 15-20 min (la etapa más larga de las 8, con toda intención).

**Puntos clave para enfatizar oralmente, en orden:**

1. **`switch` sin fall-through es una decisión de diseño DEL LENGUAJE, no una limitación de la implementación** — la sección 4.7.2 del enunciado dice literalmente "el break implícito está incluido al final de cada case". Si alguien pregunta "¿por qué no soportan fall-through como C?", la respuesta es: el enunciado no lo pide, y de hecho lo prohíbe explícitamente.

2. **`profLoop`/`profSwitch`**: dos contadores enteros que llevan cuántos `for`/`switch` anidados rodean la sentencia actual. Practicá explicar por qué se resetean a 0 al invocar una función (`invocarFuncion`): así un `break` "perdido" dentro de una función llamada desde un `for` no hereda el contexto del llamador. Esto es un detalle fino que muestra dominio del alcance semántico.

3. **Recursión y llamada hacia adelante**: `main` llama a `factorial`, declarada más abajo — funciona por el intérprete de **dos pasadas** (registra TODAS las funciones antes de ejecutar main). Conectá esto con la Etapa 0.

4. **Los 4 hallazgos de la auditoría — memorizalos, son tu mejor material de defensa:**
   - **Hallazgo 1 — `break`/`continue` fuera de contexto**: antes de la corrección, no se detectaba. Ahora `profLoop == 0 && profSwitch == 0` dispara error semántico (4.8.1/4.8.2 del enunciado).
   - **Hallazgo 2 — comparación relacional de strings**: la sección 4.4.1 dice "las comparaciones entre cadenas se hacen lexicográficamente" bajo el título de igualdad, pero describe un criterio de ORDEN. Faltaba aplicarlo también a `<`, `<=`, `>=`, `>` (función `Relacional` en `operaciones.go`, usa `strings.Compare`).
   - **Hallazgo 3 — colisión de nombres en ámbito global**: el enunciado (7.1) exige que función, variable y struct no compartan nombre. Faltaba validar la colisión ENTRE las tres categorías, no solo redeclaración dentro de la misma.
   - **Hallazgo 4 — el `switch` que silenciaba errores de tipo**: comparar tipos incompatibles dentro de un `case` (ej. `int` contra `struct`) se descartaba en silencio, como si simplemente no coincidiera — inconsistente con cómo `==` reporta ese mismo error fuera de un switch.
   - Cerralo con la frase del quiz real: *"ninguno de los 4 es un fallo de ejecución — todos son huecos de validación que el enunciado exige cerrar"*.

**Demo en vivo (la más importante de toda la charla):** corré `entradas/ejemplo6_errores.vch` — es EL archivo diseñado a propósito para disparar error léxico, sintáctico y varios semánticos en una sola ejecución (asignación de tipo incompatible, variable no declarada, campo inexistente, redeclaración, división entre cero, carácter inválido `#`, condición de `if` no booleana). Mostralo en el cliente para que se vea la tabla de errores completa con línea/columna, o por CLI para verlo más rápido. Es la prueba más contundente de que el manejo de errores (sección 8.1 del enunciado) funciona de punta a punta sin abortar la ejecución.

**Preguntas probables:**
- *"¿Por qué `factorial` puede invocarse desde `main` aunque se declara más abajo?"* → Intérprete de dos pasadas: todas las funciones se registran antes de ejecutar cualquier cuerpo.
- *"Un `break` dentro de una función `f()`, invocada desde un `for` en `main`, ¿es válido?"* → No. `invocarFuncion` resetea `profLoop`/`profSwitch` a 0 al entrar a `f()`.
- *"¿Qué tienen en común los 4 hallazgos de la auditoría?"* → Los 4 son casos donde el análisis semántico dejaba pasar (o descartaba en silencio) algo que el enunciado exige reportar como error.
- **Pregunta de defensa que casi seguro te van a hacer:** *"¿Por qué tu intérprete detecta este error?"* señalando alguno de los 4 casos → Tené lista, para CADA hallazgo, la cita exacta de la sección del enunciado que lo exige (4.8.1/4.8.2, 4.4.1, 7.1, 4.4.1) — está todo en las diapositivas 8-11 de esta etapa, calcala antes de entrar.

**Transición:** "El intérprete ya cubre todo el lenguaje. Ahora veamos cómo se expone por HTTP."

---

## Etapa 6 — El servidor REST con `net/http`

**Objetivo pedagógico:** mostrar que el servidor no usa ningún framework (decisión deliberada), que cada petición crea un `Interprete` nuevo (entorno fresco), y que hay recuperación ante pánico interno.

**Tiempo sugerido:** 8-10 min.

**Puntos clave para enfatizar oralmente:**
- Solo 3 rutas: `GET /`, `GET /salud`, `POST /interpretar`. El endpoint que importa es uno solo.
- "Entorno fresco por petición" es el mismo criterio que DataForge y ConjAnalyzer aplican en Java — remarcá que es un patrón que se repite en TODO el curso, no una idea aislada de VLangCherry.
- El `defer recover()` alrededor de traducción + interpretación es importante para un servidor de larga duración: un panic por un caso interno no contemplado no debe tirar abajo el proceso completo, solo esa petición puntual.
- La compilación cruzada (`GOOS=linux GOARCH=amd64 go build`) responde a un requisito de entrega real (restricción 10.2): el enunciado exige ejecución nativa en Linux, y el desarrollo se hizo en Windows.

**Demo en vivo:** con el servidor ya corriendo (Terminal 1 de la preparación), mostrá en el navegador o con curl: `http://localhost:4100/salud` → `{"estado":"ok"}`. Si querés algo más vistoso, hacé un POST a `/interpretar` con curl usando el contenido de un `.vch` corto como `codigo`.

**Preguntas probables:**
- *"¿Qué framework HTTP usa el servidor?"* → Ninguno — solo `net/http`, la librería estándar. Decisión deliberada.
- *"Si dos peticiones concurrentes ejecutan código distinto, ¿pueden mezclarse sus errores?"* → No — cada llamada a `Analizar(codigo)` crea su propia instancia de `Interprete` con su propia `ListaErrores`.
- *"¿Qué pasa si el código dispara un panic interno de Go no contemplado?"* → El `defer recover()` de `analizar.go` lo atrapa y lo reporta como error semántico adicional (línea/columna 0); el servidor sigue en pie.

**Transición:** "El backend ya expone todo por HTTP. Falta ver quién consume ese contrato del lado del navegador."

---

## Etapa 7 — El cliente React, reusando componentes

**Objetivo pedagógico:** mostrar que la interfaz (editor multi-archivo, consola, tablas, grafo de AST) se adaptó de otro proyecto del curso, cambiando solo la capa de comunicación (`api.js`) para apuntar al servidor Go.

**Tiempo sugerido:** 8-10 min. Si el tiempo aprieta, esta es la primera candidata a comprimir — el jurado suele centrarse más en el intérprete que en el cliente.

**Puntos clave para enfatizar oralmente:**
- El contrato JSON es `{errores, consola, consolaLineas, simbolos, ast, dot}` — el mismo (en forma) que otro proyecto del curso con arquitectura cliente-servidor equivalente, más `consolaLineas` y `dot`.
- Un solo archivo conoce la URL del backend: `api.js` vía `VITE_API_URL`. Todo el resto de los componentes solo conoce la forma del JSON.
- `AstGrafo.jsx` usa vis-network con layout jerárquico (`direction: 'UD'`) — mismo formato `{nodes, edges}` que consumía el grafo del otro proyecto.
- La suite Playwright (`client/e2e/vlangcherry.spec.js`) corre 3 pruebas reales: ejecución del ejemplo inicial (verifica consola, símbolos y AST), reporte de errores semánticos con salto de línea, y creación de archivo nuevo `.vch`. Tenerla en verde es un argumento fuerte de "el proyecto está probado, no solo demostrado a mano".

**Demo en vivo:** con cliente y servidor arriba, ejecutá el ejemplo inicial con el botón "▶ Ejecutar", mostrá las 4 pestañas (Consola, Errores, Símbolos, AST) navegando entre ellas. Si el tiempo alcanza, corré `npm run e2e` en vivo desde `client/` y mostrá que las 3 pruebas pasan en verde — es un cierre visual fuerte antes de pasar a la profundización o al cierre de la charla.

**Preguntas probables:**
- *"¿Qué archivo hay que tocar si el servidor cambia de puerto o dominio?"* → `api.js`, vía `VITE_API_URL`.
- *"¿Por qué CORS está habilitado con `*` en todas las rutas?"* → Porque en desarrollo cliente (Vite) y servidor (Go) corren en orígenes distintos; sin CORS el navegador bloquearía el `fetch`.
- *"Si agregan un campo nuevo al AST, ¿hace falta cambiar `AstGrafo.jsx`?"* → No, mientras el formato siga siendo `{nodes, edges}` — el recorrido por reflection ya lo incorpora, y vis-network solo necesita ese formato genérico.

**Transición:** "Con las 8 etapas completas, cerramos con una mirada comparativa: cómo se ve este mismo problema resuelto con las otras dos herramientas del curso."

---

## Profundización ★ — ANTLR4 vs JFlex+CUP vs Jison

**Objetivo pedagógico:** poner los 3 generadores lado a lado (herramientas, archivos fuente, lenguaje de salida, algoritmo de parsing, patrón de recorrido) para mostrar que la elección de herramienta sigue al requisito del proyecto, no al revés.

**Tiempo sugerido:** 8-10 min. Es opcional si el tiempo de la defensa es corto — priorizá Etapa 5 y el cierre antes que esta profundización.

**Puntos clave para enfatizar oralmente:**
- La tabla comparativa completa (herramientas, archivos, target, algoritmo, patrón de árbol, proyecto que lo usa) es un resumen perfecto para cerrar toda la charla del curso, no solo de VLangCherry.
- La frase de cierre de la etapa es la más citable de toda la presentación: *"la herramienta no decide el diseño — lo habilita. El requisito del enunciado sí decide."* Usala si te preguntan por qué el equipo no siguió el mismo stack que los otros proyectos.
- Mencioná los patrones que se repiten en los 5 proyectos del curso pese a usar lenguajes distintos: "entorno fresco por ejecución", "propagación de error sin cascada", "reporte de AST por recorrido genérico" — son decisiones de diseño de compiladores, no atadas a una herramienta puntual.

**Demo en vivo:** ninguna — es una etapa de síntesis comparativa, no de código nuevo.

**Preguntas probables:**
- *"¿Cuál de los tres generadores resuelve léxico y sintaxis en archivos separados?"* → Solo JFlex+CUP.
- *"¿Qué tienen en común los parsers que genera JFlex+CUP y los que genera Jison?"* → Ambos son LALR (Jison también soporta SLR/LR); ANTLR4 es LL adaptativo.
- *"Si CompInterpreter (Jison) hubiera necesitado Go en vez de JS, ¿habría podido seguir con Jison?"* → No — Jison solo genera JavaScript.

---

## Para la defensa oral — modificaciones que el equipo debería poder hacer en vivo con confianza

Esto es lo que pidió el enunciado de la defensa: van a pedirles que toquen código en vivo para probar autoría real, no memorización de diapositivas. Estos 5 puntos son los que, según la arquitectura ya documentada, más probablemente les pidan — practicalos ANTES del día, no los improvisen ahí.

### 1 · Agregar una función nativa nueva

**Dónde:** `server/internal/runtime/nativas.go`. Hay que tocar exactamente dos lugares:
1. Agregar el nombre al `switch` de `EsNombreNativa` (línea ~16-20).
2. Agregar un `case` nuevo en el `switch` de `llamarNativa` (línea ~28 en adelante), siguiendo el patrón exacto de las funciones existentes: validar cantidad de argumentos con `in.errorSemantico(...)` si falla, y devolver `(Valor, bool)` — el `bool` en `false` indica error ya reportado.

Ejemplo concreto para practicar en vivo: agregar `sum(slice)` que sume todos los elementos de un `[]int`. Miren `len` como plantilla (valida `len(args) != 1` y el tipo con un `switch v.Tipo.Base`) y calquen la estructura.

**Qué decir mientras lo hacen:** "Esto no requiere tocar la gramática ni el traductor porque las llamadas a función ya están resueltas genéricamente — `EsNombreNativa` es solo una lista de nombres reservados, y `llamarNativa` es el único lugar que despacha por nombre."

### 2 · Agregar un tipo de error semántico nuevo (una validación nueva)

**Dónde:** `server/internal/runtime/interprete.go`, siguiendo exactamente el patrón de los 4 hallazgos de la auditoría de la Etapa 5. La función clave es `in.errorSemantico(desc string, linea, col int)`.

Ejemplo concreto para practicar: agregar una validación de que una variable declarada `mut x int` no pueda reasignarse con un tipo distinto la segunda vez (si no está ya cubierta) o, más simple todavía, replicar en vivo el Hallazgo 1 completo (quitarlo de código a propósito antes de la defensa y volver a agregarlo en vivo): el patrón es

```go
if in.profLoop == 0 && in.profSwitch == 0 {
    in.errorSemantico("la sentencia \"break\" no puede usarse fuera de un ciclo o un switch", n.Linea, n.Columna)
    return Senal{Tipo: SenalNinguna}
}
```

**Qué decir mientras lo hacen:** citen la sección exacta del enunciado que exige esa validación (así demuestran que la corrección no es arbitraria) y expliquen que `errorSemantico` acumula en `ListaErrores` sin abortar la ejecución — el resto del programa sigue corriendo después de reportar el error.

### 3 · Explicar y demostrar por qué los structs son "por referencia" gratis

No es una modificación de código sino una explicación que sí les pueden pedir defender palabra por palabra, así que practiquen decirla de memoria:

"`Valor` es un struct de Go con campos `Slice *SliceVal` y `Struct *StructVal` — ya son punteros. Go copia structs por valor automáticamente al pasar un parámetro o hacer una asignación, pero copiar un puntero no copia lo apuntado. Entonces cuando 'copiamos' un `Valor` que representa un struct VLangCherry, en realidad copiamos el puntero — las dos variables terminan compartiendo el mismo `StructVal` real en memoria. No hay tabla de referencias ni mecanismo de boxing: es una consecuencia directa de cómo Go maneja sus propios punteros."

Si les piden demostrarlo en vivo: usen `entradas/ejemplo2_structs.vch`, comenten temporalmente el receptor por puntero de `Cumplir` (`func (p *Persona) Cumplir()`) y prueben con receptor por valor (`func (p Persona) Cumplir()`) — corran de nuevo y muestren que `persona.Edad` sigue reflejando el cambio de todos modos, porque `Struct` ya es puntero independientemente del receptor. Esto es un matiz real (documentado en la respuesta del quiz de la Etapa 3) que demuestra dominio fino, no superficial, del modelo de valores.

### 4 · Agregar un operador nuevo a la gramática (recorrido completo del pipeline)

Es la modificación más completa que pueden mostrar — toca los 3 archivos que van a tener abiertos según la preparación de este guion:

1. **`grammar/VLangCherry.g4`**: agregar una alternativa nueva a la regla `expr`, CON su etiqueta `#nombre` (ej. un operador de potencia `^`, alternativa `expr '^' expr # exprPotencia`, ubicada según la precedencia que quieran darle).
2. Regenerar con el comando real:
   ```
   java -jar tools/antlr.jar -Dlanguage=Go -visitor -o server/internal/parser -package parser grammar/VLangCherry.g4
   ```
3. **`server/internal/traductor/traductor.go`**: agregar el `case *parser.ExprPotenciaContext` correspondiente en `traducirExpr`.
4. **`server/internal/runtime/operaciones.go`** (o donde corresponda): implementar la operación en tiempo de ejecución.

**Qué decir mientras lo hacen:** remarquen que si olvidan la etiqueta `#nombre` en el paso 1, ANTLR4 no genera un tipo de contexto Go propio, y el `type-switch` del traductor no puede distinguir esa alternativa de las demás — es literalmente la pregunta 1 del quiz de la Etapa 1, así que si se las hacen ya tienen la respuesta ensayada con código real.

### 5 · Explicar (y si hace falta, romper/arreglar en vivo) el intérprete de dos pasadas

Concepto: `Interpretar(programa *ast.Programa)` registra TODOS los structs primero, luego TODAS las funciones y métodos, y recién después crea el entorno global y ejecuta `main()`. Esto es lo que permite llamadas hacia adelante (`main` llamando a `factorial`, declarada más abajo — `entradas/ejemplo5_funciones.vch`).

Si les piden demostrar que entienden por qué es necesario, un truco efectivo en vivo: expliquen (sin necesariamente tener que romper el código real) qué pasaría si el intérprete fuera de una sola pasada y ejecutara cada declaración según aparece en el archivo — `main()` fallaría con "función no declarada" al intentar llamar a `factorial()`, porque en ese momento el intérprete todavía no habría llegado a leer su declaración. Las dos pasadas existen exactamente para evitar esa dependencia del orden textual.

**Qué decir para conectar con el resto:** "Es el mismo criterio de 'entorno fresco por ejecución' que vimos en la Etapa 6 para el servidor — cada llamada a `Analizar()` crea un `Interprete` nuevo, y cada `Interprete` nuevo vuelve a hacer sus dos pasadas desde cero. Nunca hay estado que sobreviva entre ejecuciones ni sesgo por orden de declaración dentro de una misma ejecución."

---

## Cierre

VLangCherry es el único de los 5 proyectos del curso construido con ANTLR4 sobre Go — una elección que no fue estilística sino forzada por el enunciado (Go obligatorio → se necesita un generador que emita Go → ANTLR4 es el único de los tres del curso que lo hace). A partir de ahí, 4 decisiones de diseño atraviesan todo el proyecto: AST propio desacoplado del parse tree de ANTLR, un traductor por `type-switch` que aprovecha las etiquetas `#nombre` de la gramática en vez del patrón Visitor completo, referencia "gratis" para structs y slices vía punteros de Go (`Valor.Slice`/`Valor.Struct`), y un intérprete de dos pasadas con entorno fresco por ejecución — el mismo criterio que DataForge y ConjAnalyzer aplican en Java.

La auditoría de la Etapa 5 (4 hallazgos de validación semántica corregidos el 2026-07-21: break/continue fuera de contexto, comparación relacional de strings, colisión de nombres en ámbito global, switch que silenciaba errores de tipo) es el argumento más fuerte que tienen para la defensa: demuestra que el equipo no solo hizo que el intérprete "corriera", sino que volvió sobre el propio código con espíritu crítico para verificar que validara exactamente lo que el enunciado exige — ni más, ni menos.

Para la defensa oral en vivo, practiquen los 5 puntos de la sección anterior ANTES del día: agregar una función nativa, agregar una validación semántica nueva, explicar (y poder demostrar) por qué los structs son por referencia gratis, agregar un operador completo a la gramática (los 3 archivos), y explicar por qué el intérprete necesita dos pasadas. Si dominan esos 5 movimientos con soltura, cualquier variación que les pidan en vivo va a ser una aplicación directa del mismo patrón que ya practicaron.

Éxitos en la defensa.
