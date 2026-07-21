package runtime

import (
	"fmt"
	"strconv"
	"strings"
)

// Valor es el valor en tiempo de ejecucion de una expresion VLangCherry.
// Slice y Struct viajan por puntero: asi una copia del Valor (paso de
// parametro, asignacion) comparte la misma identidad subyacente, que es
// exactamente el "por referencia" que pide la seccion 7 para structs y
// slices (los primitivos, en cambio, se copian por valor via los campos
// escalares I/F/S/B/R).
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

// SliceVal es la representacion de un []T o [][]T. Append (5.1.5) genera
// SIEMPRE un SliceVal nuevo, igual que en Go real: la variable debe
// reasignarse con el resultado.
type SliceVal struct {
	Elems    []Valor
	TipoElem Tipo
}

// StructVal es una instancia de struct: nombre de tipo + campos con su
// orden de declaracion (para el formato de impresion 7.2.1.1).
type StructVal struct {
	NombreTipo  string
	OrdenCampos []string
	Campos      map[string]*Valor
}

func NuevoSlice(tipoElem Tipo, elems []Valor) Valor {
	return Valor{Tipo: TipoSliceDe(tipoElem), Slice: &SliceVal{Elems: elems, TipoElem: tipoElem}}
}

func NuevaInstanciaStruct(nombre string, orden []string, campos map[string]*Valor) Valor {
	return Valor{Tipo: TipoStructDe(nombre), Struct: &StructVal{NombreTipo: nombre, OrdenCampos: orden, Campos: campos}}
}

func EsNil(v Valor) bool {
	return v.Tipo.Base == TNil || ((v.Tipo.Base == TSlice) && v.Slice == nil) || (v.Tipo.Base == TStruct && v.Struct == nil)
}

// AImprimir formatea un valor para consola/print (7.2.1, 7.2.1.1).
func AImprimir(v Valor) string {
	switch v.Tipo.Base {
	case TInt:
		return strconv.FormatInt(v.I, 10)
	case TFloat64:
		return strconv.FormatFloat(v.F, 'f', -1, 64)
	case TString:
		return v.S
	case TBool:
		return strconv.FormatBool(v.B)
	case TRune:
		return string(v.R)
	case TNil:
		return "nil"
	case TSlice:
		if v.Slice == nil {
			return "nil"
		}
		partes := make([]string, len(v.Slice.Elems))
		for i, e := range v.Slice.Elems {
			partes[i] = AImprimir(e)
		}
		return "[" + strings.Join(partes, " ") + "]"
	case TStruct:
		if v.Struct == nil {
			return "nil"
		}
		partes := make([]string, len(v.Struct.OrdenCampos))
		for i, campo := range v.Struct.OrdenCampos {
			partes[i] = fmt.Sprintf("%s: %s", campo, AImprimir(*v.Struct.Campos[campo]))
		}
		return v.Struct.NombreTipo + "{" + strings.Join(partes, ", ") + "}"
	}
	return ""
}

// AValorReporte formatea un valor para la tabla de simbolos (8.2): mismo
// formato que consola, es suficientemente legible para el reporte.
func AValorReporte(v Valor) string {
	return AImprimir(v)
}
