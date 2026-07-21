---
tags: [proyecto, guia, lexico, sintactico]
aliases: [guia conjanalyzer, tutorial conjanalyzer, conjanalyzer paso a paso]
fuente: "Enunciado ConjAnalyzer (OLC1_PT1 _S2024) + Libro del Dragón + construcción real verificada: 2026-07-20 + auditoría de calidad: 2026-07-21"
fecha: 2026-07-21
---

# Guía de elaboración — [[ConjAnalyzer]]

Guía paso a paso, escrita DESPUÉS de que el código quedó construido (a diferencia de [[DataForge]], acá los roles se invirtieron: el usuario escribió el código y esta guía documenta lo que realmente se construyó, con `mvn compile` limpio y `TestInterprete` verificado sobre los 4 casos de `entradas/`). Sigue el mismo formato que la [[Guía DataForge]] (proyecto hermano, mismo stack).

## Estado de construcción

| Etapa | Contenido | Estado |
|---|---|---|
| 1 | Tabla de tokens + lexer [[JFlex]] (`Lexer.flex`) | ✅ verificado 2026-07-20 |
| 2 | Gramática BNF + parser [[CUP]] (`parser.cup`) | ✅ verificado 2026-07-20 |
| 3 | Ejecución: conjuntos, operaciones, EVALUAR (`interprete/*`) | ✅ verificado 2026-07-20 |
| 4 | Simplificación algebraica (`Simplificador.java`) + JSON (`JsonSalida.java`, Gson) | ✅ verificado 2026-07-20; **completado con distributivas** 2026-07-21 |
| 5 | GUI [[JavaFX y Scene Builder]] (`EditorApp`, `Lanzador`) | ✅ (código completo; no se relanzó la ventana en esta sesión, se corrió el pipeline por consola con `TestInterprete`) |
| 6 | Diagrama de Venn (`DiagramaVenn.java`) | ✅ (código completo; requiere GUI para verse) |
| 7 | Reportes HTML + JSON | ✅ verificado 2026-07-20 (los 3 `.html` + `.json` se generaron en `reportes/`) |
| 8 | Auditoría de calidad (código muerto, gaps del enunciado) | ✅ hecha 2026-07-21 (ver sección 9 de esta guía) |

**PROYECTO FUNCIONAL COMPLETO.** Pendiente real: empaquetado de entrega (JAR ejecutable, manuales PDF, repo `OLC1_Proyecto1_#Carnet`).

## 1. Requisitos previos (reales)

- **IntelliJ IDEA**: importa el `pom.xml` como [[Maven]], JDK 17 (el proyecto compila con `maven.compiler.source/target=17`; en esta verificación se usó JDK 25 embebido en el IDE sin problemas). Maven viene embebido en el IDE — no se instala nada por terminal.
- Dependencias reales del `pom.xml`: `java-cup-runtime:11b-20160615`, `javafx-controls:21.0.4`, **`gson:2.11.0`** (nueva respecto a DataForge — sirve para el JSON de simplificación, sección 5.4). Plugins: `cup-maven-plugin:11b-20160615-3`, `jflex-maven-plugin:1.9.1`, `javafx-maven-plugin:0.0.8` (mainClass `conjanalyzer.gui.EditorApp`), `exec-maven-plugin:3.1.0` (mainClass `conjanalyzer.TestInterprete`).
- Ciclo: editar `.flex`/`.cup` → Maven `compile` (regenera `Lexer`/`Parser`/`sym` en `target/generated-sources/`) → ▶ Run sobre `Lanzador` (GUI) o `TestInterprete` (consola).

## 2. Análisis léxico — CONSTRUIDO ✅

### 2.1 Tabla de tokens (la real, `src/main/jflex/Lexer.flex`)

- **3 palabras reservadas**: `CONJ`, `OPERA`, `EVALUAR`.
- **4 operadores de conjuntos** (notación prefija, 4.6): `U` (unión), `&` (intersección), `^` (complemento), `-` (diferencia).
- **9 símbolos estructurales**: `->` (FLECHA), `~` (VIRGULILLA, rango), `{` `}` `(` `)` `:` `;` `,`.
- **3 con patrón**: `ID` = letra(letra\|dígito)\* · `NUMERO` = dígito+ · `SIMBOLO` = comodín, cualquier ASCII 33..126 que no calzó arriba (posible elemento de conjunto: `!`, `@`, `$`, `<`, `>`, etc.).
- **Se descartan**: comentario de línea `#...\n`, comentario multilínea `<! ... !>`, blancos.
- **Universo**: ASCII 33 (`!`) a 126 (`~`) — cualquier carácter fuera de ese rango es **error léxico** (se descarta y el análisis continúa).

