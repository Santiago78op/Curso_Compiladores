# Guía de Programación — Cómo construir CompScript desde cero

Universidad de San Carlos de Guatemala · Facultad de Ingeniería · Escuela de Ciencias y Sistemas
Organización de Lenguajes y Compiladores 1 · Vacaciones de Diciembre 2024

## Cómo usar esta guía

Este documento **no es el Manual Técnico** (que describe lo que YA está construido) ni el Manual de Usuario (que describe cómo se USA el editor). Es material de clase: te lleva, tema por tema, por el mismo camino que se siguió para construir CompScript, con código real del proyecto y explicando el porqué de cada decisión antes de que la veas escrita. Está pensado para que vos programes mientras leés, no para que copies y pegues.

Cada tema cierra con una sección **"Ahora programá vos"** con un ejercicio concreto sobre el código real del repositorio, y con **"Dónde se rompe si...")** señalando el error típico que vas a cometer la primera vez.

La referencia teórica de fondo es el Dragón (Aho/Lam/Sethi/Ullman, 2ª ed.): Cap. 3 para el analizador léxico, Cap. 4 para el sintáctico, Cap. 5 para la traducción dirigida por la sintaxis y los árboles, Cap. 7 para entornos y pila de control. Si ya construiste DataForge, vas a reconocer el 70% de las herramientas (JFlex, CUP, Maven) — lo que cambia de fondo es la **arquitectura**: DataForge no necesita AST, CompScript sí, y ese "sí" es el hilo conductor de todo este documento.

---

## 1. Preparación del proyecto (Maven, JFlex, CUP)

Antes de escribir una sola línea de gramática, necesitás el mismo esqueleto de proyecto que ya armaste en DataForge: un `pom.xml` que declare JFlex y CUP como generadores de código en la fase `generate-sources`, más JavaFX para el editor.

El `pom.xml` real de CompScript (`CompScript/pom.xml`) declara exactamente cuatro plugins de build, cada uno con una responsabilidad:

```xml
<properties>
    <maven.compiler.source>17</maven.compiler.source>
    <maven.compiler.target>17</maven.compiler.target>
</properties>

<dependencies>
    <dependency> <!-- Symbol, la clase que devuelve cada token del lexer -->
        <groupId>com.github.vbmacher</groupId>
        <artifactId>java-cup-runtime</artifactId>
        <version>11b-20160615</version>
    </dependency>
    <dependency> <!-- controles de UI para el editor -->
        <groupId>org.openjfx</groupId>
        <artifactId>javafx-controls</artifactId>
        <version>21.0.4</version>
    </dependency>
</dependencies>
```

Y en `<build><plugins>`: `cup-maven-plugin` (genera `Parser.java` y `sym.java` en `target/generated-sources/cup`), `jflex-maven-plugin` (genera `Lexer.java` en `target/generated-sources/jflex`), `javafx-maven-plugin` (para `mvn clean javafx:run`) y `exec-maven-plugin` apuntando a `compscript.TestInterprete` (para probar por consola sin abrir la GUI).

**Vas a crear estas carpetas** (vacías al principio):

```
CompScript/
├── src/main/jflex/Lexer.flex     ← vos escribís esto
├── src/main/cup/parser.cup       ← vos escribís esto
└── src/main/java/compscript/
    ├── ast/A.java                ← vos escribís esto
    ├── interprete/               ← vos escribís esto (Contexto, Entorno, Tipo, Valor, Simbolo, Operaciones...)
    ├── gui/                      ← vos escribís esto (EditorApp, Lanzador)
    └── reportes/                 ← vos escribís esto (Reportes)
```

Los paquetes `compscript.analisis.Lexer`, `.Parser` y `.sym` **no los escribís vos**: Maven los regenera en cada `compile` a partir de `Lexer.flex` y `parser.cup`. Si los editás a mano, tu cambio desaparece en la próxima compilación.

**Dónde se rompe si...** olvidás correr `mvn compile` (o "Reload All Maven Projects" en IDEA) después de tocar `.flex`/`.cup`: el IDE sigue mostrando `Lexer`/`Parser`/`sym` con el contenido VIEJO, y vas a perseguir un bug que en realidad ya corregiste en el fuente.

### Ahora programá vos
Cloná el `pom.xml` de DataForge, subí `maven.compiler.source/target` a 17, y agregá la dependencia de `javafx-controls`. Corré `mvn compile` con un `Lexer.flex`/`parser.cup` mínimos (un solo token, una sola producción) y confirmá que `target/generated-sources/` aparece con las tres clases generadas antes de seguir.

---

## 2. Análisis léxico — tabla de tokens (`Lexer.flex`)

CompScript tiene más vocabulario que DataForge, pero las reglas de diseño de un `.flex` son las mismas: reservadas antes que el identificador genérico, texto más largo gana, y recuperación de errores sin abortar.

### 2.1 Las opciones de la directiva `%%` inicial

```
%class Lexer
%public
%unicode
%cup
%line
%column
%ignorecase
```

`%public` es obligatorio (sin él la clase queda package-private y nadie fuera de `compscript.analisis` puede instanciarla). `%cup` conecta el lexer con CUP (cada regla devuelve un `Symbol`). `%line`/`%column` habilitan `yyline`/`yycolumn` para reportar posición. `%ignorecase` hace que **tanto las reservadas como los identificadores** sean case-insensitive — por eso `Contexto`/`Entorno` guardan las claves de la tabla de símbolos en minúscula (vas a verlo en el Tema 5).

### 2.2 Las categorías de token que vas a declarar

CompScript reconoce **37 palabras reservadas** y **28 símbolos**, agrupados por función:

```
"int" "double" "bool" "char" "string" "void"      // tipos de dato
"let" "const" "true" "false"                       // declaracion
"if" "else" "match" "default"                      // control
"while" "for" "do"                                 // ciclos
"break" "continue" "return"                        // transferencia
"struct" "list" "cast" "as"                         // tipos compuestos
"console" "log" "push" "get" "set" "remove" "pop" "reverse"   // consola y listas
"round" "length" "tostring" "run_main"              // funciones nativas
```

