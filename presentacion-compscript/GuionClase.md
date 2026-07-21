# Guion de clase — CompScript

> Esto NO es la diapositiva. Es tu chuleta para dar la clase: qué decir, qué correr en vivo, qué responder si preguntan. Las diapositivas HTML (`etapa0.html` … `etapa7.html`, `ast.html`) están en 3ª persona para la audiencia; este documento está escrito para vos, en 2ª persona, como notas de director.

## Antes de empezar

**Cómo abrís la presentación:** doble clic en `index.html` (carpeta `presentacion-compscript/`). No necesita servidor ni internet — es el mismo motor que ya usaste en DataForge (`assets/estilo.css` + `assets/deck.js`), así que la navegación es idéntica: flechas ← → del teclado o los botones de abajo, botón ◐ arriba a la derecha para tema claro/oscuro. Desde el índice saltás a cada etapa con un clic; desde cualquier etapa, "⌂ Volver al índice" te regresa.

**Duración total estimada: ~2 horas** (117 min de contenido + margen de preguntas). Si el tiempo aprieta, las candidatas a recortar son Etapa 1 (léxico — el público ya vio JFlex en DataForge) y Etapa 7 (editor — es la parte menos conceptual). Etapa 3, Etapa 5 y la profundización `ast.html` son las que NO deberías apurar: ahí está el corazón del proyecto y los dos bugs reales.

