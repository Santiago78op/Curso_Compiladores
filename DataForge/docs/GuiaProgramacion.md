# Guía de Programación — DataForge

**Organización de Lenguajes y Compiladores 1 — USAC**
**Material de tutoría: cómo programar DataForge desde cero, por tema**

---

## Cómo usar esta guía

Esta guía **no** sigue el orden cronológico de una clase (eso ya lo tenés en `presentacion-dataforge/`, etapa por etapa). Acá el corte es **por tema técnico**: análisis léxico, análisis sintáctico, ejecución, GUI, reportes. Es más densa y asume que ya viste la teoría en las diapositivas o en el Dragón — el objetivo es que **vos programes** cada pieza, con el código real de DataForge como referencia de "a dónde tenés que llegar", no como algo para copiar sin entender.

Cada tema tiene la misma estructura:

1. **El problema** — qué te pide el enunciado, en tus palabras.
2. **La decisión de diseño** — por qué se resolvió así y no de otra forma (con alternativas descartadas).
3. **Código real** — el fragmento tal como quedó en `src/`, citado textual.
4. **Errores comunes** — lo que ya le pasó a alguien programando esto mismo.

Cuando la guía dice "ahora escribí...", el código real que sigue es la **referencia para comparar**, no el primer lugar donde deberías mirar. Intentalo primero solo.

---

## 1. Preparación del proyecto: Maven, JFlex y CUP

### 1.1 El problema

Vas a escribir DataForge en tres lenguajes de especificación distintos que conviven en un solo proyecto Java:

- `Lexer.flex` — reglas léxicas, las procesa **JFlex**.
- `parser.cup` — gramática + acciones, las procesa **CUP**.
- El resto — Java normal.

Ninguno de los dos primeros es Java: son **generadores de código**. JFlex lee `Lexer.flex` y escribe `Lexer.java`; CUP lee `parser.cup` y escribe `Parser.java` + `sym.java`. Vos nunca tocás esos `.java` generados — si necesitás cambiar algo, lo cambiás en el `.flex`/`.cup` fuente y volvés a generar.

### 1.2 La decisión de diseño: generar en `generate-sources`, no a mano

Maven tiene una fase específica para esto (`generate-sources`, antes de `compile`). Los plugins `jflex-maven-plugin` y `cup-maven-plugin` se enganchan ahí:

```xml
<plugin>
    <groupId>com.github.vbmacher</groupId>
    <artifactId>cup-maven-plugin</artifactId>
    <version>11b-20160615-3</version>
    <executions>
        <execution>
            <goals><goal>generate</goal></goals>
        </execution>
    </executions>
    <configuration>
        <className>Parser</className>
        <symbolsName>sym</symbolsName>
    </configuration>
</plugin>

<plugin>
    <groupId>de.jflex</groupId>
    <artifactId>jflex-maven-plugin</artifactId>
    <version>1.9.1</version>
    <executions>
        <execution>
            <goals><goal>generate</goal></goals>
        </execution>
    </executions>
</plugin>
```

**Por qué importa que ambos corran en la misma fase, antes de `compile`**: `Lexer.flex` referencia `sym.PROGRAM`, `sym.NUMERO`, etc. (las constantes que CUP genera a partir de tus declaraciones `terminal`). Si JFlex generara `Lexer.java` en una fase posterior a la compilación de `sym.java`, no habría problema tampoco — pero lo importante es que **los dos archivos generados (`Lexer.java` y `sym.java`) tienen que existir ANTES de que `javac` compile cualquier cosa**, porque se referencian cruzado. Poniendo ambos plugins en `generate-sources` (que corre siempre antes de `compile`), Maven garantiza el orden sin que vos tengas que pensarlo.

Convención de carpetas que ambos plugins asumen por defecto (no la inventes distinta sin razón):

```
src/main/jflex/Lexer.flex     → target/generated-sources/jflex/dataforge/analisis/Lexer.java
src/main/cup/parser.cup       → target/generated-sources/cup/dataforge/analisis/Parser.java
                                                          dataforge/analisis/sym.java
```

Las dos dependencias en tiempo de ejecución que necesitás:

```xml
<dependency>
    <groupId>com.github.vbmacher</groupId>
    <artifactId>java-cup-runtime</artifactId>
    <version>11b-20160615</version>
</dependency>
<dependency>
    <groupId>org.openjfx</groupId>
    <artifactId>javafx-controls</artifactId>
    <version>21.0.4</version>
</dependency>
```

`java-cup-runtime` aporta la clase `Symbol` que tu lexer va a construir y devolver en cada token — la necesitás desde la Etapa de léxico, mucho antes de escribir una sola línea de `.cup`. `javafx-controls` la agregás recién cuando llegués al tema de GUI (sección 5) — no hace falta desde el principio.

### 1.3 Ejercicio: arrancar el esqueleto

1. Creá `pom.xml` con lo de arriba (agregá también `exec-maven-plugin` con `mainClass` apuntando a una clase de prueba tuya — la vas a necesitar en la sección 2).
2. Creá `src/main/jflex/Lexer.flex` con el mínimo indispensable para que compile (aunque no reconozca nada útil todavía): las tres secciones separadas por `%%` (ver sección 2.2).
3. Creá `src/main/cup/parser.cup` igual de mínimo.
4. Corré `mvn compile`. Si ves errores de `sym` no encontrado desde `Lexer.flex`, revisá que **CUP haya corrido primero** — normalmente Maven resuelve el orden solo, pero si falla, un `mvn clean compile` fuerza a regenerar todo desde cero.

**Error común**: si IntelliJ subraya en rojo las clases `Parser` o `sym` aunque el proyecto compile bien por Maven, es porque el IDE no re-indexó el `target/generated-sources/`. Solución: Maven → **Reload All Maven Projects**, no hace falta tocar código.

---

## 2. Análisis léxico — JFlex

### 2.1 El problema: convertir el enunciado en una tabla de tokens

Antes de escribir una sola regla en `Lexer.flex`, tenés que responder una pregunta por cada sección del enunciado (§5.1 a §5.10): **¿qué palabra o símbolo nuevo aparece acá?** Es un trabajo de lectura, no de código. Por ejemplo:

- §5.2 (encapsulamiento) → aparecen `PROGRAM` y `END`.
- §5.5 (declarar variable) → aparecen `var`, `:`, `::`, `<-`, `end`, `;`.
- §5.6.2 (arreglos) → aparece el prefijo `@` y el símbolo `[` `]`.
- §5.10 (gráficas) → aparecen `graphBar`/`graphPie`/`graphLine`/`Histogram`/`EXEC`, y quince nombres de atributo (`titulo`, `ejeX`, `values`...) que **no** son reservados (ver 2.4).

El resultado de este barrido es la tabla de tokens real de DataForge:

| Categoría | Tokens |
|---|---|
| Encapsulamiento | `PROGRAM`, `END` |
| Declaración | `VAR`, `ARR`, `DOUBLE`, `CHAR` |
| Aritmética | `SUM`, `RES`, `MUL`, `DIV`, `MOD` |
| Estadística | `MEDIA`, `MEDIANA`, `MODA`, `VARIANZA`, `MAX`, `MIN` |
| Consola | `CONSOLE`, `PRINT`, `COLUMN` |
| Gráficas | `GRAPH_BAR`, `GRAPH_PIE`, `GRAPH_LINE`, `HISTOGRAM`, `EXEC` |
| Símbolos | `:` `::` `<-` `=` `->` `;` `,` `(` `)` `[` `]` |
| Con patrón | `NUMERO`, `CADENA`, `ID`, `ID_ARREGLO` |

25 tokens de palabra reservada (28 formas de superficie si contás los alias `graphbar`/`grapbar`, `graphpie`/`grappie`, `graphline`/`grapline` — ver 2.4), 11 símbolos, 4 con patrón léxico propio.

### 2.2 Las tres secciones de un `.flex`

```jflex
package dataforge.analisis;
import java_cup.runtime.Symbol;

%%                                    ← separador 1

%class Lexer
%public
%unicode
%cup
%line
%column
%ignorecase

%{
  // código Java auxiliar
%}

Letra   = [a-zA-Z]
Digito  = [0-9]
Id      = {Letra}({Letra}|{Digito})*

%%                                    ← separador 2 (reglas)

"program"   { return symbol(sym.PROGRAM); }
{Id}        { return symbol(sym.ID); }
[^]         { /* error léxico */ }
```

Cada directiva de la sección 2 (opciones) tiene un motivo concreto — **no las escribas de memoria, entendé qué rompe si falta cada una**:

- **`%public`**: sin ella, la clase `Lexer` generada queda *package-private*. Como vas a usarla desde `dataforge.gui` y `dataforge.interprete` (paquetes distintos), el compilador tira `'Lexer' is not public in dataforge.analisis; cannot be accessed from outside package`. Es el error más común de la Etapa 1 — si lo ves, ya sabés la causa.
- **`%cup`**: le dice a JFlex que genere un lexer compatible con la interfaz que CUP espera (`next_token()` devolviendo `Symbol`). Sin esto, JFlex genera un `Lexer` "genérico" que CUP no puede usar.
- **`%line` / `%column`**: activan `yyline`/`yycolumn`, que vas a necesitar para reportar la posición de cada token y de cada error (§6.1 y §6.2 del enunciado piden línea y columna en los reportes).
- **`%unicode`**: soporta tildes/ñ en cadenas y comentarios sin romper el conteo de caracteres.
- **`%ignorecase`**: resuelve §5.1 del enunciado (case insensitive) en **una sola directiva**, en la capa léxica — no hace falta escribir `"VaR"|"var"|"VAR"|...` a mano ni resolverlo en el parser.

### 2.3 Las definiciones regulares (macros)

```jflex
Letra        = [a-zA-Z]
Digito       = [0-9]
Id           = {Letra}({Letra}|{Digito})*
IdArreglo    = "@"{Id}
Numero       = {Digito}+("."{Digito}+)?
Cadena       = \"[^\"]*\"
ComentLinea  = "!"[^\r\n]*
ComentMulti  = "<!"~"!>"
Blancos      = [ \t\r\n]+
```

Dos decisiones de diseño concretas acá, ambas defendibles con alternativas descartadas:

1. **`IdArreglo` es un patrón completo (`@` + `Id`), no dos tokens separados.** La alternativa —tratar `@` como símbolo suelto y dejar que el parser una `@` seguido de `ID`— también era válida, pero complica la gramática sin necesidad (tendrías que escribir `arreglo ::= ARROBA ID` en vez de simplemente `ID_ARREGLO`). Se eligió la opción que deja la gramática de la sección 3 más simple.
2. **`ComentMulti` usa el operador `~` de JFlex (`"<!"~"!>"`)**, que significa "cualquier cosa hasta la primera ocurrencia de". Si en cambio hubieras escrito `"<!"[^]*"!>"` con matching goloso normal, el comentario se comería de más si hay dos bloques `<! ... !> ... <! ... !>` en el mismo archivo (el `*` sin restricción intentaría llegar hasta el **último** `!>` del archivo). El operador `~` evita ese problema sin que tengas que escribir una expresión regular más compleja a mano.

### 2.4 Las reglas: orden y prioridad

Regla de oro de JFlex (Dragón, cap. 3): **ante un empate en la longitud del lexema reconocido, gana la regla que aparece PRIMERO en el archivo.** Esto no es un detalle menor — es la razón por la que el archivo real tiene este orden exacto:

```jflex
/* 1. descartables */
{ComentLinea}   { }
{ComentMulti}   { }
{Blancos}       { }

/* 2. reservadas — ANTES que {Id} */
"program"                 { return symbol(sym.PROGRAM); }
"var"                     { return symbol(sym.VAR); }
...
"graphbar"  | "grapbar"   { return symbol(sym.GRAPH_BAR); }
...

/* 3. símbolos — el más largo primero cuando hay prefijo compartido */
"::"                      { return symbol(sym.DOBLE_DOS_PUNTOS); }
":"                       { return symbol(sym.DOS_PUNTOS); }
"<-"                      { return symbol(sym.ASIGNACION); }
...

/* 4. con patrón — DESPUÉS de las reservadas */
{Numero}                  { return symbol(sym.NUMERO, Double.valueOf(yytext())); }
{Cadena}                  { return symbol(sym.CADENA, yytext().substring(1, yylength()-1)); }
{IdArreglo}               { return symbol(sym.ID_ARREGLO); }
{Id}                      { return symbol(sym.ID); }

/* 5. cualquier otra cosa = error léxico */
[^]                       { /* ... */ }
```

Si invirtieras el orden y pusieras `{Id}` antes que `"var"`, el lexer **nunca** produciría el token `VAR`: para la entrada `var`, tanto la regla de `"var"` como la de `{Id}` matchean exactamente 4 caracteres (empate de longitud), así que ganaría la que esté primero — si es `{Id}`, `var` se reconocería como un identificador común, y tu programa nunca podría declarar una variable. Este es el error de diseño más fácil de cometer y el más difícil de diagnosticar a simple vista (compila bien, corre, pero silenciosamente ninguna palabra reservada funciona).

Con `"::"` antes que `":"` pasa lo mismo pero por un motivo distinto: acá **no hay empate**, JFlex ya usa *longest match* (la coincidencia MÁS LARGA gana siempre, sin importar el orden) — así que en rigor el orden de `::`/`:`  no importa tanto como el de las reservadas/`{Id}` (que si empatan en longitud). Vale la pena que entiendas la diferencia entre ambas reglas de desambiguación de JFlex: **longest match primero, orden de declaración como desempate**.

El typo sistemático del enunciado (abre con `graphBar(` pero pide ejecutar `EXEC grapBar`, sin la segunda "h") se resuelve con una alternancia en la misma regla, no con dos reglas separadas — así ambas formas producen el mismo token y el resto del sistema ni se entera de que hay dos grafías:

```jflex
"graphbar"  | "grapbar"   { return symbol(sym.GRAPH_BAR); }
```

**Decisión que NO tomás acá**: los atributos de gráfica (`titulo`, `ejeX`, `values`, `label`...) **no son palabras reservadas**. El enunciado los reutiliza como nombre de variable común en otros contextos (`var:char[]:: titulo <- "Hola" end;`), así que declararlos como tokens propios rompería esos casos. Se dejan como `ID` normal, y su validez como nombre de atributo de gráfica se revisa en **tiempo de ejecución** (`Entorno.registrarGrafica`, sección 4.5), no en el lexer ni en el parser.

### 2.5 Conectar el lexer con el resto del sistema (aunque todavía no exista)

Aunque en esta etapa no tengas ni parser ni intérprete, dejá el gancho preparado — te va a ahorrar tocar `Lexer.flex` otra vez más adelante:

```jflex
%{
  public dataforge.interprete.Entorno entorno;

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

El campo `entorno` empieza en `null` y así se queda mientras solo tengas `TestLexer` corriendo (por eso el `if (entorno != null)` en `registrar`: el lexer tiene que poder correr solo, sin intérprete). Cuando construyas `Interprete.ejecutar(...)` en la sección 4, vas a asignarle un valor real — y en ese momento, sin tocar una sola línea más de `Lexer.flex`, empieza a registrar cada token para el reporte §6.1.

La regla de error léxico usa el mismo patrón defensivo:

```jflex
[^]  { if (entorno != null) {
         entorno.error("Léxico", "el carácter '" + yytext()
             + "' no pertenece al lenguaje", yyline, yycolumn);
       } else {
         System.err.printf("Error léxico: '%s' no pertenece al lenguaje (línea %d, columna %d)%n",
             yytext(), yyline + 1, yycolumn + 1);
       } }