```
++ -- == != <= >= && || =>      // compuestos (longest match)
+ - * / ^ $ % ! < > =           // simples
( ) [ ] { } , ; : .             // agrupadores y puntuacion
```

Fijate en `^` y `$`: acá NO son XOR ni variable de shell — el enunciado los redefine explícitamente como potencia y raíz (5.5.5, 5.5.6). Es una decisión del enunciado, no una convención que traigas de otro lenguaje; documentala así en tus comentarios para no confundir a quien lea tu `.flex` después.

### 2.3 El orden importa: reservadas ANTES que `{Id}`

```
Letra        = [a-zA-Z_]
Digito       = [0-9]
Id           = {Letra}({Letra}|{Digito})*
...
"int"       { return symbol(sym.T_INT); }
...              /* TODAS las reservadas van aca */
...
{Id}        { return symbol(sym.ID, yytext()); }   /* esta regla va AL FINAL */
```

JFlex, ante un empate de longitud entre dos reglas que matchean el mismo texto, prefiere la regla escrita **primero** en el archivo. Si `{Id}` estuviera declarada antes que `"int"`, la palabra `int` matchearía como identificador genérico y la reservada nunca existiría — el lenguaje entero perdería sus 37 palabras clave de un plumazo.

**Ahora programá vos:** agregá al `.flex` la regla de `"int"` seguida de dos identificadores de prueba (`interno`, `integer`) y confirmá con `TestInterprete` que ambos siguen reconociéndose como `ID` completo (no como `T_INT` + resto) — es la prueba de que el *longest match* de JFlex, no solo el orden, está funcionando.

### 2.4 Tokens con patrón y helpers del lexer

```
Entero       = {Digito}+
Decimal      = {Digito}+"."{Digito}+
Cadena       = \"([^\"\\\n]|\\.)*\"
Caracter     = '([^'\\\n]|\\.)'
```

```java
{Decimal}   { return symbol(sym.DECIMAL, Double.valueOf(yytext())); }
{Entero}    { return symbol(sym.ENTERO, Integer.valueOf(yytext())); }
{Cadena}    { String v = procesarEscapes(yytext().substring(1, yylength()-1));
              return symbol(sym.CADENA, v); }
```

Tres helpers viven en el bloque `%{ ... %}` del lexer:

```java
private Symbol symbol(int type)                 // token sin valor semantico
private Symbol symbol(int type, Object value)    // token con valor semantico
private void registrar(int type)                 // notifica el token al Contexto (reporte 6.1)
private String procesarEscapes(String s)          // traduce \n \t \\ \" \' dentro de cadenas/caracteres
```

**El punto de diseño que vale la pena que entiendas bien:** las secuencias de escape se resuelven **acá, en el lexer**, no más adelante. Cuando el token `CADENA` llega al parser, el `\n` que el programador escribió en el `.cs` ya es un salto de línea real dentro del `String` de Java. Ninguna capa posterior (parser, AST, intérprete) necesita saber que la fuente usaba barras invertidas. Esa es la razón de ser de `procesarEscapes()`: mantiene el conocimiento de "cómo se escribe un escape" encapsulado en un solo lugar.

### 2.5 El lexer no imprime a `stderr`: reporta al `Contexto`

```java
public compscript.interprete.Contexto contexto;   // lo conecta Interprete.ejecutar(...)

private void registrar(int type) {
    if (contexto != null) contexto.registrarToken(yytext(), type, yyline, yycolumn);
}
```

```
[^]         { if (contexto != null) {
                contexto.error("Lexico", "el caracter '" + yytext()
                    + "' no pertenece al lenguaje", yyline, yycolumn);
              } else {
                System.err.printf("Error lexico: '%s' (linea %d, col %d)%n", ...);
              } }
```

La regla catch-all `[^]` (cualquier carácter que ninguna regla anterior reconoció) es la recuperación de errores léxicos: el carácter se **descarta** y el análisis sigue con el siguiente. El `if (contexto != null)` existe porque el lexer también puede correr standalone en pruebas unitarias sin `Interprete` de por medio; en ese caso, cae al `System.err` en vez de fallar con un `NullPointerException`.

### Dónde se rompe si...
- Ponés `{Id}` antes de las reservadas → perdés 37 palabras clave (ver 2.3).
- Te olvidás `%public` → `Lexer` no compila desde `compscript.interprete.Interprete` (paquete distinto).
- Escribís `Cadena = \"[^\"]*\"` (sin permitir `\\.`) → una cadena con `\"` adentro corta el token antes de tiempo.

---

## 3. Análisis sintáctico — gramática con precedencia, *dangling else* (`parser.cup`)

### 3.1 Por qué la precedencia importa acá y no en DataForge

En DataForge, las llamadas tipo `SUM(a, b)` delimitan sus operandos con paréntesis: no hay ambigüedad posible. CompScript usa **notación infija real** (`a + b * c`), así que la producción `expr` es ambigua en su forma BNF pura — `a + b * c` podría agruparse como `(a+b)*c` o `a+(b*c)`, y un parser LALR necesita que le digas cuál.

CUP resuelve esto con una tabla `precedence`, declarada de más débil a más fuerte (**el último declarado liga MÁS fuerte**):

```
precedence left  OR;                                                          // nivel 7 (mas debil)
precedence left  AND;                                                         // nivel 6
precedence right NOT;                                                        // nivel 5
precedence left  IGUAL_IGUAL, DIFERENTE, MENOR, MENOR_IGUAL, MAYOR, MAYOR_IGUAL; // nivel 4
precedence left  MAS, MENOS;                                                  // nivel 3
precedence left  POR, DIV, MOD;                                               // nivel 2
precedence nonassoc POT, RAIZ;                                                // nivel 1
precedence right UMENOS;                                                     // nivel 0 (mas fuerte)
```

Esta tabla es una traducción **literal** de la sección 5.9 del enunciado — no es una elección libre tuya, el enunciado ya fija los 8 niveles y vos solo los pasás al `.cup` en el orden correcto (invertido, porque CUP declara de menor a mayor precedencia). Contrastá esto con el libro clásico, que resolvería la misma ambigüedad partiendo la gramática en niveles (`expr → termino → factor`); usar la tabla de `precedence` es una elección de herramienta igual de válida, y es la que exige la práctica real con CUP.

