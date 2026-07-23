package runtime

import (
	"fmt"
	"strconv"
	"strings"
)

// EsNombreNativa detecta si un identificador de llamada es una funcion
// embebida (7.2) o una conversion de tipo estilo Go (float64(x), int(x)...
// necesaria porque 3.6/4.2.1 exige asignacion de tipo EXACTO: el enunciado
// muestra f64(10+1) como conversion explicita pero no dedica una seccion
// a "casteos" como los otros proyectos del curso - se resuelve con esta
// sintaxis de conversion tipo Go, consistente con el unico ejemplo real).
func EsNombreNativa(nombre string) bool {
	switch nombre {
	case "print", "println", "len", "append", "indexOf", "join", "Atoi", "parseFloat", "typeOf",
		"int", "float64", "string", "bool", "rune":
		return true
	}
	return false
}

// llamarNativa despacha las funciones embebidas de 7.2 (print/println,
// len, append, indexOf, join, Atoi, parseFloat, typeOf) mas las
// conversiones de tipo estilo Go, que delega en convertirTipo.
func (in *Interprete) llamarNativa(nombre string, args []Valor, linea, col int) (Valor, bool) {
	switch nombre {
	case "print", "println":
		partes := make([]string, len(args))
		for i, a := range args {
			partes[i] = AImprimir(a)
		}
		in.consola = append(in.consola, strings.Join(partes, " "))
		return Valor{Tipo: TipoVoid()}, true

	case "len":
		if len(args) != 1 {
			in.errorSemantico("len espera 1 argumento", linea, col)
			return Valor{}, false
		}
		v := args[0]
		switch v.Tipo.Base {
		case TSlice:
			if v.Slice == nil {
				return Valor{Tipo: TipoInt(), I: 0}, true
			}
			return Valor{Tipo: TipoInt(), I: int64(len(v.Slice.Elems))}, true
		case TString:
			return Valor{Tipo: TipoInt(), I: int64(len(v.S))}, true
		}
		in.errorSemantico("len solo acepta slice o string, se recibió "+v.Tipo.String(), linea, col)
		return Valor{}, false

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
		// Un slice declarado sin inicializar (ValorPorDefecto) llega con
		// v.Slice == nil pero con el tipo de elemento en v.Tipo.Elemento:
		// append sobre el es como append de Go sobre un nil slice -> uno de 1.
		var previos []Valor
		tipoElem := TipoNil()
		if v.Tipo.Elemento != nil {
			tipoElem = *v.Tipo.Elemento
		}
		if v.Slice != nil {
			previos = v.Slice.Elems
			tipoElem = v.Slice.TipoElem
		}
		nuevos := append(append([]Valor{}, previos...), args[1])
		return NuevoSlice(tipoElem, nuevos), true

	case "indexOf":
		if len(args) != 2 {
			in.errorSemantico("indexOf espera 2 argumentos", linea, col)
			return Valor{}, false
		}
		v := args[0]
		if v.Tipo.Base != TSlice || v.Slice == nil {
			in.errorSemantico("indexOf espera un slice como primer argumento", linea, col)
			return Valor{}, false
		}
		for i, e := range v.Slice.Elems {
			if igual, err := valoresIguales(e, args[1]); err == nil && igual {
				return Valor{Tipo: TipoInt(), I: int64(i)}, true
			}
		}
		return Valor{Tipo: TipoInt(), I: -1}, true

	case "join":
		if len(args) != 2 {
			in.errorSemantico("join espera 2 argumentos", linea, col)
			return Valor{}, false
		}
		v := args[0]
		// tipo de elemento: del SliceVal si existe, si no del tipo declarado
		// (un []string sin inicializar tiene v.Slice == nil).
		tipoElem := TipoNil()
		if v.Tipo.Base == TSlice && v.Tipo.Elemento != nil {
			tipoElem = *v.Tipo.Elemento
		}
		if v.Slice != nil {
			tipoElem = v.Slice.TipoElem
		}
		if v.Tipo.Base != TSlice || tipoElem.Base != TString {
			in.errorSemantico("join solo es válido para []string", linea, col)
			return Valor{}, false
		}
		sep := args[1]
		var partes []string
		if v.Slice != nil {
			partes = make([]string, len(v.Slice.Elems))
			for i, e := range v.Slice.Elems {
				partes[i] = e.S
			}
		}
		return Valor{Tipo: TipoString(), S: strings.Join(partes, sep.S)}, true

	case "Atoi":
		if len(args) != 1 || args[0].Tipo.Base != TString {
			in.errorSemantico("Atoi espera 1 argumento de tipo string", linea, col)
			return Valor{}, false
		}
		n, err := strconv.ParseInt(strings.TrimSpace(args[0].S), 10, 64)
		if err != nil {
			in.errorSemantico(fmt.Sprintf("no se pudo convertir %q a int", args[0].S), linea, col)
			return Valor{}, false
		}
		return Valor{Tipo: TipoInt(), I: n}, true

	case "parseFloat":
		if len(args) != 1 || args[0].Tipo.Base != TString {
			in.errorSemantico("parseFloat espera 1 argumento de tipo string", linea, col)
			return Valor{}, false
		}
		f, err := strconv.ParseFloat(strings.TrimSpace(args[0].S), 64)
		if err != nil {
			in.errorSemantico(fmt.Sprintf("no se pudo convertir %q a float64", args[0].S), linea, col)
			return Valor{}, false
		}
		return Valor{Tipo: TipoFloat(), F: f}, true

	case "typeOf":
		if len(args) != 1 {
			in.errorSemantico("typeOf espera 1 argumento", linea, col)
			return Valor{}, false
		}
		return Valor{Tipo: TipoString(), S: args[0].Tipo.String()}, true

	case "int", "float64", "string", "bool", "rune":
		return in.convertirTipo(nombre, args, linea, col)
	}
	return Valor{}, false
}

// convertirTipo implementa int(x)/float64(x)/string(x)/bool(x)/rune(x):
// solo permite las conversiones explicitas que tienen sentido para cada
// tipo destino (p.ej. bool(x) solo acepta un bool, no hay conversion
// booleana implicita); cualquier otra combinacion es un error semantico.
func (in *Interprete) convertirTipo(destino string, args []Valor, linea, col int) (Valor, bool) {
	if len(args) != 1 {
		in.errorSemantico(destino+"(...) espera 1 argumento", linea, col)
		return Valor{}, false
	}
	v := args[0]
	switch destino {
	case "int":
		switch v.Tipo.Base {
		case TInt:
			return v, true
		case TFloat64:
			return Valor{Tipo: TipoInt(), I: int64(v.F)}, true
		case TRune:
			return Valor{Tipo: TipoInt(), I: int64(v.R)}, true
		}
	case "float64":
		switch v.Tipo.Base {
		case TFloat64:
			return v, true
		case TInt:
			return Valor{Tipo: TipoFloat(), F: float64(v.I)}, true
		}
	case "string":
		switch v.Tipo.Base {
		case TString:
			return v, true
		case TInt, TFloat64, TBool:
			return Valor{Tipo: TipoString(), S: AImprimir(v)}, true
		case TRune:
			return Valor{Tipo: TipoString(), S: string(v.R)}, true
		}
	case "bool":
		if v.Tipo.Base == TBool {
			return v, true
		}
	case "rune":
		switch v.Tipo.Base {
		case TRune:
			return v, true
		case TInt:
			return Valor{Tipo: TipoRune(), R: rune(v.I)}, true
		}
	}
	in.errorSemantico(fmt.Sprintf("no se puede convertir %s a %s", v.Tipo, destino), linea, col)
	return Valor{}, false
}
