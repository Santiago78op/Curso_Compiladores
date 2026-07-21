---
tags: [proyecto, guia, lexico, sintactico, semantico, web]
aliases: [guia compinterpreter, compinterpreter paso a paso]
fuente: "Enunciado CompInterpreter (OLC1_PT2_2S2024) + construcción real (servidor y cliente COMPLETOS: 2026-07-20)"
fecha: 2026-07-20
---

# Guía de elaboración — [[CompInterpreter]]

Intérprete **web cliente-servidor** en JS/TS para el lenguaje `.ci`, *case insensitive*. A diferencia de los otros 3 proyectos del curso (Java + JFlex + CUP + JavaFX), acá el análisis léxico/sintáctico corre con **Jison** en un servidor **Express** y la interfaz vive en un cliente **React** que le habla por REST.

## Estado de construcción

| Etapa | Contenido | Estado |
|---|---|---|
| 1 | Gramática Jison (léxico + sintáctico, AST genérico) | ✅ 2026-07-20 |
| 2 | Intérprete: entorno, tipos, operaciones, nativas, errores | ✅ 2026-07-20 |
| 3 | Servidor Express (`POST /interpretar`) | ✅ 2026-07-20 |
| 4 | Reporte AST como grafo genérico (servidor) | ✅ 2026-07-20 |
| 5 | Cliente React: editor multi-archivo con gutter y resaltado | ✅ 2026-07-20 |
| 6 | Cliente React: panel de reportes (Consola/Errores/Símbolos/AST con vis-network) | ✅ 2026-07-20 |
| 7 | Suite Playwright (3 tests end-to-end) | ✅ 2026-07-20 (3/3 pasan) |

**PROYECTO FUNCIONAL COMPLETO** en ambos lados (servidor y cliente). Falta empaquetado: `docs/gramatica.txt`, `ManualUsuario.md`, `ManualTecnico.md` (documentación de entrega, ver `CompInterpreter/docs/`) y repo de entrega `OLC1_Proyecto2_#Carnet`.

## 1. Requisitos previos (reales)

- **Servidor** (`CompInterpreter/server/`): Node.js + `express` + `cors` + `jison` (devDependency). `package.json`: script `start` → `node server.js`, script `grammar` → `jison src/grammar.jison -o src/parser.js` (regenera el parser cuando cambia la gramática).
- **Cliente** (`CompInterpreter/client/`): React 19 + Vite 8 + `vis-network`/`vis-data` (grafo del AST) + `oxlint` (lint) + Playwright (`e2e/`). Scripts: `dev` **y** `start` apuntan ambos a `vite` — cualquiera de los dos levanta el cliente.
- Variable opcional `VITE_API_URL` (default `http://localhost:4000`) si el servidor no corre en el puerto por defecto — ver `client/src/api.js`.

## 2. Léxico y sintáctico — CONSTRUIDO ✅

### 2.1 Un solo archivo `.jison`

`server/src/grammar.jison` reúne `%lex` (tokens) y la gramática (`%%` sintáctico) en un único archivo — la convención de [[Jison]]. `%options case-insensitive flex` resuelve de un tirón la sección 5.1 del enunciado (case-insensitive) sin normalizar manualmente cada palabra reservada.

Detalles reales del lexer:
- Los comentarios (`// ...` y `/* ... */`) y los espacios se descartan sin devolver token.
- Las reservadas están declaradas **antes** que `[a-zA-Z_][a-zA-Z0-9_]* return 'ID';` — el mismo principio de "reservadas antes que el identificador" que en JFlex ([[Cap 3 - Análisis léxico]]).
- Los operadores de dos caracteres (`++ -- == != <= >= && ||`) están **antes** que sus prefijos de un carácter, para que el *longest match* no los trague de a uno.
- Las secuencias de escape en cadenas y caracteres (`\n \t \" \' \\`) se resuelven en la propia acción léxica con `.replace(/\\(.)/g, ...)` antes de devolver el token, no en el intérprete.
- Recuperación léxica (5.1/6.1): la regla catch-all `.` no descarta silenciosamente — empuja el error a `yy.errores` (ver más abajo) y continúa.