```

**Por qué el análisis CONTINÚA en vez de detenerse**: `[^]` matchea *cualquier carácter individual* que no matcheó ninguna regla anterior. La acción no hace `return` — no entrega ningún token al parser, simplemente descarta ese carácter y JFlex sigue escaneando desde el siguiente. Es la recuperación léxica más simple posible (Dragón §3, "transductores con acciones de pánico mínimas"): un carácter roto no debería tirar abajo el análisis de los 500 caracteres siguientes.

### 2.6 Verificación

Escribí una clase de prueba mínima (`TestLexer`) que recorra `lexer.next_token()` hasta `sym.EOF` e imprima lexema/token/línea/columna. Contra `entradas/ejemplo1.df` (comentarios, dos declaraciones, un arreglo con aritmética anidada, una función estadística, `print` y `column`) el resultado verificado es **73 tokens** exactos, incluyendo el `EOF` implícito de `PROGRAM`/`END`/`PROGRAM` al final. Si tu conteo da un número distinto, sospechá primero de comentarios mal descartados o de un símbolo compuesto (`::`, `<-`) partido en dos tokens sueltos.

### 2.7 Errores comunes de esta sección

- **Olvidar `%public`** → `'Lexer' is not public...` al compilar desde otro paquete.
- **Olvidar `%cup`** → el lexer generado no produce objetos `Symbol` compatibles; CUP no compila contra él.
- **Reglas en mal orden** (reservada después de `{Id}`) → la reservada nunca se reconoce, sin ningún error visible.
- **Convertir el número con `Double.valueOf(yytext())` y no guardar también el lexema original** → el reporte de tokens (§6.1) tiene que mostrar `1`, no `1.0`; si solo guardás el `Double` convertido, perdés la forma en que el usuario lo escribió. Por eso `registrar(type)` recibe `yytext()` (el String crudo), no el valor ya convertido.
- **Editar `target/generated-sources/.../Lexer.java` directamente** para "arreglar algo rápido" → se pierde en el siguiente `mvn compile`. El único archivo fuente válido es `Lexer.flex`.

---

## 3. Análisis sintáctico — CUP

### 3.1 El problema: por qué una lista de tokens no alcanza

Con el lexer tenés una secuencia plana de tokens. Pero eso no te dice si `SUM ( PAR_IZQ NUMERO COMA PAR_DER` es válido o no — necesitás una **gramática** que describa el ORDEN correcto. Y no cualquier técnica sirve: mirá este caso real del enunciado (§5.8, ejemplo anidado):

```
DIV( SUM( Max(@notas), Min(@notas) ), 2 )
```

Para validar que los paréntesis casan, hace falta **contar** la profundidad de anidación, y una expresión regular no tiene memoria para eso (un AFD tiene estados finitos; "contar paréntesis arbitrariamente anidados" necesitaría infinitos estados). Esto es exactamente la limitación que el Dragón describe en el capítulo 4.1: las gramáticas libres de contexto (con una pila, no con estados finitos) sí pueden expresar anidación arbitraria; las expresiones regulares, no.

### 3.2 Derivar la BNF del enunciado

El método es el mismo que en léxico: leer el enunciado y traducir cada construcción a una producción. Por ejemplo, §5.5 (declarar variable):

> `var:<TIPO>::<ID> <- <EXPRESION> end;`

se traduce directo a:

```
<declaracion-var> ::= "var" ":" <tipo> "::" ID "<-" <expresion> "end" ";"
```

Vas a notar un patrón que se repite en casi TODAS las instrucciones del lenguaje: terminan en `"end" ";"`. Aprovechá ese patrón — no memorices 15 producciones sueltas, memorizá la forma común y las excepciones.

La pieza que hace posible anidar funciones sin límite es la **recursión mutua** entre `<aritmetica>` y `<expresion>`:

```
<expresion>     ::= NUMERO | CADENA | ID | <aritmetica> | <estadistica>
<aritmetica>    ::= <operacion> "(" <expresion> "," <expresion> ")"
```

Una aritmética contiene expresiones, y una expresión puede SER una aritmética — esa mutua-referencia es la que te da profundidad infinita sin escribir una regla extra por cada nivel de anidación. Es el mismo mecanismo, aplicado a un lenguaje real, que resuelve el problema del punto 3.1.

El archivo entregable (`docs/gramatica.txt`) es una BNF **limpia**, pensada para un lector humano — el enunciado (§8) exige explícitamente que **no** sea una copia del `.cup`. Mantené los dos en sincronía a mano cada vez que cambies la gramática, pero no confundas "sincronizados" con "idénticos": el `.cup` tiene tipos, acciones y nombres de variable de CUP; `gramatica.txt` no debería tener nada de eso.

### 3.3 La estructura de un `.cup`

```cup
package dataforge.analisis;
import java_cup.runtime.*;
import dataforge.interprete.Entorno;
import dataforge.interprete.Operaciones;

parser code {:
    public Entorno entorno = new Entorno();
    public void syntax_error(Symbol s) {
        entorno.error("Sintáctico", "no se esperaba '" + s.value + "'", s.left, s.right);
    }
:};

terminal PROGRAM, END, VAR, ARR, DOUBLE, CHAR;
terminal SUM, RES, MUL, DIV, MOD;
/* ... el resto de los terminal, TODOS antes de cualquier producción ... */
terminal Double NUMERO;
terminal String CADENA, ID, ID_ARREGLO;

non terminal inicio, lista_instr, instruccion;
non terminal Object expr, aritmetica, estadistica;
non terminal ArrayList lista_expr;
/* ... TODAS las non terminal antes de cualquier producción ... */

start with inicio;

inicio ::= PROGRAM lista_instr END PROGRAM ;
/* ... producciones ... */
```

Convención de proyecto que **CUP exige, no es opcional**: todas las declaraciones `terminal` y `non terminal` van antes de la primera producción. Si intercalás una producción entre dos bloques de declaraciones, CUP tira error de compilación de la gramática (no un error de Java — un error al generar `Parser.java`).

Otra convención menos obvia: los `non terminal` que transportan un valor se declaran con **tipos raw** (`ArrayList`, no `ArrayList<Object>`, no arrays). CUP genera código intermedio que castea estos valores, y los genéricos en ese contexto no son confiables — es una limitación conocida de la herramienta, no un capricho del proyecto.

### 3.4 Terminales tipados: la puerta de entrada del valor

```cup
terminal Double NUMERO;
terminal String CADENA, ID, ID_ARREGLO;
```

Esto no es cosmético: le dice a CUP que el `Symbol` de un `NUMERO` trae adentro un `Double` de verdad (el que el lexer construyó con `Double.valueOf(yytext())`), no un `String` que haya que volver a convertir en el parser. El parser **confía** en el tipo que el lexer ya preparó — la conversión ocurre una sola vez, en el lugar correcto (léxico, donde tenés el lexema crudo).

### 3.5 Acciones semánticas S-atribuidas (sin AST)

Acá es donde DataForge toma su decisión de diseño más importante. Mirá la acción real de una aritmética:

```cup
aritmetica  ::= op_arit:op PAR_IZQ expr:a COMA expr:b PAR_DER
                {: RESULT = Operaciones.aritmetica(op, a, b,
                                parser.entorno, opleft, opright); :} ;
```

`a` y `b` ya llegan **calculados** cuando esta acción se ejecuta — no son nodos de un árbol que vas a recorrer después, son valores Java ya resueltos (`Double`, `String`, lo que corresponda). Esto es posible por dos razones que tienen que darse juntas:

1. **La gramática es S-atribuida**: cada producción sintetiza su propio valor (`RESULT`) a partir de los valores YA calculados de sus componentes. No hay atributos heredados (nada "baja" desde el padre) — todo "sube".
2. **El parser es LR (ascendente)**: reduce las producciones internas ANTES que las externas que las contienen. Cuando `SUM(3, 4)` reduce como `aritmetica`, ya se ejecutó y devolvió `7.0` antes de que el parser intente reducir la `decl_arr` que la contiene.

La combinación de ambas cosas es la que te permite **ejecutar directamente en las acciones**, sin construir un árbol de sintaxis abstracta (AST) intermedio. Esta decisión es válida específicamente porque **DataForge no tiene control de flujo**: no hay condicionales ni ciclos, así que cada instrucción del programa se ejecuta exactamente una vez, en el mismo orden en que aparece en el texto — ejecutar "sobre la marcha" mientras el parser reduce da exactamente el mismo resultado que construir un árbol y recorrerlo después.

**Por qué esto DEJARÍA de ser válido en un lenguaje con `if`/`while`** (como el próximo proyecto del curso): si tuvieras que ejecutar el cuerpo de un `while` cinco veces, no podés "ejecutar mientras parseás" — el parser solo ve el código fuente una vez, de arriba a abajo. Necesitás construir una representación (el AST) que puedas recorrer las veces que haga falta, en el orden que el programa decida en tiempo de ejecución (no en el orden en que aparece el texto). Guardate esta distinción — es la primera pregunta que te van a hacer al defender por qué DataForge no tiene AST y por qué el próximo proyecto sí lo va a necesitar.

### 3.6 Recuperación en modo pánico

```cup
instruccion ::= decl_var | decl_arr | imprimir | columna | grafica
              | error PUNTO_COMA ;
```

El terminal especial `error` es una palabra reservada de CUP (Dragón §4.8.3, "modo pánico"): cuando el parser encuentra un token que no puede desplazar en ningún estado válido, en vez de abortar, **descarta símbolos de la pila hasta poder desplazar `error`**, y después **descarta tokens de la entrada hasta el próximo `;`** — que actúa como "punto de sincronización" porque toda instrucción de DataForge termina en `;` sin excepción. El análisis sigue con la instrucción siguiente.

Costo real de esta recuperación: la instrucción rota se **pierde completa** — no hay forma de "salvar" parcialmente `var:double:: roto 5 end;` (le falta el `<-`); CUP descarta desde donde detectó el problema hasta el próximo `;`, así que ni `roto` llega a declararse. Contrastá esto con un error léxico (que descarta un solo carácter y dentro de la misma instrucción todo lo demás sobrevive) — son recuperaciones de capas distintas, con alcance distinto.

### 3.7 Verificación

Con `TestParser` (lexer → parser, sin ejecutar nada — solo responde "¿es válido?"), un archivo bien formado imprime `[OK]`. Pero ojo con una trampa: por el modo pánico, un archivo CON un error sintáctico **también** puede terminar en `[OK]`, porque la recuperación evita que el análisis aborte — el error se ve únicamente en el mensaje que `syntax_error()` imprime durante el proceso, no en el veredicto final. "OK" en `TestParser` significa "el análisis terminó sin abortar", no "sin errores". Es un buen ejercicio para entender la diferencia entre "detectar un error" y "abortar por un error" — DataForge hace lo primero, nunca lo segundo.

### 3.8 Errores comunes de esta sección

- **Declarar una producción entre dos bloques de `terminal`/`non terminal`** → error de compilación de la gramática al generar `Parser.java`.
- **Usar `ArrayList<Object>` en vez de `ArrayList` raw en un `non terminal`** → comportamiento indefinido o errores crípticos de CUP; usá siempre tipos raw.
- **Olvidar la producción con `error`** → `Couldn't repair and continue parse`: CUP no tiene ninguna forma de recuperarse y aborta en el primer error sintáctico, lo que te impide reportar TODOS los errores de una corrida (exigencia del §6.2 del enunciado).
- **Confundir "recuperado" con "sin error"** → como viste en 3.7, el modo pánico hace que el análisis termine, pero el error igual quedó registrado; no asumas que `[OK]` significa "programa perfecto".

---

## 4. El Entorno y la ejecución

### 4.1 El problema: dónde vive el estado de un programa corriendo

El parser valida FORMA. Alguien tiene que darle SIGNIFICADO: guardar qué variables existen, calcular el resultado de `SUM(3,4)`, acumular lo que hay que imprimir, y recordar qué errores semánticos ocurrieron. Ese "alguien" es una única instancia de `Entorno`, expuesta como campo público del parser:

```cup
parser code {:
    public Entorno entorno = new Entorno();
    public void syntax_error(Symbol s) { ... }
:};
```

Todas las acciones semánticas acceden a este estado vía `parser.entorno`, por ejemplo:

```cup
decl_var ::= VAR DOS_PUNTOS tipo:t DOBLE_DOS_PUNTOS ID:id
             ASIGNACION expr:e END PUNTO_COMA
             {: parser.entorno.declararVariable(id, t, e, idleft, idright); :} ;
```

### 4.2 La tabla de símbolos: por qué `LinkedHashMap` y no `HashMap`

```java
private final LinkedHashMap<String, Simbolo> simbolos = new LinkedHashMap<>();
```

Dos requisitos simultáneos que un `HashMap` normal no te da:

1. **Búsqueda O(1) por nombre** (para `valorDe`/`valorArreglo`) — cualquier `Map` hash te lo da.
2. **Orden de DECLARACIÓN preservado** (el reporte §6.3 tiene que listar las variables en el orden en que el programa las declaró, no alfabético) — esto es lo que un `HashMap` normal **no** garantiza, y `LinkedHashMap` sí, sin sacrificar el punto 1.

El lenguaje es *case insensitive* (§5.1) también para identificadores, así que la CLAVE del mapa se normaliza a minúsculas:

```java
public void declararVariable(String id, String tipo, Object valor, int l, int c) {
    String clave = id.toLowerCase();
    if (simbolos.containsKey(clave)) {
        error("Semántico", "'" + id + "' ya fue declarada (línea "
                + simbolos.get(clave).linea + ")", l, c);
        return;
    }
    if (valor == null) return;  // la expresión ya reportó su propio error
    boolean ok = (tipo.equals("double") && valor instanceof Double)
              || (tipo.equals("char[]") && valor instanceof String);
    if (!ok) {
        error("Semántico", "no se puede asignar " + nombreTipo(valor)
                + " a la variable '" + id + "' de tipo " + tipo, l, c);
        return;
    }
    simbolos.put(clave, new Simbolo(id, "variable", tipo, valor, l + 1, c + 1));
}
```

Pero el `Simbolo` guarda el **nombre original** que escribió el usuario (`nombre`, con mayúsculas y todo) — la normalización es solo para la búsqueda, nunca para lo que se muestra en un reporte. Si mezclás las dos cosas (buscás y mostrás con la clave en minúsculas), el reporte de símbolos le va a mostrar al usuario un nombre distinto al que escribió.

### 4.3 La convención más importante del proyecto: propagación de errores por `null`

Cuando algo no puede evaluarse, el método que detecta el problema registra el error **una sola vez** y devuelve `null`. Todo lo que reciba ese `null` como entrada, en vez de volver a reportar el mismo problema, **calla y también devuelve `null`**:

```java
// Operaciones.aritmetica
if (a == null || b == null) return null;
```

```java
// Entorno.declararVariable
if (valor == null) return;  // la expresión ya reportó su propio error
```

**Por qué esta regla y no lanzar una excepción de Java**: una excepción no capturada abortaría TODO el análisis en el primer error semántico. El enunciado (§6.2) exige que el reporte de errores acumule **todos** los errores de una corrida, no que se detenga en el primero. `null` te da eso "gratis": la expresión rota no produce valor, pero el programa sigue evaluando las líneas siguientes con normalidad.

**Por qué no reportar el error de nuevo en cada nivel**: si no tuvieras esta convención, una expresión como `SUM(fantasma, 5)` generaría dos errores (uno por `fantasma` no declarada, y otro falso "SUM recibió tipo incorrecto" al ver que el primer operando es `null`) por un solo problema real. La regla "si me llega `null`, yo también devuelvo `null` sin quejarme" corta la cascada en la fuente.

Vas a ver este patrón calcado en `Operaciones.aritmetica`, `Operaciones.estadistica`, `Entorno.declararVariable` y `Entorno.declararArreglo` — es la convención más repetida de todo el proyecto. Cuando escribas una función nueva que reciba valores ya evaluados, replicala.

### 4.4 Aritmética y estadística con chequeo de tipos

```java
public static Object aritmetica(String op, Object a, Object b, Entorno ent, int l, int c) {
    if (a == null || b == null) return null;
    if (!(a instanceof Double) || !(b instanceof Double)) {
        ent.error("Semántico", op + " solo acepta valores double (recibió "
                + Entorno.formatear(a) + " y " + Entorno.formatear(b) + ")", l, c);
        return null;
    }
    double x = (Double) a, y = (Double) b;
    switch (op) {
        case "SUM": return x + y;
        case "RES": return x - y;
        case "MUL": return x * y;
        case "MOD":
            if (y == 0) {
                ent.error("Semántico", "módulo entre cero", l, c);
                return null;
            }
            return x % y;
        case "DIV":
            if (y == 0) {
                ent.error("Semántico", "división entre cero", l, c);
                return null;
            }
            return x / y;
    }
    return null;
}
```

**Nota de auditoría real (2026-07-21)**: `MOD` no siempre validó división entre cero — durante una revisión posterior se detectó que solo `DIV` tenía el chequeo `if (y == 0)`, y `MOD` lo dejaba pasar directo a `x % y`. En Java, `x % 0.0` con `double` **no lanza excepción** (a diferencia de la división entera): devuelve silenciosamente `NaN`. Sin el chequeo explícito, ese `NaN` se hubiera colado hasta la consola o un reporte, violando la convención de la sección 4.3 (todo error semántico tiene que pasar por `entorno.error(...)` y devolver `null`, nunca un valor "raro" disfrazado de resultado válido).

**La lección general, más allá de este bug puntual**: cuando dos operaciones son semánticamente parecidas (`DIV` y `MOD` comparten "un divisor no puede ser cero"), **replicá la validación explícitamente en ambas** — no asumas que el chequeo de una alcanza para la otra solo porque están en el mismo `switch`. Es un buen hábito de revisión: cada vez que agregues una función nueva a una familia (`SUM`/`RES`/`MUL`/`DIV`/`MOD`, o `Media`/`Mediana`/`Moda`/...), repasá si las validaciones de sus "hermanas" también le aplican a la nueva.

`estadistica(...)` sigue la misma convención de entrada/salida, pero valida algo distinto: que **todos** los elementos del arreglo sean `Double`, y que el arreglo no esté vacío:

```java
public static Object estadistica(String fn, ArrayList<Object> arr, Entorno ent, int l, int c) {
    if (arr == null) return null;
    ArrayList<Double> datos = new ArrayList<>();
    for (Object v : arr) {
        if (v == null) return null;
        if (!(v instanceof Double)) {
            ent.error("Semántico", fn + " requiere un arreglo de tipo double", l, c);
            return null;
        }
        datos.add((Double) v);
    }
    if (datos.isEmpty()) {
        ent.error("Semántico", fn + " recibió un arreglo vacío", l, c);
        return null;
    }
    switch (fn) {
        case "Media":    return media(datos);
        /* ... */
    }
    return null;
}
```

### 4.5 Dos formateadores — no los mezcles

```java
/** Formato de CONSOLA: 15.0 se muestra "15"; las cadenas van sin comillas. */
public static String formatear(Object v) {
    if (v == null) return "null";
    if (v instanceof Double d) {
        if (d == Math.floor(d) && !d.isInfinite()) return String.valueOf(d.longValue());
        return String.valueOf(d);
    }
    return String.valueOf(v);
}

/** Formato de REPORTE (§6.3): cadenas CON comillas, arreglos elemento por elemento. */
public static String valorReporte(Object v) {
    if (v == null) return "null";
    if (v instanceof String s) return "\"" + s + "\"";
    if (v instanceof ArrayList<?> lista) {
        StringBuilder b = new StringBuilder("[");
        for (int i = 0; i < lista.size(); i++) {
            if (i > 0) b.append(", ");
            b.append(valorReporte(lista.get(i)));
        }
        return b.append("]").toString();
    }
    return formatear(v);
}
```

Ambos existen porque el enunciado usa convenciones DISTINTAS en `console::print` (§5.9: `"hola"` se imprime `hola`, sin comillas) y en la tabla de símbolos (§6.3, ejemplo: `"Hola Mundo"` se muestra CON comillas). Si programás un solo formateador "genérico" para los dos casos, vas a violar uno de los dos formatos del enunciado apenas tu programa tenga una cadena. Mantenelos separados a propósito, aunque parezca redundante.

### 4.6 Gráficas: "la última instrucción gana"

```java
public void registrarGrafica(String tipo, ArrayList<?> attrs, int l, int c) {
    LinkedHashMap<String, Object> mapa = new LinkedHashMap<>();
    boolean exec = false;
    for (Object o : attrs) {
        Object[] par = (Object[]) o;
        String nombre = (String) par[0];
        if (nombre.equals("EXEC")) { exec = true; continue; }
        String clave = nombre.toLowerCase();
        if (!ATTRS_VALIDOS.contains(clave)) {
            error("Semántico", "'" + nombre + "' no es un atributo de gráfica válido", l, c);
            continue;
        }
        mapa.put(clave, par[1]);   // put repetido = la última instrucción gana
    }
    if (!exec) return;             // sin EXEC no se muestra — no es error
    if (!validarGrafica(tipo, mapa, l, c)) return;
    if (tipo.equals("Histogram")) tablaHistograma(mapa);
    graficas.add(new Grafica(tipo, mapa));
}
```

El enunciado (§5.10) especifica que si un atributo se asigna dos veces antes de `EXEC`, **gana el último valor asignado**. La forma más simple de implementar esto NO es comparar manualmente "¿ya existe esta clave? ¿la reemplazo?" — es aprovechar que `LinkedHashMap.put(clave, valor)` **ya sobrescribe** el valor anterior si la clave existe. Recorrés la lista de atributos en orden de aparición y hacés `put` sin condicional — el último `put` de una clave repetida automáticamente "gana" porque pisa al anterior. Es una decisión de estructura de datos que reemplaza lógica condicional explícita.

Fijate también que **sin `EXEC` no hay error** (`if (!exec) return;`) — el enunciado dice que `EXEC` es lo que "ejecuta que la gráfica se muestre", no una validación de forma. Un bloque de gráfica sin `EXEC` es sintácticamente inválido según la gramática (mirá `attr` en 3.2: la única forma de cerrar el bloque incluye el `EXEC`), así que en la práctica este caso no ocurre salvo que decidas relajar la gramática — pero el chequeo está ahí por defensividad.

`validarGrafica` centraliza qué atributos exige cada tipo, con qué tipo de dato, y (donde aplica) que los tamaños coincidan:

```java
case "graphPie" -> {
    ok &= exigirTexto(tipo, m, "titulo", l, c);
    ok &= exigirLista(tipo, m, "label", String.class, l, c);
    ok &= exigirLista(tipo, m, "values", Double.class, l, c);
    if (ok && ((ArrayList<?>) m.get("label")).size() != ((ArrayList<?>) m.get("values")).size()) {
        error("Semántico", tipo + ": label y values deben tener la misma cantidad de elementos", l, c);
        ok = false;
    }
}
```

Nota importante: los nombres de atributo (`titulo`, `ejex`, `values`, `label`...) **no son palabras reservadas del lexer** (sección 2.4) — llegan acá como `String` comunes, y esta es la ÚNICA capa que decide si son válidos. Si escribís `titol` en vez de `titulo`, el lexer y el parser lo aceptan sin quejarse (es un `ID` válido); el error "no es un atributo de gráfica válido" recién aparece acá, en tiempo de ejecución.

### 4.7 La fachada `Interprete`: por qué existe

```java
public class Interprete {
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
}
```

Tres decisiones concentradas en pocas líneas:

1. **`StringReader` en vez de `FileReader`**: el `Lexer` generado por JFlex está escrito contra la interfaz `Reader`, no contra una implementación concreta. Programar contra la abstracción desde la Etapa de léxico es lo que te permite, más adelante, alimentar el mismo lexer con texto que viene de un `TextArea` de la GUI en vez de un archivo en disco — sin tocar `Lexer.flex`.
2. **`lexer.entorno = parser.entorno;`** — esta línea es la que conecta de verdad el lexer con el resto del sistema (el gancho que dejaste preparado en 2.5). Sin ella, el lexer sigue funcionando (produce tokens igual), pero nadie registra esos tokens para el reporte §6.1, y los errores léxicos se van por `stderr` en vez de sumarse a la lista de errores del `Entorno`. **Si programás tu propia clase de prueba de consola sin esta línea** (como pasa con `TestInterprete`/`TestParser`, ver sección 7), vas a ver ese mismo comportamiento — no es un bug del intérprete, es que esa clase específica no wireó el lexer.
3. **Un `Parser` (y por lo tanto un `Entorno`) nuevo por llamada**: como `entorno` se instancia en el campo `parser code` del `.cup` (`public Entorno entorno = new Entorno();`), cada `new Parser(lexer)` trae un `Entorno` fresco. Esto es intencional y no accidental: el §6 del enunciado exige que los reportes reflejen **solo el último análisis** — si reutilizaras el mismo `Entorno` entre corridas, las variables y errores de una ejecución anterior se mezclarían con la actual.

### 4.8 Errores comunes de esta sección

- **Reportar el mismo error dos veces** al no propagar `null` correctamente (revisar 4.3 antes de escribir una función nueva que combine varios operandos).
- **Validar solo una operación de una familia simétrica** (el bug real de `MOD`/`DIV`, sección 4.4) — cuando agregues una función, repasá las validaciones de sus "hermanas".
- **Mezclar `formatear()` y `valorReporte()`** — usar el de consola en un reporte (o viceversa) produce cadenas sin comillas donde el enunciado las pide con comillas, o arreglos con `.0` de más.
- **Comparar claves de la tabla de símbolos sin normalizar** (olvidarte del `.toLowerCase()` en una búsqueda nueva) — rompe el case-insensitive del §5.1 justo en el punto donde lo agregaste.
- **Reutilizar el mismo `Entorno` entre ejecuciones** — mezcla el estado de dos corridas distintas y viola el §6 del enunciado.

---

## 5. GUI JavaFX

### 5.1 El problema: conectar el intérprete de consola con una ventana real

Hasta acá tenés un intérprete que funciona perfecto por consola (`Interprete.ejecutar(String) → Entorno`). El enunciado (§4) pide una interfaz gráfica con pestañas, botones y una consola de solo lectura. La buena noticia, gracias a la decisión de la sección 4.7 (programar contra `Reader`), es que **no vas a tocar ni una línea del lexer o el parser** para esto — es pura ingeniería de UI.

### 5.2 El *scene graph*, construido por código (sin FXML)

```java
public class EditorApp extends Application {
    private TabPane pestanas;
    private TextArea consola;
    private Entorno ultimoEntorno;

    @Override
    public void start(Stage stage) {
        consola = new TextArea();
        consola.setEditable(false);

        pestanas = new TabPane();

        Button bEjecutar = new Button("▶ Ejecutar");
        bEjecutar.setOnAction(e -> ejecutar());
        HBox barra = new HBox(8, /* ... */ bEjecutar, /* ... */);

        SplitPane centro = new SplitPane(pestanas, consola);
        centro.setOrientation(Orientation.VERTICAL);

        BorderPane raiz = new BorderPane(centro);
        raiz.setTop(barra);

        stage.setScene(new Scene(raiz, 1000, 700));
        stage.show();
    }
}
```

La jerarquía es: `Stage` (la ventana del sistema operativo) → `Scene` (el contenido) → `BorderPane` (raíz) → `HBox` de botones arriba + `SplitPane` vertical en el centro (`TabPane` del editor arriba, `TextArea` de consola abajo). Se eligió construir esto **por código en vez de con FXML** a propósito: con solo media docena de controles, escribir el árbol a mano deja la estructura `Stage → Scene → Nodos` completamente visible en el propio código — que es justo lo que tenés que entender para este tema. FXML esconde esa jerarquía en un XML aparte.

### 5.3 `ejecutar()`: el entorno fresco en acción

```java
private void ejecutar() {
    TextArea area = areaActual();
    if (area == null) { consola.setText("No hay ninguna pestaña abierta."); return; }

    ultimoEntorno = Interprete.ejecutar(area.getText());

    StringBuilder salida = new StringBuilder(ultimoEntorno.getConsola());
    if (!ultimoEntorno.getErrores().isEmpty()) {
        salida.append("\n─── ERRORES (").append(ultimoEntorno.getErrores().size()).append(") ───\n");
        for (var e : ultimoEntorno.getErrores()) salida.append(e).append('\n');
    }
    if (!ultimoEntorno.getGraficas().isEmpty()) {
        salida.append("\n(").append(ultimoEntorno.getGraficas().size()).append(" gráfica(s) mostrada(s) en pantalla)\n");
        ultimoEntorno.getGraficas().forEach(Graficador::mostrar);
    }
    consola.setText(salida.toString());
}
```

Fijate qué poco código hace falta acá gracias a todo lo construido en las secciones 2-4: `ejecutar()` no sabe nada de lexer/parser/CUP — solo llama a `Interprete.ejecutar(texto)` y consume el `Entorno` resultante con los mismos getters que ya escribiste. El campo `ultimoEntorno` es el que van a consumir tanto el dibujado de gráficas (sección 5.4) como los reportes (sección 6) — por eso se guarda como campo de la clase, no como variable local.

**Pregunta para vos**: ¿qué pasaría si reutilizaras el mismo `Entorno` entre clics de "Ejecutar" en vez de dejar que `Interprete.ejecutar` cree uno nuevo cada vez? (Repasá la sección 4.7, punto 3, si no te sale la respuesta.)

### 5.4 El truco de `Lanzador`

```java
public class Lanzador {
    public static void main(String[] args) {
        EditorApp.main(args);
    }
}
```

Una clase de una sola línea útil, no un adorno. El *launcher* estándar de Java revisa si la clase con `main` **extiende** `javafx.application.Application`; si es así y JavaFX no está declarado en el *module-path* (el caso típico al correr con ▶ Run de IDEA sin el plugin de módulos configurado), aborta con `Error: JavaFX runtime components are missing`. Como `Lanzador` no extiende `Application` (solo llama al `main` de `EditorApp` desde adentro), ese chequeo específico del launcher nunca se dispara, y la aplicación corre con el classpath plano que arma IDEA desde el `pom.xml`. La alternativa que sí resuelve el *module-path* correctamente es `mvn clean javafx:run` — ambos caminos funcionan, elegís uno según si estás en el IDE o en terminal.

**Corolario práctico**: nunca ejecutes `EditorApp` directamente con ▶ Run — vas a reproducir el error de arriba. Es un buen error para ver una vez en vivo (y entender por qué existe `Lanzador`) antes de acostumbrarte a correr siempre sobre `Lanzador`.

### 5.5 Errores comunes de esta sección

- **Correr `EditorApp` en vez de `Lanzador`** → `JavaFX runtime components are missing`.
- **Reconstruir el `Entorno` a mano en vez de usar `Interprete.ejecutar`** → te perdés el `lexer.entorno = parser.entorno` de la fachada, y con eso el registro de tokens y errores léxicos (mismo problema de la sección 4.7, punto 2).
- **Dibujar las gráficas antes de terminar de ejecutar** (llamar a `Graficador.mostrar` dentro de una acción semántica del `.cup`, por ejemplo) — rompe la separación de capas: el intérprete no debería saber que existe una GUI. Las gráficas se acumulan como objetos `Grafica` durante la ejecución y se dibujan recién al final, desde `ejecutar()`.

---

## 6. Reportes HTML

### 6.1 El problema: generar 3 reportes sin mantener un `switch` de tokens a mano

El §6.1 del enunciado pide un reporte de tokens con el **nombre** de cada tipo (`VAR`, `NUMERO`, `ID`...), no solo su número interno. La forma obvia —un `switch` con un `case` por cada constante de `sym`— es frágil: cada vez que agregás un token nuevo en el `.cup`, tenés que acordarte de agregar también un `case` acá. `Reportes.java` evita ese mantenimiento manual con reflexión:

```java
private static LinkedHashMap<Integer, String> nombresTokens() throws Exception {
    LinkedHashMap<Integer, String> m = new LinkedHashMap<>();
    for (Field f : sym.class.getFields()) {
        if (f.getType() == int.class) m.put(f.getInt(null), f.getName());
    }
    return m;
}
```

`sym.java` (generado por CUP a partir de tus declaraciones `terminal`) tiene un campo `public static final int` por cada token — literalmente `PROGRAM`, `VAR`, `NUMERO`, etc. Recorrer esos campos con reflexión te da el mapa "número de token → nombre" automáticamente, sincronizado siempre con la gramática actual, sin que vos tengas que tocar `Reportes.java` cada vez que agregás un token.

### 6.2 Los tres reportes comparten una plantilla

```java
public static File[] generar(Entorno ent, File carpeta) throws Exception {
    carpeta.mkdirs();
    File t = new File(carpeta, "tokens.html");
    File e = new File(carpeta, "errores.html");
    File s = new File(carpeta, "simbolos.html");
    Files.writeString(t.toPath(), tokens(ent));
    Files.writeString(e.toPath(), errores(ent));
    Files.writeString(s.toPath(), simbolos(ent));
    return new File[]{ t, e, s };
}
```

Cada uno de `tokens(ent)`, `errores(ent)`, `simbolos(ent)` arma sus filas y llama a la misma plantilla `pagina(titulo, columnas, filas)` con CSS embebido — así los tres archivos son **autocontenidos** (abren sin conexión a internet, tal como pide el enunciado implícitamente al pedir "reportes en formato HTML").

Los datos siempre vienen del **último** `Entorno` ejecutado (el `ultimoEntorno` de `EditorApp`, section 5.3) — nunca se acumulan reportes de corridas anteriores, cumpliendo el §6 del enunciado.

### 6.3 Escapar HTML: un detalle que rompe el reporte si lo saltás

```java
private static String esc(String s) {
    return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;");
}
```

Algunos lexemas reales del lenguaje contienen `<` (el operador `<-`). Si insertás ese lexema sin escapar dentro de una celda `<td>...</td>`, el navegador interpreta `<-` como el inicio de una etiqueta rota y el reporte se desarma visualmente a partir de ahí. `esc(...)` se aplica a **todo texto que viene del usuario** (lexemas, descripciones de error, nombres de variable) antes de insertarlo en la tabla — nunca a texto que vos escribiste literal en la plantilla.

### 6.4 Errores comunes de esta sección

- **Mantener un `switch` manual para nombres de token** en vez de reflexión sobre `sym` — funciona, pero te vas a olvidar de actualizarlo la primera vez que agregues un token nuevo.
- **Olvidar `esc(...)` en un campo nuevo** que agregues a una tabla — cualquier lexema con `<`, `>` o `&` te rompe el HTML generado.
- **Generar los reportes antes de haber ejecutado algún programa** — chequeá siempre que `ultimoEntorno != null` antes de llamar a `Reportes.generar(...)` (ver `EditorApp.reportes()`).

---

## 7. Errores comunes reales (consolidado)

Esta sección junta los errores que ya se documentaron en el camino (y algunos adicionales) en un solo lugar, para repasar rápido antes de programar tu propio proyecto:

| Error | Capa | Síntoma | Causa |
|---|---|---|---|
| Falta `%public` | Léxico | `'Lexer' is not public...` al compilar | JFlex genera *package-private* por defecto |
| Falta `%cup` | Léxico | El lexer no produce `Symbol` compatible con CUP | Falta la directiva que activa la interfaz CUP |
| Reservada después de `{Id}` | Léxico | La palabra reservada nunca se reconoce (sin error visible) | JFlex desempata por orden de declaración cuando hay empate de longitud |
| `NUMERO` sin guardar `yytext()` | Léxico | El reporte de tokens muestra `1.0` en vez de `1` | Se guardó el valor convertido, no el lexema original |
| Editar código en `target/generated-sources/` | Léxico/Sintáctico | El cambio "desaparece" | Se pierde en el siguiente `mvn compile`; el fuente es el `.flex`/`.cup` |
| Producción entre bloques `terminal`/`non terminal` | Sintáctico | Error al generar `Parser.java` | CUP exige todas las declaraciones antes de cualquier producción |
| `non terminal` con genéricos o arrays | Sintáctico | Comportamiento indefinido / errores crípticos | CUP requiere tipos *raw* en `non terminal` |
| Sin producción `error` | Sintáctico | `Couldn't repair and continue parse`, aborta en el primer error | Falta la recuperación en modo pánico |
| Confundir `[OK]` con "sin errores" | Sintáctico | Un archivo con error sintáctico igual imprime `[OK]` | El modo pánico evita que el análisis aborte; el error se ve en el mensaje de `syntax_error()`, no en el veredicto |
| Reportar un error dos veces | Semántico | Un solo problema genera 2+ entradas en el reporte de errores | No se propagó `null` correctamente entre operaciones |
| Validar una operación de una familia simétrica y no la otra | Semántico | Un caso "hermano" (ej. `MOD` vs `DIV`) deja pasar un valor inválido silenciosamente | No se repasaron las validaciones de las funciones relacionadas al agregar una nueva |
| Mezclar `formatear()`/`valorReporte()` | Semántico | Cadenas sin comillas en un reporte, o con comillas en la consola | Son dos formatos deliberadamente distintos, no intercambiables |
| Clave de tabla de símbolos sin normalizar | Semántico | El case-insensitive del §5.1 deja de funcionar en un caso puntual | Falta `.toLowerCase()` en una búsqueda o inserción nueva |
| Reutilizar el mismo `Entorno` entre ejecuciones | Ejecución/GUI | Los reportes mezclan datos de corridas distintas | El §6 exige reportes solo del último análisis; cada ejecución necesita un `Entorno` nuevo |
| Correr `EditorApp` en vez de `Lanzador` | GUI | `JavaFX runtime components are missing` | El launcher de Java bloquea clases que extienden `Application` sin module-path configurado |
| No wirear `lexer.entorno = parser.entorno` | Integración | Tokens no registrados; errores léxicos van a `stderr` en vez del reporte | Falta la línea de conexión que sí tiene `Interprete.ejecutar` (una clase de prueba propia puede omitirla sin darse cuenta) |
| Olvidar `esc(...)` en un campo de reporte nuevo | Reportes | El HTML generado se rompe visualmente con ciertos lexemas | Falta escapar `&`, `<`, `>` antes de insertar texto de usuario en una tabla |

---

## Cierre

Si programaste los siete temas en orden, tenés el mismo resultado que DataForge: un lexer que convierte texto en tokens, un parser que valida la gramática y ejecuta directamente (sin AST, porque el lenguaje no tiene control de flujo), un `Entorno` que le da significado a cada instrucción con una única convención de propagación de errores, una GUI que reutiliza ese intérprete sin tocarlo, y tres reportes que documentan cada análisis. La pieza que vas a tener que **repensar** en tu propio proyecto (si tiene condicionales o ciclos) es la de la sección 3.5: en cuanto haya control de flujo, la ejecución directa en las acciones del parser deja de alcanzar, y necesitás migrar a un AST que se recorra después de parsear completo.
