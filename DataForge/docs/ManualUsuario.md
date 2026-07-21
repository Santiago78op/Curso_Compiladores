# Manual de Usuario — DataForge

**Universidad de San Carlos de Guatemala**
**Facultad de Ingeniería — Escuela de Ciencias y Sistemas**
**Organización de Lenguajes y Compiladores 1 — Proyecto 1**

---

## 1. Introducción

DataForge es un intérprete con interfaz gráfica para un lenguaje de dominio específico orientado al análisis de datos: permite declarar variables y arreglos, realizar operaciones aritméticas y estadísticas, imprimir resultados en consola y generar gráficas (de barras, circulares, de línea e histogramas) a partir de colecciones de valores.

Este manual describe, desde la perspectiva de un usuario final, cómo instalar el entorno de desarrollo, ejecutar la aplicación y utilizar cada una de sus funcionalidades: el editor de código, la ejecución de programas, la consola de salida, la generación de gráficas y los reportes HTML.

## 2. Requisitos previos

- **JDK 17** o superior (el proyecto fue compilado y probado con JDK 25; el `pom.xml` fija el nivel de lenguaje en 17).
- **IntelliJ IDEA** (Community o Ultimate). Maven viene incluido en el IDE — no es necesario instalarlo por separado.
- Conexión a internet la primera vez que se abre el proyecto, para que Maven descargue las dependencias (`java-cup-runtime`, `javafx-controls`) y los plugins de generación (`jflex-maven-plugin`, `cup-maven-plugin`, `javafx-maven-plugin`).

## 3. Instalación y arranque

1. Abrir IntelliJ IDEA y seleccionar **Open** sobre la carpeta `DataForge/` (la que contiene el `pom.xml`).
2. IntelliJ detecta el proyecto como Maven y descarga las dependencias automáticamente. Si no lo hace, usar el panel lateral **Maven → Reload All Maven Projects**.
3. En **File → Project Structure → SDK**, verificar que el proyecto tenga asignado un JDK 17 o superior. Si no existe uno instalado, IntelliJ permite descargarlo desde el mismo diálogo (**Download JDK…**, distribución Temurin recomendada).
4. Compilar el proyecto: panel **Maven → dataforge → Lifecycle → compile** (doble clic). Este paso ejecuta JFlex y CUP sobre los archivos `Lexer.flex` y `parser.cup`, generando el analizador léxico y sintáctico antes de compilar el resto del código Java.
5. Ejecutar la aplicación: en el árbol de proyecto, abrir `src/main/java/dataforge/gui/Lanzador.java` y presionar el botón ▶ **Run** junto a la clase `Lanzador`.

   > **Importante:** la aplicación se debe ejecutar sobre la clase `Lanzador`, **no** directamente sobre `EditorApp`. Ejecutar `EditorApp` directamente produce el error "JavaFX runtime components are missing", porque el lanzador estándar de Java bloquea las clases que extienden `Application` cuando JavaFX no está en el *module-path*. `Lanzador` es una clase auxiliar que invoca a `EditorApp.main(...)` y evita ese chequeo.
   >
   > Alternativa por terminal (si se cuenta con Maven instalado): `mvn clean javafx:run`.

`[CAPTURA DE PANTALLA: ventana principal de DataForge recién abierta, con una pestaña nueva vacía y la consola de salida debajo]`

## 4. El entorno de trabajo

La ventana principal se divide en tres zonas:

- Una **barra de herramientas** superior con los botones Nuevo, Abrir, Guardar, ▶ Ejecutar y Reportes.
- Un **editor con pestañas**, donde cada pestaña corresponde a un archivo `.df` abierto o a un archivo nuevo sin guardar.
- Una **consola de salida** en la parte inferior, de solo lectura, donde se muestran los resultados de la ejecución, los errores detectados y un resumen de las gráficas generadas.

`[CAPTURA DE PANTALLA: ventana principal señalando cada una de las tres zonas — barra de botones, editor de pestañas, consola]`

### 4.1 Nuevo archivo

El botón **Nuevo** crea una pestaña en blanco con la plantilla mínima de un programa DataForge:

```dataforge
PROGRAM


END PROGRAM
```

