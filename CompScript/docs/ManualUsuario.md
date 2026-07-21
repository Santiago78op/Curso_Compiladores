# Manual de Usuario — CompScript

Universidad de San Carlos de Guatemala
Facultad de Ingeniería
Escuela de Ciencias y Sistemas
Organización de Lenguajes y Compiladores 1
Vacaciones de Diciembre 2024

## 1. Introducción

CompScript es un intérprete de un lenguaje de programación de propósito educativo, desarrollado como proyecto del curso de Organización de Lenguajes y Compiladores 1. El presente manual describe el uso del entorno de trabajo (editor gráfico) desde la perspectiva del usuario final: cómo crear, abrir y guardar archivos fuente, cómo ejecutar el análisis e interpretación de un programa, y cómo interpretar los reportes que el sistema genera.

Este documento no cubre la implementación interna del intérprete; para ello se debe consultar el Manual Técnico.

## 2. Instalación y ejecución del entorno

El entorno de trabajo se distribuye como un proyecto Maven que produce, mediante los complementos configurados en `pom.xml`, tanto el analizador léxico como el sintáctico y la interfaz gráfica.

Para iniciar la aplicación:

- Desde un entorno de desarrollo (IntelliJ IDEA): ejecutar la clase `compscript.gui.Lanzador`.
- Desde línea de comandos: `mvn clean javafx:run`.

[CAPTURA DE PANTALLA: ventana principal del editor recién iniciada, mostrando la pestaña con la plantilla de código por defecto y el área de consola vacía en la parte inferior]

## 3. El editor de trabajo

### 3.1 Descripción general

La ventana principal se organiza en tres zonas:

1. Una barra superior de botones: **Nuevo**, **Abrir**, **Guardar**, **▶ Ejecutar**, **Reportes** y **Ver AST**.
2. Un área central con pestañas, donde cada pestaña corresponde a un archivo fuente abierto o en edición.
3. Un área de consola en la parte inferior, de solo lectura, donde se muestran los resultados de la ejecución.

[CAPTURA DE PANTALLA: ventana principal con las tres zonas señaladas o enmarcadas]

### 3.2 Nuevo archivo

El botón **Nuevo** crea una pestaña adicional con una plantilla de código inicial. Esta acción no afecta el contenido de la consola: la consola conserva el resultado de la última ejecución hasta que se presiona **▶ Ejecutar** nuevamente. Cada pestaña nueva se numera automáticamente (por ejemplo, `nuevo1.cs`, `nuevo2.cs`).

[CAPTURA DE PANTALLA: resultado de presionar Nuevo, con una pestaña adicional visible]

### 3.3 Abrir archivo

El botón **Abrir** despliega un explorador de archivos filtrado por la extensión `.cs`, propia del lenguaje CompScript. Al seleccionar un archivo, su contenido se carga en una nueva pestaña, identificada con el nombre del archivo.

[CAPTURA DE PANTALLA: cuadro de diálogo de selección de archivo, mostrando el filtro *.cs y alguno de los ejemplos de la carpeta entradas/]

### 3.4 Guardar archivo

El botón **Guardar** almacena el contenido de la pestaña activa en disco:

- Si la pestaña corresponde a un archivo ya abierto o previamente guardado, se sobrescribe ese mismo archivo.
- Si la pestaña es nueva (no asociada a ningún archivo), se despliega un cuadro de diálogo de tipo "Guardar como", donde debe indicarse el nombre y la ubicación. El nombre de la pestaña se actualiza al del archivo guardado.

[CAPTURA DE PANTALLA: cuadro de diálogo "Guardar como" con el filtro *.cs]

## 4. Ejecución del programa

El botón **▶ Ejecutar** invoca al intérprete sobre el contenido completo de la pestaña activa. El proceso realiza, en orden, el análisis léxico, el análisis sintáctico (construcción del árbol de sintaxis abstracta) y el análisis semántico junto con la ejecución de las instrucciones.

Cada ejecución parte de un estado limpio: no conserva variables, funciones ni errores de ejecuciones anteriores. Los reportes disponibles después de ejecutar corresponden siempre al último archivo ejecutado.

Al finalizar, el área de consola muestra:

- La salida generada por las instrucciones `console.log(...)` del programa, en el orden en que se ejecutaron.
- A continuación, si existieron errores léxicos, sintácticos o semánticos, un bloque separado con el conteo total y el detalle de cada uno (tipo de error, descripción, línea y columna).

[CAPTURA DE PANTALLA: consola tras ejecutar un programa sin errores, mostrando únicamente la salida de console.log]

