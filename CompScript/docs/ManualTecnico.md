# Manual Técnico — CompScript

Universidad de San Carlos de Guatemala
Facultad de Ingeniería
Escuela de Ciencias y Sistemas
Organización de Lenguajes y Compiladores 1
Vacaciones de Diciembre 2024

## 1. Introducción

Este documento describe la arquitectura interna del intérprete CompScript: el lenguaje y las herramientas utilizadas, la organización en paquetes, las clases y métodos principales de cada fase del análisis, y las decisiones de diseño no evidentes a partir de la sola lectura del código fuente. Está dirigido a un lector con conocimientos de compiladores equivalentes a los del curso.

## 2. Lenguaje y herramientas

| Herramienta | Uso en el proyecto |
|---|---|
| **Java 17** | Lenguaje de implementación de la totalidad del intérprete (`maven.compiler.source`/`target` = 17 en `pom.xml`). |
| **JFlex** (`jflex-maven-plugin` 1.9.1) | Generación del analizador léxico a partir de la especificación `src/main/jflex/Lexer.flex`. Genera `Lexer.java` en `target/generated-sources/jflex/compscript/analisis/`. |
| **CUP** (`cup-maven-plugin` 11b-20160615-3, runtime `java-cup-runtime` 11b-20160615) | Generación del analizador sintáctico ascendente (LALR) a partir de la especificación `src/main/cup/parser.cup`. Genera `Parser.java` y `sym.java` en `target/generated-sources/cup/compscript/analisis/`. |
| **JavaFX** (`javafx-controls` 21.0.4, `javafx-maven-plugin` 0.0.8) | Interfaz gráfica del editor de trabajo, construida directamente por código (sin archivos FXML). |
| **Maven** | Gestión de dependencias y orquestación del ciclo de generación de código (fase `generate-sources`) y de ejecución. |

Ambos generadores (JFlex y CUP) se ejecutan automáticamente en la fase `generate-sources` de Maven; el código que producen no debe modificarse manualmente, dado que se sobrescribe en cada compilación. Toda modificación al lenguaje se realiza sobre los archivos fuente `Lexer.flex` y `parser.cup`.

## 3. Arquitectura general

El proyecto se organiza en cinco paquetes bajo `compscript`:

```
compscript
├── analisis      (GENERADO: Lexer, Parser, sym — no editar directamente)
├── ast           (A.java: definición y ejecución del árbol de sintaxis abstracta)
├── interprete    (estado en tiempo de ejecución: entornos, tipos, valores, operaciones)
├── gui           (EditorApp, Lanzador: interfaz gráfica)
└── reportes      (Reportes: generación de los reportes HTML)
```

Adicionalmente, `compscript.TestInterprete` es una clase de prueba por consola que ejecuta el pipeline completo sin interfaz gráfica.

### 3.1 Decisión de diseño central: análisis sintáctico dirigido a la construcción de un AST

A diferencia de otros proyectos del mismo curso que emplean gramáticas S-atribuidas (evaluación directa en las acciones semánticas del análisis sintáctico, sin árbol intermedio), CompScript construye explícitamente un **árbol de sintaxis abstracta** (paquete `compscript.ast`, clase `A`) durante el análisis sintáctico, y lo recorre en una fase posterior para realizar el análisis semántico y la ejecución.

Esta decisión responde a un requisito estructural del lenguaje: CompScript incluye construcciones de control de flujo con repetición y reevaluación (`while`, `for`, `do-while`, llamadas a función, recursión) que exigen poder ejecutar el mismo fragmento de sintaxis un número variable de veces, incluyendo cero veces. Una evaluación S-atribuida realizada directamente en las acciones del analizador sintáctico solo permite un recorrido único, en el orden fijo de las reducciones del analizador; no permite retornar a evaluar un subárbol ya reducido. Por ello, en cuanto la gramática de un lenguaje incluye ciclos o funciones con recursión, la construcción de un árbol intermedio se vuelve necesaria.

## 4. Analizador léxico (`Lexer.flex`, paquete generado `compscript.analisis`)

