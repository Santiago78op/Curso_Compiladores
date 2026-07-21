# Manual Técnico — DataForge

**Universidad de San Carlos de Guatemala**
**Facultad de Ingeniería — Escuela de Ciencias y Sistemas**
**Organización de Lenguajes y Compiladores 1 — Proyecto 1**

---

## 1. Propósito de este documento

Este manual describe la arquitectura interna de DataForge, las decisiones de diseño tomadas durante su construcción y las funciones y métodos más relevantes de cada paquete, con el propósito de que cualquier persona distinta al autor original pueda dar mantenimiento al proyecto en el futuro (corregir errores, extender el lenguaje o modificar la interfaz gráfica) sin necesidad de leer todo el código fuente desde cero.

## 2. Lenguaje y herramientas

| Aspecto | Detalle |
|---|---|
| Lenguaje de implementación | Java 17 (compilado y probado también con JDK 25) |
| Gestor de dependencias / build | Maven, con `pom.xml` en la raíz del proyecto |
| Generador de analizador léxico | JFlex, plugin `de.jflex:jflex-maven-plugin:1.9.1` |
| Generador de analizador sintáctico | CUP (LALR), plugin `com.github.vbmacher:cup-maven-plugin:11b-20160615-3` |
| Runtime de CUP | `com.github.vbmacher:java-cup-runtime:11b-20160615` |
| Interfaz gráfica | JavaFX 21 (`org.openjfx:javafx-controls:21.0.4`), construida por código (sin FXML) |
| Gráficas | JavaFX Charts (`BarChart`, `PieChart`, `LineChart`), incluidos en `javafx-controls` |
| Ejecución de la GUI | `org.openjfx:javafx-maven-plugin:0.0.8` (objetivo `javafx:run`) |

El ciclo de generación de código en Maven ocurre en la fase `generate-sources`: los plugins de JFlex y CUP leen `src/main/jflex/Lexer.flex` y `src/main/cup/parser.cup` respectivamente, y depositan las clases generadas (`Lexer`, `Parser`, `sym`) en `target/generated-sources/`. Estas clases **nunca deben editarse directamente**: cualquier cambio se pierde en la siguiente compilación y debe hacerse en los archivos fuente `.flex` / `.cup`.

## 3. Estructura de paquetes

```
src/main/java/dataforge/
├── TestLexer.java             (prueba de consola: solo el lexer — tabla de tokens cruda)
├── TestParser.java            (prueba de consola: solo el parser — valida sin ejecutar)
├── TestInterprete.java        (prueba de consola: lexer → parser → ejecución)
├── gui/
│   ├── EditorApp.java         (ventana principal, editor, botones, orquestación)
│   ├── Lanzador.java          (punto de entrada real — ver sección 5.2)
│   └── Graficador.java        (Grafica → JavaFX Chart → Stage)
├── interprete/
│   ├── Entorno.java           (estado de ejecución: símbolos, consola, errores, gráficas)
│   ├── Interprete.java        (fachada String → Entorno)
│   ├── Operaciones.java       (aritmética y estadística con chequeo de tipos)
│   ├── Simbolo.java           (entrada de la tabla de símbolos)
│   ├── Grafica.java           (tipo + atributos resueltos de una gráfica)
│   └── RegistroError.java     (registro estructurado de un error)
└── reportes/
    └── Reportes.java          (generación de los 3 reportes HTML)

src/main/jflex/Lexer.flex      (especificación léxica — fuente, editable)
src/main/cup/parser.cup        (gramática + acciones — fuente, editable)
target/generated-sources/      (Lexer.java, Parser.java, sym.java — GENERADO, no editar)
```

El paquete `dataforge.analisis` (`Lexer`, `Parser`, `sym`) no tiene código fuente propio en `src/`: es enteramente generado a partir de `Lexer.flex` y `parser.cup` en cada compilación.

## 4. Arquitectura del pipeline de análisis y ejecución

### 4.1 Decisión central: sin árbol de sintaxis (AST)

DataForge ejecuta las instrucciones **directamente en las acciones semánticas del archivo `parser.cup`** (bloques `{: ... :}`), sin construir un árbol de sintaxis abstracta intermedio. Esta decisión es viable porque el lenguaje **no tiene control de flujo** (no hay condicionales ni ciclos): cada instrucción del programa fuente se ejecuta exactamente una vez, en el mismo orden en que el analizador sintáctico la reduce. Un lenguaje con control de flujo (como el de proyectos posteriores del curso) sí requeriría construir un AST para poder recorrerlo en un orden distinto al de la entrada (por ejemplo, repetir un bloque en un ciclo).

