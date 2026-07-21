# Guía de Programación — ConjAnalyzer

Universidad de San Carlos de Guatemala
Facultad de Ingeniería — Escuela de Ciencias y Sistemas
Organización de Lenguajes y Compiladores 1
Proyecto 1 — ConjAnalyzer

## Cómo usar esta guía

Esta guía es material de **tutoría**, no de referencia. Acá el rol es al revés del manual técnico: en el `ManualTecnico.md` alguien más ya construyó el proyecto y lo describe; en esta guía **vos vas a programar ConjAnalyzer desde cero**, tema por tema, y en cada paso se explica el *porqué* antes que el *cómo*. Cuando la guía dice "ahora escribí...", el código que sigue es exactamente el que existe hoy en `ConjAnalyzer/src/`, verificado contra el enunciado real (`OLC1_PT1 _S2024`, sección y número citados entre paréntesis) y contra el Libro del Dragón (Aho/Lam/Sethi/Ullman, 2ª ed.) donde corresponda.

No copies y pegues sin leer la explicación previa: la evaluación oral de OLC1 pregunta por qué tomaste cada decisión, no solo si el programa corre.

## Índice

1. Preparación del proyecto (Maven, JFlex, CUP)
2. Tema: análisis léxico — la tabla de tokens del universo de conjuntos
3. Tema: análisis sintáctico — notación prefija anidada y el árbol acotado
4. Tema: el modelo `Conjunto` y las operaciones
5. Tema: el `Simplificador` — las 5 leyes algebraicas
6. Tema: diagrama de Venn con Canvas
7. Tema: reportes JSON con Gson
8. Errores comunes reales al programar cada tema

---

## 1. Preparación del proyecto (Maven, JFlex, CUP)

### 1.1 Por qué este stack

ConjAnalyzer es un intérprete: lee texto fuente (`.ca`), lo reconoce con un analizador léxico y uno sintáctico, y **ejecuta** el resultado directamente (no genera código de otro lenguaje — eso es tema de los capítulos 8 a 12 del Dragón, fuera de alcance del curso). Para construir el léxico y la sintaxis se usan dos generadores clásicos:

- **JFlex**: a partir de una especificación declarativa de expresiones regulares (`Lexer.flex`), genera la clase Java `Lexer` que reconoce tokens. Es la implementación práctica del algoritmo de construcción de autómatas finitos a partir de expresiones regulares que ves en el Cap. 3 del Dragón (Thompson + subconjuntos, aunque JFlex lo hace por vos).
- **CUP** (Construction of Useful Parsers): a partir de una gramática BNF anotada con acciones semánticas (`parser.cup`), genera un parser **LALR** (Cap. 4, análisis ascendente) que además ejecuta código Java en cada reducción — esto es traducción dirigida por la sintaxis (Cap. 5).

Igual que en DataForge (el proyecto hermano de este mismo curso), ninguno de los dos genera un AST del programa completo por defecto: eso lo decidís vos en las acciones de la gramática. La novedad de ConjAnalyzer, que vas a entender en el Tema 3, es que acá **sí** hace falta construir un árbol, aunque acotado.

### 1.2 Arrancá el `pom.xml`

Creá un proyecto Maven vacío en IntelliJ IDEA (JDK 17, gestionado por el IDE) y armá el `pom.xml` con las dependencias y plugins que vas a necesitar. No hace falta que lo escribas de memoria — la razón de cada bloque es lo importante:

```xml
<properties>
    <maven.compiler.source>17</maven.compiler.source>
    <maven.compiler.target>17</maven.compiler.target>
    <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
</properties>

<dependencies>
    <!-- Runtime de CUP: aporta la clase Symbol que devuelve el lexer -->
    <dependency>
        <groupId>com.github.vbmacher</groupId>
        <artifactId>java-cup-runtime</artifactId>
        <version>11b-20160615</version>
    </dependency>
    <!-- JavaFX: controles de UI + Canvas para el diagrama de Venn -->
    <dependency>
        <groupId>org.openjfx</groupId>
        <artifactId>javafx-controls</artifactId>
        <version>21.0.4</version>
    </dependency>
    <!-- Gson: serializa el JSON de simplificacion (5.4) -->
    <dependency>
        <groupId>com.google.code.gson</groupId>
        <artifactId>gson</artifactId>
        <version>2.11.0</version>
    </dependency>
</dependencies>
```

`java-cup-runtime` es distinto del plugin que **genera** el parser: el runtime es la biblioteca mínima (la clase `Symbol`, principalmente) que el código generado necesita para compilar. `gson` es nuevo respecto a DataForge — lo vas a usar recién en el Tema 7, pero declaralo ahora para no interrumpir el flujo de Maven más adelante.

En `<build><plugins>` agregá, en este orden de razonamiento (no de aparición en el XML):

```xml
<!-- CUP: genera Parser.java y sym.java en target/generated-sources/cup -->
<plugin>
    <groupId>com.github.vbmacher</groupId>
    <artifactId>cup-maven-plugin</artifactId>
    <version>11b-20160615-3</version>
    <executions>
        <execution><goals><goal>generate</goal></goals></execution>
    </executions>
    <configuration>
        <className>Parser</className>
        <symbolsName>sym</symbolsName>
    </configuration>
</plugin>

<!-- JFlex: genera Lexer.java en target/generated-sources/jflex -->
<plugin>
    <groupId>de.jflex</groupId>
    <artifactId>jflex-maven-plugin</artifactId>
    <version>1.9.1</version>
    <executions>
        <execution><goals><goal>generate</goal></goals></execution>
    </executions>
</plugin>

<!-- Corre la GUI con: mvn clean javafx:run -->
<plugin>
    <groupId>org.openjfx</groupId>
    <artifactId>javafx-maven-plugin</artifactId>
    <version>0.0.8</version>
    <configuration>
        <mainClass>conjanalyzer.gui.EditorApp</mainClass>
    </configuration>
</plugin>

<!-- Prueba por consola con: mvn compile exec:java -->
<plugin>
    <groupId>org.codehaus.mojo</groupId>
    <artifactId>exec-maven-plugin</artifactId>
    <version>3.1.0</version>
    <configuration>
        <mainClass>conjanalyzer.TestInterprete</mainClass>
    </configuration>
</plugin>
```

Los dos primeros plugins corren en la fase `generate-sources`, **antes** de que Maven compile tus clases propias — así, cuando tu código (por ejemplo `Interprete.java`) hace `import conjanalyzer.analisis.Lexer;`, esa clase ya existe. Regla que no vas a romper en todo el proyecto: **nunca edites a mano** lo que aparece en `target/generated-sources/` — si necesitás cambiar un token o una producción, el cambio va en `Lexer.flex` o `parser.cup`, y Maven regenera el resto.

### 1.3 Estructura de carpetas

