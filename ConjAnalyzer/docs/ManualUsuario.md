# Manual de Usuario — ConjAnalyzer

Universidad de San Carlos de Guatemala
Facultad de Ingeniería — Escuela de Ciencias y Sistemas
Organización de Lenguajes y Compiladores 1
Proyecto 1 — ConjAnalyzer

## Índice

1. Introducción
2. Requisitos para la ejecución
3. El entorno de trabajo (editor)
4. Funcionalidades del editor
5. El lenguaje ConjAnalyzer (resumen de uso)
6. Ejecución del análisis
7. Reportes generados
   7.1 Reporte de tokens
   7.2 Reporte de errores
   7.3 Reporte de conjuntos y operaciones
   7.4 Diagrama de Venn
   7.5 JSON de simplificación
8. Ejemplo completo guiado
9. Solución de problemas frecuentes

---

## 1. Introducción

ConjAnalyzer es una herramienta desarrollada para el curso de Organización de Lenguajes y Compiladores 1 que permite definir conjuntos, realizar operaciones entre ellos (unión, intersección, diferencia y complemento) y evaluar la pertenencia de elementos a los conjuntos resultantes. El propósito de la herramienta es que estudiantes del curso de Matemática para Computación 1 puedan verificar de forma práctica las respuestas de sus tareas y exámenes de teoría de conjuntos.

El sistema implementa un analizador léxico y un analizador sintáctico construidos con las herramientas JFlex y CUP, respectivamente, y presenta una interfaz gráfica desarrollada en JavaFX que permite editar el código fuente, ejecutarlo y visualizar los resultados mediante diagramas de Venn y reportes en formato HTML y JSON.

## 2. Requisitos para la ejecución

- Java 17 o superior instalado en el equipo (o el entorno de ejecución empaquetado con el archivo `.jar`, según cómo se distribuya la entrega).
- El archivo ejecutable `ConjAnalyzer.jar` (o el comando de ejecución equivalente provisto por el desarrollador).
- No se requiere conexión a internet: el programa y sus reportes son completamente autocontenidos.

Para iniciar la aplicación se ejecuta el archivo `.jar` correspondiente. Al abrirse aparecerá la ventana principal del editor.

`[CAPTURA DE PANTALLA: ventana principal de ConjAnalyzer recién abierta, mostrando la pestaña de ejemplo precargada, la consola vacía y el panel de diagramas de Venn en su estado inicial]`

## 3. El entorno de trabajo (editor)

La ventana principal se divide en tres zonas:

- **Barra de herramientas** (parte superior): contiene los botones Nuevo, Abrir, Guardar, ▶ Ejecutar y Reportes.
- **Panel izquierdo**: dividido verticalmente entre el editor de código fuente (con soporte de pestañas para trabajar con varios archivos a la vez) y la consola de salida, que muestra los resultados de la ejecución y no puede ser editada manualmente por el usuario.
- **Panel derecho**: muestra los diagramas de Venn generados para cada operación definida en el último análisis ejecutado, con controles de navegación para desplazarse entre ellos.

`[CAPTURA DE PANTALLA: ventana principal señalando con anotaciones las tres zonas — barra de herramientas, editor + consola, panel de Venn]`

## 4. Funcionalidades del editor

| Botón | Acción |
|---|---|
| **Nuevo** | Crea una nueva pestaña de edición vacía, con la extensión de archivo `.ca`. |
| **Abrir** | Muestra un cuadro de diálogo para seleccionar un archivo `.ca` existente y cargarlo en una nueva pestaña. |
| **Guardar** | Guarda el contenido de la pestaña activa. Si la pestaña no tiene un archivo asociado todavía, se solicita el nombre y la ubicación mediante un cuadro de "Guardar como". |
| **▶ Ejecutar** | Envía el contenido de la pestaña activa al intérprete: se realiza el análisis léxico, sintáctico y la evaluación de las instrucciones. El resultado se muestra en la consola y se actualiza el panel de diagramas de Venn. |
| **Reportes** | Genera los reportes HTML (tokens, errores y conjuntos/operaciones) y el JSON de simplificación correspondientes al último análisis ejecutado, y los abre en el navegador predeterminado del sistema. |