La gramática es **S-atribuida** (atributos sintetizados únicamente): cada producción calcula su propio valor (`RESULT`) a partir de los valores ya calculados de sus componentes, y ese valor sube hacia la producción que la contiene. El orden de evaluación queda garantizado por el orden de las reducciones del analizador ascendente LALR: los símbolos internos siempre se reducen antes que los externos que los contienen.

Ejemplo real, tomado de `src/main/cup/parser.cup`:

```java
aritmetica ::= op_arit:op PAR_IZQ expr:a COMA expr:b PAR_DER
               {: RESULT = Operaciones.aritmetica(op, a, b,
                               parser.entorno, opleft, opright); :} ;
```

Los no terminales que transportan un valor se declaran con tipo (`non terminal Object expr;`, `non terminal ArrayList lista_expr;`), usando siempre tipos *raw* (sin genéricos ni arreglos), ya que CUP no admite de forma confiable genéricos en las declaraciones de no terminales.

El estado completo de una ejecución (tabla de símbolos, consola, errores y gráficas registradas) vive en una única instancia de `Entorno`, expuesta en el parser como campo público:

```java
parser code {:
    public Entorno entorno = new Entorno();
    public void syntax_error(Symbol s) { ... }
:};
```

Las acciones de las producciones acceden a este estado a través de `parser.entorno`, por ejemplo `parser.entorno.declararVariable(id, t, e, idleft, idright)`.

### 4.2 Pipeline completo

```
código fuente (.df / String)
        │
        ▼
   Lexer (JFlex)  ──► produce Symbol (token + valor + línea/columna)
        │
        ▼
   Parser (CUP, LALR)  ──► reduce producciones; cada acción ejecuta
        │                  contra parser.entorno
        ▼
     Entorno  ──► tabla de símbolos, consola, errores, gráficas
```

La clase `dataforge.interprete.Interprete` es la fachada que expone este pipeline a la interfaz gráfica (y a cualquier otro cliente futuro) sin exponer los detalles de JFlex/CUP:

```java
public static Entorno ejecutar(String codigo) {
    Lexer lexer = new Lexer(new StringReader(codigo));
    Parser parser = new Parser(lexer);
    lexer.entorno = parser.entorno;   // el lexer registra tokens y errores léxicos
    try {
        parser.parse();
    } catch (Exception e) {
        // error sintáctico irrecuperable ya registrado por syntax_error()
    }
    return parser.entorno;
}
```

Nótese la línea `lexer.entorno = parser.entorno;`: conecta el lexer generado con el mismo `Entorno` que usa el parser, de modo que el lexer pueda registrar cada token reconocido (para el reporte de tokens, sección 7.1) y cada error léxico, en la misma estructura de datos que usarán el parser y el resto del sistema.

### 4.3 Entorno fresco por ejecución

Cada llamada a `Interprete.ejecutar(...)` crea una instancia **nueva** de `Parser` (y por lo tanto de `Entorno`, ya que `Entorno` se instancia en el campo `parser code` del `.cup`). Esto garantiza que el estado de una ejecución no se mezcle con el de la siguiente: los reportes (sección 7) siempre reflejan exclusivamente el último análisis realizado, tal como exige el enunciado del proyecto. En `gui/EditorApp.java`, el resultado de cada ejecución se guarda en el campo `ultimoEntorno`, que es el que consumen tanto el dibujado de gráficas como la generación de reportes.

### 4.4 Propagación de errores por valor nulo

Cuando una expresión no puede evaluarse (por ejemplo, una variable no declarada, o una operación aritmética con un operando de tipo incorrecto), el método correspondiente registra el error en el `Entorno` y devuelve `null` en vez de un valor. Las operaciones que reciben `null` como operando **no vuelven a reportar el error**: simplemente propagan `null` hacia arriba. Esta convención evita cascadas de errores redundantes (un mismo problema de origen no genera un error distinto por cada expresión que lo contiene) y aparece de forma consistente en `Operaciones.aritmetica`, `Operaciones.estadistica` y en los métodos de declaración de `Entorno`:

```java
// Operaciones.aritmetica
if (a == null || b == null) return null;
```

```java
// Entorno.declararVariable
if (valor == null) return;  // la expresión ya reportó su propio error
```

### 4.5 Recuperación de errores por fase

| Fase | Mecanismo | Efecto |
|---|---|---|
| Léxico | El carácter no reconocido se descarta (regla `[^]` al final del `.flex`, registrada vía `entorno.error(...)`) | El análisis continúa desde el siguiente carácter |
| Sintáctico | **Modo pánico** (Dragón, cap. 4): producción `instruccion ::= error PUNTO_COMA` | CUP descarta símbolos de la pila hasta poder desplazar el terminal especial `error`, y descarta tokens de la entrada hasta el próximo `;` — la instrucción defectuosa se sacrifica completa, pero el análisis continúa con la siguiente |
| Semántico | Propagación por `null` (sección 4.4) | La expresión afectada no produce valor; el resto del programa continúa normalmente |

En los tres casos el objetivo es el mismo: el reporte de errores (sección 7.2, correspondiente al §6.2 del enunciado) debe acumular **todos** los errores de un análisis, no detenerse en el primero.

## 5. Paquete `dataforge.gui`

### 5.1 `EditorApp`

Clase principal de la interfaz, construida **por código** (sin FXML), para que la estructura del *scene graph* sea explícita: `BorderPane` en la raíz, una `HBox` de botones en la zona superior, y un `SplitPane` vertical en el centro que separa el `TabPane` del editor (arriba) de la `TextArea` de consola, no editable (abajo).

Métodos relevantes:

- `nuevaPestana(String titulo, String contenido, File archivo)`: crea una pestaña; el archivo asociado viaja en `Tab.userData` (`null` si el archivo aún no se ha guardado).
- `abrir(Stage)` / `guardar(Stage)`: usan `FileChooser` con filtro `*.df`; `guardar` reutiliza el archivo asociado a la pestaña o pide uno nuevo mediante "Guardar como" si la pestaña es nueva.
- `ejecutar()`: obtiene el texto de la pestaña activa, invoca `Interprete.ejecutar(texto)`, arma el texto de salida (consola + resumen de errores + resumen de gráficas) y, si corresponde, llama a `Graficador.mostrar(...)` por cada gráfica registrada. Al ejecutarse desde el hilo de eventos de JavaFX (el clic del botón), es seguro abrir ventanas nuevas directamente.
- `reportes()`: invoca `Reportes.generar(ultimoEntorno, new File("reportes"))` y abre cada archivo generado con `getHostServices().showDocument(...)`.

### 5.2 `Lanzador`

Clase auxiliar cuyo único propósito es servir como punto de entrada real de la aplicación:

```java
public class Lanzador {
    public static void main(String[] args) {
        EditorApp.main(args);
    }
}
```

Existe porque el lanzador estándar de Java verifica si la clase con `main` extiende `javafx.application.Application`; si es así y JavaFX no está declarado en el *module-path*, aborta con el error "JavaFX runtime components are missing". Como `Lanzador` **no** extiende `Application`, ese chequeo no se activa, y la aplicación puede ejecutarse con el classpath plano que arma IntelliJ a partir del `pom.xml`. La alternativa equivalente por terminal es `mvn clean javafx:run`, que sí configura correctamente el *module-path* de JavaFX.

### 5.3 `Graficador`

Traduce cada instancia de `Grafica` (tipo + mapa de atributos ya validados por `Entorno.validarGrafica`) a un `javafx.scene.chart.Chart` concreto y lo muestra en una ventana (`Stage`) propia:

| Tipo de gráfica | Chart de JavaFX | Ejes / datos |
|---|---|---|
| `graphBar` | `BarChart<String, Number>` | `CategoryAxis` (ejeX) / `NumberAxis` (ejeY) |
| `graphPie` | `PieChart` | pares `label` / `values` |
| `graphLine` | `LineChart<String, Number>` | mismos ejes que `graphBar` |
| `Histogram` | `BarChart<String, Number>` de frecuencias | eje X = valores agrupados, eje Y = frecuencia (vía `Operaciones.frecuencias`) |

Como los atributos ya fueron validados semánticamente antes de llegar a esta clase (`Entorno.validarGrafica`, ver sección 6.5), los *casts* dentro de `Graficador` (`(String)`, `(ArrayList<?>)`, `(Double)`) se hacen sin verificación adicional: la responsabilidad de esta clase es únicamente dibujar.

