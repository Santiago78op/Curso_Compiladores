---
tags: [proyecto, guia, lexico, sintactico]
aliases: [guia dataforge, tutorial dataforge, dataforge paso a paso]
fuente: "Enunciado DataForge + Libro del Dragón + construcción real (Etapas 1-6 COMPLETAS: 2026-07-10)"
fecha: 2026-07-10
---

# Guía de elaboración — [[DataForge]]

Guía paso a paso. Cada paso enlaza al **porqué** teórico. Se actualiza con el material REAL a medida que el proyecto se construye (curso guiado: presentación en `Hades\presentacion-dataforge\`, código en `Hades\DataForge\`).

## Estado de construcción

| Etapa | Contenido | Estado |
|---|---|---|
| 0 | Conceptos base (intérprete, [[Token, lexema y patrón]]) | ✅ 2026-07-10 |
| 1 | Tabla de tokens + lexer [[JFlex]] funcionando | ✅ 2026-07-10 |
| 2 | Gramática BNF + parser [[CUP]] validando | ✅ 2026-07-10 |
| 3 | Ejecución ([[Traducción dirigida por la sintaxis]], [[Tabla de símbolos]]) | ✅ 2026-07-10 |
| 4 | GUI [[JavaFX y Scene Builder]] (editor, pestañas, ejecutar, consola) | ✅ 2026-07-10 |
| 5 | Gráficas con JavaFX Charts | ✅ 2026-07-10 |
| 6 | Reportes HTML + recuperación de errores | ✅ 2026-07-10 |

**PROYECTO FUNCIONAL COMPLETO** — cumple todos los requisitos mínimos del enunciado (§7). Falta empaquetado: manuales (usuario/técnico en PDF) y repo de entrega `OLC1_Proyecto1_#Carnet`.

## 1. Requisitos previos (reales)

- **IntelliJ IDEA** Community: importa el `pom.xml` como proyecto [[Maven]] y descarga el **JDK 17** desde el IDE (Project Structure → SDK → Download JDK, Temurin). Maven viene embebido — no se instala nada por terminal.
- Ciclo de trabajo en IDEA: editar `.flex`/`.cup` → ventana Maven → `compile` (regenera) → ▶ Run. Tras tocar el `pom.xml`: 🔄 Reload All Maven Projects.

## 2. Análisis léxico — CONSTRUIDO ✅

### 2.1 Tabla de tokens (la real)

- **24 palabras reservadas** (case insensitive): `program end var arr double char` · `sum res mul div mod` · `media mediana moda varianza max min` · `console print column` · `graphBar graphPie graphLine histogram exec`.
- **11 símbolos**: `:` `::` `<-` `=` `->` `;` `,` `(` `)` `[` `]`.
- **4 con patrón**: `ID` = letra(letra|dígito)* · `ID_ARREGLO` = `@`+ID · `NUMERO` = dígitos(.dígitos)? · `CADENA` = `"[^"]*"`.
- **Se descartan**: comentarios `!…\n` y `<!…!>`, blancos.

### 2.2 Decisiones de diseño tomadas

1. `@datos` es **UN token** `ID_ARREGLO` (el `@` nunca aparece solo → gramática más simple).
2. Los atributos de gráfica (`titulo`, `ejeX`, `values`…) **NO son reservados**: el enunciado §5.9.2 usa `titulo` como variable. Son `ID` y su validez se revisa en ejecución (`Entorno.registrarGrafica`).
3. El enunciado abre `graphBar(` pero ejecuta `EXEC grapBar` (typo sistemático) → el patrón acepta ambas: `"graphbar" | "grapbar"`.
4. Case insensitive con `%ignorecase` — y alcanza también a los **identificadores**: la tabla de símbolos usa claves en minúsculas conservando el nombre original.

### 2.3 El `.flex` real (fragmentos clave)

Archivo: `src/main/jflex/Lexer.flex`, tres secciones separadas por `%%` (ver [[JFlex]]):

```jflex
%class Lexer
%public      // ⚠ SIN esto: "Lexer is not public... cannot be accessed
%unicode     //   from outside package" al usarlo desde otro paquete
%cup
%line
%column
%ignorecase

Id        = {Letra}({Letra}|{Digito})*
IdArreglo = "@"{Id}
Numero    = {Digito}+("."{Digito}+)?
Cadena    = \"[^\"]*\"

%%
{ComentLinea}  { }                  // descartar
"var"          { return symbol(sym.VAR); }        // reservadas ANTES que {Id}
"::"           { return symbol(sym.DOBLE_DOS_PUNTOS); }
{Numero}       { return symbol(sym.NUMERO, Double.valueOf(yytext())); }
{Id}           { return symbol(sym.ID); }
[^]            { /* error léxico con línea/columna */ }
```

Orden crítico: a igual longitud de match gana la **primera regla** — reservadas antes que `{Id}` o no existen ([[Cap 3 - Análisis léxico]]).

### 2.4 Verificación

`TestLexer.java` recorre `next_token()` e imprime la tabla (embrión del reporte §6.1). Verificado: 73 tokens correctos sobre `entradas/ejemplo1.df`, comentarios descartados, `%ignorecase` y longest match confirmados.

## 3. Análisis sintáctico — CONSTRUIDO ✅

### 3.1 La gramática

- **BNF limpia** (entregable §8, NO copia del .cup): `docs/gramatica.txt`. Núcleo:

```
<inicio>      ::= "PROGRAM" <instrucciones> "END" "PROGRAM"
<instruccion> ::= <declaracion-var> | <declaracion-arr> | <imprimir>
                | <imprimir-columna> | <grafica>
<expresion>   ::= NUMERO | CADENA | ID | <aritmetica> | <estadistica>
<aritmetica>  ::= <operacion> "(" <expresion> "," <expresion> ")"
```

- La **recursión mutua** `<aritmetica>` ↔ `<expresion>` da la anidación infinita con paréntesis balanceados (lo que las regex no pueden: no cuentan — [[Cap 4 - Análisis sintáctico]] §4.1).
- Operaciones como funciones (`SUM(a,b)`) → sin ambigüedad, **sin tabla de [[Ambigüedad, precedencia y asociatividad|precedencia]]**.
- El token `END` se desambigua por contexto en la gramática (`end;` vs `END PROGRAM`), no en el lexer.
- `valor_attr ::= expr | arreglo` — los atributos de gráfica aceptan `@nombre` además de `[...]`.

### 3.2 El `.cup` real

Archivo: `src/main/cup/parser.cup` (estructura en [[CUP]]). Claves:
- Las declaraciones `terminal` **generan la clase `sym`** → se BORRÓ el `sym.java` temporal de la Etapa 1; el `.flex` no cambió ni una letra.
- `syntax_error(Symbol s)` registra el error en el Entorno con `s.left+1`/`s.right+1`.

### 3.3 Verificación

Verificado 2026-07-10: `ejemplo1.df` y `ejemplo2.df` → `[OK]`; quitando un `;` → `Error sintáctico: no se esperaba 'arr' (línea 5, columna 5)`. El error apunta al token donde se DESCUBRE el problema, no a la causa (naturaleza [[Análisis sintáctico ascendente LR|LALR]]).

## 4. Ejecución — CONSTRUIDO ✅

### 4.1 El diseño (Cap. 5: [[Traducción dirigida por la sintaxis]])

**Evaluación directa en las acciones de CUP** (sin AST — viable porque DataForge no tiene control de flujo). Gramática **S-atribuida**: cada producción sintetiza su valor con `{: RESULT = … :}` ([[Atributos sintetizados y heredados]]); el orden de evaluación lo garantiza el orden de reducciones del parser LALR (interno antes que externo).

```java
aritmetica ::= op_arit:op PAR_IZQ expr:a COMA expr:b PAR_DER
               {: RESULT = Operaciones.aritmetica(op, a, b,
                               parser.entorno, opleft, opright); :} ;
```

- No terminales **tipados**: `non terminal Object expr;` · `non terminal ArrayList lista_expr;` (tipos RAW — generics en `non terminal` son terreno frágil en CUP).
- `parser code {: public Entorno entorno = new Entorno(); :}` — las acciones acceden vía `parser.entorno`.
- `Xleft`/`Xright` de cada etiqueta dan línea/columna para los errores.

### 4.2 El paquete `dataforge.interprete`

| Clase | Responsabilidad |
|---|---|
| `Entorno` | [[Tabla de símbolos]] (`LinkedHashMap`, claves lowercase), consola (`StringBuilder`), errores, gráficas registradas |
| `Operaciones` | SUM/RES/MUL/DIV/MOD + Media/Mediana/Moda/Varianza/Max/Min con chequeo de tipos (double only, §5.7) |
| `Simbolo` | ficha del reporte §6.3: nombre, categoría, tipo, valor, línea, columna |
| `Grafica` | tipo + atributos resueltos (「la última gana」 vía put repetido en LinkedHashMap; solo se registra con EXEC) |

**Errores semánticos** ([[Manejo de errores (léxicos, sintácticos, semánticos)]]): variable no declarada, tipos incompatibles, división entre cero, redeclaración, atributo de gráfica inválido. **No detienen la ejecución** (el reporte §6.2 quiere todos) y la **propagación por null** evita cascadas: expresión fallida → null → las operaciones que reciben null callan.

### 4.3 Verificación (2026-07-10)

`TestInterprete.java` sobre los 3 ejemplos: `ejemplo1.df` → media 3.5 + columna + 4 símbolos; `ejemplo2.df` → 3 gráficas registradas; `ejemplo3_errores.df` → 5 errores semánticos con línea/columna y la ejecución sobrevive. Formato de números como el enunciado: `15.0` se muestra `15`, `15.7` queda `15.7`.

### 4.4 `pom.xml` (versión Etapa 3)

`jflex-maven-plugin:1.9.1` + `cup-maven-plugin:11b-20160615-3` + runtime `java-cup-runtime:11b-20160615` + `exec-maven-plugin` (mainClass `dataforge.TestInterprete`). Falta agregar: JavaFX (Etapa 4).

## 5. GUI — CONSTRUIDA ✅ (2026-07-10)

- `gui/EditorApp.java`: [[JavaFX y Scene Builder|JavaFX]] **por código** (sin FXML — decisión: 6 controles, árbol visible). Scene graph: BorderPane → HBox de botones (Nuevo/Abrir/Guardar/▶ Ejecutar) + SplitPane vertical (TabPane editor / TextArea consola `setEditable(false)`).
- Pestañas: `Tab.userData` = File asociado (null → «Guardar como»); FileChooser con filtro `*.df`. Cumple §4.2-4.4 y 4.6 del enunciado.
- `interprete/Interprete.java`: fachada `String → Entorno` con **StringReader** (mismo pipeline, otra fuente). **Entorno fresco por ejecución** (§6: reportes solo del último análisis). `ultimoEntorno` guardado para Etapas 5-6.
- `gui/Lanzador.java`: main en clase que NO extiende Application → esquiva «JavaFX runtime components are missing» con classpath plano. **En IDEA se corre Lanzador**, o `mvn clean javafx:run`.
- pom: + `org.openjfx:javafx-controls:21.0.4` y `javafx-maven-plugin:0.0.8`. Los WARNING de arranque (unnamed module, native access, Unsafe) son esperados e inofensivos.

## 5b. Gráficas — CONSTRUIDAS ✅ (2026-07-10)

- `gui/Graficador.java`: cada `Grafica` → Chart → Stage propio. Mapeo: graphBar→`BarChart`, graphPie→`PieChart`, graphLine→`LineChart`, Histogram→`BarChart` de frecuencias. Los Charts vienen en `javafx-controls` (sin dependencias nuevas).
- `Entorno.validarGrafica`: chequeo semántico ANTES de registrar — atributos requeridos por tipo, tipo de elementos, tamaños coherentes (ejeX vs ejeY, label vs values). *Validar temprano, dibujar confiado* (los casts del Graficador son directos).
- Histograma (5.10.3): `Operaciones.frecuencias()` (ordena y cuenta) + tabla Valor/Frec./Acum./Rel. escrita por el Entorno en consola (semántica, existe sin GUI). Formato de tabla = decisión propia (la imagen del PDF no se convirtió). El `values::char[]` del enunciado se salda a favor de doubles.
- Se dibuja DESPUÉS de ejecutar, desde el hilo FX del botón: separación de capas (el intérprete no conoce la GUI) + regla del hilo único de JavaFX.
- Verificado: `ejemplo2.df` (con graphLine agregado) abre las 4 ventanas.

## 5c. Reportes HTML — CONSTRUIDOS ✅ (2026-07-10)

- `reportes/Reportes.java`: genera `tokens.html`, `errores.html`, `simbolos.html` (plantilla con CSS embebido, autocontenida). Nombres de token por **reflexión sobre `sym`**; lexemas **escapados** (`<-` rompería la tabla del navegador). Botón **Reportes** en la GUI → `getHostServices().showDocument`.
- **Tokens (§6.1)**: el lexer los registra al nacer con `yytext()` (lexema ORIGINAL — `1` se reporta `1`, no `1.0`), vía campo público `entorno` que setea el `Interprete`. Errores léxicos también van al entorno (adiós stderr).
- **Errores (§6.2)**: `RegistroError` estructurado (tipo/descripción/línea/columna, base 1; `toString()` compatible con la consola). Las 3 familias en una lista.
- **Recuperación en modo pánico** ([[Manejo de errores (léxicos, sintácticos, semánticos)]], Dragón §4.8.3): producción `instruccion ::= error PUNTO_COMA` — CUP poda la pila hasta poder desplazar `error` y descarta tokens hasta el `;` (punto de sincronización natural). El reporte acumula TODOS los errores sintácticos; la instrucción rota se sacrifica entera.
- **Símbolos (§6.3)**: formato del enunciado vía `Entorno.valorReporte()` — cadenas con comillas (`"Hola Mundo"`), arreglos sin los `.0` de Java (`[1, 2.5, 7]`). ⚠ Formato de reporte ≠ formato de consola (donde las cadenas van sin comillas): son DOS formateadores explícitos.
- Verificado con `ejemplo4_mixto.df`: los 3 tipos de error en un programa, todos reportados con línea/columna, y la ejecución sobrevive (`precio` declarada pese al `$`; `roto` sacrificada por el pánico).

## 6. Casos de prueba

- `entradas/ejemplo1.df`: variables, arreglo con `SUM` anidada, `Media`, `console::print/column`, ambos comentarios.
- `entradas/ejemplo2.df`: las 4 gráficas con `EXEC grapX`, `DIV(SUM(Max,Min),2)` anidada, atributo con `@arreglo`.
- `entradas/ejemplo3_errores.df`: 5 errores semánticos a propósito (léxica/sintácticamente perfecto).
- `entradas/ejemplo4_mixto.df`: los 3 tipos de error en un solo programa — prueba de las 3 recuperaciones y de los reportes.
- Negativos verificados: quitar `;` / `<-` → error sintáctico; carácter `$` → error léxico.

## 7. Errores comunes (los que ya nos pasaron)

- **Olvidar `%public`** → `'Lexer' is not public ...` (JFlex genera package-private por defecto).
- Olvidar `%cup` en el `.flex` → el lexer no produce `Symbol` para [[CUP]].
- Reglas en mal orden → un patrón goloso se traga a otro (*longest match* / prioridad).
- Editar código generado en vez del `.flex`/`.cup` → se pierde en el siguiente `compile`.
- `NUMERO` convertido con `Double.valueOf` imprime `1.0` para el lexema `1` → guardar también `yytext()` para el reporte de tokens.
- `Couldn't repair and continue parse` = CUP sin recuperación de errores; para acumular todos, terminal `error` (etapa de reportes).
- Declaraciones `non terminal` van TODAS antes de las producciones; tipos raw (`ArrayList`), no genéricos ni arrays.
- IDEA subraya `Parser`/`sym` en rojo → Maven `compile` + 🔄 Reload.

## Relacionadas
- [[DataForge]]
- [[JFlex]] · [[CUP]] · [[Maven]]
- [[Cap 3 - Análisis léxico]] · [[Cap 4 - Análisis sintáctico]]
- [[Traducción dirigida por la sintaxis]] · [[Atributos sintetizados y heredados]]
- [[Análisis sintáctico ascendente LR]]
