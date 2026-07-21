# Manual de Usuario — CompInterpreter

**Universidad de San Carlos de Guatemala**
**Facultad de Ingeniería — Escuela de Ciencias y Sistemas**
**Organización de Lenguajes y Compiladores 1**
**Proyecto 2 — Segundo Semestre 2024**

---

## 1. Introducción

CompInterpreter es un intérprete web, de arquitectura cliente-servidor, para el lenguaje de programación `.ci`. El presente manual describe, desde la perspectiva del usuario final, cómo instalar, ejecutar y utilizar la aplicación: la puesta en marcha del servidor y del cliente, el uso del editor de código, la ejecución de un programa y la interpretación de cada uno de los reportes generados (consola, tabla de errores, tabla de símbolos y árbol de sintaxis abstracta).

## 2. Requisitos previos

Para ejecutar la aplicación es necesario contar con:

- Node.js (versión 18 o superior) y npm.
- Un navegador web moderno (Google Chrome, Microsoft Edge o Firefox), dado que toda la interacción con el intérprete se realiza desde el navegador.
- Conexión de red local entre el cliente y el servidor (ambos pueden ejecutarse en la misma máquina durante el desarrollo).

## 3. Instalación

El proyecto está dividido en dos componentes independientes: el servidor (`server/`) y el cliente (`client/`). Cada uno cuenta con su propio archivo `package.json` y debe instalarse por separado.

### 3.1 Instalación del servidor

1. Abrir una terminal en la carpeta `server/`.
2. Ejecutar:

   ```
   npm install
   ```

   Esto instala las dependencias del servidor: `express` (framework web) y `cors` (habilitación de peticiones entre orígenes distintos), además de `jison` como herramienta de desarrollo para regenerar el analizador léxico-sintáctico.

### 3.2 Instalación del cliente

1. Abrir una terminal en la carpeta `client/`.
2. Ejecutar:

   ```
   npm install
   ```

   Esto instala React, las herramientas de construcción (Vite) y la librería de graficación de grafos (`vis-network` / `vis-data`), utilizada para el reporte del árbol de sintaxis abstracta.

## 4. Puesta en marcha de la aplicación

La aplicación requiere que **ambos componentes estén en ejecución simultáneamente**: el servidor atiende las peticiones de análisis e interpretación, y el cliente provee la interfaz gráfica en el navegador.

### 4.1 Iniciar el servidor

Desde la carpeta `server/`:

```
npm start
```

Este comando ejecuta `node server.js`. Si el inicio es exitoso, la terminal muestra el mensaje:

```
CompInterpreter server escuchando en http://localhost:4000
```

El servidor queda escuchando en el puerto 4000 por defecto (configurable mediante la variable de entorno `PORT`).

### 4.2 Iniciar el cliente

Desde la carpeta `client/`, con el servidor ya en ejecución, ejecutar:

```
npm run dev
```

(El script `npm start` es equivalente, ya que ambos invocan la misma herramienta de desarrollo, Vite). La terminal mostrará una dirección local, típicamente:

```
http://localhost:5173
```

Se debe abrir dicha dirección en el navegador para acceder a la aplicación.

> Nota: si el servidor se ejecuta en una dirección distinta a `http://localhost:4000`, debe definirse la variable de entorno `VITE_API_URL` antes de iniciar el cliente, apuntando a la dirección correcta.

`[CAPTURA DE PANTALLA: terminal mostrando el servidor iniciado ("CompInterpreter server escuchando en http://localhost:4000") y, en otra terminal, el cliente Vite iniciado con la URL local]`

## 5. Descripción de la interfaz

Al abrir la aplicación en el navegador se presentan cuatro áreas principales:

1. **Barra de herramientas** (parte superior): título de la aplicación y los botones de acción.
2. **Pestañas de archivos**: permite alternar entre los distintos archivos abiertos.
3. **Editor de código** (panel izquierdo): área de edición del código fuente `.ci`.
4. **Panel de reportes** (panel derecho): pestañas de Consola, Errores, Símbolos y AST.

`[CAPTURA DE PANTALLA: vista general de la aplicación recién cargada, señalando las cuatro áreas descritas]`

## 6. Uso del editor

### 6.1 Crear un archivo nuevo

El botón **Nuevo** de la barra de herramientas crea una pestaña con un archivo en blanco (nombrado automáticamente `sin-titulo-N.ci`), que queda disponible para editar de inmediato.

### 6.2 Abrir un archivo existente

El botón **Abrir** despliega el selector de archivos del sistema operativo, filtrado por la extensión `.ci`. Al seleccionar un archivo, su contenido se carga en una nueva pestaña del editor.

### 6.3 Guardar un archivo

El botón **Guardar** descarga el contenido del archivo actualmente activo como un archivo de texto con extensión `.ci`, utilizando el mecanismo de descarga del propio navegador. El nombre de la pestaña muestra un punto (•) cuando existen cambios sin guardar.

`[CAPTURA DE PANTALLA: barra de herramientas con los botones Nuevo, Abrir, Guardar y ▶ Ejecutar visibles]`

### 6.4 Trabajar con múltiples archivos (pestañas)

El editor permite tener abiertos varios archivos `.ci` de forma simultánea. Cada archivo se representa como una pestaña independiente; para cambiar de archivo basta con hacer clic sobre la pestaña deseada. Cada pestaña, salvo que sea la única abierta, incluye un botón de cierre (×).