Especificación en `src/main/jflex/Lexer.flex`. Opciones relevantes de la directiva `%%` inicial: `%class Lexer`, `%public` (necesario para que la clase sea accesible desde otros paquetes), `%unicode`, `%cup` (integración con el analizador CUP), `%line` y `%column` (seguimiento de posición), `%ignorecase` (el lenguaje no distingue mayúsculas de minúsculas, incluyendo identificadores).

### 4.1 Categorías de tokens reconocidos

- **37 palabras reservadas**: tipos de dato (`int double bool char string void`), declaración (`let const true false`), control (`if else match default`), ciclos (`while for do`), transferencia (`break continue return`), tipos compuestos (`struct list cast as`), consola y listas (`console log push get set remove pop reverse`), funciones nativas (`round length tostring run_main`).
- **28 símbolos**: operadores de incremento/decremento e igualdad compuestos (`++ -- == != <= >= && || =>`), operadores simples (`+ - * / ^ $ % ! < > =`), agrupadores y puntuación (`( ) [ ] { } , ; : .`).
- **Tokens con patrón**: `ID` (identificador), `ENTERO`, `DECIMAL`, `CADENA`, `CARACTER`.
- **Elementos descartados**: comentarios de línea (`// ...`), comentarios de bloque (`/* ... */`) y espacios en blanco.

Las palabras reservadas están declaradas en la especificación antes que el patrón general `{Id}`, condición necesaria en JFlex para que, ante coincidencias de igual longitud, la primera regla declarada tenga prioridad sobre las siguientes.

### 4.2 Métodos auxiliares del lexer

```java
private Symbol symbol(int type)                 // token sin valor semantico asociado
private Symbol symbol(int type, Object value)    // token con valor semantico (numero, cadena, caracter)
private void registrar(int type)                 // notifica el token al Contexto para el reporte 6.1
private String procesarEscapes(String s)         // traduce \n \t \\ \" \' dentro de cadenas y caracteres
```

El lexer mantiene una referencia pública `contexto` (de tipo `compscript.interprete.Contexto`), asignada externamente por `Interprete.ejecutar(...)`. Cada token reconocido se reporta a través de `contexto.registrarToken(...)`, y cada carácter no reconocido (regla `[^]` al final de la especificación) se reporta como error léxico mediante `contexto.error("Lexico", ...)`, sin detener el análisis: el carácter se descarta y el reconocimiento continúa con el siguiente.

## 5. Analizador sintáctico (`parser.cup`, paquete generado `compscript.analisis`)

### 5.1 Precedencia y asociatividad de operadores

Dado que las expresiones de CompScript son infijas, la gramática de la producción `expr` es ambigua en su forma BNF pura y se desambigua mediante las declaraciones de precedencia de CUP, en orden de menor a mayor prioridad:

```
precedence left  OR;
precedence left  AND;
precedence right NOT;
precedence left  IGUAL_IGUAL, DIFERENTE, MENOR, MENOR_IGUAL, MAYOR, MAYOR_IGUAL;
precedence left  MAS, MENOS;
precedence left  POR, DIV, MOD;
precedence nonassoc POT, RAIZ;
precedence right UMENOS;
```

El terminal `UMENOS` es un pseudo-terminal sin patrón léxico propio, declarado únicamente para asignar a la negación unaria (`MENOS expr %prec UMENOS`) una precedencia distinta de la resta binaria, ambas asociadas al mismo símbolo léxico `MENOS`.

El conflicto de la sentencia `if` sin bloque `else` correspondiente (*dangling else*) se resuelve mediante el orden de las producciones de `if_stmt`, de forma que el analizador desplaza (*shift*) en lugar de reducir, asociando cada `else` con el `if` abierto más próximo.

### 5.2 Construcción del árbol

Cada acción semántica de una producción construye un nodo del paquete `compscript.ast`, por ejemplo:

```java
expr ::= expr:a MAS expr:b  {: RESULT = new A.Binaria("+", a, b, aleft, aright); :}
       | ID:id               {: RESULT = new A.AccesoVariable(id, idleft, idright); :}
       ;
```

