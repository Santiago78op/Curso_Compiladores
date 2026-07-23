# Manual Técnico — CompInterpreter

**Universidad de San Carlos de Guatemala**
**Facultad de Ingeniería — Escuela de Ciencias y Sistemas**
**Organización de Lenguajes y Compiladores 1**
**Proyecto 2 — Segundo Semestre 2024**

---

## 1. Introducción

El presente manual describe la arquitectura, las decisiones de diseño y los componentes técnicos más relevantes del proyecto CompInterpreter, con el propósito de servir de referencia para su mantenimiento futuro. CompInterpreter es un intérprete del lenguaje `.ci`, implementado con una arquitectura cliente-servidor: un servidor realiza el análisis léxico, sintáctico y semántico, y ejecuta el programa; un cliente web provee el entorno de edición y la visualización de los reportes generados.

## 2. Lenguaje de programación y herramientas utilizadas

| Componente | Tecnología |
|---|---|
| Lenguaje de implementación | JavaScript (Node.js en el servidor, JavaScript/JSX en el cliente) |
| Análisis léxico y sintáctico | Jison (generador de analizadores LALR(1), equivalente en el ecosistema JavaScript a herramientas como Yacc/Bison) |
| Servidor | Express (framework web para Node.js) |
| Comunicación cliente-servidor | API REST sobre HTTP, formato JSON |
| Cliente | React 19, construido y servido con Vite |
| Visualización del árbol de sintaxis abstracta | vis-network y vis-data (grafo interactivo) |
| Pruebas automatizadas | Playwright (pruebas de extremo a extremo sobre el cliente y el servidor en conjunto) |
| Control de concurrencia entre orígenes | CORS (paquete `cors`, habilitado en el servidor) |

## 3. Arquitectura general

La aplicación está organizada en dos módulos independientes, cada uno con su propio ciclo de vida y su propio archivo de dependencias (`package.json`):

```
CompInterpreter/
├── server/     Servidor Express + analizador Jison + intérprete
│   └── src/
│       ├── grammar.jison        Gramática léxica y sintáctica; construye el AST
│       ├── parser.js            Analizador generado por Jison a partir de grammar.jison
│       ├── analizar.js          Orquestador: código fuente -> resultado del análisis
│       ├── interprete/          Módulos del análisis semántico y la ejecución
│       └── reportes/            Generación del reporte de AST como grafo
└── client/     Cliente React (editor + reportes)
    └── src/
        ├── App.jsx              Componente raíz: estado de la aplicación
        ├── api.js                Comunicación con el servidor (fetch)
        └── components/           Editor, pestañas, tablas de reportes, grafo del AST
```

El servidor no mantiene estado entre peticiones: cada solicitud de análisis crea sus propias estructuras (lista de errores, entorno de ejecución), de modo que ejecuciones sucesivas no se contaminan entre sí.

### 3.1 Contrato de la API REST

El servidor expone un único endpoint relevante para la interpretación de programas:

```
POST /interpretar
Cuerpo de la petición (JSON):  { "codigo": "<código fuente .ci>" }

Respuesta (JSON):
{
  "errores":        Array<{ tipo, descripcion, linea, columna }>,
  "consola":        String (líneas de salida unidas con salto de línea),
  "consolaLineas":  Array<String> (cada línea de salida como elemento independiente),
  "simbolos":       Array<{ id, categoria, tipoDato, entorno, valor, linea, columna }>,
  "ast":            { nodes: Array<{ id, label }>, edges: Array<{ from, to, label }> },
  "dot":            String (representación equivalente en formato DOT de Graphviz)
}
```

Adicionalmente, el servidor expone `GET /` (información general) y `GET /salud` (verificación de disponibilidad, utilizada por la suite de pruebas automatizadas para esperar a que el servidor esté listo antes de ejecutar las pruebas).

El servidor habilita CORS de forma global (`app.use(cors())`) para permitir que el cliente, servido desde un origen distinto durante el desarrollo (por ejemplo, `http://localhost:5173`), pueda invocar el endpoint sin restricciones del navegador.

### 3.2 Decisión de diseño: entorno fresco por cada llamada

Cada invocación a `POST /interpretar` construye una nueva instancia de la lista de errores y del intérprete (y, con él, un nuevo entorno global de ejecución). Esta decisión evita que variables, funciones o errores de una ejecución anterior persistan o interfieran en la siguiente, garantizando que los reportes correspondan exclusivamente al último análisis realizado.

## 4. Análisis léxico y sintáctico

### 4.1 Gramática

