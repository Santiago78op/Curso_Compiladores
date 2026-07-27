# Manual Técnico — VLangCherry

Universidad de San Carlos de Guatemala
Facultad de Ingeniería — Escuela de Ingeniería en Ciencias y Sistemas
Organización de Lenguajes y Compiladores 2
Proyecto 1 — V-Lang Cherry (VLangCherry)
Proyecto grupal (3 integrantes)

## Índice

1. Introducción
2. Lenguaje de programación y herramientas utilizadas
3. Estructura del proyecto
4. Arquitectura general (cliente-servidor REST)
5. Gramática y generación del analizador léxico/sintáctico (ANTLR)
6. El AST propio (`internal/ast`)
7. El traductor: parse tree → AST (`internal/traductor`)
8. El intérprete (`internal/runtime`)
9. Manejo de errores
10. El servidor REST (`cmd/servidor`, `internal/analizar`)
11. El cliente web (`client/`)
12. Decisiones de diseño no evidentes
13. Compilación y ejecución
14. Mantenimiento futuro

---

## 1. Introducción

Este manual describe la arquitectura técnica del proyecto VLangCherry, con el objetivo de que un desarrollador distinto al autor original pueda comprender su estructura, dar mantenimiento al código o extender sus funcionalidades. Se documentan el lenguaje y las herramientas empleadas, los paquetes reales del código fuente, las estructuras de datos y funciones más relevantes (citando sus firmas reales) y las decisiones de diseño tomadas frente a inconsistencias del enunciado del proyecto.

## 2. Lenguaje de programación y herramientas utilizadas

| Herramienta | Uso en el proyecto |
|---|---|
| **ANTLR4** | Generador del analizador léxico y sintáctico a partir de la especificación `grammar/VLangCherry.g4`. Se invoca con destino Go (`-Dlanguage=Go -visitor`), generando el paquete `internal/parser`. |
| **Go** (módulo `vlangcherry`, `go 1.25.3`) | Lenguaje de implementación de todo el intérprete, el AST propio y el servidor. Es el requisito obligatorio del enunciado (sección 2.3: "toda la lógica del compilador debe estar implementada en GO"). |
| **`github.com/antlr4-go/antlr/v4`** | Biblioteca en tiempo de ejecución de ANTLR para Go; provee `InputStream`, `CommonTokenStream`, `DefaultErrorListener`, entre otros, usados por el lexer y el parser generados. |
| **`net/http`** (librería estándar de Go) | Servidor REST del proyecto. Se optó deliberadamente por no incorporar un framework HTTP externo. |
| **React 19 + Vite** | Cliente web: editor de código, ejecución y visualización de reportes. |
| **vis-network / vis-data** | Renderizado interactivo del reporte de AST como grafo, en el cliente. |
| **Playwright** | Pruebas end-to-end sobre el cliente (`client/e2e/vlangcherry.spec.js`). |

## 3. Estructura del proyecto

```
VLangCherry/
├── grammar/
│   └── VLangCherry.g4          (gramática ANTLR4, fuente de verdad del léxico/sintaxis)
├── tools/
│   └── antlr.jar                (herramienta de generación)
├── docs/
│   ├── gramatica.txt             (gramática BNF entregable + decisiones de diseño)
│   ├── ManualUsuario.md
│   └── ManualTecnico.md
├── entradas/                     (6 casos de prueba .vch + 1 de cobertura de gramática)
├── server/
│   ├── go.mod                    (module vlangcherry)
│   ├── cmd/
│   │   ├── cli/                  (ejecución de un archivo .vch por línea de comandos)
│   │   ├── servidor/              (servidor REST, punto de entrada de producción)
│   │   ├── asttest/, parsetest/   (utilidades de verificación puntual)
│   └── internal/
│       ├── parser/                (100% GENERADO por ANTLR — no editar a mano)
│       ├── ast/                   (tipos de nodo del AST propio + grafo por reflection)
│       ├── traductor/             (parse tree de ANTLR -> AST propio)
│       └── runtime/                (modelo de valores, tabla de símbolos, intérprete,
│                                     funciones nativas, manejo de errores)
└── client/                       (aplicación React + Vite)
    ├── src/
    │   ├── App.jsx                 (estado de archivos abiertos, panel de resultados)
    │   ├── api.js                  (fetch a POST /interpretar)
    │   └── components/             (Toolbar, FileTabs, Editor, Consola,
    │                                 TablaErrores, TablaSimbolos, AstGrafo)
    └── e2e/vlangcherry.spec.js     (pruebas Playwright)
```

El paquete `internal/parser` se regenera con el comando de la sección 13; sus cuatro archivos (`vlangcherry_lexer.go`, `vlangcherry_parser.go`, `vlangcherry_base_visitor.go`, `vlangcherry_visitor.go`) **nunca se editan manualmente**: cualquier cambio de léxico o sintaxis se realiza en `grammar/VLangCherry.g4`.

## 4. Arquitectura general (cliente-servidor REST)

El sistema sigue una arquitectura cliente-servidor, con el mismo contrato de comunicación (en su forma) que otros intérpretes web del curso desarrollados con Node/Jison, aunque aquí el backend está implementado íntegramente en Go:

```
código fuente (.vch, String)
   → Lexer (ANTLR, generado)      — reconoce tokens; reporta errores léxicos
   → Parser (ANTLR, generado)     — valida la gramática (LL adaptativo); reporta
                                     errores sintácticos; produce un parse tree
   → Traductor                    — convierte el parse tree en un AST propio
                                     (paquete ast), desacoplado de los tipos
                                     generados por ANTLR
   → Interprete (dos pasadas)     — registra structs/funciones/métodos, declara
                                     variables globales, ejecuta main(); acumula
                                     errores semánticos, consola y tabla de símbolos
   → Resultado (JSON)             — { errores, consola, consolaLineas, simbolos,
                                     ast, dot }
```

