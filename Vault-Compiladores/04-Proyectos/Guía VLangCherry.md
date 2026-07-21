---
tags: [proyecto, guia, lexico, sintactico, semantico, olc2]
aliases: [guia vlangcherry, tutorial vlangcherry, vlangcherry paso a paso]
fuente: "Enunciado OLC2_P1_EV2025 + construcción real (código verificado: 2026-07-21)"
fecha: 2026-07-21
---

# Guía de elaboración — [[VLangCherry]]

Guía de referencia sobre el código REAL de VLangCherry, en `Hades\VLangCherry\`. A diferencia de las guías de los proyectos de OLC1 (donde Claude construyó paso a paso con el usuario tutorando), este proyecto es de **OLC2** — curso distinto — y ya llegó **completo y verificado** en una sola sesión: gramática [[ANTLR]] + intérprete en Go + servidor REST + cliente React, con 6 ejemplos en `entradas/` y una suite Playwright de 3 pruebas end-to-end, todo corriendo. Esta guía documenta la arquitectura real para estudio y para responder preguntas de autoría (sección 11.1 del enunciado).

⚠️ Es un proyecto **grupal de 3 personas** (10.4/12.3) y debe **ejecutarse nativamente en Linux** (10.2) — ver sección 8 (Build cruzado) más abajo.

## Estado de construcción

| Área | Contenido | Estado |
|---|---|---|
| Léxico + sintáctico | Gramática [[ANTLR]] (`grammar/VLangCherry.g4`), lexer/parser generados para Go | ✅ 2026-07-21 |
| AST propio | `internal/ast/ast.go` (~25 tipos de nodo) + `internal/ast/grafo.go` (grafo genérico por reflection) | ✅ 2026-07-21 |
| Traductor | `internal/traductor/traductor.go`: parse tree (ANTLR) → AST propio, sin Visitor generado | ✅ 2026-07-21 |
| Modelo de valores | `internal/runtime/valores.go`: `Valor` con slices/structs por puntero | ✅ 2026-07-21 |
| Intérprete | `internal/runtime/interprete.go`: dos pasadas (registro global → globales → `main()`) | ✅ 2026-07-21 |
| Errores | `internal/runtime/errores.go`: recolección de léxico+sintáctico+semántico sin abortar | ✅ 2026-07-21 |
| Servidor REST | `cmd/servidor/main.go`: `net/http`, sin frameworks, puerto 4100 | ✅ 2026-07-21 |
| Cliente web | `client/` (React + Vite), adaptado de [[CompInterpreter]] | ✅ 2026-07-21 |
| Pruebas | 6 `.vch` en `entradas/` (CLI) + 3 tests Playwright (`e2e/vlangcherry.spec.js`) | ✅ 2026-07-21, todos verificados en verde |
| Manuales/gramática entregable | `docs/gramatica.txt`, `docs/ManualUsuario.md`, `docs/ManualTecnico.md` | ✅ 2026-07-21 |

Pendiente (fuera del alcance de esta sesión): empaquetado final del repo de entrega `OLC2_Proyecto1_#Carnet` con colaborador `vallit0` agregado.

## 1. La gramática ([[ANTLR]])

`grammar/VLangCherry.g4` es una gramática ANTLR4 de un solo archivo (léxico + sintáctico juntos, a diferencia de [[JFlex]]+[[CUP]] que son dos herramientas separadas). Se genera para Go:

```
java -jar tools/antlr.jar -Dlanguage=Go -visitor -o server/internal/parser -package parser grammar/VLangCherry.g4
```

Produce en `server/internal/parser/`: `vlangcherry_lexer.go`, `vlangcherry_parser.go`, `vlangcherry_base_visitor.go`, `vlangcherry_visitor.go` — **100% generado**, no se edita nunca (misma regla que `Lexer`/`Parser`/`sym` en [[DataForge]]).

El orden de las alternativas en la regla recursiva `expr` codifica la precedencia (tabla 4.6 del enunciado, de mayor a menor):

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

Cada alternativa etiquetada con `# nombre` genera un tipo de contexto Go concreto (`ExprMultiplicativaContext`, `ExprLlamadaContext`, ...) — esto es lo que hace posible el traductor por type-switch de la sección 3 (en vez de implementar el Visitor completo). Ver [[ANTLR]] para más detalle de la herramienta en sí.

La gramática BNF limpia (entregable, sección 9.2 del enunciado) con las decisiones de diseño reales está en `docs/gramatica.txt` — **no** es una copia mecánica del `.g4`.

## 2. El AST propio (`internal/ast/`)

### 2.1 Los nodos (`ast.go`)

Cada nodo implementa una interfaz mínima:

```go
type Nodo interface {
    Tipo() string
    Pos() (int, int)
}
```