```
ConjAnalyzer/
├── pom.xml
├── docs/                       (gramática BNF entregable, manuales)
├── entradas/                   (tus casos de prueba .ca)
├── reportes/                   (salida generada — no la creás a mano, mkdirs() lo hace)
└── src/main/
    ├── jflex/Lexer.flex
    ├── cup/parser.cup
    └── java/conjanalyzer/
        ├── TestInterprete.java
        ├── analisis/           (vacío al inicio — JFlex/CUP lo llenan en target/, no acá)
        ├── interprete/
        ├── reportes/
        └── gui/
```

Con esto armado y `mvn compile` corriendo sin errores (aunque `Lexer.flex`/`parser.cup` estén casi vacíos todavía), estás listo para el Tema 2.

---

## 2. Tema: análisis léxico — la tabla de tokens del universo de conjuntos

### 2.1 Por qué el vocabulario de ConjAnalyzer es distinto al de DataForge

Antes de escribir una sola regla, mirá qué vocabulario necesita reconocer el lenguaje leyendo la sección 4 del enunciado: 3 palabras reservadas (`CONJ`, `OPERA`, `EVALUAR`), 4 operadores de conjuntos (`U`, `&`, `^`, `-`), símbolos estructurales (`->`, `~`, `{`, `}`, `(`, `)`, `:`, `;`, `,`) y dos tipos de comentario. La diferencia de diseño más importante frente a DataForge es la sección **4.1 — Case Sensitive**: acá `conjuntoA` y `conjuntoa` son identificadores **distintos**. Esa única decisión se propaga a tres lugares del proyecto (lo vas a ver reaparecer en el Tema 3 y en el Tema 4), y arranca en el `.flex`.

### 2.2 Escribí el encabezado del `.flex`

```jflex
package conjanalyzer.analisis;

import java_cup.runtime.Symbol;

%%

%class Lexer
%public
%unicode
%cup
%line
%column
```

Tres detalles que no son opcionales:

- **`%public`**: sin esta directiva la clase `Lexer` queda `package-private` y no la vas a poder instanciar desde `conjanalyzer.interprete.Interprete`. Es el mismo error que en DataForge — anotalo la primera vez que compiles y te dé "Lexer is not public".
- **NO uses `%ignorecase`**: es la contraparte directa de la 4.1. Si lo agregás "por costumbre" (viniendo de DataForge, que sí era case-insensitive), rompés un requisito explícito del enunciado.
- **`%line` / `%column`**: sin esto, `yyline`/`yycolumn` no existen y no podés reportar la posición de un token ni de un error.

### 2.3 Conectá el lexer con el `Entorno` (todavía no existe — lo construís en el Tema 4)

El lexer necesita un canal para avisar "reconocí este token" y "encontré un carácter inválido" sin usar `System.out`/`System.err` directamente (los reportes del Tema 7 necesitan esa información estructurada). La técnica es un campo público que el `Interprete` asigna después de crear el lexer:

```jflex
%{
  public conjanalyzer.interprete.Entorno entorno;

  private Symbol symbol(int type) {
    registrar(type);
    return new Symbol(type, yyline, yycolumn, yytext());
  }
  private Symbol symbol(int type, Object value) {
    registrar(type);
    return new Symbol(type, yyline, yycolumn, value);
  }
  private void registrar(int type) {
    if (entorno != null) entorno.registrarToken(yytext(), type, yyline, yycolumn);
  }
%}
```

El chequeo `if (entorno != null)` no es defensivo por las dudas: importa porque vas a querer poder instanciar un `Lexer` suelto (por ejemplo en una prueba unitaria) sin que explote por falta de `Entorno`.

### 2.4 Definí las macros y las reglas, EN ESTE ORDEN

```jflex
Letra        = [a-zA-Z]
Digito       = [0-9]
Id           = {Letra}({Letra}|{Digito})*
Numero       = {Digito}+
ComentLinea  = "#"[^\r\n]*
ComentMulti  = "<!"~"!>"
Blancos      = [ \t\r\n]+

%%

{ComentLinea}   { /* ignorar */ }
{ComentMulti}   { /* ignorar */ }
{Blancos}       { /* ignorar */ }

"CONJ"                    { return symbol(sym.CONJ); }
"OPERA"                   { return symbol(sym.OPERA); }
"EVALUAR"                 { return symbol(sym.EVALUAR); }

"U"                       { return symbol(sym.UNION); }
"&"                       { return symbol(sym.INTERSECCION); }
"^"                       { return symbol(sym.COMPLEMENTO); }
"-"                       { return symbol(sym.DIFERENCIA); }

"->"                      { return symbol(sym.FLECHA); }
"~"                       { return symbol(sym.VIRGULILLA); }
"{"                       { return symbol(sym.LLAVE_IZQ); }
"}"                       { return symbol(sym.LLAVE_DER); }
"("                       { return symbol(sym.PAR_IZQ); }
")"                       { return symbol(sym.PAR_DER); }
":"                       { return symbol(sym.DOS_PUNTOS); }
";"                       { return symbol(sym.PUNTO_COMA); }
","                       { return symbol(sym.COMA); }

{Id}                      { return symbol(sym.ID); }
{Numero}                  { return symbol(sym.NUMERO); }

[\x21-\x7E]               { return symbol(sym.SIMBOLO); }

[^]                       { if (entorno != null) {
                              entorno.error("Lexico", "el caracter '" + yytext()
                                  + "' no pertenece al lenguaje", yyline, yycolumn);
                            } }
```

Pensá el orden como una serie de **prioridades**, de más específico a más genérico — es la aplicación directa de la regla del Dragón "ante un empate de longitud de match, gana la regla declarada primero":

1. **Comentarios y blancos primero**: se reconocen y se descartan sin generar token.
2. **Reservadas y operadores**: `"CONJ"`, `"U"`, `"&"`, etc. son literales exactos. Si los pusieras después de `{Id}`, JFlex nunca los alcanzaría — `{Id}` ya habría matcheado `"CONJ"` como un identificador cualquiera.
3. **`"->"` antes que `"-"`**: acá es donde el "empate de longitud" se vuelve real. Ante la entrada `->`, tanto la regla de `"-"` (match de 1 carácter) como la de `"->"` (match de 2) son candidatas; JFlex se queda con el match **más largo**, así que en la práctica el orden entre estas dos reglas no importa para este caso puntual — pero sí importa el principio general: escribí siempre el patrón más largo lo antes posible dentro de su "familia", para no depender de la sutileza de cuándo un empate de longitud se resuelve por orden y cuándo por longitud.
4. **`{Id}` y `{Numero}` después de las reservadas**: ahora sí, cualquier identificador que no sea una palabra reservada cae acá.
5. **`SIMBOLO` como comodín**: cualquier carácter ASCII imprimible (33 a 126) que no calzó en ninguna regla anterior — esto es lo que permite `CONJ : signos -> !,?,@,$,%;` en `entradas/ejemplo4_nosimplificable.ca`. Sin este comodín, un conjunto de puntuación sería imposible de declarar.
6. **`[^]` al final**: literalmente "cualquier cosa" — como JFlex ya probó todas las reglas anteriores y ninguna calzó, lo que llega acá es CUALQUIER carácter fuera del universo ASCII 33-126 (por ejemplo una `ñ`, que en UTF-8/Unicode existe pero no es un símbolo imprimible del rango declarado). Fijate que esta regla **no aborta**: registra el error y sigue — es tu primera aplicación del principio "ningún error detiene el análisis" que vas a repetir en el Tema 3.