El resultado final del análisis sintáctico, `parser.raiz` (de tipo `ArrayList`, tipo *raw* por restricción de CUP sobre los no terminales), contiene la lista de instrucciones de nivel superior del programa.

### 5.3 Recuperación de errores sintácticos

El método `syntax_error(Symbol s)`, sobrescrito en la sección `parser code {: ... :}`, registra el error en el `Contexto` compartido con la posición del símbolo no esperado. La producción `instruccion ::= error PUNTO_COMA` implementa el modo pánico: ante un error, el analizador descarta símbolos hasta encontrar un punto y coma, y continúa el análisis de la instrucción siguiente, permitiendo acumular múltiples errores sintácticos en una sola pasada.

## 6. Árbol de sintaxis abstracta (paquete `compscript.ast`, clase `A`)

### 6.1 Interfaces base

```java
public interface Nodo {
    String etiquetaAst();
    default List<Nodo> hijosAst() { return new ArrayList<>(); }
}
public interface Instruccion extends Nodo { void ejecutar(Entorno e); }
public interface Expresion  extends Nodo { Valor evaluar(Entorno e); }
```

Todas las clases de nodo (aproximadamente 25, entre ellas `A.Declaracion`, `A.If`, `A.While`, `A.For`, `A.DoWhile`, `A.Match`, `A.Binaria`, `A.Unaria`, `A.Cast`, `A.Llamada`, `A.OperacionLista`, `A.DeclaracionFuncion`, `A.DeclaracionStruct`) implementan `Instruccion` o `Expresion` según corresponda, siguiendo el patrón de intérprete de árbol de sintaxis (*tree-walking interpreter*): cada nodo conoce el procedimiento para ejecutarse o evaluarse a sí mismo, en lugar de delegar en una clase visitante externa. Los métodos `etiquetaAst()` y `hijosAst()` permiten recorrer el árbol de manera uniforme para la generación del reporte de AST, sin necesidad de una jerarquía de visitantes adicional.

### 6.2 Ejecución en tres pasadas (`Interprete.ejecutar(String codigo)`)

```java
public static Contexto ejecutar(String codigo) {
    Contexto ctx = new Contexto();
    Lexer lexer = new Lexer(new StringReader(codigo));
    lexer.contexto = ctx;
    Parser parser = new Parser(lexer);
    parser.contexto = ctx;
    parser.parse();
    ...
    // 1a pasada: registro de funciones, metodos y structs
    // 2a pasada: declaraciones y asignaciones globales
    // 3a pasada: ejecucion de RUN_MAIN
}
```

Siguiendo la recomendación del enunciado del proyecto, la ejecución se organiza en tres recorridos del árbol: el primero registra las funciones, métodos y structs declarados, para permitir referencias hacia adelante (una función puede invocar a otra declarada más abajo en el archivo); el segundo ejecuta las declaraciones y asignaciones de ámbito global; el tercero ejecuta la sentencia `RUN_MAIN`, que constituye el punto de entrada del programa.

### 6.3 Validar la categoría de tipo antes de castear: decisión de diseño no obvia (bug real corregido)

`A.Declaracion.validarVector(Valor crudo, Tipo declarado, Entorno e)` recibe el valor ya evaluado de la expresión inicializadora de una variable declarada como vector (`let v: int[] = <expresion>;`). Esa expresión inicializadora es una `<expresion>` general de la gramática (5.11), no exclusivamente un literal de vector (`[ ... ]`): nada impide sintácticamente escribir `let v: int[] = 5;`, una expresión aritméticamente válida pero incompatible con el tipo declarado.

Antes de la corrección (auditoría de código, 2026-07-21), `validarVector` asumía sin comprobación que `crudo` ya era un vector e invocaba directamente `crudo.lista()` (que hace un cast interno a `List<Valor>`). Ante una inicialización no vectorial, ese cast lanzaba un `ClassCastException` de Java sin controlar en ningún punto del pipeline: el proceso completo de la máquina virtual terminaba abruptamente, no solo el programa `.cs` que se estaba analizando. Esto contradecía directamente el requisito de recuperación de errores del enunciado (4.3): un error semántico debe quedar registrado en la tabla de errores y terminar la ejecución de forma ordenada, nunca interrumpir el intérprete mismo.