## 6. Paquete `dataforge.interprete`

### 6.1 `Entorno`

Contiene todo el estado mutable de una ejecución:

- `LinkedHashMap<String, Simbolo> simbolos` — tabla de símbolos. Las claves se normalizan a minúsculas (`id.toLowerCase()`) porque el lenguaje es *case insensitive* también para los identificadores, pero cada `Simbolo` conserva el nombre original tal como lo escribió el usuario, para mostrarlo en los reportes.
- `StringBuilder consola` — texto acumulado de salida.
- `List<RegistroError> errores` — errores léxicos, sintácticos y semánticos, en el orden en que se detectaron.
- `List<Grafica> graficas` — gráficas registradas (solo las que incluyeron `EXEC`).
- `List<Object[]> tokens` — lexema, tipo, línea y columna de cada token reconocido por el lexer (para el reporte de tokens).

Métodos principales:

- `declararVariable(String id, String tipo, Object valor, int l, int c)` / `declararArreglo(...)`: validan que el identificador no exista ya en la tabla (si existe, error semántico de redeclaración) y que el valor recibido corresponda al tipo declarado (`double` ⇔ `Double`, `char[]` ⇔ `String`); si la validación pasa, insertan el `Simbolo` correspondiente.
- `valorDe(String id, int l, int c)` / `valorArreglo(String idArr, int l, int c)`: consultan la tabla de símbolos; si el identificador no existe o no es de la categoría esperada (variable vs. arreglo), registran un error semántico y devuelven `null`.
- `imprimir(ArrayList<Object> exprs)`: implementa `console::print`, uniendo los valores con `", "`.
- `imprimirColumna(Object titulo, ArrayList<Object> arr)`: implementa `console::column`, imprimiendo el arreglo como una tabla de una sola columna.
- `registrarGrafica(String tipo, ArrayList<?> attrs, int l, int c)`: recorre los atributos del bloque de gráfica aplicando la regla **"la última instrucción gana"** (si un mismo atributo se asigna más de una vez, el `LinkedHashMap.put` sobrescribe el valor anterior); solo si el bloque contuvo `EXEC` se valida (`validarGrafica`) y, de ser válida, se agrega a la lista de gráficas.
- `formatear(Object v)` (estático) — formato de **consola**: un `Double` entero se muestra sin decimales (`15.0` → `"15"`), uno fraccionario conserva sus decimales (`15.7` → `"15.7"`); las cadenas se muestran **sin** comillas.
- `valorReporte(Object v)` (estático) — formato de **reporte de símbolos** (§6.3 del enunciado): las cadenas se muestran **con** comillas (`"Hola Mundo"`) y los arreglos se formatean elemento por elemento con la misma regla de enteros/decimales que `formatear`. Estos dos formateadores son intencionalmente distintos y **no deben combinarse**: la consola y el reporte de símbolos tienen convenciones de formato diferentes, ambas tomadas del enunciado.

### 6.2 `Operaciones`

Implementa la semántica de las operaciones aritméticas (`SUM`, `RES`, `MUL`, `DIV`, `MOD`) y estadísticas (`Media`, `Mediana`, `Moda`, `Varianza`, `Max`, `Min`), todas restringidas a operandos de tipo `double` (§5.7-5.8 del enunciado). Ambos métodos de entrada (`aritmetica(...)` y `estadistica(...)`) siguen la misma convención: si algún operando llega `null`, devuelven `null` sin reportar un nuevo error (propagación descrita en la sección 4.4); si el tipo es incorrecto, registran el error semántico correspondiente en el `Entorno` recibido como parámetro.

`DIV` y `MOD` tienen una validación adicional: si el divisor es `0`, se registra el error semántico correspondiente (`"división entre cero"` / `"módulo entre cero"`) y se devuelve `null`, en vez de dejar propagar el `NaN` silencioso que produce Java al hacer `x % 0.0` (a diferencia de la división entera, el operador `%` de `double` no lanza excepción: sin este chequeo el `NaN` se habría colado hasta la consola o un reporte, violando la convención de la sección 4.4).

El método `frecuencias(ArrayList<Double> datos)` ordena los datos y cuenta las repeticiones de cada valor en un `LinkedHashMap<Double, Integer>` (por eso el orden de iteración resultante es ascendente); lo usan tanto la tabla de consola del histograma (`Entorno.tablaHistograma`) como su representación gráfica (`Graficador.histograma`).