### 2.5 Verificá con un archivo de prueba

Antes de tocar una sola línea de `parser.cup`, escribí un pequeño lexer de prueba o corré `mvn compile` y revisá que `target/generated-sources/jflex/conjanalyzer/analisis/Lexer.java` se generó sin errores. Un error típico en este punto: si el `.cup` todavía no existe, `sym` tampoco existe y `Lexer.flex` no compila — es normal, y es la señal de que te toca el Tema 3.

---

## 3. Tema: análisis sintáctico — notación prefija anidada y el árbol acotado

### 3.1 La gramática, en prosa antes que en CUP

La sección 4.6 del enunciado define las operaciones en **notación prefija (polaca)**: el operador va antes que sus operandos. La razón práctica (que vas a redescubrir vos mismo al escribir la gramática) es que la notación prefija **nunca necesita paréntesis para desambiguar precedencia** — a diferencia de la notación infija de DataForge, que sí necesitaba una jerarquía de no terminales (`expresion` → `termino` → `factor`) para que `2 + 3 * 4` se agrupe bien. Acá, como cada operador "sabe" cuántos operandos consume, la estructura es no ambigua por construcción:

```
U U {A} {B} {C}      =  (A U B) U C
& U {C} {A} {B}      =  (C U A) & B
```

Escribí primero la gramática en BNF (esto es exactamente lo que vas a dejar documentado en `docs/gramatica.txt`, el entregable de la sección 9 — nunca una copia mecánica del `.cup`):

```
<inicio>        ::= "{" <sentencias> "}"
<sentencias>    ::= <sentencias> <sentencia> | <sentencia>
<sentencia>     ::= <definicion-conjunto> | <definicion-operacion> | <evaluacion>

<definicion-conjunto> ::= "CONJ" ":" ID "->" <notacion> ";"
<notacion>      ::= <elemento> "~" <elemento>
                  | <lista-elementos>

<definicion-operacion> ::= "OPERA" ":" ID "->" <operacion> ";"
<operacion>     ::= "U" <operacion> <operacion>
                  | "&" <operacion> <operacion>
                  | "-" <operacion> <operacion>
                  | "^" <operacion>
                  | "{" ID "}"

<evaluacion>    ::= "EVALUAR" "(" "{" <lista-elementos> "}" "," ID ")" ";"
```

### 3.2 Por qué acá SÍ hace falta un árbol de expresión (aunque acotado)

Esta es la decisión arquitectónica más importante de todo el proyecto, así que pensala con calma antes de escribir el `.cup`.

En DataForge, cada instrucción se ejecutaba **una sola vez**, directamente en la acción semántica de su producción: no había control de flujo, así que no existía ninguna razón para guardar una representación intermedia de la instrucción — evaluarla y tirarla era suficiente. Es traducción dirigida por la sintaxis "pura" (Cap. 5 del Dragón): la gramática es **S-atribuida**, cada no terminal sintetiza su valor hacia arriba y no queda nada más que el resultado final.

En ConjAnalyzer, una `OPERA` (una vez definida) se necesita **tres veces, de tres formas distintas**:

1. **Evaluarla** contra los conjuntos base y el universo, para obtener el `Set<Character>` resultante (Tema 4).
2. **Simplificarla** reescribiendo su estructura con las leyes de la sección 7 — esto es una transformación de árbol, no se puede hacer sobre un `Set<Character>` ya calculado, porque el resultado simplificado necesita expresarse en la MISMA notación prefija que el original, solo que más corta (Tema 5).
3. **Servir de función booleana** al diagrama de Venn, para decidir píxel por píxel si un punto pertenece al resultado (Tema 6).

Si `operacion` sintetizara directamente un `Set<Character>` (como hacía DataForge con sus expresiones aritméticas), la segunda y la tercera reutilización serían imposibles — ya no tendrías la ESTRUCTURA de la operación, solo su resultado. La regla general que te vas a llevar de esta etapa (y que va a reaparecer completa en CompScript con control de flujo real): **si algo se recorre una sola vez, ejecutalo directo en la acción semántica; si se recorre más de una vez, de formas distintas, hace falta materializarlo como estructura de datos — un árbol**.

Fijate que el árbol de ConjAnalyzer es **acotado**: solo las operaciones (`OPERA`) se representan como árbol. Las sentencias `CONJ` y `EVALUAR` se siguen ejecutando directo en la acción del `.cup`, exactamente igual que en DataForge — no hay AST del programa completo.

### 3.3 Escribí `NodoOperacion` antes que el `.cup`

Como el no terminal `operacion` va a sintetizar un `NodoOperacion`, escribí esa clase primero. Un nodo tiene solo 3 formas posibles:

```java
public class NodoOperacion {
    public final char op;             // 'U' '&' '-' '^'  o  '\0' si es hoja
    public final String nombreConj;   // solo hojas
    public final NodoOperacion izq;
    public final NodoOperacion der;    // null en unario y hoja

    public static NodoOperacion hoja(String nombre) { return new NodoOperacion('\0', nombre, null, null); }
    public static NodoOperacion unario(char op, NodoOperacion hijo) { return new NodoOperacion(op, null, hijo, null); }
    public static NodoOperacion binario(char op, NodoOperacion a, NodoOperacion b) { return new NodoOperacion(op, null, a, b); }

    public boolean esHoja()   { return op == '\0'; }
    public boolean esUnario() { return op == '^'; }
}
```

Con esas 3 fábricas estáticas ya podés escribir las acciones del `.cup` sin que el constructor quede expuesto (`private`).

### 3.4 Escribí `parser.cup`

Empezá por el bloque `parser code` y las declaraciones de terminales/no terminales — TODAS las declaraciones `non terminal` van antes de cualquier producción (restricción real de CUP, no estilística):

```cup
package conjanalyzer.analisis;

import java_cup.runtime.*;
import java.util.ArrayList;
import conjanalyzer.interprete.Entorno;
import conjanalyzer.interprete.NodoOperacion;

parser code {:
    public Entorno entorno = new Entorno();

    public void syntax_error(Symbol s) {
        entorno.error("Sintactico",
            "no se esperaba '" + s.value + "'", s.left, s.right);
    }
:};

terminal CONJ, OPERA, EVALUAR;
terminal UNION, INTERSECCION, COMPLEMENTO, DIFERENCIA;
terminal FLECHA, VIRGULILLA;
terminal LLAVE_IZQ, LLAVE_DER, PAR_IZQ, PAR_DER;
terminal DOS_PUNTOS, PUNTO_COMA, COMA;
terminal String ID, NUMERO, SIMBOLO;

non terminal inicio, lista_sent, sentencia;
non terminal def_conj, def_opera, evaluar_stmt;
non terminal NodoOperacion operacion;
non terminal String elemento;
non terminal ArrayList lista_elem;

start with inicio;
```

