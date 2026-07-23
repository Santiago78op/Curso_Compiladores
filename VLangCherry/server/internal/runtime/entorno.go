package runtime

import "vlangcherry/internal/ast"

// Entorno es una pila de ambitos anidados (4.1: bloques, funciones).
// Las variables se guardan como *Valor para que la asignacion mute in
// place la misma celda (necesario para lugar[i]/lugar.campo tambien).
type Entorno struct {
	Nombre    string // nombre de la funcion contenedora, para el reporte 8.2
	variables map[string]*Valor
	// mutable[n] es false para una variable declarada SIN 'mut': reasignarla
	// (=, +=, -=, ++, --) es error semantico (A1). Los parametros, variables
	// de for y receptores se declaran con Declarar -> mutables por defecto.
	mutable map[string]bool
	Padre   *Entorno
}

func NuevoEntorno(nombre string, padre *Entorno) *Entorno {
	return &Entorno{Nombre: nombre, variables: make(map[string]*Valor), mutable: make(map[string]bool), Padre: padre}
}

// Declarar crea una variable MUTABLE en ESTE ambito (4.1: los bloques
// anidados pueden reusar un nombre del ambito superior, ocultandolo -
// "shadowing"). Se usa para lo que siempre es reasignable: parametros,
// receptores y las variables de control del for.
func (e *Entorno) Declarar(nombre string, v Valor) {
	e.DeclararMut(nombre, v, true)
}

// DeclararMut crea una variable indicando si es reasignable (lleva 'mut').
func (e *Entorno) DeclararMut(nombre string, v Valor, mutable bool) {
	copia := v
	e.variables[nombre] = &copia
	e.mutable[nombre] = mutable
}

// EsMutable indica si el enlace mas cercano de "nombre" es reasignable.
// Camina la cadena de ambitos igual que Buscar (respeta el shadowing). El
// segundo retorno es false si la variable no existe en ningun ambito.
func (e *Entorno) EsMutable(nombre string) (mutable bool, existe bool) {
	for amb := e; amb != nil; amb = amb.Padre {
		if _, ok := amb.variables[nombre]; ok {
			return amb.mutable[nombre], true
		}
	}
	return false, false
}

// Buscar resuelve un identificador subiendo por la cadena de ambitos.
func (e *Entorno) Buscar(nombre string) (*Valor, bool) {
	for amb := e; amb != nil; amb = amb.Padre {
		if v, ok := amb.variables[nombre]; ok {
			return v, true
		}
	}
	return nil, false
}

// DeclaradoLocal indica si el nombre ya existe en ESTE ambito (para
// detectar redeclaracion en el mismo bloque).
func (e *Entorno) DeclaradoLocal(nombre string) bool {
	_, ok := e.variables[nombre]
	return ok
}

// Global agrupa el registro de structs y funciones/metodos: viven en un
// unico espacio de nombres global (6.1/7.1: "solo pueden declararse en el
// ambito global"), independiente de los ambitos de ejecucion.
type Global struct {
	Structs   map[string]*ast.DeclStruct
	Funciones map[string]*ast.DeclFuncion
	Metodos   map[string]map[string]*ast.DeclFuncion // tipoStruct -> nombreMetodo -> decl
}

func NuevoGlobal() *Global {
	return &Global{
		Structs:   make(map[string]*ast.DeclStruct),
		Funciones: make(map[string]*ast.DeclFuncion),
		Metodos:   make(map[string]map[string]*ast.DeclFuncion),
	}
}