~25 structs concretos (`Programa`, `DeclStruct`, `DeclFuncion`, `SentenciaIf`, `ExprBinaria`, `LiteralStruct`, ...), todos embebiendo `base{Linea, Columna int}` para el reporte de errores con posición. `ExprBinaria.Tipo()` es un caso interesante: no devuelve un string fijo, sino `"BINARIA_" + n.Operador` — el operador queda embebido en la etiqueta del nodo del reporte AST sin necesitar un campo aparte.

### 2.2 El grafo por reflection (`grafo.go`)

`ConstruirGrafo(raiz Nodo) Grafo` recorre CUALQUIER nodo usando `reflect`, sin necesitar mantenimiento cuando se agrega un tipo de nodo nuevo — **mismo patrón que `CompInterpreter/server/src/reportes/ast-grafo.js`, pero en Go**: donde el JS original inspecciona objetos dinámicamente, acá se inspeccionan structs vía `reflect.Value`/`reflect.Type`. La función distingue:

- **Punteros/interfaces** que implementan `Nodo`: se convierten en una arista hacia un nuevo nodo del grafo.
- **Slices**: cada elemento se procesa individualmente con un nombre de campo indexado (`Campos[0]`, `Campos[1]`, ...).
- **Structs "contenedor transparente"** (`RamaIf`, `CasoSwitch`, `Parametro`, `CampoStruct`, `CampoValorLiteral`): no son `Nodo` en sí mismos, así que sus campos se recorren propagando el nombre compuesto (`Ramas[0].Condicion`).
- **`TipoAST`**: caso especial, se imprime como texto (`t.String()`) pegado a la etiqueta del nodo padre en vez de generar un nodo aparte — un tipo declarado no necesita su propio nodo en el AST visual.
- **Escalares**: se agregan como texto a la etiqueta del nodo (`"nombreCampo: valor"`).

El resultado (`Grafo{Nodes, Edges}`) viaja como JSON al cliente, que lo dibuja con [[vis-network]] (`AstGrafo.jsx`) — el mismo formato `{nodes, edges}` que ya usaba [[CompInterpreter]]. También se genera un `.dot` equivalente (`ADot`) por si se quiere graficar del lado del servidor con [[Graphviz]] (sugerencia del enunciado, sección 8.3/reportes).

## 3. El traductor (`internal/traductor/traductor.go`)

Convierte el parse tree que genera ANTLR (`parser.IProgramaContext`) al AST propio (`ast.Programa`). Decisión de diseño explícita en el comentario de cabecera del archivo:

```go
// Package traductor convierte el parse tree que genera ANTLR en el AST
// propio de internal/ast. No usa el patron Visitor generado: como cada
// alternativa etiquetada de la gramatica produce un tipo de contexto Go
// concreto, alcanza con un type-switch directo (mas corto que implementar
// las ~60 firmas de BaseVLangCherryVisitor).
```

Por ejemplo, `traducirExpr` hace un `switch` sobre el tipo concreto de contexto en vez de sobreescribir cada método `Visit...`:

```go
func traducirExpr(ctx parser.IExprContext) ast.Nodo {
    switch ec := ctx.(type) {
    case *parser.ExprLlamadaContext:
        ...
    case *parser.ExprMultiplicativaContext:
        return traducirBinaria(ec, ec.AllExpr(), ec.GetOp().GetText())
    ...
    }
}
```

Esto es una alternativa legítima al patrón *Visitor* clásico del Dragón (cap. 5): en vez de doble despacho (`accept`/`visit`), se aprovecha que Go ya sabe distinguir el tipo dinámico de la interfaz con un type-switch — funciona porque ANTLR generó un tipo por alternativa etiquetada, no un único nodo genérico.

Otro detalle no evidente: `traducirLugar` y `traducirExpr` producen **los mismos nodos** (`ast.Identificador`, `ast.ExprIndice`, `ast.ExprCampo`) tanto para lectura como para el lado izquierdo de una asignación — el intérprete después decide si necesita el valor (`evaluar`) o el puntero a la celda (`resolverLugar`) según el contexto.

## 4. El modelo de valores (`internal/runtime/valores.go`)

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

La decisión central, documentada en el propio archivo:

```go
// Valor es el valor en tiempo de ejecucion de una expresion VLangCherry.
// Slice y Struct viajan por puntero: asi una copia del Valor (paso de
// parametro, asignacion) comparte la misma identidad subyacente, que es
// exactamente el "por referencia" que pide la seccion 7 para structs y
// slices (los primitivos, en cambio, se copian por valor via los campos
// escalares I/F/S/B/R).
```