La corrección agrega una guarda de tipo explícita antes de castear:

```java
private Valor validarVector(Valor crudo, Tipo declarado, Entorno e) {
    if (crudo == null || crudo.tipo.cat != Tipo.Cat.VECTOR)
        e.errorSemantico("el vector '" + id + "' espera un literal de vector y recibio "
                + Entorno.tipoDe(crudo), linea, columna);
    List<Valor> vals = crudo.lista();
    ...
}
```

Con la guarda, `let v: int[] = 5;` produce un error semántico ordenado (`el vector 'v' espera un literal de vector y recibio Entero`) en lugar de derribar el proceso. La lección de diseño generalizable, propia de un intérprete sin gramática que garantice el tipo desde afuera del nodo: **ningún cast hacia una representación interna del valor (`.lista()`, `.campos()`) debe ejecutarse sin comprobar antes la categoría (`Tipo.Cat`) del valor recibido** — el mismo principio, aplicado consistentemente, es el que ya siguen `AccesoVector`, `AccesoCampo` y `OperacionLista` al validar `s.tipo.cat` antes de llamar a `.lista()` o `.campos()`.

## 7. Análisis semántico y entorno de ejecución (paquete `compscript.interprete`)

### 7.1 `Contexto`: estado global de una ejecución

Contiene la consola de salida (`StringBuilder consola`), la lista de errores acumulados (`List<RegistroError> errores`), el registro de tokens reconocidos (`List<Object[]> tokens`), el registro plano de todos los símbolos declarados durante la ejecución (`List<Simbolo> simbolos`) y los mapas de funciones y structs registrados en la primera pasada (`LinkedHashMap<String, A.DeclaracionFuncion> funciones`, `LinkedHashMap<String, A.DeclaracionStruct> structs`). Se instancia una única vez por llamada a `Interprete.ejecutar(...)`, garantizando que cada ejecución parta de un estado limpio.

### 7.2 `Entorno`: pila de ámbitos enlazados

```java
public class Entorno {
    public final Contexto contexto;
    public final Entorno padre;       // null en el entorno global
    private final LinkedHashMap<String, Simbolo> tabla = new LinkedHashMap<>();
    public Entorno crearHijo(String nombre) { return new Entorno(contexto, this, nombre); }
    public Simbolo buscar(String id) { /* recorre this, this.padre, this.padre.padre, ... */ }
}
```

Cada bloque de instrucciones que introduce un nuevo ámbito (cuerpo de un `if`, de un ciclo, de una función) crea un `Entorno` hijo mediante `crearHijo(String nombre)`. La búsqueda de un identificador (`buscar`, `obtener`) recorre la cadena de entornos padre hasta encontrarlo o llegar al entorno global, implementando alcance estático anidado. Es relevante señalar que el entorno de una función se crea siempre como hijo del entorno global (`ctx.global.crearHijo(f.id)`), y no del entorno del invocador, de modo que una función no tiene visibilidad sobre las variables locales de quien la llama.

### 7.3 `Tipo`, `Valor` y `Simbolo`

`Tipo` modela la categoría de dato (`enum Cat { INT, DOUBLE, BOOL, CHAR, STRING, VOID, NULL, VECTOR, LIST, STRUCT }`) junto con la información adicional necesaria para los tipos compuestos (`elemento` y `dimensiones` para vectores y listas, `structName` para structs), e implementa igualdad estructural mediante `equals(Object)`, utilizada en la comprobación de compatibilidad de tipos.

`Valor` empareja un `Tipo` con el objeto Java que lo representa en tiempo de ejecución (`Integer`, `Double`, `Boolean`, `Character`, `String`, `List<Valor>` para vectores y listas, `LinkedHashMap<String,Valor>` para structs), y expone métodos de formateo: `texto()` (formato de consola, usado por `console.log`) y `reporte()` (formato de la tabla de símbolos, con comillas para cadenas y caracteres).

