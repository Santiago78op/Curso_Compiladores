# Manual Técnico — ConjAnalyzer

Universidad de San Carlos de Guatemala
Facultad de Ingeniería — Escuela de Ciencias y Sistemas
Organización de Lenguajes y Compiladores 1
Proyecto 1 — ConjAnalyzer

## Índice

1. Introducción
2. Lenguaje de programación y herramientas utilizadas
3. Estructura del proyecto (Maven)
4. Arquitectura general
5. Paquete `conjanalyzer.analisis` (léxico y sintáctico)
6. Paquete `conjanalyzer.interprete` (evaluación y simplificación)
7. Paquete `conjanalyzer.reportes`
8. Paquete `conjanalyzer.gui`
9. Decisiones de diseño no evidentes
10. Compilación y ejecución
11. Mantenimiento futuro

---

## 1. Introducción

Este manual describe la arquitectura técnica del proyecto ConjAnalyzer, con el objetivo de que un desarrollador distinto al autor original pueda comprender su estructura, dar mantenimiento al código o extender sus funcionalidades. Se documentan el lenguaje y herramientas empleadas, los paquetes reales del código fuente, las clases y métodos más relevantes (citando sus firmas reales) y las decisiones de diseño que no son evidentes a simple vista.

## 2. Lenguaje de programación y herramientas utilizadas

| Herramienta | Uso en el proyecto |
|---|---|
| **Java 17** | Lenguaje de implementación de todo el proyecto (`maven.compiler.source`/`target = 17` en `pom.xml`). |
| **JFlex** (`de.jflex:jflex-maven-plugin:1.9.1`) | Generador del analizador léxico a partir de la especificación `src/main/jflex/Lexer.flex`. Produce la clase `Lexer` en `target/generated-sources/jflex`. |
| **CUP** (`com.github.vbmacher:cup-maven-plugin:11b-20160615-3`) | Generador del analizador sintáctico (LALR) a partir de la especificación `src/main/cup/parser.cup`. Produce las clases `Parser` y `sym` en `target/generated-sources/cup`. |
| **java-cup-runtime** (`11b-20160615`) | Biblioteca en tiempo de ejecución de CUP; aporta la clase `Symbol` que el lexer retorna al parser. |
| **JavaFX** (`org.openjfx:javafx-controls:21.0.4`, `javafx-maven-plugin:0.0.8`) | Interfaz gráfica del entorno de trabajo (editor, consola, panel de Venn). |
| **Gson** (`com.google.code.gson:gson:2.11.0`) | Serialización del reporte de simplificación (`simplificacion.json`) al formato JSON exacto exigido por el enunciado. |
| **Apache Maven** | Gestor de dependencias y de construcción del proyecto; orquesta la generación de código (JFlex/CUP), la compilación y el empaquetado del `.jar`. |

## 3. Estructura del proyecto (Maven)

```
ConjAnalyzer/
├── pom.xml
├── docs/
│   ├── gramatica.txt          (gramática BNF entregable)
│   ├── ManualUsuario.md
│   └── ManualTecnico.md
├── entradas/                  (casos de prueba .ca)
├── reportes/                  (salida generada: tokens.html, errores.html,
│                                operaciones.html, simplificacion.json)
└── src/main/
    ├── jflex/Lexer.flex
    ├── cup/parser.cup
    └── java/conjanalyzer/
        ├── TestInterprete.java        (prueba de consola, sin GUI)
        ├── analisis/                  (paquete de destino de Lexer/Parser/sym generados)
        ├── interprete/                (evaluación, tabla de conjuntos, simplificación)
        ├── reportes/                  (generación de HTML y JSON)
        └── gui/                       (interfaz JavaFX)
```

Los plugins `jflex-maven-plugin` y `cup-maven-plugin` se ejecutan en la fase `generate-sources`, de modo que `Lexer.java`, `Parser.java` y `sym.java` quedan disponibles antes de la compilación de las clases propias del proyecto. **Estas tres clases generadas nunca se editan manualmente**: cualquier cambio en el análisis léxico o sintáctico debe hacerse en `Lexer.flex` o `parser.cup`.