El archivo `server/src/grammar.jison` integra, en un único documento, la especificación léxica (sección `%lex`) y la especificación sintáctica (gramática de producciones). El lenguaje es insensible a mayúsculas y minúsculas (`%options case-insensitive`), característica aplicada uniformemente a palabras reservadas e identificadores.

El archivo de gramática independiente del contexto, en notación BNF, requerido como entregable, se encuentra en `docs/gramatica.txt`; dicho archivo constituye una redacción propia, distinta de la especificación de Jison, según lo exigido por el enunciado del proyecto.

### 4.2 Generación del árbol de sintaxis abstracta

A diferencia de un esquema de traducción totalmente dirigido por la sintaxis sin árbol intermedio, CompInterpreter construye explícitamente un árbol de sintaxis abstracta (AST), dado que el lenguaje incluye estructuras de control de flujo (condicionales, ciclos, funciones) que requieren ser recorridas más de una vez o en un orden distinto al de su aparición textual. Cada producción de la gramática construye su nodo mediante una función auxiliar:

```javascript
function nodo(tipo, props, loc) {
  var n = { tipo: tipo };
  if (props) { for (var k in props) { if (props.hasOwnProperty(k)) n[k] = props[k]; } }
  if (loc) { n.linea = loc.first_line; n.columna = loc.first_column + 1; }
  return n;
}
```

Todo nodo del árbol posee, como mínimo, la propiedad `tipo` (identificador textual de la construcción sintáctica) y las propiedades `linea`/`columna` (empleadas por el reporte de errores), además de las propiedades específicas de cada construcción.

### 4.3 Manejo de errores léxicos y sintácticos

El módulo `server/src/analizar.js` actúa como punto de integración entre el analizador generado por Jison y el resto del sistema. Antes de invocar `parser.parse(codigo)`, se asigna al analizador un objeto compartido `yy`, cuya propiedad `errores` referencia el mismo arreglo utilizado por la lista de errores del intérprete:

```javascript
parser.yy = { errores: errores.errores };
```

Dado que en Jison el analizador léxico y el analizador sintáctico comparten el objeto `yy`, el analizador léxico puede registrar directamente sus errores (por ejemplo, un carácter no reconocido) en dicho arreglo, sin necesidad de un canal de comunicación adicional. Los errores sintácticos se capturan sobrescribiendo la función `parser.yy.parseError`, la cual construye una descripción a partir de la información que provee Jison sobre el token encontrado y los tokens esperados, limitando esta última lista a un máximo de seis elementos para mantener la legibilidad del mensaje.

## 5. Análisis semántico y ejecución

El módulo `server/src/interprete/interprete.js` implementa la clase `Interprete`, responsable de recorrer el AST realizando simultáneamente la comprobación de tipos (análisis semántico) y la ejecución de las instrucciones.

### 5.1 Recorrido en dos fases

Conforme a la recomendación del enunciado del proyecto, el intérprete recorre el árbol en dos fases sucesivas: en la primera se registran todas las funciones y métodos declarados en el ámbito global; en la segunda se ejecutan las declaraciones de variables globales y la sentencia `ejecutar`. Este esquema permite que una función o método sea invocado antes de su declaración textual dentro del archivo fuente (referencia adelantada).

### 5.2 Entorno de ejecución y tabla de símbolos

El módulo `server/src/interprete/entorno.js` define la clase `Entorno`, que representa un ámbito de ejecución mediante una tabla (`Map`) y una referencia a su ámbito padre, formando una cadena de ámbitos: global, luego función o método, luego bloque (dentro de estructuras de control). La búsqueda de un identificador (método `obtener`) recorre la cadena desde el ámbito actual hacia el ámbito global. Las claves de la tabla se normalizan a minúsculas, dado que el lenguaje es insensible a mayúsculas y minúsculas también en los identificadores.

La **tabla de símbolos del reporte** (`registrarSimbolo`) usa como clave `ámbito::id` y se **actualiza en cada asignación**, de modo que la columna de valor refleja el **valor final** de cada símbolo, no un registro histórico de sus sucesivos valores. Es una decisión de diseño deliberada y distinta de la del proyecto hermano CompScript, cuya tabla es un log cronológico; aquí una variable global reasignada dentro de una función actualiza su única fila en lugar de duplicarla.

### 5.3 Representación de valores

El módulo `server/src/interprete/valor.js` define la clase `Valor`, que encapsula todo dato manipulado en tiempo de ejecución junto con su tipo (`int`, `double`, `bool`, `char`, `string`, `null` o `vector`). Los valores de tipo `vector` incluyen además el tipo base de sus elementos y su dimensión (una o dos). Esta representación uniforme permite que las operaciones y funciones nativas verifiquen la compatibilidad de tipos antes de operar.