`Simbolo` representa una entrada de la tabla de símbolos: nombre, categoría (`Variable`, `Vector`, `Lista` o `Struct`, calculada por `Simbolo.categoriaDe(Tipo)`), tipo, mutabilidad (`let` frente a `const`), ámbito, valor (mutable, actualizado en cada asignación) y posición de declaración.

### 7.4 `Operaciones`: comprobación de tipos y semántica de operadores

Implementa las tablas de compatibilidad especificadas en el enunciado del proyecto para los operadores aritméticos (suma, resta, multiplicación, división, potencia, raíz, módulo), relacionales y lógicos, así como los cinco casteos permitidos entre tipos primitivos (`int↔double`, `int↔char`, `char→double`). Toda combinación de tipos no contemplada en las tablas produce un error semántico. Firma representativa:

```java
public static Valor aritmetica(String op, Valor a, Valor b, Entorno e, int l, int c)
public static Valor relacional(String op, Valor a, Valor b, Entorno e, int l, int c)
public static Valor cast(Valor v, Tipo destino, Entorno e, int l, int c)
```

`A.Binaria.evaluar()` (paquete `compscript.ast`) es quien decide, para cada operador, qué operandos evaluar antes de invocar a `Operaciones`: la aritmética y las relacionales evalúan siempre ambos lados (no dependen del valor del otro operando para decidir si evaluarlo), pero `&&` y `||` sí implementan cortocircuito — si `izq` ya determina el resultado (`false` para `&&`, `true` para `||`), `der` no se evalúa. Esto no es una optimización cosmética: es semánticamente observable en expresiones de guarda como `x != 0 && 10/x > 1`, donde evaluar siempre ambos lados dispararía una división entre cero en vez de evitarla.


### 7.5 Manejo de errores semánticos: decisión de diseño no obvia

```java
public void errorSemantico(String descripcion, int linea, int columna) {
    contexto.error("Semantico", descripcion, linea, columna);
    throw new ErrorSemantico(descripcion, linea, columna);
}
```

A diferencia de un esquema de propagación de errores por valores nulos (en el cual una expresión inválida se sustituye por un valor vacío y la ejecución continúa), CompScript adopta el comportamiento exigido por el enunciado del proyecto: un error semántico detectado en tiempo de ejecución se registra en la tabla de errores y, mediante el lanzamiento de la excepción interna `ErrorSemantico` (sin captura de traza, por tratarse de control de flujo y no de una condición excepcional real), interrumpe de forma ordenada el resto de la ejecución. La excepción es capturada en el nivel más externo, dentro de `Interprete.ejecutar(...)`.

Las sentencias `break`, `continue` y `return` se implementan mediante el mismo mecanismo de excepciones no capturadas por el usuario (clase `Senales`, con las subclases `Break`, `Continue` y `Retorno`), interpretadas por los ciclos y por el mecanismo de invocación de funciones respectivamente. Si alguna de estas señales llega sin ser capturada hasta `Interprete.ejecutar(...)` (por ejemplo, un `break` fuera de cualquier ciclo), se traduce en un error semántico adicional.

### 7.6 Invocación de funciones y métodos

La resolución de llamadas (`A.invocar(String id, List<Argumento> args, Entorno caller, int l, int c)`, método privado de la clase `A`) asocia los argumentos por identificador, no por posición, conforme al enunciado del proyecto. Valida que no existan argumentos repetidos ni desconocidos, aplica los valores por defecto declarados en los parámetros cuando el argumento correspondiente se omite, y comprueba la compatibilidad de tipos tanto de los argumentos como del valor de retorno frente al tipo declarado de la función.

## 8. Reportes (paquete `compscript.reportes`, clase `Reportes`)

```java
public static File[] generar(Contexto ctx, File carpeta) throws Exception
```

