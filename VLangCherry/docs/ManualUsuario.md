# Manual de Usuario — VLangCherry

Universidad de San Carlos de Guatemala
Facultad de Ingeniería — Escuela de Ingeniería en Ciencias y Sistemas
Organización de Lenguajes y Compiladores 2
Proyecto 1 — V-Lang Cherry (VLangCherry)
Proyecto grupal (3 integrantes)

## Índice

1. Introducción
2. Requisitos para la ejecución
3. Levantar el sistema (servidor y cliente)
4. El entorno de trabajo (editor)
5. Funcionalidades del editor
6. El lenguaje VLangCherry (resumen de uso)
7. Ejecución del análisis
8. Reportes generados
   8.1 Consola
   8.2 Reporte de errores
   8.3 Reporte de tabla de símbolos
   8.4 Reporte de AST
9. Ejemplo completo guiado
10. Solución de problemas frecuentes

---

## 1. Introducción

VLangCherry es un intérprete para el lenguaje V-lang Cherry, desarrollado para el curso de Organización de Lenguajes y Compiladores 2. El lenguaje tiene una sintaxis inspirada en Go, reducida a un subconjunto de instrucciones pensado para explorar los conceptos fundamentales de un compilador (análisis léxico, sintáctico, semántico y construcción/recorrido del Árbol de Sintaxis Abstracta).

El sistema se organiza en dos componentes: un servidor que implementa el analizador léxico y sintáctico (generado con ANTLR) y el intérprete (implementado en Go), y un cliente web (React) que provee el editor de código, la ejecución y la visualización de los reportes generados. La comunicación entre ambos se realiza mediante una API REST.

## 2. Requisitos para la ejecución

- El binario del servidor, compilado para el sistema operativo de destino. La entrega del proyecto exige ejecución **nativa en Linux**; el binario se genera mediante compilación cruzada desde el entorno de desarrollo (ver sección 3.1).
- Node.js (para servir la interfaz web construida, o para ejecutar el cliente en modo desarrollo).
- Un navegador web moderno con soporte de JavaScript habilitado.
- No se requiere una base de datos ni conexión a internet: toda la información de una ejecución (errores, consola, símbolos, AST) se genera y se muestra en el momento, sin persistencia entre sesiones.

## 3. Levantar el sistema (servidor y cliente)

### 3.1 Servidor

El servidor expone la API REST en el puerto `4100` por defecto. Para levantarlo desde el código fuente (equipo de desarrollo):

```
cd server
go run ./cmd/servidor
```

Para generar el binario de entrega, apto para ejecutarse de forma nativa en Linux según lo exige el enunciado del proyecto, se realiza una compilación cruzada desde el entorno de desarrollo:

```
cd server
GOOS=linux GOARCH=amd64 go build -o vlangcherry-servidor ./cmd/servidor
```

El binario resultante (`vlangcherry-servidor`) se transfiere al equipo o servidor Linux de destino y se ejecuta directamente, sin necesidad de instalar Go en dicho equipo:

```
./vlangcherry-servidor
```

Al iniciar correctamente, el servidor indica en la salida estándar el puerto en el que quedó escuchando.

`[CAPTURA DE PANTALLA: terminal mostrando el mensaje "VLangCherry server escuchando en http://localhost:4100"]`

### 3.2 Cliente

El cliente es una aplicación web construida con React y Vite. Para levantarlo en modo desarrollo:

```
cd client
npm install
npm run dev
```

La terminal indicará la dirección local en la que quedó disponible la aplicación (por defecto un puerto asignado por Vite). Se abre dicha dirección en el navegador.

`[CAPTURA DE PANTALLA: terminal mostrando la URL local que expone "npm run dev"]`

Si el cliente y el servidor no se ejecutan en la misma máquina o en el mismo puerto por defecto, la dirección del servidor puede indicarse mediante la variable de entorno `VITE_API_URL` al construir o levantar el cliente.

## 4. El entorno de trabajo (editor)

Al abrir la aplicación en el navegador aparece la ventana principal del editor, dividida en las siguientes zonas:

- **Barra de herramientas** (parte superior): título de la aplicación y los botones Nuevo, Abrir, Guardar y ▶ Ejecutar.
- **Pestañas de archivo** (debajo de la barra de herramientas): permite tener abiertos varios archivos `.vch` de forma simultánea.
- **Panel izquierdo**: el editor de código fuente del archivo actualmente seleccionado, con una columna lateral que numera las líneas y resalta en color aquellas donde se detectó un error durante la última ejecución.
- **Panel derecho**: los reportes de la última ejecución, organizados en pestañas (Consola, Errores, Símbolos, AST).

`[CAPTURA DE PANTALLA: ventana principal de VLangCherry recién abierta, mostrando el archivo de ejemplo precargado, el panel de reportes en la pestaña Consola y las cuatro pestañas de resultado]`

## 5. Funcionalidades del editor

