# Guion de clase — CompInterpreter

Nota de trabajo para vos (quien da la clase), no para la audiencia. Las diapositivas HTML están en 3ª persona porque son para los estudiantes; este documento es en 2ª persona porque es tu chuleta.

Fuentes que usé para armar esto: las 8 `etapaN.html` + `index.html` + `README.md` + `arquitectura-rest.html` de `presentacion-compinterpreter/`, el `docs/ManualTecnico.md` de `CompInterpreter/`, los 5 `entradas/*.ci`, `client/playwright.config.js`, `client/package.json` y `server/package.json`, y `server/test-cli.js`.

---

## Antes de empezar

**Abrir la presentación:** doble clic en `presentacion-compinterpreter\index.html`. Es autocontenida (mismo CSS/JS que las otras presentaciones del curso), no necesita servidor ni internet — no dependas de que haya wifi en el salón.

**Duración total estimada:** entre 1h45 y 2h15, con preguntas. Es el proyecto más largo de explicar de los 4 porque tiene una capa más (cliente-servidor) que los otros tres no tienen. Si vas corto de tiempo, la Etapa 4 (auditoría de bugs) y la profundización ★ (arquitectura REST) son las dos piezas que más justifican no recortar; si hay que sacrificar algo, la Etapa 6 (detalles de React) es la que más se puede resumir sin perder el argumento central.

**Preparar las demos en vivo — este proyecto es cliente-servidor, así que necesitás DOS terminales además de la presentación:**

1. **Terminal servidor** (déjala corriendo toda la clase):
   ```
   cd CompInterpreter/server
   npm start
   ```
   (equivale a `node server.js`; escucha en `http://localhost:4000`). Vas a ver los WARNING normales de Node si los hay — no son problema. Si tocás `grammar.jison` alguna vez antes de clase, regenerá el parser con `npm run grammar` antes de arrancar.

2. **Terminal cliente** (también corriendo toda la clase):
   ```
   cd CompInterpreter/client
   npm run dev
   ```
   (Vite, puerto `5173` por defecto). Abrí `http://localhost:5173` en el navegador y dejá esa pestaña lista — es la app real que vas a mostrar en la Etapa 6.

3. **Demo liviana sin levantar nada de lo anterior** (útil para las Etapas 1 a 4, cuando todavía no llegaste a hablar del servidor/cliente): desde `CompInterpreter/server`, corré directamente el pipeline por consola con el script de prueba que ya existe en el proyecto:
   ```
   node test-cli.js ../entradas/ejemplo_errores.ci
   ```
   Imprime consola, errores (con línea/columna), símbolos y el conteo de nodos/aristas del AST — sin necesidad de servidor Express ni navegador. Cambiá el argumento por cualquiera de los 5 archivos de `entradas/` según qué quieras mostrar (ver qué demo corresponde a cada etapa, abajo).

4. **Para la Etapa 7 (Playwright):** desde `CompInterpreter/client`:
   ```
   npx playwright test
   ```
   Ojo: el `playwright.config.js` ya levanta servidor y cliente automáticamente si no están corriendo (`webServer` con `reuseExistingServer: true`), así que si dejaste las Terminales 1 y 2 abiertas, Playwright las reutiliza en vez de levantar copias nuevas — más rápido y evita choques de puerto.

**Ensayá las 4 demos vos mismo ANTES de la clase, una vez completa de punta a punta.** Es el proyecto con más piezas móviles del curso (Node + Express + Vite + navegador); si algo falla en vivo (puerto ocupado, `node_modules` desactualizado), preferís descubrirlo ahora y no con el grupo mirando.

---

## Etapa 0 — Conceptos base y por qué cambia el stack

**Objetivo pedagógico:** que quede claro que la teoría no cambia entre proyectos — cambia la herramienta, y por primera vez en el curso, la arquitectura.

**Tiempo sugerido:** 8 minutos.

**Puntos clave a enfatizar oralmente (no repitas la tabla completa, solo estos tres):**
- La tabla DataForge-vs-CompInterpreter la podés dejar que la lean solos; vos verbalizá el punto que más importa: **el AST es la diferencia real de diseño**, no el lenguaje de implementación. Java-vs-JS es superficial; construir árbol o no, no.
- Marcá bien que la arquitectura cliente-servidor **no es una elección de estudiante**: citá la sección 9 del enunciado y el objetivo específico — es un requisito. Esto evita la pregunta "¿por qué no hicieron algo más simple tipo DataForge?".
- El repaso de token/lexema/patrón es rápido a propósito — es literalmente el mismo vocabulario del Dragón cap. 3 que ya vieron en DataForge. No te detengas ahí, es solo para "enganchar" a quien no se acuerde.