### 2.2 Precedencia real (sección 5.10)

```
%right IF DOSPUNTOS          /* operador ternario: mas bajo que todo */
%left OR                     /* nivel 7 */
%left AND                    /* nivel 6 */
%right NOT                   /* nivel 5 */
%left IGUAL DIFERENTE MENOR MENORIGUAL MAYOR MAYORIGUAL IS   /* nivel 4 */
%left MAS MENOS              /* nivel 3 */
%left POR DIV MODULO         /* nivel 2 */
%nonassoc POTENCIA RAIZ      /* nivel 1 */
%right UMINUS                /* nivel 0 (mayor importancia) */
```

Coincide exactamente con la tabla 5.10 del enunciado. La [[Gramática libre de contexto (BNF)|gramática BNF entregable]] (`CompInterpreter/docs/gramatica.txt`) no usa declaraciones de precedencia — en su lugar reescribe las expresiones en capas (`expresion-or > expresion-and > ... > expresion-unaria`) para quedar **libre de ambigüedad por construcción**, como exige la sección 8 ("solamente se debe escribir la gramática independiente del contexto"). Ver [[Ambigüedad, precedencia y asociatividad]].

### 2.3 El AST se construye en las propias acciones (sin pases intermedios)

A diferencia de DataForge (S-atribuida, sin AST), CompInterpreter sí construye un [[Árbol de sintaxis abstracta (AST)|AST]] explícito — lo exige el control de flujo del lenguaje (`if/switch/while/for/do-until/loop/function`). Cada producción llama a una fábrica de nodos genérica:

```js
function nodo(tipo, props, loc) {
  var n = { tipo: tipo };
  if (props) { for (var k in props) { if (props.hasOwnProperty(k)) n[k] = props[k]; } }
  if (loc) { n.linea = loc.first_line; n.columna = loc.first_column + 1; }
  return n;
}
```

Todo nodo tiene `tipo` (string) + `linea`/`columna` (para el reporte de errores) + las propiedades propias de esa producción (`cond`, `entonces`, `sino` para un `IF`; `izq`/`op`/`der` para una `BINARIA`; etc.). Esta uniformidad es la que permite que el reporte de AST (más abajo) funcione **sin mantenimiento** cuando se agregan nodos nuevos.

Una extensión propia documentada en el `.jison` (no está en 5.15.1 literal): `vector_init1`/`vector_init2` aceptan también una `expresion` genérica como inicializador (no solo `new vector` o lista literal), para poder escribir `let v: int[] = reverse(otro);` — necesaria porque varias nativas (`reverse`, `toCharArray`) devuelven vectores.

## 3. Intérprete (semántico + ejecución) — CONSTRUIDO ✅

Paquete `server/src/interprete/`:

| Archivo | Responsabilidad |
|---|---|
| `entorno.js` | [[Tabla de símbolos]] como cadena de [[Entornos y alcance|entornos]] (`Entorno { padre, nombre, tabla: Map }`); claves en minúscula (case-insensitive también para identificadores, igual que DataForge) |
| `tipos.js` | catálogo `TIPO` (`int double bool char string null vector`) + nombres en español para reportes + `valorPorDefecto` (tabla 5.3) |
| `valor.js` | clase `Valor { tipo, valor, tipoBase, dimension }` — todo dato en tiempo de ejecución va etiquetado con su tipo; `aTexto()`/`aNumero()` son los conversores comunes |
| `operaciones.js` | tablas de compatibilidad de `+ - * / ^ $ %` **transcritas literalmente** de la sección 5.5, más relacionales y lógicos |
| `nativas.js` | las 12 funciones nativas (5.22-5.25) |
| `errores.js` | `ListaErrores` (acumulador) + `ErrorSemantico` |
| `interprete.js` | `Interprete`: recorrido del AST en dos pasadas, ejecución de instrucciones, evaluación de expresiones, coerción |

### 3.1 Dos pasadas (sección 5.21, forward-reference)

El enunciado recomienda explícitamente "realizar 2 pasadas del AST generado: la primera para almacenar todas las funciones, y la segunda para las variables globales y la función ejecutar". Así quedó implementado en `Interprete.interpretar()`:

```js
// PASADA 1: funciones y metodos
for (const g of ast.globales) {
  if (g.tipo === 'FUNCION' || g.tipo === 'METODO') {
    const clave = g.id.toLowerCase();
    if (this.funciones.has(clave)) {
      this.errores.semantico('Ya existe una función o método con el identificador "' + g.id + '"...');
      continue;
    }
    this.funciones.set(clave, g);
    ...
  }
}
// PASADA 2: variables globales y ejecutar
for (const g of ast.globales) {
  if (g.tipo === 'DECLARACION' || g.tipo === 'DECLARACION_VECTOR') this.ejecutarInstruccion(g, this.global);
  else if (g.tipo === 'EJECUTAR') this.ejecutarEjecutar(g);
}
```

Esto es lo que permite llamar una función antes de su declaración textual (`ejemplo_funciones.ci` llama `factorial` y `calculaFuerza` desde `main()`, ambas declaradas más abajo en el archivo — verificado, ver sección Casos de prueba).

### 3.2 Control de flujo por señales, no por excepciones

`break`/`continue`/`return` se implementan como un valor de señal (`{ tipo: 'BREAK' }`, `{ tipo: 'RETURN', valor, tieneValor }`) que cada `ejecutarInstruccion`/`ejecutarBloque` propaga hacia arriba, en vez de lanzar excepciones de JS. Cada ciclo (`execWhile`, `execFor`, `execDoUntil`, `execLoop`) decide qué hacer con la señal recibida: `BREAK` corta el ciclo, `CONTINUE` salta a la siguiente iteración (en `for`, cae en la actualización), y `RETURN` se re-propaga hasta la función que la originó. `execSwitch` usa el mismo mecanismo para implementar el **fall-through** (5.16.2): una vez que un `case` coincide, sigue ejecutando los siguientes sin volver a evaluar su expresión, hasta toparse con una señal `BREAK`.

Guardas anti-bucle-infinito y anti-recursión: `MAX_ITER = 1_000_000` (cada ciclo) y `MAX_DEPTH = 2000` (profundidad de `invocar`), ambas decisiones propias del proyecto (no están en el enunciado) para que un `.ci` mal escrito no cuelgue el servidor.

### 3.3 Parámetros nombrados y valores por defecto (5.19)

Las llamadas (`llamada ::= ID ( args )`, `arg ::= ID ASIGNA expresion`) son siempre por nombre. `Interprete.invocar()` arma un `Map` `nombre → expresión` a partir de los argumentos provistos y, para cada parámetro formal, resuelve en este orden: (1) el argumento nombrado si vino en la llamada, evaluado **en el entorno del llamador** (para que pueda referenciar variables locales del que llama); (2) si no vino, el valor por defecto del parámetro, evaluado en el entorno local de la función; (3) si tampoco hay valor por defecto, error semántico "Falta el parámetro...".

### 3.4 Coerción de asignación — decisión de diseño documentada en el código

```js
/* DECISION DE DISENO: los tipos numericos son mutuamente asignables con
   coercion (double->int trunca). Sin esto, el propio ejemplo oficial
   del Anexo 11.1 seria invalido: `let modulo:int = x % 3` donde `%`
   produce DECIMAL segun la tabla literal del enunciado (5.5.6). */
```

`coercionar(tipoDest, v, l, c)` en `interprete.js` permite `double→int` (trunca), `bool→int`/`bool→double` (numérico), pero **no** deja mezclar `char`/`string`/`bool` entre sí en una asignación: cada uno solo acepta su propio tipo salvo las conversiones numéricas explícitas arriba. Esto es distinto del `cast()` explícito (sección 3.6) — la coerción es implícita y solo en asignaciones/parámetros/retorno; el `cast` es explícito y cubre combinaciones adicionales (`int↔char`, `→string`). Ver [[Conversión de tipos (coerción y cast)]].

### 3.5 Tablas de compatibilidad de operadores — transcripción literal + una decisión anotada

