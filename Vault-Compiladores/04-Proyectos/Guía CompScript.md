---
tags: [proyecto, guia, lexico, sintactico, semantico, ast]
aliases: [guia compscript, tutorial compscript, compscript paso a paso]
fuente: "Enunciado CompScript (OLC1_PT1_ VD2024.clean.md) + Libro del Dragón + construcción real verificada 2026-07-20"
fecha: 2026-07-20
---

# Guía de elaboración — [[CompScript]]

A diferencia de [[DataForge]] (que lo construyó Claude como demostración guiada), **CompScript lo escribiste vos**: acá el rol de Claude fue tutorear y revisar. Esta guía documenta el material REAL una vez que el código ya compila y corre (`mvn compile` limpio, `TestInterprete` verificado sobre los 6 ejemplos de `entradas/`).

## Estado de construcción

| Etapa | Contenido | Estado |
|---|---|---|
| Léxico | 37 palabras reservadas + 28 símbolos, `Lexer.flex` | ✅ verificado 2026-07-20 |
| Sintáctico | `parser.cup`, construye AST (NO ejecuta en las acciones) | ✅ verificado 2026-07-20 |
| AST | `ast/A.java` — patrón tree-walk, interfaces `Nodo/Instruccion/Expresion` | ✅ verificado 2026-07-20 |
| Semántico | `interprete/*` — pila de entornos, tablas de compatibilidad de tipos | ✅ verificado 2026-07-20 |
| Ejecución | `Interprete.java` — 3 pasadas (funciones/structs → globales → RUN_MAIN) | ✅ verificado 2026-07-20 |
| GUI | `gui/EditorApp.java` — pestañas, ejecutar, reportes, AST en TreeView | ✅ (código completo; no se abrió la ventana en esta sesión) |
| Reportes | `reportes/Reportes.java` — tokens/errores/símbolos en HTML + AST en HTML y `.dot` | ✅ (código completo) |
| Fase 1 (enunciado §8) | tipos, expresiones, control, ciclos, transferencia sin `return`, impresión, reportes tokens/errores/símbolos | ✅ cubierto por el código actual |
| Fase 2 (enunciado §8) | vectores, listas, structs, funciones, métodos, llamadas, `RUN_MAIN`, resto de nativas | ✅ cubierto por el código actual |

**Ambas fases del enunciado ya están implementadas en el código actual** (Fase 1 y Fase 2 no quedaron separadas en el repo: el proyecto se entrega completo). Pendiente real: `docs/gramatica.txt` recién se generó (ver [[CompScript]]), y falta el empaquetado de entrega (`.jar`, manuales en PDF, repo `OLC1_VD24_#Carnet`).

## 1. Requisitos previos (reales)

- **IntelliJ IDEA**: importa `pom.xml` como Maven; JDK 17 (`maven.compiler.source/target`), aunque esta verificación se corrió con el JDK 25 embebido de IDEA sin problema (el `pom.xml` no fija una toolchain estricta).
- `mvn compile` regenera `Lexer.java`/`Parser.java`/`sym.java` en `target/generated-sources/` vía `jflex-maven-plugin:1.9.1` y `cup-maven-plugin:11b-20160615-3` — igual que en DataForge.
- GUI: `mvn clean javafx:run`, o ▶ Run sobre `gui.Lanzador` en IDEA (mismo truco que DataForge: `Lanzador` no extiende `Application`, esquiva el error "JavaFX runtime components are missing").
- Prueba por consola: `mvn compile exec:java -Dexec.mainClass=compscript.TestInterprete -Dexec.args="entradas/ejemplo1.cs"`.

## 2. Análisis léxico — CONSTRUIDO ✅

### 2.1 Tabla de tokens (la real, `Lexer.flex`)

