# Guía de Programación — CompInterpreter

**Universidad de San Carlos de Guatemala**
**Facultad de Ingeniería — Escuela de Ciencias y Sistemas**
**Organización de Lenguajes y Compiladores 1**
**Proyecto 2 — Segundo Semestre 2024**

---

## Cómo usar esta guía

Esta guía es material de tutor: está pensada para que un instructor te acompañe mientras **vos programás** CompInterpreter desde cero, no para que la leas pasivamente. Cada tema explica el problema teórico, muestra el código real del proyecto (el que vive hoy en `CompInterpreter/server/` y `CompInterpreter/client/`) y te propone ejercicios concretos de "ahora escribí...". Los fragmentos de código citados son exactos, copiados del proyecto funcionando — cuando compiles el tuyo, el comportamiento observable debería coincidir.

La referencia técnica exhaustiva de decisiones de diseño es `docs/ManualTecnico.md` (léelo en paralelo); esta guía es su versión "paso a paso, en el orden en que se construye", no un reemplazo.

---

## Tema 1 — Preparación del proyecto

### 1.1 Por qué este stack y no Java + JFlex/CUP + JavaFX

Ya construiste (o viste construir) DataForge, ConjAnalyzer y/o CompScript con Java + JFlex + CUP + JavaFX. La teoría de compiladores no cambia un ápice acá: léxico, sintáctico, semántico, tabla de símbolos y propagación de errores son los mismos conceptos del Dragón. Lo que cambia es:

1. **La arquitectura es un requisito del enunciado, no una preferencia.** La sección 9 exige explícitamente arquitectura cliente-servidor y visualización desde el navegador. JavaFX es una tecnología de escritorio: no sirve para ese requisito.
2. **El lenguaje de implementación es JavaScript** (Node.js en el servidor, JS/JSX en el cliente) en vez de Java. Node ejecuta el intérprete y el analizador Jison; el navegador no puede correr esos módulos directamente — por eso hace falta el servidor como intermediario.
3. **Jison reemplaza a JFlex+CUP.** Es la misma clase de analizador LALR(1); la diferencia es de ergonomía de la herramienta: Jison junta la especificación léxica y sintáctica en **un solo archivo** (`grammar.jison`), mientras que JFlex y CUP son dos herramientas separadas con dos archivos.
4. **Esta vez SÍ hace falta un AST explícito.** En DataForge no había control de flujo (cada instrucción corre exactamente una vez, en el orden textual), así que la ejecución podía ir directamente en las acciones semánticas de CUP. Acá el lenguaje tiene `if`, `while`, `for`, `switch`, funciones invocables antes de su declaración textual — todas estructuras que se recorren más de una vez o en un orden distinto al textual. Eso exige construir un árbol y recorrerlo aparte.

### 1.2 Instalar el servidor

```bash
cd CompInterpreter/server
npm init -y
npm install express cors
npm install --save-dev jison
```

Tu `package.json` del servidor debe terminar viéndose así (es el real del proyecto):

```json
{
  "name": "compinterpreter-server",
  "version": "1.0.0",
  "main": "server.js",
  "type": "commonjs",
  "scripts": {
    "start": "node server.js",
    "grammar": "jison src/grammar.jison -o src/parser.js",
    "build": "npm run grammar"
  },
  "dependencies": {
    "cors": "^2.8.5",
    "express": "^4.19.2"
  },
  "devDependencies": {
    "jison": "^0.4.18"
  }
}
```

El script `grammar` es el que vas a correr **cada vez** que modifiques `grammar.jison`: Jison no regenera el parser automáticamente, y `src/parser.js` es un archivo generado que nunca vas a editar a mano (misma disciplina que con `Lexer`/`Parser`/`sym` en DataForge).

### 1.3 Instalar el cliente

```bash
cd CompInterpreter/client
npm create vite@latest . -- --template react
npm install vis-network vis-data
npm install --save-dev @playwright/test
```

`vis-network`/`vis-data` son la librería que vas a usar para dibujar el AST como grafo interactivo (Tema 7); `@playwright/test` es para la suite E2E (Tema 8).

### 1.4 Poner ambos a correr

Necesitás **dos terminales** simultáneas — a diferencia de un monolito JavaFX de un solo proceso:

```bash
# Terminal 1
cd CompInterpreter/server
npm start          # escucha en http://localhost:4000

# Terminal 2
cd CompInterpreter/client
npm run dev         # Vite, http://localhost:5173
```

Si en algún momento vas a probar solo el motor (léxico+sintáctico+semántico+ejecución) sin levantar servidor ni navegador, el proyecto real incluye `server/test-cli.js` para eso — volvé a este script en el Tema 4, es tu mejor herramienta de depuración mientras programás el intérprete.

---

## Tema 2 — La gramática Jison: tokens y `%lex`

### 2.1 Un solo archivo, dos secciones

Vos vas a crear `server/src/grammar.jison` con esta forma general:

```
%lex
%options case-insensitive flex
%%
... reglas lexicas ...
/lex

%right IF DOSPUNTOS
... declaraciones de precedencia ...
%start inicio

%%
... gramatica sintactica ...
%%
... codigo de soporte (funcion nodo()) ...
```

`%options case-insensitive` resuelve en una sola línea la sección 5.1 del enunciado (el lenguaje es insensible a mayúsculas/minúsculas) — sin ese `%options`, tendrías que escribir cada palabra reservada como una expresión regular carácter por carácter (`[iI][fF]` en vez de `if`). Es la misma clase de decisión que en JFlex, pero con ergonomía distinta.

### 2.2 Orden de las reglas: el mismo principio que en JFlex

Todo analizador léxico basado en *longest match* necesita que, ante un empate de longitud, gane la regla escrita **primero**. Esto se traduce en dos reglas de oro que ya viste en DataForge y que acá aplican igual:

**Reservadas antes que el identificador genérico.** Empezá escribiendo las palabras reservadas y dejá `ID` al final:

```
"true"                              return 'TRUE';
"false"                             return 'FALSE';
"let"                               return 'LET';
"const"                             return 'CONST';
"int"                               return 'TIPO_INT';
...
"function"                          return 'FUNCTION';
"echo"                              return 'ECHO';

/* --- identificador --- */
[a-zA-Z_][a-zA-Z0-9_]*              return 'ID';
```

**Ahora escribí vos el ejercicio:** ¿qué pasaría si movieras `[a-zA-Z_][a-zA-Z0-9_]* return 'ID';` **antes** de `"if" return 'IF';`? Respuesta: ninguna reservada se reconocería jamás — `if`, `let`, `function` pasarían a matchear como `ID` porque el motor generado encuentra la primera regla que coincide con esa longitud, y `ID` coincide con `if` igual de bien que la regla específica.

**Operadores de dos caracteres antes que los de uno.** Mismo principio:

```
"++"                                return 'INCR';
"--"                                return 'DECR';
"=="                                return 'IGUAL';
"!="                                return 'DIFERENTE';
"<="                                return 'MENORIGUAL';
">="                                return 'MAYORIGUAL';
"&&"                                return 'AND';
"||"                                return 'OR';
"+"                                 return 'MAS';
"-"                                 return 'MENOS';
...
"<"                                 return 'MENOR';
">"                                 return 'MAYOR';
"="                                 return 'ASIGNA';
```