## 4. Arquitectura general

El proyecto se organiza en 4 paquetes con responsabilidades claramente separadas:

1. **`conjanalyzer.analisis`**: contiene únicamente el código generado por JFlex y CUP (`Lexer`, `Parser`, `sym`) más el paquete declarado en las cabeceras de `Lexer.flex`/`parser.cup`.
2. **`conjanalyzer.interprete`**: el núcleo semántico. Recibe las llamadas desde las acciones del `.cup` y mantiene todo el estado de una ejecución (`Entorno`), la representación de los conjuntos y operaciones, y el motor de simplificación algebraica.
3. **`conjanalyzer.reportes`**: transforma el `Entorno` de la última ejecución en los artefactos de salida exigidos por el enunciado (HTML y JSON).
4. **`conjanalyzer.gui`**: la interfaz gráfica de usuario, construida sobre JavaFX, que conecta el editor de texto con el intérprete y presenta los resultados (consola y diagramas de Venn).

El flujo de ejecución completo es:

```
código fuente (String)
   → Lexer (JFlex)       — produce tokens, registra cada uno en el Entorno
   → Parser (CUP, LALR)  — valida la gramática y, en sus acciones, llama a Entorno
                            para construir conjuntos, árboles de operación y evaluarlos
   → Entorno              — queda con: conjuntos, operaciones (ya evaluadas y
                            simplificadas), consola, errores y tokens
   → Reportes / JsonSalida — a demanda, generan los archivos de salida
   → GUI (DiagramaVenn)    — a demanda, dibuja el resultado
```

La ejecución sintáctica y la evaluación semántica **no están separadas en dos pasadas**: la gramática es S-atribuida (traducción dirigida por la sintaxis) y las acciones del `.cup` invocan directamente a los métodos del `Entorno` en el momento en que cada producción se reduce. La única estructura intermedia que se materializa como árbol es `NodoOperacion`, exclusiva de las operaciones definidas con `OPERA` (ver sección 6.2).

## 5. Paquete `conjanalyzer.analisis` (léxico y sintáctico)

### 5.1 Especificación léxica (`Lexer.flex`)

Declarado con las directivas `%class Lexer`, `%public`, `%unicode`, `%cup`, `%line`, `%column`. El lenguaje **no** utiliza `%ignorecase`: es sensible a mayúsculas y minúsculas por requerimiento del enunciado (sección 4.1).

Tokens reconocidos:

| Categoría | Tokens / patrón |
|---|---|
| Palabras reservadas | `CONJ`, `OPERA`, `EVALUAR` |
| Operadores de conjunto | `U` (unión), `&` (intersección), `^` (complemento), `-` (diferencia) |
| Símbolos estructurales | `->` (FLECHA), `~` (VIRGULILLA), `{`, `}`, `(`, `)`, `:`, `;`, `,` |
| Con patrón | `ID` = `{Letra}({Letra}|{Digito})*`, `NUMERO` = `{Digito}+`, `SIMBOLO` = comodín `[\x21-\x7E]` (cualquier ASCII imprimible no reconocido por las reglas anteriores) |
| Descartados | comentario de línea `"#"[^\r\n]*`, comentario multilínea `"<!"~"!>"`, blancos `[ \t\r\n]+` |
| Error léxico | `[^]` — cualquier carácter fuera del universo ASCII 33-126; se registra el error y el análisis continúa |

El código de usuario incluido en la sección 1 del `.flex` declara un campo público `Entorno entorno`, asignado externamente por `Interprete.ejecutar(...)`, que permite al lexer:
- Registrar cada token reconocido mediante `entorno.registrarToken(yytext(), type, yyline, yycolumn)` (para el reporte de tokens).
- Reportar los errores léxicos al `Entorno` en lugar de imprimirlos por `stderr`.

### 5.2 Especificación sintáctica (`parser.cup`)