- **37 palabras reservadas** (case insensitive vía `%ignorecase`): tipos `int double bool char string void`; declaración `let const true false`; control `if else match default`; ciclos `while for do`; transferencia `break continue return`; compuestos `struct list cast as`; consola/listas `console log push get set remove pop reverse`; nativas `round length tostring run_main`.
- **28 símbolos**: `++ -- == != <= >= && || =>` · `+ - * / ^ $ % ! < > =` · `( ) [ ] { } , ; : .`
- **Con patrón**: `ID` = `[a-zA-Z_]([a-zA-Z_]|[0-9])*` · `ENTERO` = dígitos · `DECIMAL` = dígitos"."dígitos · `CADENA` = `"([^"\\\n]|\\.)*"` · `CARACTER` = `'([^'\\\n]|\\.)'`.
- **Se descartan**: `// ...` (línea), `/* ... */` (multilínea), blancos.
- **Secuencias de escape** (5.4) resueltas en el propio lexer, no en el parser: `procesarEscapes()` traduce `\n \t \\ \" \'` dentro de `CADENA`/`CARACTER` antes de devolver el símbolo.

### 2.2 Decisiones de diseño

1. **Orden crítico**: las 37 reservadas están declaradas ANTES que `{Id}` en el `.flex` — si no, `"int"` matchearía como `ID` (mismo motivo que en [[DataForge]], [[Cap 3 - Análisis léxico]]).
2. El lexer no imprime a `stderr`: si `contexto` está conectado (lo setea `Interprete`), cada error léxico va a `contexto.error("Lexico", ...)` — la tabla de errores (6.2) reúne los 3 tipos en un solo lugar.
3. Cada token reconocido se registra con `contexto.registrarToken(yytext(), type, yyline, yycolumn)` — esto alimenta el reporte de tokens (6.1) sin que el lexer sepa nada de HTML ni de la GUI.

### 2.3 Verificación real

`ejemplo1.cs` (tipos + expresiones): **209 tokens reconocidos**, 0 errores. `ejemplo6_errores.cs` (con un `#` suelto): el carácter se descarta y el análisis sigue — el error queda en la tabla como `[Lexico] el caracter '#' no pertenece al lenguaje (linea 5, columna 1)`.

## 3. Análisis sintáctico — CONSTRUIDO ✅

### 3.1 La gramática

- **BNF limpia** (entregable, NO copia del `.cup`): `docs/gramatica.txt`. A diferencia de [[DataForge]] (llamadas tipo `SUM(a,b)`, sin ambigüedad) y de [[ConjAnalyzer]] (notación prefija), CompScript tiene **expresiones infijas** (`a + b * c`), así que la BNF de `<expresion>` es ambigua por diseño y se desambigua con una **tabla de precedencia y asociatividad** ([[Ambigüedad, precedencia y asociatividad]]) — igual que hace el propio enunciado en su §5.9.
- Tabla real (`parser.cup`, de más débil a más fuerte, el ÚLTIMO declarado liga más fuerte en CUP):

```
precedence left  OR;                                              // nivel 7 (mas debil)
precedence left  AND;                                             // nivel 6
precedence right NOT;                                             // nivel 5
precedence left  IGUAL_IGUAL, DIFERENTE, MENOR, MENOR_IGUAL, MAYOR, MAYOR_IGUAL;  // nivel 4
precedence left  MAS, MENOS;                                       // nivel 3
precedence left  POR, DIV, MOD;                                    // nivel 2
precedence nonassoc POT, RAIZ;                                     // nivel 1
precedence right UMENOS;                                           // nivel 0 (mas fuerte)
```

- El **else colgante** ([[Conflictos shift-reduce y reduce-reduce]]) se resuelve con la gramática misma: `if_stmt ::= IF (...) bloque ELSE if_stmt` reduce el `else if` como parte del `if` más cercano — CUP no reportó conflictos shift-reduce al generar.
- `UMENOS` es un **pseudo-terminal** (`terminal UMENOS;` sin lexema): solo existe para que `%prec UMENOS` le dé a la negación unaria (`MENOS expr`) una precedencia distinta de la resta binaria, aunque ambas usan el mismo token `MENOS`. Truco estándar de CUP/Yacc para el problema clásico "unario vs. binario" del mismo símbolo.

### 3.2 El `.cup` real — construye AST, no ejecuta

Esta es la decisión de diseño que distingue a CompScript de [[DataForge]] y [[ConjAnalyzer]] (ambos S-atribuidos, ejecutan directo en las acciones `{: RESULT = ... :}`). Acá cada acción arma un nodo:

```java
expr ::= expr:a MAS expr:b  {: RESULT = new A.Binaria("+", a, b, aleft, aright); :}
       | ...
       | ID:id              {: RESULT = new A.AccesoVariable(id, idleft, idright); :}
       ;
```