Cada pestaña nueva se numera automáticamente (`nuevo1.df`, `nuevo2.df`, …) hasta que se guarda con un nombre definitivo.

### 4.2 Abrir archivo

El botón **Abrir** despliega un selector de archivos filtrado por la extensión `.df`. Al elegir un archivo, su contenido se carga en una nueva pestaña cuyo título es el nombre del archivo. Si el archivo no puede leerse, el mensaje de error se muestra en la consola.

`[CAPTURA DE PANTALLA: diálogo "Abrir archivo DataForge" mostrando los archivos de ejemplo con extensión .df]`

### 4.3 Guardar archivo

El botón **Guardar** almacena el contenido de la pestaña actualmente seleccionada:

- Si la pestaña corresponde a un archivo ya abierto desde disco, se sobrescribe ese mismo archivo.
- Si la pestaña es nueva (no tiene un archivo asociado), se abre un diálogo **Guardar como**, donde se debe indicar el nombre y la ubicación; a partir de ese momento la pestaña queda asociada a ese archivo.

### 4.4 Cerrar pestañas

Cada pestaña puede cerrarse en cualquier momento haciendo clic en su botón de cierre. Si el contenido no se ha guardado, se descarta sin advertencia adicional — se recomienda guardar antes de cerrar una pestaña con cambios importantes.

### 4.5 Ejecutar un programa

El botón **▶ Ejecutar** envía el contenido de la pestaña actualmente seleccionada al intérprete. El proceso realiza, en orden:

1. Análisis léxico y sintáctico del texto.
2. Ejecución de las instrucciones válidas (declaraciones, impresiones, cálculos, registro de gráficas).
3. Escritura del resultado en la consola: primero la salida generada por las instrucciones `console::print` y `console::column`, y luego, si existieron, un resumen de los errores encontrados y de las gráficas mostradas.
4. Si el programa registró gráficas con `EXEC`, cada una se abre automáticamente en su propia ventana.

Cada clic en **Ejecutar** reinicia el estado del intérprete: las variables, arreglos, errores y gráficas de una ejecución anterior no se conservan en la siguiente. Esto es intencional, ya que los reportes (sección 6) siempre reflejan **únicamente el último análisis realizado**.

`[CAPTURA DE PANTALLA: pestaña con el contenido de ejemplo1.df y la consola mostrando la salida "la media es:, 3.5" y la tabla de columna "Datos"]`

**Ejemplo de programa** (`entradas/ejemplo1.df`, incluido en el proyecto):

```dataforge
PROGRAM
    var:double:: numero <- 2.5 end;
    var:char[]:: saludo <- "hola mundo" end;

    arr:double::@datos <- [1, 2.5, SUM(3, 4)] end;

    var:double:: m <- Media(@datos) end;

    console::print = "la media es:", m end;
    console::column = "Datos" -> @datos end;
END PROGRAM
```

Al ejecutarlo, la consola muestra:

```
la media es:, 3.5
--------------
Datos
--------------
1
2.5
7
```

### 4.6 Gráficas

DataForge soporta cuatro tipos de gráfica: barras (`graphBar`), circular (`graphPie`), de línea (`graphLine`) e histograma (`Histogram`). Cada bloque de gráfica define sus atributos (título, ejes, valores, etiquetas) y solamente se dibuja si el bloque incluye la instrucción `EXEC` correspondiente.

Al ejecutar un programa con bloques de gráfica válidos, cada gráfica se abre en una ventana independiente, superpuesta a la ventana principal.

`[CAPTURA DE PANTALLA: las cuatro ventanas de gráfica abiertas simultáneamente al ejecutar entradas/ejemplo2.df — barras, pastel, línea e histograma]`

Si un bloque de gráfica tiene un atributo faltante, de tipo incorrecto, o con listas de tamaño incompatible (por ejemplo, `ejeX` y `ejeY` con distinta cantidad de elementos), la gráfica no se dibuja y el problema se reporta como un error semántico en la consola y en el reporte de errores.

### 4.7 Reportes

El botón **Reportes** genera tres archivos HTML a partir del **último programa ejecutado** y los abre automáticamente en el navegador predeterminado:

- **`tokens.html`** — todos los tokens reconocidos por el analizador léxico, en el orden en que aparecieron, con su lexema original, el nombre del token y su posición (línea, columna).
- **`errores.html`** — todos los errores detectados (léxicos, sintácticos y semánticos), con su tipo, descripción y posición.
- **`simbolos.html`** — la tabla de símbolos resultante: cada variable o arreglo declarado, con su categoría, tipo, valor y la posición donde fue declarado.

Los tres archivos se generan en la carpeta `reportes/` del proyecto y son autocontenidos (CSS embebido), por lo que pueden abrirse sin conexión a internet.

Si se presiona **Reportes** antes de ejecutar algún programa, la consola indica que primero debe ejecutarse un programa.

`[CAPTURA DE PANTALLA: reporte tokens.html abierto en el navegador, mostrando la tabla de tokens de un programa de ejemplo]`

`[CAPTURA DE PANTALLA: reporte errores.html abierto en el navegador, mostrando los errores léxico/sintáctico/semántico de entradas/ejemplo4_mixto.df]`

`[CAPTURA DE PANTALLA: reporte simbolos.html abierto en el navegador, mostrando la tabla de símbolos de un programa de ejemplo]`

## 5. Manejo de errores

DataForge distingue tres tipos de error, y en los tres casos la ejecución **continúa** en vez de detenerse:

- **Léxico**: un carácter no reconocido por el lenguaje (por ejemplo, `$`). El carácter se descarta y el análisis continúa desde el siguiente.
- **Sintáctico**: una instrucción mal formada (por ejemplo, falta un `;` o un `<-`). El analizador descarta tokens hasta encontrar el siguiente `;` y continúa con la instrucción siguiente; la instrucción defectuosa se pierde por completo.
- **Semántico**: un problema detectado durante la ejecución (variable no declarada, tipos incompatibles, división entre cero, redeclaración de una variable, atributo de gráfica inválido). La expresión afectada se resuelve como ausente y las instrucciones que dependen de ella no producen una salida adicional, pero el resto del programa continúa normalmente.

Todos los errores encontrados, sin importar su tipo, se acumulan y se muestran juntos al final de la consola y en el reporte de errores — el programa completo se reporta en un único análisis.

`[CAPTURA DE PANTALLA: consola mostrando la ejecución de entradas/ejemplo4_mixto.df con los tres tipos de error listados y el mensaje "final:, 10, precio:, 99" probando que la ejecución sobrevivió]`

## 6. Referencia rápida del lenguaje

Esta sección resume la sintaxis del lenguaje DataForge para consulta rápida del usuario. La gramática formal completa, en formato BNF, se encuentra en el archivo `docs/gramatica.txt` que acompaña este manual.

- El lenguaje **no distingue mayúsculas de minúsculas** (`case insensitive`), incluyendo los nombres de variables.
- Todo programa debe estar encerrado entre `PROGRAM` y `END PROGRAM`.
- Comentarios de una línea: `! texto`. Comentarios multilínea: `<! texto !>`.
- Tipos de dato: `double` (numérico) y `char[]` (cadena de texto).
- Declarar una variable: `var:double:: nombre <- expresión end;`
- Declarar un arreglo: `arr:double::@nombre <- [v1, v2, ...] end;` (el nombre de un arreglo siempre lleva el prefijo `@`).
- Operaciones aritméticas (sobre `double`): `SUM`, `RES`, `MUL`, `DIV`, `MOD` — se usan como funciones: `SUM(a, b)`.
- Operaciones estadísticas (sobre arreglos de `double`): `Media`, `Mediana`, `Moda`, `Varianza`, `Max`, `Min` — por ejemplo, `Media(@datos)`.
- Imprimir en consola: `console::print = expr1, expr2, ... end;`
- Imprimir un arreglo como columna: `console::column = título -> @arreglo end;`
- Bloque de gráfica: `graphBar( atributo1::tipo = valor end; ... EXEC graphBar end; ) end;` (análogo para `graphPie`, `graphLine`, `Histogram`).

Los cuatro archivos de prueba incluidos en `entradas/` (`ejemplo1.df`, `ejemplo2.df`, `ejemplo3_errores.df`, `ejemplo4_mixto.df`) sirven como referencia práctica de la sintaxis y pueden abrirse directamente desde el botón **Abrir**.