**Qué preparar antes de arrancar:**
1. Tené **CompScript abierto en IntelliJ IDEA** (el proyecto en `Hades\CompScript\`), para correr demos en vivo sin salir de la clase.
2. El comando exacto para correr un ejemplo por consola con el Maven embebido de IntelliJ (usa el JDK 25 gestionado por el IDE, no el `java` 1.8 del PATH):
   ```bash
   JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemploN.cs"
   ```
   Cambiá `ejemploN.cs` por el archivo que toque en cada etapa (`ejemplo1.cs` … `ejemplo6_errores.cs`).
3. **Qué imprime este comando, para que no te agarre de sorpresa en vivo:** `TestInterprete` muestra, en este orden, (a) la consola del programa (lo que generó `console.log`), (b) la tabla de símbolos completa (reporte 6.4: id, tipo, ámbito, valor, línea/columna), (c) la lista de errores registrados, y (d) el conteo total de tokens reconocidos. **No** genera los HTML/`.dot` de reportes ni el AST visual — eso solo sale de la GUI (`Lanzador`, Etapa 7). Si querés mostrar el AST dibujado o el `TreeView`, tenés que correr la GUI aparte (▶ Run sobre `Lanzador`, nunca sobre `EditorApp`).
4. Si vas a demostrar el bug de `validarVector` (Etapa 3) o el de cortocircuito (Etapa 5), preparate un archivo `.cs` chico ad-hoc (te doy el contenido exacto en cada sección) — no viene en `entradas/`, hay que escribirlo en el momento o tenerlo ya guardado.

---

## Etapa 0 — Conceptos base: por qué CompScript necesita un AST

**Objetivo pedagógico:** que quede clarísimo, antes de ver una sola línea de código, POR QUÉ este proyecto no puede resolverse como DataForge.

**Tiempo sugerido:** 8 min. No hay demo — es la etapa más conceptual, apoyate solo en las diapositivas y en el pizarrón/tablero si querés dibujar la comparación.

**Puntos clave a enfatizar oralmente (no repitas el texto de la slide):**
- Arrancá con la pregunta directa a la clase: "¿alguien recuerda por qué en DataForge las acciones del `.cup` calculaban directo, tipo `RESULT = Operaciones.aritmetica(...)`?" — dejá que respondan, es repaso.
- El quiebre está en UNA palabra: **repetición**. Un parser LALR reduce cada producción una sola vez. Si tu lenguaje tiene `while`, necesitás ejecutar el mismo pedazo de gramática 0, 1 o N veces — y esa producción ya "pasó" y no vuelve.
- Remarcá que esto NO es una limitación de CUP en particular: es una propiedad de CUALQUIER análisis S-atribuido evaluado en el momento de la reducción. El AST es la solución general del libro (Dragón, cap. 5).
- Mencioná el detalle de las 3 pasadas acá mismo (aunque se profundiza en Etapa 6) porque es un beneficio directo de tener el árbol como estructura de datos: no se puede "pasar dos veces" por una acción de CUP ya ejecutada, pero sí por un árbol en memoria.

**Preguntas probables (del quiz real + dudas típicas):**
- *"¿Por qué una gramática S-atribuida no alcanza para un `while`?"* → Porque el parser reduce cada producción una vez, y el cuerpo del `while` necesita re-evaluarse según una condición que cambia en tiempo de ejecución — algo que ya pasó no puede "volver a pasar".
- *"¿Qué cambia en las acciones del `.cup`?"* → En vez de `RESULT = Operaciones.aritmetica(...)` (calcula ya), ahora es `RESULT = new A.Binaria("+", a, b, ...)` (construye un objeto que sabe calcularlo después).
- *Duda típica adicional: "¿por qué no usar directamente recursión en Java sin construir un árbol, evaluando el código fuente carácter a carácter cada vez que hace falta repetir?"* → Porque re-parsear el texto fuente en cada iteración sería carísimo y frágil; el AST separa "entender la sintaxis" (una vez) de "ejecutar" (las veces que haga falta), que es exactamente la idea de una fase de compilación separada.
- *Duda típica: "¿el árbol se guarda en disco?"* → No, vive en memoria durante toda la ejecución (es el valor de retorno de `parser.parse()`); lo que sí se puede volcar a disco es su *representación* para los reportes (Etapa 7: `ast.html`, `ast.dot`).

**Transición a Etapa 1:** "Bueno, ya sabemos el POR QUÉ. Ahora empecemos por el principio del pipeline real: cómo CompScript convierte texto plano en tokens."

---

## Etapa 1 — Tabla de tokens + analizador léxico con JFlex

**Objetivo pedagógico:** mostrar que el vocabulario de CompScript es más grande que el de DataForge, pero las reglas de diseño del lexer (orden de reservadas, match más largo, recuperación sin abortar) son las mismas que ya vieron.

**Tiempo sugerido:** 12 min (podés recortar a 8 si el público ya domina JFlex de DataForge — no repitas la teoría de autómatas, andá directo a lo específico de CompScript).

**Puntos clave a enfatizar oralmente:**
- No leas la tabla de las 37 reservadas en voz alta — remarcá solo lo nuevo: `^` es potencia (no XOR) y `$` es raíz cuadrada, ambos por decisión explícita del enunciado, no por convención de otro lenguaje.
- El punto más importante de la etapa es que **las secuencias de escape se resuelven en el lexer, no después**: cuando `CADENA` llega al parser, ya viene con el `\n` traducido a salto de línea real. Esto es una decisión de diseño que vale la pena que la clase entienda: ninguna capa posterior necesita saber que la fuente usaba barras invertidas.
- Recuperación de errores léxicos: el carácter inválido se descarta y el análisis SIGUE. Contrastalo bien con el error semántico (que sí aborta) — esa asimetría es a propósito (§4.3) y se retoma en Etapa 5.

**Demo en vivo:**
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo1.cs"
```
Señalá al final de la salida la línea `--- TOKENS: 209 reconocidos ---` — es el número verificado en la sesión de auditoría. Después corré:
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo6_errores.cs"
```
y mostrá en la sección `--- ERRORES ---` cómo aparece el error léxico del `#` suelto (línea 5) sin que el resto del archivo deje de analizarse (vas a ver también el error sintáctico y el semántico de ese mismo archivo — decile a la clase que por ahora se fijen solo en el primero, los otros dos vienen en Etapas 2 y 5).

**Preguntas probables:**
- *"¿Por qué `^` no puede ser XOR acá?"* → Porque el enunciado (§5.5.5) lo define específicamente como potencia para este lenguaje; el lexer lo tokeniza como `POT` y `Operaciones.aritmetica` lo trata como exponenciación.
- *"¿Cuántos tokens produce `console.log(\"ab\\n\");` y qué pasó con el `\n`?"* → Unos 8 tokens; el `\n` no es un token aparte, es parte del contenido de `CADENA`, ya traducido por `procesarEscapes()` antes de salir del lexer.
- *"Si el lexer encuentra un error, ¿se detiene todo el análisis?"* → No. Descarta el carácter, registra el error en `Contexto`, y sigue. Solo un error semántico detiene la ejecución.
- *Duda típica adicional: "¿por qué el orden de las reservadas importa tanto en el `.flex`?"* → Porque JFlex, ante un empate de longitud de match, prefiere la regla escrita primero. Si `{Id}` estuviera antes que `int`, la palabra `int` matchearía como identificador genérico y la reservada nunca existiría.

**Transición a Etapa 2:** "Ya tenemos los tokens. El siguiente problema es agruparlos en una gramática — y acá CompScript se pone más interesante que DataForge, porque usa notación infija de verdad."

---

## Etapa 2 — Gramática BNF + analizador sintáctico con CUP

**Objetivo pedagógico:** mostrar cómo CompScript resuelve la ambigüedad de la notación infija (algo que DataForge no tenía) usando la tabla de precedencia de CUP, sin partir la gramática en niveles.

**Tiempo sugerido:** 12 min.

**Puntos clave a enfatizar oralmente:**
- Arrancá contrastando con DataForge: ahí las llamadas tipo `SUM(a,b)` delimitaban solazamente los operandos con paréntesis — sin ambigüedad posible. Acá, `a + b * c` es ambigua por diseño de la BNF, y se resuelve con la tabla `precedence` del `.cup`, no con niveles `expr → termino → factor` como en el libro.
- El truco de `UMENOS` es el punto más rico de la etapa — dedicale tiempo. Es el patrón clásico de CUP/Yacc para "mismo símbolo léxico, dos precedencias distintas" (resta binaria vs. negación unaria). Si el público ya vio esto en DataForge o en teoría, preguntá primero si alguien recuerda cómo se resuelve el mismo problema en Yacc/Bison.
- El *dangling else* se resuelve por ORDEN de producciones, sin directiva de precedencia — vale la pena mostrar el orden exacto de las tres producciones de `<if>` en la slide y explicar por qué CUP nunca reportó conflicto shift-reduce ahí.
- Machacá el punto central otra vez (viene de Etapa 0): cada acción del `.cup` ahora CONSTRUYE un nodo, no calcula.

**Demo en vivo:** reusá la corrida de `ejemplo6_errores.cs` (si no la hiciste en Etapa 1, corré ahora):
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo6_errores.cs"
```
Señalá en `--- ERRORES ---` el mensaje sintáctico de la línea 8 (`let w: int = ;`) y explicá que el modo pánico descartó hasta el próximo `;` y siguió analizando la línea 11 sin abortar todo el archivo.

**Preguntas probables:**
- *"¿Cómo se agrupa `2 + 3 * 4`?"* → `2 + (3 * 4)`, porque `POR` está en un nivel de precedencia más fuerte que `MAS` en la tabla del `.cup`.
- *"¿Por qué `UMENOS` no aparece nunca como lexema en un `.cs`?"* → Es un pseudo-terminal: se declara solo para darle precedencia distinta a la negación unaria vía `%prec UMENOS`. El lexer nunca produce ese token; el símbolo real siempre es `MENOS`.
- *"¿La acción de `expr ::= expr MAS expr` calcula un número o construye un objeto?"* → Un objeto (`new A.Binaria("+", a, b, ...)`). El cálculo se pospone hasta que alguien llame a `evaluar()`.
- *Duda típica adicional: "¿qué pasa si me olvido el `%prec UMENOS`?"* → CUP reporta un conflicto shift-reduce al generar el parser, porque no sabe que esa producción específica de la resta unaria debe ligar distinto que la binaria. Es un error de compilación del `.cup`, no de ejecución — se detecta antes de correr nada.

**Transición a Etapa 3:** "Ya tenemos tokens agrupados en una gramática que construye objetos en vez de calcular. Ahora toca ver esos objetos de cerca: el árbol mismo."

---

## Etapa 3 — El árbol de sintaxis abstracta: `ast/A.java`

**Objetivo pedagógico:** esta es LA etapa distintiva de CompScript frente a los otros 3 proyectos de OLC1 — ninguno de DataForge, ConjAnalyzer o CompInterpreter construye un AST real. Si el público ya vio la presentación de DataForge, este es el momento de contrastar explícitamente.

**Tiempo sugerido:** 18 min — no la apures, es la más densa conceptualmente y tiene el primer bug real.

**Puntos clave a enfatizar oralmente:**
- **Contraste explícito con DataForge (decilo en voz alta, no des por sentado que lo recuerdan):** "En DataForge, cada acción de CUP evaluaba directo, sin AST, porque no había ninguna instrucción que necesitara ejecutarse una cantidad variable de veces. Acá, en cambio, cada nodo del árbol implementa una de dos interfaces — `Instruccion` (hace algo) o `Expresion` (produce un valor) — y sabe ejecutarse o evaluarse A SÍ MISMO. Es el único de los 4 proyectos que necesita esto."
- El patrón es *tree-walking interpreter* — una variante del Visitor sin doble despacho. Si la clase vio Visitor en algún curso de POO, es un buen gancho.
- `etiquetaAst()` / `hijosAst()` como "beneficio gratis": la misma estructura que usan `evaluar()`/`ejecutar()` alimenta directamente el reporte de AST (Etapa 7) sin necesitar una clase visitante aparte.
- Por qué los constructores reciben `Object` y castean adentro: porque las acciones del `.cup` trabajan con no-terminales de tipo *raw* (igual que en DataForge) — no hay genéricos seguros ahí.

**El bug real — dale tiempo, es un ejemplo pedagógico fuerte de "qué pasa cuando el árbol no valida lo que asume":**
`A.Declaracion.validarVector(...)` recibía el valor evaluado de la expresión inicializadora de un vector y asumía sin comprobar que YA era un vector. Antes de la corrección, `let v: int[] = 5;` (una expresión válida pero no vectorial) provocaba un `ClassCastException` de Java sin controlar que tumbaba el intérprete ENTERO — no solo el programa `.cs`, sino el proceso completo. Eso viola directamente el enunciado (§4.3: los errores semánticos deben reportarse y terminar de forma ordenada, no reventar). La corrección agrega la guarda de tipo antes de castear.

**Demo en vivo (opcional pero recomendada — muestra el bug YA corregido en acción):** creá un archivo chico, por ejemplo `entradas/demo_bug.cs`, con una sola línea:
```
let v: int[] = 5;
```
y corré:
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/demo_bug.cs"
```
Mostrá que el resultado es un error semántico ordenado en la sección `--- ERRORES ---` (algo como `el vector 'v' espera un literal de vector y recibio Entero`), NO un crash de Java. Si te queda tiempo, contá que antes de la corrección esto tiraba abajo el proceso completo.

**Preguntas probables:**
- *"¿Qué interfaz implementa `A.If`? ¿Y `A.Binaria`?"* → `A.If` implementa `Instruccion` (no produce valor, decide qué ejecutar). `A.Binaria` implementa `Expresion` (produce un `Valor`).
- *"¿Cuál habría sido el síntoma del bug ANTES de la corrección?"* → Un `ClassCastException` sin capturar que interrumpía TODO el intérprete, no solo el programa que se estaba ejecutando.
- *"¿Por qué no hace falta una clase visitante para el reporte de AST?"* → Porque cada nodo ya sabe describirse (`etiquetaAst()`) y enumerar a sus hijos (`hijosAst()`); el generador de reportes solo recorre esa estructura recursivamente.
- *Duda típica adicional (aprovechá para el contraste con DataForge): "¿por qué en DataForge el chequeo de tipos era más simple?"* → Porque con evaluación directa el chequeo ocurre en el mismo lugar donde se usa el valor, una sola vez. Con AST, cada nodo repite esa responsabilidad en su propio `evaluar()`/`ejecutar()` — no hay una gramática externa que lo garantice por vos. Es exactamente el motivo por el que apareció este bug: nadie "afuera" del nodo estaba validando por él.

**Transición a Etapa 4:** "El árbol ya sabe ejecutarse a sí mismo, pero necesita dónde guardar y buscar variables. Eso es la tabla de símbolos — y acá CompScript necesita algo que DataForge no: una pila de ámbitos."

---

## Etapa 4 — Tabla de símbolos: `Contexto` y `Entorno`, ámbitos anidados

**Objetivo pedagógico:** mostrar que la diferencia entre un mapa plano (DataForge) y una cadena de ámbitos enlazados (CompScript) es consecuencia directa de tener bloques y funciones.

**Tiempo sugerido:** 12 min.

**Puntos clave a enfatizar oralmente:**
- Separá bien las dos responsabilidades: `Contexto` es ÚNICO por ejecución (consola, errores, tokens, símbolos históricos, funciones/structs registrados); `Entorno` es MÚLTIPLE — uno por cada bloque (`if`, `while`, cuerpo de función).
- La decisión de diseño que más vale la pena explicar despacio: **el entorno de una función SIEMPRE cuelga del entorno global, nunca del entorno de quien la invoca.** Esto es lo que le da a CompScript alcance estático. Si alguien preguntara "¿y si en vez de eso colgara del invocador?" — ya está en el quiz, tenés la respuesta lista abajo.
- El dato que más impacta a la clase suele ser el de las 186 entradas de `n` en el reporte de símbolos de `fib(10)`: usalo como evidencia tangible, no solo como afirmación teórica, de que cada llamada recursiva tiene su propio `Entorno`.
- Aclará que `ctx.simbolos` es un LOG histórico, no un snapshot final — nunca se sobrescribe. Esto puede sorprender a quien esperaría "ver el valor final de cada variable".

**Demo en vivo:**
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo5.cs"
```
En la sección `--- TABLA DE SIMBOLOS ---`, contá (o aproximá) cuántas filas del parámetro `n` aparecen para `fib` — es el mismo dato que cita la diapositiva (186 entradas). Señalá que cada fila tiene un `Entorno`/ámbito distinto aunque el nombre `n` se repita.

**Preguntas probables:**
- *"¿Por qué DataForge no necesitaba una pila de `Entorno`s?"* → Porque no tiene bloques anidados ni funciones — todas sus variables viven en un único ámbito plano.
- *"Si `fibonacci(n=10)` hace ~177 llamadas recursivas, ¿cuántos `Entorno`s distintos existieron?"* → Uno por cada llamada — cada invocación crea su propio hijo del global.
- *"¿Qué pasaría si el entorno de una función colgara del invocador en vez del global?"* → La función podría ver variables locales de quien la invoca — *dynamic scoping* en vez de *static scoping* — rompiendo el modelo mental esperado.
- *Duda típica adicional: "¿qué pasa si dos variables en ámbitos distintos tienen el mismo nombre?"* → Es shadowing normal: `buscar()` sube por la cadena de padres y encuentra primero la del ámbito más cercano; la de más afuera queda "tapada" mientras dure el ámbito interno, sin perderse.

**Transición a Etapa 5:** "Con el árbol y la pila de ámbitos ya armados, llegamos al momento en que TODO esto se justifica: el flujo de control de verdad."

---

## Etapa 5 — Flujo de control: if / match / while / for / do-while

**Objetivo pedagógico:** demostrar en código real que `if`/`while`/`for`/`do-while` re-evalúan su condición cuantas veces haga falta — la razón estructural completa del proyecto, cerrando el arco que abrió la Etapa 0. Y presentar el segundo bug real de la auditoría, uno bueno para hablar de calidad de código.

**Tiempo sugerido:** 15 min — no la recortes, tiene el bug más interesante de exponer en clase.

**Puntos clave a enfatizar oralmente:**
- Mostrá el patrón `while(true) { evaluar cond; if(!cond) break; ejecutar cuerpo; catch Break/Continue }` como el mecanismo Java literal que hace posible recorrer el mismo subárbol N veces — la contraparte concreta de lo que en Etapa 0 quedó solo como argumento teórico.
- `break`/`continue`/`return` como excepciones de control (no banderas booleanas) es una decisión elegante: dejá que la clase note que esto resuelve solo el anidamiento dentro de varios `if` sin código manual de "revisar bandera después de cada instrucción".
- `match` sin fall-through: el `return` explícito dentro del bucle de casos ES la implementación literal de "no hace falta `break`" del enunciado.

**El bug del cortocircuito — el momento fuerte de la etapa, dale el tiempo que pide el enunciado del guion:**
La versión original de `A.Binaria.evaluar()` evaluaba `izq` Y `der` ANTES de entrar al `switch` que decide el operador — es decir, `&&` y `||` evaluaban SIEMPRE los dos lados, sin cortocircuito. Es un ejemplo perfecto de **bug sutil que pasa desapercibido**: el enunciado no exige cortocircuito explícitamente, así que nada en los requisitos lo iba a atrapar; y en los 6 ejemplos de `entradas/` no hay ningún caso que dependa de él, así que las pruebas normales tampoco lo detectaban. Recién con un caso ad-hoc de guarda típica (`x != 0 && 10/x > 1`) se hace visible: sin cortocircuito, `10/x` se evaluaba incluso cuando `x==0`, y explotaba con división entre cero — justo lo que la guarda estaba tratando de evitar.

**Demo en vivo (2 partes):**

1. Ciclos con `break`/`continue`:
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo3.cs"
```
En `--- CONSOLA ---` señalá la salida `0 10 30` del `for` y explicá en vivo por qué no aparece `20` (el `continue` en `i==2`) ni `40` (el `break` en `i==4`).

2. El bug de cortocircuito, YA corregido — armá un archivo `entradas/demo_cortocircuito.cs`:
```
void main() {
    let x: int = 0;
    console.log(x != 0 && 10 / x > 1);
}
RUN_MAIN main();
```
y corré:
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/demo_cortocircuito.cs"
```
Mostrá que imprime `false` sin error de división entre cero — y contale a la clase que ANTES de la corrección esto SÍ tiraba un error semántico de división por cero, porque `10/x` se evaluaba igual aunque `x != 0` ya fuera falso.

**Preguntas probables:**
- *"¿Qué diferencia hay entre `while` y `do-while` en el número mínimo de ejecuciones?"* → `while` puede correr 0 veces (evalúa ANTES); `do-while` corre al menos 1 vez (evalúa DESPUÉS) — confirmado por `ejemplo3.cs` con `k` imprimiendo `0 1`.
- *"¿Por qué `Senales.Break` extiende `RuntimeException` en vez de una bandera booleana?"* → Porque un `break` puede estar anidado dentro de varios `if` o ciclos; la excepción deja que Java "desenrolle" automáticamente esas capas sin revisar una bandera después de cada instrucción.
- *"Si `verificar()` tiene efectos secundarios y se evalúa `false && verificar()`, ¿corre `verificar()`?"* → No, con la corrección aplicada. `izq` se evalúa primero, ya decide `false`, y `der` nunca se evalúa. (Antes de la corrección, sí corría siempre — era el bug.)
- *Duda típica adicional: "¿esto se considera parte del enunciado o fue puramente un hallazgo de auditoría?"* → Fue un hallazgo de auditoría de código (2026-07-21), no algo pedido literalmente por el enunciado — pero sí es un comportamiento esperado en cualquier lenguaje con `&&`/`||`, y su ausencia SÍ era observable semánticamente (rompía guardas típicas), no solo una pérdida de eficiencia. Buen ejemplo de por qué conviene auditar código aunque los tests pasen.

**Transición a Etapa 6:** "Ya vimos que el árbol se puede recorrer N veces. Las funciones llevan esa idea un paso más allá: el árbol se recorre desde ADENTRO de sí mismo — eso es recursión."

---

## Etapa 6 — Funciones, argumentos por nombre y recursión

**Objetivo pedagógico:** cerrar el argumento de "el árbol como estructura reutilizable" mostrando recursión real, y la decisión de diseño de argumentos nombrados (no posicionales).

**Tiempo sugerido:** 12 min.

**Puntos clave a enfatizar oralmente:**
- Las 3 pasadas de `Interprete.ejecutar(...)` ya se mencionaron en Etapa 0 — acá es donde se ve el PORQUÉ concreto: la primera pasada registra funciones/structs para permitir referencias hacia adelante, algo que la evaluación directa de una sola pasada jamás podría resolver.
- Argumentos por nombre (`potencia(base = 4)`) es una decisión explícita del enunciado (§5.23), no una elección libre — remarcá que el ORDEN no importa, solo el nombre.
- El detalle sutil que vale la pena señalar en voz alta: el argumento se evalúa en el entorno del LLAMADOR (`caller`), pero se declara en el entorno NUEVO de la función (`local`) — los dos entornos coexisten solo durante esa resolución puntual.
- `return` viaja como excepción (`Senales.Retorno`) con el valor adentro — mismo mecanismo que `break`/`continue` de la Etapa 5, no algo nuevo.

**Demo en vivo:**
```bash
JAVA_HOME="/c/Users/72358/.jdks/openjdk-25.0.1" "/c/Program Files/JetBrains/IntelliJ IDEA 2025.2.4/plugins/maven/lib/maven3/bin/mvn" compile exec:java -Dexec.args="entradas/ejemplo5.cs"
```
Mostrá en `--- CONSOLA ---` los resultados reales: `factorial(n=5)` → 120, `fib(n=10)` → 55, `potencia(base=4)` → 16 (usa el `exp` por defecto), `potencia(base=2, exp=5)` → 32. Es un buen momento para preguntarle a la clase, ANTES de mostrar la salida, cuánto creen que va a dar `potencia(base=4)` sin decirles que hay un valor por defecto — para que la sorpresa del default quede grabada.

**Preguntas probables:**
- *"¿Por qué la primera pasada es necesaria antes de ejecutar cualquier instrucción global?"* → Para permitir referencias hacia adelante: una función puede llamar a otra declarada más abajo en el archivo.
- *"En `potencia(base = 4)`, ¿de dónde sale `exp`?"* → Del valor por defecto declarado en el parámetro (`exp: int = 2`), evaluado cuando no aparece ningún argumento `exp=...` en la llamada.
- *"¿Qué evidencia real demuestra que cada llamada recursiva tiene su propio `Entorno`?"* → El reporte de símbolos de `ejemplo5.cs` con las 186 entradas de `n` en distintos ámbitos (ya lo vieron en Etapa 4 — es buen momento para conectar ambas etapas).
- *Duda típica adicional: "¿qué pasa si omito un argumento sin valor por defecto?"* → `A.invocar(...)` reporta un error semántico ("falta el argumento '...'") antes de intentar ejecutar el cuerpo de la función.

**Transición a Etapa 7:** "El pipeline entero ya funciona por consola. Lo que falta es la cara visible: el editor y los reportes que pide el enunciado."

---

## Etapa 7 — Editor JavaFX + los 5 reportes

**Objetivo pedagógico:** cerrar el recorrido mostrando la interfaz gráfica y las 3 formas distintas en que CompScript expone el AST (algo que ningún proyecto hermano necesita mostrar).

**Tiempo sugerido:** 10 min — es la etapa menos conceptual, no te demores de más acá.

**Puntos clave a enfatizar oralmente:**
- `Lanzador` resuelve el mismo problema de siempre de JavaFX ("JavaFX runtime components are missing") — si el público vio DataForge, es el mismo patrón exacto, decilo explícitamente para que no piensen que es algo nuevo de CompScript.
- El punto realmente nuevo frente a DataForge: **5 reportes en vez de 4** — el extra es el AST, y encima en TRES formas distintas (`ast.html` lista anidada, `ast.dot` para Graphviz, `TreeView` embebido en la GUI). Explicá por qué tres y no una: el enunciado (§4.5) exige mostrarlo "desde la interfaz" (eso lo cumple el `TreeView`); el HTML y el `.dot` son un plus de documentación externa, y las tres reutilizan el mismo `etiquetaAst()`/`hijosAst()` de la Etapa 3 — no triplican lógica.
- Nombres de token por reflexión sobre `sym.class`: mismo truco que DataForge, mencionalo rápido si ya lo vieron.

**Demo en vivo (si el tiempo alcanza):** en IntelliJ, ▶ Run sobre `Lanzador` (NUNCA sobre `EditorApp`). Escribí o abrí `entradas/ejemplo5.cs` en el editor, presioná ▶ Ejecutar, y mostrá la consola. Si querés mostrar el AST, buscá el botón "Ver AST" — te muestra el `TreeView` embebido con la misma estructura que viste dibujada a mano en la profundización (próxima sección). Si no da el tiempo, decilo explícitamente a la clase ("no vamos a abrir la ventana ahora por tiempo, pero el código está completo y lo pueden correr ustedes") — es más honesto que apurar una demo a medias.

**Preguntas probables:**
- *"¿Por qué `Lanzador` no extiende `Application` directamente?"* → Porque si la clase con el `main` extiende `Application` directo, la JVM aplica una verificación de módulos que falla con "JavaFX runtime components are missing" cuando JavaFX no está en el *module path*.
- *"Si se agrega un terminal nuevo en `sym.java`, ¿hay que tocar el código de reportes?"* → No — los nombres se leen por reflexión sobre los campos públicos de `sym.class`, cualquier terminal nuevo aparece solo.
- *"¿Qué archivo tiene CompScript que DataForge no genera?"* → `ast.html` (y `ast.dot`). DataForge no construye AST, así que no tiene árbol que reportar.
- *Duda típica adicional: "¿por qué no alcanza con el `TreeView` solo, si ya cumple el requisito de la interfaz?"* → Porque el HTML y el `.dot` sirven para documentación FUERA de la GUI (por ejemplo, para el manual entregable o para imprimir el árbol con Graphviz) — el `TreeView` no se puede "adjuntar" a un PDF de entrega tal cual.

**Transición a la profundización:** "Con el pipeline completo ya armado, vale la pena parar y ver TODO junto sobre un ejemplo real — esta es la parte obligatoria de la clase."

---

## Profundización ★ — El AST real de CompScript (`ast.html`)

**Objetivo pedagógico:** consolidar las 5 etapas anteriores sobre un único ejemplo real (`factorial`), y dejar a la clase con una frase de defensa de arquitectura que puedan repetir en su propia calificación. Es la sección más importante para cerrar bien el proyecto.

**Tiempo sugerido:** 18 min. No la sacrifiques por tiempo — el enunciado de esta clase la marca como obligatoria, y es donde más rinde el contraste con DataForge si el público ya la vio.

**Puntos clave a enfatizar oralmente:**
- Esta página tiene una demo interactiva tipo "stepper" (botón ▶ Paso) que reconstruye, reducción por reducción, cómo el parser arma el árbol de `n * factorial(n = n - 1)` sin calcular nada. Dejá que la clase pida los pasos uno por uno en vez de mostrarlo todo de una — el ritmo importa acá, es la misma demo pedagógica del patrón que ya usaste en `gramaticas.html` de DataForge.
- Cuando llegues al árbol completo dibujado en ASCII (`function Entero factorial / ├── param n ...`), señalá el detalle de que las llaves, paréntesis y punto y coma NO aparecen en ningún lado del árbol — ya cumplieron su función guiando al parser y no son parte de la sintaxis *abstracta*.
- **Este es el momento de decir en voz alta, sin rodeos, la frase de defensa que da la propia diapositiva** (está pensada literalmente para que la repitan en su calificación): *"CompScript tiene control de flujo real —if, ciclos, funciones con recursión— por lo que cada instrucción puede necesitar ejecutarse cero, una o N veces. Eso es imposible con evaluación directa en las acciones de un parser LALR, que reduce cada producción una única vez. Por eso el análisis sintáctico construye un AST y la ejecución es una fase separada que recorre ese árbol tantas veces como el propio programa lo exija."*
- **Si el público ya vio la presentación de DataForge, hacé el contraste EXPLÍCITO acá, no lo des por sobreentendido:** mostrá la tabla comparativa de la slide 5 (evaluación directa vs. recorrido de árbol) y remarcá el dato más contundente — el árbol de `factorial` se construye UNA sola vez (durante el parseo) pero se recorre 5 veces (una por cada llamada recursiva con `n=5,4,3,2,1`), cada vez con un `Entorno` distinto. Eso es literalmente lo que DataForge no podría reproducir jamás con su arquitectura.

**Demo en vivo (complementaria a la demo HTML, opcional):** si te queda tiempo, corré de nuevo `ejemplo5.cs` y, mientras la clase mira la salida de consola con `factorial(n=5) → 120`, dibujá en el pizarrón (o señalá en la tabla de la slide 4) la cadena de 5 llamadas anidadas, cada una con su propio `n` — es la misma idea que el reporte de símbolos de Etapa 4, pero ahora conectada al árbol dibujado.

**Preguntas probables:**
- *"¿Por qué hay DOS nodos `AccesoVariable` distintos con el mismo nombre `n` en `return n * factorial(n = n - 1)`?"* → Porque cada aparición TEXTUAL de `n` en el código fuente produce su propio nodo al reducirse — el árbol no deduplica identificadores. Ambos nodos, al evaluarse, consultan el mismo `Entorno` y encuentran el mismo valor, pero son dos objetos distintos en memoria.
- *"¿Cuántas veces se construye el árbol de `factorial` durante `factorial(n=5)`? ¿Y cuántas veces se recorre?"* → Se construye UNA sola vez (al parsear). Se recorre 5 veces — una por cada llamada recursiva, cada una con un `Entorno` distinto.
- *"Si a DataForge le agregaran un `while` en la calificación, ¿qué tendría que cambiar?"* → Migrar a AST: crear interfaces `Instruccion`/`Expresion` y una clase por construcción (igual que `ast.A`), cambiar cada acción de CUP de "calcular" a "construir nodo", y después del parse recorrer la raíz llamando a `ejecutar(entorno)`. El lexer y la gramática de operadores quedarían intactos — solo cambia QUIÉN ejecuta y CUÁNDO.
- *Duda típica adicional (si alguien pregunta por qué no partieron la gramática de expresiones en niveles como el libro clásico recomienda): "¿por qué usar tabla de precedencia de CUP y no `expr → termino → factor`?"* → Porque el propio enunciado (§5.9) ya opta por ese mecanismo, y es exactamente lo que el CUP real de CompScript usa (visto en Etapa 2) — es una elección de herramienta válida, no una desviación del libro; el resultado gramatical es equivalente, solo cambia dónde se resuelve la ambigüedad (en la tabla, no en la estructura de no-terminales).

**Transición al cierre:** "Con esto ya recorrimos el proyecto completo, de punta a punta. Cerremos conectándolo con lo que sigue."

---

## Cierre

**Tiempo sugerido:** 5 min.

**Síntesis final (guion sugerido, adaptalo a tu estilo):**
"CompScript demostró algo que DataForge no necesitaba: que en cuanto un lenguaje tiene control de flujo real —`if`, ciclos, funciones con recursión— hace falta separar 'entender la sintaxis' de 'ejecutarla'. Esa separación es el AST: se construye una vez, durante el análisis sintáctico, y se recorre las veces que el programa lo pida. Vimos el pipeline completo — lexer (37 reservadas + 28 símbolos), parser (precedencia CUP + `UMENOS` + *dangling else*), el árbol mismo (~25 clases de nodo, cada una sabiendo ejecutarse a sí misma), la pila de ámbitos con alcance estático, el flujo de control con excepciones de control, y funciones con recursión real. Y de yapa, encontramos dos bugs reales en la auditoría — el de `validarVector` y el de cortocircuito en `&&`/`||` — que son un buen recordatorio de que 'pasa todos los tests' no es lo mismo que 'está bien auditado'."

**Conexión con el próximo proyecto:**
"El siguiente proyecto del curso es **CompInterpreter**, y ahí el stack cambia por completo: en vez de Java + JFlex + CUP + JavaFX, es JavaScript/TypeScript + Jison + Express (servidor) + React (cliente). La lógica de fondo — lexer, parser, y la necesidad de un AST si hay control de flujo — es la MISMA teoría que acaban de ver acá; lo que cambia es la herramienta (Jison en vez de JFlex/CUP) y que la interfaz ahora es web, no de escritorio. Si algo de esta clase les quedó claro, ya tienen el 80% del modelo mental para entender CompInterpreter sin partir de cero."

Si el público pregunta por qué no se armó el mismo tipo de comparación con ConjAnalyzer: mencioná que ConjAnalyzer es S-atribuido igual que DataForge (teoría de conjuntos, sin control de flujo real), así que el contraste relevante para AST es específicamente contra DataForge, no contra los tres proyectos por igual.

---

*Fuentes usadas para armar este guion: `index.html`, `etapa0.html`…`etapa7.html`, `ast.html` y `README.md` de esta carpeta; `CompScript/entradas/ejemplo1.cs`…`ejemplo6_errores.cs`; `CompScript/src/main/java/compscript/TestInterprete.java` (para confirmar qué imprime el comando de consola).*