El servidor expone este pipeline mediante un único endpoint REST (sección 10); el cliente web lo consume y presenta los resultados en paneles independientes.

## 5. Gramática y generación del analizador léxico/sintáctico (ANTLR)

La especificación completa está en `grammar/VLangCherry.g4`. A diferencia de otras herramientas del curso que separan el generador léxico del sintáctico (JFlex + CUP), ANTLR describe ambos en un solo archivo. El parser generado es **LL adaptativo** (algoritmo ALL(\*)), no LALR: decide la alternativa correcta examinando la cantidad de entrada que haga falta.

El orden de las alternativas de la regla recursiva `expr` codifica la precedencia de operadores (de mayor a menor), siguiendo la tabla del enunciado (sección 4.6):

```antlr
expr
    : expr '(' listaArgumentos? ')'            # exprLlamada
    | expr '[' expr ']'                        # exprIndice
    | expr '.' ID                              # exprCampo
    | '(' expr ')'                             # exprParentesis
    | op=('!'|'-') expr                        # exprUnario
    | expr op=('*'|'/'|'%') expr               # exprMultiplicativa
    | expr op=('+'|'-') expr                   # exprAditiva
    | expr op=('<'|'<='|'>='|'>') expr         # exprRelacional
    | expr op=('=='|'!=') expr                 # exprIgualdad
    | expr '&&' expr                           # exprAnd
    | expr '||' expr                           # exprOr
    | literalSlice                             # exprSliceLit
    | literalStruct                            # exprStructLit
    | literal                                  # exprLiteral
    | ID                                       # exprIdentificador
    ;
```

Cada alternativa etiquetada con `# nombre` hace que ANTLR genere un **tipo de contexto Go concreto** (`ExprMultiplicativaContext`, `ExprLlamadaContext`, ...), en vez de un único tipo genérico de contexto para toda la regla. Esta característica es la que habilita el diseño del traductor descrito en la sección 7: no se necesita implementar el patrón *Visitor* completo (las decenas de métodos de `BaseVLangCherryVisitor`), sino que basta un `switch` sobre el tipo dinámico del contexto.

Los terminales léxicos relevantes (fragmento del `.g4`):

```antlr
CADENA
    : '"' (ESC | ~["\\\r\n])* '"'
    ;

RUNE
    : '\'' (ESC | ~['\\\r\n]) '\''
    ;

fragment ESC
    : '\\' ["\\ntr]
    ;

COMENTARIO_LINEA
    : '//' ~[\r\n]* -> skip
    ;

COMENTARIO_BLOQUE
    : '/*' .*? '*/' -> skip
    ;
```

La gramática BNF entregable (con las decisiones de diseño frente a las inconsistencias del enunciado documentadas punto por punto) está en `docs/gramatica.txt`; no es una transcripción mecánica del `.g4`.

## 6. El AST propio (`internal/ast`)

El proyecto no opera directamente sobre el parse tree que produce ANTLR: lo traduce a un AST propio, definido en `internal/ast/ast.go`. Todo nodo implementa:

```go
type Nodo interface {
    Tipo() string
    Pos() (int, int)
}
```

Se definen alrededor de 25 tipos de nodo concretos (`Programa`, `DeclStruct`, `DeclFuncion`, `Bloque`, `DeclVariable`, `Asignacion`, `SentenciaIf`, `SentenciaSwitch`, `SentenciaFor`, `ExprBinaria`, `ExprLlamada`, `LiteralStruct`, `LiteralSlice`, etc.), todos embebiendo un struct `base{Linea, Columna int}` que resuelve `Pos()`.

El reporte de AST (`internal/ast/grafo.go`) se construye recorriendo **cualquier** `Nodo` mediante reflection (paquete `reflect`), sin necesidad de mantenimiento adicional cuando se agrega un tipo de nodo nuevo:

```go
func ConstruirGrafo(raiz Nodo) Grafo
```

La función distingue, campo por campo de cada struct, los siguientes casos:

- **Punteros/interfaces que implementan `Nodo`**: generan un nodo hijo y una arista etiquetada con el nombre del campo.
- **Slices**: cada elemento se procesa individualmente, con el nombre de campo indexado (`Campos[0]`, `Campos[1]`, ...).
- **Structs "contenedor transparente"** (`RamaIf`, `CasoSwitch`, `Parametro`, `CampoStruct`, `CampoValorLiteral`): no implementan `Nodo` por sí mismos, así que sus campos se recorren propagando un nombre compuesto (por ejemplo, `Ramas[0].Condicion`).
- **`TipoAST`**: se serializa como texto (`ta.String()`) y se agrega a la etiqueta del nodo padre, en vez de crear un nodo de grafo aparte.
- **Valores escalares**: se agregan como texto (`"campo: valor"`) a la etiqueta del nodo.

El resultado (`Grafo{Nodes []NodoGrafo, Edges []AristaGrafo}`) se serializa a JSON y viaja al cliente, que lo renderiza con **vis-network** (mismo formato `{nodes, edges}` de otros proyectos web del curso). También se genera la representación equivalente en formato **Graphviz DOT** (`ADot(g Grafo) string`), útil para graficar del lado del servidor si se requiere.