`parser.raiz` termina siendo un `ArrayList` de nodos `A.Instruccion`; el `Interprete` los recorre después. Por qué: el enunciado pide explícitamente "construir un AST" (objetivo específico #3) y CompScript SÍ tiene control de flujo (`if/while/for/match/funciones`) — evaluar en las acciones del parser, con reducciones LALR de abajo hacia arriba, no se lleva bien con "ejecutar un `while` cero o N veces" (necesitás poder re-evaluar el mismo subárbol). [[Traducción dirigida por la sintaxis]] documenta esta tensión: S-atribuido = una pasada, un valor; AST = tantas pasadas como haga falta.

### 3.3 Verificación real

Los 6 ejemplos de `entradas/` parsean y ejecutan. Caso negativo (`ejemplo6_errores.cs`, línea 8: `let w: int = ;`): `syntax_error()` registra `[Sintactico] error de sintaxis, no se esperaba ';' (linea 8, columna 14)` y el modo pánico (`instruccion ::= error PUNTO_COMA`) descarta hasta el siguiente `;`, permitiendo que el análisis siga con la línea 11.

## 4. El AST — `ast/A.java` (CONSTRUIDO ✅, la pieza que no tienen los proyectos hermanos)

### 4.1 Diseño: una sola clase contenedora, dos interfaces

```java
public interface Nodo {
    String etiquetaAst();
    default List<Nodo> hijosAst() { return new ArrayList<>(); }
}
public interface Instruccion extends Nodo { void ejecutar(Entorno e); }
public interface Expresion  extends Nodo { Valor evaluar(Entorno e); }
```

Cada nodo de sintaxis (`A.If`, `A.While`, `A.Binaria`, `A.Llamada`, ~25 clases en total) implementa `Instruccion` o `Expresion` y sabe **auto-ejecutarse o auto-evaluarse** — el patrón clásico "tree-walking interpreter" (variante del Visitor sin doble despacho: en vez de un visitante externo, el propio nodo tiene el método). Ver [[Árbol de sintaxis abstracta (AST)]].

`etiquetaAst()` + `hijosAst()` son el segundo propósito de cada nodo: alimentan el **reporte de AST** (6.3) sin necesitar una clase visitante aparte — `Reportes.astHtmlNodo()` y `EditorApp.construir()` simplemente recorren `hijosAst()` recursivamente.

### 4.2 Por qué los constructores reciben `Object`

```java
public Binaria(String op, Object izq, Object der, int l, int c) {
    this.op = op; this.izq = (Expresion) izq; this.der = (Expresion) der; ...
}
```

Las acciones del `.cup` no pueden usar genéricos con seguridad (los `non terminal` son tipos raw, igual que en DataForge/ConjAnalyzer) — recibir `Object` y castear adentro del constructor mantiene las acciones del `.cup` en una sola línea, sin casts repetidos ahí.

### 4.3 Verificación real (reporte 6.3)

`Reportes.java` genera el AST de dos formas: `ast.html` (árbol anidado en `<ul>/<li>`, autocontenido, sin dependencias) y `ast.dot` (fuente [[Graphviz]] para quien quiera `dot -Tpng ast.dot -o ast.png`, alta fidelidad). La GUI además lo muestra en un `TreeView` de JavaFX (`EditorApp.verAst()`) — o sea 3 vistas del mismo árbol para cumplir "mostrarlo desde la interfaz" (§4.5) sin triplicar lógica.

## 5. Análisis semántico y ejecución — CONSTRUIDO ✅

### 5.1 Entornos y alcance: pila enlazada, NO un solo mapa global

```java
public class Entorno {
    public final Contexto contexto;
    public final Entorno padre;       // null = entorno global
    private final LinkedHashMap<String, Simbolo> tabla = new LinkedHashMap<>();
    public Entorno crearHijo(String nombre) { return new Entorno(contexto, this, nombre); }
    public Simbolo buscar(String id) {
        for (Entorno e = this; e != null; e = e.padre) { ... }  // sube por los padres
    }
}
```

Cada `if`, `while`, `for`, `match`, cuerpo de función… crea un `Entorno` hijo (`e.crearHijo("if")`), y `buscar()` sube por la cadena de padres hasta encontrar la variable o llegar al global. Esto es exactamente [[Entornos y alcance|alcance estático anidado]] y [[Registro de activación y pila de control]]: cada llamada a función (`ctx.global.crearHijo(f.id)`) es un nuevo registro de activación cuyo padre es SIEMPRE el entorno global (no el entorno del llamador) — así CompScript implementa **alcance estático** (una función no ve las variables locales de quien la llama), verificado con la recursión de `fibonacci` (ver §7).

### 5.2 `Contexto` vs `Entorno`: quién es compartido y quién no

`Entorno` es UN ámbito de la pila de scopes (privado a su rama del árbol de llamadas). `Contexto` es el estado global único de la corrida completa: consola (`StringBuilder`), lista de errores, lista de TODOS los tokens, lista de TODOS los símbolos alguna vez declarados (`ctx.simbolos`, plana, para el reporte 6.4), y los mapas `funciones`/`structs` registrados en la primera pasada. Mismo patrón que [[DataForge]] pero con un nivel más (ahí no había pila de entornos porque no hay funciones ni bloques anidados).

### 5.3 Comprobación de tipos: `Tipo` + las tablas de `Operaciones`

`Tipo` (`Cat` enum: `INT DOUBLE BOOL CHAR STRING VOID NULL VECTOR LIST STRUCT`) implementa **igualdad estructural** (`equals()`): dos vectores son el mismo tipo si tienen igual `elemento` y `dimensiones`; dos structs, si tienen el mismo `structName`. Esto es [[Comprobación de tipos]] aplicada a tipos compuestos, no solo primitivos.

`Operaciones.java` implementa LITERALMENTE las 8 tablas de compatibilidad del enunciado (§5.5-5.7: suma, resta, multiplicación, división, potencia, raíz, módulo, relacionales) más los 5 casteos válidos de §5.13 ([[Conversión de tipos (coerción y cast)]]). Ejemplo real, la tabla de suma es la única que admite `bool`/`char`/`string` mezclados porque el enunciado así lo pide:

```java
private static Valor suma(Valor a, Valor b, Entorno e, int l, int c) {
    Cat x = a.tipo.cat, y = b.tipo.cat;
    if (x == Cat.STRING || y == Cat.STRING) { if (esPrimitivo(x) && esPrimitivo(y)) res = Cat.STRING; }
    else if (x == Cat.CHAR && y == Cat.CHAR) { res = Cat.STRING; }        // char + char -> cadena
    else if (x == Cat.INT && y == Cat.INT) { res = Cat.INT; }
    else if (/* numerico/bool/char */) { ... }
    ...
}
```

Toda combinación no contemplada cae en `error(...)`, que llama a `e.errorSemantico(...)`.

### 5.4 La decisión que MÁS distingue a CompScript de DataForge: el error semántico ABORTA

```java
// Entorno.java
public void errorSemantico(String descripcion, int linea, int columna) {
    contexto.error("Semantico", descripcion, linea, columna);
    throw new ErrorSemantico(descripcion, linea, columna);   // desenrolla la pila
}
```

En DataForge, una expresión con error semántico se convierte en `null` y las operaciones que lo reciben simplemente callan (propagación por null, sin abortar). En CompScript el enunciado (§4.3) pide lo contrario: *"el intérprete debe tener la capacidad de recuperarse de errores... y poder TERMINAR la ejecución del programa"*. Por eso `ErrorSemantico extends RuntimeException` (sin stack trace, es control de flujo) y `Interprete.ejecutar()` la atrapa en el nivel más alto: el error ya quedó en la tabla (6.2), y la ejecución corta ordenadamente ahí — no sigue con las instrucciones siguientes. Verificado con `ejemplo6_errores.cs`: la línea `"a" - "b"` lanza el error semántico y el `console.log(z)` de la línea siguiente JAMÁS se ejecuta (no aparece en consola).

### 5.5 `break`/`continue`/`return` como excepciones de control (`Senales.java`)

```java
public static class Break extends RuntimeException {
    public Break(int l, int c) { super(null, null, false, false); ... }  // sin stacktrace, mensaje ni supresion
}
```

`while`/`for`/`do-while` atrapan `Senales.Break`/`Continue` alrededor de cada ejecución del cuerpo; las llamadas a función atrapan `Senales.Retorno` para extraer el valor de un `return`. Si una de estas señales se escapa hasta `Interprete.ejecutar()` (p. ej. un `break;` fuera de cualquier ciclo), se convierte en un error semántico ("'break' fuera de un ciclo"). Ver [[Flujo de control y switch]].

### 5.6 Ejecución en 3 pasadas (`Interprete.java`, recomendación del enunciado §5.23)

```java
// 1a pasada: registro de funciones/metodos/structs (permite forward-reference)
for (A.Instruccion i : raiz) {
    if (i instanceof A.DeclaracionFuncion f) ctx.registrarFuncion(f);
    else if (i instanceof A.DeclaracionStruct s) ctx.registrarStruct(s);
}
// 2a pasada: declaraciones y asignaciones globales
// 3a pasada: RUN_MAIN (punto de entrada)
```

Esto permite llamar una función antes de que aparezca textualmente en el archivo — verificado en `ejemplo2.cs`, donde `clasificar` se define antes de `main`, pero el diseño soporta el orden inverso porque el registro ocurre en una pasada previa a cualquier ejecución.

### 5.7 `match` sin fall-through

```java
for (CasoMatch cm : casos) {
    ...
    if (eq != null && (Boolean) eq.valor) { A.ejecutar(cm.cuerpo, e.crearHijo("match")); return; }
}
```

El `return` explícito dentro del bucle de casos es la implementación literal de "no requiere `break`, no hay fall-through" del enunciado §5.15.2 — apenas un caso matchea, se ejecuta su cuerpo y se sale del método sin mirar los casos restantes.

## 6. GUI y reportes — CONSTRUIDOS (código completo, no verificado visualmente en esta sesión)

- `gui/EditorApp.java`: JavaFX **por código** (sin FXML, igual que DataForge). `TabPane` con pestañas (`Tab.userData = File`, igual patrón que DataForge para "Guardar" vs "Guardar como"), `TextArea` de consola no editable, barra con **Nuevo / Abrir / Guardar / ▶ Ejecutar / Reportes / Ver AST**. `gui/Lanzador.java` esquiva el mismo problema de siempre ("JavaFX runtime components are missing") con una clase `main` que no extiende `Application`.
- **Entorno fresco por ejecución**: `ultimo = Interprete.ejecutar(area.getText())` crea un `Contexto` nuevo cada vez — los reportes siempre reflejan el ÚLTIMO análisis (cumple la exigencia general de "reportes solo del último archivo ejecutado" que también aplica en DataForge).
- `reportes/Reportes.java` genera 4 HTML autocontenidos (`tokens.html`, `errores.html`, `simbolos.html`, `ast.html`) + `ast.dot`, con **nombres de token por reflexión sobre `sym.class`** (mismo truco que DataForge: `sym` la genera CUP a partir de las declaraciones `terminal`, así que nunca hay que mantener una tabla de nombres a mano).
- El AST también se ve embebido en la GUI vía `TreeView` (`EditorApp.verAst()`), sin depender de abrir el HTML — cumple "mostrarlo desde la interfaz" (§4.5) de una segunda forma.

## 7. Casos de prueba (reales, corridos en esta sesión — `mvn compile exec:java -Dexec.mainClass=compscript.TestInterprete -Dexec.args="entradas/ejemploN.cs"`)

| Archivo | Qué prueba | Resultado real |
|---|---|---|
| `ejemplo1.cs` | tipos, expresiones aritméticas/relacionales/lógicas, casteos, `console.log` | 209 tokens, 0 errores; `a/b` (10/3) → `3.3333333333333335` (división siempre Decimal); `a%b` → `1.0`; `a^b` → `1000` (potencia entero-entero da entero); `9$2` → `3.0` (raíz cuadrada); `'a'+'b'` → `"ab"` (char+char→cadena) |
| `ejemplo2.cs` | `match` con y sin `default`, funciones con retorno `string`, llamada con argumento por nombre | consola: `uno / tres / otro numero / cayo en default`; tabla de símbolos muestra `n` y `r` re-declarados en CADA invocación de `clasificar` (ámbito distinto por llamada) |
| `ejemplo3.cs` | `while`, `for` con `continue`/`break`, `do-while` | consola: `0 1 2 / 0 10 30 / 0 1` — confirma que `continue` en `i==2` salta el `console.log` y `break` en `i==4` corta el `for` antes de imprimir 40 |
| `ejemplo4.cs` | vectores 1D/2D, listas dinámicas (`push/get/set/remove/reverse`), structs (instanciación por nombre, no por posición) | consola coincide con los comentarios del archivo; símbolo `p` de tipo `Struct persona` reporta `persona { nombre: "Ana", edad: 31 }` (formato `toString`, 5.27) |
| `ejemplo5.cs` | recursión (`factorial`, `fibonacci`), parámetro con valor por defecto | `factorial(5)`→120, `fib(10)`→55, `potencia(base=4)` (usa `exp` por defecto=2)→16, `potencia(base=2,exp=5)`→32. La tabla de símbolos muestra **186 entradas** para `n` en distintos ámbitos `fib` — evidencia directa de que cada llamada recursiva crea su propio `Entorno` y el reporte 6.4 es un log histórico plano, no solo el estado final |
| `ejemplo6_errores.cs` | los 3 tipos de error en un solo archivo | `[Lexico]` carácter `#` descartado, línea 5; `[Sintactico]` `;` inesperado, línea 8 (modo pánico); `[Semantico]` resta entre dos cadenas, línea 11 — y la ejecución **se detiene ahí**: el `console.log(z)` de la línea 12 nunca corre (ver §5.4) |

## 8. Errores comunes (los que ya aparecieron o son previsibles por el diseño)

- **Confundir la propagación de errores con la de DataForge**: acá un error semántico ABORTA (lanza `ErrorSemantico`), no propaga `null` en silencio. Si el reporte de errores muestra solo 1 error semántico aunque el programa tenga varios problemas de tipos, es esperado: el intérprete se detiene en el primero.
- **`UMENOS` sin `%prec`**: si se te olvida `%prec UMENOS` en la producción de negación unaria, CUP no sabe que `MENOS expr` (unario) debe ligar distinto de `expr MENOS expr` (resta) — conflicto shift-reduce en la gramática.
- **Alcance de funciones**: `local = ctx.global.crearHijo(f.id)` cuelga SIEMPRE del entorno global, nunca del entorno del llamador — si por error se cuelga del `caller`, las funciones dejarían de tener alcance estático y verían variables locales de quien las invoca (bug clásico de "dynamic scoping" no deseado).
- **Reflexión sobre `sym.class`** para nombres de token en los reportes: si el `.cup` no declara el terminal con `terminal X;` (por ejemplo un pseudo-terminal como `UMENOS`), igual aparece en `sym`, así que los reportes seguirán funcionando aunque `UMENOS` nunca aparezca como lexema real.
- **Non terminal tipados raw**: igual que en DataForge, los `non terminal ArrayList ...;` no usan genéricos — usar `List<X>` ahí rompe la generación de CUP.
- **`error PUNTO_COMA`** (modo pánico) sacrifica la instrucción completa hasta el próximo `;`; si el error ocurre dentro de un bloque sin `;` cercano (por ejemplo un `if` mal cerrado), el pánico puede consumir más de lo esperado — vale la pena probarlo con casos reales como se hizo acá.

## Relacionadas
- [[CompScript]]
- [[JFlex]] · [[CUP]] · [[Maven]] · [[JavaFX y Scene Builder]] · [[Graphviz]]
- [[Árbol de sintaxis abstracta (AST)]] · [[Traducción dirigida por la sintaxis]] · [[Atributos sintetizados y heredados]]
- [[Comprobación de tipos]] · [[Conversión de tipos (coerción y cast)]] · [[Flujo de control y switch]]
- [[Entornos y alcance]] · [[Tabla de símbolos]] · [[Registro de activación y pila de control]]
- [[Manejo de errores (léxicos, sintácticos, semánticos)]]
- [[Ambigüedad, precedencia y asociatividad]] · [[Conflictos shift-reduce y reduce-reduce]] · [[Análisis sintáctico ascendente LR]]
- [[Cap 4 - Análisis sintáctico]] · [[Cap 5 - Traducción dirigida por la sintaxis]] · [[Cap 7 - Entornos en tiempo de ejecución]]
- [[Guía DataForge]] (proyecto hermano, mismo stack, decisiones contrastantes)