Si pusieras `"="` antes que `"=="`, `x == y` se tokenizaría como `ASIGNA ASIGNA` en vez de `IGUAL` — un error real y común de quien no respeta este orden.

### 2.3 Literales y escapes resueltos en el propio lexer

Los literales numéricos son directos:

```
[0-9]+"."[0-9]+                     return 'DECIMAL';
[0-9]+                              return 'ENTERO';
```

(Fijate que `DECIMAL` va **antes** que `ENTERO`: si un `123.45` se tokenizara con la regla de `ENTERO` primero, el analizador reconocería `123` como `ENTERO` y dejaría `.45` sin consumir apropiadamente. En la práctica Jison ya prioriza por longitud de coincidencia, pero declarar la regla más específica primero es la práctica defensiva correcta.)

Las cadenas y caracteres resuelven sus secuencias de escape (`\n`, `\t`, `\"`, `\'`, `\\`) **dentro de la propia acción léxica**, no más adelante en el intérprete:

```
\"(\\.|[^\\"])*\"    %{
                       yytext = yytext.slice(1, -1).replace(/\\(.)/g, function (m, c) {
                         switch (c) { case 'n': return '\n'; case 't': return '\t';
                           case '"': return '"'; case "'": return "'"; case '\\': return '\\';
                           default: return c; }
                       });
                       return 'CADENA';
                    %}
```

**Por qué acá y no después:** resolver el escape es un fenómeno puramente léxico (no depende de contexto sintáctico ni semántico). Si lo dejaras para el intérprete, cada consumidor posterior de esa cadena (ejecución, reportes) tendría que repetir la misma lógica de reemplazo.

### 2.4 Recuperación de errores léxicos

El enunciado (6.1) exige que un carácter no reconocido no aborte el análisis. La regla catch-all `.` (que solo matchea si ninguna regla anterior lo hizo) empuja el error a un arreglo compartido en vez de detener el lexer:

```
.                                   %{
                                       if (yy.errores) {
                                         yy.errores.push({ tipo: 'Léxico',
                                           descripcion: 'El carácter "' + yytext + '" no pertenece al lenguaje',
                                           linea: yylloc.first_line, columna: yylloc.first_column + 1 });
                                       }
                                    %}
```

`yy` es un objeto que Jison comparte automáticamente entre el lexer y el parser — es la pieza que hace posible que **todos** los errores (léxicos, sintácticos, y después semánticos) terminen en un único arreglo ordenable por línea/columna. Vas a volver a ver `yy` en el Tema 6, cuando conectes este lexer con el resto del sistema desde `analizar.js`.

**Ejercicio:** agregá vos la regla léxica para reconocer comentarios de línea (`// ...`) y de bloque (`/* ... */`), que deben descartarse sin devolver ningún token (no llevan `return`):

```
"//".*                              /* comentario de linea */
"/*"([^*]|\*+[^*/])*\*+"/"          /* comentario multilinea */
\s+                                 /* espacios en blanco */
```

---

## Tema 3 — El parser: gramática sintáctica y construcción del AST

### 3.1 La tabla de precedencia (sección 5.10 del enunciado)

Jison declara la precedencia de menor a mayor importancia (al revés de como uno a veces la piensa):

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

El caso más preguntado en clase es `%nonassoc POTENCIA RAIZ`: no es un capricho de la herramienta, es literal del enunciado (tabla 5.10) que potencia (`^`) y raíz (`$`) **no son asociativas** — `2^3^2` sin paréntesis debe ser un error de gramática, no algo que Jison resuelva por defecto a izquierda o a derecha. Compará con `%left MAS MENOS`, donde si encadenás `2+3+4` sin ambigüedad porque la asociatividad izquierda está declarada explícitamente.

**Importante:** esta tabla de precedencia es un mecanismo de la herramienta Jison. El entregable real (`docs/gramatica.txt`) es una gramática BNF **libre de ambigüedad por construcción**, reescribiendo las expresiones en capas anidadas donde cada capa solo combina con la de abajo:

```
<expresion-or>     ::= <expresion-or> "||" <expresion-and> | <expresion-and>
<expresion-and>    ::= <expresion-and> "&&" <expresion-not> | <expresion-not>
<expresion-not>    ::= "!" <expresion-not> | <expresion-relacional>
<expresion-relacional> ::= ... | <expresion-suma>
<expresion-suma>   ::= <expresion-suma> "+" <expresion-termino> | ... | <expresion-termino>
<expresion-termino> ::= <expresion-termino> "*" <expresion-potencia> | ... | <expresion-potencia>
<expresion-potencia> ::= <expresion-unaria> "^" <expresion-unaria> | ... | <expresion-unaria>
<expresion-unaria>  ::= "-" <expresion-unaria> | <expresion-primaria>
```

No confundas una cosa con la otra: la tabla `%left`/`%right`/`%nonassoc` de Jison **no** es la gramática BNF que hay que entregar (sección 8 del enunciado exige explícitamente que no sea copia de la herramienta).

### 3.2 La fábrica de nodos: la pieza central de este tema

Cada producción de la gramática construye su propio nodo del AST llamando a una única función auxiliar:

```javascript
function nodo(tipo, props, loc) {
  var n = { tipo: tipo };
  if (props) { for (var k in props) { if (props.hasOwnProperty(k)) n[k] = props[k]; } }
  if (loc) { n.linea = loc.first_line; n.columna = loc.first_column + 1; }
  return n;
}
```

Todo nodo tiene, como mínimo, `tipo` (string) y `linea`/`columna` (para el reporte de errores), más las propiedades propias de esa construcción. Por ejemplo, un `if`:

```
if_stmt
    : IF PARIZQ expresion PARDER bloque
        { $$ = nodo('IF', { cond: $3, entonces: $5, sino: null }, @1); }
    | IF PARIZQ expresion PARDER bloque ELSE bloque
        { $$ = nodo('IF', { cond: $3, entonces: $5, sino: $7 }, @1); }
    | IF PARIZQ expresion PARDER bloque ELSE if_stmt
        { $$ = nodo('IF', { cond: $3, entonces: $5, sino: [$7] }, @1); }
    ;
```

**Por qué importa tanto esta convención:** mientras cualquier nodo nuevo respete `{ tipo, ... }`, el resto del sistema —el reporte del AST como grafo (Tema 7)— lo dibuja automáticamente, sin que nadie tenga que tocar ese código. Es el mismo argumento de "abierto a extensión, cerrado a modificación" aplicado a un árbol sintáctico.

### 3.3 Evitar conflictos shift/reduce con un prefijo compartido

Las cuatro formas de declaración (`let x: int;`, `let x: int = 5;`, `let v: int[] = [...]`, `let m: int[][] = [...]`) comparten el mismo prefijo `mutabilidad lista_ids : tipo` en la gramática, precisamente para que Jison no tenga que decidir con lookahead limitado cuál de las cuatro está viendo:

```
declaracion
    : mutabilidad lista_ids DOSPUNTOS tipo
        { $$ = nodo('DECLARACION', { mutable: $1, ids: $2, tipoDato: $4, valor: null }, @2); }
    | mutabilidad lista_ids DOSPUNTOS tipo ASIGNA expresion
        { $$ = nodo('DECLARACION', { mutable: $1, ids: $2, tipoDato: $4, valor: $6 }, @2); }
    | mutabilidad lista_ids DOSPUNTOS tipo CORIZQ CORDER ASIGNA vector_init1
        { $$ = nodo('DECLARACION_VECTOR', { mutable: $1, ids: $2, tipoDato: $4, dimension: 1, init: $8 }, @2); }
    | mutabilidad lista_ids DOSPUNTOS tipo CORIZQ CORDER CORIZQ CORDER ASIGNA vector_init2
        { $$ = nodo('DECLARACION_VECTOR', { mutable: $1, ids: $2, tipoDato: $4, dimension: 2, init: $10 }, @2); }
    ;
```

Los vectores usan `lista_ids` igual que las variables simples, aunque el lenguaje solo permite un identificador por vector — esa regla (semántica, no sintáctica) se valida después en el intérprete (`execDeclaracionVector`), no acá. Es una decisión deliberada: mantenerla sintáctica hubiera forzado una gramática separada solo para ese caso, a cambio de nada.

### 3.4 Recuperación de errores sintácticos

Igual que en CUP con `error`, Jison soporta un token especial `error` para modo pánico:

```
elemento_global
    : declaracion PUNTOCOMA          { $$ = $1; }
    | funcion                        { $$ = $1; }
    | ejecutar_stmt                  { $$ = $1; }
    | error PUNTOCOMA                { $$ = nodo('ERROR', {}); }
    ;
```

Cuando el parser no puede continuar con la entrada actual, descarta tokens hasta encontrar el próximo `;` y sigue desde ahí — el mismo principio de `instruccion ::= error PUNTO_COMA` que ya viste en DataForge, adaptado a la sintaxis de Jison.

### 3.5 Verificación aritmética barata del AST

**Ahora corré vos** (con el proyecto ya armado) esta prueba sobre un archivo con vectores:

```bash
cd CompInterpreter/server
node test-cli.js ../entradas/ejemplo_vectores2d.ci
```

La última línea imprime `AST --- nodos: 106  aristas: 105`. En cualquier árbol sin ciclos, `aristas = nodos - 1`. Es una verificación aritmética barata de que tu recorrido realmente construyó un árbol (y no un grafo con referencias cruzadas) — antes de mirar el contenido, esta cuenta ya te dice si algo está estructuralmente mal.

---

## Tema 4 — El intérprete: entorno, tipos y evaluación en dos pasadas

El módulo `server/src/interprete/interprete.js` recorre el AST haciendo, al mismo tiempo, comprobación de tipos y ejecución — no hay una fase separada de "chequeo semántico" y otra de "ejecución"; ambas ocurren en el mismo recorrido, sobre los mismos nodos.

### 4.1 El entorno como cadena de ámbitos

Empezá por `server/src/interprete/entorno.js`. Es exactamente el mismo concepto de tabla de símbolos que en DataForge, con la representación de cadena de ámbitos:

```javascript
class Entorno {
  constructor(padre, nombre) {
    this.padre = padre || null;
    this.nombre = nombre || 'global';
    this.tabla = new Map();
  }

  clave(id) { return String(id).toLowerCase(); }

  declarar(simbolo) {
    const k = this.clave(simbolo.id);
    if (this.tabla.has(k)) return false;
    this.tabla.set(k, simbolo);
    return true;
  }

  obtener(id) {
    const k = this.clave(id);
    let e = this;
    while (e) {
      if (e.tabla.has(k)) return e.tabla.get(k);
      e = e.padre;
    }
    return null;
  }
}
```

Cada ámbito (`global`, luego función/método, luego bloque dentro de `if`/`while`/`for`) es un `Entorno` con una referencia a su padre. `obtener()` sube por la cadena hasta encontrar el identificador o llegar al final. Las claves se normalizan a minúsculas porque el lenguaje es insensible a mayúsculas también en los identificadores (sección 5.1).

**Ejercicio:** ¿por qué `declarar()` devuelve `false` en vez de sobreescribir silenciosamente? Porque así el llamador (`execDeclaracion` en el intérprete) puede distinguir "ya existía en **este mismo** ámbito" (error semántico: redeclaración) de "existe en un ámbito padre" (válido: sombreado normal de variables).

### 4.2 Valor: todo dato en tiempo de ejecución va etiquetado con su tipo

`server/src/interprete/valor.js` define la representación uniforme:

```javascript
class Valor {
  constructor(tipo, valor) {
    this.tipo = tipo;      // 'int' | 'double' | 'bool' | 'char' | 'string' | 'null' | 'vector'
    this.valor = valor;    // JS number/boolean/string/null, o Array<Valor> para vectores
    this.tipoBase = null;  // solo vectores: tipo de los elementos
    this.dimension = 0;    // solo vectores: 1 o 2
  }
  static int(v) { return new Valor(TIPO.INT, Math.trunc(v)); }
  static double(v) { return new Valor(TIPO.DOUBLE, v); }
  ...
}
```

Cada operación y cada función nativa puede verificar `.tipo` antes de operar, sin necesidad de `instanceof` contra clases distintas por tipo. Es la misma idea que un `Valor`/`Object` etiquetado en cualquier intérprete didáctico.

### 4.3 Dos pasadas: por qué, y qué resuelven

El enunciado (5.21) recomienda **literalmente** recorrer el AST en dos pasadas. Así se implementa en `Interprete.interpretar()`:

```javascript
interpretar(ast) {
  if (!ast || !ast.globales) return;

  // PASADA 1: funciones y metodos
  for (const g of ast.globales) {
    if (g.tipo === 'FUNCION' || g.tipo === 'METODO') {
      const clave = g.id.toLowerCase();
      if (this.funciones.has(clave)) {
        this.errores.semantico('Ya existe una función o método con el identificador "' + g.id + '"...', g.linea, g.columna);
        continue;
      }
      this.funciones.set(clave, g);
      ...
    }
  }

  // PASADA 2: variables globales y ejecutar
  for (const g of ast.globales) {
    if (g.tipo === 'DECLARACION' || g.tipo === 'DECLARACION_VECTOR') {
      this.ejecutarInstruccion(g, this.global);
    } else if (g.tipo === 'EJECUTAR') {
      this.ejecutarEjecutar(g);
    }
  }
}
```

**Lo que esto resuelve:** referencia adelantada. `entradas/ejemplo_funciones.ci` tiene `main()` llamando a `factorial()`, declarada más abajo en el archivo. Para cuando la pasada 2 ejecuta `ejecutar main();`, la pasada 1 ya registró **todas** las funciones y métodos en el `Map` — sin importar el orden textual en que aparecen.

**Corré vos esta demo** para verlo funcionando:

```bash
node test-cli.js ../entradas/ejemplo_funciones.ci
```

Consola real:
```
Fuerza (masa=10, acel default 9.8): 98.0
Fuerza (masa=10, acel=2.5): 25.0
5! = 120
cast(18.6 as int) = 18
...
```