**Demo en vivo:** no hace falta código todavía. Si querés algo visual, mostrá en el explorador de archivos las carpetas `CompInterpreter/server` y `CompInterpreter/client` lado a lado — refuerza visualmente "dos procesos independientes" antes de que aparezca en las diapositivas de más adelante.

**Preguntas probables:**
- *"¿Por qué no armaron esto con JavaFX como los otros tres?"* — Porque JavaFX es una tecnología de escritorio y el enunciado exige explícitamente que se visualice todo desde el navegador (sección 9) y que se use arquitectura cliente-servidor como objetivo específico. No es una preferencia, es un requisito de proyecto.
- *"¿Jison es tan bueno como JFlex+CUP?"* — Es la misma clase de analizador (LALR(1)), el mismo principio teórico. La diferencia práctica es que Jison junta léxico y sintáctico en un solo archivo; JFlex+CUP los separa en dos herramientas. Ninguno es "mejor" en teoría de compiladores — cambia la ergonomía de la herramienta.
- *"Si no hay AST en DataForge, ¿por qué acá sí?"* — Porque acá hay control de flujo real: ciclos que repiten su cuerpo N veces, condicionales que se recorren selectivamente, funciones invocables antes de su declaración textual. Nada de eso se resuelve ejecutando una sola vez en el orden que reconoce el parser.

**Transición:** "Ya mapeamos el terreno — ahora vamos a la primera pieza concreta: cómo Jison arma el léxico completo en un solo archivo."

---

## Etapa 1 — La gramática Jison: tabla de tokens y léxico en `%lex`

**Objetivo pedagógico:** mostrar que el mismo principio de longest-match y orden de reglas de JFlex aplica igual en Jison — la herramienta cambia, la teoría no.

**Tiempo sugerido:** 10 minutos.

**Puntos clave a enfatizar oralmente:**
- `%options case-insensitive flex` resuelve en una sola línea algo que en JFlex necesitaba tratamiento carácter por carácter — buen ejemplo de "una herramienta distinta, mismo problema, distinta ergonomía".
- El orden de las reglas importa exactamente igual que en JFlex: reservadas antes de `ID`, operadores de dos caracteres antes que los de uno (`==` antes que `=`). Si alguien pregunta "¿por qué ese orden?", la respuesta es siempre la misma: a igual longitud de coincidencia, gana la regla escrita primero.
- Los escapes de cadenas (`\n`, `\t`, etc.) se resuelven **dentro de la propia regla léxica**, no más adelante. Vale la pena remarcar el "por qué": es un fenómeno puramente léxico, no depende de contexto sintáctico ni semántico.
- El objeto `yy` compartido entre lexer y parser es la pieza que hace posible que TODOS los errores (léxicos, sintácticos, y después semánticos) terminen en un solo arreglo ordenable por línea/columna. Es un buen anticipo de la Etapa 5.

**Demo en vivo:**
```
cd CompInterpreter/server
node test-cli.js ../entradas/ejemplo_errores.ci
```
Este archivo tiene un `#` intencional en la línea 5. La salida real confirma el desglose exacto que muestra la diapositiva: 1 error léxico + 2 sintácticos + 5 semánticos = 8 total, con línea y columna en cada uno. Es la misma verificación que después reaparece en el test de Playwright de la Etapa 7 — mencionalo, ayuda a que la clase vea la continuidad.

**Preguntas probables:**
- *"¿Qué pasa si `ID` se declara antes que las reservadas?"* — Ninguna reservada se reconocería jamás; todo `if`, `let`, `function` pasaría a ser un identificador común. Es el quiz 1 de la etapa, tal cual.
- *"¿Por qué no dejar que el intérprete resuelva los escapes de `\n`?"* — Porque cada consumidor posterior (intérprete, reportes) tendría que repetir la misma lógica de reemplazo. Resolverlo en la regla léxica es hacerlo una sola vez, en el lugar correcto.
- *"¿Qué gana el proyecto compartiendo `yy` entre lexer y parser?"* — Un único arreglo de errores en el orden real en que ocurrieron, sin tener que fusionar fuentes distintas al armar el reporte final.

**Transición:** "El léxico ya entrega una secuencia de tokens limpia. Ahora — ¿cómo se convierte eso en un árbol?"

---

## Etapa 2 — El parser: gramática sintáctica, precedencia y el AST