Este patrón — recorrido genérico del AST por reflection para construir un reporte visual — es análogo al empleado en `reportes/ast-grafo.js` de otro proyecto del curso desarrollado en JavaScript; aquí se logra el mismo resultado usando la reflection nativa de Go en vez de la introspección dinámica de objetos de JavaScript.

## 7. El traductor: parse tree → AST (`internal/traductor`)

`Traducir(ctx parser.IProgramaContext) *ast.Programa` es el punto de entrada. El comentario de cabecera del archivo documenta la decisión de diseño central:

```go
// Package traductor convierte el parse tree que genera ANTLR en el AST
// propio de internal/ast. No usa el patron Visitor generado: como cada
// alternativa etiquetada de la gramatica produce un tipo de contexto Go
// concreto, alcanza con un type-switch directo (mas corto que implementar
// las ~60 firmas de BaseVLangCherryVisitor).
```

Ejemplo representativo, la función `traducirExpr`:

```go
func traducirExpr(ctx parser.IExprContext) ast.Nodo {
    switch ec := ctx.(type) {
    case *parser.ExprLlamadaContext:
        l, c := pos(ec)
        n := &ast.ExprLlamada{Callee: traducirExpr(ec.Expr())}
        n.Linea, n.Columna = l, c
        if la := ec.ListaArgumentos(); la != nil {
            n.Argumentos = traducirListaArgumentos(la.(*parser.ListaArgumentosContext))
        }
        return n
    case *parser.ExprMultiplicativaContext:
        return traducirBinaria(ec, ec.AllExpr(), ec.GetOp().GetText())
    // ... un case por cada alternativa etiquetada de la regla expr
    }
    return nil
}
```

Un detalle de diseño reutilizado en varios puntos del traductor: la función `traducirLugar` (para el lado izquierdo de una asignación) y `traducirExpr` (para lectura) producen exactamente los **mismos tipos de nodo** (`ast.Identificador`, `ast.ExprIndice`, `ast.ExprCampo`) para las mismas construcciones sintácticas (`ID`, `lugar[expr]`, `lugar.campo`). Es el intérprete, no el traductor, quien decide más adelante si necesita el valor (`evaluar`) o un puntero mutable a la celda (`resolverLugar`), según si el nodo aparece del lado derecho o izquierdo de una asignación.

La posición (línea/columna) de cada nodo se obtiene con un helper compartido:

```go
func pos(ctx antlr.ParserRuleContext) (int, int) {
    t := ctx.GetStart()
    return t.GetLine(), t.GetColumn() + 1
}
```

## 8. El intérprete (`internal/runtime`)

### 8.1 Modelo de valores (`valores.go`)

```go
type Valor struct {
    Tipo   Tipo
    I      int64
    F      float64
    S      string
    B      bool
    R      rune
    Slice  *SliceVal
    Struct *StructVal
}
```

El comentario del propio archivo documenta por qué este diseño resuelve, sin mecanismos adicionales, el requisito de paso por referencia de structs y slices (enunciado, sección 7.1):

```go
// Valor es el valor en tiempo de ejecucion de una expresion VLangCherry.
// Slice y Struct viajan por puntero: asi una copia del Valor (paso de
// parametro, asignacion) comparte la misma identidad subyacente, que es
// exactamente el "por referencia" que pide la seccion 7 para structs y
// slices (los primitivos, en cambio, se copian por valor via los campos
// escalares I/F/S/B/R).
```

Como `Slice` y `Struct` ya son punteros dentro de `Valor`, copiar un `Valor` (lo que Go hace automáticamente al pasar un parámetro o hacer una asignación) copia el puntero, no la estructura apuntada: dos variables que "comparten" un struct realmente apuntan al mismo `StructVal` subyacente. Los tipos primitivos (`int`, `float64`, `string`, `bool`, `rune`) se copian por valor a través de sus campos escalares (`I`, `F`, `S`, `B`, `R`), consistente con el requisito de la sección 7.1.1 del enunciado ("los structs y slices se pasan por referencia; los demás tipos... se pasan por valor").

La función `append` es la excepción intencional a la mutación en el lugar: siempre produce un `SliceVal` **nuevo**, igual que en Go real, de modo que la variable debe reasignarse explícitamente con el resultado (`numeros = append(numeros, 4)`).

### 8.2 Entorno y alcance (`entorno.go`)

```go
type Entorno struct {
    Nombre    string
    variables map[string]*Valor
    mutable   map[string]bool
    Padre     *Entorno
}
```

Las variables se almacenan como `*Valor` (no `Valor`), de modo que asignar `numeros[i] = x` o `p.Campo = x` (resueltos mediante `resolverLugar`, sección 8.4) mutan la misma celda en memoria en vez de una copia. `Buscar` resuelve un identificador subiendo por la cadena de entornos padres (alcance léxico anidado); `DeclaradoLocal` detecta redeclaración dentro del mismo bloque. El mapa paralelo `mutable` guarda, por variable, si fue declarada con `mut` (ver sección 8.8, hallazgo A1): `EsMutable` lo consulta subiendo por la misma cadena de entornos que `Buscar`, y `Declarar` (usado para parámetros, receptores y variables de `for`) declara siempre como mutable, mientras que `DeclararMut` respeta la bandera capturada del fuente. El registro de `structs`, `funciones` y `métodos` vive aparte, en la estructura `Global`, independiente de la pila de entornos de ejecución (el enunciado exige que structs y funciones solo se declaren en el ámbito global, secciones 6.1 y 7.1).

### 8.3 El intérprete de dos pasadas (`interprete.go`)