Declara los terminales `CONJ, OPERA, EVALUAR, UNION, INTERSECCION, COMPLEMENTO, DIFERENCIA, FLECHA, VIRGULILLA, LLAVE_IZQ, LLAVE_DER, PAR_IZQ, PAR_DER, DOS_PUNTOS, PUNTO_COMA, COMA` (sin valor semántico) y `ID, NUMERO, SIMBOLO` (terminales con valor `String`).

Los no terminales `operacion` (tipo `NodoOperacion`), `elemento` (tipo `String`) y `lista_elem` (tipo `ArrayList`, sin genéricos, por restricción de CUP) sintetizan valores que se usan en las producciones superiores. El bloque `parser code {: ... :}` declara `public Entorno entorno = new Entorno();` y sobrescribe `syntax_error(Symbol s)` para registrar los errores sintácticos en el `Entorno` con la posición del símbolo que causó el conflicto.

La recuperación de errores sintácticos se implementa mediante el patrón de modo pánico del Dragón (sección 4.8.3): la producción

```
sentencia ::= def_conj | def_opera | evaluar_stmt | error PUNTO_COMA ;
```

permite a CUP descartar símbolos hasta poder desplazar el terminal especial `error`, y luego descartar tokens hasta encontrar el siguiente `;` — el análisis continúa y el reporte de errores acumula todas las fallas sintácticas de una misma ejecución, en lugar de detenerse en la primera.

## 6. Paquete `conjanalyzer.interprete` (evaluación y simplificación)

### 6.1 `Entorno`

Clase que concentra todo el estado de una ejecución. Constantes: `UNIVERSO_MIN = 33`, `UNIVERSO_MAX = 126` (ASCII `!` a `~`). Estructuras internas: `Set<Character> universo`, `LinkedHashMap<String, Conjunto> conjuntos`, `LinkedHashMap<String, Operacion> operaciones`, `StringBuilder consola`, `List<RegistroError> errores`, `List<Object[]> tokens`.

Métodos principales invocados desde las acciones del `.cup`:

```java
public void definirConjuntoRango(String id, String a, String b, int l, int c)
public void definirConjuntoLista(String id, ArrayList<String> elementos, int l, int c)
public void definirOperacion(String id, NodoOperacion arbol, int l, int c)
public void evaluar(ArrayList<String> datos, String operacion, int l, int c)
public void error(String tipo, String descripcion, int linea, int columna)
public void registrarToken(String lexema, int tipo, int linea, int columna)
```

`definirOperacion` es el punto donde convergen validación, evaluación y simplificación: valida que todos los conjuntos referenciados por el árbol existan (`arbol.referencias(refs)`), evalúa el árbol contra el universo y los conjuntos definidos (`arbol.evaluar(mapa, universo)`), invoca al `Simplificador` sobre el árbol original, y almacena el resultado como un objeto `Operacion`.

Las claves de los mapas `conjuntos` y `operaciones` **no se normalizan** (a diferencia de DataForge, proyecto hermano del mismo curso): se respeta estrictamente el requisito de sensibilidad a mayúsculas y minúsculas del enunciado.

### 6.2 `NodoOperacion`: el árbol de una operación

Representa el árbol prefijo de una definición `OPERA`. Es la única estructura de tipo árbol de sintaxis abstracta del proyecto (acotada a las operaciones, no al programa completo). Un nodo puede ser hoja (`op = '\0'`, referencia a un conjunto por nombre), unario (`op = '^'`, complemento) o binario (`op` en `{'U','&','-'}`).

```java
public static NodoOperacion hoja(String nombre)
public static NodoOperacion unario(char op, NodoOperacion hijo)
public static NodoOperacion binario(char op, NodoOperacion a, NodoOperacion b)
public Set<Character> evaluar(Map<String, Set<Character>> conjuntos, Set<Character> universo)
public boolean pertenenciaRegion(Map<String, Boolean> region)
public void referencias(Set<String> acumulador)
public String toPrefijo()
public NodoOperacion copia()
```