### 6.3 `Simbolo`

Clase de datos inmutable que representa una entrada de la tabla de símbolos: `nombre` (tal como lo escribió el usuario), `categoria` (`"variable"` o `"arreglo"`), `tipo` (`"double"` o `"char[]"`), `valor` (`Double`, `String` o `ArrayList<Object>`), y la posición de declaración (`linea`, `columna`, en base 1).

### 6.4 `Grafica`

Clase de datos inmutable: `tipo` (uno de `"graphBar"`, `"graphPie"`, `"graphLine"`, `"Histogram"`) y `atributos` (`LinkedHashMap<String, Object>` con los atributos ya resueltos y validados). El método `Graficador.mostrar` es el único consumidor de esta clase.

### 6.5 Validación semántica de gráficas

El método privado `Entorno.validarGrafica(String tipo, LinkedHashMap<String, Object> m, int l, int c)` centraliza el chequeo de que cada tipo de gráfica reciba los atributos que necesita, con el tipo correcto:

- `graphBar` / `graphLine`: requieren `titulo`, `titulox`, `tituloy` (cadenas), `ejex` (lista de `char[]`) y `ejey` (lista de `double`) de igual tamaño.
- `graphPie`: requiere `titulo` (cadena), `label` (lista de `char[]`) y `values` (lista de `double`) de igual tamaño.
- `Histogram`: requiere `titulo` (cadena) y `values` (lista de `double`).

Los atributos de gráfica (`titulo`, `ejeX`, `values`, etc.) **no son palabras reservadas** del lenguaje: en la gramática son identificadores comunes (`ID`), porque el enunciado también los usa como nombre de variable en otros contextos. Su validez como atributo de gráfica se revisa en tiempo de ejecución, no en el analizador léxico ni sintáctico.

### 6.6 `RegistroError`

Registra un error con su `tipo` (`"Léxico"`, `"Sintáctico"` o `"Semántico"`), `descripcion`, y posición (`linea`, `columna`, ya en base 1). Su `toString()` produce el formato `[tipo] descripción (línea l, columna c)`, usado tanto en la consola como en el reporte de errores.

## 7. Paquete `dataforge.reportes`

### 7.1 `Reportes`

Clase estática que genera los tres archivos HTML exigidos por el enunciado (§6), siempre a partir del `Entorno` de la **última** ejecución:

- `tokens(Entorno ent)` — recorre `ent.getTokens()` (lexema original, id de tipo de token, línea, columna). El **nombre** de cada tipo de token (por ejemplo, `VAR`, `NUMERO`, `ID`) se obtiene por **reflexión** sobre los campos públicos `static int` de la clase generada `dataforge.analisis.sym` (`nombresTokens()`), en vez de mantener un `switch` manual que debería actualizarse cada vez que cambia la gramática.
- `errores(Entorno ent)` — recorre `ent.getErrores()`, mostrando las tres familias de error (léxico, sintáctico, semántico) en una sola tabla.
- `simbolos(Entorno ent)` — recorre `ent.getSimbolos().values()`, usando `Entorno.valorReporte(...)` para el formato de los valores (sección 6.1).

Los tres métodos comparten la plantilla HTML (`pagina(String titulo, String[] columnas, String filas)`), con CSS embebido para que los archivos generados sean autocontenidos y puedan abrirse sin conexión a internet. El método `esc(String s)` escapa los caracteres `&`, `<` y `>` en los lexemas y descripciones antes de insertarlos en la tabla, porque algunos lexemas del lenguaje (por ejemplo, `<-`) romperían el HTML si se insertaran sin escapar.

## 8. Gramática y archivos generados

La gramática formal del lenguaje, en notación BNF, se encuentra en `docs/gramatica.txt`. Es un documento independiente escrito para ser leído por una persona (siguiendo la convención del enunciado), y **no** es una copia del archivo `src/main/cup/parser.cup`; sin embargo, ambos son consistentes entre sí: cada producción de `docs/gramatica.txt` corresponde a una o más producciones reales de `parser.cup`.

Los archivos `Lexer.java`, `Parser.java` y `sym.java` bajo `target/generated-sources/` se regeneran en cada `mvn compile` a partir de `src/main/jflex/Lexer.flex` y `src/main/cup/parser.cup`. No se versiona ni se edita manualmente ningún archivo dentro de `target/`.