[CAPTURA DE PANTALLA: consola tras ejecutar un programa con errores léxicos, sintácticos y semánticos, mostrando el bloque de errores al final]

### 4.1 Recuperación ante errores

El entorno está diseñado para recuperarse de los tres tipos de error sin detener el proceso de análisis:

- **Error léxico**: el carácter no reconocido se descarta y el análisis continúa con el siguiente carácter.
- **Error sintáctico**: el analizador descarta los símbolos siguientes hasta encontrar un punto y coma (`;`), y retoma el análisis desde la instrucción siguiente. Esto permite que un único error de sintaxis no impida detectar el resto de errores del archivo.
- **Error semántico**: a diferencia de los dos anteriores, un error semántico detectado durante la ejecución (por ejemplo, una operación entre tipos incompatibles) queda registrado en el reporte de errores y **finaliza la ejecución del programa de forma ordenada**. Las instrucciones posteriores al punto del error no se ejecutan.

## 5. Reportes

El botón **Reportes** genera y abre en el navegador predeterminado del sistema cuatro páginas HTML correspondientes al último programa ejecutado. El botón **Ver AST** despliega, adicionalmente, una ventana con el árbol de sintaxis abstracta en forma de árbol navegable, sin necesidad de abrir el navegador.

Si no se ha ejecutado ningún programa todavía, ambos botones muestran un mensaje indicándolo en la consola.

### 5.1 Reporte de tokens

Enumera, en el orden en que fueron reconocidos, todos los tokens (lexemas) identificados por el analizador léxico durante la última ejecución. Cada fila indica el número consecutivo, el lexema tal como aparece en el código fuente, el nombre del tipo de token, la línea y la columna donde fue encontrado.

[CAPTURA DE PANTALLA: reporte de tokens abierto en el navegador]

### 5.2 Reporte de errores

Enumera todos los errores detectados durante el análisis léxico, sintáctico y semántico del último programa ejecutado. Cada fila indica el tipo de error (Léxico, Sintáctico o Semántico), su descripción, y la línea y columna donde ocurrió. Si el programa no presentó errores, el reporte lo indica explícitamente.

[CAPTURA DE PANTALLA: reporte de errores abierto en el navegador, mostrando al menos un error de cada tipo]

### 5.3 Reporte del árbol de sintaxis abstracta (AST)

Representa gráficamente, en forma de árbol, la estructura sintáctica reconocida del último programa analizado. Cada nodo describe la construcción del lenguaje que representa (por ejemplo, una declaración, una operación aritmética, un ciclo) y sus nodos hijos corresponden a las subexpresiones o subinstrucciones que lo componen.

Este reporte está disponible en dos formas equivalentes: una página HTML autocontenida (accesible mediante el botón **Reportes**) y una ventana con árbol interactivo dentro de la misma aplicación (accesible mediante el botón **Ver AST**).

[CAPTURA DE PANTALLA: ventana "Ver AST" mostrando el árbol expandido de un programa de ejemplo]

[CAPTURA DE PANTALLA: reporte de AST en HTML, abierto en el navegador]

### 5.4 Reporte de tabla de símbolos

Enumera todas las variables, vectores, listas y structs declarados durante la ejecución del programa, indicando su identificador, su categoría (Variable, Vector, Lista o Struct), su tipo de dato, el ámbito (entorno) en el que fueron declarados, su valor al momento de generarse el reporte, y la línea y columna de su declaración.

Cuando el programa incluye funciones recursivas, es normal observar múltiples entradas para un mismo identificador: cada invocación de la función crea su propio ámbito, y el reporte conserva el registro histórico completo de todas las declaraciones realizadas durante la ejecución, no únicamente el estado final.

[CAPTURA DE PANTALLA: reporte de tabla de símbolos abierto en el navegador, incluyendo al menos una variable, un vector y un struct]

## 6. Área de consola

El área de consola, ubicada en la parte inferior de la ventana principal, es de solo lectura y concentra toda la salida generada por el programa en ejecución: los mensajes producidos por la instrucción `console.log(...)` del lenguaje, y, cuando corresponde, el listado de errores léxicos, sintácticos y semánticos detectados.

## 7. Notas de uso

- La extensión de archivo propia del lenguaje es `.cs`.
- El lenguaje es *case insensitive*: no distingue mayúsculas de minúsculas, tanto en palabras reservadas como en identificadores.
- Se recomienda guardar el archivo de trabajo antes de cerrar la aplicación, dado que no existe una función de recuperación automática de contenido no guardado.