- `evaluar(...)` recorre el árbol en postorden y calcula el `Set<Character>` resultante: unión (`addAll`), intersección (`retainAll`), diferencia (`removeAll`) y complemento (universo menos el conjunto evaluado del único hijo).
- `pertenenciaRegion(Map<String,Boolean>)` evalúa el árbol como función booleana sobre una asignación hipotética de pertenencia a cada conjunto base; es la función que usa `DiagramaVenn` para decidir, píxel por píxel, si una región debe sombrearse, sin necesidad de recorrer el universo completo.
- `toPrefijo()` serializa el árbol de vuelta a la notación prefija del lenguaje (usada en los reportes y en el JSON de simplificación).
- `copia()` produce una copia profunda; la usa el `Simplificador` para no mutar el árbol original de la operación.

### 6.3 `Simplificador`: motor de simplificación algebraica (sección 7 del enunciado)

```java
public ResultadoSimplificacion simplificar(NodoOperacion original)
```

Aplica repetidamente, en una pasada bottom-up (`paso(NodoOperacion n)`), a lo sumo una regla por nodo, hasta que el árbol serializado (`toPrefijo()`) deja de cambiar entre una pasada y la siguiente, con un límite de seguridad de 100 iteraciones (`MAX_ITER`). Las leyes implementadas como transformación real son:

- **Ley del doble complemento**: `^^X → X`.
- **Leyes de DeMorgan**, con una condición de guarda: `^(X U Y) → ^X & ^Y` (y su dual con `&`) se aplica únicamente cuando `X` o `Y` ya es a su vez un complemento, de modo que la transformación genere un `^^` que luego se cancela por doble complemento en la siguiente pasada. Esta guarda evita aplicar DeMorgan en casos donde no produce ninguna simplificación neta.
- **Propiedades idempotentes**: `X U X → X`, `X & X → X`.
- **Propiedades de absorción**: `X U (X & Y) → X` (y las variantes conmutadas).
- **Propiedades distributivas** (`distributiva(char op, NodoOperacion a, NodoOperacion b)`): `(X & Y) U (X & Z) → X & (Y U Z)` y su dual `(X U Y) & (X U Z) → X U (Y & Z)`. Es decir, solo se aplica en el sentido que **factoriza** un término común (reduce la cantidad de hojas del árbol), nunca en el sentido inverso de "expandir" un factor (`X & (Y U Z) → (X & Y) U (X & Z)` aumenta el árbol y no simplifica nada) — la misma lógica de guarda que DeMorgan. Detecta el factor común entre los cuatro pares posibles de operandos (`a.izq`/`a.der` contra `b.izq`/`b.der`) usando `equivalentes(...)`, así que reconoce el factor aunque aparezca en cualquier posición o anidado.

Las **propiedades conmutativa y asociativa no se aplican como reglas de reescritura independientes** (reordenar o reagrupar el árbol no lo simplifica por sí solo, y podría producir ciclos); en su lugar, se usan como criterio de comparación estructural mediante el método privado `canon(NodoOperacion n)`, que produce una forma canónica de un subárbol: los operandos de cadenas de `U` o de `&` se aplanan recursivamente, se ordenan alfabéticamente y se eliminan duplicados. Dos subárboles son considerados equivalentes (método `equivalentes`) si su forma canónica coincide, lo que permite reconocer idempotencia, absorción o el factor común de la distributiva aunque los operandos aparezcan en distinto orden o con distinta agrupación. Cuando la equivalencia detectada no es trivial (los árboles no son textualmente idénticos), el motor añade también "Propiedades conmutativas" y/o "Propiedades asociativas" a la lista de leyes reportadas.

El resultado se encapsula en `ResultadoSimplificacion(List<String> leyes, NodoOperacion simplificado, boolean seSimplifico)`.

### 6.4 Otras clases del paquete