```go
type Interprete struct {
    errores   *ListaErrores
    global    *Global
    entGlobal *Entorno
    consola   []string
    simbolos  map[string]*FilaSimbolo
}

func (in *Interprete) Interpretar(programa *ast.Programa)
```

Documentado en el propio código:

```go
// Interprete es el "ambiente fresco por ejecucion" (analogo a los otros
// proyectos del curso): dos pasadas (2.2 del enunciado) - registrar
// structs/funciones primero (permite llamadas hacia adelante, 7.1), luego
// declarar globales y correr main().
```

`Interpretar` ejecuta, en orden: (1) registro de todos los `structs` declarados, detectando redeclaración; (2) registro de todas las funciones y métodos (los métodos se indexan en `Metodos map[string]map[string]*ast.DeclFuncion`, por tipo de struct receptor y luego por nombre de método); (3) creación del entorno global y evaluación de las declaraciones de variables globales; (4) localización de la función `main` y su invocación — reportando error semántico si no existe o si declara parámetros (sección 3.2 del enunciado).

La separación en dos pasadas (registro global antes de ejecutar cualquier cuerpo) es lo que permite llamadas hacia adelante: una función puede invocar a otra declarada más abajo en el archivo fuente.

### 8.4 Resolución de "lugares" mutables

```go
func (in *Interprete) resolverLugar(nodo ast.Nodo, ent *Entorno) (*Valor, bool)
```

Dado un nodo `ast.Identificador`, `ast.ExprIndice` o `ast.ExprCampo`, devuelve un puntero (`*Valor`) a la celda real referenciada, en vez de una copia. Se usa tanto en asignaciones simples (`=`) como en `+=`/`-=` y en `++`/`--`, garantizando que `numeros[2] = 100` o `persona.Edad = 30` muten el dato real (en el slice o en el struct), y no una copia efímera.

### 8.5 Evaluación de expresiones y compatibilidad de tipos

```go
func (in *Interprete) evaluar(nodo ast.Nodo, ent *Entorno) (Valor, bool)
func tipoCompatible(declarado Tipo, v Valor) bool
func coercionar(declarado Tipo, v Valor) Valor
```

`evaluar` recorre el AST con un `switch` sobre el tipo dinámico del nodo (literales, identificador, unaria, binaria, índice, campo, llamada, literal de slice, literal de struct), devolviendo siempre un par `(Valor, bool)`: el booleano en `false` indica que ya se reportó un error y que la expresión contenedora debe abandonar su evaluación sin generar errores adicionales en cascada.

`tipoCompatible`/`coercionar` implementan la regla de compatibilidad de tipos, documentada en el propio código y ampliada en `docs/gramatica.txt` (sección de decisiones de diseño, punto 6): se exige tipo idéntico salvo para la combinación `int`/`float64`, que se admite en cualquier dirección y se resuelve mediante coerción explícita al tipo declarado.

### 8.6 Funciones nativas y conversión de tipo (`nativas.go`)

```go
func EsNombreNativa(nombre string) bool
func (in *Interprete) llamarNativa(nombre string, args []Valor, linea, col int) (Valor, bool)
func (in *Interprete) convertirTipo(destino string, args []Valor, linea, col int) (Valor, bool)
```

`EsNombreNativa` reconoce tanto las funciones embebidas del enunciado (`print`, `println`, `len`, `append`, `indexOf`, `join`, `Atoi`, `parseFloat`, `typeOf` — todas sin el prefijo de paquete que aparece únicamente en los títulos de sección del enunciado, ver `docs/gramatica.txt` punto 4) como los nombres de conversión de tipo estilo Go (`int`, `float64`, `string`, `bool`, `rune`), resolviendo la ausencia de una sección dedicada a "casteos" en el enunciado (`docs/gramatica.txt`, punto 5) mediante el único ejemplo real que sí aparece (`f64(10 + 1)`).

#### `print` y `println` comparten implementación — y por qué

Las dos comparten la misma rama del `switch` de `llamarNativa` (`case "print", "println":`) y hacen exactamente lo mismo: convierten cada argumento con `AImprimir`, los unen con un espacio y **agregan una entrada** a `in.consola`.

No es un descuido, es consecuencia del modelo de consola. `in.consola` es un `[]string` y viaja al cliente como `ConsolaLineas []string` (más un `Consola` de conveniencia con las líneas unidas por `\n`). En ese modelo **la unidad mínima es la línea**: no existe la noción de "escribir sin avanzar de línea", que es justamente lo que distinguiría un `print` de un `println` en una consola de flujo de caracteres. Implementar la diferencia obligaría a que la consola fuera un buffer de texto y a definir cuándo se corta una línea — un rediseño que el enunciado no pide.

El enunciado (§7.2.1) solo define **`fmt.Println`**; `print` se acepta como sinónimo por comodidad y no aparece en sus ejemplos. Si en una versión futura hiciera falta distinguirlas, el cambio se localiza en `llamarNativa` y en el tipo del campo `consola`, no en el resto del intérprete.

### 8.7 Validaciones semánticas reforzadas (auditoría 2026-07-21)

Una auditoría posterior a la primera versión funcional detectó y corrigió 3 huecos de validación semántica que cambian el comportamiento en tiempo de ejecución (un cuarto hallazgo de esa misma auditoría, la comparación relacional lexicográfica de strings, ya está documentado en `docs/gramatica.txt`, punto 8, e implementado en la función `Relacional` de `operaciones.go`):