### 3.2 El truco de `UMENOS`

```
terminal UMENOS;   /* pseudo-terminal: NUNCA lo produce el lexer */
...
| MENOS:m expr:a  {: RESULT = new A.Unaria("NEG", a, mleft, mright); :} %prec UMENOS
```

`-5 + 3` y `10 - 5` usan el mismo token léxico `MENOS`, pero necesitan precedencias distintas: la resta binaria liga con `+`/`-` (nivel 3), pero la negación unaria debe ligar más fuerte que la potencia (nivel 0, para que `-2^2` se lea `-(2^2)` y no `(-2)^2`... según lo que decida tu enunciado). `UMENOS` es un **pseudo-terminal**: no tiene patrón léxico propio (el lexer jamás lo produce), existe solo para que `%prec UMENOS` le asigne a esa producción específica una precedencia diferente de la resta. Es el patrón estándar de CUP/Yacc para "mismo símbolo léxico, dos precedencias distintas".

**Ahora programá vos:** quitá el `%prec UMENOS` de tu copia local y volvé a generar el parser (`mvn compile`). CUP va a reportar un conflicto shift-reduce al compilar — confirmá que lo ves en la consola de Maven, y después restaurá la línea. Es la mejor forma de entender que ese `%prec` no es decorativo.

### 3.3 El *dangling else*

```
if_stmt ::= IF:f PAR_IZQ expr:c PAR_DER bloque:b
            {: RESULT = new A.If(c, b, null, null, fleft, fright); :}
          | IF:f PAR_IZQ expr:c PAR_DER bloque:b ELSE bloque:e
            {: RESULT = new A.If(c, b, e, null, fleft, fright); :}
          | IF:f PAR_IZQ expr:c PAR_DER bloque:b ELSE if_stmt:ei
            {: RESULT = new A.If(c, b, null, ei, fleft, fright); :}
          ;
```

El conflicto clásico de `if (a) if (b) X; else Y;` (¿el `else` es del `if` de adentro o del de afuera?) **no se resuelve con una directiva de precedencia acá**: se resuelve con el orden y la forma de las tres producciones de `if_stmt`. CUP, ante la ambigüedad, prefiere por defecto *shift* sobre *reduce* — es decir, sigue leyendo en busca de un `else` en vez de cerrar el `if` ya. Con esta gramática, eso asocia automáticamente cada `else` con el `if` abierto más cercano, y CUP no reporta ningún conflicto shift-reduce al generar (podés verificarlo vos mismo mirando la salida de `mvn compile`).

### 3.4 Cada acción CONSTRUYE, no calcula

Esta es la diferencia de fondo con DataForge. Mirá la misma clase de producción en los dos proyectos:

```java
// DataForge (S-atribuida): calcula ya
expr ::= expr:a MAS expr:b {: RESULT = Operaciones.aritmetica("+", a, b, ...); :}

// CompScript (AST): construye un objeto que sabe calcular DESPUES
expr ::= expr:a MAS expr:b {: RESULT = new A.Binaria("+", a, b, aleft, aright); :}
```

`parser.raiz` (un `ArrayList` de tipo *raw*, restricción de CUP sobre los no terminales) termina siendo la lista de instrucciones de nivel superior. Nadie ejecutó nada todavía: el árbol es una estructura de datos en memoria, lista para que otra fase la recorra. Ese "por qué" es exactamente el tema siguiente.

### 3.5 Recuperación de errores sintácticos

```java
parser code {:
    public Contexto contexto;
    public void syntax_error(Symbol s) {
        contexto.error("Sintactico",
            "error de sintaxis, no se esperaba '" + s.value + "'", s.left, s.right);
    }
:};
...
instruccion ::= ... | error PUNTO_COMA {: RESULT = null; :} ;
```

`syntax_error(Symbol s)` registra el error con la posición del símbolo inesperado. La producción `instruccion ::= error PUNTO_COMA` es el **modo pánico**: ante un error, CUP descarta símbolos hasta encontrar un `;` y retoma el análisis desde la instrucción siguiente. Esto permite que un solo error de sintaxis no tape el resto de errores del archivo — el mismo mecanismo de DataForge.

### Dónde se rompe si...
- Olvidás `%prec UMENOS` → conflicto shift-reduce al generar (3.2).
- Escribís `if_stmt` con el `else` en el orden equivocado → el conflicto de *dangling else* deja de resolverse como lo esperás (probalo invirtiendo el orden de las dos últimas alternativas y mirá si CUP reporta un conflicto).
- Te olvidás la producción `error PUNTO_COMA` → un solo error sintáctico aborta TODO el análisis en vez de acumular errores.

---

## 4. El AST — por qué acá SÍ hace falta (`ast/A.java`)

### 4.1 El argumento completo

DataForge no construye AST porque no tiene ninguna instrucción que deba ejecutarse una cantidad variable de veces: cada producción de su gramática se reduce una vez, y en ese mismo momento se calcula el resultado final. CompScript tiene `if`, `while`, `for`, `do-while` y **funciones con recursión** — construcciones que necesitan ejecutar el mismo fragmento de sintaxis 0, 1 o N veces, con N decidido en tiempo de ejecución.

Un parser LALR reduce cada producción **una sola vez**, en el orden fijo en que ocurren los shifts/reduces. Si evaluaras directamente en la acción de `while_stmt`, ya "pasaste" por esa reducción: no hay forma de volver a evaluar el mismo cuerpo del ciclo una segunda vez sin re-parsear el texto fuente completo (carísimo y frágil). La solución general — la que documenta el Cap. 5 del Dragón para lenguajes con control de flujo — es separar **"entender la sintaxis"** (una vez, durante el parseo) de **"ejecutar"** (las veces que el programa lo pida): eso es un árbol de sintaxis abstracta.

### 4.2 Las dos interfaces que vas a escribir primero

```java
public interface Nodo {
    String etiquetaAst();
    default List<Nodo> hijosAst() { return new ArrayList<>(); }
}
public interface Instruccion extends Nodo { void ejecutar(Entorno e); }
public interface Expresion  extends Nodo { Valor evaluar(Entorno e); }
```