- **`Conjunto`**: registro inmutable de un conjunto definido con `CONJ` (nombre, `Set<Character>` de elementos, cadena de definición legible, línea y columna de declaración).
- **`Operacion`**: registro inmutable de una operación definida con `OPERA` (nombre, árbol `NodoOperacion`, resultado ya evaluado, conjuntos base referenciados y el `ResultadoSimplificacion` correspondiente).
- **`RegistroError`**: estructura de un error (tipo — Léxico, Sintáctico o Semántico —, descripción, línea y columna en base 1).
- **`Interprete`**: fachada estática `public static Entorno ejecutar(String codigo)` que instancia `Lexer` y `Parser` sobre un `StringReader`, conecta `lexer.entorno = parser.entorno` para que el lexer pueda registrar tokens y errores léxicos, invoca `parser.parse()` dentro de un bloque `try/catch` que evita que una excepción de un error sintáctico irrecuperable propague hasta detener la interfaz gráfica, y retorna el `Entorno` resultante.

## 7. Paquete `conjanalyzer.reportes`

### 7.1 `Reportes`

```java
public static File[] generar(Entorno ent, File carpeta) throws Exception
```

Genera tres archivos HTML autocontenidos (con CSS embebido en cada página, sin dependencias externas): `tokens.html` (reporte de tokens, sección 5.2 del enunciado), `errores.html` (reporte de errores, sección 5.3) y `operaciones.html` (reporte adicional de conjuntos y operaciones definidos, con su resultado o estado de simplificación).

El nombre de cada tipo de token se obtiene por **reflexión** sobre la clase `sym` generada por CUP (`nombresTokens()` recorre `sym.class.getFields()` y construye un mapa `id numérico → nombre del campo`), evitando mantener un `switch` manual que se desincronizaría cada vez que cambien los terminales declarados en el `.cup`. Todos los lexemas y descripciones se escapan (método `esc(String s)`) para que caracteres como `&` o `<` (parte legítima del lenguaje) no rompan la tabla HTML resultante.

### 7.2 `JsonSalida`

```java
public static String construir(Entorno ent)
public static File generar(Entorno ent, File carpeta) throws Exception
```

Construye, usando Gson (`new GsonBuilder().setPrettyPrinting().disableHtmlEscaping().create()`), el objeto JSON exigido por la sección 5.4 del enunciado: para cada operación con `seSimplifico == true`, un objeto con las claves `"leyes"` (arreglo) y `"conjunto simplificado"` (la expresión prefija ya simplificada); para las que no se pudieron simplificar, el valor asociado es el literal `"No se puede simplificar la operacion"`. Se escribe en `reportes/simplificacion.json`.

## 8. Paquete `conjanalyzer.gui`

- **`EditorApp extends Application`**: construye la interfaz completa por código (sin FXML). Estructura del grafo de escena: `BorderPane` raíz, con una `HBox` de botones en la parte superior y, en el centro, un `SplitPane` horizontal que separa el área de trabajo (izquierda: `SplitPane` vertical con `TabPane` de edición y `TextArea` de consola no editable) del panel de diagramas de Venn (derecha). Mantiene el `Entorno` de la última ejecución (`ultimoEntorno`) para que los reportes y el panel de Venn siempre reflejen el análisis más reciente, cumpliendo el requisito de la sección 5 del enunciado.
- **`Lanzador`**: clase con un método `main` que no extiende `Application`, cuyo único propósito es invocar `EditorApp.main(args)`. Existe porque el lanzador de Java revisa si la clase con el método `main` extiende `Application`; si es así y JavaFX no está declarado en el module-path, aborta con el error "JavaFX runtime components are missing". Ejecutar `Lanzador` en su lugar evita esa validación cuando se corre con el classpath plano que arma un IDE.
- **`DiagramaVenn`**: método estático `public static VBox crear(Operacion op)`, que construye el panel visual (título, `Canvas` y etiquetas de texto) para una operación. El sombreado se calcula píxel por píxel: para cada punto del lienzo se determina en qué círculos (conjuntos base) cae, y se evalúa `op.arbol.pertenenciaRegion(region)` para decidir si ese punto pertenece al conjunto resultado. Soporta de 1 a 3 conjuntos base; con más de 3, no existe una disposición geométrica razonable de círculos y se muestra el resultado únicamente en forma de texto.

