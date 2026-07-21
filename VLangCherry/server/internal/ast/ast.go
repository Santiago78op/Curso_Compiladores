// Package ast define los tipos de nodo del AST propio de VLangCherry
// (uno por construccion del lenguaje: declaraciones, sentencias y
// expresiones) y, en grafo.go, la conversion generica de cualquier Nodo
// a un grafo {nodes, edges} para el reporte AST (8.3).
package ast

// Nodo es la interfaz minima que cualquier nodo del AST debe cumplir.
// ConstruirGrafo (grafo.go) usa reflection sobre los structs concretos,
// asi que un nodo nuevo aparece en el reporte AST (8.3) sin mantenimiento
// extra, igual que en CompInterpreter (reportes/ast-grafo.js).
type Nodo interface {
	Tipo() string
	Pos() (int, int)
}

type base struct {
	Linea, Columna int
}

func (b base) Pos() (int, int) { return b.Linea, b.Columna }

// ---------------- Programa y declaraciones globales ----------------

type Programa struct {
	base
	Structs   []*DeclStruct
	Funciones []*DeclFuncion
	Globales  []*DeclVariable
}

func (n *Programa) Tipo() string { return "PROGRAMA" }

type CampoStruct struct {
	Nombre string
	Tipo   TipoAST
}

type DeclStruct struct {
	base
	Nombre string
	Campos []CampoStruct
}

func (n *DeclStruct) Tipo() string { return "STRUCT" }

type Parametro struct {
	Nombre string
	Tipo   TipoAST
}

type DeclFuncion struct {
	base
	Nombre          string
	ReceptorNombre  string // "" si es funcion libre
	ReceptorTipo    string // nombre del struct receptor
	ReceptorPuntero bool
	Parametros      []Parametro
	TipoRetorno     *TipoAST // nil = sin retorno
	Cuerpo          *Bloque
}

func (n *DeclFuncion) Tipo() string { return "FUNCION" }

// TipoAST es la representacion sintactica de un tipo (antes de resolverlo
// contra runtime.Tipo, que necesita saber que structs existen).
type TipoAST struct {
	EsSlice    bool
	Elemento   *TipoAST
	NombreBase string // primitivo o nombre de struct
}

// ---------------- Sentencias ----------------

type Bloque struct {
	base
	Sentencias []Nodo
}

func (n *Bloque) Tipo() string { return "BLOQUE" }

type DeclVariable struct {
	base
	Nombre   string
	TipoVar  *TipoAST // nil si es inferido (:=)
	Inferido bool
	Valor    Nodo // puede ser nil
}

func (n *DeclVariable) Tipo() string { return "DECLARACION" }

type Asignacion struct {
	base
	Lugar    Nodo
	Operador string // "=", "+=", "-="
	Valor    Nodo
}

func (n *Asignacion) Tipo() string { return "ASIGNACION" }

type IncrementoDecremento struct {
	base
	Lugar    Nodo
	Operador string // "++", "--"
}

func (n *IncrementoDecremento) Tipo() string { return "INCDEC" }

type ExpresionSentencia struct {
	base
	Expr Nodo
}

func (n *ExpresionSentencia) Tipo() string { return "EXPR_SENTENCIA" }

type RamaIf struct {
	Condicion Nodo
	Cuerpo    *Bloque
}

type SentenciaIf struct {
	base
	Ramas []RamaIf // if, else-if...
	Else  *Bloque  // nil si no hay else
}

func (n *SentenciaIf) Tipo() string { return "IF" }

type CasoSwitch struct {
	Valor      Nodo // nil = default
	Sentencias []Nodo
}

type SentenciaSwitch struct {
	base
	Expr  Nodo
	Casos []CasoSwitch
}

func (n *SentenciaSwitch) Tipo() string { return "SWITCH" }

// SentenciaFor cubre las 3 formas (4.7.3): condicion / clasico / rango.
type SentenciaFor struct {
	base
	Forma         string // "condicion" | "clasico" | "rango"
	Condicion     Nodo   // forma condicion/clasico
	Init          Nodo   // forma clasico (DeclVariable | Asignacion), puede ser nil
	Actualizacion Nodo   // forma clasico (Asignacion | IncrementoDecremento), puede ser nil
	VarIndice     string // forma rango
	VarValor      string // forma rango
	Iterable      Nodo   // forma rango
	Cuerpo        *Bloque
}

func (n *SentenciaFor) Tipo() string { return "FOR" }

type SentenciaBreak struct{ base }

func (n *SentenciaBreak) Tipo() string { return "BREAK" }

type SentenciaContinue struct{ base }

func (n *SentenciaContinue) Tipo() string { return "CONTINUE" }

type SentenciaReturn struct {
	base
	Valor Nodo // nil = return sin valor
}

func (n *SentenciaReturn) Tipo() string { return "RETURN" }

// ---------------- Expresiones ----------------

type Identificador struct {
	base
	Nombre string
}

func (n *Identificador) Tipo() string { return "IDENTIFICADOR" }

type LiteralEntero struct {
	base
	Valor int64
}

func (n *LiteralEntero) Tipo() string { return "LITERAL_ENTERO" }

type LiteralDecimal struct {
	base
	Valor float64
}

func (n *LiteralDecimal) Tipo() string { return "LITERAL_DECIMAL" }

type LiteralCadena struct {
	base
	Valor string
}

func (n *LiteralCadena) Tipo() string { return "LITERAL_CADENA" }

type LiteralRune struct {
	base
	Valor rune
}

func (n *LiteralRune) Tipo() string { return "LITERAL_RUNE" }

type LiteralBool struct {
	base
	Valor bool
}

func (n *LiteralBool) Tipo() string { return "LITERAL_BOOL" }

type LiteralNil struct{ base }

func (n *LiteralNil) Tipo() string { return "LITERAL_NIL" }

type LiteralSlice struct {
	base
	TipoElem TipoAST
	Filas    [][]Nodo // 1 fila = slice simple; N filas = [][]T
}

func (n *LiteralSlice) Tipo() string { return "LITERAL_SLICE" }

type CampoValorLiteral struct {
	Nombre string
	Valor  Nodo
}

type LiteralStruct struct {
	base
	NombreStruct string
	Campos       []CampoValorLiteral
}

func (n *LiteralStruct) Tipo() string { return "LITERAL_STRUCT" }

type ExprUnaria struct {
	base
	Operador string // "!" | "-"
	Operando Nodo
}

func (n *ExprUnaria) Tipo() string { return "UNARIA" }

type ExprBinaria struct {
	base
	Operador  string
	Izquierda Nodo
	Derecha   Nodo
}

func (n *ExprBinaria) Tipo() string { return "BINARIA_" + n.Operador }

type ExprIndice struct {
	base
	Base   Nodo
	Indice Nodo
}

func (n *ExprIndice) Tipo() string { return "INDICE" }

type ExprCampo struct {
	base
	Base   Nodo
	Nombre string
}

func (n *ExprCampo) Tipo() string { return "CAMPO" }

type ExprLlamada struct {
	base
	Callee     Nodo
	Argumentos []Nodo
}

func (n *ExprLlamada) Tipo() string { return "LLAMADA" }