Esto es todo el contrato. Cada una de las ~25 clases de nodo (`A.If`, `A.While`, `A.Binaria`, `A.Llamada`, `A.DeclaracionFuncion`, etc.) implementa `Instruccion` (hace algo, no produce valor) o `Expresion` (produce un `Valor`). El patrón es *tree-walking interpreter*, una variante de Visitor sin doble despacho: en vez de una clase visitante externa que sepa procesar cada tipo de nodo, **cada nodo sabe procesarse a sí mismo**. `etiquetaAst()`/`hijosAst()` son un beneficio casi gratis de esta misma estructura: alimentan el reporte de AST (Tema 8) sin necesitar una jerarquía de visitantes aparte — el generador de reportes solo recorre `hijosAst()` recursivamente.

### 4.3 Por qué los constructores reciben `Object`

```java
public Binaria(String op, Object izq, Object der, int l, int c) {
    this.op = op; this.izq = (Expresion) izq; this.der = (Expresion) der; ...
}
```

Los no terminales del `.cup` son tipos *raw* (`Object`, `ArrayList` sin genéricos — la misma restricción que ya viste en DataForge). Recibir `Object` y castear DENTRO del constructor mantiene la acción semántica del `.cup` en una sola línea (`new A.Binaria("+", a, b, aleft, aright)`), sin casts repetidos ahí. El costo de esta comodidad es justamente el tema de los errores comunes (Tema 9): un cast mal guardado revienta en tiempo de ejecución.

### 4.4 Un nodo de instrucción real, completo

```java
public static class If implements Instruccion {
    public final Expresion cond; public final List<Instruccion> cuerpo, sino;
    public final If sinoIf; public final int linea, columna;

    public void ejecutar(Entorno e) {
        Valor v = cond.evaluar(e);
        if (v.tipo.cat != Tipo.Cat.BOOL)
            e.errorSemantico("la condicion del if debe ser Booleana", linea, columna);
        if ((Boolean) v.valor) A.ejecutar(cuerpo, e.crearHijo("if"));
        else if (sinoIf != null) sinoIf.ejecutar(e);
        else if (!sino.isEmpty()) A.ejecutar(sino, e.crearHijo("else"));
    }
    public String etiquetaAst() { return "if"; }
    public List<Nodo> hijosAst() { return hijos(cond, cuerpo, sino, sinoIf); }
}
```

Fijate en la secuencia: 1) evalúa la condición, 2) valida que sea `Booleana` (chequeo semántico, no sintáctico — la gramática no puede saber el tipo de `cond`), 3) decide qué rama ejecutar creando un `Entorno` hijo nuevo para esa rama. Este último punto es el gancho al Tema 5: cada bloque que se ejecuta vive en su propio ámbito.

### Ahora programá vos
Escribí `A.Unaria` completo (ya lo tenés en el repo real, pero intentalo antes de mirarlo): la clase implementa `Expresion`, guarda un `op` (`"NEG"` o `"NOT"`) y un `Expresion expr`, y en `evaluar()` decide entre `Operaciones.negacion(...)` y `Operaciones.negacionLogica(...)` según `op`. Después comparalo con `ast/A.java` línea 189 en adelante.

### Dónde se rompe si...
- Casteás `(Expresion) o` sobre un objeto que en realidad es un `A.LiteralStruct` (por ejemplo) → `ClassCastException` en el momento del cast, no antes. El compilador no te avisa: los `Object` del `.cup` no llevan tipo estático.

---

## 5. Tabla de símbolos y ámbitos (`Contexto` / `Entorno`)

### 5.1 Dos responsabilidades, dos clases

`Contexto` es **único** por ejecución: consola de salida, lista de errores, lista de tokens reconocidos, registro plano de TODOS los símbolos alguna vez declarados, y los mapas de funciones/structs registrados en la primera pasada. `Entorno` es **múltiple**: uno por cada bloque que introduce un ámbito nuevo (`if`, `while`, cuerpo de función...).

```java
public class Entorno {
    public final Contexto contexto;
    public final Entorno padre;       // null = entorno global
    private final LinkedHashMap<String, Simbolo> tabla = new LinkedHashMap<>();

    public Entorno crearHijo(String nombre) { return new Entorno(contexto, this, nombre); }

    public Simbolo buscar(String id) {
        String clave = id.toLowerCase();
        for (Entorno e = this; e != null; e = e.padre) {
            Simbolo s = e.tabla.get(clave);
            if (s != null) return s;
        }
        return null;
    }
}
```

`buscar()` sube por la cadena de padres hasta encontrar el identificador o llegar al ámbito global — eso es **alcance estático anidado**: si dos ámbitos distintos declaran una variable con el mismo nombre, la del ámbito más cercano "tapa" (*shadowing*) a la de más afuera mientras dure ese bloque, sin perderla.

### 5.2 La decisión que más distingue a CompScript de un lenguaje sin funciones

```java
// A.java, dentro de invocar(...)
Entorno local = ctx.global.crearHijo(f.id);   // SIEMPRE del global, NUNCA del caller
```

El entorno de una función cuelga **siempre** del entorno global, nunca del entorno de quien la invoca. Esto es lo que le da a CompScript alcance estático real: una función no puede ver las variables locales de quien la llama. Si en cambio colgara del `caller`, tendrías *dynamic scoping* — cada función vería el estado local de su invocador, un modelo mental mucho más difícil de razonar y que rompería las expectativas de cualquier lenguaje "normal".

**Evidencia real, no solo teórica:** `ejemplo5.cs` corre `fibonacci(n=10)`, que dispara ~177 llamadas recursivas. El reporte de la tabla de símbolos (6.4) muestra **186 entradas** distintas para el parámetro `n` — una por cada invocación, cada una en su propio `Entorno` hijo del global. Esa cifra es la prueba tangible de que cada llamada crea un ámbito nuevo.

### 5.3 El registro de símbolos es un LOG, no un snapshot

```java
public void declarar(Simbolo s, int linea, int columna) {
    String clave = s.nombre.toLowerCase();
    if (tabla.containsKey(clave))
        errorSemantico("'" + s.nombre + "' ya fue declarada en este ambito", linea, columna);
    tabla.put(clave, s);
    contexto.simbolos.add(s);   // registro plano para el reporte 6.4 -- NUNCA se sobreescribe
}
```