## 9. Decisiones de diseño no evidentes

1. **Alcance del árbol de sintaxis abstracta**: el proyecto no construye un AST del programa completo; solo las operaciones (`OPERA`) se representan como árbol (`NodoOperacion`), porque son la única construcción del lenguaje que necesita reescritura estructural (simplificación) y evaluación diferida. Las sentencias `CONJ` y `EVALUAR` se ejecutan directamente en las acciones semánticas del `.cup`, en el mismo estilo de traducción dirigida por la sintaxis que el proyecto hermano DataForge.
2. **Sensibilidad a mayúsculas y minúsculas real**: a diferencia de otros proyectos del mismo curso, aquí no se aplica ninguna normalización de identificadores; se declaró explícitamente en el `.flex` la ausencia de la directiva `%ignorecase`.
3. **Guarda de DeMorgan**: aplicar la ley de DeMorgan sin restricción puede transformar árboles sin ganancia neta de simplicidad; el motor solo la aplica cuando al menos uno de los operandos ya es un complemento, garantizando que el resultado se reduzca en un paso posterior por la ley del doble complemento.
4. **Comparación estructural en vez de textual**: para que las propiedades idempotente y de absorción reconozcan operandos en distinto orden o agrupación, el `Simplificador` compara los subárboles por una forma canónica (aplanado + orden + deduplicación de cadenas de `U`/`&`) en lugar de comparar directamente su representación en texto.
5. **Discrepancia documentada del enunciado**: el ejemplo de la sección 4.8 del enunciado presenta una salida de consola que no corresponde a la semántica real de la operación descrita (evalúa como "exitoso" un elemento que, dada la definición de los conjuntos involucrados, no pertenece al resultado). El código respeta la semántica correcta (pertenencia real al conjunto resultante evaluado) y documenta esta discrepancia en el propio código fuente (`Entorno.evaluar`).
6. **Entorno nuevo por ejecución**: cada llamada a `Interprete.ejecutar(...)` crea una instancia nueva de `Entorno`, de modo que los reportes y diagramas nunca mezclan información de ejecuciones distintas, cumpliendo el requisito explícito de la sección 5 del enunciado.

## 10. Compilación y ejecución

Desde IntelliJ IDEA (Maven embebido en el IDE):

- Regenerar el lexer y el parser y compilar: acción `compile` del ciclo de vida de Maven (regenera `Lexer`, `Parser` y `sym` a partir de los archivos `.flex`/`.cup`).
- Ejecutar la interfaz gráfica: ejecutar la clase `Lanzador` (nunca `EditorApp` directamente, para evitar el error de módulos de JavaFX faltantes).
- Probar el intérprete por consola: ejecutar la clase `TestInterprete`, opcionalmente con el argumento de la ruta de un archivo `.ca`.
- Generar el archivo ejecutable: acción `package` del ciclo de vida de Maven.

Desde línea de comandos (equivalente):

```
mvn clean compile
mvn javafx:run
mvn compile exec:java -Dexec.mainClass=conjanalyzer.TestInterprete -Dexec.args="entradas/ejemplo1.ca"
mvn clean package
```

## 11. Mantenimiento futuro

- Cualquier cambio en los tokens del lenguaje debe hacerse en `src/main/jflex/Lexer.flex`; cualquier cambio en la gramática, en `src/main/cup/parser.cup`. Las clases `Lexer`, `Parser` y `sym` se regeneran automáticamente en cada `compile` y no deben editarse a mano.
- Si se agregan nuevas propiedades de simplificación (sección 7 del enunciado), deben implementarse como una nueva rama dentro de `Simplificador.paso(...)`, siguiendo el mismo patrón: aplicar la regla solo cuando produce una reducción real del árbol, y registrar el nombre de la ley aplicada en el conjunto `leyes`.
- Si se agregan nuevos tipos de reporte, seguir el patrón de `Reportes.pagina(...)`: una plantilla HTML con CSS embebido, sin dependencias externas, y con escape de caracteres especiales en los datos insertados dinámicamente.