`factorial(5)` da `120` aunque esté declarada después de `main()` en el archivo — la prueba de que la referencia adelantada funciona de verdad.

### 4.4 Coerción implícita vs. casteo explícito

Son dos mecanismos distintos y no hay que confundirlos:

- **Coerción implícita** (`coercionar`, dentro de `Interprete`): automática al asignar, al pasar un argumento o al retornar, solo entre tipos numéricamente compatibles.
- **Casteo explícito** (`cast(expr as tipo)`): pedido a propósito por el programador, cubre combinaciones adicionales (`int` a `char` vía código ASCII, cualquier numérico a `string`).

```javascript
/* DECISION DE DISENO: los tipos numericos son mutuamente asignables con
   coercion (double->int trunca). Sin esto, el propio ejemplo oficial del
   Anexo 11.1 seria invalido: `let modulo: int = x % 3` donde `%` produce
   DECIMAL segun la tabla literal del enunciado (5.5.6). */
coercionar(tipoDest, v, l, c) {
  if (v === null) return null;
  const o = v.tipo;
  if (o === tipoDest) return v;
  switch (tipoDest) {
    case TIPO.INT:
      if (o === TIPO.DOUBLE) return Valor.int(Math.trunc(v.valor));
      if (o === TIPO.BOOL) return Valor.int(v.valor ? 1 : 0);
      break;
    ...
  }
  this.errores.semantico('No se puede asignar un valor de tipo ' + ops.nombreTipo(o) +
    ' a una variable de tipo ' + ops.nombreTipo(tipoDest), l, c);
  return null;
}
```

**Por qué la coerción tiene que existir, no es un capricho:** sin `double→int`, el propio ejemplo oficial del Anexo 11.1 sería inválido. `%` (módulo) entre dos `int` produce `DECIMAL` según la tabla 5.5.6 del enunciado (`T_MOD`), y el ejemplo asigna ese resultado a una variable `int` (`let modulo: int = x % 3;`). Sin coerción implícita, esa línea sería un error de tipos en el ejemplo que el propio enunciado da como válido.

### 4.5 Propagación de errores por `null`

Cuando una expresión falla semánticamente, la función que la evalúa devuelve `null` de JavaScript — **no confundir con el `Valor` de tipo `'null'` propio del lenguaje `.ci`**, que es una instancia válida de `Valor`. Toda operación que recibe un operando `null` de JS lo propaga en silencio:

```javascript
function aritmetica(op, izq, der, l, c, errores) {
  if (izq === null || der === null) return null;   // propagacion por null
  ...
}
```

Esto evita que un único error de causa raíz produzca una cascada de mensajes derivados sobre la misma expresión — el mismo principio que en DataForge, con la misma trampa de nombres para explicar con cuidado en clase (`null` de JS ≠ `Valor` de tipo `'null'`).

### 4.6 Tablas de compatibilidad de tipos: transcripción literal

`server/src/interprete/operaciones.js` reproduce las tablas del enunciado (5.5.1-5.5.7) como objetos JS. Fijate en el caso del módulo, documentado explícitamente porque es contraintuitivo:

```javascript
// 5.5.6 Modulo  -> NOTA: el enunciado define Entero % Entero = DECIMAL
// (identico a la tabla de Raiz). Es inusual (uno esperaria Entero), pero
// se implementa TAL CUAL lo dice la fuente, sin "corregirlo".
const T_MOD = {
  int:    { int: 'double', double: 'double' },
  double: { int: 'double', double: 'double' }
};
```

**Ejercicio:** ¿por qué documentar explícitamente una decisión "rara" en vez de simplemente implementarla? Porque quien mantenga el código después (vos mismo, en dos semanas) va a leer `int % int -> double` y va a pensar que es un bug — el comentario evita que alguien "corrija" algo que en realidad es fidelidad literal al enunciado.

---

## Tema 5 — Señales de control: break/continue/return, y el caso real de la auditoría

Este es el tema con más contenido de diseño del proyecto: cómo implementar control de flujo sin usar las excepciones nativas de JavaScript, y dos bugs reales que se encontraron y corrigieron en la auditoría de código del **21 de julio de 2026**.

### 5.1 El mecanismo: señales, no excepciones

**Preguntate primero:** ¿cómo implementarías vos `break`/`continue`/`return` en JavaScript? La respuesta más común es `throw`/`try-catch`. CompInterpreter elige otra cosa a propósito: cada instrucción devuelve, hacia la instrucción que la contiene, un objeto de señal o `null`:

```javascript
ejecutarInstruccion(nodo, entorno) {
  if (!nodo) return null;
  switch (nodo.tipo) {
    ...
    case 'BREAK': return { tipo: 'BREAK', linea: nodo.linea, columna: nodo.columna };
    case 'CONTINUE': return { tipo: 'CONTINUE', linea: nodo.linea, columna: nodo.columna };
    case 'RETURN': {
      const v = nodo.valor ? this.evaluar(nodo.valor, entorno) : null;
      return { tipo: 'RETURN', valor: nodo.valor ? v : null, tieneValor: !!nodo.valor };
    }
    ...
  }
}

ejecutarBloque(instrucciones, entorno) {
  for (const inst of instrucciones) {
    const senal = this.ejecutarInstruccion(inst, entorno);
    if (senal) return senal;
  }
  return null;
}
```

Cada estructura de control interpreta la señal según su propia semántica. Un `while`:

```javascript
execWhile(nodo, entorno) {
  let n = 0;
  while (true) {
    const cond = this.evaluar(nodo.cond, entorno);
    if (cond === null) return null;
    if (!cond.valor) break;
    if (++n > MAX_ITER) { this.errores.semantico('Posible ciclo infinito (while)', nodo.linea, nodo.columna); break; }
    const senal = this.ejecutarBloque(nodo.cuerpo, new Entorno(entorno, entorno.nombre));
    if (senal) {
      if (senal.tipo === 'BREAK') break;
      if (senal.tipo === 'CONTINUE') continue;
      return senal;   // RETURN se re-propaga hacia la funcion que llamo
    }
  }
  return null;
}
```

**Por qué un objeto de señal y no una excepción:** las excepciones funcionarían perfectamente acá — no es una limitación técnica. Es una decisión de diseño: la señal es flujo de control **explícito y visible** en cada punto donde alguien decide qué hacer con ella, mientras que una excepción viaja implícita a través de la pila de llamadas. Y, de paso, resuelve gratis otro problema.

### 5.2 El fall-through de switch "sale gratis"

```javascript
execSwitch(nodo, entorno) {
  const expr = this.evaluar(nodo.expr, entorno);
  if (expr === null) return null;
  const local = new Entorno(entorno, entorno.nombre);
  let coincidio = false;
  for (const caso of nodo.casos) {
    if (!coincidio) {
      const ce = this.evaluar(caso.expr, local);
      if (ce === null) continue;
      const igual = ops.relacional('==', expr, ce, nodo.linea, nodo.columna, this.errores);
      if (igual !== null && igual.valor === true) coincidio = true;
    }
    if (coincidio) {
      const senal = this.ejecutarBloque(caso.cuerpo, local);
      if (senal) {
        if (senal.tipo === 'BREAK') return null;
        return senal;   // RETURN / CONTINUE se propagan hacia afuera del switch
      }
    }
  }
  ...
}
```

