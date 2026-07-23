# Guía de Programación — VLangCherry

Universidad de San Carlos de Guatemala — Organización de Lenguajes y Compiladores 2
Material de clase para programar VLangCherry desde cero, tema por tema, con el código real del proyecto.

## A quién está dirigida esta guía

Esta guía asume que vos (o tu equipo) van a **escribir o modificar código real** de VLangCherry, no solo leerlo. Está en segunda persona porque cada sección termina pidiéndote que hagas algo: agregues una regla a la gramática, un tipo de nodo al AST, una función nativa, una validación semántica. Es, a la vez, la preparación más directa para la defensa oral: el enunciado del proyecto (11.1) exige que el equipo demuestre autoría real tocando código en vivo, y cada tema de esta guía es justo el tipo de cosa que les van a pedir que modifiquen frente al tribunal.

No reemplaza `docs/ManualTecnico.md` (que describe la arquitectura ya construida) ni `docs/gramatica.txt` (la BNF entregable con las decisiones de diseño frente al enunciado). Esta guía es el "cómo se llega ahí, paso a paso, escribiendo el código".

## Índice

1. [Preparación del proyecto](#1-preparación-del-proyecto)
2. [Tema: la gramática ANTLR4](#2-tema-la-gramática-antlr4)
3. [Tema: del parse tree al AST propio](#3-tema-del-parse-tree-al-ast-propio)
4. [Tema: tipos estáticos, structs y métodos "por referencia" gratis](#4-tema-tipos-estáticos-structs-y-métodos-por-referencia-gratis)
5. [Tema: slices 1D/2D](#5-tema-slices-1d2d)
6. [Tema: control de flujo y recursión](#6-tema-control-de-flujo-y-recursión)
7. [Tema: el servidor REST](#7-tema-el-servidor-rest)
8. [Tema: el cliente React](#8-tema-el-cliente-react)
9. [Errores comunes reales al programar cada tema](#9-errores-comunes-reales-al-programar-cada-tema)
10. [Para poder defender esto en vivo](#10-para-poder-defender-esto-en-vivo)

---

## 1. Preparación del proyecto

### 1.1 Qué necesitás instalado

- **Go 1.25** o compatible (el `go.mod` de `server/` fija `go 1.25.3`). Verificá con `go version`.
- **Java** (cualquier JRE moderna) — únicamente para correr `tools/antlr.jar`, el generador. No hace falta para compilar ni correr el proyecto Go en sí.
- **Node.js** — para el cliente React (`client/`).

No necesitás instalar ANTLR4 "de verdad": el `.jar` ya está commiteado en `tools/antlr.jar`, así que alcanza con `java -jar`.

### 1.2 Estructura que vas a tocar

```
VLangCherry/
├── grammar/VLangCherry.g4        ← vas a EDITAR este archivo (Tema 2)
├── tools/antlr.jar               ← generador, no se edita
├── docs/gramatica.txt            ← BNF entregable + decisiones de diseño
├── entradas/*.vch                ← casos de prueba, usalos para probar tu código
└── server/
    ├── go.mod                    ← module vlangcherry, go 1.25.3
    ├── cmd/
    │   ├── cli/                  ← ejecuta un .vch por consola (tu mejor herramienta de prueba rápida)
    │   ├── servidor/              ← servidor REST real
    │   └── asttest/, parsetest/  ← utilidades puntuales de verificación
    └── internal/
        ├── parser/                ← 100% GENERADO, nunca se edita a mano (Tema 2)
        ├── ast/                   ← tipos de nodo (Tema 3)
        ├── traductor/              ← parse tree → AST (Tema 3)
        └── runtime/                ← el intérprete (Temas 4-6)
```

### 1.3 El ciclo de trabajo que vas a repetir todo el curso

1. Editás `grammar/VLangCherry.g4` (si el cambio toca sintaxis) **o** código Go en `internal/`.
2. Si tocaste el `.g4`, regenerás el parser:
   ```
   java -jar tools/antlr.jar -Dlanguage=Go -visitor -o server/internal/parser -package parser grammar/VLangCherry.g4
   ```
3. Compilás y probás rápido con un archivo de `entradas/`, sin levantar servidor ni cliente:
   ```
   cd server
   go run ./cmd/cli entradas/ejemplo1_basico.vch
   ```
   Esto imprime `=== ERRORES ===`, `=== CONSOLA ===`, `=== SIMBOLOS ===` y `=== AST ===` directo en la terminal — es tu ciclo de retroalimentación más rápido, úsalo antes de tocar el servidor o el cliente.
4. Cuando el cambio funciona por CLI, recién ahí probás con el servidor + cliente si hace falta ver el reporte visual.

**Ejercicio para practicar el ciclo antes de seguir:** corré los 6 ejemplos de `entradas/` uno por uno con `go run ./cmd/cli entradas/ejemploN_*.vch` y mirá la salida completa de cada uno. Vas a reconocer en la salida real las mismas estructuras (`ErrorReporte`, `FilaSimbolo`, `Grafo`) que vamos a explicar en esta guía.

---

## 2. Tema: la gramática ANTLR4

### 2.1 Por qué un solo archivo `.g4`

En los proyectos de OLC1 del curso (DataForge, ConjAnalyzer, CompScript) el léxico y la sintaxis son **dos herramientas separadas**: JFlex genera el lexer, CUP genera el parser a partir de los tokens de JFlex. ANTLR4 hace las dos cosas en un solo archivo, `grammar/VLangCherry.g4`: las reglas en minúscula (`programa`, `expr`, `sentencia`...) son reglas sintácticas; las reglas en MAYÚSCULA (`ID`, `ENTERO`, `CADENA`...) son reglas léxicas. No es una preferencia de estilo: es la única opción de los tres generadores del curso que emite código **Go**, y el enunciado de OLC2 (sección 2.3) exige que toda la lógica del compilador esté en Go.

Vos vas a trabajar casi siempre en la mitad sintáctica del archivo; la mitad léxica ya cubre todos los tokens del lenguaje.

### 2.2 La regla `expr`: donde vive la precedencia

Abrí `grammar/VLangCherry.g4` y mirá la regla `expr`:

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

Vos ya sabés, del Dragón (cap. 4), que la precedencia de operadores en una gramática LL/LR clásica se resuelve **factorizando** en una cascada de reglas (`expr → term → factor → unario → primario`, cada nivel de precedencia una regla). ANTLR4 no te obliga a eso: la regla `expr` es recursiva izquierda de verdad, y **el orden de las alternativas ES la precedencia** — la que está más arriba se resuelve primero (liga más fuerte). ANTLR4 resuelve la recursión izquierda por vos en tiempo de generación. Esto es lo que hace la regla tan compacta comparada con una cascada manual.

La tabla de precedencia real (`docs/gramatica.txt`, sección 6, de mayor a menor):

```
( ) [ ] . llamada  >  ! - (unario)  >  * / %  >  + -  >  < <= >= >  >  == !=  >  &&  >  ||
```

### 2.3 Las etiquetas `# nombre`: la pieza que habilita todo lo demás

Cada alternativa etiquetada (`# exprMultiplicativa`, `# exprLlamada`, ...) hace que ANTLR4 genere, al compilar la gramática, un **tipo de contexto Go concreto y distinto** para esa alternativa: `ExprMultiplicativaContext`, `ExprLlamadaContext`, etc., todos implementando la interfaz `parser.IExprContext`. Sin la etiqueta, todas las alternativas de `expr` compartirían un único contexto genérico, indistinguible entre sí en tiempo de ejecución.

Esta decisión de la gramática es la que hace posible todo el Tema 3 (el traductor sin Visitor). Memorizá esta frase, te la van a preguntar en la defensa: *"si agregás una alternativa a `expr` sin etiquetarla con `#nombre`, ANTLR4 no genera un tipo de contexto propio, y el `type-switch` del traductor no puede distinguirla de las demás."*

### 2.4 Ejercicio guiado: agregá un operador de potencia (`^`)

Este es el ejercicio más completo que podés practicar en esta etapa — toca gramática, generación y (más adelante, en el Tema 3) el traductor.

**Paso 1 — la gramática.** En `grammar/VLangCherry.g4`, agregá una alternativa nueva a `expr`, con su propia etiqueta, en la posición de precedencia que quieras darle (por ejemplo, más fuerte que `*`/`/`/`%`):

```antlr
expr
    : expr '(' listaArgumentos? ')'            # exprLlamada
    | expr '[' expr ']'                        # exprIndice
    | expr '.' ID                              # exprCampo
    | '(' expr ')'                             # exprParentesis
    | op=('!'|'-') expr                        # exprUnario
    | expr '^' expr                            # exprPotencia
    | expr op=('*'|'/'|'%') expr               # exprMultiplicativa
    ...
```

Y agregá el token léxico junto a los demás operadores:

```antlr
POTENCIA: '^';
```

**Paso 2 — regenerar.** Corré el comando de la sección 1.3. Vas a ver que `server/internal/parser/` cambió por completo (los 4 archivos se sobrescriben) y que ahora existe un `ExprPotenciaContext`.

**Paso 3 — el traductor y el runtime.** Estos los vas a completar recién cuando llegues al Tema 3 (agregar el `case` en `traducirExpr`) y al Tema 4 (implementar `Potencia` en `operaciones.go`, siguiendo el patrón de `Multiplicacion`). Volvé a este ejercicio después de leer esas secciones.

### 2.5 Preguntas para practicar de memoria

- ¿Qué generador de los tres del curso (JFlex+CUP, Jison, ANTLR4) separa léxico de sintaxis en archivos distintos? *(Solo JFlex+CUP.)*
- ¿Qué algoritmo de parsing usa ANTLR4, y en qué se diferencia de LALR? *(LL adaptativo, ALL(\*): decide la alternativa en tiempo de análisis examinando la entrada que haga falta, en vez de tablas LALR precomputadas en la generación.)*
- ¿Por qué `func` y no `fn`? ¿Por qué comillas dobles y no interpolación con comilla simple? Repasá `docs/gramatica.txt`, puntos 1 y 2 — son desambiguaciones reales frente a inconsistencias del enunciado, no gustos de diseño.

---

## 3. Tema: del parse tree al AST propio

### 3.1 Por qué no interpretar directo sobre el parse tree de ANTLR

El parse tree que produce ANTLR (`parser.IProgramaContext` y toda su jerarquía) está atado a los tipos que genera la herramienta — si regenerás la gramática, esos tipos pueden cambiar. VLangCherry traduce ese parse tree a un AST propio, definido en `server/internal/ast/ast.go`, desacoplado de ANTLR y con control total sobre el formato del reporte visual (sección 8.3 del enunciado).

Todo nodo del AST implementa una interfaz mínima:

```go
type Nodo interface {
    Tipo() string
    Pos() (int, int)
}
```

Hay ~25 tipos concretos (`Programa`, `DeclStruct`, `DeclFuncion`, `Bloque`, `SentenciaIf`, `ExprBinaria`, `LiteralStruct`, ...), todos embebiendo un `base{Linea, Columna int}` que resuelve `Pos()` automáticamente.

### 3.2 La decisión central: type-switch en vez de Visitor generado

ANTLR4, con `-visitor`, genera una interfaz `VLangCherryVisitor` con un método `Visit...` por cada regla y alternativa etiquetada — unas 60 firmas. VLangCherry **no implementa esa interfaz**. En vez de eso, `server/internal/traductor/traductor.go` aprovecha exactamente lo que viste en el Tema 2 (cada alternativa etiquetada ya es un tipo Go concreto) y usa un `switch ctx.(type)` directo:

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
    // ... un case por cada alternativa etiquetada de expr
    }
    return nil
}
```

Es funcionalmente equivalente al patrón *Visitor* del Dragón (doble despacho `accept`/`visit`), pero con menos código repetido: en vez de una llamada indirecta a través de `accept`, Go ya sabe distinguir el tipo dinámico de la interfaz con un `switch`.

**Práctica: escribí vos el helper de posición.** Es una función de una línea que vas a usar en cada función `traducir...`:

```go
func pos(ctx antlr.ParserRuleContext) (int, int) {
    t := ctx.GetStart()
    return t.GetLine(), t.GetColumn() + 1
}
```

`GetColumn()` de ANTLR es 0-based; el `+1` es para reportar columnas empezando en 1, como espera el reporte de errores (sección 8.1 del enunciado).

### 3.3 Un mismo nodo para leer y para escribir

Fijate en `traducirLugar` (para el lado izquierdo de una asignación) y en `traducirExpr` (para lectura): ambas producen exactamente los mismos tipos de nodo (`ast.Identificador`, `ast.ExprIndice`, `ast.ExprCampo`) para las mismas construcciones sintácticas (`ID`, `lugar[expr]`, `lugar.campo`):

```go
func traducirLugar(ctx parser.ILugarContext) ast.Nodo {
    switch lc := ctx.(type) {
    case *parser.LugarIdContext:
        l, c := pos(lc)
        n := &ast.Identificador{Nombre: lc.ID().GetText()}
        n.Linea, n.Columna = l, c
        return n
    case *parser.LugarIndiceContext:
        l, c := pos(lc)
        n := &ast.ExprIndice{Base: traducirLugar(lc.Lugar()), Indice: traducirExpr(lc.Expr())}
        n.Linea, n.Columna = l, c
        return n
    // ...
    }
}
```

El AST **no distingue** "leer" de "mutar". Esa decisión la toma el intérprete más adelante (Tema 4, función `resolverLugar`), no el traductor ni el AST. Si te preguntan por qué el AST es tan simple, esta es la respuesta.

### 3.4 El reporte de AST es gratis: recorrido por `reflect`

`server/internal/ast/grafo.go` define `ConstruirGrafo(raiz Nodo) Grafo`, que recorre **cualquier** nodo usando el paquete `reflect` de Go: no hace falta tocar este archivo cuando agregás un tipo de nodo nuevo, siempre que sus campos sean exportados (mayúscula inicial) y de un tipo reconocido:

- **Puntero/interfaz que implementa `Nodo`** → genera un nodo hijo y una arista con el nombre del campo.
- **Slice** → cada elemento se procesa individualmente (`Campos[0]`, `Campos[1]`, ...).
- **Struct "contenedor transparente"** (`RamaIf`, `CasoSwitch`, `Parametro`, `CampoStruct`, `CampoValorLiteral` — no implementan `Nodo` por sí mismos) → se recorren sus campos propagando un nombre compuesto (`Ramas[0].Condicion`).
- **`TipoAST`** → se serializa como texto y se pega a la etiqueta del nodo padre.
- **Escalar** → se agrega como texto (`"campo: valor"`).

### 3.5 Ejercicio guiado: agregá un tipo de nodo nuevo

Seguí con el ejercicio del operador `^` del Tema 2. Como `expr '^' expr` es una operación binaria más, **no necesitás un tipo de nodo nuevo**: reutilizá `ast.ExprBinaria` (que ya tiene `Operador string`), igual que `+`, `-`, `*`, etc. Agregá el `case` en `traducirExpr`:

```go
case *parser.ExprPotenciaContext:
    return traducirBinaria(ec, ec.AllExpr(), "^")
```

Si en cambio quisieras un nodo **realmente nuevo** (por ejemplo, una sentencia `do-while` que DataForge y VLangCherry no tienen), el patrón sería: (1) definir el struct en `ast.go` embebiendo `base` e implementando `Tipo()`; (2) agregar el `case` correspondiente en el traductor; (3) no tocar `grafo.go` — el recorrido por reflection ya lo incorpora.

---

## 4. Tema: tipos estáticos, structs y métodos "por referencia" gratis

### 4.1 El modelo de valores

`server/internal/runtime/valores.go` define el valor en tiempo de ejecución de **cualquier** expresión VLangCherry:

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

El enunciado (sección 7.1) exige que structs y slices se pasen **por referencia**, y los tipos primitivos **por valor**. VLangCherry no implementa dos mecanismos de paso de parámetros distintos: aprovecha que `Slice` y `Struct` **ya son punteros** dentro de `Valor`. Cuando Go copia un `Valor` (que hace automáticamente al pasar un parámetro o hacer una asignación, porque `Valor` es un struct normal), copia el puntero — no la estructura apuntada. Dos variables VLangCherry que "comparten" un struct terminan apuntando al mismo `StructVal` real en memoria. Los primitivos (`int`, `float64`, `string`, `bool`, `rune`) se copian por valor a través de los campos escalares `I`/`F`/`S`/`B`/`R`.

Practicá decir esto de memoria — es la pregunta más probable de toda la defensa (ver Tema 10).

### 4.2 Receptor por valor vs. por puntero: mismo resultado, matiz distinto

Con `entradas/ejemplo2_structs.vch`:

```go
func (p Persona) Saludar() string {
    return "Hola, soy " + p.Nombre
}

func (p *Persona) Cumplir() {
    p.Edad = p.Edad + 1
}
```

`invocarFuncion` (en `interprete.go`) declara el receptor en el entorno de la función igual que cualquier parámetro:

```go
entFn := NuevoEntorno(fn.Nombre, in.entGlobal)
if receptor != nil && fn.ReceptorNombre != "" {
    rec := *receptor
    if !fn.ReceptorPuntero {
        rec = ClonarPorValor(rec) // receptor por valor: trabaja sobre una copia
    }
    entFn.Declarar(fn.ReceptorNombre, rec)
}
```

`Saludar` no muta nada. `Cumplir` (receptor **por puntero**, `(p *Persona)`) incrementa `Edad` y el cambio persiste afuera, porque el receptor comparte el mismo `StructVal` que el llamador. La diferencia con un receptor **por valor** (`(p Persona)`) es real y observable: ahí `invocarFuncion` clona el struct con `ClonarPorValor` —copia en profundidad, igual que Go copia un struct al pasarlo por valor—, de modo que cualquier mutación dentro del método (reasignar el struct, un campo, o un campo anidado) **no** se ve afuera. Si te piden demostrarlo en vivo: cambiá `func (p *Persona) Cumplir()` por `func (p Persona) Cumplir()`, volvé a correr `entradas/ejemplo2_structs.vch`, y mostrá que ahora `persona.Edad` **no** refleja el incremento. (Este comportamiento correcto es reciente: hasta la auditoría del 2026-07-23, el receptor por valor compartía el `StructVal` por error —`ReceptorPuntero` no se leía—; ver Manual Técnico 8.8, hallazgo A4.)

### 4.3 La excepción: `append` siempre es nuevo

`append` (`nativas.go`) es la única excepción intencional a "todo se comparte por puntero": siempre devuelve un `SliceVal` **nuevo**, igual que en Go real:

```go
case "append":
    if len(args) != 2 {
        in.errorSemantico("append espera 2 argumentos", linea, col)
        return Valor{}, false
    }
    v := args[0]
    if v.Tipo.Base != TSlice {
        in.errorSemantico("append espera un slice como primer argumento", linea, col)
        return Valor{}, false
    }
    nuevos := append(append([]Valor{}, v.Slice.Elems...), args[1])
    return NuevoSlice(v.Slice.TipoElem, nuevos), true
```

Por eso hace falta reasignar explícitamente: `numeros = append(numeros, 60)`. Si no marcás esta excepción en la defensa, alguien te puede acorralar con "¿entonces por qué necesito reasignar si todo es por referencia?".

### 4.4 Compatibilidad de tipos: exacta, salvo `int`↔`float64`

`tipoCompatible`/`coercionar` (`interprete.go`) implementan la regla real de asignación:

```go
func tipoCompatible(declarado Tipo, v Valor) bool {
    if declarado.Igual(v.Tipo) {
        return true
    }
    if EsNumerico(declarado) && EsNumerico(v.Tipo) {
        return true
    }
    if v.Tipo.Base == TNil && (declarado.Base == TSlice || declarado.Base == TStruct) {
        return true
    }
    return false
}
```

El enunciado (3.6) exige tipo **exacto**, pero las tablas aritméticas (4.3) muestran promoción automática `int + float64 -> float64`. Se decidió relajar la mezcla `int`/`float64` en cualquier asignación — documentado en `docs/gramatica.txt` y en el propio código. Esto NO se extiende a ninguna otra combinación: asignar `"hola"` a una variable `int` sigue siendo error semántico.

### 4.5 Ejercicio guiado: agregá una función nativa

**Dónde:** `server/internal/runtime/nativas.go`. Tocás exactamente dos lugares:

1. Agregá el nombre al `switch` de `EsNombreNativa`.
2. Agregá un `case` nuevo en el `switch` de `llamarNativa`, siguiendo el patrón de `len`: validar cantidad de argumentos con `in.errorSemantico(...)` si falla, validar el tipo con un `switch v.Tipo.Base`, y devolver `(Valor, bool)`.

Practicá con `sum(numeros)` que sume todos los elementos de un `[]int`:

```go
func EsNombreNativa(nombre string) bool {
    switch nombre {
    case "print", "println", "len", "append", "indexOf", "join", "Atoi", "parseFloat", "typeOf",
        "int", "float64", "string", "bool", "rune", "sum": // <- agregado
        return true
    }
    return false
}
```

```go
case "sum":
    if len(args) != 1 {
        in.errorSemantico("sum espera 1 argumento", linea, col)
        return Valor{}, false
    }
    v := args[0]
    if v.Tipo.Base != TSlice || v.Slice.TipoElem.Base != TInt {
        in.errorSemantico("sum solo acepta []int", linea, col)
        return Valor{}, false
    }
    var total int64
    for _, e := range v.Slice.Elems {
        total += e.I
    }
    return Valor{Tipo: TipoInt(), I: total}, true
```

No hace falta tocar la gramática ni el traductor: las llamadas a función ya están resueltas genéricamente en `evaluarLlamada` (`interprete.go`) — `EsNombreNativa` es solo una lista de nombres reservados, y `llamarNativa` es el único punto que despacha por nombre.

**Cerrá el ejercicio del operador `^`:** implementá `Potencia` en `operaciones.go` siguiendo el patrón exacto de `Multiplicacion`, y conectala en `evaluarBinaria` (`interprete.go`) agregando el `case "^"`.

---

## 5. Tema: slices 1D/2D

### 5.1 No hay un tipo especial para "2D"

`SliceVal` (`valores.go`):

```go
type SliceVal struct {
    Elems    []Valor
    TipoElem Tipo
}
```

Un `[][]int` es simplemente un `SliceVal` cuyo `TipoElem` es a su vez `Tipo{Base: TSlice, Elemento: &Tipo{Base: TInt}}` — la misma estructura recursiva, aplicada dos veces. Lo mismo pasa en la gramática: `tipoSlice: '[' ']' tipo` está definida en términos de sí misma (`tipo` puede volver a ser un `tipoSlice`). No hay ninguna rama de código separada para "2D" ni en el runtime ni en la gramática.

### 5.2 El literal de slice: filas, no una lista plana

`evaluarLiteralSlice` (`interprete.go`) distingue el caso 1D del 2D mirando `n.TipoElem.EsSlice`:

```go
func (in *Interprete) evaluarLiteralSlice(n *ast.LiteralSlice, ent *Entorno) (Valor, bool) {
    if n.TipoElem.EsSlice {
        tipoFila := in.resolverTipoAST(*n.TipoElem.Elemento, n.Linea, n.Columna)
        var filas []Valor
        for _, fila := range n.Filas {
            elems, ok := in.evaluarListaConTipo(fila, tipoFila, ent, n.Linea, n.Columna)
            if !ok {
                return Valor{}, false
            }
            filas = append(filas, NuevoSlice(tipoFila, elems))
        }
        tipoDeFila := in.resolverTipoAST(n.TipoElem, n.Linea, n.Columna)
        return NuevoSlice(tipoDeFila, filas), true
    }
    // ... caso 1D
}
```

En la gramática, el literal 2D usa **filas** (`literalSlice: '[' ']' tipo '{' filaSlice (',' filaSlice)* ','? '}'`, con `filaSlice: '{' listaArgumentos? '}'`), no una lista plana de números. Con `entradas/ejemplo3_slices.vch`:

```go
mtx := [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}
for f, fila in mtx {
    for c, val in fila {
        print("mtx[", f, "][", c, "] =", val)
    }
}
```

Cada `fila` en el `for` externo es, en sí misma, un `[]int` completo — no hay ningún caso especial para "el elemento de un slice 2D".

### 5.3 Las 4 funciones nativas sobre slices

`len`, `append`, `indexOf`, `join` (todas en `nativas.go`) son la sección 7.2 del enunciado. Repasalas con atención porque son terreno fértil para el ejercicio "agregá una función nativa" (Tema 4.5) en la defensa — `indexOf` usa `valoresIguales` (la misma función que implementa `==`) para comparar cada elemento, y `join` valida explícitamente que el `TipoElem` del slice sea `string` antes de concatenar.

### 5.4 Ejercicio guiado: extendé `sum` a slices 2D

Tomá la función `sum` que escribiste en el Tema 4.5 y extendela para que, si recibe un `[][]int`, sume todos los elementos de todas las filas. Vas a necesitar revisar `v.Slice.TipoElem.Base` para decidir si estás ante un `[]int` o un `[][]int`, y recursar sobre cada fila (`e.Slice.Elems`) en el segundo caso. Este ejercicio conecta directamente el Tema 4 (cómo se agrega una nativa) con el Tema 5 (cómo se representa un slice 2D) — es exactamente el tipo de combinación que un jurado puede pedir en vivo.

---

## 6. Tema: control de flujo y recursión

### 6.1 `if/else if/else`

`ejecutarIf` (`interprete.go`) evalúa las ramas en orden y ejecuta la primera cuya condición sea `true`, cada una en su propio entorno anidado:

```go
func (in *Interprete) ejecutarIf(n *ast.SentenciaIf, ent *Entorno) Senal {
    for _, rama := range n.Ramas {
        cond, ok := in.evaluar(rama.Condicion, ent)
        if !ok {
            return Senal{Tipo: SenalNinguna}
        }
        if cond.Tipo.Base != TBool {
            l, c := rama.Condicion.Pos()
            in.errorSemantico("la condición del if debe ser de tipo bool", l, c)
            return Senal{Tipo: SenalNinguna}
        }
        if cond.B {
            return in.ejecutarBloque(rama.Cuerpo, NuevoEntorno(ent.Nombre, ent))
        }
    }
    if n.Else != nil {
        return in.ejecutarBloque(n.Else, NuevoEntorno(ent.Nombre, ent))
    }
    return Senal{Tipo: SenalNinguna}
}
```

Notá el tipo `Senal`: cada función `ejecutar...` devuelve una `Senal{Tipo, Valor}` que indica si el flujo normal se interrumpió (`SenalBreak`, `SenalContinue`, `SenalReturn`) o no (`SenalNinguna`). `ejecutarBloque` corta la ejecución del resto de sentencias en cuanto recibe cualquier señal distinta de `SenalNinguna`, y la propaga hacia arriba — así es como un `return` dentro de un `if` anidado dentro de un `for` termina saliendo de la función completa.

### 6.2 `switch` sin fall-through

`ejecutarSwitch` compara (`==`) la expresión del switch contra cada `case` en orden; el primero que matchea gana, y no hay encadenamiento entre casos:

```go
func (in *Interprete) ejecutarSwitch(n *ast.SentenciaSwitch, ent *Entorno) Senal {
    in.profSwitch++
    defer func() { in.profSwitch-- }()

    val, ok := in.evaluar(n.Expr, ent)
    if !ok {
        return Senal{Tipo: SenalNinguna}
    }
    for _, caso := range n.Casos {
        if caso.Valor == nil {
            continue
        }
        cv, ok := in.evaluar(caso.Valor, ent)
        if !ok {
            continue
        }
        igual, err := valoresIguales(val, cv)
        if err != nil {
            l, c := caso.Valor.Pos()
            in.errorSemantico(err.Error(), l, c)
            continue
        }
        if !igual {
            continue
        }
        senal := in.ejecutarListaSentencias(caso.Sentencias, NuevoEntorno(ent.Nombre, ent))
        if senal.Tipo == SenalBreak {
            return Senal{Tipo: SenalNinguna}
        }
        return senal
    }
    // ... corre el default si ningun case matcheo
}
```

Esto **no** es una limitación de la implementación: la sección 4.7.2 del enunciado dice explícitamente "el break implícito está incluido al final de cada case" — es una decisión de diseño real del lenguaje, distinta del switch de CompInterpreter (otro proyecto del curso, con fall-through real).

### 6.3 Las tres formas de `for`

`ejecutarFor` cubre "condición" (equivalente a un `while`), "clásico" (`init; condición; actualización`) y "rango" (`for i, v in slice`), las tres con su propio manejo de entornos:

- **condición**: reevalúa la condición en cada vuelta, un solo entorno para el cuerpo por iteración.
- **clásico**: un entorno `entFor` compartido por todas las iteraciones (para que `i` sobreviva entre vueltas), más un entorno nuevo por cuerpo de iteración.
- **rango**: un entorno **nuevo por iteración** (`entIter`) que declara `VarIndice` y `VarValor` en cada vuelta.

### 6.4 `profLoop`/`profSwitch` y la recursión

```go
type Interprete struct {
    // ...
    profLoop   int
    profSwitch int
}
```

Estos dos contadores llevan cuántos `for`/`switch` anidados rodean la sentencia que se está ejecutando ahora mismo. `ejecutarFor` y `ejecutarSwitch` los incrementan al entrar y los decrementan con `defer` al salir. `invocarFuncion` los **resetea a 0** al entrar a una función o método, y los restaura al salir:

```go
loopPrevio, switchPrevio := in.profLoop, in.profSwitch
in.profLoop, in.profSwitch = 0, 0
senal := in.ejecutarBloque(fn.Cuerpo, entFn)
in.profLoop, in.profSwitch = loopPrevio, switchPrevio
```

¿Por qué? Para que un `break` "perdido" dentro del cuerpo de una función invocada desde dentro de un `for` no herede el contexto de ciclo del llamador. Practicá explicar este detalle: es fino, y demuestra dominio real del alcance semántico, no solo "el break funciona".

La recursión (`entradas/ejemplo5_funciones.vch`, `factorial`/`fibonacci`) funciona sin ningún mecanismo especial: cada llamada a `invocarFuncion` arma su propio entorno (`entFn`) con `Padre: in.entGlobal`, así que cada nivel de recursión tiene sus propias variables locales, aisladas de las demás. La llamada hacia adelante (`main` invocando `factorial`, declarada más abajo en el archivo) funciona porque `Interpretar` registra **todas** las funciones en `global.Funciones` antes de ejecutar cualquier cuerpo (ver Tema 6.6).

### 6.5 Caso de estudio: los 4 hallazgos de la auditoría (2026-07-21)

Esta es la evidencia más fuerte de que el equipo entiende la diferencia entre "el programa corre" y "el programa valida lo que el enunciado exige". Los 4 hallazgos, con el código antes/después real:

**Hallazgo 1 — `break`/`continue` fuera de contexto.** Antes de la corrección, `ejecutarSentencia` simplemente devolvía la señal sin verificar nada:

```go
// ANTES (sin validar)
case *ast.SentenciaBreak:
    return Senal{Tipo: SenalBreak}
```

```go
// DESPUES (real, interprete.go)
case *ast.SentenciaBreak:
    // 4.8.1: "la sentencia break solo se puede usar dentro de un bucle
    // (for) o un switch... si se encuentra fuera se considerará un error".
    if in.profLoop == 0 && in.profSwitch == 0 {
        in.errorSemantico("la sentencia \"break\" no puede usarse fuera de un ciclo o un switch", n.Linea, n.Columna)
        return Senal{Tipo: SenalNinguna}
    }
    return Senal{Tipo: SenalBreak}
```

`continue` sigue el mismo patrón, pero solo revisa `profLoop` (4.8.2: `continue` solo es válido dentro de un `for`, no dentro de un `switch` suelto).

**Hallazgo 2 — comparación relacional de strings.** La sección 4.4.1 del enunciado dice "las comparaciones entre cadenas se hacen lexicográficamente" bajo el título de igualdad, pero describe un criterio de ORDEN. Faltaba aplicarlo a `<`, `<=`, `>=`, `>` (no solo a `==`/`!=`). La función `Relacional` (`operaciones.go`) ya lo cubre:

```go
func Relacional(op string, a, b Valor) (Valor, error) {
    var cmp int
    switch {
    case EsNumerico(a.Tipo) && EsNumerico(b.Tipo):
        cmp = compararFloats(aFloat(a), aFloat(b))
    case a.Tipo.Base == TRune && b.Tipo.Base == TRune:
        cmp = int(a.R) - int(b.R)
    case a.Tipo.Base == TString && b.Tipo.Base == TString:
        cmp = strings.Compare(a.S, b.S)
    default:
        return Valor{}, fmt.Errorf("no se puede comparar (%s) entre %s y %s", op, a.Tipo, b.Tipo)
    }
    // ...
}
```

Este hallazgo está documentado con más detalle en `docs/gramatica.txt`, punto 8.

**Hallazgo 3 — colisión de nombres en el ámbito global.** El enunciado (7.1) exige que función, variable y struct no compartan nombre. Antes solo se detectaba redeclaración **dentro** de la misma categoría. `Interpretar` ahora chequea la colisión **entre** categorías al registrar cada función:

```go
if _, dup := in.global.Structs[f.Nombre]; dup {
    in.errorSemantico("el nombre \""+f.Nombre+"\" ya está en uso por un struct", f.Linea, f.Columna)
    continue
}
```

y `ejecutarDeclaracion` la chequea de nuevo al declarar una variable global, contra structs y funciones:

```go
if ent == in.entGlobal {
    if _, dup := in.global.Structs[n.Nombre]; dup {
        in.errorSemantico("el nombre \""+n.Nombre+"\" ya está en uso por un struct", n.Linea, n.Columna)
        return
    }
    if _, dup := in.global.Funciones[n.Nombre]; dup {
        in.errorSemantico("el nombre \""+n.Nombre+"\" ya está en uso por una función", n.Linea, n.Columna)
        return
    }
}
```

**Hallazgo 4 — el `switch` que silenciaba errores de tipo.** Comparar un `int` contra un `struct` dentro de un `case` se descartaba en silencio (como si simplemente no coincidiera), inconsistente con cómo `==` reporta ese mismo error fuera de un switch. La corrección real está en el fragmento de `ejecutarSwitch` de la sección 6.2 de esta guía: en vez de ignorar el `error` que devuelve `valoresIguales`, ahora se reporta con `in.errorSemantico(err.Error(), l, c)`.

**Lo que tienen en común los 4:** ninguno es un fallo de ejecución — todos son huecos de validación que el enunciado exige cerrar. Si te preguntan "¿por qué tu intérprete detecta este error?" para cualquiera de los 4, tené lista la cita exacta de la sección del enunciado (4.8.1/4.8.2, 4.4.1, 7.1, 4.4.1).

### 6.6 El intérprete de dos pasadas

```go
func (in *Interprete) Interpretar(programa *ast.Programa) {
    for _, s := range programa.Structs { /* pasada 1: registra structs */ }
    for _, f := range programa.Funciones { /* pasada 2: registra funciones y metodos */ }
    in.entGlobal = NuevoEntorno("global", nil)
    for _, g := range programa.Globales { in.ejecutarDeclaracion(g, in.entGlobal) }
    mainFn, ok := in.global.Funciones["main"]
    // ... localizar y ejecutar main()
}
```

Sin las dos pasadas, un intérprete de una sola pasada que ejecutara cada declaración según aparece en el archivo fallaría al llamar `factorial()` desde `main()` si `factorial` está declarada más abajo — en ese punto todavía no la habría "visto". Las dos pasadas existen exactamente para eliminar esa dependencia del orden textual.

### 6.7 Ejercicio guiado: agregá una validación semántica nueva

Siguiendo el patrón exacto de los 4 hallazgos, la función clave es `in.errorSemantico(desc string, linea, col int)`, que acumula en `ListaErrores` sin abortar la ejecución. Practicá replicando el Hallazgo 1 completo: quitá la validación de `break` a propósito (dejá solo `return Senal{Tipo: SenalBreak}`), corré `entradas/ejemplo4_control.vch` (no debería fallar, porque ahí los `break`/`continue` están bien puestos), después escribí un `.vch` nuevo con un `break` fuera de cualquier `for`/`switch`, comprobá que NO se reporta error, y volvé a agregar la validación real. Mientras lo hacés, citá la sección exacta del enunciado (4.8.1) que la exige.

---

## 7. Tema: el servidor REST

### 7.1 Sin framework, a propósito

`server/cmd/servidor/main.go` usa únicamente `net/http` de la librería estándar de Go. Tres rutas:

| Método y ruta | Qué hace |
|---|---|
| `GET /` | Info básica (`{"nombre": "VLangCherry", "estado": "ok", "endpoint": "..."}`). |
| `GET /salud` | `{"estado": "ok"}` — para verificar que el servidor está arriba. |
| `POST /interpretar` | Recibe `{"codigo": "<fuente .vch>"}`, responde el resultado del análisis. |

El puerto es configurable por la variable de entorno `PORT`, con `"4100"` como default hardcodeado:

```go
puerto := os.Getenv("PORT")
if puerto == "" {
    puerto = "4100"
}
```

Todas las rutas pasan por el middleware `conCORS`, que agrega `Access-Control-Allow-Origin: *` — necesario porque en desarrollo el cliente (Vite) y el servidor (Go) corren en orígenes distintos.

### 7.2 El pipeline completo por petición

`server/internal/analizar/analizar.go` es donde se arma el pipeline entero: lexer → parser → traductor → intérprete, con un `Entorno` **nuevo en cada llamada** — nada de estado sobrevive entre peticiones:

```go
func Analizar(codigo string) Resultado {
    errores := runtime.NuevaListaErrores()

    entrada := antlr.NewInputStream(codigo)
    lexer := parser.NewVLangCherryLexer(entrada)
    lexer.RemoveErrorListeners()
    lexer.AddErrorListener(&oyenteErrores{errores: errores, esLexer: true})

    tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
    p := parser.NewVLangCherryParser(tokens)
    p.RemoveErrorListeners()
    p.AddErrorListener(&oyenteErrores{errores: errores, esLexer: false})

    arbolParse := p.Programa()

    var grafo ast.Grafo
    interprete := runtime.NuevoInterprete(errores)

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
    // ... arma el Resultado con errores, consola, simbolos, ast, dot
}
```

Notá dos cosas: (1) los listeners de error por defecto de ANTLR se **reemplazan** por `oyenteErrores`, que traduce cada error léxico/sintáctico de ANTLR al mismo formato `ErrorReporte` que usan los errores semánticos; (2) la traducción y la interpretación corren dentro de un `defer recover()` — si algo interno no contemplado hace panic (por ejemplo, un tipo de nodo del AST sin manejar en algún `switch`), se reporta como error semántico adicional en vez de tumbar el proceso completo. Es un detalle relevante en un servidor de larga duración que atiende múltiples peticiones.

### 7.3 El contrato JSON

```go
type Resultado struct {
    Errores       []runtime.ErrorReporte `json:"errores"`
    Consola       string                 `json:"consola"`
    ConsolaLineas []string               `json:"consolaLineas"`
    Simbolos      []runtime.FilaSimbolo  `json:"simbolos"`
    AST           ast.Grafo              `json:"ast"`
    Dot           string                 `json:"dot"`
}
```

Es, en su forma, el mismo contrato que usa CompInterpreter (otro proyecto del curso, con servidor Node/Express) — `{errores, consola, simbolos, ast}` — más `consolaLineas` (la consola ya partida en líneas, para que el cliente no tenga que hacer `split('\n')`) y `dot` (el equivalente Graphviz del AST, generado del lado del servidor).

### 7.4 Ejercicio guiado: agregá un endpoint nuevo

Practicá agregando `GET /version` que devuelva `{"version": "1.0"}`. Vas a necesitar: (1) un `mux.HandleFunc("/version", conCORS(...))` en `main.go`, siguiendo el patrón de `/salud`; (2) nada más — no hace falta tocar `analizar.go` porque este endpoint no ejecuta el pipeline. Esto te sirve para practicar la mecánica de agregar rutas sin arriesgar romper el endpoint principal antes de un ensayo de defensa.

---

## 8. Tema: el cliente React

### 8.1 Un solo archivo conoce el servidor

`client/src/api.js` es el único lugar que sabe la URL del backend:

```js
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:4100';

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

Todo el resto de los componentes solo conoce la **forma** del JSON que devuelve `interpretar(codigo)`, no de dónde viene. Si el servidor cambia de puerto o dominio, este es el único archivo que hay que tocar.

### 8.2 `App.jsx`: el estado que orquesta todo

`App.jsx` mantiene el estado de los archivos abiertos (`{id, nombre, contenido, sinGuardar}`), el archivo activo, el resultado de la última ejecución y el panel seleccionado. La función clave es `manejarEjecutar`:

```js
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

Fijate en la selección automática de pestaña: si hay errores, se muestra la pestaña Errores; si no, Consola — así el usuario ve inmediatamente lo relevante sin tener que buscar.

### 8.3 Componentes de responsabilidad única

- **`Toolbar.jsx`** / **`FileTabs.jsx`**: acciones (Nuevo/Abrir/Guardar/Ejecutar) y pestañas multi-archivo.
- **`Editor.jsx`**: numeración de línea, resalta las líneas con error (recibe `erroresPorLinea`, un `Set` calculado en `App.jsx` a partir de `resultado.errores`).
- **`Consola.jsx`**, **`TablaErrores.jsx`**, **`TablaSimbolos.jsx`**: presentación tabular; las dos tablas soportan clic en fila para saltar a la línea (`onIrALinea` → `editorRef.current.irALinea(num)`).
- **`AstGrafo.jsx`**: renderiza `{nodes, edges}` con **vis-network**, layout jerárquico `direction: 'UD'` (arriba hacia abajo), leyendo los colores de las variables CSS del tema activo (claro/oscuro) para que el grafo combine con el resto de la interfaz.

Este patrón de componentes está adaptado de otro proyecto del curso con arquitectura cliente-servidor equivalente: se reutiliza el diseño de interfaz, cambiando únicamente `api.js` para apuntar al servidor Go en vez de a un servidor Node/Express.

### 8.4 Ejercicio guiado: agregá un botón "Limpiar consola"

Practicá el patrón de estado de React que ya usa `App.jsx`: agregá un botón en `Toolbar.jsx` que dispare un `onLimpiarConsola` prop, manejado en `App.jsx` con un `setResultado` que preserve `errores`/`simbolos`/`ast` pero vacíe `consolaLineas`. Es un ejercicio chico pero real: te obliga a entender que `resultado` es un único objeto inmutable que se reemplaza completo en cada actualización de estado, no un conjunto de piezas sueltas.

---

## 9. Errores comunes reales al programar cada tema

- **Tema 2 (gramática):** agregar una alternativa a `expr` sin la etiqueta `# nombre` — el traductor no puede reconocerla y el `type-switch` simplemente no matchea nunca ese caso (falla en silencio, sin panic, porque `traducirExpr` devuelve `nil` al final si ningún `case` matchea).
- **Tema 2:** olvidar regenerar el parser después de tocar el `.g4` — vas a seguir viendo el comportamiento viejo aunque el archivo fuente ya esté actualizado, porque `internal/parser/` sigue siendo el generado anterior hasta que corrés el comando de ANTLR de nuevo.
- **Tema 3 (traductor):** confundir `traducirLugar` con `traducirExpr` — si usás `traducirExpr` para el lado izquierdo de una asignación, vas a perder la distinción que necesita `resolverLugar` más adelante en el intérprete.
- **Tema 4 (tipos):** olvidar coercionar (`coercionar(tipoDeclarado, v)`) después de validar `tipoCompatible` — la validación pasa, pero el valor queda con el tipo original (por ejemplo, un `int` sin convertir a `float64`), y las operaciones subsiguientes fallan de forma confusa.
- **Tema 4:** asumir que un receptor `(p Persona)` hace una copia "segura" — como `Valor.Struct` es puntero, el struct interior se sigue compartiendo aunque el receptor sea por valor.
- **Tema 5 (slices):** olvidar reasignar el resultado de `append` — es la excepción a la regla de referencia compartida; si no reasignás, la variable original queda intacta.
- **Tema 6 (control de flujo):** agregar una validación semántica nueva y reportar el error, pero **también** seguir ejecutando la rama que causó el error — revisá siempre que tu `errorSemantico(...)` vaya seguido de un `return` (o del control de flujo equivalente) para no generar errores en cascada.
- **Tema 6:** olvidar el `defer` al incrementar `profLoop`/`profSwitch` — si el contador no se decrementa al salir (por ejemplo, por un `return` temprano sin `defer`), un `for` que ya terminó sigue "contando" como si estuviera activo para el resto del programa.
- **Tema 7 (servidor):** agregar un endpoint nuevo y olvidar envolverlo en `conCORS` — el cliente lo va a poder llamar con `curl` pero el navegador lo va a bloquear por CORS.
- **Tema 8 (cliente):** actualizar `resultado` con un objeto parcial en vez de reemplazarlo completo — como varios componentes leen distintas partes de `resultado`, un estado parcial deja paneles desincronizados entre sí.

---

## 10. Para poder defender esto en vivo

El enunciado de la defensa oral pide que el equipo demuestre autoría real modificando código en vivo, no recitando diapositivas. Estas son las partes del código que cualquiera del equipo debería poder explicar **línea por línea, sin dudar**, porque son las que con más probabilidad les van a pedir que toquen:

1. **La regla `expr` de la gramática** (`grammar/VLangCherry.g4`) y por qué el orden de sus alternativas es la precedencia — y qué pasa si agregás una alternativa sin la etiqueta `# nombre`.
2. **El `type-switch` de `traducirExpr`** (`traductor.go`) y por qué reemplaza al patrón Visitor generado — poder nombrar la interfaz `VLangCherryVisitor` y explicar por qué NO se implementó.
3. **La frase completa sobre `Valor.Slice`/`Valor.Struct` como punteros** (Tema 4.1) — practicada de memoria, palabra por palabra, incluyendo la excepción de `append`.
4. **El bloque `profLoop == 0 && profSwitch == 0`** de la validación de `break` (Tema 6.5, Hallazgo 1) — y poder citar la sección exacta del enunciado (4.8.1) que lo exige, más el motivo de por qué `invocarFuncion` resetea los contadores.
5. **Las dos pasadas de `Interpretar`** (Tema 6.6) — y poder explicar, sin el código delante, qué pasaría si el intérprete fuera de una sola pasada (`main` fallaría al llamar a una función declarada más abajo).
6. **El `defer recover()` de `analizar.go`** (Tema 7.2) — por qué existe en un servidor de larga duración y qué pasa si no estuviera.

Si el equipo domina estos 6 puntos con soltura, cualquier variación que les pidan en vivo (agregar un operador, una función nativa, una validación semántica, un endpoint) va a ser una aplicación directa de un patrón que ya practicaron en esta guía. La chuleta completa para el día de la defensa, con tiempos sugeridos por etapa y las preguntas más probables del quiz real, está en `presentacion-vlangcherry/GuionClase.md` — esta guía es el material para **aprender a programar** cada tema; ese guion es el material para **presentarlo y defenderlo** el día del examen.