`[CAPTURA DE PANTALLA: cuadro de diálogo "Abrir archivo ConjAnalyzer" mostrando el filtro de extensión *.ca]`

`[CAPTURA DE PANTALLA: cuadro de diálogo "Guardar como" al presionar Guardar sobre una pestaña nueva sin archivo asociado]`

Es importante notar que cada vez que se presiona **Ejecutar**, el análisis se realiza sobre un entorno de ejecución completamente nuevo: los reportes y diagramas mostrados corresponden siempre al último análisis realizado, y no se conserva información de ejecuciones anteriores.

## 5. El lenguaje ConjAnalyzer (resumen de uso)

Todo programa debe encerrarse entre llaves `{ }`. El lenguaje es sensible a mayúsculas y minúsculas (`conjuntoA` es distinto de `conjuntoa`), y admite dos tipos de comentario: de línea, iniciado con `#`, y multilínea, delimitado por `<!` y `!>`.

### 5.1 Definición de un conjunto

```
CONJ : nombre -> notacion;
```

La notación puede ser un **rango** (`a~z`, todos los caracteres entre `a` y `z` inclusive) o una **lista de elementos** separados por coma (`m, j, d, 1`). Cada elemento debe ser un único carácter dentro del universo (símbolos ASCII entre `!` y `~`, es decir, códigos 33 a 126).

### 5.2 Definición de una operación

```
OPERA : nombre -> operacion;
```

Las operaciones se escriben en **notación prefija (polaca)**: el operador se escribe antes que sus operandos. Los operadores disponibles son:

| Operador | Significado |
|---|---|
| `U {A} {B}` | Unión |
| `& {A} {B}` | Intersección |
| `- {A} {B}` | Diferencia |
| `^ {A}` | Complemento |

Las operaciones se pueden anidar sin límite, por ejemplo `U U {A} {B} {C}` equivale a `(A U B) U C`.

### 5.3 Evaluación de pertenencia

```
EVALUAR ( {e1, e2, ...} , nombre_operacion );
```

Evalúa, elemento por elemento, si cada uno pertenece al conjunto resultante de la operación indicada. El resultado se muestra en la consola indicando "exitoso" o "fallo" para cada elemento.

## 6. Ejecución del análisis

Para ejecutar un programa:

1. Escribir o cargar el código fuente en una pestaña del editor.
2. Presionar el botón **▶ Ejecutar**.
3. Observar el resultado en la consola de salida (parte inferior del panel izquierdo).
4. Revisar, si corresponde, los diagramas de Venn generados en el panel derecho.

`[CAPTURA DE PANTALLA: consola de salida mostrando el resultado de un EVALUAR, con las líneas "elemento -> exitoso" / "elemento -> fallo"]`

Si el código fuente contiene errores léxicos, sintácticos o semánticos, la ejecución **no se detiene**: el sistema continúa el análisis hasta el final e informa todos los errores encontrados al final de la consola, indicando su tipo (Léxico, Sintáctico o Semántico), su descripción y su ubicación (línea y columna).

`[CAPTURA DE PANTALLA: consola de salida mostrando la sección "--- ERRORES ---" con varios errores de distinto tipo listados]`

## 7. Reportes generados

Al presionar el botón **Reportes** (después de haber ejecutado al menos un análisis), el sistema genera los siguientes archivos dentro de la carpeta `reportes/` del proyecto y los abre automáticamente en el navegador:

### 7.1 Reporte de tokens (`tokens.html`)

Presenta, en una tabla, cada token reconocido por el analizador léxico durante la última ejecución: número correlativo, lexema (el texto exacto reconocido), tipo de token, línea y columna donde aparece.

`[CAPTURA DE PANTALLA: reporte de tokens abierto en el navegador, mostrando varias filas de la tabla]`

### 7.2 Reporte de errores (`errores.html`)

Presenta, en una tabla, cada error detectado durante el análisis léxico, sintáctico o semántico: número correlativo, tipo de error, descripción, línea y columna. Si no se detectó ningún error, la tabla lo indica explícitamente.

`[CAPTURA DE PANTALLA: reporte de errores abierto en el navegador, mostrando errores de los tres tipos]`

### 7.3 Reporte de conjuntos y operaciones (`operaciones.html`)