Notá el tipo raw `ArrayList` (sin `<String>`) — es una restricción práctica de CUP con genéricos, la misma que en DataForge.

Ahora la producción central, la de `operacion` — es la que sintetiza el árbol:

```cup
operacion   ::= UNION operacion:a operacion:b
                {: RESULT = NodoOperacion.binario('U', a, b); :}
              | INTERSECCION operacion:a operacion:b
                {: RESULT = NodoOperacion.binario('&', a, b); :}
              | DIFERENCIA operacion:a operacion:b
                {: RESULT = NodoOperacion.binario('-', a, b); :}
              | COMPLEMENTO operacion:a
                {: RESULT = NodoOperacion.unario('^', a); :}
              | LLAVE_IZQ ID:id LLAVE_DER
                {: RESULT = NodoOperacion.hoja(id); :} ;
```

Compará esto con la producción de `def_conj` y `evaluar_stmt`, que NO sintetizan nada — solo llaman directo al `Entorno` (vas a escribir `Entorno` recién en el Tema 4, así que por ahora dejá esas acciones como comentario o con un método placeholder):

```cup
def_conj    ::= CONJ DOS_PUNTOS ID:id FLECHA elemento:a VIRGULILLA elemento:b PUNTO_COMA
                {: parser.entorno.definirConjuntoRango(id, a, b, idleft, idright); :}
              | CONJ DOS_PUNTOS ID:id FLECHA lista_elem:l PUNTO_COMA
                {: parser.entorno.definirConjuntoLista(id, l, idleft, idright); :} ;

def_opera   ::= OPERA DOS_PUNTOS ID:id FLECHA operacion:op PUNTO_COMA
                {: parser.entorno.definirOperacion(id, op, idleft, idright); :} ;

evaluar_stmt ::= EVALUAR PAR_IZQ LLAVE_IZQ lista_elem:datos LLAVE_DER COMA ID:op PAR_DER PUNTO_COMA
                {: parser.entorno.evaluar(datos, op, opleft, opright); :} ;
```

Esta diferencia (`def_conj`/`evaluar_stmt` sin tipo, `operacion` con tipo `NodoOperacion`) es la prueba visual de la decisión de diseño del punto 3.2: solo lo que se recorre más de una vez necesita quedar representado como estructura.

### 3.5 Recuperación de errores: modo pánico (Dragón §4.8.3)

Agregá la alternativa de error a `sentencia`, no a `inicio` ni a `operacion`:

```cup
sentencia   ::= def_conj | def_opera | evaluar_stmt
              | error PUNTO_COMA ;
```

Cuando el parser encuentra un error sintáctico, CUP descarta símbolos de la pila hasta encontrar un estado donde pueda desplazar el terminal especial `error` — en este caso, eso ocurre al nivel de `sentencia` — y después descarta tokens de la entrada hasta el próximo `;`. El efecto práctico: una sentencia mal formada no aborta todo el análisis, solo se pierde ESA sentencia, y las siguientes se procesan con normalidad. Si ponés `error` en el lugar equivocado (por ejemplo, solo dentro de `operacion`), un error sintáctico en una sentencia va a "burbujear" mucho más arriba en la gramática antes de encontrar dónde sincronizar, y vas a perder de más.

### 3.6 Verificá con errores a propósito

Escribí un `.ca` con un error sintáctico deliberado (por ejemplo, `CONJ : -> a~z;`, sin el `ID`) y confirmá dos cosas: que el error se reporta con línea/columna razonables, y que una sentencia válida ESCRITA DESPUÉS del error se sigue procesando. Si no se procesa, revisá que `error PUNTO_COMA` esté en el nivel correcto de la gramática.

---

## 4. Tema: el modelo `Conjunto` y las operaciones

### 4.1 `Entorno`: todo el estado de una ejecución

Antes de escribir un solo método, decidí qué necesita vivir junto: constantes del universo, los conjuntos definidos, las operaciones ya evaluadas, la consola, los errores y los tokens. Todo eso es el `Entorno`:

```java
public class Entorno {
    public static final int UNIVERSO_MIN = 33;   // '!'
    public static final int UNIVERSO_MAX = 126;  // '~'

    private final Set<Character> universo = new LinkedHashSet<>();
    private final LinkedHashMap<String, Conjunto> conjuntos = new LinkedHashMap<>();
    private final LinkedHashMap<String, Operacion> operaciones = new LinkedHashMap<>();
    private final StringBuilder consola = new StringBuilder();
    private final List<RegistroError> errores = new ArrayList<>();
    private final List<Object[]> tokens = new ArrayList<>();

    public Entorno() {
        for (int i = UNIVERSO_MIN; i <= UNIVERSO_MAX; i++) universo.add((char) i);
    }
}
```

Tres decisiones a fijar ANTES de escribir los métodos:

- **`LinkedHashMap`, no `HashMap`**: los reportes (Tema 7) necesitan iterar los conjuntos y operaciones en el mismo orden en que se definieron en el código fuente. `HashMap` no te garantiza eso.
- **Las claves NO se normalizan**: acá es donde la sección 4.1 (case sensitive) se vuelve código real. En DataForge hubieras escrito `id.toLowerCase()` antes de usar la clave del mapa; en ConjAnalyzer, **no** — `conjuntos.put(id, ...)` usa el `id` tal cual vino del lexer.
- **Un `Entorno` nuevo por ejecución**: nunca reutilices una instancia entre dos corridas. La sección 5 del enunciado exige que los reportes sean SOLO del último análisis — la forma más simple de garantizar eso es que cada ejecución empiece con un objeto fresco, sin arrastrar estado. Vas a ver esto en el Tema 3 de la GUI: `Interprete.ejecutar(codigo)` crea un `Entorno` nuevo en cada llamada.

### 4.2 `Conjunto`: registro simple

```java
public class Conjunto {
    public final String nombre;
    public final Set<Character> elementos;
    public final String definicion;   // "a~z" o "1, 2, 3, a, b" (para mostrar en reportes)
    public final int linea;
    public final int columna;
}
```

Guardá la `definicion` como texto legible (no solo el `Set` calculado) porque el reporte de operaciones (Tema 7) quiere mostrar CÓMO se declaró el conjunto, no solo su contenido.

### 4.3 Definí conjuntos: rango y lista

```java
public void definirConjuntoRango(String id, String a, String b, int l, int c) {
    if (existe(id, l, c)) return;
    Character ini = charUnico(a, l, c);
    Character fin = charUnico(b, l, c);
    if (ini == null || fin == null) return;
    if (ini > fin) {
        error("Semantico", "rango invalido en '" + id + "': '" + a
                + "' debe ser menor o igual que '" + b + "'", l, c);
        return;
    }
    LinkedHashSet<Character> set = new LinkedHashSet<>();
    for (char ch = ini; ch <= fin; ch++) set.add(ch);
    conjuntos.put(id, new Conjunto(id, set, a + "~" + b, l + 1, c + 1));
}
```