### 5.4 Operaciones y tablas de compatibilidad de tipos

El módulo `server/src/interprete/operaciones.js` implementa las operaciones aritméticas (suma, resta, multiplicación, división, potencia, raíz, módulo), relacionales y lógicas, mediante tablas de compatibilidad de tipos que reproducen las especificaciones del enunciado del proyecto (sección 5.5). Toda combinación de tipos no contemplada en dichas tablas produce un error semántico.

### 5.5 Conversión de tipos: coerción implícita y casteo explícito

El lenguaje distingue dos mecanismos de conversión de tipos:

- **Coerción implícita** (método `coercionar` de la clase `Interprete`): se aplica automáticamente al asignar un valor a una variable, al pasar un argumento a un parámetro o al retornar un valor desde una función, cuando el tipo de origen y el tipo de destino son numéricamente compatibles (por ejemplo, de `double` a `int`, truncando la parte decimal; o de `bool` a `int`/`double`).
- **Casteo explícito** (palabra reservada `cast`, sección 5.13 del enunciado): permite conversiones adicionales solicitadas explícitamente por el programador (por ejemplo, `int` a `char` mediante el valor ASCII correspondiente, o cualquier tipo numérico a `string`).

### 5.6 Propagación de errores mediante valores nulos

Cuando la evaluación de una expresión produce un error semántico, la función responsable retorna el valor `null` de JavaScript (a no confundir con el valor de tipo `null` propio del lenguaje `.ci`, que es una instancia válida de la clase `Valor`). Toda operación que recibe un operando `null` de JavaScript propaga dicho valor sin generar un nuevo mensaje de error, evitando que un único error de causa raíz produzca una cascada de errores derivados sobre la misma expresión.

### 5.7 Control de flujo mediante señales

Las instrucciones `break`, `continue` y `return` se implementan mediante un valor de señal (un objeto con la propiedad `tipo` igual a `'BREAK'`, `'CONTINUE'` o `'RETURN'`) que cada instrucción retorna hacia la instrucción que la contiene, en lugar de emplear excepciones del lenguaje JavaScript. Cada estructura de control (`while`, `for`, `do-until`, `loop`, `switch`) interpreta la señal recibida de acuerdo con su semántica particular. En el caso de la instrucción `switch`, este mecanismo permite implementar el comportamiento de caída consecutiva entre casos (*fall-through*, sección 5.16.2 del enunciado): una vez que un caso resulta verdadero, la ejecución continúa en los casos siguientes hasta encontrar una señal `BREAK`.

Como medida de protección ante ciclos infinitos o recursión excesiva —no exigida por el enunciado, pero necesaria para la estabilidad del servidor— se establecen límites máximos de iteración (1,000,000 por ciclo) y de profundidad de llamadas a funciones (2,000 niveles).

Una señal `BREAK` o `CONTINUE` que llega sin ser consumida hasta el cuerpo completo de la función o método (es decir, ningún ciclo que la contuviera la interceptó) se reporta como error semántico —"La sentencia 'continue'/'break' debe estar dentro de un ciclo"— en el punto donde `invocar` recibe la señal final, en vez de truncar en silencio el resto de la ejecución. Esto corresponde literalmente a la sección 5.18.2 del enunciado para `continue` ("siempre debe de estar dentro de un ciclo, de lo contrario será un error"); se aplicó el mismo criterio a `break` por consistencia, ya que antes de esta corrección una instrucción de este tipo fuera de cualquier ciclo detenía la ejecución de las instrucciones restantes sin generar ningún error en el reporte.

Adicionalmente, al construir un vector a partir de una lista literal (`vector_init1`/`vector_init2` de tipo `VECTOR_LISTA`/`VECTOR_LISTA2`), si un elemento no puede coercionarse al tipo base declarado, la posición correspondiente se rellena con el valor por defecto de ese tipo (`coercionarElementoOValorDefecto`) en vez de dejar el `null` de JavaScript propio de la propagación por error dentro del arreglo. Esto es necesario porque, a diferencia de una expresión suelta, los elementos de un vector deben ser siempre instancias de `Valor`: las funciones nativas sobre vectores (`sum`, `max`, `min`, `average`) y `aTexto()` leen `.tipo` de cada elemento sin volver a comprobar `null`, por lo que un `null` crudo en esa posición producía una excepción no controlada que abortaba toda la interpretación (reportada como "Error interno durante la ejecución") en vez de un único error semántico localizado.

### 5.8 Funciones nativas