El enunciado (sección 7.1, considerandos) pide que structs y slices se pasen **por referencia** y los primitivos **por valor**. En vez de implementar dos mecanismos de paso de parámetros distintos, se aprovecha que `Slice`/`Struct` ya son punteros (`*SliceVal`/`*StructVal`) dentro del `Valor`: copiar un `Valor` (que Go hace por valor, como cualquier struct) copia el puntero, no la estructura apuntada — dos variables que "comparten" un struct realmente apuntan al mismo `StructVal`. Es la misma idea de **entornos y alcance con referencia compartida** de [[Entornos y alcance]], resuelta con punteros nativos de Go en vez de una tabla de referencias explícita.

`append` (sección 5.1.5) es la excepción intencional: siempre genera un `SliceVal` **nuevo** — igual que en Go real, la variable debe reasignarse con el resultado (`numeros = append(numeros, 4)`), reflejado en la firma de `llamarNativa("append", ...)`.

## 5. El intérprete de dos pasadas (`internal/runtime/interprete.go`)

Documentado en el comentario de la struct `Interprete`:

```go
// Interprete es el "ambiente fresco por ejecucion" (analogo a los otros
// proyectos del curso): dos pasadas (2.2 del enunciado) - registrar
// structs/funciones primero (permite llamadas hacia adelante, 7.1), luego
// declarar globales y correr main().
```

`Interpretar(programa *ast.Programa)` hace, en orden:

1. Registra TODOS los `structs` en `global.Structs` (detectando redeclaración).
2. Registra TODAS las funciones/métodos en `global.Funciones`/`global.Metodos` (los métodos se indexan por tipo de receptor, `map[string]map[string]*ast.DeclFuncion`).
3. Crea el entorno global (`NuevoEntorno("global", nil)`) y ejecuta las declaraciones de variables globales.
4. Busca `main` y la invoca — si no existe, error semántico; si tiene parámetros, error semántico (sección 3.2: "VLangCherry contará con una función main... el intérprete deberá localizar esta función").

La primera pasada es lo que permite **llamadas hacia adelante**: una función puede llamar a otra declarada más abajo en el archivo, porque para cuando se ejecuta cualquier cuerpo, `global.Funciones` ya tiene registradas todas.

### 5.1 Propagación de errores: `(Valor, bool)` en vez de excepciones

Cada función de evaluación devuelve `(Valor, bool)` — el booleano indica si hubo error. Es la misma filosofía de **propagación por null** que usa [[DataForge]] (Dragón, manejo de errores): un error se reporta una sola vez, en el punto exacto donde ocurre, y la expresión que lo contiene simplemente deja de evaluarse (`if !ok { return Valor{}, false }`) sin generar errores en cascada ni abortar la ejecución completa — el enunciado (8.1) exige recolectar TODOS los errores de una ejecución.

### 5.2 Compatibilidad de tipos relajada (sección 1 de `docs/gramatica.txt`)

```go
// tipoCompatible/coercionar: 3.6 exige tipo exacto, pero se relaja la
// mezcla int<->float64 (documentado en docs/gramatica.txt): el enunciado
// muestra ambos casos (el ejemplo de conversion explicita f64(10+1) Y
// tablas aritmeticas con promocion automatica int+float64->float64
// asignada de vuelta a variables float64), asi que se prioriza usabilidad.
func tipoCompatible(declarado Tipo, v Valor) bool { ... }
```

Ver [[Comprobación de tipos]] y [[Conversión de tipos (coerción y cast)]] para la teoría general; el detalle de por qué se tomó esta decisión frente al enunciado está en `docs/gramatica.txt`.

### 5.3 `resolverLugar`: puntero mutable a la celda

```go
func (in *Interprete) resolverLugar(nodo ast.Nodo, ent *Entorno) (*Valor, bool)
```

Da acceso a la celda real (`*Valor`) referenciada por un "lugar" (`ID` | `lugar[expr]` | `lugar.campo`), reutilizando los mismos nodos AST que las expresiones de lectura (`ExprIndice`, `ExprCampo`). Es lo que permite que `numeros[2] = 100` o `p.Nombre = "Bob"` muten el valor real en la [[Tabla de símbolos|tabla de símbolos]] / dentro del slice o struct, en vez de mutar una copia.

## 6. El servidor REST (`cmd/servidor/main.go`)

Contrato idéntico en forma al de [[CompInterpreter]], pero implementado con la librería estándar de Go (`net/http`), sin ningún framework:

```
POST /interpretar { "codigo": "<fuente .vch>" }
  -> { errores, consola, consolaLineas, simbolos, ast, dot }
```

`internal/analizar/analizar.go` orquesta el pipeline completo por cada petición — **entorno fresco por llamada** (mismo criterio que [[DataForge]] y [[ConjAnalyzer]]: los reportes son siempre del último análisis):