Reporte adicional que lista, en una sola tabla, todos los conjuntos definidos con `CONJ` (con su definición original y sus elementos) y todas las operaciones definidas con `OPERA` (con su expresión en notación prefija y, si aplica, su forma simplificada).

`[CAPTURA DE PANTALLA: reporte de conjuntos y operaciones abierto en el navegador]`

### 7.4 Diagrama de Venn

Para cada operación definida en el análisis, el panel derecho de la ventana principal muestra un diagrama de Venn con un círculo por cada conjunto base involucrado (hasta un máximo de 3 conjuntos base; si la operación involucra más de 3, se muestra el resultado únicamente en forma de texto). La región que corresponde al conjunto resultante de la operación aparece sombreada. Debajo del diagrama se indica la expresión de la operación, el conjunto resultante y, si fue posible, su simplificación.

Es posible navegar entre los diagramas de las distintas operaciones definidas usando los botones ◀ / ▶ o el menú desplegable ubicado en la parte superior del panel.

`[CAPTURA DE PANTALLA: panel de diagramas de Venn mostrando un diagrama con 2 o 3 conjuntos y su región sombreada, junto con los controles de navegación ◀ ▶ y el selector]`

### 7.5 JSON de simplificación (`simplificacion.json`)

Archivo en formato JSON, generado en la misma carpeta `reportes/`, que documenta para cada operación definida si fue posible simplificarla aplicando las propiedades de la teoría de conjuntos (ley de doble complemento, leyes de DeMorgan, propiedades idempotentes, de absorción, distributivas, conmutativas y asociativas). Si la operación se pudo simplificar, se listan las leyes aplicadas y la expresión resultante en notación prefija; si no fue posible simplificarla, el valor asociado es el texto `"No se puede simplificar la operacion"`.

Ejemplo real generado por el sistema:

```json
{
  "demorgan": {
    "leyes": [
      "Leyes de DeMorgan",
      "Ley del doble complemento"
    ],
    "conjunto simplificado": "U & {conjA} {conjB} {conjC}"
  },
  "union1": "No se puede simplificar la operacion"
}
```

## 8. Ejemplo completo guiado

El siguiente programa ilustra el uso conjunto de todas las funcionalidades del lenguaje:

```
{
    # Definicion de conjuntos
    CONJ : conjuntoA -> 1,2,3,a,b;
    CONJ : conjuntoB -> a~z;
    CONJ : conjuntoC -> 0~9;

    # Definicion de operaciones
    OPERA : operacion1 -> & {conjuntoA} {conjuntoB};
    OPERA : operacion2 -> & U {conjuntoB} {conjuntoC} {conjuntoA};

    # Evaluamos conjuntos de datos
    EVALUAR ( {a, b, c} , operacion1 );
    EVALUAR ( {1, b} , operacion1 );
}
```

Al ejecutarlo, la consola muestra el resultado de cada evaluación (indicando "exitoso" o "fallo" según la pertenencia real de cada elemento al conjunto resultante de `operacion1`), el panel de Venn permite consultar el diagrama de `operacion1` y `operacion2`, y el botón Reportes genera los archivos HTML y el JSON descritos en la sección 7.

`[CAPTURA DE PANTALLA: ejecución completa del ejemplo anterior, mostrando editor con el código, consola con el resultado y panel de Venn con operacion1 seleccionada]`

## 9. Solución de problemas frecuentes

- **La consola no muestra ningún resultado tras presionar Ejecutar**: revisar la sección de errores al final de la consola; es posible que el programa tenga errores léxicos o sintácticos que impidieron reconocer alguna instrucción.
- **El botón Reportes no genera nada**: los reportes solo están disponibles después de haber ejecutado al menos un análisis con el botón ▶ Ejecutar.
- **Un conjunto o elemento no se reconoce como se esperaba**: verificar que el elemento sea un único carácter y que se encuentre dentro del universo permitido (símbolos ASCII entre el 33 y el 126); verificar también que el nombre del conjunto u operación no coincida con una palabra reservada (`CONJ`, `OPERA`, `EVALUAR`) ni con un operador (`U`, `&`, `^`, `-`).
- **El diagrama de Venn muestra únicamente texto en lugar de un diagrama gráfico**: esto ocurre cuando la operación involucra más de tres conjuntos base distintos, caso en el que no existe una representación geométrica clara de Venn.