**Objetivo pedagógico:** que entiendan la fábrica de nodos como la pieza que permite un AST genérico, y por qué la tabla de precedencia de Jison no es lo mismo que la gramática BNF entregable.

**Tiempo sugerido:** 12 minutos.

**Puntos clave a enfatizar oralmente:**
- La tabla de precedencia (`%right`/`%left`/`%nonassoc`) reproduce nivel por nivel la tabla 5.10 del enunciado — mostralo en paralelo si tenés tiempo, es un buen ejercicio de "leer una tabla de precedencia de una herramienta real".
- `%nonassoc` en potencia/raíz es la pieza más preguntada: remarcá que **no es un capricho de la herramienta**, es literal del enunciado — `2^3^2` sin paréntesis debe ser error de gramática, no algo resuelto por defecto a izquierda o derecha.
- La función `nodo(tipo, props, loc)` es LA pieza central de esta etapa. Insistí en la idea: "cualquier nodo nuevo que se agregue al lenguaje, mientras respete `{tipo, ...}`, el resto del sistema (grafo del AST, Etapa 6) lo dibuja sin que nadie tenga que tocar ese código". Es el mismo argumento de "abierto a extensión, cerrado a modificación" que probablemente ya vieron en otras materias.
- Marcá el `docs/gramatica.txt` como el entregable real — la tabla de precedencia de Jison es un mecanismo de la herramienta, NO la gramática BNF que piden entregar. Esto es un punto de examen típico: alguien puede confundir "la gramática de Jison" con "la gramática BNF del proyecto".

**Demo en vivo:**
```
cd CompInterpreter/server
node test-cli.js ../entradas/ejemplo_vectores2d.ci
```
Mirá la línea final: `AST --- nodos: 106 aristas: 105`. Aprovechá para explicar en vivo la verificación aritmética barata: en cualquier árbol sin ciclos, aristas = nodos − 1. Es un buen momento para preguntarle a la clase "¿por qué esto prueba que es un árbol y no un grafo con ciclos?" antes de responder vos mismo.

**Preguntas probables:**
- *"¿Por qué potencia y raíz son `%nonassoc` y no `%left`/`%right`?"* — Porque el enunciado (5.10) dice explícitamente que no son asociativos; encadenarlos sin paréntesis debe ser error de gramática.
- *"Si la gramática BNF entregable no usa `%left`/`%right`, ¿cómo resuelve la precedencia?"* — Reescribiendo las expresiones en capas anidadas (`expresion-or → expresion-and → ... → expresion-primaria`), donde cada capa solo combina con la de abajo. Es ambigüedad resuelta por construcción, sin depender de la herramienta.
- *"¿Qué pasa si agrego un tipo de nodo nuevo al lenguaje — hay que tocar el código del grafo?"* — No, mientras el nodo nuevo tenga la propiedad `tipo` como string. El recorrido del grafo (Etapa 6) es genérico.

**Transición:** "Con el árbol ya construido, la pregunta que sigue es: ¿quién lo recorre, y cuántas veces?"

---

## Etapa 3 — El intérprete: entorno, tipos y dos pasadas

**Objetivo pedagógico:** entender por qué dos pasadas resuelven la referencia adelantada, y la diferencia entre coerción y cast (necesaria para que el ejemplo oficial del enunciado ni siquiera falle).

**Tiempo sugerido:** 12 minutos.

**Puntos clave a enfatizar oralmente:**
- Las dos pasadas están en el **enunciado mismo** (sección 5.21: "se recomienda realizar 2 pasadas") — no es una decisión libre del proyecto, es seguir al pie de la letra la sugerencia oficial. Vale la pena decir esto explícitamente para que no parezca sobre-ingeniería.
- El entorno como cadena de ámbitos (`Map` + referencia al padre) es el mismo concepto exacto que la tabla de símbolos de DataForge — no reinventes la explicación, apoyate en lo que ya vieron.
- La distinción coerción-vs-cast tiene un gancho perfecto: **sin coerción `double→int`, el propio ejemplo oficial del Anexo 11.1 sería inválido** (`let modulo: int = x % 3`, donde `%` entre enteros da `DECIMAL` según la tabla 5.5.6). Es un argumento fuerte para justificar por qué la coerción existe — no es capricho, es necesidad del propio enunciado.
- La propagación por `null` es idéntica en teoría a DataForge, pero con una trampa de nombres real: el `null` de JavaScript (señal interna de error) NO es lo mismo que el valor de tipo `'null'` del lenguaje `.ci` (una instancia válida de `Valor`). Vale la pena escribir ambos en el pizarrón/proyector si notás caras de confusión — es el tipo de detalle que confunde en el parcial.