| Botón | Acción |
|---|---|
| **Nuevo** | Crea una nueva pestaña de edición vacía, con un nombre autogenerado y la extensión `.vch`. |
| **Abrir** | Muestra el cuadro de diálogo del sistema operativo para seleccionar un archivo `.vch` existente desde el equipo, y lo carga en una nueva pestaña. |
| **Guardar** | Descarga el contenido de la pestaña activa como un archivo `.vch` mediante el mecanismo de descarga del navegador. Al tratarse de una aplicación web, el archivo se guarda como una nueva descarga en la carpeta configurada del navegador; no se sobrescribe directamente el archivo original en el sistema de archivos. |
| **▶ Ejecutar** | Envía el contenido de la pestaña activa al servidor: se realiza el análisis léxico, sintáctico, la traducción a AST y la evaluación semántica de las instrucciones. Los reportes del panel derecho se actualizan con el resultado. |

`[CAPTURA DE PANTALLA: cuadro de diálogo de selección de archivo al presionar Abrir, mostrando un archivo .vch]`

Cada pestaña de archivo muestra un indicador (`•`) cuando su contenido tiene cambios sin guardar. Es posible cerrar una pestaña con el botón `×`, siempre que quede al menos un archivo abierto.

Cada vez que se presiona **▶ Ejecutar**, el análisis se realiza sobre un entorno completamente nuevo en el servidor: los reportes mostrados corresponden siempre a la última ejecución, sin mezclar información de ejecuciones anteriores.

## 6. El lenguaje VLangCherry (resumen de uso)

El lenguaje es sensible a mayúsculas y minúsculas. Todo programa debe declarar una función `main`, punto de entrada de la ejecución:

```
func main() {
    println("Hola, VLangCherry")
}
```

Aspectos principales (la gramática completa y las decisiones de diseño frente al enunciado se documentan en `docs/gramatica.txt`):

- **Variables**: declaración explícita (`mut edad int = 25`) o inferida (`nombre := "Ana"`); `mut` indica que la variable puede reasignarse. Reasignar una variable declarada **sin** `mut` (con `=`, `+=`, `-=`, `++` o `--`) es un error semántico; `mut` gobierna solo la reasignación de la variable completa: mutar un campo o un elemento (`persona.Edad = 30`, `numeros[0] = 9`) se permite aunque la variable no sea `mut`.
- **Tipos primitivos**: `int`, `float64`, `string`, `bool`, `rune`.
- **Slices**: `numeros := []int{1, 2, 3}`; funciones nativas `len`, `append`, `indexOf`, `join`.
- **Structs y métodos**: `struct Persona { string Nombre; int Edad }`, con métodos asociados por valor o por referencia (`func (p Persona) Saludar() string` / `func (p *Persona) Cumplir()`).
- **Control de flujo**: `if/else if/else`, `switch/case` (sin encadenamiento entre casos — cada `case` finaliza automáticamente), y `for` en sus tres formas (con condición, con inicialización/incremento, y `for indice, valor in slice`). Las sentencias `break` y `continue` solo son válidas dentro de un `for` (`break` también dentro de un `switch`); usarlas fuera de ese contexto se reporta como error semántico.
- **Funciones nativas**: `print`/`println`, `len`, `append`, `indexOf`, `join`, `Atoi`, `parseFloat`, `typeOf`, y conversión explícita de tipo (`int(x)`, `float64(x)`, `string(x)`, `bool(x)`, `rune(x)`).
  - **`print` y `println` se comportan igual: ambos emiten una línea.** El enunciado (§7.2.1) solo define `Println`; `print` se acepta como sinónimo por comodidad. No existe una variante que escriba *sin* salto de línea, porque la consola de este intérprete es una **lista de líneas** y no un flujo de caracteres (ver 8.1). Para armar una línea a partir de varios valores, pasálos como argumentos —`println("x =", x, "y =", y)`, que los separa con un espacio— o concatenalos con `+`.
- **Nombres en el ámbito global**: una función, una variable global y un `struct` no pueden compartir el mismo nombre; el sistema lo reporta como error semántico si ocurre.

Los archivos de ejemplo de la carpeta `entradas/` en la raíz del proyecto ilustran cada uno de estos aspectos por separado.

## 7. Ejecución del análisis

Para ejecutar un programa:

1. Escribir o cargar el código fuente en una pestaña del editor.
2. Presionar el botón **▶ Ejecutar**.
3. El sistema selecciona automáticamente la pestaña de resultado adecuada: si se detectaron errores, se muestra la pestaña **Errores**; en caso contrario, se muestra la pestaña **Consola**.
4. Revisar el resultado en el panel derecho, y consultar las pestañas Símbolos y AST según se necesite.

`[CAPTURA DE PANTALLA: ejecución de un archivo sin errores, mostrando la pestaña Consola activa con el resultado de varias sentencias print]`

Si el código fuente contiene errores léxicos, sintácticos o semánticos, la ejecución **no se detiene** en el primer error encontrado: el sistema recolecta todos los errores detectados en las distintas fases del análisis y los muestra juntos en la pestaña Errores, indicando su tipo, descripción, línea y columna.