Una vez que un `case` coincide (`coincidio = true`), el bucle sigue ejecutando los casos siguientes **sin volver a evaluar su expresión**, hasta toparse con una señal `BREAK` (que corta el `switch` devolviendo `null`) o con `RETURN`/`CONTINUE` (que se re-propagan hacia afuera). Esto implementa el *fall-through* exigido por la sección 5.16.2 del enunciado sin ningún código adicional — es la misma señal que ya usás para los ciclos.

**Corré esta demo** para verlo en consola:

```bash
node test-cli.js ../entradas/ejemplo_switch.ci
```

```
dos (sin break, cae)
tres (fall-through)
cuatro (aqui si hay break)
```

`case 2` no tiene `break`, entonces cae en `case 3` y sigue hasta `case 4`, que sí rompe.

### 5.3 Guardas de estabilidad: MAX_ITER y MAX_DEPTH

```javascript
const MAX_ITER = 1000000;   // guarda anti-bucle-infinito
const MAX_DEPTH = 2000;     // guarda anti-recursion desbordada
```

No están exigidas por el enunciado — son una decisión propia necesaria para la estabilidad del servidor. Un `loop {}` sin salida o una recursión infinita en un intérprete de escritorio solo afecta a ese proceso; acá, colgarían el **proceso Node completo**, afectando a todas las peticiones concurrentes (dos pestañas, dos estudiantes usando el mismo servidor a la vez), no solo a la que causó el problema. Volvés a este punto en el Tema 6, cuando veas por qué el servidor crea un entorno fresco por petición.

### 5.4 Caso de estudio real — Bug 1: vector con elemento incompatible tumbaba TODO el análisis

**El síntoma:** con una entrada como

```
let v: int[] = [1, "hola", 3];
```

antes de la corrección, el análisis completo abortaba con un mensaje genérico "Error interno durante la ejecución" en vez de un error semántico localizado.

**La causa raíz:** al construir un vector desde una lista literal, cada elemento se evalúa y se intenta coercionar al tipo base declarado (`int`, en el ejemplo). Cuando la coerción falla (`"hola"` no es coercionable a `int`), `coercionar()` reporta el error semántico y devuelve `null` de JavaScript — el mecanismo normal de propagación por error que ya viste en el Tema 4. El problema es que un vector **no es una expresión suelta**: sus elementos van a ser leídos después por las funciones nativas sobre vectores (`sum`, `max`, `min`, `average`) y por `aTexto()`, y esas funciones leen `.tipo` de cada elemento **sin volver a comprobar si es null**, porque asumen que todo elemento de un vector es siempre una instancia de `Valor`. Reconstruyendo la línea tal como la describe la sección 5.7 del Manual Técnico, el código anterior a la corrección hacía, en esencia:

```javascript
// ANTES (reconstruido a partir de la descripción de la auditoría):
for (const e of init.valores) {
  const ev = this.evaluar(e, entorno);
  arr.push(ev === null ? Valor.nulo() : this.coercionar(nodo.tipoDato, ev, nodo.linea, nodo.columna));
}
```

Si `ev` no era `null` (el literal `"hola"` se evalúa perfectamente bien como un `Valor` de tipo `string`) pero `coercionar()` fallaba por incompatibilidad de tipos, el resultado de `coercionar()` era `null` de JS — y ese `null` crudo quedaba guardado en el arreglo. Cuando más adelante `sum(v)` o `aTexto(v)` recorrían el arreglo y hacían `elemento.tipo`, JavaScript lanzaba `TypeError: Cannot read properties of null (reading 'tipo')` — una excepción **no controlada** que `analizar.js` atrapaba en su bloque `try/catch` general y reportaba como el mensaje genérico "Error interno durante la ejecución", abortando toda la interpretación por un solo elemento mal escrito.

**La corrección:** un método nuevo que, si la coerción falla, en vez de dejar el `null` crudo, rellena la posición con el **valor por defecto** del tipo base declarado:

```javascript
// DESPUES (código real, interprete.js):
coercionarElementoOValorDefecto(tipoBase, v, l, c) {
  const coer = this.coercionar(tipoBase, v, l, c);
  return coer !== null ? coer : new Valor(tipoBase, valorPorDefecto(tipoBase));
}
```

Y en la construcción del vector:

```javascript
} else if (init.tipo === 'VECTOR_LISTA') {
  const arr = [];
  for (const e of init.valores) {
    const ev = this.evaluar(e, entorno);
    arr.push(ev === null ? Valor.nulo() : this.coercionarElementoOValorDefecto(nodo.tipoDato, ev, nodo.linea, nodo.columna));
  }
  vectorVal = Valor.vector(nodo.tipoDato, 1, arr);
}
```

El error semántico sigue reportándose (lo reporta `coercionar()` internamente, antes de que `coercionarElementoOValorDefecto` decida qué guardar), pero ahora la posición del vector queda siempre ocupada por un `Valor` válido — nunca por un `null` crudo de JavaScript. El resto de la interpretación puede continuar con normalidad.

**La lección general, no solo de este bug:** siempre que un valor vaya a ser almacenado en una estructura que otro código va a leer después asumiendo una forma fija (acá, "todo elemento de vector es un `Valor`"), tenés que coercionar **o** rellenar con un valor por defecto — nunca dejar pasar un `null` crudo de JavaScript al runtime que sigue. La propagación por `null` es perfecta para una expresión suelta que nadie más va a tocar; es peligrosa en cuanto ese resultado se guarda en un lugar con un contrato de forma fija.

### 5.5 Caso de estudio real — Bug 2: `break`/`continue` fuera de un ciclo no generaba error

**El síntoma:** un `break;` puesto directamente en el cuerpo de una función (sin ningún ciclo alrededor) truncaba en silencio el resto de las instrucciones de esa función, sin dejar ningún rastro en la consola ni en la tabla de errores.

**La causa raíz:** el enunciado (5.18.2) exige literalmente que `continue` fuera de un ciclo sea un error ("siempre debe estar dentro de un ciclo, de lo contrario será un error"). La validación en `invocar()` (donde termina de resolverse una llamada a función o método, después de ejecutar su cuerpo) originalmente solo comprobaba `CONTINUE`:

```javascript
// ANTES (reconstruido a partir de la descripción de la auditoría):
if (senal && senal.tipo === 'CONTINUE') {
  this.errores.semantico('La sentencia "continue" debe estar dentro de un ciclo', ...);
  return null;
}
```

Si la señal que llegaba hasta acá era `BREAK` (por ejemplo, un `break;` suelto en el cuerpo de un método, sin ningún `while`/`for`/`switch` que la hubiera interceptado antes), ese `if` no se cumplía, y la ejecución simplemente seguía de largo sin reportar nada — el `break` se "tragaba" en silencio.

**La corrección:** extender la misma validación a ambas señales, por el mismo criterio del enunciado:

```javascript
// DESPUES (código real, interprete.js):
if (senal && (senal.tipo === 'BREAK' || senal.tipo === 'CONTINUE')) {
  const cual = senal.tipo === 'CONTINUE' ? 'continue' : 'break';
  this.errores.semantico('La sentencia "' + cual + '" debe estar dentro de un ciclo',
    senal.linea != null ? senal.linea : linea, senal.columna != null ? senal.columna : columna);
  return null;
}
```

El enunciado solo menciona `continue` explícitamente, pero el mismo criterio se aplicó a `break` por consistencia y porque el síntoma anterior (truncar en silencio) era peor que el que exige corregir el enunciado.

**La lección que conecta los dos bugs:** ninguno de los dos falla en el "camino feliz" — los ejemplos de prueba que ya existían antes de la auditoría pasaban sin problema. Ambos son fallas del **manejo de errores**, que solo se manifiestan con una entrada específica que ningún caso de prueba anterior ejercitaba a propósito. Y el bug 2 es, en cierto sentido, peor que el bug 1: un error que aborta con un mensaje (aunque sea genérico) al menos te avisa que algo salió mal; un error que se traga en silencio hace que el programa simplemente deje de comportarse como se esperaba, sin ninguna pista de por qué.

---

## Tema 6 — El servidor Express: contrato REST y entorno fresco

### 6.1 `server.js`: transporte, nada más

```javascript
const express = require('express');
const cors = require('cors');
const { analizar } = require('./src/analizar');

const app = express();
const PUERTO = process.env.PORT || 4000;

app.use(cors());
app.use(express.json({ limit: '5mb' }));

app.get('/', (req, res) => {
  res.json({ nombre: 'CompInterpreter', estado: 'ok', endpoint: 'POST /interpretar { codigo }' });
});

app.get('/salud', (req, res) => res.json({ estado: 'ok' }));

app.post('/interpretar', (req, res) => {
  const codigo = (req.body && typeof req.body.codigo === 'string') ? req.body.codigo : '';
  try {
    const resultado = analizar(codigo);
    res.json(resultado);
  } catch (e) {
    res.status(500).json({
      errores: [{ tipo: 'Interno', descripcion: e.message || String(e), linea: 0, columna: 0 }],
      consola: '', simbolos: [], ast: { nodes: [], edges: [] }, dot: ''
    });
  }
});

app.listen(PUERTO, () => {
  console.log('CompInterpreter server escuchando en http://localhost:' + PUERTO);
});
```

Fijate qué corto es: `server.js` no tiene casi lógica propia a propósito. Separa el protocolo de transporte (HTTP, JSON, CORS) del análisis y la ejecución, que viven en `analizar.js` y `interprete/`. Esa separación es la que te permite invocar `analizar()` directo desde una consola (`test-cli.js`) sin levantar ningún servidor — la usaste todo el Tema 4 y 5.

### 6.2 `analizar.js`: el orquestador, y el canal `yy` compartido

```javascript
function analizar(codigo) {
  const errores = new ListaErrores();
  const parser = obtenerParser();

  parser.yy = { errores: errores.errores };

  parser.yy.parseError = function (msg, hash) {
    let descripcion;
    if (hash && hash.token) {
      const encontrado = hash.text !== undefined && hash.text !== '' ? ('"' + hash.text + '"') : hash.token;
      descripcion = 'Se encontró ' + encontrado;
      if (hash.expected && hash.expected.length) {
        descripcion += ' y se esperaba: ' + hash.expected.slice(0, 6).join(', ');
      }
    } else {
      descripcion = msg;
    }
    const l = hash && hash.loc ? hash.loc.first_line : (hash ? hash.line + 1 : 0);
    const c = hash && hash.loc ? hash.loc.first_column + 1 : 0;
    errores.sintactico(descripcion, l, c);
    if (!hash || !hash.recoverable) {
      throw new Error(descripcion);
    }
  };
  parser.parseError = parser.yy.parseError;

  let ast = null;
  try {
    ast = parser.parse(codigo);
  } catch (e) {
    if (!errores.hay()) errores.sintactico(e.message || 'Error de análisis', 0, 0);
    ast = ast || null;
  }

  const interprete = new Interprete(errores);
  if (ast) {
    try {
      interprete.interpretar(ast);
    } catch (e) {
      errores.semantico('Error interno durante la ejecución: ' + (e.message || e), 0, 0);
    }
  }
  ...
}
```

Volvé al Tema 2: `parser.yy = { errores: errores.errores }` **antes** de `parser.parse()` es lo que conecta el arreglo compartido con el que el lexer ya sabía empujar errores léxicos. Si te olvidás de esta línea (u la ponés después de `parser.parse()`), el lexer no tiene dónde reportar — es uno de los errores más fáciles de cometer al conectar Jison con el resto del sistema.

Notá también el `try/catch` alrededor de `interprete.interpretar(ast)`: es la última red de seguridad — si un bug como el del Tema 5 (Bug 1) volviera a colarse, en vez de tumbar el proceso Node completo, se reporta como un error semántico genérico y la petición HTTP igual responde con un JSON válido.

### 6.3 El contrato JSON: la frontera real entre servidor y cliente

```
POST /interpretar
Cuerpo:  { "codigo": "<código fuente .ci>" }

Respuesta:
{
  "errores":        Array<{ tipo, descripcion, linea, columna }>,
  "consola":        String,
  "consolaLineas":  Array<String>,
  "simbolos":       Array<{ id, categoria, tipoDato, entorno, valor, linea, columna }>,
  "ast":            { nodes: Array<{ id, label }>, edges: Array<{ from, to, label }> },
  "dot":            String
}
```

El cliente (Tema 7) **no sabe nada** de Jison, de las dos pasadas, ni de las señales de control — solo conoce la forma de este JSON. Es el mismo tipo de desacoplamiento que hay entre un compilador real y su IDE.

### 6.4 Entorno fresco por petición: dos razones, no una

```javascript
const interprete = new Interprete(errores);
```

Cada `POST /interpretar` construye su propia `ListaErrores` y su propio `Interprete` (con su propio `Entorno` global) desde cero. En DataForge, "entorno fresco" existía solo para que los reportes correspondieran al último análisis (enunciado §6). **Acá hay una segunda razón que no existe en un monolito de escritorio de un solo usuario:** un servidor puede recibir dos peticiones simultáneas (dos pestañas del navegador, dos estudiantes usando el mismo servidor). Si el `Interprete` fuera una única instancia compartida entre peticiones, dos ejecuciones concurrentes se contaminarían entre sí — variables de una aparecerían en el reporte de símbolos de la otra.

### 6.5 CORS, en una frase

`app.use(cors())` existe porque, durante el desarrollo, el cliente (`http://localhost:5173`, servido por Vite) y el servidor (`http://localhost:4000`) son **orígenes distintos** para el navegador. Sin `cors()`, el navegador bloquearía el `fetch` del cliente hacia el servidor aunque ambos corran en la misma máquina — es una política del navegador, no del servidor ni de la red.

**Ejercicio con curl** (para ver el contrato real en acción):

```bash
curl -X POST http://localhost:4000/interpretar \
  -H "Content-Type: application/json" \
  -d '{"codigo":"let x:int=1; ejecutar main(); function void main(){echo x;}"}'
```