Fijate el patrón `if (... == null) return;` repetido: cada verificación que falla **ya registró su propio error** dentro de `charUnico`/`existe` antes de salir. Este es el patrón de **propagación por null** que vas a usar en todo el proyecto: una función que detecta un problema semántico registra el error y devuelve `null` (o simplemente no hace nada), y quien la llamó revisa el resultado antes de seguir. Nunca lances una excepción para un error semántico esperable — eso mataría el análisis completo, violando "ningún error aborta la ejecución".

### 4.4 Definí operaciones: acá se conectan sintaxis, evaluación y simplificación

```java
public void definirOperacion(String id, NodoOperacion arbol, int l, int c) {
    if (existeOperacion(id, l, c)) return;

    LinkedHashSet<String> refs = new LinkedHashSet<>();
    arbol.referencias(refs);
    for (String ref : refs) {
        if (!conjuntos.containsKey(ref)) {
            error("Semantico", "la operacion '" + id + "' referencia el conjunto '"
                    + ref + "', que no ha sido definido", l, c);
            return;
        }
    }

    Map<String, Set<Character>> mapa = new LinkedHashMap<>();
    for (var e : conjuntos.entrySet()) mapa.put(e.getKey(), e.getValue().elementos);

    Set<Character> resultado = arbol.evaluar(mapa, universo);
    if (resultado == null) resultado = new LinkedHashSet<>();

    ResultadoSimplificacion simpl = new Simplificador().simplificar(arbol);

    operaciones.put(id, new Operacion(id, arbol, resultado, refs, simpl, l + 1, c + 1));
}
```

Este método hace, en orden, las tres cosas que motivaron el árbol acotado del Tema 3: valida referencias, evalúa, y simplifica — sobre el MISMO `arbol`, sin volver a parsear nada. Escribilo así, en ese orden, porque simplificar antes de validar podría intentar reescribir un árbol con referencias rotas.

### 4.5 `NodoOperacion.evaluar`: el primer recorrido

```java
public Set<Character> evaluar(Map<String, Set<Character>> conjuntos, Set<Character> universo) {
    if (esHoja()) {
        Set<Character> base = conjuntos.get(nombreConj);
        return (base == null) ? null : new LinkedHashSet<>(base);
    }
    Set<Character> a = izq.evaluar(conjuntos, universo);
    if (a == null) return null;
    if (esUnario()) {
        Set<Character> r = new LinkedHashSet<>(universo);
        r.removeAll(a);
        return r;
    }
    Set<Character> b = der.evaluar(conjuntos, universo);
    if (b == null) return null;
    Set<Character> r = new LinkedHashSet<>(a);
    switch (op) {
        case 'U' -> r.addAll(b);
        case '&' -> r.retainAll(b);
        case '-' -> r.removeAll(b);
    }
    return r;
}
```

Es un recorrido en **postorden**: para calcular el resultado de un nodo binario necesitás el resultado de AMBOS hijos primero — no podés calcular `A U B` sin tener antes los conjuntos evaluados de `A` y de `B`. El complemento es el único caso unario: universo menos el conjunto del único hijo.

### 4.6 `EVALUAR`: pertenencia, y el caso de estudio de la sección 4.8

```java
public void evaluar(ArrayList<String> datos, String operacion, int l, int c) {
    Operacion op = operaciones.get(operacion);
    if (op == null) {
        error("Semantico", "la operacion '" + operacion + "' no ha sido definida", l, c);
        return;
    }
    consola.append("===============\n").append("Evaluar: ").append(operacion).append('\n').append("===============\n");
    for (String dato : datos) {
        Character ch = charUnico(dato, l, c);
        if (ch == null) continue;
        boolean pertenece = op.resultado.contains(ch);
        consola.append(ch).append(" -> ").append(pertenece ? "exitoso" : "fallo").append('\n');
    }
    consola.append('\n');
}
```

Este método parece trivial, pero escondé acá una decisión de diseño que vale la pena que discutas en tu revisión oral: la sección 4.7 del enunciado dice, literalmente, "validar que cierto conjunto de datos pertenezca al conjunto resultante" — así que la pertenencia se mide contra `op.resultado`, el `Set<Character>` YA evaluado. Guardalo en la cabeza para el Tema 8, donde vas a ver que el propio ejemplo del enunciado (sección 4.8) contradice su propia definición de conjuntos.

---

## 5. Tema: el `Simplificador` — las 5 leyes algebraicas

### 5.1 El contrato: punto fijo, postorden, a lo sumo una regla por nodo

La sección 7 del enunciado lista **13 propiedades** de teoría de conjuntos, agrupadas en 7 familias: doble complemento, DeMorgan, conmutativa, asociativa, distributiva, idempotente y absorción. Antes de escribir código, tomá una decisión de diseño explícita — **no todas se implementan como reescritura**:

- **5 familias SÍ reescriben el árbol** (reducen su tamaño): doble complemento, DeMorgan, idempotentes, absorción, distributivas.
- **2 familias NO reescriben nada por sí solas**: conmutativa y asociativa. Reordenar o reagrupar un árbol no lo simplifica — en el peor caso, podría hacer que el motor cicle sin converger. En cambio, se usan como **ayuda de comparación**: para reconocer que `A U B` y `B U A` son "el mismo árbol" a pesar de estar escritos distinto.

El motor central es un bucle de punto fijo:

```java
public ResultadoSimplificacion simplificar(NodoOperacion original) {
    leyes = new LinkedHashSet<>();
    NodoOperacion actual = original.copia();
    String antes;
    int guarda = 0;
    do {
        antes = actual.toPrefijo();
        actual = paso(actual);
    } while (!actual.toPrefijo().equals(antes) && ++guarda < MAX_ITER);
    return new ResultadoSimplificacion(new ArrayList<>(leyes), actual, !leyes.isEmpty());
}
```

Fijate en `original.copia()`: el `Simplificador` trabaja sobre una COPIA, nunca sobre el árbol original — porque ese mismo árbol original lo necesita el diagrama de Venn (Tema 6) sin modificar. Escribí `copia()` en `NodoOperacion` antes de seguir:

```java
public NodoOperacion copia() {
    if (esHoja()) return hoja(nombreConj);
    if (esUnario()) return unario(op, izq.copia());
    return binario(op, izq.copia(), der.copia());
}
```

El bucle compara `actual.toPrefijo()` antes y después de cada `paso(...)`: mientras la serialización cambie, seguí pasando. `MAX_ITER = 100` es una guarda de seguridad, no un límite que esperes alcanzar con las reglas bien escritas — es la misma prudencia que un compilador real aplica ante cualquier reescritura iterativa que en teoría debería converger pero que un bug podría hacer oscilar.

### 5.2 `paso(...)`: postorden, una regla por nodo