1. **`break`/`continue` fuera de un ciclo o `switch` (4.8.1/4.8.2)**: antes de la corrección no se validaba el contexto — un `break` o `continue` sueltos en cualquier parte del programa producían igualmente su señal, sin verificar que estuvieran dentro de un `for`/`switch`. Se agregaron los contadores `profLoop`/`profSwitch` en `Interprete` (incrementados al entrar a `ejecutarFor`/`ejecutarSwitch` y decrementados con `defer` al salir) y la validación correspondiente en `ejecutarSentencia`:
   ```go
   case *ast.SentenciaBreak:
       if in.profLoop == 0 && in.profSwitch == 0 {
           in.errorSemantico("la sentencia \"break\" no puede usarse fuera de un ciclo o un switch", n.Linea, n.Columna)
           return Senal{Tipo: SenalNinguna}
       }
       return Senal{Tipo: SenalBreak}
   ```
   `invocarFuncion` resetea ambos contadores a 0 al entrar a una función o método (restaurando los valores previos al salir), de modo que un `break`/`continue` suelto en el cuerpo de una función invocada desde dentro de un `for` no herede el contexto de ciclo del llamador.

2. **Colisión de nombres función/variable/struct en el ámbito global (7.1)**: antes de la corrección solo se detectaba redeclaración dentro de la misma categoría (dos `structs` con igual nombre, dos funciones con igual nombre), pero no la colisión **entre** categorías que exige 7.1 ("las funciones, variables o structs no pueden tener el mismo nombre"). Ahora, al registrar una función en `Interpretar` se verifica que su nombre no coincida con un struct ya registrado, y al declarar una variable global (`ejecutarDeclaracion`, cuando `ent == in.entGlobal`) se verifica que su nombre no coincida con un struct ni con una función ya registrados.

3. **`switch` que silenciaba errores de tipo entre `case`s (4.4.1)**: comparar el valor del `switch` contra un `case` de tipo incompatible (por ejemplo, un `int` contra un `struct`) se descartaba en silencio en `ejecutarSwitch`, como si el caso simplemente no coincidiera. Ahora el error que devuelve `valoresIguales` se reporta como error semántico — el mismo tratamiento que ya recibía `==` fuera de un `switch` — en vez de ignorarse con un `continue`.

Ninguno de los 4 hallazgos de la auditoría (los 3 anteriores más la comparación lexicográfica de strings) es un fallo de ejecución: los cuatro son huecos de validación semántica que el enunciado exige cerrar, no errores que impidieran correr el programa.

### 8.8 Segunda auditoría (2026-07-23): semántica, robustez y reportes

Una segunda auditoría —parte de una revisión cruzada de los proyectos del curso— leyó el runtime completo y detectó una nueva tanda de hallazgos, todos latentes (ningún ejemplo los ejercitaba). Se corrigieron en bloque y se agregó `entradas/ejemplo7_semantica.vch`, que ejercita los caminos afectados y debe correr sin errores, más dos casos nuevos en `entradas/ejemplo6_errores.vch`. Igual que la auditoría anterior, ninguno era un fallo que impidiera correr el intérprete; eran huecos de validación o de robustez.

**Hallazgos de semántica del lenguaje:**

1. **`mut` dejó de ser decorativo (A1).** La gramática acepta `'mut'?` en ambas formas de declaración (`VLangCherry.g4:92-94`) y el `ManualUsuario` documenta que "`mut` indica que la variable puede reasignarse", pero el runtime nunca leía la bandera: reasignar una variable declarada sin `mut` se aceptaba en silencio. Se agregó el campo `Mutable` a `ast.DeclVariable`, capturado en el traductor desde `MUT()`; el `Entorno` rastrea la mutabilidad (`DeclararMut`/`EsMutable`, sección 8.2); y `verificarReasignable` valida `=`, `+=`, `-=`, `++` y `--` en `ejecutarAsignacion`/`ejecutarIncDec`. **Alcance decidido y documentado:** `mut` gobierna únicamente reasignar la **variable entera**; mutar un campo o un elemento (`p.Edad = 30`, `xs[0] = 9`) se permite siempre, porque no reenlaza la variable. Las variables de control del `for` se fuerzan mutables en el traductor (una `for i := 0; …; i++` sin `mut` seguiría siendo válida), lo mismo que parámetros y receptores.

2. **Cortocircuito de `&&` y `||` (A2).** `evaluar` evaluaba siempre ambos operandos antes de despachar a los operadores lógicos. Ahora, en el caso `*ast.ExprBinaria`, si el operando izquierdo es booleano y ya decide el resultado (`false` para `&&`, `true` para `||`), se retorna sin evaluar el derecho — la semántica habitual, y lo que permite que una guarda como `x != 0 && 10/x > 1` no dispare "división entre cero".

3. **Receptor por valor que compartía el struct (A4).** `ReceptorPuntero` se capturaba en el AST pero el runtime nunca lo leía: `invocarFuncion` declaraba el receptor con una copia superficial del `Valor`, y como `Valor.Struct` es un puntero, un método con receptor **por valor** terminaba mutando el `StructVal` del llamador. Se agregó `ClonarPorValor` (`valores.go`), que replica la copia por valor de Go —clona el struct y sus structs anidados en profundidad; los slices internos se comparten, igual que en Go real— y se aplica en `invocarFuncion` cuando `!fn.ReceptorPuntero`. Ahora un receptor por valor trabaja sobre una copia y no afecta al llamador; uno por puntero sí comparte el struct.