### 2.2 Decisiones de diseño tomadas

1. **NO se usa `%ignorecase`** — el propio `.flex` lo aclara con un comentario: *"OJO: NO se usa %ignorecase — el lenguaje es CASE SENSITIVE (4.1)"*. Contraste explícito con [[DataForge]], que sí era case-insensitive. `conjuntoA` y `conjuntoa` son conjuntos distintos.
2. Los **operadores de conjunto son tokens reservados** (`U`, `&`, `^`, `-`), lo que automáticamente impide usarlos como nombre de conjunto/operación o como elemento (`"U"` matchea `UNION` antes de llegar a la regla de `{Id}` — el orden de las reglas en el `.flex` es lo que hace cumplir esta restricción, documentada aparte en `docs/gramatica.txt`).
3. **`SIMBOLO` es el token comodín** que hace que *cualquier* carácter ASCII imprimible no reservado sirva como elemento de un conjunto (ver 4.4 y la restricción de diseño). Esto es lo que permite `CONJ : signos -> !,?,@,$,%;` (verificado en `ejemplo4_nosimplificable.ca`).
4. El orden de las reglas es crítico (igual que en DataForge): comentarios y blancos primero (se descartan), luego reservadas, luego operadores, luego símbolos estructurales (`->` ANTES que `-` para que el longest-match no parta la flecha), luego `{Id}`/`{Numero}`, y al final el comodín `SIMBOLO` y el error léxico universal `[^]`.

### 2.3 El `.flex` real (fragmentos clave)

```jflex
%class Lexer
%public
%unicode
%cup
%line
%column
/* OJO: NO se usa %ignorecase — el lenguaje es CASE SENSITIVE (4.1) */

Id           = {Letra}({Letra}|{Digito})*
Numero       = {Digito}+
ComentLinea  = "#"[^\r\n]*
ComentMulti  = "<!"~"!>"

%%
"CONJ"     { return symbol(sym.CONJ); }
"U"        { return symbol(sym.UNION); }
"->"       { return symbol(sym.FLECHA); }   // ANTES que "-" (DIFERENCIA)
{Id}       { return symbol(sym.ID); }
{Numero}   { return symbol(sym.NUMERO); }
[\x21-\x7E] { return symbol(sym.SIMBOLO); } // comodin: cualquier otro ASCII imprimible
[^]        { entorno.error("Lexico", "el caracter '" + yytext() + "' no pertenece al lenguaje", ...); }
```

El campo público `entorno` (seteado por `Interprete.ejecutar`) es la misma técnica de [[DataForge]]: el lexer registra CADA token vía `registrarToken()` (reporte 5.2) y manda los errores léxicos al `Entorno` en vez de a `stderr` (reporte 5.3).

### 2.4 Verificación

Corrido sobre `ejemplo2_errores.ca`, que a propósito mete una `ñ` (fuera del universo ASCII 33-126): el error se reportó como
```
[Lexico] el caracter '�' no pertenece al lenguaje (linea 9, columna 17)
```
y el análisis **continuó** (el nombre del conjunto quedó truncado a `conju`, que sí se definió correctamente). `reportes/tokens.html` confirma el registro token por token con nombre de token real (columna "Tipo" = nombre del campo en `sym`, por reflexión).

## 3. Análisis sintáctico — CONSTRUIDO ✅

### 3.1 La gramática

`docs/gramatica.txt` — **verificado contra el `.cup` real y el enunciado**: coincide token por token con las declaraciones `terminal` del `.cup` (CONJ/OPERA/EVALUAR, UNION/INTERSECCION/COMPLEMENTO/DIFERENCIA, FLECHA/VIRGULILLA, LLAVE_IZQ/DER, PAR_IZQ/DER, DOS_PUNTOS/PUNTO_COMA/COMA, y los terminales con valor ID/NUMERO/SIMBOLO), y con las 5 producciones de `<operacion>`. Es BNF limpio con prosa (no una copia mecánica): explica el rango `~`, la notación prefija con ejemplos (`U U {A} {B} {C} = (A U B) U C`), y documenta explícitamente las restricciones de diseño (caracteres reservados no válidos como elemento/nombre; un conjunto no puede definirse en términos de otro). **No requirió cambios** en la auditoría de 2026-07-21 — quedó igual (gramática y `.cup`/`.flex` sin bugs encontrados).