El módulo `server/src/interprete/nativas.js` implementa las doce funciones nativas del lenguaje (`lower`, `upper`, `round`, `len`, `truncate`, `toString`, `toCharArray`, `reverse`, `max`, `min`, `sum`, `average`), cada una validando el tipo de su argumento antes de operar y retornando un error semántico descriptivo en caso de incompatibilidad.

### 5.9 Validaciones semánticas reforzadas (auditoría cruzada 2026-07-23)

Una revisión cruzada de los proyectos del curso (el mismo lote que auditó VLangCherry) detectó en CompInterpreter una tanda de huecos de validación, todos latentes: ningún archivo de prueba los ejercitaba. Se corrigieron en bloque y se agregó `entradas/ejemplo_semantica.ci`, que ejercita los caminos afectados y debe correr sin errores, más casos nuevos en `entradas/ejemplo_errores.ci`. Los dos primeros son transversales —el mismo defecto apareció en más de un proyecto del curso—.

1. **Cortocircuito de `&&` y `||` (transversal).** En el caso `LOGICA` de `evaluar`, antes se evaluaban siempre ambos operandos antes de llamar a `ops.logica`. Ahora, si el operando izquierdo es booleano y ya decide el resultado (`false` para `&&`, `true` para `||`), se retorna sin evaluar el derecho. Esto es lo que permite que una guarda como `x != 0 && 10 / x > 1` no dispare "División entre cero". El operador unario `!` no se ve afectado. (Es el mismo defecto que CompScript ya había corregido en su propia auditoría y que VLangCherry también tenía.)

2. **`return <valor>` dentro de un método (transversal).** Un método es *void*; `invocar` validaba la señal `RETURN` solo para las funciones. Ahora, si el cuerpo de un método produce una señal `RETURN` con valor, se reporta un error semántico ("el método es void; return no puede llevar una expresión"), en el punto exacto del `return` (la señal ahora transporta línea y columna). Un `return;` sin valor sigue siendo válido como salida anticipada.

3. **Método (void) usado como expresión.** Un método no produce valor; usarlo donde se espera uno (`let x: int = miMetodo();`) devolvía `null` en silencio, y `x` quedaba con un valor vacío que se propagaba como si "ya hubiera habido un error antes", cuando nunca lo hubo. Ahora `invocar` recibe una bandera `comoExpresion`: los dos puntos donde se despacha una `LLAMADA` la distinguen —como **instrucción** (`false`, un método suelto es válido) o como **expresión** (`true`, un método aquí es error semántico)—.

4. **Índices de vector no enteros.** `evalAcceso` (lectura 1D/2D) y `execAsignacionVector` (asignación 1D/2D) truncaban el índice con `Math.trunc(aNumero(idx))` sin verificar su tipo. Un índice de tipo cadena o `null` daba `NaN` —y como `NaN < 0` y `NaN >= longitud` son ambos falsos, *pasaba* los chequeos de rango—, y un `double` se truncaba en silencio. Se agregó el auxiliar `indiceEntero`, aplicado en los cuatro puntos, que exige tipo `int` —la misma validación que la creación `new vector int[n]` ya hacía sobre el tamaño—.

5. **Funciones nativas numéricas sobre elementos sin asignar.** `new vector int[n]` inicializa sus posiciones con el valor `null` del lenguaje. Las nativas `sum`, `average`, `max` y `min` hacían `aNumero(elemento)`, que para un `null` da `NaN`: `sum` sobre un vector recién creado imprimía `NaN` sin error alguno. Ahora, mediante el auxiliar `tieneNulos`, las cuatro reportan un error semántico si el vector contiene elementos nulos sin asignar.

**Política documentada (argumentos faltantes).** Cuando una llamada omite un argumento que no tiene valor por defecto, se reporta el error semántico correspondiente pero la función **igualmente se ejecuta**, tomando el valor por defecto del tipo del parámetro. Es la política de "continuar y acumular errores" propia de este proyecto —opuesta al aborto inmediato de CompScript—, coherente con la propagación por `null` de la sección 5.6.

## 6. Reporte del árbol de sintaxis abstracta como grafo genérico

El módulo `server/src/reportes/ast-grafo.js` transforma el AST en una estructura de nodos y aristas adecuada para su visualización, mediante un recorrido genérico que no requiere conocimiento previo de los distintos tipos de nodo del árbol: toda propiedad de un nodo cuyo valor sea a su vez un objeto con una propiedad `tipo` de tipo cadena se interpreta como un nodo hijo (generando una arista); toda propiedad cuyo valor sea un dato escalar (número, cadena, booleano) se añade como texto descriptivo a la etiqueta del nodo padre.