**Demo en vivo:**
```
cd CompInterpreter/server
node test-cli.js ../entradas/ejemplo_funciones.ci
```
Este archivo tiene forward-reference real: `main()` llama a `factorial()` que está declarada más abajo en el archivo, y funciona. También tiene parámetros con valor por defecto y casteos (`cast(18.6 as int)` → `18`). Es la demo más completa para esta etapa porque toca 3 conceptos en un solo archivo.

**Preguntas probables:**
- *"¿Por qué `main()` puede llamar a `factorial()` si está declarada después en el texto?"* — Porque la pasada 1 recorre TODOS los elementos globales y registra cada función/método en un `Map` antes de que la pasada 2 ejecute nada. Para cuando `ejecutar main();` realmente corre, `factorial` ya está disponible.
- *"¿`x % 3` no debería dar error de tipos si `%` produce `DECIMAL` y se asigna a un `int`?"* — No, porque la asignación aplica coerción implícita `double→int` (trunca). Es la decisión de diseño documentada en el propio código para que el ejemplo oficial del enunciado sea válido.
- *"¿El `null` de JavaScript es lo mismo que un valor `null` de `.ci`?"* — No, para nada. El de JS es una señal interna de "ya hubo un error, no reportes otro". El de `.ci` es un `Valor` normal que el programador puede asignar y comparar a propósito.

**Transición:** "El intérprete ya recorre, verifica tipos y ejecuta. Pero hay una pieza que todavía no expliqué: ¿cómo funcionan `break`, `continue` y `return` sin que el intérprete use excepciones de JavaScript? Acá es donde entra la etapa más importante de mostrar bien."

---

## Etapa 4 — Señales de control: break/continue/return y la auditoría de código

**Esta es la etapa que más tenés que cuidar en tiempo y énfasis — no la apures.**

**Objetivo pedagógico:** entender el mecanismo de señales (y cómo resuelve gratis el fall-through de switch), y usar los 2 bugs reales de la auditoría como caso de estudio de por qué el manejo de errores importa tanto como el "camino feliz".

**Tiempo sugerido:** 15-18 minutos — es la etapa con más contenido narrativo y la que más vale la pena dramatizar un poco.

**Puntos clave a enfatizar oralmente:**
- Arrancá preguntando a la clase: *"¿cómo implementarían ustedes `break`/`continue`/`return` en JavaScript?"* — casi seguro alguien va a decir `throw`/`try-catch`. Dejá que lo digan, y después mostrá que el proyecto eligió otra cosa a propósito: un objeto de señal (`{tipo: 'BREAK', ...}`) que cada instrucción retorna hacia arriba. Esto genera más enganche que simplemente mostrar el código de entrada.
- El fall-through de `switch` "sale gratis" con el mismo mecanismo — es un buen momento para remarcar que una buena decisión de diseño a veces resuelve un problema que ni estabas buscando resolver todavía.
- Las guardas `MAX_ITER` (1,000,000) y `MAX_DEPTH` (2,000) no están en el enunciado — son una decisión propia de estabilidad del servidor. Explicá el motivo real: un `loop{}` sin salida o una recursión infinita colgaría el proceso Node **completo**, afectando a TODAS las peticiones concurrentes, no solo a la que causó el problema. Esto conecta directo con la Etapa 5 (por qué importa la concurrencia acá y no en un monolito).

**El caso de estudio de la auditoría — destacalo bien, es el corazón de esta etapa:**

El 21 de julio de 2026 el proyecto pasó una auditoría de código y se encontraron 2 bugs reales. Contalos como una historia, no como una lista de bullets:

1. **Bug 1 — vector con elemento incompatible tumbaba TODO el análisis.** Entrada: `let v: int[] = [1, "hola", 3];`. Antes de la corrección, el `null` de JavaScript (de la propagación normal por error) quedaba guardado crudo en esa posición del vector. El problema: las funciones nativas sobre vectores (`sum`, `max`, `min`, `average`, `aTexto()`) leen `.tipo` de cada elemento **sin volver a comprobar null** — asumen que todo elemento es siempre un `Valor`. Resultado: una excepción de JS no controlada abortaba TODA la interpretación, reportada genéricamente como "Error interno durante la ejecución" — un solo elemento mal escrito tumbaba el análisis completo del archivo. La corrección: `coercionarElementoOValorDefecto()` rellena con el valor por defecto del tipo (0, "", etc.) en vez de dejar el hueco en `null`.
2. **Bug 2 — `break`/`continue` fuera de un ciclo no generaba error.** El enunciado (5.18.2) exige literalmente que `continue` fuera de un ciclo sea un error. La validación en `invocar()` solo se aplicaba a `CONTINUE`, no a `BREAK` — por el mismo criterio, debía aplicarse a ambos. Consecuencia real: un `break;` mal puesto (por ejemplo directo en el cuerpo de una función, sin ciclo alrededor) truncaba en silencio el resto de las instrucciones de esa función, sin dejar rastro en consola ni en la tabla de errores. La corrección extiende la misma validación de `CONTINUE` a `BREAK`.