4. **`return` mal validado contra el tipo declarado (M1/M2).** El zero-value de `Tipo` es `int`, así que un `return` sin expresión en una función tipada devolvía `0` sin error, y un `return <valor>` en una función sin tipo de retorno se descartaba en silencio. Se agregó la bandera `TieneValor` a la señal `Senal`: ahora `invocarFuncion` reporta error si una función tipada trae un `return` sin valor, y si una función `void` trae un `return` con valor.

**Hallazgos de robustez del servidor:**

5. **Guardas anti-desborde (A3).** No existían límites de iteración ni de profundidad. Un `for true {}` dejaba la goroutine de la petición girando para siempre, y —más grave— una recursión sin caso base provocaba un **stack overflow de Go, que es un `fatal error` no recuperable por `recover()`** (ver sección 9): mataba el proceso entero del servidor, no solo esa petición. Se agregaron dos contadores, con los mismos valores que usa el proyecto hermano CompInterpreter: `maxIteraciones` (1 000 000) por ciclo en `ejecutarFor`, y `maxProfundidad` (2000) de llamadas anidadas en `invocarFuncion`. Al superarse, se reporta un error semántico ("posible ciclo infinito" / "recursión demasiado profunda") en vez de colgar o derribar el servidor.

6. **`append`/`join` sobre un slice `nil` (A5/M3).** `ValorPorDefecto` devolvía para un slice o struct sin inicializar un `Valor{Tipo: TipoNil()}` que **perdía el tipo declarado**, de modo que una asignación posterior (`xs = []int{…}` tras `mut xs []int`) se rechazaba como incompatible con `nil`. Ahora devuelve `Valor{Tipo: t}` conservando el tipo (con `Slice`/`Struct` en `nil`, forma que `EsNil` ya contemplaba). Como ese fix hace alcanzable un slice `nil` real, se endurecieron además `append` y `join` en `nativas.go` para tratar el slice `nil` como Go (append de un elemento; join de cadena vacía) en vez de hacer panic.

**Hallazgos de reportes y defensa (menores):**

7. **Tabla de símbolos por sitio de declaración (M4).** `registrarSimbolo` usaba como clave `ambito::nombre` en un mapa, de modo que dos bloques hermanos de la misma función que declararan el mismo nombre se pisaban mutuamente (comparten el `Nombre` del entorno de la función). La clave ahora incluye la posición en el fuente (`ambito::nombre::linea:col`), y `Simbolos()` ordena las filas por ámbito y luego por línea/columna. **Decisión de diseño documentada:** la tabla es un registro de **declaraciones** —una fila por sitio de declaración, con el valor **al momento de declarar** (snapshot)—, no una traza de ejecución: una variable declarada dentro de un ciclo o de una función recursiva ocupa siempre la misma posición y por lo tanto una sola fila, y una reasignación posterior no actualiza la columna Valor.

8. **Errores defensivos y campo duplicado (B1/B2).** Los casos `default` de `resolverLugar` y `evaluar` devolvían `false` sin mensaje ante un nodo inesperado (fallo silencioso); ahora reportan un error con posición ("destino de asignación no válido" / "expresión no soportada"). Y un literal de struct con el mismo campo repetido (`Persona{Nombre:"a", Nombre:"b"}`), que antes dejaba ganar al último en silencio, ahora es error semántico.

## 9. Manejo de errores

`internal/runtime/errores.go` define una estructura única para las tres familias de error que exige el enunciado (sección 8.1):

```go
type ErrorReporte struct {
    Tipo        string `json:"tipo"`
    Descripcion string `json:"descripcion"`
    Linea       int    `json:"linea"`
    Columna     int    `json:"columna"`
}

type ListaErrores struct {
    Errores []ErrorReporte
}

func (l *ListaErrores) Lexico(desc string, linea, columna int)
func (l *ListaErrores) Sintactico(desc string, linea, columna int)
func (l *ListaErrores) Semantico(desc string, linea, columna int)
```

Los errores léxicos y sintácticos se capturan reemplazando el listener de error por defecto de ANTLR:

```go
type oyenteErrores struct {
    *antlr.DefaultErrorListener
    errores *runtime.ListaErrores
    esLexer bool
}

func (o *oyenteErrores) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{},
    line, column int, msg string, e antlr.RecognitionException) {
    if o.esLexer {
        o.errores.Lexico(msg, line, column+1)
    } else {
        o.errores.Sintactico(msg, line, column+1)
    }
}
```

Los errores semánticos se generan directamente desde el intérprete en el punto exacto donde se detectan (variable no declarada, tipo incompatible, campo o struct inexistente, redeclaración en el mismo ámbito, división entre cero, condición no booleana, etc.), sin interrumpir la ejecución del resto del programa: cada función de evaluación devuelve `(Valor, bool)`, y el llamador que recibe `false` simplemente detiene la evaluación de la expresión que lo contiene — un error se reporta una única vez y no genera errores adicionales en cascada.

Adicionalmente, la orquestación completa (`internal/analizar/analizar.go`) envuelve la traducción y la interpretación en un bloque protegido con `recover()`:

```go
func() {
    defer func() {
        if r := recover(); r != nil {
            errores.Semantico(fmt.Sprintf("error interno durante la ejecución: %v", r), 0, 0)
        }
    }()
    programa := traductor.Traducir(arbolParse)
    grafo = ast.ConstruirGrafo(programa)
    interprete.Interpretar(programa)
}()
```