Vas a recibir un JSON con `consola: "1"` y un `ast` con unos pocos nodos — la forma exacta del contrato que el cliente consume.

---

## Tema 7 — El cliente React: editor multipestaña y visualización

### 7.1 `api.js`: el único punto de contacto con el servidor

```javascript
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:4000';

export async function interpretar(codigo) {
  const res = await fetch(`${API_URL}/interpretar`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ codigo }),
  });
  if (!res.ok) {
    throw new Error('El servidor respondió ' + res.status + '. ¿Está corriendo en ' + API_URL + '?');
  }
  return res.json();
}
```

Toda la comunicación con el servidor pasa por esta única función. `VITE_API_URL` es la variable de entorno que le permite al cliente apuntar a un servidor en otro puerto u otra máquina sin recompilar nada más.

### 7.2 `App.jsx`: el estado global de la aplicación

`App.jsx` mantiene los archivos abiertos como un arreglo de objetos `{ id, nombre, contenido, sinGuardar }`:

```javascript
function archivoNuevo(nombre, contenido) {
  return { id: nuevoId(), nombre, contenido, sinGuardar: false };
}

const [archivos, setArchivos] = useState(() => [
  archivoNuevo('principal.ci', EJEMPLO_ANEXO),
]);
```

Al ejecutar, el resultado decide a qué pestaña de reporte saltar automáticamente:

```javascript
const manejarEjecutar = async () => {
  if (!archivoActivo) return;
  setEjecutando(true);
  setErrorRed(null);
  try {
    const res = await interpretar(archivoActivo.contenido);
    setResultado(res);
    setPanel(res.errores && res.errores.length ? 'errores' : 'consola');
  } catch (e) {
    setErrorRed(e.message || String(e));
  } finally {
    setEjecutando(false);
  }
};
```

Esta decisión de UX —saltar a "Errores" si hubo alguno, quedarse en "Consola" si no— parece un detalle chico, pero reaparece verificada literalmente por la suite de Playwright del Tema 8: no es un capricho visual, es un comportamiento que se prueba automáticamente.

**Ejercicio:** ¿por qué `erroresPorLinea` se calcula con `useMemo` en vez de recalcularse en cada render?

```javascript
const erroresPorLinea = useMemo(() => {
  const set = new Set();
  if (resultado && resultado.errores) {
    resultado.errores.forEach((e) => set.add(e.linea));
  }
  return set;
}, [resultado]);
```

Porque ese `Set` solo necesita recalcularse cuando cambia `resultado` (después de ejecutar), no en cada tecla que el usuario escribe en el editor — evita reconstruir la estructura en cada render sin ningún cambio real de datos.

### 7.3 `Editor.jsx`: `forwardRef` + `useImperativeHandle`

Este es el patrón más particular del cliente. El componente padre (`App.jsx`) necesita poder decirle al editor "andá a la línea N" cuando el usuario hace clic en una fila de la tabla de errores o de símbolos — sin que esa acción puntual (mover el cursor, hacer scroll) tenga que modelarse como parte del estado normal de props.

```jsx
const Editor = forwardRef(function Editor({ value, onChange, erroresPorLinea }, ref) {
  const textareaRef = useRef(null);
  const gutterRef = useRef(null);
  const [linea, setLinea] = useState(1);
  const [scrollTop, setScrollTop] = useState(0);

  useImperativeHandle(ref, () => ({
    irALinea(num) {
      const ta = textareaRef.current;
      if (!ta) return;
      const offsets = value.split('\n').reduce(
        (acc, l) => { acc.push(acc[acc.length - 1] + l.length + 1); return acc; },
        [0]
      );
      const pos = offsets[Math.max(0, num - 1)] || 0;
      ta.focus();
      ta.setSelectionRange(pos, pos);
      setLinea(num);
      const destino = (num - 1) * LINE_HEIGHT - ta.clientHeight / 2;
      ta.scrollTop = Math.max(0, destino);
      setScrollTop(ta.scrollTop);
      if (gutterRef.current) gutterRef.current.scrollTop = ta.scrollTop;
    },
  }));
  ...
});
```

Y en `App.jsx`, el padre simplemente lo invoca a través de la `ref`:

```javascript
const editorRef = useRef(null);
const irALinea = (num) => {
  if (editorRef.current) editorRef.current.irALinea(num);
};
// ...
<Editor ref={editorRef} value={...} onChange={...} erroresPorLinea={...} />
```

**Por qué `forwardRef` en vez de una prop más:** "ir a una línea" es una acción puntual, no un cambio de estado que deba disparar un re-render de toda la jerarquía. Exponer un método vía `ref` permite dispararlo justo cuando ocurre el clic, sin acoplar el estado interno del editor (posición del cursor, scroll) al estado del componente padre.

El **gutter** (columna de números de línea) es un `<div>` aparte, sincronizado a mano con el `scrollTop` del `<textarea>` — HTML no tiene numeración de línea nativa en un textarea, así que hay que armarlo:

```javascript
const sincronizarScroll = (e) => {
  const top = e.target.scrollTop;
  setScrollTop(top);
  if (gutterRef.current) gutterRef.current.scrollTop = top;
};
```

### 7.4 Los reportes: Consola, TablaErrores, TablaSimbolos

Los tres siguen el mismo patrón: reciben datos ya calculados por el servidor y los renderizan; `TablaErrores` y `TablaSimbolos` además reciben `onIrALinea` y lo invocan al hacer clic en una fila:

```jsx
// TablaErrores.jsx
{errores.map((e, i) => (
  <tr key={i} className={'fila-error fila-error-' + e.tipo.toLowerCase()}
      onClick={() => onIrALinea && onIrALinea(e.linea)}>
    <td>{i + 1}</td><td>{e.tipo}</td><td>{e.descripcion}</td><td>{e.linea}</td><td>{e.columna}</td>
  </tr>
))}
```

La clase `'fila-error-' + e.tipo.toLowerCase()` (`fila-error-léxico`, `fila-error-sintáctico`, `fila-error-semántico`) es la que la suite de Playwright usa para contar cuántos errores hay de cada tipo sin depender del texto exacto del mensaje.

### 7.5 `AstGrafo.jsx`: dibujar, no reconstruir

El punto clave de compiladores acá no son los detalles de React — es que el cliente **no reconstruye el AST**. Solo recibe la estructura `{ nodes, edges }` que el servidor ya armó (Tema 6, sección "Reporte del AST como grafo" en `ManualTecnico.md` §6) y la dibuja con `vis-network`:

```jsx
const options = {
  layout: {
    hierarchical: {
      enabled: true,
      direction: 'UD',
      sortMethod: 'directed',
      levelSeparation: 90,
      nodeSpacing: 150,
    },
  },
  ...
  physics: false,
  ...
};

const nodes = new DataSet(ast.nodes.map((n) => ({ id: n.id, label: n.label })));
const edges = new DataSet(ast.edges.map((e, i) => ({ id: i, from: e.from, to: e.to, label: e.label })));
networkRef.current = new Network(containerRef.current, { nodes, edges }, options);
```