```java
private NodoOperacion paso(NodoOperacion n) {
    if (n.esHoja()) return n;

    if (n.esUnario()) {
        NodoOperacion hijo = paso(n.izq);
        if (hijo.esUnario()) {                // ^^X -> X
            leyes.add("Ley del doble complemento");
            return hijo.izq;
        }
        if (esBinarioUnionInter(hijo) && (hijo.izq.esUnario() || hijo.der.esUnario())) {
            leyes.add("Leyes de DeMorgan");
            char nuevo = (hijo.op == 'U') ? '&' : 'U';
            return NodoOperacion.binario(nuevo,
                    NodoOperacion.unario('^', hijo.izq),
                    NodoOperacion.unario('^', hijo.der));
        }
        return NodoOperacion.unario('^', hijo);
    }

    NodoOperacion a = paso(n.izq);
    NodoOperacion b = paso(n.der);

    if (n.op == 'U' || n.op == '&') {
        if (equivalentes(a, b)) {
            leyes.add("Propiedades idempotentes");
            return a;
        }
        NodoOperacion abs = absorcion(n.op, a, b);
        if (abs != null) return abs;
        NodoOperacion dist = distributiva(n.op, a, b);
        if (dist != null) return dist;
    }
    return NodoOperacion.binario(n.op, a, b);
}
```

Escribí primero **doble complemento** e **idempotentes** — son las más simples y te dan la forma del método. Corré `entradas/ejemplo3_simplificacion.ca` (o un archivo propio con `OPERA : doble -> ^ ^ {conjA};`) después de cada ley nueva que agregues, no esperes a tener las 5.

### 5.3 La guarda de DeMorgan: el patrón "solo si reduce"

Fijate que DeMorgan **no se aplica siempre** que ve `^(X U Y)`. La condición extra es:

```java
if (esBinarioUnionInter(hijo) && (hijo.izq.esUnario() || hijo.der.esUnario())) {
```

Esto exige que **al menos uno** de los operandos (`X` o `Y`) ya sea, a su vez, un complemento. La razón: si aplicás DeMorgan sin esta guarda, `^(A U B)` se reescribe a `^A & ^B` — dos negaciones nuevas, ningún término menos. El árbol NO se simplificó, solo cambió de forma, y en la próxima pasada podría volver a aplicarse la ley inversa y oscilar para siempre. La guarda garantiza que, si dispara, genera un `^^` que la ley del doble complemento va a cancelar en la pasada siguiente (o en la misma, según el orden de evaluación) — una ganancia neta real.

Guardate esta pregunta para tu revisión oral: *"¿por qué DeMorgan necesita una guarda y doble complemento no?"* — porque doble complemento SIEMPRE reduce (dos nodos por uno), mientras que DeMorgan por sí sola no reduce nada; solo lo hace cuando destapa una cancelación posterior.

### 5.4 Absorción

```java
private NodoOperacion absorcion(char op, NodoOperacion a, NodoOperacion b) {
    char interno = (op == 'U') ? '&' : 'U';
    if (esBinario(b, interno) && (equivalentes(b.izq, a) || equivalentes(b.der, a))) {
        leyes.add("Propiedades de absorcion");
        return a;
    }
    if (esBinario(a, interno) && (equivalentes(a.izq, b) || equivalentes(a.der, b))) {
        leyes.add("Propiedades de absorcion");
        return b;
    }
    return null;
}
```

`A U (A & B) = A`: el patrón es "un operando simple, y el otro es una operación del signo OPUESTO que contiene al primero". Fijate que se revisan las dos variantes conmutadas (`A op (A interno B)` y `(A interno B) op A`) — sin esa segunda rama, `(A & B) U A` no se reconocería como absorción, solo `A U (A & B)`.

### 5.5 Distributiva — la ley que se agregó en la auditoría de calidad (2026-07-21)

Esta es la última de las 5, y la más instructiva porque tiene una historia real: en una primera versión de este proyecto, el `Simplificador` implementaba 4 de las 5 leyes y **omitía la distributiva completamente** — ni siquiera como caso de "no simplificable". Una auditoría de calidad posterior, cruzando el grafo de llamadas del proyecto contra el enunciado, encontró el gap: la sección 7 lista la distributiva explícitamente, y la sección 5.4 (JSON) pide aplicar TODAS las leyes de la sección 7, no un subconjunto. A diferencia de conmutativa/asociativa (que están ausentes como regla de reescritura **a propósito**, con una razón matemática documentada), la distributiva sí tenía un sentido que reduce el árbol — omitirla era un bug real, no una decisión de diseño.

Escribila así, siguiendo el mismo patrón de guarda que DeMorgan:

```java
private NodoOperacion distributiva(char op, NodoOperacion a, NodoOperacion b) {
    char interno = (op == 'U') ? '&' : 'U';
    if (!esBinario(a, interno) || !esBinario(b, interno)) return null;

    NodoOperacion factor, restoA, restoB;
    if (equivalentes(a.izq, b.izq))      { factor = a.izq; restoA = a.der; restoB = b.der; }
    else if (equivalentes(a.izq, b.der)) { factor = a.izq; restoA = a.der; restoB = b.izq; }
    else if (equivalentes(a.der, b.izq)) { factor = a.der; restoA = a.izq; restoB = b.der; }
    else if (equivalentes(a.der, b.der)) { factor = a.der; restoA = a.izq; restoB = b.izq; }
    else return null;

    leyes.add("Propiedades distributivas");
    return NodoOperacion.binario(interno, factor, NodoOperacion.binario(op, restoA, restoB));
}
```

El patrón que se repite en TODA la clase de leyes de este `Simplificador` — grabátelo, porque te lo van a preguntar en la defensa oral: **una ley se implementa como reescritura solo si existe una dirección en la que reduce el árbol; nunca en la dirección que lo expande**. Acá, `(X & Y) U (X & Z) → X & (Y U Z)` factoriza (reduce hojas repetidas de `X`); la dirección inversa, `X & (Y U Z) → (X & Y) U (X & Z)`, "expande" un factor y agrega hojas — el método `distributiva(...)` deliberadamente NO implementa esa dirección. Es la misma lógica exacta que la guarda de DeMorgan, aplicada a una ley distinta.

Para verificar que quedó bien conectada, agregá un caso a tu archivo de pruebas:

```
OPERA : distributiva -> U & {conjA} {conjB} & {conjA} {conjC};
```

y confirmá en el JSON de salida (Tema 7) que aparece `"leyes": ["Propiedades distributivas"]` y `"conjunto simplificado": "& {conjA} U {conjB} {conjC}"`.

### 5.6 `equivalentes` y `canon`: por qué conmutativa/asociativa SÍ aparecen en el reporte

Comparar dos subárboles por su texto (`toPrefijo()`) no alcanza: `U {A} {B}` y `U {B} {A}` son el MISMO conjunto matemáticamente, pero textos distintos. La solución es una forma canónica:

```java
private String canon(NodoOperacion n) {
    if (n.esHoja()) return "'" + n.nombreConj + "'";
    if (n.esUnario()) return "^(" + canon(n.izq) + ")";
    if (n.op == '-') return "-(" + canon(n.izq) + "," + canon(n.der) + ")";
    List<String> ops = new ArrayList<>();
    aplanar(n, n.op, ops);
    List<String> unicos = new ArrayList<>(new LinkedHashSet<>(ops));
    Collections.sort(unicos);
    return n.op + "[" + String.join(",", unicos) + "]";
}
```