`[CAPTURA DE PANTALLA: barra de pestañas con al menos dos o tres archivos abiertos, uno de ellos marcado como "sin guardar"]`

### 6.5 Numeración de líneas y línea actual

El editor muestra, en su margen izquierdo, la numeración de todas las líneas del archivo (columna de gutter). La línea donde se encuentra el cursor se resalta visualmente, tanto en el propio texto como en su número correspondiente en el margen. Esto permite ubicar con precisión la posición actual dentro del código, especialmente útil al revisar errores reportados con número de línea.

`[CAPTURA DE PANTALLA: editor con el cursor en una línea intermedia del archivo, mostrando el resaltado de esa línea y de su número en el margen]`

## 7. Ejecutar un programa

El botón **▶ Ejecutar**, ubicado en la barra de herramientas, envía el contenido del archivo actualmente activo al servidor para su análisis léxico, sintáctico y semántico, y para su ejecución. Mientras la operación está en curso, el botón muestra el texto "Ejecutando…" y queda deshabilitado.

Al finalizar la ejecución:

- Si el análisis **no encontró errores**, el panel de reportes se ubica automáticamente en la pestaña **Consola**, mostrando la salida generada por las instrucciones `echo` del programa.
- Si el análisis **encontró uno o más errores** (léxicos, sintácticos o semánticos), el panel de reportes se ubica automáticamente en la pestaña **Errores**, y dicha pestaña muestra entre paréntesis la cantidad total de errores encontrados.

`[CAPTURA DE PANTALLA: botón ▶ Ejecutar antes y durante la ejecución ("Ejecutando…")]`

## 8. Reportes generados

El panel derecho de la aplicación contiene cuatro pestañas de reportes. A continuación se describe la información que presenta cada una.

### 8.1 Consola

Muestra, en orden, cada una de las salidas producidas por las instrucciones `echo` del programa ejecutado. Si el programa no produjo ninguna salida, se indica que no hay contenido para mostrar.

`[CAPTURA DE PANTALLA: pestaña Consola con la salida de un programa de ejemplo ejecutado exitosamente]`

### 8.2 Errores

Presenta una tabla con las siguientes columnas: número de fila, tipo de error (Léxico, Sintáctico o Semántico), descripción del error, línea y columna donde fue detectado. Al hacer clic sobre cualquier fila de la tabla, el editor salta automáticamente a la línea correspondiente y la resalta, facilitando la localización y corrección del error en el código fuente.

`[CAPTURA DE PANTALLA: pestaña Errores mostrando varios errores de distinto tipo, con un clic en una fila y el editor saltando a esa línea]`

### 8.3 Símbolos

Presenta una tabla con las variables, vectores, funciones, métodos y parámetros declarados durante la ejecución del programa, con las columnas: número de fila, identificador, tipo (mostrado tal como lo etiqueta el servidor: Variable, Vector, Funcion, Metodo o Parametro), tipo de dato, entorno donde fue declarado, valor actual, línea y columna de declaración. Al igual que en la tabla de errores, al hacer clic sobre una fila el editor salta a la línea correspondiente.

`[CAPTURA DE PANTALLA: pestaña Símbolos con varias filas de variables y al menos una función/método declarados]`

### 8.4 AST (árbol de sintaxis abstracta)

Muestra, en forma de grafo interactivo, la estructura sintáctica del programa analizado. Cada nodo del grafo representa una construcción del lenguaje (una declaración, una expresión, una instrucción de control, etc.) y las flechas indican la relación de composición entre nodos padre e hijo. El grafo puede desplazarse (arrastrando) y aumentarse o reducirse su tamaño (zoom) mediante la rueda del mouse, para facilitar la inspección de árboles grandes.

`[CAPTURA DE PANTALLA: pestaña AST mostrando el grafo generado a partir de un programa de ejemplo, con varios niveles de profundidad visibles]`

## 9. Flujo de trabajo recomendado

1. Iniciar el servidor (`npm start` en `server/`).
2. Iniciar el cliente (`npm run dev` en `client/`) y abrir la dirección indicada en el navegador.
3. Escribir o abrir el archivo `.ci` a analizar.
4. Presionar **▶ Ejecutar**.
5. Revisar la pestaña **Consola** (si no hubo errores) o la pestaña **Errores** (si los hubo, corrigiendo cada uno en el editor mediante el salto automático de línea).
6. Consultar, según sea necesario, la tabla de **Símbolos** y el reporte de **AST** para verificar el comportamiento del programa.

## 10. Solución de problemas comunes

- **El cliente muestra un mensaje indicando que el servidor respondió con un error o no está disponible**: verificar que el servidor Express se encuentre en ejecución y que la dirección configurada en el cliente (`VITE_API_URL`, si aplica) corresponda con el puerto real en el que escucha el servidor.
- **El botón Guardar no descarga el archivo**: verificar la configuración de descargas del navegador; el mecanismo utilizado es el estándar del navegador para descarga de archivos generados dinámicamente.
- **El árbol de sintaxis (AST) aparece vacío**: esto ocurre cuando el programa no pudo ejecutarse en absoluto (por ejemplo, un error sintáctico que impide construir el árbol); revisar primero la pestaña de Errores.