`operaciones.js` reproduce las tablas 5.5.1-5.5.7 del enunciado como objetos JS (`T_SUMA`, `T_RESTA`, `T_MULT`, `T_DIV`, `T_POT`, `T_RAIZ`, `T_MOD`), con `'X'` marcando combinación inválida. Un caso queda anotado explícitamente en el código porque el dato del enunciado es contraintuitivo pero se respetó tal cual:

```js
// 5.5.6 Modulo  -> NOTA: el enunciado define Entero % Entero = DECIMAL
// (identico a la tabla de Raiz). Es inusual (uno esperaria Entero), pero
// se implementa TAL CUAL lo dice la fuente, sin "corregirlo".
```

Los operadores relacionales (`==`, `!=`, `<`, `<=`, `>`, `>=`) no tenían la matriz de compatibilidad recuperable del PDF original (celdas vacías tras la conversión — documentado en `doc/Projects/OLC1_PT2_2S2024.clean.md`); el criterio adoptado quedó igualmente anotado en el código como decisión propia y no como dato del enunciado: igualdad admite `null` contra cualquier tipo, numéricos (`int/double/char`) entre sí, `string` solo con `string`, `bool` solo con `bool`; orden admite numéricos entre sí y `string` con `string` (lexicográfico), cualquier otra combinación (incluidos `bool` y `null`) es error semántico. Verificado con `ejemplo_errores.ci`: `bandera > true` (dos `BOOLEAN`) produce el error "No se puede comparar (>) entre BOOLEAN y BOOLEAN".

### 3.6 Propagación por null (mismo principio que DataForge)

Igual que en [[DataForge]]: una expresión que falla devuelve JS `null` (¡distinto del `Valor` de tipo `'null'` del propio lenguaje, que es un objeto `Valor` válido!) y las operaciones que reciben `null` como operando lo propagan en silencio (`if (izq === null || der === null) return null;`), evitando cascadas de errores por una sola causa raíz. Ver [[Manejo de errores (léxicos, sintácticos, semánticos)]].

### 3.7 Errores compartidos entre lexer y parser vía `yy`

`server/src/analizar.js` es el orquestador: crea una `ListaErrores` y conecta su arreglo interno directamente como `parser.yy.errores` **antes** de llamar a `parser.parse(codigo)`. Como el lexer generado por Jison y las acciones sintácticas comparten el mismo objeto `yy`, el lexer puede empujar errores léxicos (`yy.errores.push(...)`, ver `grammar.jison` línea ~126) al mismo arreglo que después usa el resto del pipeline — sin pasar el error por ningún valor de retorno.

Los errores sintácticos se capturan sobreescribiendo `parser.yy.parseError`: arma una descripción a partir de `hash.token`/`hash.text`/`hash.expected` (los primeros 6 tokens esperados, para no generar una lista gigante) y, si el error no es recuperable, relanza para que el `try/catch` de `analizar()` corte el `parse()` — pero los errores ya quedaron registrados en la lista compartida antes de relanzar, así que no se pierden.

```js
function analizar(codigo) {
  const errores = new ListaErrores();
  const parser = obtenerParser();
  parser.yy = { errores: errores.errores };
  parser.yy.parseError = function (msg, hash) { /* arma descripcion + errores.sintactico(...) */ };
  parser.parseError = parser.yy.parseError;
  let ast = null;
  try { ast = parser.parse(codigo); }
  catch (e) { if (!errores.hay()) errores.sintactico(e.message || 'Error de análisis', 0, 0); }
  const interprete = new Interprete(errores);
  if (ast) { try { interprete.interpretar(ast); } catch (e) { errores.semantico(...); } }
  ...
}
```

**Entorno fresco por llamada**: cada `POST /interpretar` crea su propia `ListaErrores`, su propio `Interprete` (con su propio `Entorno` global) — no hay estado compartido entre ejecuciones, igual que la regla de DataForge.

## 4. Reporte de AST como grafo genérico — CONSTRUIDO ✅

`server/src/reportes/ast-grafo.js` recorre el AST **sin conocer los tipos de nodo de antemano**: para cada objeto con propiedad `tipo` (string), sus propiedades escalares (`op`, `nombre`, `valor`...) se concatenan a la etiqueta del nodo, y sus propiedades que son otro nodo (o arreglo de nodos) se convierten en aristas hijas, recursivamente:

```js
function esNodo(x) { return x !== null && typeof x === 'object' && typeof x.tipo === 'string'; }

function construirGrafo(raiz) {
  const nodes = []; const edges = []; let contador = 0;
  function visitar(n) {
    const miId = 'n' + (contador++);
    let etiqueta = n.tipo;
    const hijos = [];
    for (const clave of Object.keys(n)) {
      if (clave === 'tipo' || clave === 'linea' || clave === 'columna') continue;
      const v = n[clave];
      if (esNodo(v)) hijos.push({ etiqueta: clave, nodo: v });
      else if (Array.isArray(v)) { /* nodos, o escalares a la etiqueta */ }
      else if (v !== null && v !== undefined) etiqueta += '\n' + clave + ': ' + String(v);
    }
    nodes.push({ id: miId, label: etiqueta });
    for (const h of hijos) { const hijoId = visitar(h.nodo); edges.push({ from: miId, to: hijoId, label: h.etiqueta }); }
    return miId;
  }
  if (raiz) visitar(raiz);
  return { nodes, edges };
}
```

Consecuencia práctica: si mañana se agrega un nodo de AST nuevo (por ejemplo, para una futura sentencia), el grafo lo dibuja automáticamente sin tocar `ast-grafo.js`. También genera el `.dot` equivalente (`aDot()`) para [[Graphviz]] por si se quisiera graficar fuera del navegador, aunque el reporte real usa [[vis-network]] del lado del cliente (sección 5).

Verificado con `entradas/ejemplo_vectores2d.ci`: el AST resultante tiene **106 nodos y 105 aristas** (un árbol: aristas = nodos - 1).

## 5. Servidor Express — CONSTRUIDO ✅

`server/server.js`: un único endpoint relevante.

```js
app.post('/interpretar', (req, res) => {
  const codigo = (req.body && typeof req.body.codigo === 'string') ? req.body.codigo : '';
  try {
    const resultado = analizar(codigo);
    res.json(resultado);
  } catch (e) {
    res.status(500).json({ errores: [{ tipo: 'Interno', descripcion: e.message || String(e), linea: 0, columna: 0 }], ... });
  }
});
```

Contrato REST real: `POST /interpretar { codigo: string }` → `{ errores, consola, consolaLineas, simbolos, ast, dot }`. `cors()` está habilitado globalmente porque el cliente Vite corre en otro puerto/origen en desarrollo.

## 6. Cliente React — CONSTRUIDO ✅

Estructura real en `client/src/`: `App.jsx` (estado y orquestación) + `api.js` (fetch a `/interpretar`) + `components/{Toolbar,FileTabs,Editor,Consola,TablaErrores,TablaSimbolos,AstGrafo}.jsx`.

### 6.1 Multi-archivo con pestañas (4.1/4.2)

`App.jsx` mantiene `archivos` como un arreglo de `{ id, nombre, contenido, sinGuardar }`; `FileTabs.jsx` las renderiza y permite cerrar (deshabilitado si solo queda una). Nuevo/Abrir/Guardar (`Toolbar.jsx`) son las tres funcionalidades de la sección 4.2:
- **Nuevo**: agrega un archivo vacío `sin-titulo-N.ci` y lo activa.
- **Abrir**: `<input type="file" accept=".ci,text/plain">` oculto + `FileReader.readAsText`.
- **Guardar**: arma un `Blob` y dispara una descarga (`URL.createObjectURL` + `<a download>`), marca `sinGuardar: false`.

### 6.2 `Editor.jsx` — gutter sincronizado + línea actual (4.1: "deberá mostrar la línea actual")

Es un `forwardRef` con `useImperativeHandle` que expone un método `irALinea(num)` — así `App.jsx` puede pedirle al editor que salte a una línea desde afuera (por ejemplo, al hacer clic en una fila de la tabla de errores) sin acoplar el estado del cursor al componente padre:

```jsx
const Editor = forwardRef(function Editor({ value, onChange, erroresPorLinea }, ref) {
  ...
  useImperativeHandle(ref, () => ({
    irALinea(num) {
      const ta = textareaRef.current;
      const offsets = value.split('\n').reduce((acc, l) => { acc.push(acc[acc.length-1] + l.length + 1); return acc; }, [0]);
      const pos = offsets[Math.max(0, num - 1)] || 0;
      ta.focus(); ta.setSelectionRange(pos, pos);
      setLinea(num);
      ta.scrollTop = Math.max(0, (num - 1) * LINE_HEIGHT - ta.clientHeight / 2);
      ...
    },
  }));
```

El gutter es un `<div>` separado con un número por línea, cuyo `scrollTop` se sincroniza a mano con el del `<textarea>` en el handler `onScroll` (no hay `<textarea>` con números de línea nativos en HTML). La línea actual se recalcula contando saltos de línea hasta `selectionStart` (`onKeyUp`/`onClick`/`onFocus`) y se resalta con una franja absoluta (`editor-linea-resaltada`) posicionada por `(linea-1) * LINE_HEIGHT - scrollTop`. Las líneas con error reciben la clase `tiene-error` a partir de un `Set<number>` (`erroresPorLinea`, calculado con `useMemo` en `App.jsx` a partir de `resultado.errores`).

### 6.3 Panel de reportes con pestañas (4.5, 6.1-6.4)

`App.jsx` mantiene un estado `panel` (`'consola' | 'errores' | 'simbolos' | 'ast'`) y cambia automáticamente de pestaña según el resultado de ejecutar: si hay errores, salta a "Errores"; si no, se queda en "Consola" (`setPanel(res.errores && res.errores.length ? 'errores' : 'consola')`) — verificado por la suite Playwright (sección 8).

- `Consola.jsx` (6.4): `<pre>` de solo lectura sobre `consolaLineas` (el arreglo, no el string ya unido — permite `.join('\n')` limpio).
- `TablaErrores.jsx` (6.1): columnas `# / Tipo / Descripción / Línea / Columna`; cada fila tiene `onClick` que llama `onIrALinea(e.linea)` → `editorRef.current.irALinea(num)`.
- `TablaSimbolos.jsx` (6.2): columnas `# / ID / Tipo / Tipo de Dato / Entorno / Valor / Línea / Columna`, mismo patrón de clic-para-ir-a-línea.
- `AstGrafo.jsx` (6.3): usa `vis-network` + `vis-data` directamente sobre `{ nodes, edges }` que ya llegan armados del servidor (sección 4) — el cliente no reconstruye el árbol, solo lo dibuja. Layout `hierarchical` (`direction: 'UD'`, `sortMethod: 'directed'`), `physics: false` (árbol ya jerárquico, no necesita simulación de fuerzas), colores tomados de las custom properties CSS del tema (`--text-h`, `--code-bg`, `--accent`, `--border`) para que el grafo respete claro/oscuro.

## 7. Casos de prueba (reales, verificados con `server/src/analizar.js` vía consola)

Comando de verificación usado:
```bash
cd CompInterpreter/server
node -e "const {analizar}=require('./src/analizar'); const fs=require('fs'); \
  console.log(JSON.stringify(analizar(fs.readFileSync('../entradas/ARCHIVO.ci','utf-8')).errores))"
```

- **`entradas/ejemplo_funciones.ci`** (parámetros por defecto, forward-reference, recursión, cast, nativas, `is`): 0 errores, 10 símbolos, consola real:
  ```
  Fuerza (masa=10, acel default 9.8): 98.0
  Fuerza (masa=10, acel=2.5): 25.0
  5! = 120
  cast(18.6 as int) = 18
  cast(70 as char) = F
  cast(16 as double) = 16.0
  upper: HOLA
  lower: mundo
  round(15.51) = 16
  truncate(9.99) = 9
  10 is int -> true
  PI is double -> true
  ```