`physics: false` es una decisión simple: un AST ya es jerárquico por definición (el layout `hierarchical` ya calcula la disposición de arriba hacia abajo); simular física de repulsión/atracción entre nodos sería una animación de acomodo completamente innecesaria para una estructura que ya sabe cómo ordenarse.

---

## Tema 8 — Pruebas E2E con Playwright

### 8.1 Qué prueba Playwright que `test-cli.js` no prueba

`test-cli.js` (que usaste en los Temas 4 y 5) prueba solo el motor de análisis y ejecución, sin servidor HTTP ni interfaz. Playwright levanta el **sistema completo integrado**: servidor Express real + cliente Vite real + navegador real, y valida que el `fetch` llegue, que el JSON se parsee, y que el DOM refleje el resultado. Son formas de prueba complementarias, no redundantes.

### 8.2 `playwright.config.js`: levantar ambos procesos automáticamente

```javascript
export default defineConfig({
  testDir: './e2e',
  timeout: 30000,
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'retain-on-failure',
  },
  webServer: [
    { command: 'node server.js', cwd: '../server', url: 'http://localhost:4000/salud',
      timeout: 20000, reuseExistingServer: true },
    { command: 'npm run dev -- --port 5173 --strictPort', cwd: '.', url: 'http://localhost:5173',
      timeout: 20000, reuseExistingServer: true },
  ],
});
```

`GET /salud` (que viste en el Tema 6) reaparece acá con su propósito real: Playwright lo usa para saber cuándo el servidor terminó de arrancar, antes de dejar correr los tests — evita condiciones de carrera donde el primer test se ejecuta contra un servidor que todavía no escucha. `reuseExistingServer: true` significa que si ya tenés el servidor y el cliente corriendo en tus propias terminales (Tema 1), Playwright los reutiliza en vez de levantar copias nuevas.

### 8.3 Los 3 tests reales

```javascript
test('ejecuta el ejemplo del anexo y muestra consola/símbolos/AST', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('button', { name: '▶ Ejecutar' }).click();
  await expect(page.locator('.pestana.es-activa')).toHaveText('Consola');
  const consola = page.locator('.consola');
  await expect(consola).toContainText('X es mayor que Y y no es cero');
  ...
});

test('reporta errores léxicos/sintácticos/semánticos y salta a la línea', async ({ page }) => {
  ...
  await textarea.fill(EJEMPLO_ERRORES);
  await page.getByRole('button', { name: '▶ Ejecutar' }).click();
  await expect(page.locator('.pestana.es-activa')).toHaveText(/Errores/);
  const filas = page.locator('.tabla-reporte tbody tr');
  await expect(filas).toHaveCount(8);
  await expect(page.locator('.fila-error-léxico')).toHaveCount(1);
  await expect(page.locator('.fila-error-sintáctico')).toHaveCount(2);
  await expect(page.locator('.fila-error-semántico')).toHaveCount(5);

  await filas.first().click();
  await expect(page.locator('.editor-gutter-linea.es-actual')).toHaveText('5');
});

test('nuevo archivo agrega una pestaña editable', async ({ page }) => {
  ...
});
```

El segundo test reutiliza `entradas/ejemplo_errores.ci`, el mismo archivo que ya usaste manualmente contra `test-cli.js` en el Tema 2 (donde confirmaste el desglose exacto: 1 léxico + 2 sintácticos + 5 semánticos = 8). Reutilizarlo en vez de inventar un caso nuevo es deliberado: ya se conoce el resultado exacto esperado de haberlo corrido contra el motor aislado, así que este test confirma que ese mismo resultado, ya auditado, también llega bien hasta la interfaz — sin introducir una fuente de verdad nueva sin verificar.

**Corré la suite vos mismo:**

```bash
cd CompInterpreter/client
npx playwright test
```

Deberías ver `3 passed`. Con `--headed` podés ver el navegador abriéndose y haciendo los clics automáticamente.

---

## Tema 9 — Errores comunes reales al programar cada tema

Esta lista junta los errores que de verdad aparecen al construir un proyecto de esta forma — no son hipotéticos:

1. **`ID` declarado antes que las palabras reservadas en el lexer** (Tema 2): ninguna reservada se reconoce jamás; todo pasa a ser identificador. Siempre reservadas primero.
2. **Operador de un carácter declarado antes que su versión de dos caracteres** (Tema 2): `x == y` se tokeniza como dos `ASIGNA` en vez de un `IGUAL`. Los de dos caracteres van primero.
3. **Confundir la tabla de precedencia de Jison (`%left`/`%right`/`%nonassoc`) con la gramática BNF entregable** (Tema 3): son cosas distintas; el enunciado exige explícitamente que la BNF no sea copia de la herramienta.
4. **Olvidar conectar `parser.yy` antes de `parser.parse()`** (Tema 6): si `parser.yy = { errores: ... }` se asigna después de llamar a `parse()`, o no se asigna, el lexer no tiene dónde reportar sus errores léxicos — se pierden en silencio.
5. **Perder los errores acumulados al relanzar la excepción de `parseError`**: hay que registrar el error en el arreglo compartido **antes** de hacer `throw`, no después.
6. **No diferenciar el `null` de JavaScript (señal interna de error) del `Valor` de tipo `'null'` del lenguaje** (Tema 4): son cosas distintas. El de JS significa "ya hubo un error, no reportes otro"; el de `.ci` es un dato normal que el programador puede asignar y comparar a propósito.
7. **Dejar pasar un `null` crudo de JavaScript a una estructura que otro código va a leer asumiendo una forma fija** (Tema 5, Bug 1): la lección central de este tema. La propagación por `null` es correcta para una expresión suelta; es peligrosa en cuanto ese resultado se **guarda** — en un vector, en una tabla de símbolos, en cualquier lugar que otro código vaya a recorrer después sin volver a chequear. La regla práctica: siempre coercioná el valor al tipo esperado, o rellená con el valor por defecto de ese tipo; nunca dejes que ese `null` crudo llegue al resto del runtime.
8. **Validar una señal de control (`BREAK`/`CONTINUE`) solo parcialmente** (Tema 5, Bug 2): si el enunciado exige una validación para `continue` fuera de un ciclo, pensá si el mismo criterio no debería aplicarse también a `break` — y a cualquier otra señal simétrica. Un caso "no mencionado explícitamente" no es lo mismo que un caso "que no puede pasar".
9. **Compartir el mismo `Interprete`/`Entorno` entre peticiones del servidor** (Tema 6): en un servidor con más de un usuario concurrente, esto contamina el reporte de símbolos y errores de una ejecución con el estado de otra. Siempre entorno fresco por petición.
10. **Modelar una acción puntual de UI (como "saltar a una línea") como estado elevado al padre** (Tema 7) en vez de exponer un método imperativo con `forwardRef`+`useImperativeHandle`: funciona, pero fuerza un re-render de más componentes de los necesarios por un cambio que no es realmente parte del flujo de datos normal de la aplicación.
11. **No poner guardas anti-ciclo-infinito/anti-recursión en un intérprete que corre en un servidor compartido** (Tema 5): un error que en un programa de escritorio de un solo usuario solo cuelga esa ventana, en un servidor Node de un solo proceso cuelga a todos los usuarios conectados en ese momento.