Núcleo:
```
<inicio>              ::= "{" <sentencias> "}"
<definicion-conjunto> ::= "CONJ" ":" ID "->" <notacion> ";"
<definicion-operacion> ::= "OPERA" ":" ID "->" <operacion> ";"
<operacion>           ::= "U" <operacion> <operacion>
                        | "&" <operacion> <operacion>
                        | "-" <operacion> <operacion>
                        | "^" <operacion>
                        | "{" ID "}"
<evaluacion>          ::= "EVALUAR" "(" "{" <lista-elementos> "}" "," ID ")" ";"
```

### 3.2 El `.cup` real

Archivo: `src/main/cup/parser.cup`. Claves:
- `non terminal NodoOperacion operacion;` — a diferencia de DataForge (que no tipaba árboles), acá el no-terminal `operacion` sintetiza un **`NodoOperacion`** real: cada reducción arma un nodo (`NodoOperacion.binario('U', a, b)`, `.unario('^', a)`, `.hoja(id)`).
- La regla de recuperación en modo pánico está en el nivel de `sentencia`: `sentencia ::= def_conj | def_opera | evaluar_stmt | error PUNTO_COMA ;` — igual patrón que DataForge (Dragón §4.8.3).
- `syntax_error(Symbol s)` registra en el entorno con `s.left`/`s.right` (columna 0-based; el entorno les suma 1 al guardar).

### 3.3 Verificación

Corrido sobre `ejemplo2_errores.ca` (a propósito con 2 errores sintácticos):
```
Error sintactico: no se esperaba '->' (linea 12, columna 12)     ← CONJ : -> a~z;  (falta el ID)
Error sintactico: no se esperaba ';' (linea 15, columna 33)      ← OPERA : mala -> U {vocales} ;  (falta el 2o operando)
```
El análisis sobrevivió ambos y siguió procesando las sentencias válidas después (`OPERA : buena -> ^ {vocales};` se definió y evaluó normalmente) — confirma que el modo pánico funciona.

## 4. Ejecución — CONSTRUIDO ✅

### 4.1 El diseño: un mini-AST solo para operaciones

A diferencia de DataForge (sin AST, todo se ejecuta directo en las acciones del CUP), ConjAnalyzer **sí construye un [[Árbol de sintaxis abstracta (AST)|árbol]]**, pero acotado: `NodoOperacion` (paquete `interprete`) representa exclusivamente el árbol de una operación (`OPERA`). Las sentencias de nivel superior (`CONJ`, `EVALUAR`) siguen ejecutándose directo en las acciones del `.cup` — sin AST propio —, igual que en DataForge ([[Traducción dirigida por la sintaxis]], gramática S-atribuida). El árbol de `NodoOperacion` es necesario porque una operación se tiene que:
1. **Evaluar** (postorden, `evaluar(Map<String,Set<Character>>, Set<Character> universo)`).
2. **Reescribir** para simplificar (`Simplificador`, sección 5 de esta guía) — algo que una evaluación directa sin árbol no puede hacer.
3. **Serializar** de vuelta a notación prefija (`toPrefijo()`) para el reporte y el JSON (5.4).

```java
public Set<Character> evaluar(Map<String, Set<Character>> conjuntos, Set<Character> universo) {
    if (esHoja()) { ... }
    Set<Character> a = izq.evaluar(conjuntos, universo);
    if (esUnario()) { // complemento: universo menos A
        Set<Character> r = new LinkedHashSet<>(universo); r.removeAll(a); return r;
    }
    Set<Character> b = der.evaluar(conjuntos, universo);
    switch (op) { case 'U' -> r.addAll(b); case '&' -> r.retainAll(b); case '-' -> r.removeAll(b); }
    return r;
}
```