`tabla` (privada de cada `Entorno`) sí se comporta como un mapa normal: una clave, un valor, se puede sobreescribir. Pero `contexto.simbolos` (compartido, del `Contexto`) solo **crece**: cada `declarar(...)` agrega una entrada nueva, nunca reemplaza una vieja. Por eso el reporte de tabla de símbolos (6.4) muestra el historial completo de una ejecución con recursión, no solo el estado final — es una decisión de diseño explícita para cumplir lo que pide el reporte, no un descuido.

### 5.4 `Tipo`, `Valor` y `Simbolo`: quién guarda qué

- `Tipo` (enum `Cat`: `INT DOUBLE BOOL CHAR STRING VOID NULL VECTOR LIST STRUCT`) modela la categoría, con `elemento`/`dimensiones` para vectores y listas, y `structName` para structs. Implementa **igualdad estructural** (`equals()`): dos vectores son el mismo tipo si comparten `elemento` y `dimensiones`.
- `Valor` empareja un `Tipo` con el objeto Java que lo representa (`Integer`, `Double`, `Boolean`, `Character`, `String`, `List<Valor>` para vectores/listas, `LinkedHashMap<String,Valor>` para structs).
- `Simbolo` es una fila de la tabla: nombre, categoría (`Simbolo.categoriaDe(Tipo)` → `Variable`/`Vector`/`Lista`/`Struct`), tipo, mutabilidad (`let` vs `const`), ámbito, valor (mutable) y posición.

### Ahora programá vos
Escribí a mano, en un papel o comentario, la secuencia de `Entorno`s que se crean al ejecutar `factorial(5)` de `ejemplo5.cs` (cada llamada recursiva). Confirmá cuántos `Entorno`s hijos del global existen simultáneamente en el punto de máxima recursión (`n=1`), y qué pasa con ellos cuando cada llamada retorna.

### Dónde se rompe si...
- Colgás `local` del `caller` en vez del global → *dynamic scoping* no deseado (5.2). Es el bug de diseño más peligroso de esta sección porque el programa "funciona" en casos simples y falla silenciosamente en recursión con nombres de variable repetidos.

---

## 6. Flujo de control sobre el árbol — if/while/for, y el cortocircuito como caso de estudio real

### 6.1 El patrón Java que hace posible "repetir"

```java
public static class While implements Instruccion {
    public void ejecutar(Entorno e) {
        while (true) {
            Valor v = cond.evaluar(e);
            if (v.tipo.cat != Tipo.Cat.BOOL)
                e.errorSemantico("la condicion del while debe ser Booleana", linea, columna);
            if (!(Boolean) v.valor) break;
            try {
                A.ejecutar(cuerpo, e.crearHijo("while"));
            } catch (Senales.Break b) { break; }
            catch (Senales.Continue ct) { /* siguiente iteracion */ }
        }
    }
}
```

Este `while(true) { evaluar cond; if(!cond) break; ejecutar cuerpo; catch Break/Continue }` es el mecanismo Java LITERAL que hace posible recorrer el mismo subárbol (`cuerpo`) tantas veces como la condición lo permita — la contraparte concreta de lo que en el Tema 4 quedó como argumento teórico ("necesitás poder re-evaluar el mismo subárbol"). `A.For` y `A.DoWhile` siguen la misma estructura, con dos matices: `For` crea su propio `Entorno` para la variable de control (persiste entre iteraciones) y un `Entorno` hijo nuevo por cada ejecución del cuerpo; `DoWhile` evalúa la condición **después** del cuerpo, garantizando al menos una ejecución (`ejemplo3.cs`, variable `k`, imprime `0 1`: dos ejecuciones mínimo).

### 6.2 `break`/`continue`/`return` como excepciones, no banderas

```java
public static class Break implements Instruccion {
    public void ejecutar(Entorno e) { throw new Senales.Break(linea, columna); }
}
```

```java
public static class Break extends RuntimeException {
    public Break(int l, int c) { super(null, null, false, false); linea = l; columna = c; }
}
```

`Senales.Break/Continue/Retorno` extienden `RuntimeException` pero se construyen con `super(null, null, false, false)` — **sin mensaje, sin causa, sin stack trace, sin supresión**. Son control de flujo, no errores: el costo de capturar un stack trace en Java es real, y acá no hace falta ninguno. La ventaja de usar excepciones en vez de una bandera booleana que revisás después de cada instrucción: un `break` puede estar anidado dentro de varios `if` sin código manual — Java "desenrolla" automáticamente esas capas hasta el primer `catch (Senales.Break b)` que encuentra (el ciclo más cercano).

### 6.3 El bug real de cortocircuito — antes y después

`A.Binaria.evaluar()` decide, para cada operador, **qué operandos evaluar antes de invocar a `Operaciones`**. La aritmética y las relacionales siempre evalúan ambos lados (no hay forma de saltarse ninguno: `a + b` necesita los dos valores). Pero `&&` y `||` son distintos: si el operando izquierdo ya determina el resultado, evaluar el derecho es innecesario — y en algunos casos, **peligroso**.

**Antes de la corrección** (auditoría de código, 2026-07-21), el código evaluaba `izq` y `der` SIEMPRE, antes de entrar al `switch` que decide el operador:

```java
// ANTES (bug): evalua los dos lados sin importar el operador
public Valor evaluar(Entorno e) {
    Valor a = izq.evaluar(e);
    Valor b = der.evaluar(e);          // <- esto corria SIEMPRE, incluso para && / ||
    switch (op) {
        case "&&": case "||": return Operaciones.logica(op, a, b, e, linea, columna);
        ...
    }
}
```

Este bug es un ejemplo perfecto de "pasa los tests pero está mal": el enunciado no exige cortocircuito explícitamente, así que ningún requisito lo iba a atrapar; y ninguno de los 6 ejemplos de `entradas/` tenía un caso que dependiera de él. Recién con una guarda típica se hace visible: `x != 0 && 10 / x > 1`. Si `x` vale `0`, el cortocircuito de cualquier lenguaje normal evita evaluar `10/x`; sin cortocircuito, `10/x` se evalúa SIEMPRE, y explota con "división entre cero" — justo lo que la guarda estaba tratando de evitar.