```go
entrada := antlr.NewInputStream(codigo)
lexer := parser.NewVLangCherryLexer(entrada)
lexer.RemoveErrorListeners()
lexer.AddErrorListener(&oyenteErrores{errores: errores, esLexer: true})
...
arbolParse := p.Programa()
programa := traductor.Traducir(arbolParse)
grafo = ast.ConstruirGrafo(programa)
interprete.Interpretar(programa)
```

Un detalle defensivo: la traducción + interpretación corre dentro de un `defer recover()` — si algo interno panickea (por ejemplo, un caso no contemplado de tipo), se convierte en un error semántico reportado en vez de tumbar el servidor.

## 7. El cliente (`client/`)

React + Vite, reusando/adaptando componentes del patrón de [[CompInterpreter]]: `Toolbar` (Nuevo/Abrir/Guardar/▶ Ejecutar), `FileTabs` (multi-archivo, pestañas con indicador de cambios sin guardar `•`), `Editor` (con gutter que resalta líneas con error), `Consola`, `TablaErrores`, `TablaSimbolos`, `AstGrafo` (vis-network, layout jerárquico). `api.js` hace el `fetch` a `POST /interpretar`.

## 8. Build cruzado a Linux (restricción 10.2 del enunciado)

El desarrollo se hizo en Windows, pero el enunciado exige ejecución nativa en Linux. Go cross-compila sin dependencias adicionales:

```
cd server
GOOS=linux GOARCH=amd64 go build -o vlangcherry-servidor ./cmd/servidor
```

El binario resultante corre en Linux sin necesidad de tener Go instalado ahí (Go produce binarios estáticos por defecto). Detalle completo en `docs/ManualTecnico.md`.

## 9. Casos de prueba

- `entradas/ejemplo1_basico.vch`: aritmética, comparación, lógicos, rune, `main` con variable global. Verificado: 0 errores, 121 nodos / 120 aristas en el AST.
- `entradas/ejemplo2_structs.vch`: structs anidados (`Persona` con campo `Direccion`), método por valor (`Saludar`) y por puntero (`Cumplir`, muta `Edad`).
- `entradas/ejemplo3_slices.vch`: slices, multidimensionales, `append`/`indexOf`/`join`/`len`.
- `entradas/ejemplo4_control.vch`: `if/else if/else`, `switch` sin fall-through, las 3 formas de `for`.
- `entradas/ejemplo5_funciones.vch`: recursión (`factorial`, `fibonacci`), funciones sin efectos colaterales (`esPar`). Verificado: 0 errores, 70 nodos / 69 aristas.
- `entradas/ejemplo6_errores.vch`: 7 errores a propósito (1 léxico, 6 semánticos: tipo incompatible, variable no declarada, campo inexistente, redeclaración, división entre cero, condición de `if` no booleana) — verificado que el intérprete reporta los 7 con línea/columna y la ejecución sobrevive (símbolos declarados antes del error siguen apareciendo en la tabla).
- `e2e/vlangcherry.spec.js` (Playwright, 3 tests): ejecutar el ejemplo inicial y ver consola/símbolos/AST; reportar errores semánticos y saltar a la línea marcada; crear un archivo nuevo con pestaña `.vch`.

## 10. Errores comunes / puntos de atención

- El `.g4` es la única fuente de verdad del léxico/sintaxis — el código en `internal/parser/` se regenera con ANTLR, nunca se edita a mano (idéntico criterio a [[JFlex]]/[[CUP]] en los proyectos hermanos).
- `traductor.go` depende de que cada alternativa de `expr` tenga su etiqueta `# nombre` en el `.g4`: si se agrega una alternativa sin etiquetar, ANTLR no genera un tipo de contexto propio y el type-switch de `traducirExpr` no la reconoce.
- La mezcla int↔float64 en asignaciones es una decisión DELIBERADA de relajar 3.6 (ver `docs/gramatica.txt`), no un bug — está documentada tanto ahí como en el comentario de `tipoCompatible`.
- El switch de VLangCherry NO tiene fall-through (a propósito, sección 4.7.2) — no confundir con el switch de [[CompInterpreter]], que sí lo tiene: son lenguajes distintos con semántica de switch distinta.
- Recordar el requisito de ejecución en Linux (10.2) al empaquetar la entrega: compilar con `GOOS=linux GOARCH=amd64`.

## Relacionadas
- [[VLangCherry]]
- [[ANTLR]] · [[vis-network]] · [[Graphviz]]
- [[Árbol de sintaxis abstracta (AST)]] · [[Comprobación de tipos]] · [[Conversión de tipos (coerción y cast)]]
- [[Tabla de símbolos]] · [[Entornos y alcance]] · [[Manejo de errores (léxicos, sintácticos, semánticos)]]
- [[CompInterpreter]] (mismo contrato REST, cliente hermano)
- [[DataForge]] · [[ConjAnalyzer]] (mismo criterio de "entorno fresco por ejecución")