**El punto pedagógico que tenés que dejar bien claro, con tus propias palabras (no leas la diapositiva):** los dos bugs son fallas del *manejo de errores*, no del camino feliz — solo se manifestaban con una entrada específica que ningún ejemplo de prueba anterior había ejercitado a propósito. "Los ejemplos de prueba pasan" no es lo mismo que "el manejo de errores es correcto en todos los casos". Y el bug 2 es peor que un error visible: un error que no se reporta es el peor tipo de error para quien depura, porque el código simplemente deja de funcionar sin ninguna pista.

**Demo en vivo (dos partes):**

Parte A — fall-through del switch, funcionando:
```
node test-cli.js ../entradas/ejemplo_switch.ci
```
La consola muestra en vivo cómo `case 2` sin `break` cae en `case 3` y `case 4` hasta encontrar el primer `break`. Señalá cada línea de la salida contra el código fuente del `.ci` mientras corre.

Parte B — reproducir el bug 1 en vivo, ya corregido (opcional pero de alto impacto si tenés tiempo): escribí en un archivo temporal o pegá en el editor del cliente (si ya tenés el stack levantado) algo como:
```
let v: int[] = [1, "hola", 3];
```
y mostrá que hoy produce un único error semántico localizado con línea y columna — no una excepción genérica ni un cuelgue. Es la mejor forma de que la clase entienda "antes vs. después" sin tener que mostrarles el commit del fix.

**Preguntas probables:**
- *"¿Por qué no usar `throw`/`try-catch` de JS, que es más simple?"* — El código no lo documenta como limitación técnica (las excepciones funcionarían perfectamente); es una decisión de diseño: un objeto de señal es flujo de control explícito y visible en cada punto donde se decide qué hacer con él, mientras que una excepción viaja implícita. Y de paso resuelve gratis el fall-through de switch.
- *"En el bug 1, ¿por qué no alcanzaba con propagar `null` como en cualquier otra expresión?"* — Porque una posición de un vector no es una expresión suelta: la van a leer después funciones (`sum`, `max`, `aTexto()`...) que asumen que todo elemento es un `Valor` y leen `.tipo` sin chequear null. Un `null` crudo ahí provoca una excepción no controlada, no un error semántico limpio.
- *"¿Qué tienen en común los dos bugs?"* — Los dos son fallas del manejo de errores, no del camino feliz, y en ambos el síntoma era peor que un mensaje visible: silenciar el problema (bug 2) o abortar todo sin localizar la causa (bug 1).

**Transición:** "Con el control de flujo ya robusto y auditado, toca la pieza que distingue arquitectónicamente a este proyecto de los otros tres: el servidor que expone todo esto por HTTP."

---

## Etapa 5 — El servidor Express: contrato REST y entorno fresco

**Objetivo pedagógico:** entender `analizar.js` como el orquestador central, y por qué "entorno fresco por petición" acá tiene una razón adicional (concurrencia) que no existe en DataForge.

**Tiempo sugerido:** 10 minutos.

**Puntos clave a enfatizar oralmente:**
- `server.js` es corto A PROPÓSITO — remarcá la separación de responsabilidades: transporte HTTP en `server.js`, análisis/ejecución en `analizar.js` + `interprete/`. Esto es lo que permite que `analizar()` se pueda invocar directo desde consola (como en `test-cli.js`, que ya usaste en las etapas anteriores) sin levantar ningún servidor.
- El contrato JSON completo (`errores`, `consola`, `consolaLineas`, `simbolos`, `ast`, `dot`) es LA frontera entre servidor y cliente — insistí en que el cliente NO sabe nada de Jison, de las dos pasadas, ni de las señales de control; solo conoce la forma de este JSON. Es el mismo tipo de desacoplamiento que un compilador real y su IDE.
- "Entorno fresco por petición" ya lo vieron en DataForge, pero acá agregá el matiz nuevo: en DataForge es solo por correctitud de reportes (enunciado §6). Acá ADEMÁS es por concurrencia real — un servidor puede recibir dos peticiones simultáneas (dos pestañas, dos estudiantes), y si el `Interprete` fuera compartido, se contaminarían entre sí.
- CORS: explicá en una frase por qué existe — cliente (`:5173`) y servidor (`:4000`) son orígenes distintos para el navegador, y sin `cors()` el navegador bloquearía el `fetch` aunque ambos corran en la misma máquina.