**Después de la corrección**, el código real (`ast/A.java`):

```java
case "&&": {
    // Cortocircuito: si izq es un Booleano falso, el resultado es false
    // sin necesidad de evaluar der -- importa cuando der tiene efectos
    // secundarios o falla en tiempo de ejecucion (p. ej. "x != 0 && 10/x > 1").
    Valor a = izq.evaluar(e);
    if (a.tipo.cat == Tipo.Cat.BOOL && !(Boolean) a.valor) return Valor.vBool(false);
    return Operaciones.logica(op, a, der.evaluar(e), e, linea, columna);
}
case "||": {
    // Cortocircuito simetrico: si izq es Booleano verdadero, el resultado
    // es true sin evaluar der.
    Valor a = izq.evaluar(e);
    if (a.tipo.cat == Tipo.Cat.BOOL && (Boolean) a.valor) return Valor.vBool(true);
    return Operaciones.logica(op, a, der.evaluar(e), e, linea, columna);
}
```

La clave está en que `der.evaluar(e)` ahora aparece **dentro** de la llamada a `Operaciones.logica(...)`, como argumento — es decir, solo se ejecuta si la línea anterior no hizo `return` antes de llegar ahí. Esta es la lección para programar el resto del intérprete: **el orden textual del código Java determina qué se evalúa y cuándo**; no asumas que "evaluar ambos lados antes del switch" es gratis, porque en un lenguaje con efectos secundarios (impresión, errores en tiempo de ejecución, futuras funciones con side effects) ese orden es observable.

### Ahora programá vos
Escribí `entradas/demo_cortocircuito.cs` con:
```
void main() {
    let x: int = 0;
    console.log(x != 0 && 10 / x > 1);
}
RUN_MAIN main();
```
Corré `mvn compile exec:java -Dexec.args="entradas/demo_cortocircuito.cs"` y confirmá que imprime `false` sin ningún error de división entre cero en la sección `--- ERRORES ---`.

### Dónde se rompe si...
- Evaluás `der` antes del `switch` (como en la versión con bug) → cualquier guarda típica (`x != null && x.algo`, `i < n && arreglo[i] > 0`) deja de proteger nada.
- Olvidás el `catch (Senales.Continue ct)` en `While`/`For`/`DoWhile` → un `continue` se propaga como una señal no capturada hasta `Interprete.ejecutar(...)`, que lo reporta como error semántico ("`continue` fuera de un ciclo") aunque sí estuviera dentro de uno.

---

## 7. Funciones y recursión — 3 pasadas, argumentos por nombre

### 7.1 Por qué 3 pasadas y no 1

```java
// Interprete.java
// 1a pasada: registro de funciones/metodos/structs (permite forward-reference)
for (A.Instruccion i : raiz) {
    if (i instanceof A.DeclaracionFuncion f) ctx.registrarFuncion(f);
    else if (i instanceof A.DeclaracionStruct s) ctx.registrarStruct(s);
}
// 2a pasada: declaraciones y asignaciones globales
// 3a pasada: RUN_MAIN (punto de entrada)
```

El enunciado (5.23) permite llamar una función antes de que aparezca textualmente en el archivo. Una sola pasada de ejecución de arriba hacia abajo no podría resolver eso: si `main()` (definida primero) llama a `ayuda()` (definida después), y ejecutás instrucción por instrucción en orden textual, `ayuda` todavía no existiría en la tabla de funciones cuando `main` la necesita. La primera pasada resuelve esto **registrando** (no ejecutando) todas las funciones y structs antes de correr una sola instrucción real — así, para cuando la segunda y tercera pasada ejecutan código, el mapa `ctx.funciones` ya está completo.

### 7.2 Argumentos por identificador, no por posición

```
<llamada> ::= ID "(" <lista-argumentos> ")"
<argumento> ::= ID "=" <expresion>
```

```
factorial(n = 5);
potencia(base = 4);                 // exp usa su valor por defecto
potencia(base = 2, exp = 5);        // orden invertido tambien valido
```

A diferencia de la mayoría de los lenguajes que conocés, en CompScript **el orden de los argumentos no importa**, solo el nombre. Esto es una decisión explícita del enunciado (5.23), no una elección libre tuya al implementar.

### 7.3 `invocar(...)`: el corazón de la resolución de llamadas

```java
private static Valor invocar(String id, List<Argumento> args, Entorno caller, int l, int c) {
    DeclaracionFuncion f = ctx.funciones.get(id.toLowerCase());
    if (f == null) caller.errorSemantico("no existe la funcion o metodo '" + id + "'", l, c);

    // 1) detectar argumentos repetidos o desconocidos (validacion antes de ejecutar nada)
    Map<String, Argumento> mapaArgs = new LinkedHashMap<>();
    for (Argumento a : args) { ...
        if (mapaArgs.containsKey(a.id.toLowerCase()))
            caller.errorSemantico("argumento '" + a.id + "' repetido en la llamada a '" + id + "'", ...);
        mapaArgs.put(a.id.toLowerCase(), a);
    }

    // 2) nuevo Entorno COLGADO DEL GLOBAL (Tema 5.2), no del caller
    Entorno local = ctx.global.crearHijo(f.id);
    for (Parametro p : f.params) {
        Argumento a = mapaArgs.get(p.id.toLowerCase());
        Valor v = (a != null) ? a.expr.evaluar(caller)          // <- evaluado en el CALLER
                : (p.defecto != null) ? p.defecto.evaluar(caller)
                : /* error: falta el argumento */;
        local.declarar(new Simbolo(p.id, ..., v, ...), l, c);   // <- declarado en el LOCAL
    }

    Valor retorno = null;
    try { A.ejecutar(f.cuerpo, local); }
    catch (Senales.Retorno r) { retorno = r.valor; }
    ...
    return retorno;
}
```