Para cadenas de `U`/`&`, `canon` aplana los operandos recursivamente (así `(A U B) U C` y `A U (B U C)` producen la misma lista), los ordena alfabéticamente y elimina duplicados. Dos árboles son `equivalentes(...)` si su `canon()` coincide. Cuando esa comparación "no trivial" (los textos no eran idénticos) es la que habilita idempotencia, absorción o distributiva, el motor agrega también `"Propiedades conmutativas"` y/o `"Propiedades asociativas"` a la lista de leyes — eso es lo que documenta, con evidencia, que se usaron sin ser reglas de reescritura independientes.

---

## 6. Tema: diagrama de Venn con Canvas

### 6.1 La idea: sombrear por función booleana, no por dibujo aproximado

La sección 5.1 pide un diagrama de Venn navegable por operación. La tentación es "dibujar círculos y sombrear a ojo" — no lo hagas así. La forma correcta (y la que ya usaste sin saberlo en el Tema 4) es preguntarte, por cada píxel del `Canvas`: *¿a qué conjuntos base pertenece este punto?*, y evaluar el árbol de la operación como una función booleana sobre esa respuesta. Escribí ese segundo recorrido en `NodoOperacion`:

```java
public boolean pertenenciaRegion(Map<String, Boolean> region) {
    if (esHoja()) return region.getOrDefault(nombreConj, false);
    if (esUnario()) return !izq.pertenenciaRegion(region);
    boolean a = izq.pertenenciaRegion(region);
    boolean b = der.pertenenciaRegion(region);
    return switch (op) {
        case 'U' -> a || b;
        case '&' -> a && b;
        case '-' -> a && !b;
        default  -> false;
    };
}
```

Es el mismo árbol, la misma estructura de `evaluar(...)` del Tema 4 — pero acá los "valores" que fluyen no son `Set<Character>`, son `boolean`. Es la prueba concreta de por qué hacía falta guardar el árbol como estructura: estás recorriéndolo de una SEGUNDA forma completamente distinta, sin volver a parsear nada.

### 6.2 Sombreado píxel a píxel

```java
private static void dibujar(GraphicsContext gc, Operacion op, List<String> bases) {
    int n = bases.size();
    double[][] centros = centros(n);
    double r = (n == 3) ? 100 : 110;

    PixelWriter pw = gc.getPixelWriter();
    Map<String, Boolean> region = new HashMap<>();
    for (int y = 0; y < ALTO; y++) {
        for (int x = 0; x < ANCHO; x++) {
            boolean dentroDeAlguno = false;
            for (int i = 0; i < n; i++) {
                double dx = x - centros[i][0], dy = y - centros[i][1];
                boolean adentro = dx * dx + dy * dy <= r * r;
                region.put(bases.get(i), adentro);
                dentroDeAlguno |= adentro;
            }
            if (op.arbol.pertenenciaRegion(region)) {
                pw.setColor(x, y, dentroDeAlguno ? DENTRO : FUERA);
            }
        }
    }
}
```

Dos detalles de diseño que vas a necesitar explicar:

- **Dos tonos (`DENTRO`/`FUERA`) del mismo color**: el complemento incluye elementos del universo que no caen en NINGÚN círculo dibujado — si solo pintaras `DENTRO`, el complemento de un conjunto se vería vacío en el interior de los círculos y no reflejaría que "todo lo de afuera" también pertenece al resultado.
- **Límite geométrico honesto**: con 1, 2 o 3 conjuntos base hay una disposición de círculos razonable (docta, tipo trébol para 3). Con 4 o más, NO existe una forma geométrica simple de dibujar todas las intersecciones posibles con círculos — en vez de forzar un dibujo engañoso, mostrá el resultado solo como texto (`bases.size() > 3`).

### 6.3 Usá el árbol ORIGINAL, no el simplificado

`DiagramaVenn.crear(Operacion op)` recibe `op.arbol` — el árbol **antes** de la simplificación del Tema 5 — no `op.simplificacion.simplificado`. Esto es deliberado: el panel de Venn quiere mostrar el resultado real de la operación tal como el usuario la escribió; la simplificación se muestra aparte, como texto informativo (`op.simplificacion.simplificado.toPrefijo()`), sin afectar el dibujo. Si el `Simplificador` reescribe el árbol de forma incorrecta en algún caso raro, el diagrama de Venn seguiría siendo correcto porque no depende de esa reescritura.

---

## 7. Tema: reportes JSON con Gson

### 7.1 El formato exacto de la sección 5.4

El enunciado especifica el JSON con dos formas posibles de valor por operación: un objeto `{"leyes": [...], "conjunto simplificado": "..."}`, o el **string literal** `"No se puede simplificar la operacion"` si no se aplicó ninguna ley. No es un objeto con arreglo vacío — es el string tal cual, como valor directo de la clave.

```java
public static String construir(Entorno ent) {
    LinkedHashMap<String, Object> raiz = new LinkedHashMap<>();
    for (Operacion op : ent.getOperaciones().values()) {
        if (op.simplificacion.seSimplifico) {
            Map<String, Object> detalle = new LinkedHashMap<>();
            detalle.put("leyes", op.simplificacion.leyes);
            detalle.put("conjunto simplificado", op.simplificacion.simplificado.toPrefijo());
            raiz.put(op.nombre, detalle);
        } else {
            raiz.put(op.nombre, "No se puede simplificar la operacion");
        }
    }
    return GSON.toJson(raiz);
}
```

El truco de tipos es `Map<String, Object>`: la misma estructura (`raiz`) necesita aceptar tanto un `Map` anidado como un `String` plano según el caso — por eso el valor se declara `Object`, no `String` ni `Map`.

### 7.2 `disableHtmlEscaping()` no es opcional

```java
private static final Gson GSON = new GsonBuilder().setPrettyPrinting().disableHtmlEscaping().create();
```

Por defecto, Gson escapa caracteres como `&` a `&` pensando que el JSON se va a incrustar en HTML. En ConjAnalyzer, `&` es un carácter LEGÍTIMO del lenguaje (el operador de intersección) que aparece todo el tiempo en `"conjunto simplificado"`. Si te olvidás de `disableHtmlEscaping()`, tu JSON va a tener `&` en vez de `&` y no vas a poder copiar-pegar la expresión de vuelta como código válido — probalo una vez sin la opción para ver el problema con tus propios ojos, y después agregala.

### 7.3 Reportes HTML: nombres de token por reflexión

Para `tokens.html` necesitás el NOMBRE de cada token (`"CONJ"`, `"ID"`, etc.), no solo su valor numérico interno. En vez de escribir un `switch` manual que se desincroniza cada vez que agregás un terminal al `.cup`, usá reflexión sobre la clase `sym` que CUP genera:

```java
private static LinkedHashMap<Integer, String> nombresTokens() throws Exception {
    LinkedHashMap<Integer, String> m = new LinkedHashMap<>();
    for (Field f : sym.class.getFields()) {
        if (f.getType() == int.class) m.put(f.getInt(null), f.getName());
    }
    return m;
}
```