Esto evita que un **panic** interno no contemplado (por ejemplo, un tipo de nodo del AST sin manejar, o un índice fuera de rango dentro del propio intérprete) derribe el proceso del servidor; en su lugar, se reporta como un error semántico adicional. Con una salvedad importante: `recover()` solo atrapa *panics*, no los `fatal error` de Go. El caso típico es el **stack overflow** por recursión sin caso base, que Go termina de forma no recuperable y que mataría el proceso entero pese al `recover()`. Por eso la robustez del servidor **no** descansa solo en este bloque: las guardas `maxProfundidad`/`maxIteraciones` de la sección 8.8 (hallazgo A3) atajan esos dos casos —recursión infinita y ciclo infinito— antes de que lleguen a ser fatales, reportándolos como error semántico.

## 10. El servidor REST (`cmd/servidor`, `internal/analizar`)

El servidor (`cmd/servidor/main.go`) se implementa únicamente con la librería estándar de Go (`net/http`), sin frameworks. Expone:

| Método y ruta | Descripción |
|---|---|
| `GET /` | Información básica del servicio (nombre, estado, endpoint disponible). |
| `GET /salud` | Verificación de estado (`{"estado": "ok"}`). |
| `POST /interpretar` | Endpoint principal. Recibe `{ "codigo": "<fuente .vch>" }` y responde el resultado del análisis. |

Todas las rutas incluyen encabezados CORS (`Access-Control-Allow-Origin: *`), habilitando que el cliente web (servido desde un origen distinto en desarrollo) consuma la API sin restricciones.

`internal/analizar/analizar.go` implementa el pipeline completo por cada petición, con un `Entorno` (en el sentido de instancia de ejecución) **nuevo en cada llamada** — no se conserva estado entre peticiones:

```go
type Resultado struct {
    Errores       []runtime.ErrorReporte `json:"errores"`
    Consola       string                 `json:"consola"`
    ConsolaLineas []string               `json:"consolaLineas"`
    Simbolos      []runtime.FilaSimbolo  `json:"simbolos"`
    AST           ast.Grafo              `json:"ast"`
    Dot           string                 `json:"dot"`
}

func Analizar(codigo string) Resultado
```

## 11. El cliente web (`client/`)

Aplicación React (v19) construida con Vite, estructurada en componentes con responsabilidad única:

- **`App.jsx`**: mantiene el estado de los archivos abiertos (arreglo de objetos `{id, nombre, contenido, sinGuardar}`), el archivo activo, el resultado de la última ejecución y el panel de resultado seleccionado. Orquesta las acciones Nuevo/Abrir/Guardar/Ejecutar.
- **`api.js`**: función `interpretar(codigo)`, realiza el `fetch` a `POST {VITE_API_URL}/interpretar` y devuelve el JSON de respuesta.
- **`Toolbar.jsx`**, **`FileTabs.jsx`**: barra de acciones y pestañas de archivo multi-documento.
- **`Editor.jsx`**: área de edición de texto con numeración de línea y resaltado de las líneas reportadas con error (recibe el conjunto de líneas con error calculado en `App.jsx` a partir de `resultado.errores`).
- **`Consola.jsx`**, **`TablaErrores.jsx`**, **`TablaSimbolos.jsx`**: presentación tabular de los reportes correspondientes; las dos tablas soportan clic en fila para saltar a la línea en el editor.
- **`AstGrafo.jsx`**: renderiza `resultado.ast` (`{nodes, edges}`) con la librería **vis-network**, en un layout jerárquico (`direction: 'UD'`), reutilizando variables CSS del tema activo (claro/oscuro) para los colores de nodo, borde y texto.

El cliente reutiliza y adapta el patrón de interfaz de otro proyecto del curso con arquitectura cliente-servidor equivalente (mismo contrato `{errores, consola, simbolos, ast}`), reemplazando únicamente la capa de comunicación para apuntar al servidor Go en vez de a un servidor Node/Express.

Las pruebas end-to-end (`client/e2e/vlangcherry.spec.js`, Playwright) verifican: la ejecución del ejemplo inicial y la visualización de consola/símbolos/AST; el reporte de errores semánticos y el salto a línea; y la creación de un archivo nuevo con extensión `.vch`.

## 12. Decisiones de diseño no evidentes