`[CAPTURA DE PANTALLA: ejecución de un archivo con errores, mostrando la pestaña Errores activa con varias filas de distinto tipo — Léxico, Sintáctico, Semántico]`

## 8. Reportes generados

### 8.1 Consola

Muestra, línea por línea, la salida producida por las llamadas a `print`/`println` durante la ejecución del programa.

La consola es una **secuencia de líneas**, no un flujo de caracteres: cada llamada a `print` o `println` agrega exactamente **una** entrada, y el servidor la envía al cliente como un arreglo de cadenas. Esa es la razón por la que las dos funciones son equivalentes y por la que no hay forma de escribir "a media línea".

### 8.2 Reporte de errores

Tabla con una fila por cada error detectado, con las columnas: número correlativo, tipo (Léxico, Sintáctico o Semántico), descripción, línea y columna. Al hacer clic sobre una fila, el editor salta automáticamente a la línea correspondiente y la resalta en la columna de numeración. Si no se detectó ningún error, el panel lo indica explícitamente.

`[CAPTURA DE PANTALLA: reporte de errores con varias filas, una de ellas seleccionada y el editor mostrando la línea resaltada]`

### 8.3 Reporte de tabla de símbolos

Tabla con una fila por cada variable, función o método declarado durante la ejecución: identificador, tipo de símbolo (Variable, Función, Método), tipo de dato, entorno (ámbito) en el que fue declarado, valor, línea y columna. Al igual que en el reporte de errores, hacer clic sobre una fila lleva el editor a la línea de la declaración.

`[CAPTURA DE PANTALLA: reporte de tabla de símbolos mostrando variables globales y variables locales a una función, con su entorno correspondiente]`

### 8.4 Reporte de AST

Presenta el Árbol de Sintaxis Abstracta generado a partir del código de entrada, como un grafo interactivo: es posible desplazarse (arrastrar), acercar o alejar el zoom, y pasar el cursor sobre un nodo para identificar su contenido. La disposición es jerárquica, de arriba hacia abajo, siguiendo la estructura del programa.

`[CAPTURA DE PANTALLA: reporte de AST mostrando el árbol completo de un programa de ejemplo, con la raíz PROGRAMA en la parte superior]`

## 9. Ejemplo completo guiado

El siguiente programa (equivalente a `entradas/ejemplo2_structs.vch`) ilustra el uso conjunto de structs, métodos por valor y por referencia:

```
struct Direccion {
    string Ciudad
    int CodigoPostal
}

struct Persona {
    string Nombre
    int Edad
    Direccion Domicilio
}

func (p Persona) Saludar() string {
    return "Hola, soy " + p.Nombre
}

func (p *Persona) Cumplir() {
    p.Edad = p.Edad + 1
}

func main() {
    mut casa := Direccion{Ciudad: "Guatemala", CodigoPostal: 1010}
    mut persona := Persona{Nombre: "Ana", Edad: 25, Domicilio: casa}

    print(persona.Saludar())
    print("Edad antes:", persona.Edad)
    persona.Cumplir()
    print("Edad despues:", persona.Edad)
    print("Ciudad:", persona.Domicilio.Ciudad)
}
```

Al ejecutarlo, la consola muestra el saludo, la edad antes y después de invocar `Cumplir()` (que la incrementa por referencia) y la ciudad del domicilio anidado. El reporte de símbolos muestra las variables `casa` y `persona` en el entorno `main`, así como `Saludar` y `Cumplir` como métodos en el entorno global. El reporte de AST permite inspeccionar visualmente la estructura completa del programa, incluyendo la definición de ambos structs y de ambos métodos.

`[CAPTURA DE PANTALLA: ejecución completa del ejemplo anterior, mostrando el editor con el código, la consola con el resultado y la tabla de símbolos]`

## 10. Solución de problemas frecuentes

- **El botón ▶ Ejecutar muestra un aviso de error de red**: verificar que el servidor esté en ejecución y que la dirección configurada (`VITE_API_URL`, o `http://localhost:4100` por defecto) sea la correcta.
- **La consola no muestra ningún resultado tras ejecutar**: revisar la pestaña Errores; es posible que el programa tenga errores léxicos o sintácticos que impidieron reconocer alguna instrucción.
- **El botón Guardar no sobrescribe el archivo original**: es el comportamiento esperado de una aplicación web — el navegador únicamente permite iniciar una descarga nueva, no escribir directamente sobre un archivo ya existente en el disco.
- **El reporte de AST aparece vacío**: el AST solo se genera después de presionar ▶ Ejecutar al menos una vez sobre el archivo activo.
- **Un valor no se reconoce como se esperaba (por ejemplo, una asignación entre `int` y `float64`)**: revisar la sección de decisiones de diseño en `docs/gramatica.txt`, donde se documentan los casos en los que el lenguaje permite conversión implícita entre tipos numéricos.