**Demo en vivo:** con la Terminal servidor ya corriendo (`npm start` en `server/`), mostrá una petición real con curl o con las herramientas de red del navegador contra el cliente ya levantado. Si preferís algo más simple sin salir de la terminal:
```
curl -X POST http://localhost:4000/interpretar -H "Content-Type: application/json" -d "{\"codigo\":\"let x:int=1; ejecutar main(); function void main(){echo x;}\"}"
```
(en PowerShell, las comillas dobles internas necesitan escape con `` ` `` o usar comillas simples según la shell — probalo antes de clase). Alternativa más segura: simplemente abrí las DevTools del navegador (pestaña Network) mientras usás el cliente en la Etapa 6 y mostrá ahí la petición/respuesta real — es más ilustrativo y no depende de sintaxis de curl en vivo.

**Preguntas probables:**
- *"¿Por qué `server.js` no tiene casi lógica propia?"* — Separa protocolo de transporte (HTTP/JSON/CORS) de análisis y ejecución. Esa separación es la que permite probar `analizar()` directo desde consola, como hicimos en las etapas anteriores.
- *"¿Qué pasaría si se reutilizara el mismo `Interprete` entre peticiones?"* — Variables, funciones y errores de una ejecución anterior seguirían presentes en la siguiente — el reporte de símbolos mostraría cosas que ya no existen en el código actual.
- *"¿Para qué sirve `GET /salud` si el cliente nunca lo usa en uso normal?"* — Lo usa la suite Playwright (Etapa 7) para esperar a que el servidor termine de arrancar antes de interactuar con el cliente — evita condiciones de carrera en las pruebas automatizadas.

**Transición:** "El servidor ya responde por contrato. Falta la otra mitad: quién consume ese JSON desde el navegador — el cliente React."

---

## Etapa 6 — El cliente React: editor multipestaña, reportes y el AST con vis-network

**Objetivo pedagógico:** que vean la app funcionando de verdad y entiendan `forwardRef`+`useImperativeHandle` como un patrón puntual (no algo que hay que memorizar en detalle si no vieron React antes).

**Tiempo sugerido:** 12 minutos — si vas corto de tiempo, esta es la etapa donde más podés resumir sin perder el argumento del curso (los detalles de React son menos "teoría de compiladores" que las etapas anteriores).

**Puntos clave a enfatizar oralmente:**
- No te enredes explicando React a fondo si la audiencia no lo conoce — el punto de compiladores acá es otro: el cliente **no reconstruye el AST**, solo dibuja la estructura `{nodes, edges}` que el servidor ya le mandó armada. Esa es la idea que importa para el curso, no los detalles de hooks.
- El gutter de líneas como `<div>` aparte sincronizado a mano con el `scrollTop` del textarea es un buen ejemplo de "HTML no tiene esto nativo, hay que armarlo" — mencionalo brevemente, no hace falta profundizar en el cálculo de posición.
- El cambio automático de pestaña ("Errores" si hubo alguno, "Consola" si no) es una decisión de UX chica pero que después reaparece verificada literalmente por los primeros dos tests de Playwright — es un buen gancho hacia la Etapa 7.
- `physics: false` en el grafo del AST: la razón es simple y vale la pena decirla — un AST ya es jerárquico por definición, así que simular física de repulsión/atracción sería una animación de acomodo innecesaria.

**Demo en vivo — esta es LA demo central del curso, dale tiempo:** con el cliente abierto en `http://localhost:5173`:
1. Mostrá el archivo inicial (el ejemplo del Anexo 11.1) y hacé clic en ▶ Ejecutar. Señalá que el panel salta solo a "Consola" porque no hubo errores.
2. Pegá el contenido de `entradas/ejemplo_errores.ci` en una pestaña nueva y ejecutá — el panel salta solo a "Errores", mostrá la tabla de 8 errores, y hacé clic en una fila para que la clase vea el salto automático a la línea exacta en el editor (`irALinea`).
3. Andá a la pestaña "AST" y mostrá el grafo interactivo con `vis-network` — dejá que alguien lo mueva/haga zoom si el proyector lo permite.
4. Si hay tiempo, mostrá Nuevo/Abrir/Guardar — remarcá que esto NO toca ningún sistema de archivos del servidor, son APIs del propio navegador (`FileReader`, `Blob` + `<a download>`).

**Preguntas probables:**
- *"¿Por qué usar `forwardRef` en vez de pasar la línea como una prop más?"* — Porque "ir a una línea" es una acción puntual (mover cursor, hacer scroll), no un cambio de estado que deba re-renderizar todo. Exponer un método vía `ref` permite dispararlo justo cuando ocurre el clic, sin modelarlo como parte del flujo normal de props/estado.
- *"¿El cliente arma el árbol del AST él mismo?"* — No, recibe `{nodes, edges}` ya calculado por el servidor en la respuesta JSON. El cliente solo dibuja con `vis-network`.
- *"¿Por qué no simula física el grafo del AST, como cualquier grafo genérico?"* — Porque un AST ya es una jerarquía por definición; el layout jerárquico ya calcula la disposición de arriba hacia abajo, simular física sería redundante.

**Transición:** "Ya vimos todo el sistema funcionando en vivo con nuestros propios clics. La pregunta que sigue es: ¿cómo se prueba esto de forma automática, sin que alguien tenga que hacer clic a mano cada vez?"

---

## Etapa 7 — Pruebas de extremo a extremo con Playwright

**Objetivo pedagógico:** distinguir "probar el motor aislado" (lo que veníamos haciendo con `test-cli.js`) de "probar el sistema completo integrado" — y por qué ambas formas de prueba son complementarias, no redundantes.

**Tiempo sugerido:** 10 minutos.

**Puntos clave a enfatizar oralmente:**
- Contrastá explícitamente con lo que ya hicieron juntos en clase: `test-cli.js` (Etapas 1-4) prueba solo el motor, sin servidor ni interfaz. Playwright levanta servidor Express real + cliente Vite real + navegador real, y valida la integración completa — que el fetch llegue, que el JSON se parsee, que el DOM refleje el resultado.
- Los 3 tests son: (1) camino feliz con el ejemplo del Anexo, verificando los 4 reportes; (2) los 3 tipos de error con `ejemplo_errores.ci` — reutilizando el mismo caso ya verificado manualmente en la Etapa 1, y confirmando además el clic-a-línea; (3) crear un archivo nuevo. Remarcá que el test 2 reutiliza el resultado ya conocido (8 errores: 1+2+5) en vez de inventar un caso nuevo — es más valioso porque no introduce una fuente de verdad sin auditar.
- `GET /salud` (Etapa 5) reaparece acá con su propósito real: Playwright lo usa para saber cuándo el servidor terminó de arrancar antes de dejar correr los tests — sin condiciones de carrera.

**Demo en vivo:**
```
cd CompInterpreter/client
npx playwright test
```
Vas a ver la salida `3 passed` en la terminal. Si el proyector lo permite, corré con `npx playwright test --headed` para que la clase vea el navegador abriéndose y ejecutando los clics automáticamente — es más impactante que solo ver texto verde en la terminal.

**Preguntas probables:**
- *"¿Qué diferencia hay entre correr `test-cli.js` y correr Playwright?"* — `test-cli.js` prueba solo el motor de análisis/ejecución, sin servidor HTTP ni interfaz. Playwright levanta el sistema completo y verifica la integración real entre todas las piezas.
- *"¿Por qué reutilizar `ejemplo_errores.ci` en vez de un caso nuevo para el test automatizado?"* — Porque ya se conoce el resultado exacto esperado de haberlo corrido directo contra el motor; reutilizarlo confirma que ese resultado, ya auditado, también llega bien hasta la interfaz, sin introducir una fuente de verdad nueva sin verificar.
- *"¿Por qué `GET /salud` importa específicamente para esta suite?"* — Porque Playwright arranca servidor y cliente como parte de la propia corrida de tests, y necesita saber cuándo el servidor ya escucha antes de dejar avanzar las pruebas.

**Transición hacia la profundización:** "Con esto se cierran las 8 etapas del proyecto: léxico, sintáctico, semántico, servidor, cliente y pruebas. Antes de cerrar del todo, vale la pena profundizar en LA decisión arquitectónica que más distingue a este proyecto de los otros tres del curso."

---

## Profundización ★ — Arquitectura cliente-servidor REST

**Esta profundización es la que mejor resume "qué hace diferente a CompInterpreter" — no la trates como opcional/relleno.**

**Objetivo pedagógico:** consolidar, en una sola narrativa, por qué la frontera HTTP obliga a dos decisiones de diseño (entorno fresco por concurrencia, CORS) que ningún otro proyecto del curso necesita.

**Tiempo sugerido:** 10 minutos.

**Puntos clave a enfatizar oralmente:**
- Arrancá con la tabla monolito-vs-cliente-servidor, pero no la leas entera — el punto que hay que dejar grabado es uno solo: en un monolito la UI llama directo a métodos (misma memoria, mismo proceso); acá la UI cruza una frontera HTTP y solo comparte un contrato JSON. El servidor **podría** correr en otra máquina sin que el cliente cambie nada — eso es imposible en JavaFX.
- Usá la demo interactiva de la propia diapositiva (el botón "▶ Paso" que arma el ciclo completo petición→respuesta en 8 pasos) — no lo expliques vos de memoria, dejá que la clase vea el paso a paso mientras vos narrás cada uno. Es la pieza más lograda de esta profundización, aprovechala.
- Las dos consecuencias de diseño (entorno fresco por concurrencia, CORS) ya las tocaste en la Etapa 5 — acá el valor agregado es conectar el **por qué architectónico**: no son detalles sueltos de implementación, son consecuencias directas y necesarias de haber elegido cliente-servidor.
- Cerrá con el paralelo al compilador real y su IDE: ninguno de los dos lados necesita conocer los detalles internos del otro, solo el formato de mensaje acordado. Es una buena forma de conectar esto con el mundo profesional fuera del curso.

**Demo en vivo:** ninguna adicional — la demo ES la diapositiva interactiva (botón "▶ Paso" / "↺ Reiniciar"). Practicá el timing: dar clic, dejar que se lea el paso, narrar el detalle, siguiente clic. No lo aceleres.

**Preguntas probables:**
- *"¿Por qué no llamar directo a las funciones del intérprete desde el cliente, como en DataForge?"* — Porque el intérprete corre sobre Node.js (usa el parser de Jison, una librería de Node) y el cliente corre en el navegador — son runtimes distintos, el navegador no ejecuta módulos de Node directamente. Además el enunciado exige la arquitectura cliente-servidor como objetivo, no como opción.
- *"¿Qué pasaría si el servidor no creara un intérprete nuevo por petición, específicamente en este contexto?"* — Dos peticiones concurrentes (dos pestañas, dos usuarios) se pisarían el estado mutuamente — algo que en un monolito de escritorio de un solo usuario simplemente no puede pasar.
- *"¿En qué se parece esto a un compilador real con su IDE?"* — En que coordinan a través de un formato de mensaje acordado sin que ninguno conozca los detalles internos del otro — el cliente no sabe que hay dos pasadas del AST ni cómo se implementan las señales de control.

**Transición al cierre:** "Con esto completamos el recorrido de los 4 proyectos de OLC1 — y este es, de los cuatro, el que más se parece en arquitectura a una aplicación web real del mundo profesional."

---

## Cierre

**Síntesis final (2-3 minutos):** repasá en una frase por etapa, sin releer las diapositivas: léxico en un archivo (Jison), AST con fábrica de nodos genérica, intérprete de dos pasadas con propagación por null, señales de control auditadas con 2 bugs reales corregidos, servidor REST con entorno fresco por concurrencia, cliente React que solo dibuja lo que el servidor ya calculó, y todo verificado de punta a punta con Playwright.

**El mensaje que tenés que dejar bien claro al cerrar:** la teoría de compiladores no cambió ni una vez en las 8 etapas — léxico, sintáctico, semántico, tabla de símbolos, propagación de errores son exactamente los mismos conceptos del Dragón que vieron en DataForge. Lo que cambió fue la arquitectura (cliente-servidor en vez de monolito) y, como consecuencia directa de esa arquitectura, dos decisiones que ningún otro proyecto del curso necesitó: entorno fresco por razones de concurrencia (no solo de reportes) y CORS.

**Mencioná que es el último de los 4 proyectos de OLC1.** DataForge (construido con demostración guiada), ConjAnalyzer y CompScript (con los roles invertidos: el estudiante escribe, el instructor tutorea) y CompInterpreter cierran el recorrido completo del curso.

**Si el público también va a ver VLangCherry (OLC2):** es un buen momento para tender el puente. VLangCherry es el salto de "intérprete" a un proyecto con generación de código real (fuera del alcance de OLC1) — mencionalo como el "próximo nivel" sin entrar en detalle si no corresponde a esta clase, simplemente para que la audiencia entienda que este no es el techo, es la base sobre la que se construye lo que viene.

---