## 9. Casos de prueba

El proyecto incluye cuatro programas de prueba en `entradas/`, verificados contra la versión actual del intérprete:

| Archivo | Propósito | Resultado verificado |
|---|---|---|
| `ejemplo1.df` | Variables, arreglo con `SUM` anidada, `Media`, `console::print` / `console::column`, ambos tipos de comentario | 4 símbolos declarados; consola con `"la media es:, 3.5"` y la columna de `@datos` |
| `ejemplo2.df` | Las 4 gráficas (`graphBar`, `graphPie`, `graphLine`, `Histogram`), anidación de funciones (`DIV(SUM(Max, Min), 2)`), atributo con arreglo por nombre (`@actividades`) | 3 símbolos, 4 gráficas registradas (con `TestInterprete` solo se listan por consola; el dibujado real con JavaFX Charts requiere el editor — `Lanzador`/`mvn javafx:run`), tabla de frecuencias del histograma en consola |
| `ejemplo3_errores.df` | 5 errores semánticos intencionales (léxica y sintácticamente correcto) | Los 5 errores se reportan con línea y columna, y la ejecución continúa hasta el final (`bueno` se imprime correctamente pese a los errores previos) |
| `ejemplo4_mixto.df` | Los tres tipos de error (léxico, sintáctico y semántico) en un mismo programa | Con el pipeline completo (`Interprete.ejecutar`, usado por la GUI) los 3 errores quedan en la lista del `Entorno` y en el reporte; las variables `bueno` y `precio` quedan correctamente declaradas y se imprimen al final |

Estas ejecuciones fueron reproducidas desde línea de comandos con la clase `dataforge.TestInterprete` (`mainClass` alternativo al de la GUI, útil para depurar el intérprete sin levantar JavaFX), invocada como `exec:java -Dexec.mainClass=dataforge.TestInterprete -Dexec.args="entradas/<archivo>"`. Con esta clase, el error léxico de `ejemplo4_mixto.df` no aparece en la lista `Entorno.getErrores()` — se imprime aparte, por `stderr` (ver nota en la sección 10): la reproducción fiel de "los 3 errores en un mismo reporte" requiere el pipeline con `Interprete.ejecutar`, es decir, correr el editor (`Lanzador`) y generar los reportes.

## 10. Errores conocidos y notas de mantenimiento

- Olvidar la directiva `%public` en `Lexer.flex` genera una clase `Lexer` *package-private*, que no puede usarse desde otros paquetes (`dataforge.gui`, `dataforge.interprete`); el error de compilación es `'Lexer' is not public...`.
- Omitir `%cup` en `Lexer.flex` produce un lexer que no genera objetos `Symbol` compatibles con CUP.
- El orden de las reglas léxicas importa: ante un empate en la longitud del match, JFlex prioriza la **primera** regla que aparece en el archivo — las palabras reservadas deben declararse antes que el patrón general `{Id}`.
- Las declaraciones `non terminal` en `parser.cup` deben ir **todas** antes de las producciones, y deben usar tipos *raw* (`ArrayList`, no `ArrayList<Object>` ni arreglos).
- Si IntelliJ subraya en rojo las clases `Parser` o `sym`, normalmente basta con correr Maven **compile** y luego **Reload All Maven Projects**.
- Al ejecutar la aplicación desde consola en Windows, la salida de caracteres acentuados puede mostrarse incorrectamente si la consola no está configurada en UTF-8; esto es un problema de la consola del sistema operativo, no del intérprete (los archivos `.df` y los reportes HTML se leen/escriben en UTF-8 explícitamente).
- `TestInterprete.java` (y `TestParser.java`) construyen el `Lexer` y el `Parser` por separado y **no** ejecutan `lexer.entorno = parser.entorno;` como sí hace `Interprete.ejecutar(...)`. Consecuencia práctica: al depurar por esa vía, el lexer no registra tokens (`Entorno.getTokens()` queda vacío) y los errores léxicos se imprimen directo a `stderr` en vez de sumarse a `Entorno.getErrores()`. No es un bug del intérprete — es una limitación conocida de estas dos clases de prueba mínimas; el pipeline real (usado por la GUI vía `Interprete.ejecutar`) sí conecta ambos objetos correctamente.