```javascript
function esNodo(x) {
  return x !== null && typeof x === 'object' && typeof x.tipo === 'string';
}
```

Esta estrategia tiene la ventaja de que la incorporación de nuevos tipos de nodo al AST no requiere modificar el módulo de generación del grafo: cualquier nodo que siga la convención `{ tipo, ... }` se integra automáticamente a la visualización. Adicionalmente, el módulo genera una representación equivalente en formato DOT (`aDot`), compatible con la herramienta Graphviz, como alternativa de graficación fuera del navegador.

Del lado del cliente, el componente `AstGrafo.jsx` recibe la estructura `{ nodes, edges }` ya construida por el servidor y la renderiza mediante la librería `vis-network`, con una disposición jerárquica (`layout.hierarchical`, dirección de arriba hacia abajo) y sin simulación de física, dado que la jerarquía del árbol ya define su disposición.

## 7. Componentes principales del cliente

| Componente | Responsabilidad |
|---|---|
| `App.jsx` | Estado global de la aplicación: lista de archivos abiertos, archivo activo, resultado de la última ejecución, pestaña de reporte activa |
| `api.js` | Función `interpretar(codigo)`: invoca `POST /interpretar` sobre el servidor y retorna la respuesta en formato JSON |
| `Toolbar.jsx` | Botones Nuevo, Abrir, Guardar y Ejecutar |
| `FileTabs.jsx` | Pestañas de archivos abiertos (edición multi-archivo) |
| `Editor.jsx` | Área de edición de código, con numeración de línea, resaltado de la línea actual y marcado de líneas con error |
| `Consola.jsx` | Reporte de salida de consola (instrucciones `echo`) |
| `TablaErrores.jsx` | Reporte tabular de errores léxicos, sintácticos y semánticos, con navegación a la línea correspondiente |
| `TablaSimbolos.jsx` | Reporte tabular de la tabla de símbolos resultante de la última ejecución |
| `AstGrafo.jsx` | Reporte gráfico del árbol de sintaxis abstracta mediante `vis-network` |

### 7.1 Sincronización entre el editor y los reportes

El componente `Editor.jsx` se implementa utilizando las funciones `forwardRef` y `useImperativeHandle` de React, exponiendo un método `irALinea(numero)` que el componente padre (`App.jsx`) invoca cuando el usuario hace clic sobre una fila de la tabla de errores o de la tabla de símbolos. Esta técnica permite que un componente padre controle un aspecto interno del estado de un componente hijo (la posición del cursor y del desplazamiento vertical del editor) sin necesidad de elevar dicho estado por completo al componente padre.

## 8. Pruebas automatizadas

El proyecto incluye una suite de pruebas de extremo a extremo (`client/e2e/compinterpreter.spec.js`), implementada con Playwright, que ejecuta simultáneamente el servidor Express y el cliente Vite (configuración en `client/playwright.config.js`) y valida, contra la aplicación en funcionamiento:

1. La ejecución correcta de un programa sin errores, con verificación del contenido de la consola, de la tabla de símbolos y de la generación del grafo del AST.
2. El reporte de errores léxicos, sintácticos y semánticos de un archivo de prueba con errores intencionales, incluyendo la verificación de la navegación automática hacia la línea del error al hacer clic sobre una fila del reporte.
3. La creación de una nueva pestaña de archivo en blanco.

## 9. Archivos de prueba de referencia

El directorio `entradas/` contiene archivos `.ci` utilizados como casos de prueba durante el desarrollo, cubriendo distintos aspectos del lenguaje: el ejemplo del anexo del enunciado, un archivo con los tres tipos de error (léxico, sintáctico y semántico), funciones con parámetros por defecto y recursión, la sentencia `switch` con caída consecutiva entre casos, y el manejo de vectores de una y dos dimensiones.

## 10. Consideraciones para mantenimiento futuro

- La gramática (`server/src/grammar.jison`) debe regenerarse mediante el comando `npm run grammar` (equivalente a `jison src/grammar.jison -o src/parser.js`) cada vez que se modifique; el archivo `parser.js` es generado y no debe editarse manualmente.
- Cualquier nuevo tipo de nodo agregado al AST se integrará automáticamente al reporte gráfico, siempre que respete la convención de tener una propiedad `tipo` de tipo cadena.
- Las tablas de compatibilidad de tipos (`server/src/interprete/operaciones.js`) están documentadas en el propio código fuente mediante comentarios que referencian la sección correspondiente del enunciado del proyecto, incluyendo los casos en los que el comportamiento especificado, aunque puede resultar poco intuitivo, se implementó de forma literal por fidelidad a dicho enunciado.