El detalle sutil que vale la pena señalar dos veces: **el argumento se evalúa en el entorno del llamador (`caller`)**, porque ahí es donde viven las variables que el llamador pasó como argumento; **pero se declara en el entorno nuevo de la función (`local`)**, porque ahí es donde el cuerpo de la función va a buscarlo por nombre. Los dos entornos coexisten solo durante esa resolución puntual — ni antes ni después.

`return` viaja como `Senales.Retorno`, con el valor adentro de la excepción (mismo mecanismo del Tema 6.2, no algo nuevo): si el cuerpo de la función termina sin ejecutar ningún `return`, la excepción nunca se lanza, `retorno` queda `null`, y eso se traduce en error semántico si la función NO es `void` ("la función '...' debe retornar un ...").

### 7.4 Recursión: nada especial en el código, todo en la pila de llamadas

```
int factorial(n: int) {
    if (n <= 1) { return 1; }
    return n * factorial(n = n - 1);
}
```

No hay ningún caso especial para la recursión en `invocar(...)`: cada llamada a `factorial` simplemente vuelve a entrar en el mismo método Java, que vuelve a crear un `Entorno local` hijo del global, con su propio `n`. La pila de llamadas de Java (la de verdad, la del sistema operativo) es la que soporta la recursión de CompScript — por eso `Interprete.ejecutar(...)` atrapa `StackOverflowError` y lo traduce en un error semántico ordenado ("desbordamiento de pila, posible recursión infinita") en vez de dejar que la JVM entera se caiga.

### Ahora programá vos
Con `ejemplo5.cs` ya corriendo, cambiá `potencia(base = 4)` por una llamada que **omita** un parámetro sin valor por defecto (por ejemplo, agregá un tercer parámetro `unidad: string` sin default a `potencia` y llamala sin pasarlo). Confirmá que el error que aparece es exactamente `"falta el argumento 'unidad' en la llamada a 'potencia'"`, y que ese error aborta la ejecución antes de imprimir nada del cuerpo de la función.

### Dónde se rompe si...
- Evaluás el argumento en `local` en vez de `caller` → una llamada como `f(x = x + 1)` no podría resolver el `x` de la derecha (todavía no existe en `local`).
- Te olvidás de la primera pasada → una función que se llama a sí misma ANTES de su propia declaración textual (imposible en la práctica) o que llama a una función declarada después falla con "no existe la función".

---

## 8. Editor JavaFX y los 5 reportes

### 8.1 La estructura del editor, por código (sin FXML)

```java
public class EditorApp extends Application {
    private TabPane pestanas;
    private TextArea consola;
    private Contexto ultimo;   // resultado de la ULTIMA ejecucion, para los reportes

    public void start(Stage stage) {
        HBox barra = new HBox(8, bNuevo, bAbrir, bGuardar, new Separator(),
                bEjecutar, new Separator(), bReportes, bAst);
        SplitPane centro = new SplitPane(pestanas, consola);
        centro.setOrientation(Orientation.VERTICAL);
        BorderPane raiz = new BorderPane(centro);
        raiz.setTop(barra);
        ...
    }
}
```

Seis botones: **Nuevo, Abrir, Guardar, ▶ Ejecutar, Reportes, Ver AST**. `pestanas` es un `TabPane` (múltiples archivos `.cs` abiertos a la vez); `consola` es un `TextArea` de solo lectura. El campo `ultimo` es la clave de todo el capítulo de reportes: guarda el `Contexto` de la última vez que se presionó ▶ Ejecutar, y **nada más** que eso — si no ejecutaste nada, `ultimo == null` y los botones de reportes lo avisan en la consola en vez de fallar.

### 8.2 `Lanzador`: el mismo truco de siempre de JavaFX

```java
public class Lanzador {
    public static void main(String[] args) { EditorApp.main(args); }
}
```

Si la clase con el `main` **extiende** `Application` directamente, la JVM aplica una verificación de módulos que falla con "JavaFX runtime components are missing" cuando JavaFX no está en el *module path* (típicamente al correr desde el IDE en vez de `mvn javafx:run`). `Lanzador` no extiende nada — solo delega a `EditorApp.main(args)`, que sí llama a `launch(args)` por dentro — y esa indirección esquiva el chequeo. Por eso en IDEA corrés ▶ Run sobre `Lanzador`, **nunca sobre `EditorApp`**.

### 8.3 Entorno fresco por ejecución

```java
private void ejecutar() {
    TextArea area = areaActual();
    ultimo = Interprete.ejecutar(area.getText());   // Contexto FRESCO, no reutiliza el anterior
    ...
}
```

Cada clic en ▶ Ejecutar crea un `Contexto` **nuevo** desde cero (no reutiliza variables, funciones ni errores de la corrida anterior). Esto es lo que garantiza que los reportes de la sección 8.4 correspondan siempre al ÚLTIMO archivo ejecutado, tal como exige el enunciado — nunca a una mezcla de ejecuciones.

### 8.4 Los 5 reportes, y por qué el AST aparece en 3 formas

```java
public static File[] generar(Contexto ctx, File carpeta) throws Exception {
    Files.writeString(new File(carpeta, "tokens.html").toPath(), tokens(ctx));
    Files.writeString(new File(carpeta, "errores.html").toPath(), errores(ctx));
    Files.writeString(new File(carpeta, "simbolos.html").toPath(), simbolos(ctx));
    Files.writeString(new File(carpeta, "ast.html").toPath(), astHtml(ctx));
    Files.writeString(new File(carpeta, "ast.dot").toPath(), astDot(ctx));
    return new File[]{ t, e, s, a };   // el boton "Reportes" abre estos 4; ast.dot queda en disco
}
```

Cinco archivos generados, cuatro se abren automáticamente en el navegador al presionar **Reportes** (`tokens.html`, `errores.html`, `simbolos.html`, `ast.html`); el quinto (`ast.dot`) es un extra para quien quiera renderizarlo con Graphviz (`dot -Tpng ast.dot -o ast.png`) fuera de la aplicación. El botón **Ver AST** muestra una tercera forma más, embebida: un `TreeView` de JavaFX construido recorriendo `hijosAst()` de forma recursiva —

```java
private TreeItem<String> construir(A.Nodo nodo) {
    TreeItem<String> item = new TreeItem<>(nodo.etiquetaAst());
    for (A.Nodo h : nodo.hijosAst()) if (h != null) item.getChildren().add(construir(h));
    return item;
}
```