Genera cinco archivos a partir del `Contexto` de la última ejecución: `tokens.html`, `errores.html`, `simbolos.html`, `ast.html` (árbol de sintaxis abstracta representado como lista anidada HTML) y `ast.dot` (fuente en formato Graphviz, para su renderizado externo mediante `dot -Tpng ast.dot -o ast.png`). Los nombres de los tipos de token se obtienen por reflexión sobre los campos públicos de tipo `int` de la clase generada `compscript.analisis.sym`, evitando mantener una tabla de nombres duplicada manualmente. Todo el contenido HTML es autocontenido (estilos incluidos en línea), de forma que los reportes pueden abrirse sin conexión a internet.

## 9. Interfaz gráfica (paquete `compscript.gui`)

`EditorApp` (extiende `javafx.application.Application`) construye la interfaz completa por código, sin archivos FXML: una barra de botones (Nuevo, Abrir, Guardar, Ejecutar, Reportes, Ver AST), un `TabPane` para múltiples archivos abiertos simultáneamente y un `TextArea` de consola de solo lectura, organizados en un `SplitPane` vertical dentro de un `BorderPane`.

`Lanzador` contiene el método `main` efectivamente utilizado para iniciar la aplicación desde un entorno de desarrollo: al no extender `javafx.application.Application`, evita que la máquina virtual de Java aplique la verificación de módulos de JavaFX que provoca el error "JavaFX runtime components are missing" cuando la clase de entrada extiende `Application` directamente y JavaFX no se encuentra en el *module path*.

Cada ejecución desde la interfaz crea un `Contexto` nuevo (`Interprete.ejecutar(area.getText())`), de modo que los reportes disponibles corresponden siempre al contenido de la pestaña activa en el momento de presionar "Ejecutar", y no a ejecuciones anteriores.

## 10. Resumen de decisiones de diseño no evidentes

1. **AST en lugar de evaluación S-atribuida**: necesario por la presencia de ciclos, funciones y recursión (sección 3.1).
2. **El error semántico aborta la ejecución**, en lugar de propagarse como valor nulo, conforme lo exige el enunciado del proyecto (sección 7.5).
3. **El entorno de una función cuelga del entorno global**, no del entorno del invocador, para preservar el alcance estático (sección 7.2).
4. **`UMENOS` es un pseudo-terminal** utilizado exclusivamente para resolver la precedencia distinta entre la resta binaria y la negación unaria, ambas construidas sobre el mismo símbolo léxico (sección 5.1).
5. **Los argumentos de una llamada se asocian por identificador**, no por posición, lo que permite omitir parámetros con valor por defecto en cualquier posición de la lista (sección 7.6).
6. **La tabla de símbolos del reporte es un registro histórico acumulativo**, no un estado final: conserva una entrada por cada declaración ocurrida durante la ejecución, incluyendo las de cada invocación de una función recursiva.
7. **Doble representación del AST** (HTML autocontenido y archivo `.dot` de Graphviz), además de la vista embebida en la interfaz gráfica, para satisfacer el requisito de mostrar el reporte "desde la interfaz" sin renunciar a una representación de mayor fidelidad para su documentación externa.
8. **`&&` y `||` implementan cortocircuito real** en `A.Binaria.evaluar()`: si el operando izquierdo ya determina el resultado, el derecho no se evalúa. No es una optimización cosmética: es observable semánticamente en guardas típicas (`x != 0 && 10/x > 1`), donde evaluar siempre ambos lados dispararía el error que la guarda pretende evitar (sección 7.4; corregido en la auditoría de código del 2026-07-21).
9. **Cada nodo del AST valida la categoría de tipo (`Tipo.Cat`) antes de castear a su representación interna** (`.lista()`, `.campos()`), en lugar de asumir el tipo del valor recibido: la corrección de `A.Declaracion.validarVector` (sección 6.3) es el ejemplo real de qué ocurre cuando esa guarda falta — un `ClassCastException` sin control que interrumpe el proceso completo, no solo el programa analizado (corregido en la misma auditoría del 2026-07-21).