- **`entradas/ejemplo_switch.ci`** (fall-through, do-until, loop): 0 errores; confirma fall-through real (`case 2` sin `break` cae en `case 3` y `case 4`, "cinco" nunca se alcanza porque `case 4` sí rompe).
- **`entradas/ejemplo_vectores2d.ci`** (vectores 1D/2D, `reverse/max/min/sum/average`): 0 errores; AST con 106 nodos / 105 aristas.
- **`entradas/ejemplo_errores.ci`** (los 3 tipos de error a propósito): **8 errores** — 1 léxico (`#` no pertenece al lenguaje, línea 5), 2 sintácticos (expresión mal formada / `;` inesperado), 5 semánticos (resta entre CADENA y CADENA, variable no declarada, reasignar constante, condición de `if` no booleana, comparación `>` entre dos BOOLEAN). Todos con línea/columna. Esta es la misma entrada que usa el test Playwright de errores (sección 8).
- **`entradas/ejemplo_anexo.ci`**: el ejemplo literal del Anexo 11.1 del enunciado; es el contenido inicial (`ejemplo-anexo.js`) que carga el editor al abrir el cliente por defecto.

## 8. Suite Playwright — CONSTRUIDA ✅ (3/3 tests pasan)

`client/e2e/compinterpreter.spec.js`, corrida real contra el cliente:

1. **"ejecuta el ejemplo del anexo y muestra consola/símbolos/AST"**: clic en ▶ Ejecutar, sin errores el panel queda en "Consola" solo, verifica 3 líneas de salida esperadas, cambia a "Errores" y ve "Sin errores", a "Símbolos" y encuentra la fila de `x`, a "AST" y confirma que el `<canvas>` de vis-network se renderiza.
2. **"reporta errores léxicos/sintácticos/semánticos y salta a la línea"**: carga `ejemplo_errores.ci` en el editor, ejecuta, confirma que el panel salta solo a "Errores", cuenta **8 filas** con el desglose exacto por tipo (1 léxico / 2 sintácticos / 5 semánticos vía selectores `.fila-error-léxico` etc.), hace clic en la primera fila y confirma que el gutter marca la línea 5 como actual.
3. **"nuevo archivo agrega una pestaña editable"**: clic en "Nuevo", cuenta que las pestañas suben en 1 y que la nueva pestaña activa contiene `sin-titulo-`.

## 9. Errores comunes (los que ya nos pasaron o son fáciles de repetir)

- Confundir la gramática BNF entregable con el `.jison`: el enunciado (sección 8) exige explícitamente que NO sea copia — por eso `docs/gramatica.txt` reescribe las expresiones en capas de precedencia en vez de usar `%left`/`%right`.
- Olvidar que `yy` es el canal compartido entre lexer y parser: si no se asigna `parser.yy = { errores: ... }` **antes** de `parser.parse()`, el lexer no tiene dónde reportar sus errores léxicos.
- Perder los errores ya acumulados al relanzar la excepción de `parseError`: por eso se registran en `errores` **antes** de hacer `throw`.
- No diferenciar el `Valor` de tipo `'null'` (dato del lenguaje) del `null` de JavaScript (señal de "hubo un error antes") — son cosas distintas y mezclarlas rompe la propagación por null.
- Vectores declarados con más de un identificador (`let a, b: int[] = ...`): el intérprete lo detecta y marca error semántico ("un vector solo puede declararse con un identificador"), pero la gramática lo permite sintácticamente — hay que probarlo para no llevarse una sorpresa en el reporte de errores del proyecto real.
- `switch` sin ningún `break`: cae en fall-through hasta el final (o hasta `default`), fácil de confundir con el comportamiento de otros lenguajes si no se lee 5.16.2 con cuidado.
- CORS entre cliente (Vite, puerto propio) y servidor Express (puerto 4000): resuelto con `app.use(cors())` global en `server.js`; si se cambia el puerto del servidor hay que actualizar `VITE_API_URL` en el cliente.

## Relacionadas
- [[CompInterpreter]]
- [[Jison]]
- [[Árbol de sintaxis abstracta (AST)]] · [[Traducción dirigida por la sintaxis]]
- [[Tabla de símbolos]] · [[Entornos y alcance]]
- [[Comprobación de tipos]] · [[Conversión de tipos (coerción y cast)]]
- [[Manejo de errores (léxicos, sintácticos, semánticos)]]
- [[Flujo de control y switch]]
- [[vis-network]] · [[Graphviz]]
- [[DataForge]] (proyecto hermano, mismo curso, stack Java)