1. **AST propio en vez de operar sobre el parse tree de ANTLR**: se traduce explícitamente el parse tree a un AST propio (`internal/ast`) para desacoplar la lógica de interpretación de los tipos generados por ANTLR (que cambiarían si se regenera la gramática) y para tener control total sobre el formato del reporte de AST (sección 8.3 del enunciado).
2. **Traductor por type-switch en vez de Visitor generado**: ANTLR genera una interfaz `VLangCherryVisitor` con un método por regla/alternativa etiquetada. En vez de implementarla, el traductor aprovecha que cada alternativa etiquetada ya tiene su propio tipo de contexto Go concreto, y usa un `switch ctx.(type)` directo — funcionalmente equivalente, pero con menos código repetido.
3. **Referencia "gratis" para structs y slices vía punteros de Go**: en vez de implementar una tabla de referencias o un mecanismo de "boxing" explícito para cumplir el requisito de paso por referencia de la sección 7.1, se aprovecha que `Valor.Slice` y `Valor.Struct` ya son punteros; copiar un `Valor` copia el puntero, logrando la semántica de referencia compartida sin código adicional.
4. **Compatibilidad de tipos relajada entre `int` y `float64`**: la sección 3.6 del enunciado exige tipo exacto en toda asignación, pero las tablas aritméticas de la sección 4.3 muestran promoción automática `int + float64 -> float64` asignada de vuelta a variables `float64`. Se decidió permitir la mezcla `int`/`float64` en cualquier asignación (no solo para el literal `0`, que es la única excepción que menciona explícitamente el enunciado), documentado en el código (`tipoCompatible`/`coercionar`) y en `docs/gramatica.txt`.
5. **Funciones nativas y conversión de tipo sin prefijo de paquete**: los títulos de la sección 7.2 del enunciado usan nombres con prefijo de paquete al estilo Go (`strconv.Atoi`, `reflect.TypeOf`), pero todos los ejemplos de código ejecutable de esa misma sección invocan la función sin prefijo. Se implementó sin prefijo, consistente con los ejemplos reales.
6. **Switch sin fall-through**: la sección 4.7.2 del enunciado indica explícitamente que "el break implícito está incluido al final de cada case". Es una decisión de diseño real del lenguaje, distinta de otros lenguajes del curso con semántica de switch con fall-through — no una limitación de la implementación.
7. **Recuperación ante pánico interno**: la orquestación del análisis (`internal/analizar/analizar.go`) envuelve la traducción e interpretación en un `defer recover()`, de modo que un caso interno no contemplado se reporte como error semántico adicional en vez de derribar el proceso del servidor completo — relevante en un servicio de larga duración que atiende múltiples peticiones.
8. **Entorno nuevo por ejecución**: cada llamada a `Analizar(codigo string)` crea una instancia nueva de `runtime.Interprete` y de sus estructuras internas (consola, tabla de símbolos, lista de errores), de modo que los reportes de una ejecución nunca se mezclan con los de otra.
9. **Auditoría de validación semántica (2026-07-21)**: una revisión posterior a la primera versión funcional encontró 4 huecos donde el intérprete dejaba pasar (o descartaba en silencio) casos que el enunciado exige reportar como error: `break`/`continue` fuera de un ciclo o `switch` sin validar, faltaba la comparación relacional lexicográfica entre strings, faltaba el chequeo de colisión de nombres entre función/variable/struct en el ámbito global, y el `switch` silenciaba errores de tipo entre `case`s en vez de reportarlos. Los 3 primeros que cambian comportamiento en tiempo de ejecución están detallados en la sección 8.7; la comparación de strings está documentada en `docs/gramatica.txt` (punto 8) y en `operaciones.go` (función `Relacional`).

## 13. Compilación y ejecución

### 13.1 Regenerar el analizador léxico/sintáctico

Cuando se modifica `grammar/VLangCherry.g4`:

```
java -jar tools/antlr.jar -Dlanguage=Go -visitor -o server/internal/parser -package parser grammar/VLangCherry.g4
```

Esto regenera los cuatro archivos de `server/internal/parser/`, que no deben editarse manualmente en ningún otro momento.

### 13.2 Compilar y ejecutar el servidor (desarrollo)

```
cd server
go build ./...
go run ./cmd/servidor
```

### 13.3 Compilación cruzada para la entrega (ejecución nativa en Linux)

El enunciado exige que la aplicación se ejecute de forma nativa en un sistema operativo Linux (restricción 10.2). El desarrollo se realizó en un entorno Windows, pero Go permite compilar de forma cruzada sin herramientas adicionales, generando un binario estático apto para Linux:

```
cd server
GOOS=linux GOARCH=amd64 go build -o vlangcherry-servidor ./cmd/servidor
```

El archivo `vlangcherry-servidor` resultante se transfiere al equipo Linux de destino y se ejecuta directamente (`./vlangcherry-servidor`), sin necesidad de instalar Go ni ninguna dependencia adicional en dicho equipo.

### 13.4 Probar el intérprete por línea de comandos

```
cd server
go run ./cmd/cli entradas/ejemplo1_basico.vch
```

Imprime en la salida estándar las secciones `=== ERRORES ===`, `=== CONSOLA ===`, `=== SIMBOLOS ===` y `=== AST ===` (con el conteo de nodos y aristas), útil para verificar rápidamente el resultado de un archivo sin levantar el servidor ni el cliente.

### 13.5 Cliente web

```
cd client
npm install
npm run dev      # modo desarrollo
npm run build    # construcción de producción
npm run e2e      # pruebas Playwright
```

## 14. Mantenimiento futuro

- Cualquier cambio en los tokens o en la sintaxis del lenguaje debe hacerse en `grammar/VLangCherry.g4`, seguido de la regeneración descrita en la sección 13.1. Los archivos de `internal/parser/` se sobrescriben por completo en cada regeneración.
- Si se agrega una nueva alternativa a la regla `expr` (o a cualquier otra regla con alternativas etiquetadas), debe declararse con su propia etiqueta `# nombre` en el `.g4`; de lo contrario, ANTLR no genera un tipo de contexto propio y el `traductor` no podrá distinguirla en su type-switch.
- Si se agrega un nuevo tipo de nodo al AST (`internal/ast/ast.go`), no es necesario modificar `grafo.go`: el recorrido por reflection lo incorpora automáticamente al reporte de AST, siempre que sus campos sean valores exportados (con mayúscula inicial) de tipos reconocidos (`Nodo`, slice de `Nodo`, struct contenedor, `TipoAST`, o escalar).
- Si se agrega una nueva función nativa, debe registrarse tanto en `EsNombreNativa` como en el `switch` de `llamarNativa` (`internal/runtime/nativas.go`), documentando su comportamiento y las validaciones de tipo de sus argumentos, siguiendo el patrón de las funciones existentes.
- Si se agrega un nuevo campo al reporte de errores o de símbolos, debe reflejarse en las estructuras `ErrorReporte`/`FilaSimbolo` (con su etiqueta JSON) y en los componentes correspondientes del cliente (`TablaErrores.jsx`/`TablaSimbolos.jsx`), que consumen esas mismas claves.