`sym` declara cada terminal como una constante `public static final int`, con el mismo nombre que le diste en el `.cup` (`CONJ`, `UNION`, `ID`, ...). Recorrer sus campos te da el mapa completo sin mantenimiento manual — si mañana agregás un token `XOR`, este método lo reconoce automáticamente, sin tocar una línea de `Reportes.java`.

### 7.4 Escapá HTML en los valores insertados

Un lexema real puede contener `&` o `<` (por ejemplo, el token `->` no, pero un `SIMBOLO` como `<` sí). Si lo insertás crudo en una celda de tabla HTML, rompés el documento:

```java
private static String esc(String s) {
    if (s == null) return "";
    return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;");
}
```

Aplicá `esc(...)` a TODO lexema o descripción de error que vaya a una celda — no solo a los que "parecen" peligrosos.

---

## 8. Errores comunes reales al programar cada tema

Estos son errores que vas a cometer probablemente (o que se cometieron en la construcción real de este proyecto), agrupados por tema:

**Léxico (Tema 2)**
- Olvidar `%public` en el `.flex` → `Lexer` queda package-private, no compila desde `conjanalyzer.interprete.Interprete`.
- Agregar `%ignorecase` "por costumbre" (veniendo de DataForge) → rompe la sección 4.1, que exige case sensitive real.
- Poner `{Id}` antes que las palabras reservadas → `CONJ` se reconocería como un `ID` cualquiera, nunca como la palabra reservada.
- Olvidar el comodín `SIMBOLO` → conjuntos de puntuación (`!,?,@,$,%`) se volverían imposibles de declarar.

**Sintáctico (Tema 3)**
- Declarar una producción de `non terminal` DESPUÉS de usarla en una producción → error de compilación de CUP; todas las declaraciones van primero.
- Poner la alternativa `error PUNTO_COMA` en el nivel equivocado de la gramática (por ejemplo, dentro de `operacion` en vez de `sentencia`) → el modo pánico sincroniza en el lugar equivocado y perdés más contenido del necesario tras un error.
- Tipar `operacion` como `void` en vez de `NodoOperacion` "para simplificar" → sin el árbol, la simplificación (Tema 5) y el diagrama de Venn (Tema 6) se vuelven imposibles de implementar después sin reescribir el parser.

**Ejecución (Tema 4)**
- Normalizar las claves de `conjuntos`/`operaciones` a minúsculas (copiando el hábito de DataForge) → rompe la sección 4.1.
- Lanzar una excepción Java ante un error semántico (por ejemplo, `throw new RuntimeException(...)` cuando un conjunto no existe) → mata el análisis completo; la forma correcta es `error(...)` + `return` (propagación por null).
- Reutilizar un mismo `Entorno` entre dos ejecuciones sucesivas de la GUI → los reportes mezclan conjuntos/operaciones de corridas distintas, violando la sección 5.

**Simplificador (Tema 5)**
- Aplicar DeMorgan o distributiva SIN guarda → el árbol puede "inflarse" en vez de reducirse, o el motor puede oscilar sin llegar nunca a un punto fijo.
- Comparar subárboles con `.equals()` de texto en vez de `canon(...)` → `A U B` y `B U A` no se reconocen como el mismo árbol, y absorción/idempotencia fallan en casos válidos.
- **Omitir por completo una ley que el enunciado especifica** (el caso real de la distributiva en la sección 7 de esta guía) — es el hallazgo más serio de una auditoría de calidad real: no es un detalle cosmético, es una propiedad exigida por el enunciado que faltaba en el JSON de salida.

**Reportes (Tema 7)**
- Olvidar `disableHtmlEscaping()` en el `GsonBuilder` → el `&` de la notación prefija sale como `&`.
- Devolver un objeto `{"leyes": [], ...}` vacío en vez del string literal `"No se puede simplificar la operacion"` para las operaciones no simplificables → no cumple el formato exacto de la sección 5.4.
- No escapar HTML en lexemas/descripciones → un token `<` o `&` rompe la tabla del reporte.

### Caso de estudio: cuando el propio enunciado se equivoca (sección 4.8)

Este es el error más instructivo de todo el proyecto, porque no es un bug de programación — es una decisión sobre qué hacer cuando la fuente de verdad (el PDF del enunciado) se contradice a sí misma.

La sección 4.8 define:

```
CONJ : conjuntoA -> 1,2,3,a,b;
CONJ : conjuntoB -> a~z;
OPERA : operacion1 -> & {conjuntoA} {conjuntoB};
EVALUAR ( {1, b} , operacion1 );
```

y el PDF muestra como salida:

```
1 -> exitoso
b -> exitoso
```

Ahora hacé la cuenta vos mismo, con la matemática real: `conjuntoA = {1, 2, 3, a, b}`, `conjuntoB` es el rango `a~z` (no incluye dígitos). La intersección `conjuntoA & conjuntoB` es `{a, b}` — el elemento `1` **no** pertenece a esa intersección, porque `conjuntoB` no lo contiene. La salida correcta, matemáticamente, es:

```
1 -> fallo
b -> exitoso
```

El PDF tiene un error. ¿Qué hacés? La tentación (mala) es "copiar la salida del PDF tal cual, para que coincida" — eso significaría literalmente escribir un bug a propósito, con la excusa de "parecerse al enunciado". La decisión correcta, la que quedó implementada en `Entorno.evaluar(...)`, es:

1. **Priorizar la semántica matemática real** sobre la salida textual del PDF.
2. **Documentar la discrepancia en el propio código**, con un comentario que explique el porqué — no dejarlo como una coincidencia silenciosa que un revisor futuro (o vos mismo, en 6 meses) podría confundir con un bug.

```java
/**
 * ... (El segundo bloque de salida del ejemplo 4.8 del enunciado muestra
 * "1 -> exitoso" para una interseccion que no contiene al 1: eso es una
 * inconsistencia del enunciado; aca se respeta la semantica de conjunto
 * resultante.)
 */
```

La lección general, más allá de ConjAnalyzer: **leé el enunciado con sentido crítico**. Un enunciado de curso es una especificación escrita por personas, y puede tener errores — tu criterio para detectarlos y tu capacidad de defender la decisión con evidencia matemática (no "porque sí") es exactamente lo que un evaluador de OLC1 espera ver en la revisión oral. Si alguna vez tenés dudas sobre si tu programa "debería" reproducir un comportamiento que el enunciado muestra, hacé la cuenta a mano primero — no asumas que el enunciado siempre acierta.

---

*Guía elaborada a partir del código real de `ConjAnalyzer/src/` (paquetes `analisis`, `interprete`, `reportes`, `gui`), verificada contra el enunciado real (`doc/Projects/OLC1_PT1 _S2024.clean.md`, secciones 4 a 7) y contra `ManualTecnico.md` de este mismo proyecto. Refleja el estado del código tras la auditoría de calidad del 2026-07-21 (ley distributiva agregada, `Entorno.getUniverso()` eliminado por no tener llamadores reales).*