Las tres vistas (`ast.html`, `ast.dot`, `TreeView`) reutilizan exactamente el mismo `etiquetaAst()`/`hijosAst()` que ya escribiste en el Tema 4 — no triplican lógica, solo cambian el formato de salida. El nombre de cada tipo de token en `tokens.html` sale por **reflexión** sobre los campos públicos `int` de la clase generada `sym`, así que agregar un terminal nuevo al `.cup` nunca requiere tocar `Reportes.java`.

### Ahora programá vos
Corré `Lanzador`, escribí un `.cs` con al menos un error de cada tipo (léxico, sintáctico, semántico — como `ejemplo6_errores.cs`), presioná ▶ Ejecutar y después Reportes. Confirmá que `errores.html` lista los tres, en el mismo orden que aparecen en la consola de la GUI.

### Dónde se rompe si...
- Corrés ▶ Run sobre `EditorApp` en vez de `Lanzador` → "JavaFX runtime components are missing" (8.2).
- Presionás Reportes o Ver AST sin haber ejecutado nada antes → si no chequeás `ultimo == null` primero, te vas a un `NullPointerException` en vez de un mensaje amigable en consola.

---

## 9. Errores comunes reales al programar cada tema

Esta sección junta los bugs reales encontrados en la auditoría de código del 2026-07-21 (ambos ya corregidos en el repositorio) con los errores previsibles por el diseño mismo del proyecto. Todos comparten una lección de fondo: **en un AST donde cada nodo valida su propia semántica, no hay una gramática externa que te cubra las espaldas** — si un nodo no comprueba algo, nadie más lo va a comprobar por él.

### 9.1 El `ClassCastException` de `validarVector` — validá la CATEGORÍA antes de castear

La expresión inicializadora de un vector (`let v: int[] = <expresion>;`) es una `<expresion>` general de la gramática, no exclusivamente un literal `[ ... ]`. Nada en el `.cup` impide escribir `let v: int[] = 5;` — una expresión aritméticamente válida (`5` es un `Entero` perfectamente correcto) pero semánticamente incompatible con el tipo declarado.

```java
// ANTES (bug real): asume que crudo YA es un vector, sin comprobar
private Valor validarVector(Valor crudo, Tipo declarado, Entorno e) {
    List<Valor> vals = crudo.lista();   // Valor.lista() hace (List<Valor>) valor por dentro
    ...                                  // con crudo = Entero(5), esto lanza ClassCastException
}
```

Con `crudo` siendo un `Entero`, `crudo.lista()` intenta castear un `Integer` a `List<Valor>` — y el `ClassCastException` resultante **no estaba capturado en ningún punto del pipeline**. El resultado no era "el programa `.cs` reporta un error": era la JVM entera terminando el proceso, tirando abajo también los programas `.cs` que no tenían nada que ver. Eso contradice directamente lo que pide el enunciado (4.3): un error semántico debe **reportarse y terminar ordenadamente**, nunca reventar el intérprete.

```java
// DESPUES (corregido): guarda de tipo explicita, ANTES de castear
private Valor validarVector(Valor crudo, Tipo declarado, Entorno e) {
    if (crudo == null || crudo.tipo.cat != Tipo.Cat.VECTOR)
        e.errorSemantico("el vector '" + id + "' espera un literal de vector y recibio "
                + Entorno.tipoDe(crudo), linea, columna);
    List<Valor> vals = crudo.lista();   // para cuando llega aca, ya se garantizo que es VECTOR
    ...
}
```

**La regla general para todo lo que programes de acá en adelante:** cada vez que tu código vaya a llamar a `.lista()` o `.campos()` sobre un `Valor` que viene de evaluar una expresión arbitraria (no un literal que vos mismo construiste), comprobá primero `valor.tipo.cat` contra la categoría esperada. Es exactamente lo que ya hacen `AccesoVector` (`if (s.tipo.cat != Tipo.Cat.VECTOR) e.errorSemantico(...)`), `AccesoCampo` y `OperacionLista` — todos validan la categoría ANTES de castear, nunca "dan por sentado" el tipo de lo que reciben.

### 9.2 Otros errores previsibles, por tema

- **Tema 3 (sintáctico):** olvidar `%prec UMENOS` → conflicto shift-reduce silencioso hasta que corrés `mvn compile` y leés la consola con atención.
- **Tema 4 (AST):** castear `(Expresion) o` sobre un `Object` que en realidad es una `List` (por ejemplo, confundir un argumento simple con una lista de argumentos) → `ClassCastException` en el constructor del nodo, no en el `.cup` — el error aparece "más lejos" de donde realmente está la causa.
- **Tema 5 (ámbitos):** colgar el entorno de una función del `caller` en vez del global → *dynamic scoping* no deseado; el programa "funciona" en casos simples y falla de forma confusa apenas hay dos funciones con un parámetro de mismo nombre.
- **Tema 6 (flujo de control):** evaluar ambos lados de `&&`/`||` antes de decidir el operador (el bug de cortocircuito, ver 6.3) → guardas típicas (`x != 0 && ...`) dejan de proteger nada.
- **Tema 7 (funciones):** evaluar los argumentos en el entorno `local` en vez del `caller` → una llamada que reutiliza el nombre de una variable externa como argumento falla en resolver esa variable.
- **Tema 8 (GUI):** no comprobar `ultimo == null` antes de generar reportes o mostrar el AST → `NullPointerException` en vez de un mensaje de consola razonable.

### Ahora programá vos, para cerrar
Elegí una de las clases de nodo que todavía no tiene guarda de tipo explícita antes de un cast interno (repasá `AccesoCampo`, `OperacionLista` en `ast/A.java`) y confirmá, leyendo el código, que efectivamente valida `s.tipo.cat` antes de llamar a `.campos()` o a los métodos de lista. Si algún día agregás una construcción nueva al lenguaje (una sentencia propia, por ejemplo), aplicá la misma disciplina: **valida la categoría antes de castear, siempre** — es la lección más barata de aplicar y la más cara de omitir, como ya demostró el bug de `validarVector`.