Nodo también expone `pertenenciaRegion(Map<String,Boolean>)`: evalúa el árbol como función booleana sobre "¿el elemento hipotético pertenece a cada conjunto base?" — la usa el [[#6b. Diagrama de Venn — CONSTRUIDO ✅|diagrama de Venn]] para pintar regiones sin recorrer el universo entero.

### 4.2 El paquete `conjanalyzer.interprete`

| Clase | Responsabilidad |
|---|---|
| `Entorno` | Estado completo de una ejecución: universo (`Set<Character>`, ASCII 33-126), [[Tabla de símbolos|conjuntos]] y operaciones (`LinkedHashMap`, claves **sin normalizar** — case sensitive real), consola (`StringBuilder`), errores (`RegistroError`), tokens del lexer |
| `Conjunto` | nombre + `Set<Character>` + `definicion` legible (`"a~z"` o `"1, 2, 3, a, b"`) + línea/columna |
| `NodoOperacion` | el mini-AST de una `OPERA` (hoja/unario/binario) — evaluación, `pertenenciaRegion`, `toPrefijo`, `copia` |
| `Operacion` | nombre + árbol + resultado ya evaluado + referencias (para el Venn) + `ResultadoSimplificacion` |
| `Simplificador` | motor de reescritura algebraica (sección 5) |
| `ResultadoSimplificacion` | leyes aplicadas + árbol simplificado + `seSimplifico` (booleano que decide el formato del JSON) |
| `RegistroError` | tipo (Léxico/Sintáctico/Semántico) + descripción + línea/columna, base 1 |
| `Interprete` | fachada `String → Entorno`: crea `Lexer` + `Parser`, conecta `lexer.entorno = parser.entorno` (mismo patrón que DataForge), atrapa la excepción del parser sin abortar la GUI |

**Errores semánticos reales** ([[Manejo de errores (léxicos, sintácticos, semánticos)]]): conjunto/operación redefinida, elemento que no es un único carácter, elemento fuera del universo, rango inválido (`a` > `b`), operación que referencia un conjunto inexistente. **No detienen la ejecución** — se acumulan en `errores` y el análisis sigue (verificado: `ejemplo2_errores.ca` reporta 4 errores de los 3 tipos y aun así define y evalúa `buena` al final).

### 4.3 Un detalle fino documentado en el propio código

El ejemplo 4.8 del enunciado muestra, para `EVALUAR ( {1, b} , operacion1 )` con `operacion1 = conjuntoA & conjuntoB`, la salida `1 -> exitoso`. Pero `conjuntoA = {1,2,3,a,b}` y `conjuntoB = a~z`, así que la intersección real es `{a, b}` — el `1` **no** pertenece. `Entorno.java` documenta esto explícitamente como una inconsistencia del propio enunciado y decide respetar la semántica real (pertenencia al conjunto resultante evaluado). La verificación lo confirma: corriendo `ejemplo1.ca` tal cual, la salida real es
```
1 -> fallo
b -> exitoso
```
en vez de lo que dice el PDF. Es una decisión de diseño defendible y documentada, no un bug — vale la pena mencionarla en la defensa oral. **Re-confirmado en la auditoría de 2026-07-21**, sin tocar el código: sigue siendo una inconsistencia del enunciado, no del proyecto.

## 5. Simplificación algebraica — CONSTRUIDA ✅ (completada 2026-07-21)

`Simplificador.java` implementa la sección 7 del enunciado (propiedades de teoría de conjuntos) como **reescritura de árbol en postorden, repetida hasta punto fijo** (máx. 100 iteraciones de guarda). Aplica **5 leyes** como transformación real:

| Ley aplicada | Regla |
|---|---|
| Doble complemento | `^^X → X` |
| DeMorgan (con guarda) | `^(X U Y) → ^X & ^Y` (y su dual) — **solo si** `X` o `Y` ya es un complemento, para que se dispare un `^^` que luego cancela; evita "inflar" árboles que no ganan nada |
| Idempotentes | `X U X`, `X & X → X` |
| Absorción | `X U (X & Y) → X` (y variantes conmutadas) |
| **Distributivas (con guarda) — AGREGADA 2026-07-21** | `(X & Y) U (X & Z) → X & (Y U Z)` (y su dual `(X U Y) & (X U Z) → X U (Y & Z)`) — **solo en el sentido que factoriza** (reduce hojas); nunca en el sentido de "expandir" un factor, misma lógica de guarda que DeMorgan |

**Hallazgo de la auditoría 2026-07-21**: la sección 7 del enunciado lista explícitamente las propiedades distributivas (`A U (B ∩ C) = (A U B) ∩ (A U C)` y su dual) como parte de las "propiedades de la teoría de conjuntos" que la sección 5.4 pide aplicar para simplificar. El código construido originalmente (2026-07-20) implementaba 4 de las 5 leyes con transformación real y dejaba la distributiva completamente sin implementar — ni siquiera como detección "no simplificable". Era un gap real y specificado, no una decisión de diseño documentada (a diferencia de conmutativa/asociativa, que sí estaban documentadas como "solo ayuda de comparación", con razón: aplicarlas standalone no reduce nada). La distributiva, en cambio, SÍ tiene un sentido que reduce el árbol (factorizar), análogo a absorción — se agregó el método `distributiva(char op, NodoOperacion a, NodoOperacion b)` en `Simplificador.java`, activado desde `paso(...)` después de intentar absorción. Se agregó un 5º caso a `entradas/ejemplo3_simplificacion.ca` (`distributiva -> U & {conjA} {conjB} & {conjA} {conjC}`) que verifica el comportamiento end-to-end: simplifica a `& {conjA} U {conjB} {conjC}` con `leyes=[Propiedades distributivas]`, resultado `{c, d, e}` (verificado matemáticamente correcto: `conjA ∩ (conjB ∪ conjC)`).

**Decisión de diseño más interesante** (ya existía, sigue vigente): conmutativa y asociativa **no se aplican solas** (reordenar no simplifica y podría ciclar) — se usan como *ayuda de comparación*. Dos subárboles se consideran equivalentes si comparten una **forma canónica**: los operandos de cadenas `U`/`&` se aplanan, se ordenan alfabéticamente y se deduplican (método `canon()`), de modo que `A U B` y `B U A` son iguales y `(A U B) U C` y `A U (B U C)` también. Cuando esa comparación no trivial es la que habilita una idempotencia, absorción o (ahora también) distributiva, el motor reporta también "Propiedades conmutativas" y/o "Propiedades asociativas".

Verificado con `ejemplo3_simplificacion.ca` (ahora 5 operaciones, cada una dispara una ley distinta) y `ejemplo4_nosimplificable.ca` (3 operaciones que a propósito no simplifican, sin cambios de comportamiento tras el agregado). Salida real:

```
demorgan = ^ & ^ & {conjA} {conjB} ^ {conjC}
  → simplificado: U & {conjA} {conjB} {conjC}   leyes=[Leyes de DeMorgan, Ley del doble complemento]
doble = ^ ^ {conjA}          → {conjA}          leyes=[Ley del doble complemento]
idempotente = U {conjA} {conjA} → {conjA}       leyes=[Propiedades idempotentes]
absorcion = U {conjA} & {conjA} {conjB} → {conjA}  leyes=[Propiedades de absorcion]
distributiva = U & {conjA} {conjB} & {conjA} {conjC} → & {conjA} U {conjB} {conjC}  leyes=[Propiedades distributivas]
```

`JsonSalida.java` (paquete `reportes`, usa **Gson** con `setPrettyPrinting().disableHtmlEscaping()`) arma el JSON exacto de la sección 5.4: objeto `{leyes, "conjunto simplificado"}` si `seSimplifico`, o el string literal `"No se puede simplificar la operacion"` si no. Verificado que ambos formatos salen byte por byte como pide el enunciado.

## 6. GUI — CONSTRUIDA ✅

- `gui/EditorApp.java`: JavaFX por código (sin FXML). `BorderPane` raíz → barra de botones (Nuevo/Abrir/Guardar/▶ Ejecutar/Reportes) + `SplitPane` horizontal: mitad izquierda es otro `SplitPane` vertical (pestañas de editor arriba, consola no editable abajo), mitad derecha es el panel de Venn navegable.
- Pestañas con `Tab.userData = File` asociado (mismo patrón que DataForge): `null` dispara "Guardar como"; `FileChooser` filtra `*.ca`.
- `Interprete.ejecutar(codigo)` crea un **`Entorno` fresco por ejecución** (cumple la sección 5: reportes solo del último análisis) — se guarda en `ultimoEntorno` para que Reportes y Venn lo usen.
- `gui/Lanzador.java`: mismo truco que DataForge — clase `main` que NO extiende `Application`, para esquivar el chequeo "JavaFX runtime components are missing" del launcher de Java al correr desde IDEA con classpath plano.

## 6b. Diagrama de Venn — CONSTRUIDO ✅

`gui/DiagramaVenn.java`: dibuja un círculo por cada conjunto BASE referenciado por la operación (1, 2 o 3 — con más de 3 no hay Venn geométrico razonable y se muestra el resultado en texto). El sombreado es **exacto, no aproximado**: para cada píxel del `Canvas` se calcula en qué círculos cae (`dx²+dy² ≤ r²` por conjunto) y se evalúa `NodoOperacion.pertenenciaRegion(region)` como función booleana — si da `true`, el píxel se pinta. El complemento también sombrea la región "fuera de todos los círculos" (con un tono más claro) porque su resultado incluye elementos del universo que no pertenecen a ningún conjunto base. Navegación entre operaciones con `◀`/`▶`/`ComboBox`, igual criterio de "un panel por operación, navegable" que pide 5.1.

## 6c. Reportes — CONSTRUIDOS ✅

`reportes/Reportes.java` genera 3 HTML autocontenidos (CSS embebido) a partir del `Entorno` de la ÚLTIMA ejecución:
- **`tokens.html`** (5.2): lexema + nombre de token (por **reflexión sobre la clase `sym`** generada por CUP — mismo truco que DataForge) + línea/columna.
- **`errores.html`** (5.3): tipo/descripción/línea/columna de los 3 tipos de error en una sola tabla.
- **`operaciones.html`** (extra, no pedido explícitamente pero útil): conjuntos y operaciones definidos con su resultado o su estado de simplificación.

`reportes/JsonSalida.java` genera `simplificacion.json` (5.4) con Gson. Todos escapan HTML (`&`, `<`, `>`) porque lexemas como `&` o `->` romperían la tabla sin escapar — verificado en `tokens.html` real (`-&gt;` para el lexema `->`).

## 7. Casos de prueba (verificados 2026-07-20, re-verificados 2026-07-21 tras la auditoría)

Comando de verificación (Maven embebido del IDE + JDK 25):
```
mvn -q compile exec:java -Dexec.mainClass=conjanalyzer.TestInterprete -Dexec.args="entradas/ejemplo1.ca"
```

- **`ejemplo1.ca`**: caso básico de la sección 4.8 del enunciado. 2 operaciones, 2 `EVALUAR`. Confirma la corrección de la semántica frente a la inconsistencia del PDF (sección 4.3 de esta guía).
- **`ejemplo2_errores.ca`**: 1 error léxico (carácter fuera de universo), 2 sintácticos (a propósito, recuperación en modo pánico), 1 semántico (conjunto inexistente) — los 4 en un solo archivo, análisis sobrevive a todos.
- **`ejemplo3_simplificacion.ca`**: dispara las **5** leyes de simplificación una por una (DeMorgan+doble complemento combinadas, doble complemento solo, idempotente, absorción, y desde 2026-07-21 también distributiva).
- **`ejemplo4_nosimplificable.ca`**: 3 operaciones que a propósito NO simplifican — confirma que el JSON devuelve el string literal exacto para cada una (sin cambios tras el agregado de distributiva, verificado).

Los 3 `.html` + `simplificacion.json` se generaron correctamente en `reportes/` en cada corrida, incluida la corrida post-auditoría.

## 8. Errores comunes (los reales de este proyecto)

- **No usar `%ignorecase`**: a diferencia de DataForge, acá es DELIBERADO (el lenguaje es case sensitive, 4.1) — si alguien lo agrega por costumbre, rompe la sección 4.1 del enunciado.
- El ejemplo 4.8 del PDF tiene una salida inconsistente con su propia definición de conjuntos — no copiarla ciegamente al reporte, la semántica correcta es la de pertenencia al conjunto resultante (documentado en `Entorno.evaluar`).
- Aplicar DeMorgan sin la guarda (`X` o `Y` ya complemento) infla árboles sin necesidad — por eso el `Simplificador` la condiciona explícitamente. La misma guarda se replicó para distributiva (solo factoriza, nunca expande).
- Comparar subárboles por `toPrefijo()` literal NO detecta `A U B` == `B U A`; hace falta la forma canónica (`canon()`, aplanado + orden + dedup) para que idempotencia/absorción/distributiva disparen con operandos en cualquier orden.
- Los caracteres reservados (`U & ^ - ~ { } ( ) : ; ,`, y el inicio de comentario `#`/`<!`) no pueden ser elemento de conjunto ni nombre — lo impone el orden de reglas del `.flex` (reservadas antes de `{Id}`/comodín), no una validación aparte.
- `Couldn't repair and continue parse`: falta el terminal `error` en alguna alternativa de `sentencia` para que el modo pánico tenga dónde sincronizar.

## 9. Auditoría de calidad — 2026-07-21

Auditoría solicitada explícitamente (no cosmética: solo bugs reales o features especificadas y faltantes), usando `codebase-memory-mcp` (grafo de llamadas de las 73 funciones/métodos del proyecto, cruzado contra el enunciado y el `ManualTecnico.md`).

**Encontrado y arreglado:**
1. **Propiedades distributivas de la sección 7 no implementadas** (`Simplificador.java`) — gap real y especificado (ver sección 5 de esta guía para el detalle completo). Se agregó `distributiva(...)`, guardada para solo factorizar. Verificado con un 5º caso en `ejemplo3_simplificacion.ca`.
2. **Código muerto**: `Entorno.getUniverso()` (`interprete/Entorno.java`) no tenía ningún llamador en todo el proyecto (confirmado con `search_graph`/`query_graph` sobre el grafo de llamadas Y con grep manual de respaldo) — ni GUI, ni reportes, ni Venn lo usaban (el universo se pasa directamente como parámetro donde hace falta, p.ej. `NodoOperacion.evaluar(conjuntos, universo)`). Se eliminó el getter.

**Revisado y SIN cambios (no eran bugs):**
- El grafo de llamadas marca `Entorno.registrarToken`, `definirConjuntoRango`, `definirConjuntoLista`, `definirOperacion` y `Entorno.evaluar` como "sin llamadores" — esto es un falso positivo esperado: sus únicos llamadores reales están en `Parser.java`/`Lexer.java` **generados** por CUP/JFlex en `target/generated-sources/`, que quedó explícitamente fuera del índice (instrucción del enunciado de la tarea, no del proyecto). No es código muerto.
- `EditorApp.start`, `TestInterprete.main`, `Lanzador.main` marcados "sin llamadores": son puntos de entrada reales (JavaFX `Application.launch()` y el JVM launcher), no llamados desde código propio — comportamiento esperado, no un hallazgo.
- `RegistroError.toString()`: no tiene llamador explícito en el grafo (los reportes arman las filas leyendo `tipo`/`descripcion`/`linea`/`columna` directamente, no vía `toString()`), pero SÍ se usa implícitamente por concatenación de string en `TestInterprete.main` (`"  " + (e++) + ". " + err`) y en `EditorApp.ejecutar()` (`salida.append(e)`) — Java invoca `toString()` automáticamente ahí. No es código muerto, es una limitación de detección del grafo con invocaciones implícitas.
- Gramática (`docs/gramatica.txt`), lexer (`Lexer.flex`) y parser (`parser.cup`) revisados línea por línea contra el enunciado y el `ManualTecnico.md`: sin inconsistencias. El orden de reglas del `.flex`, la convención `left`/`right` = línea/columna en los `Symbol` de CUP, y la recuperación en modo pánico funcionan como está documentado.
- La discrepancia del ejemplo 4.8 del enunciado (ya conocida, ver sección 4.3) se re-confirmó pero NO se tocó — es una inconsistencia del PDF, no del código.

**Verificación final**: `mvn compile` limpio, los 4 archivos de `entradas/` (`ejemplo1.ca`, `ejemplo2_errores.ca`, `ejemplo3_simplificacion.ca` con el nuevo caso, `ejemplo4_nosimplificable.ca`) corridos con `TestInterprete` sin excepciones y con salidas correctas (incluida la nueva ley distributiva y sin regresiones en los demás casos).

## Relacionadas
- [[ConjAnalyzer]]
- [[DataForge]] · [[Guía DataForge]] (proyecto hermano, mismo stack)
- [[JFlex]] · [[CUP]] · [[Maven]] · [[JavaFX y Scene Builder]]
- [[Gramática libre de contexto (BNF)]]
- [[Cap 3 - Análisis léxico]] · [[Cap 4 - Análisis sintáctico]]
- [[Árbol de sintaxis abstracta (AST)]] · [[Traducción dirigida por la sintaxis]] · [[Atributos sintetizados y heredados]]
- [[Análisis sintáctico ascendente LR]] · [[Derivaciones y árbol de análisis sintáctico]]
- [[Tabla de símbolos]] · [[Manejo de errores (léxicos, sintácticos, semánticos)]]
- [[Caps 8-12 - Panorama (fuera de alcance)]] (la simplificación algebraica es un primo cercano de la optimización de código)
