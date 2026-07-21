# Guion de clase — DataForge

> Nota para vos (quien da la clase): esto NO es la presentación. Es tu chuleta de mano — qué decir en voz alta, qué correr en vivo, qué responder si preguntan. Las diapositivas (`index.html` + 13 páginas) están escritas en 3ª persona para la audiencia; este documento está en 2ª persona, para vos.

---

## Antes de empezar

**Cómo abrís la presentación**: doble clic en `C:\Users\72358\Desktop\Hades\presentacion-dataforge\index.html`. Es autocontenido (CSS y JS embebidos en `assets/`), no necesita servidor ni internet — el navegador la abre directo desde `file://`. Navegás con las flechas ← → dentro de cada etapa, y el botón ◐ cambia el tema claro/oscuro (probalo antes de la clase para saber en qué modo lo vas a dejar según la luz del salón/proyector).

**Duración total estimada**: contando las 7 secciones del roadmap principal (Etapa 0 a 6) MÁS las 6 profundizaciones ★, esto es material para **dos sesiones de clase, no una**. Estimación por bloque:

| Bloque | Tiempo sugerido |
|---|---|
| Etapa 0 (conceptos base) | 15 min |
| Etapa 1 (léxico) + demo en vivo | 25 min |
| Etapa 2 (sintáctico) | 20 min |
| Etapa 3 (ejecución) + demo en vivo | 25 min |
| Etapa 4 (editor JavaFX) + demo en vivo | 15 min |
| Etapa 5 (gráficas) + demo en vivo | 15 min |
| Etapa 6 (reportes) + demo en vivo | 20 min |
| **Subtotal roadmap principal** | **~135 min (2h15)** |
| gramaticas.html (la más pesada, 7 demos) | 35 min |
| automatas.html (4 demos con SVG) | 25 min |
| tabla-simbolos.html | 15 min |
| ast.html | 15 min |
| fases.html (sin demos, más corta) | 10 min |
| semantica.html (sin demos, más corta) | 15 min |
| **Subtotal profundizaciones** | **~115 min** |

Si tenés una sola sesión de 2 horas, cubrí Etapa 0 a 6 (el roadmap principal) y dejá las 6 profundizaciones ★ para una segunda sesión de "teoría a fondo" — de hecho así están pensadas (el roadmap ya las separa en dos secciones distintas en el índice).

**Qué preparar de antemano**:

1. Tené **DataForge abierto en IntelliJ IDEA** desde antes de que empiece la clase (el proyecto en `C:\Users\72358\Desktop\Hades\DataForge\`), con el proyecto ya indexado (Maven ya resuelto) para que las demos en vivo no se traben con "downloading dependencies" en pleno directo.
2. Verificá que compila UNA vez antes de la clase: en IDEA, ventana Maven → Lifecycle → `compile`. Si vas a usar terminal en vez del botón ▶ Run de IDEA, el JDK del PATH es 1.8 y NO sirve — usá el JDK embebido que gestiona el IDE con el comando exacto de abajo.
3. **Comando exacto para correr un ejemplo por consola** (Git Bash, con el Maven embebido de IDEA y el JDK 25 que tenés instalado vía IDEA):
   ```bash
   JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemploN.df"
   ```
   Reemplazá `ejemploN.df` por el archivo que quieras demostrar (ver tabla abajo). Corré este comando parado en `C:\Users\72358\Desktop\Hades\DataForge\`.
4. Si preferís no usar terminal: en IDEA, ▶ Run sobre la clase `TestInterprete` (podés cambiar el argumento en la configuración de Run/Debug para apuntar a otro `.df`).
5. Para la demo de Etapa 4 (editor) y Etapa 5 (gráficas): ▶ Run sobre `Lanzador` — **nunca sobre `EditorApp`** directamente, porque da el error "JavaFX runtime components are missing" (es justo el ejemplo que vas a explicar en etapa4.html, así que si te pasa en vivo, mejor — es la demostración perfecta del problema).
6. Los 4 archivos de prueba y qué ejercitan, para saber cuál correr en cada momento:

| Archivo | Qué demuestra | Útil en |
|---|---|---|
| `entradas/ejemplo1.df` | Caso básico: comentarios (línea y multilínea), declaración `double`/`char[]`, arreglo con aritmética embebida (`SUM(3,4)` dentro de la lista), función estadística `Media`, `print` y `column` | Etapa 1 (tabla de tokens), Etapa 3 (ejecución básica) |
| `entradas/ejemplo2.df` | Anidación profunda (`DIV(SUM(Max(@notas),Min(@notas)),2)`) y las 4 gráficas completas (bar/pie/line/histogram), incluido el typo `grapBar`/`grapPie`/`grapLine` (sin h) que el lexer acepta a propósito | Etapa 2 (recursión de la gramática), Etapa 5 (gráficas) |
| `entradas/ejemplo3_errores.df` | 5 errores SEMÁNTICOS a propósito (no declarada, tipos incompatibles, división entre cero, asignar cadena a double, redeclaración) — léxica y sintácticamente perfecto, y la última línea SÍ imprime pese a los 5 errores previos | Etapa 3 (propagación por null) |
| `entradas/ejemplo4_mixto.df` | Los 3 tipos de error en un solo programa: léxico (`$` suelto), sintáctico (falta `<-`, dispara modo pánico), semántico (`fantasma` no declarada) — y demuestra que `precio` sí queda declarada pese al error léxico pero `roto` NO por el modo pánico | Etapa 6 (reportes + modo pánico) — es el ejemplo estrella de esa etapa |

---

## Roadmap principal (Etapa 0 a 6)

### Etapa 0 · Conceptos base

**Objetivo pedagógico**: instalar la idea de que un intérprete lee texto plano y lo convierte en 3 pasos mentales (agrupar, ordenar, ejecutar) — el mismo proceso que usamos leyendo una oración en español.

**Tiempo sugerido**: 15 min.

**Puntos clave a enfatizar oralmente**:
- Arrancá mostrando la línea de caracteres sueltos (`v a r : d o u b l e...`) — es tu gancho visual: "esto es TODO lo que Java ve al principio, ni una palabra, ni un espacio con significado".
- La analogía de «El gato negro duerme» es la columna vertebral de la etapa — dedicale tiempo, no la apures. Los 3 pasos (agrupar letras en palabras / verificar orden / entender y reaccionar) son literalmente léxico/sintáctico/ejecución, y de ahí en adelante vas a reusar esos 3 nombres toda la clase.
- Remarcá la tabla de las 3 capas con la herramienta de cada una (JFlex/CUP/Java) — es el mapa que van a necesitar para ubicarse en cada etapa siguiente.
- No te saltes la anatomía de la línea real (`var:double:: numero <- 2.5 end;` con 9 tokens coloreados) — es la primera vez que ven "esto de acá es un token" sobre código real de DataForge, no sobre español.

**Demo en vivo**: no aplica todavía (es la única etapa 100% conceptual, sin código corriendo).

**Preguntas probables y cómo responderlas**:
- *"¿Por qué `cadena` es identificador y no palabra reservada?"* (pregunta real del quiz) — porque es un nombre que el usuario inventó; las palabras reservadas son una lista fija y cerrada que el lenguaje ya trae (como `var`, `end`). Si el usuario pudo escribir cualquier cosa ahí, es identificador.
- *"¿Por qué `"hola"` y `"adios"` son del mismo token?"* (pregunta real del quiz) — porque el token es la CATEGORÍA (cadena), no el contenido. Un solo patrón (`comillas ... contenido ... comillas`) reconoce infinitos textos distintos.
- Duda típica que no está en el quiz: *"¿Por qué `var` y `VAR` valen lo mismo?"* — esto SÍ está como pregunta bonus del quiz, respondé que se resuelve en la capa léxica (JFlex con `%ignorecase`), no en la sintáctica — el orden de las palabras no cambia, solo cómo reconocemos la palabra en sí.
- Duda típica adicional: *"¿Por qué separar en capas si al final todo es un solo programa Java?"* — porque cada capa resuelve un tipo de pregunta distinto y usa una herramienta distinta (JFlex sabe de patrones, CUP sabe de gramáticas, Java sabe de ejecutar) — separar es lo que hace manejable escribir un lenguaje nuevo.

**Transición**: "Ya vimos que existen tokens. Ahora toca escribir, de verdad, la lista completa de tokens de DataForge — eso es la Etapa 1."

---

### Etapa 1 · Análisis léxico

**Objetivo pedagógico**: construir la tabla de tokens completa de DataForge recorriendo el enunciado sección por sección, y ver el `Lexer.flex` real generando el analizador.

**Tiempo sugerido**: 25 min (incluye demo en vivo).

**Puntos clave a enfatizar oralmente**:
- El método importa tanto como el resultado: mostrá que fuiste sección por sección del enunciado (5.1 a 5.10) preguntando "¿qué palabra nueva aparece acá?" — es una técnica que quieren copiar para sus propios proyectos (ConjAnalyzer, etc.).
- De las 4 decisiones de diseño (slide 5), la que más engancha en clase es el typo `graphBar`/`grapBar` — mostralo como ejemplo real de "el enunciado tiene un error tipográfico sistemático y el lexer lo tiene que tolerar a propósito, no por descuido".
- Cuando llegues al fragmento real de `Lexer.flex`, remarcá el ORDEN: reservadas antes de `{Id}`. Esto es un error clásico que van a cometer ellos mismos — enfatizalo fuerte, con el ejemplo de qué pasaría si invirtieran el orden (pregunta 1 del quiz parte 2).
- No dejes pasar el callout de "match más largo" para `::` — es la primera vez que aparece la palabra "AFD" indirectamente (se profundiza en `automatas.html`, podés adelantar que van a volver a este tema).

**Demo en vivo**: corré el ejemplo1.df por consola para mostrar la tabla de tokens real. Ojo: el `mainClass` por defecto del `pom.xml` es `TestInterprete` (que NO imprime tokens, solo consola/símbolos/errores) — para esta demo hay que apuntar explícitamente a `TestLexer`:
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.mainClass=dataforge.TestLexer -Dexec.args="entradas/ejemplo1.df"
```
Alternativa igual de válida: abrir `TestLexer.java` en IDEA y correrlo con ▶ Run — es fiel a como se construyó originalmente (sin parser todavía). Cualquiera de las dos formas sirve; la primera es más rápida de preparar si ya tenés la terminal abierta.

**Preguntas probables y cómo responderlas**:
- *"¿Cuántos tokens produce `console::print = "hola", 15 end;`?"* (pregunta real del quiz) — 9: `console` `::` `print` `=` `"hola"` `,` `15` `end` `;`. Contalos en pizarra/pantalla uno por uno, es más convincente que decir el número solo.
- *"¿Por qué `@datos` es un solo token y no dos (`@` + `datos`)?"* (pregunta real del quiz) — decisión de diseño: se definió `ID_ARREGLO` como patrón completo (`@` seguido de identificador) en vez de tratar `@` como símbolo aparte. Ambas opciones eran válidas; se eligió la más simple para el parser de la Etapa 2.
- Duda típica que no está en el quiz: *"¿Por qué no usamos simplemente `String.split()` de Java en vez de JFlex?"* — porque `split` no maneja bien casos como cadenas con comas adentro, comentarios anidados, o el match más largo (`::` vs `:` `:`); JFlex genera un autómata que resuelve todo eso de forma consistente y es el estándar de la industria (Lex/Flex/JFlex son la misma familia de herramienta hace 40+ años).
- Si preguntan por el troubleshooting de IntelliJ (`Cannot resolve symbol 'Lexer'`) — es porque falta correr Maven → Lifecycle → compile antes; el Lexer.java se GENERA, no existe hasta que se corre esa fase.

**Transición**: "Ya sabemos separar el texto en tokens. El problema es que una lista plana de tokens no nos dice si están en el ORDEN correcto — eso es gramática, Etapa 2."

---

### Etapa 2 · Análisis sintáctico

**Objetivo pedagógico**: entender por qué necesitamos una gramática libre de contexto (no regex) para validar la estructura, y ver la BNF completa de DataForge traducida a CUP.

**Tiempo sugerido**: 20 min.

**Puntos clave a enfatizar oralmente**:
- Arrancá con el ejemplo de paréntesis anidados (`DIV(SUM(Max(@notas),Min(@notas)),2)`) — es el argumento MÁS fuerte de por qué regex no alcanza (necesita "contar" profundidad arbitraria, y regex no tiene memoria para eso). Citá Dragón §4.1 acá, es el momento correcto.
- La derivación paso a paso de `var:double:: numero <- 2.5 end;` conviene proyectarla y leerla en voz alta flecha por flecha — no la resumas, es la primera vez que ven una derivación completa sobre SU propio lenguaje.
- Mostrá la BNF completa (las dos slides, núcleo + expresiones/gráficas) pero no leas cada línea — señalá los patrones repetidos (todo termina en `"end" ";"`, las gráficas anidan atributos igual que las instrucciones anidan). Esto les ahorra memorizar 15 producciones distintas.
- El fragmento real del `.cup` con `terminal`/`non terminal`/`start with` es el momento de decir en voz alta la convención del proyecto: "TODAS las declaraciones `non terminal` van antes que cualquier producción — si no, CUP tira error de compilación".

**Demo en vivo**: sí existe un ejecutable de "solo parser": `TestParser.java` (valida contra la gramática sin ejecutar nada — imprime `[OK]` o `[X]`). Corré:
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.mainClass=dataforge.TestParser -Dexec.args="entradas/ejemplo4_mixto.df"
```
`ejemplo4_mixto.df` ya trae el error exacto que conviene mostrar (línea `var:double:: roto 5 end;`, sin `<-`) — mostralo ahora aunque sea antes de tiempo, y anunciá que van a volver a este mismo archivo en la Etapa 6 para ver el reporte completo. **Ojo con el resultado**: por el modo pánico (`instruccion ::= error PUNTO_COMA`), `TestParser` de todos modos imprime `[OK] Análisis sintáctico exitoso` al final — el error se ve solo en el mensaje de `syntax_error()` impreso por stderr durante el análisis, no en el veredicto final. Es un buen gancho: "OK" no significa "sin errores", significa "el análisis terminó sin abortar" — la recuperación en modo pánico es justo la razón por la que no revienta.

**Preguntas probables y cómo responderlas**:
- *"`console::print = "hola" 15 end;` (falta la coma). ¿Es error léxico o sintáctico?"* (pregunta real del quiz) — sintáctico: todos los tokens son válidos individualmente (`"hola"` y `15` están bien formados), lo que falta es el orden esperado por la gramática (falta la COMA entre expresiones).
- *"¿Qué hace posible anidar `SUM` cualquier cantidad de veces?"* (pregunta real del quiz) — recursión mutua entre `<aritmetica>` y `<expresion>`: una aritmética contiene expresiones, y una expresión puede SER una aritmética. Esa mutua-referencia es lo que permite profundidad infinita sin escribir reglas extra.
- Duda típica que no está en el quiz: *"¿Por qué la gramática no valida también que `numero` sea double y no char[]? ¿Eso no debería ir acá?"* — buena entrada para adelantar la Etapa 3: la gramática solo valida FORMA (orden de piezas), no significado (tipos, valores). Eso es trabajo semántico, próxima etapa.

**Transición**: "La gramática ya nos dice si el programa está bien formado. Pero 'bien formado' no es lo mismo que 'funciona' — hay que hacer que ejecute de verdad. Etapa 3."

---

### Etapa 3 · Ejecución

**Objetivo pedagógico**: mostrar cómo las acciones semánticas de CUP (`{: RESULT = ... :}`) ejecutan código Java real en el momento en que el parser reduce, sin necesidad de un árbol intermedio.

**Tiempo sugerido**: 25 min (incluye demo en vivo).

**Puntos clave a enfatizar oralmente**:
- La analogía de cajas anidadas (`SUM(SUM(1,2),3)`, la interna se cierra primero) es el corazón teórico de toda la clase — de acá sale la decisión de "sin AST" que van a defender en la calificación. No la apures.
- Cuando muestres el fragmento real de CUP (`RESULT = Operaciones.aritmetica(op, a, b, parser.entorno, opleft, opright)`), señalá que `a` y `b` YA vienen calculados — es la prueba concreta de que gramática S-atribuida + parser LR (ascendente) es la pareja perfecta para ejecutar sin AST.
- La tabla del Entorno (tabla de símbolos, consola, errores, gráficas) conviene dibujarla en pizarra o remarcarla fuerte: es la clase que van a reusar mentalmente en CADA etapa siguiente.
- Las dos reglas de los errores semánticos — "no detiene ejecución" y "propagación por null sin cascadas" — son la regla de oro del proyecto entero. Repetilas explícitamente, se preguntan seguido.

**Demo en vivo**: corré `ejemplo3_errores.df` para mostrar los 5 errores semánticos y que la última línea SÍ imprime pese a todo:
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo3_errores.df"
```
Después corré `ejemplo1.df` para mostrar la salida "feliz" completa (consola con la media, columna de datos, tabla de símbolos con 4 entradas) — es la que aparece literal en la diapositiva, así que confirma que lo que ven en pantalla coincide con lo prometido.

**Preguntas probables y cómo responderlas**:
- *"`console::print = fantasma end;` sin declarar antes. ¿Qué capa detecta esto?"* (pregunta real del quiz) — semántica, en tiempo de ejecución, vía `valorDe("fantasma")` que no la encuentra en la tabla de símbolos.
- *"¿Por qué la acción de `aritmetica` puede confiar en que `a` y `b` ya están calculados?"* (pregunta real del quiz) — porque el parser LR reduce las cajas internas ANTES que las externas; cuando se ejecuta la acción de la caja externa, la interna ya produjo su valor.
- *"¿Por qué devolver `null` en vez de lanzar una excepción de Java?"* (pregunta real del quiz) — una excepción abortaría todo el análisis en el primer error; el enunciado (§6.2) exige reportar TODOS los errores de una corrida, no solo el primero.
- Duda típica que no está en el quiz: *"¿Por qué la tabla de símbolos usa `LinkedHashMap` y no un `HashMap` normal?"* — porque el reporte de símbolos (§6.3) necesita mantener el ORDEN de declaración, y `LinkedHashMap` te da hash O(1) + orden de inserción en una sola estructura.

**Transición**: "Ya tenemos un intérprete que funciona por consola. Pero el enunciado pide un editor de verdad, con pestañas y botón de ejecutar — eso es interfaz, Etapa 4."

---

### Etapa 4 · Editor JavaFX

**Objetivo pedagógico**: ver cómo se conecta el intérprete de consola (Etapa 3) a una interfaz gráfica real, sin tocar ni una línea del lexer/parser.

**Tiempo sugerido**: 15 min (incluye demo en vivo).

**Puntos clave a enfatizar oralmente**:
- Advertí desde el arranque: "acá no hay Dragón — esto es ingeniería de UI, no teoría de compiladores". Baja la ansiedad de quienes esperan más BNF.
- El árbol ASCII de Stage→Scene→Nodos conviene dibujarlo rápido en pizarra mientras señalás la ventana real de DataForge abierta — conecta lo abstracto con lo que están viendo.
- El punto más importante técnicamente es `StringReader` vs `FileReader`: remarcá que el Lexer nunca cambió, porque estaba programado contra la abstracción `Reader` desde el principio (una decisión de diseño de la Etapa 1 que recién ahora rinde fruto). Es un buen ejemplo de "por qué programar contra interfaces".
- El error real de "JavaFX runtime components are missing" y el truco de `Lanzador` (una clase de una sola línea que NO extiende `Application`) es la anécdota más memorable de esta etapa — contala como historia real de debugging, no como teoría.

**Demo en vivo**: abrí el editor completo:
- En IDEA: ▶ Run sobre `Lanzador` (NUNCA sobre `EditorApp` directamente — si alguien pregunta "¿y si corro EditorApp?", corré `EditorApp` a propósito para que vean el error real en pantalla, y después corré `Lanzador` para mostrar que sí funciona).
- Terminal equivalente:
  ```bash
  JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" clean javafx:run
  ```
- Una vez abierto, abrí `entradas/ejemplo1.df` con el FileChooser del editor, dale a Ejecutar, y mostrá la consola de la GUI llenándose en vivo.
- Los WARNING de JavaFX en la consola (unnamed module, native access, `sun.misc.Unsafe`) son esperados — adelantá que los van a ver y que no son error.

**Preguntas probables y cómo responderlas**:
- *"¿Por qué bastó con `StringReader` en vez de `FileReader`?"* (pregunta real del quiz) — el Lexer estaba escrito contra la interfaz `Reader`, no contra una implementación concreta; cualquier fuente de texto (archivo o String en memoria) sirve igual.
- *"¿Qué pasaría si reutilizáramos el MISMO Entorno entre clics de Ejecutar?"* (pregunta real del quiz) — el estado (variables, errores, gráficas) se iría acumulando entre corridas, violando el requisito §6 de que los reportes muestren SOLO el último análisis. Por eso cada clic crea un Entorno nuevo.
- *"¿Para qué existe `Lanzador.java` si tiene una sola línea?"* (pregunta real del quiz) — esquiva el chequeo del launcher estándar de Java, que exige que la clase main extienda `Application` directamente cuando corrés JavaFX sin el plugin de módulos configurado.
- Duda típica que no está en el quiz: *"¿Por qué no usaron FXML como en los tutoriales?"* — decisión pedagógica documentada: construir el scene graph por código deja visible la jerarquía Stage→Scene→Nodos, que es justo lo que se quiere enseñar acá.

**Transición**: "El editor ya ejecuta y muestra consola. Falta la parte visual de verdad: convertir números en gráficas. Etapa 5."

---

### Etapa 5 · Gráficas

**Objetivo pedagógico**: entender por qué el renderizado ocurre DESPUÉS de que termina toda la ejecución (no durante), y ver el mapeo directo entre las 4 instrucciones de gráfica del lenguaje y las clases de JavaFX Charts.

**Tiempo sugerido**: 15 min (incluye demo en vivo).

**Puntos clave a enfatizar oralmente**:
- El "por qué se dibuja DESPUÉS" es el punto conceptual más importante de la etapa: el intérprete (capa de lenguaje) no sabe que existe una GUI — solo acumula objetos `Grafica` en el Entorno; recién al final, la capa de GUI los recorre y los dibuja. Esto es la misma separación de capas que vienen viendo desde la Etapa 0, aplicada a gráficas.
- La tabla de mapeo (graphBar→BarChart, graphPie→PieChart, graphLine→LineChart, Histogram→BarChart de frecuencias) conviene mostrarla completa de un vistazo — es simple y visual, no necesita mucha explicación oral.
- El histograma merece un comentario aparte: el formato de la tabla de frecuencias (frecuencia/acumulada/relativa) NO estaba especificado en el PDF original (una imagen que no se pudo convertir a texto) — es una decisión propia documentada. Es un buen ejemplo real de "qué hacer cuando el enunciado tiene un hueco".
- Remarcá que NO hizo falta agregar dependencia nueva al `pom.xml` para esto — JavaFX Charts ya viene incluido en `javafx-controls`, que se agregó en la Etapa 4.

**Demo en vivo**: con el editor ya abierto (Lanzador corriendo), abrí `entradas/ejemplo2.df` y ejecutalo — se deberían abrir 4 ventanas nuevas (bar, pie, line, histograma). Si preferís consola en vez de GUI:
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo2.df"
```
(por consola no vas a ver las ventanas de gráfica — para eso necesitás correr el editor JavaFX. Si el objetivo es mostrar las 4 gráficas, usá el editor.)

**Preguntas probables y cómo responderlas**:
- *"¿Por qué la tabla de frecuencias del histograma la calcula el Entorno y no el Graficador?"* (pregunta real del quiz) — es semántica del LENGUAJE (debería poder calcularse aunque no hubiera GUI, por ejemplo si algún día se pidiera solo por consola), no un detalle de dibujo.
- *"`graphPie` con 3 labels pero solo 2 values: ¿quién lo rechaza?"* (pregunta real del quiz) — el Entorno, en `validarGrafica` — es un error semántico (tamaños incoherentes), detectado antes de intentar dibujar nada.
- *"¿Por qué no hizo falta ninguna dependencia nueva en el pom?"* (pregunta real del quiz) — JavaFX Charts vive dentro del mismo artefacto `javafx-controls` que ya se agregó en la Etapa 4 para el editor.
- Duda típica que no está en el quiz: *"¿Por qué el enunciado escribe `values::char[]` en el histograma si son números?"* — es otra inconsistencia real del enunciado (como el typo `grapBar`); se resuelve a favor de tratarlos como doubles, priorizando lo que tiene sentido semántico sobre lo literal del texto.

**Transición**: "Ya vimos, oímos y ejecutamos DataForge. Lo único que falta es dejar constancia en papel (bueno, en HTML) de todo lo que pasó — los reportes. Etapa 6."

---

### Etapa 6 · Reportes HTML

**Objetivo pedagógico**: cerrar el ciclo completo del enunciado mostrando los 3 reportes obligatorios (tokens, errores, símbolos) y el modo pánico como mecanismo de recuperación sintáctica.

**Tiempo sugerido**: 20 min (incluye demo en vivo — la más rica de todas).

**Puntos clave a enfatizar oralmente**:
- Anunciá desde el arranque que esta es la ETAPA QUE CIERRA el proyecto: "al terminar esto, DataForge cumple todos los requisitos mínimos del enunciado". Da un cierre de tensión narrativa a la clase completa.
- El modo pánico (`instruccion ::= error PUNTO_COMA`) es el concepto teórico más denso de esta etapa — dedicale tiempo. Los 3 pasos (reportar → podar la pila → descartar hasta `;`) conviene explicarlos con el ejemplo real de `ejemplo4_mixto.df` en pantalla mientras hablás, no en abstracto.
- El contraste `precio` (SÍ se declara) vs `roto` (NO se declara) es EL momento clave de toda la etapa — es la pregunta real del quiz y la que más cuesta entender. Armá el contraste bien claro: error léxico descarta 1 carácter (la instrucción sigue viva), error sintáctico en modo pánico descarta la INSTRUCCIÓN ENTERA.
- Mencioná de pasada la reflexión sobre `sym.getFields()` en `Reportes.java` — es un detalle elegante (evita un switch de 40 casos) que vale la pena señalar aunque no se profundice.

**Demo en vivo** (la más completa de la clase): con el editor abierto, cargá `entradas/ejemplo4_mixto.df`, ejecutalo, y abrí los 3 reportes con el botón correspondiente del editor:
- `errores.html` → debería mostrar exactamente 3 errores: 1 léxico (`$`), 1 sintáctico (`<-` faltante), 1 semántico (`fantasma`).
- `tokens.html` → la lista completa de tokens reconocidos.
- `simbolos.html` → debería aparecer `bueno` y `precio`, pero NO `roto` — este es el momento de oro de la demo, señalalo explícitamente.

Alternativa por consola (no genera los HTML pero sí muestra los 3 errores en la salida):
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo4_mixto.df"
```

**Preguntas probables y cómo responderlas**:
- *"¿Por qué `precio` SÍ queda declarada pese al `$` pero `roto` NO pese a que también se recuperó?"* (pregunta real del quiz) — son recuperaciones de capas distintas: el error léxico descarta solo el carácter `$` y la instrucción completa sobrevive; el modo pánico (sintáctico) descarta la instrucción ENTERA donde ocurrió el error.
- *"¿Por qué se eligió `;` como punto de sincronización del modo pánico?"* (pregunta real del quiz) — porque toda instrucción de DataForge termina en `;` sin excepción; es un token "de sincronización" confiable (Dragón §4.8.1) — cualquier `;` marca un punto seguro para retomar.
- *"¿Por qué el reporte de tokens necesitó tocar el LEXER y no alcanzaba con el parser?"* (pregunta real del quiz) — el parser consume y DESCARTA los Symbols a medida que reduce; además el Symbol lleva el valor ya procesado (ej. el double convertido), no el lexema original de texto — por eso el lexer necesita un campo `entorno` propio para ir registrando cada token tal como aparece.
- Duda típica que no está en el quiz: *"¿Qué pasa si hay dos errores léxicos seguidos sin `;` entre medio?"* — cada carácter inválido se descarta uno por uno de forma independiente (no hay "modo pánico léxico"); el modo pánico es un mecanismo exclusivamente sintáctico.

**Transición al cierre**: "Con esto, DataForge cumple el enunciado completo. Lo que queda pendiente es solo empaquetado de entrega — no hay más teoría nueva en el proyecto básico. Si tienen tiempo, ahí es donde entran las 6 profundizaciones ★, que retoman cada etapa con más rigor del Dragón."

---

## Profundizaciones ★ (segunda sesión, si aplica)

Estas 6 páginas no siguen el hilo narrativo de "construir DataForge etapa por etapa" — vuelven sobre el MISMO proyecto ya terminado para profundizar en la teoría formal del Dragón. Conviene presentarlas así: "ya vimos que DataForge funciona; ahora vamos a entender POR QUÉ funciona, con el rigor que pide el libro".

### gramaticas.html

**Objetivo pedagógico**: dar el aparato formal completo de gramáticas libres de contexto — derivación, LL vs LR, conflictos, ambigüedad — usando DataForge como ejemplo running.

**Tiempo sugerido**: 35 min (es la más pesada — 7 demos con stepper).

**Puntos clave a enfatizar oralmente**:
- Es la profundización con MÁS demos animadas (7 en total) — no te apures, cada una construye sobre la anterior. El orden importa: derivación en español → derivación DataForge → cajas anidadas → recursión → error sintáctico → shift-reduce → tabla predictiva LL.
- La Demo 6 (shift-reduce simulando `SUM(1,2)` con tabla Pila/Entrada/Acción) es la más importante técnicamente — es la primera vez que ven mecánicamente CÓMO reduce un parser LR, que es la base de todo lo que hace CUP.
- El bloque LL/LR (4 slides intercaladas) es denso — recursión izquierda como "veneno para LL, alimento para LR" es la frase que más se repite en evaluaciones, remarcala.
- Cerrá con el ejemplo de ambigüedad (`10-4-3`, dos árboles posibles) — conecta directo con "por qué CompInterpreter sí va a necesitar declarar `precedence left`" (DataForge no lo necesita porque su gramática no es ambigua).

**Demo en vivo**: no aplica (es 100% conceptual/teórica, sin correr DataForge) — el foco está en las demos SVG/stepper embebidas en el HTML mismo. Andá clic por clic en cada botón "▶ Paso", no los saltees con las flechas del deck.

**Preguntas probables y cómo responderlas**:
- *"¿Qué significa el «(1)» de LL(1)?"* (pregunta real del quiz) — que la decisión de qué producción usar se toma mirando solo 1 token de anticipación, y que eso funciona porque los conjuntos PRIMERO de las alternativas son disjuntos.
- *"¿Por qué `lista_instr ::= lista_instr instruccion` (recursión izquierda) funciona en CUP pero reventaría escrito a mano con descenso recursivo?"* (pregunta real del quiz) — CUP es LALR (bottom-up, no predice desde arriba); un parser LL a mano llamaría a la misma función infinitamente antes de leer un token, ciclando sin consumir entrada.
- *"¿Qué operación del algoritmo shift-reduce corresponde a 'cerrar una caja'?"* (pregunta real del quiz) — reduce.
- Duda típica que no está en el quiz: *"¿Por qué DataForge no eligió ANTLR si es más moderno que CUP?"* — es una decisión de stack del curso (CUP/JFlex forman el par clásico Lex/Yacc en Java), no una limitación técnica; ANTLR es LL y hubiera exigido evitar la recursión izquierda que la gramática de DataForge usa libremente.

**Transición**: "Ya vimos CÓMO se decide la estructura de un programa. Ahora bajemos un nivel más: cómo se reconocen los tokens mismos, con autómatas."

### automatas.html

**Objetivo pedagógico**: mostrar el pipeline completo ER→AFN→AFD→AFD mínimo del capítulo 3 del Dragón, con el token NUMERO de DataForge como hilo conductor.

**Tiempo sugerido**: 25 min (4 demos con diagramas SVG).

**Puntos clave a enfatizar oralmente**:
- Arrancá con la Demo A (AFD de NUMERO consumiendo `57.75` carácter por carácter) — es la más intuitiva y sirve de ancla para todo lo que sigue.
- El ejemplo canónico del libro `(a|b)*abb` (Demo B, Thompson) es célebre — si alguien ya vio el Dragón sabe que es LA figura de referencia del capítulo 3; nombralo así para que conecten con el libro directamente.
- La Demo C (subconjuntos) y Demo D (simulación del AFD resultante) conviene correrlas una después de otra sin pausa larga — son las dos mitades de un mismo argumento (constrí el AFD, después probalo).
- Cerrá fuerte con el dato de la última slide: JFlex fusiona los ~39 patrones del `Lexer.flex` real en UN SOLO autómata gigante — es un dato concreto y sorprendente que ancla toda la teoría abstracta en el proyecto real.

**Demo en vivo**: no aplica (sin código DataForge corriendo) — foco 100% en las 4 demos SVG del HTML.

**Preguntas probables y cómo responderlas**:
- *"¿Qué pasa con `abab`? ¿Y con `bbabb`?"* (pregunta real del quiz, sobre el AFD de `(a|b)*abb`) — `abab` termina en un estado no final (rechazada); `bbabb` sí termina en el estado de aceptación E.
- *"¿Por qué Thompson genera transiciones ε si después hay que eliminarlas con subconjuntos?"* (pregunta real del quiz) — es una división de trabajo deliberada: Thompson construye FÁCIL (mecánico, composicional), subconjuntos optimiza DESPUÉS — separar construcción de optimización simplifica cada paso.
- *"El Lexer.flex tiene unos 39 patrones distintos. ¿Cuántos autómatas separados ejecuta JFlex en tiempo real?"* (pregunta real del quiz) — uno solo; todos los patrones se fusionan en un único AFD.
- Duda típica que no está en el quiz: *"¿Esto lo hace JFlex automáticamente o hay que programarlo?"* — 100% automático; uno escribe expresiones regulares en el `.flex` y JFlex corre Thompson+subconjuntos+minimización por dentro al generar el `Lexer.java`.

**Transición**: "Los autómatas resuelven CÓMO se reconoce un token. Ahora veamos DÓNDE se guardan los nombres que el programa va declarando: la tabla de símbolos."

### tabla-simbolos.html

**Objetivo pedagógico**: mostrar la tabla de símbolos como caso general (con ámbitos anidados) y explicar honestamente que DataForge usa el caso más simple posible (un solo ámbito).

**Tiempo sugerido**: 15 min.

**Puntos clave a enfatizar oralmente**:
- Sé transparente con la comparación honesta (slide 6): DataForge no necesita cadena de ámbitos porque solo tiene un bloque (`PROGRAM...END PROGRAM`). No es una limitación oculta, es correcto para ESTE lenguaje.
- La demo de la cadena de entornos (código con bloque interno que redeclara `x`) es la pieza central — remarcá cómo `get` "sube" por la referencia al padre buscando el nombre.
- Cerrá adelantando CompScript: mostrá el fragmento de código Java con el campo `padre` y `buscar()` recursivo — es literalmente el próximo proyecto que van a construir ellos mismos (con roles invertidos: ellos escriben, vos revisás).

**Demo en vivo**: no aplica directamente sobre DataForge (el ejemplo de la demo es pseudocódigo genérico) — si querés anclarlo al proyecto real, mencioná que la `Entorno.simbolos` de DataForge es la versión de UN SOLO nivel de esta misma idea general.

**Preguntas probables y cómo responderlas**:
- *"¿Por qué `print(y)` dentro del bloque interno igual imprime 5 (el valor externo)?"* (pregunta real del quiz) — porque `y` no existe en el ámbito interno, entonces `get` sube por la referencia al entorno padre y lo encuentra ahí.
- *"¿Por qué en DataForge redeclarar la misma variable es error, pero en CompScript no necesariamente lo será?"* (pregunta real del quiz) — porque el conflicto de nombres solo existe dentro del MISMO ámbito; con bloques anidados, una variable interna puede legítimamente "tapar" (shadow) una externa sin ser error.
- *"¿Qué campo tiene el `Simbolo` de DataForge que un compilador clásico no necesitaría?"* (pregunta real del quiz) — el VALOR — porque DataForge es un intérprete (ejecuta ya mismo), no un compilador (que solo necesitaría tipo y ubicación de memoria, no el dato en sí).
- Duda típica que no está en el quiz: *"¿Por qué `LinkedHashMap` y no un `TreeMap` ordenado alfabéticamente?"* — porque el reporte de símbolos (§6.3) exige el orden de DECLARACIÓN, no el alfabético; `LinkedHashMap` preserva exactamente eso.

**Transición**: "La tabla de símbolos guarda nombres y valores. Pero ¿por qué DataForge no construye un árbol para representar el programa completo, como hacen la mayoría de libros de texto? Eso es la pregunta del AST."

### ast.html

**Objetivo pedagógico**: justificar con rigor por qué DataForge NO construye un árbol de sintaxis abstracta, y qué tendría que cambiar si algún día necesitara uno (control de flujo).

**Tiempo sugerido**: 15 min.

**Puntos clave a enfatizar oralmente**:
- Anunciá desde el inicio que esta es LA profundización que más ayuda a defender el proyecto en calificación — decilo explícitamente, motiva prestar atención.
- La demo de traza real (8 pasos reduciendo `var:double:: m <- SUM(2,3) end;`) conviene leerla completa y despacio — muestra la mecánica exacta que ya explicaron en Etapa 3, pero ahora con más detalle paso a paso.
- El argumento decisivo (slide 6, con la clase hipotética `Mientras`/while) es EL momento de la clase para practicar la respuesta de examen: leé textual la frase de defensa citada en la diapositiva, pedí que la memoricen casi literal.
- Cerrá con la frase resumen "Acción una vez · árbol N veces" — es compacta y se las van a acordar.

**Demo en vivo**: no aplica (es análisis de la traza teórica, no requiere correr DataForge) — aunque podés, si querés, correr `ejemplo1.df` en paralelo y pausar en la línea `m <- Media(@datos)` para que vean que el patrón de la demo se repite en cualquier expresión real.

**Preguntas probables y cómo responderlas**:
- *"¿Por qué la acción de `aritmetica` pudo ejecutar la suma ANTES de que se redujera `decl_var` completa?"* (pregunta real del quiz) — LR reduce de adentro hacia afuera: las cajas internas (la aritmética) se cierran antes que las externas (la declaración completa).
- *"¿Cuál es la diferencia entre el árbol de análisis (cajas/derivación) y el AST?"* (pregunta real del quiz) — el árbol de análisis es CONCRETO (incluye paréntesis, comas, todo el detalle sintáctico); el AST es ABSTRACTO (solo la operación y sus operandos, sin ruido de puntuación).
- *"Si en la calificación les pidieran agregar un `IF` a DataForge, ¿qué tendrían que cambiar?"* (pregunta real del quiz) — migrar a AST: crear interfaces `Expresion`/`Instruccion`, cambiar las acciones de `RESULT = calcular(...)` a `RESULT = new NodoX(...)`, y recorrer el árbol al final con `ejecutar(entorno)` en vez de ejecutar todo durante el parseo.
- Duda típica que no está en el quiz: *"¿Esto significa que DataForge 'hizo trampa' al no tener AST?"* — no; es una decisión de diseño VÁLIDA y defendible porque el lenguaje no tiene control de flujo (cada instrucción corre exactamente una vez) — el AST resuelve un problema (ejecutar código condicionalmente o repetidamente) que DataForge simplemente no tiene.

**Transición**: "Ya vimos por qué no hace falta un árbol para ESTE lenguaje. Demos un paso atrás y ubiquemos todo esto dentro del mapa completo de fases de un compilador — la figura que abre cualquier parcial de este curso."

### fases.html

**Objetivo pedagógico**: ubicar las 3 etapas que SÍ tiene DataForge (léxico/sintáctico/semántico) dentro del mapa completo de 6 fases de un compilador, y explicar por qué las últimas 3 no aplican.

**Tiempo sugerido**: 10 min (sin demos, la más corta).

**Puntos clave a enfatizar oralmente**:
- Dibujá o proyectá el diagrama de las 6 fases con las transversales (tabla de símbolos, manejo de errores) — es LA figura que van a ver una y otra vez en el resto del curso, vale la pena que se las graben visualmente.
- El ejemplo del libro (`position = initial + rate * 60`, fig 1.7) conviene recorrerlo fase por fase completo, incluyendo la coerción `inttofloat(60)` y el ensamblador final — es contraste directo con DataForge (que se detiene en la fase 3).
- La tabla de mapeo a DataForge (Lexer.flex/parser.cup/Entorno+Operaciones) es el resumen ejecutivo de las 6 etapas anteriores en una sola tabla — usala como repaso rápido si notás que el grupo perdió el hilo.
- Remarcá la ironía de la slide 5: DataForge es un intérprete (se queda en front-end) pero está escrito en Java, que es un HÍBRIDO (bytecode + JVM + JIT) — "capas de compilador sobre capas de compilador", una buena línea para cerrar con humor.

**Demo en vivo**: no aplica (sin código).

**Preguntas probables y cómo responderlas**:
- *"¿Por qué la tabla de símbolos y el manejo de errores NO son 'fases' del diagrama?"* (pregunta real del quiz) — porque no TRANSFORMAN la representación del programa de una forma a otra (que es lo que define una fase); son infraestructura transversal que todas las fases consultan y alimentan.
- *"¿En qué fase aparecería `inttofloat(60)` y por qué justo ahí?"* (pregunta real del quiz) — en la fase semántica, porque es la única fase que conoce los TIPOS de las variables y puede decidir que hace falta una conversión.
- *"¿DataForge tiene fase de generación de código intermedio?"* (pregunta real del quiz) — no, deliberadamente — es un intérprete, el front-end (léxico+sintáctico+semántico) desemboca directo en ejecución, sin pasar por representación intermedia ni optimización ni generación de código de máquina.
- Duda típica que no está en el quiz: *"¿Esto significa que DataForge es 'menos compilador' que uno completo?"* — no es menos, es OTRO TIPO: un intérprete resuelve el mismo problema (entender el lenguaje) sin necesitar las fases de back-end, porque en vez de generar código para correr después, ejecuta directamente ahora.

**Transición**: "Ya ubicamos las 3 fases del front-end en el mapa grande. Cerremos con la más filosófica de las profundizaciones: qué es exactamente 'semántica' y por qué DataForge no necesita coerciones de tipos."

### semantica.html

**Objetivo pedagógico**: formalizar el sistema de tipos completo de DataForge y explicar por qué, al no tener control de flujo, la diferencia entre chequeo estático y dinámico es inobservable en este proyecto (aunque dejará de serlo en CompScript).

**Tiempo sugerido**: 15 min (sin demos, cierre de la sesión de profundizaciones).

**Puntos clave a enfatizar oralmente**:
- La escalera "regex ⊂ GLC ⊂ reglas con contexto" es el resumen de TODO el curso hasta acá en una sola imagen — conviene literalmente dibujarla como jerarquía de círculos anidados.
- La tabla exhaustiva del sistema de tipos (slide 3) es densa — no la leas entera en voz alta, mejor señalá 2-3 filas representativas (aritmética double×double→double, DIV con 0 como ⊥) y dejá el resto como referencia de consulta.
- "La sorpresa" (slide 4: estático vs dinámico es inobservable en DataForge) es el punto más sutil de toda la clase — armá el argumento completo: sin control de flujo no hay ramas muertas, y sin ramas muertas no hay diferencia entre revisar antes de correr o revisar mientras corre, porque TODO va a correr sí o sí.
- Cerrá con la tabla de errores consolidada (slide 6) — es el repaso final de las 3 familias de error que atravesaron toda la clase, con su referencia exacta al capítulo del Dragón cada una.

**Demo en vivo**: no aplica (sin código nuevo) — si querés anclar el punto de "sin coerciones", podés señalar en cualquiera de los `.df` que TODOS los números son `double`, nunca hay un `int` separado que necesite conversión.

**Preguntas probables y cómo responderlas**:
- *"¿Por qué la regla 'declarar antes de usar' no puede expresarse en la gramática (GLC)?"* (pregunta real del quiz) — porque una GLC solo tiene una pila como memoria; esa regla exige recordar un CONJUNTO arbitrario de nombres ya declarados, que no cabe en la memoria de pila de un autómata de pila.
- *"¿DataForge hace comprobación de tipos estática o dinámica? ¿Importa la diferencia?"* (pregunta real del quiz) — dinámica (chequea con `instanceof Double` durante la ejecución); y no importa en la práctica porque no existe código muerto (sin `if`/`while`, toda instrucción se ejecuta exactamente una vez, así que no hay diferencia entre "revisar antes" y "revisar mientras corre").
- *"¿Por qué `DIV(x, 0)` es un error de VALOR (⊥) y no un error de TIPOS?"* (pregunta real del quiz) — porque los tipos están perfectamente bien (double × double es una operación válida); el problema es el valor concreto del segundo operando, no su tipo.
- Duda típica que no está en el quiz: *"Si DataForge no tiene coerciones, ¿cómo funciona `SUM(1, 2.5)` si 1 no tiene punto decimal?"* — en DataForge NO existe el tipo entero por separado; el lexer convierte todo literal numérico a `double` desde el análisis léxico (`Double.valueOf(yytext())`), así que `1` y `2.5` ya son el mismo tipo antes de llegar a la aritmética — no hay coerción porque nunca hubo dos tipos numéricos distintos.

---

## Cierre

**Mensaje de síntesis final** (decilo casi textual al terminar, sea que hayas cubierto solo el roadmap principal o también las profundizaciones):

> "Recorrimos las 3 capas clásicas de un compilador — léxico, sintáctico, semántico — construyendo cada una sobre DataForge real: un `Lexer.flex` que convierte texto en tokens, un `parser.cup` que valida la estructura con una gramática libre de contexto, y un `Entorno` que ejecuta directamente porque el lenguaje no tiene control de flujo. No hay AST, no hay fases de back-end, y esa AUSENCIA es una decisión de diseño defendible, no un atajo. Los 3 tipos de error (léxico, sintáctico, semántico) se recuperan cada uno a su manera, sin abortar la ejecución. Con Etapa 6 completa, DataForge cumple el enunciado entero — reportes, gráficas, editor, todo."

**Qué conecta con el próximo proyecto (ConjAnalyzer)**:

- Cambia el MODO de trabajo: en DataForge, Claude/el instructor construyó el proyecto como demostración guiada; en ConjAnalyzer (y CompScript, CompInterpreter) los roles se invierten — el estudiante escribe el código, el instructor tutorea y revisa. Anticipá esto explícitamente para que no esperen la misma dinámica.
- Reusan el MISMO stack técnico: Java 17 + JFlex + CUP + JavaFX + Maven, con la misma estructura de `pom.xml` (jflex-maven-plugin + cup-maven-plugin) que ya vieron funcionando en DataForge — DataForge ES la referencia de "esto sí compila y corre", así que pueden copiar patrones de ahí sin miedo.
- Lo que SÍ va a cambiar y conviene adelantar: ConjAnalyzer trabaja con conjuntos y probablemente necesite estructuras de datos y validaciones semánticas distintas a las de DataForge (números/estadística/gráficas) — pero el pipeline léxico→sintáctico→semántico y la lógica de recuperación de errores (léxico descarta carácter, sintáctico modo pánico, semántico propagación de error) son directamente trasladables.
- Si alcanzás a mencionar CompScript de pasada: ahí la tabla de símbolos SÍ necesitará ámbitos anidados de verdad (el "caso general" que se adelantó en `tabla-simbolos.html`) y probablemente SÍ haga falta un AST (por el control de flujo) — es la continuación natural de las dos profundizaciones que dejaron la puerta abierta a propósito.
