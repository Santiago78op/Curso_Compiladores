package runtime

// ListaErrores recolecta TODOS los errores (lexico/sintactico/semantico)
// de una ejecucion, tal como pide 8.1 ("todos los errores se deberan
// recolectar").
type ErrorReporte struct {
	Tipo        string `json:"tipo"`
	Descripcion string `json:"descripcion"`
	Linea       int    `json:"linea"`
	Columna     int    `json:"columna"`
}

type ListaErrores struct {
	Errores []ErrorReporte
}

func NuevaListaErrores() *ListaErrores {
	return &ListaErrores{Errores: []ErrorReporte{}}
}

func (l *ListaErrores) agregar(tipo, desc string, linea, columna int) {
	l.Errores = append(l.Errores, ErrorReporte{Tipo: tipo, Descripcion: desc, Linea: linea, Columna: columna})
}

func (l *ListaErrores) Lexico(desc string, linea, columna int) {
	l.agregar("Léxico", desc, linea, columna)
}
func (l *ListaErrores) Sintactico(desc string, linea, columna int) {
	l.agregar("Sintáctico", desc, linea, columna)
}
func (l *ListaErrores) Semantico(desc string, linea, columna int) {
	l.agregar("Semántico", desc, linea, columna)
}
func (l *ListaErrores) Hay() bool { return len(l.Errores) > 0 }

// HayDeEntrada informa si ya se reporto algun error lexico o sintactico,
// es decir, si el parse tree que sigue esta incompleto. Lo consulta el
// pipeline para no convertir el panic esperable de traducir un arbol roto
// en un "error interno" que confunde al usuario (hallazgo B3).
func (l *ListaErrores) HayDeEntrada() bool {
	for _, e := range l.Errores {
		if e.Tipo == "Léxico" || e.Tipo == "Sintáctico" {
			return true
		}
	}
	return false
}
