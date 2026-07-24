package runtime

// TipoBase enumera las categorias de tipo de VLangCherry (seccion 3.7-3.8).
type TipoBase int

const (
	TInt TipoBase = iota
	TFloat64
	TString
	TBool
	TRune
	TSlice
	TStruct
	TNil  // valor nulo (3.9), compatible con slice/struct/nil
	TVoid // funcion sin retorno
)

// Tipo describe un tipo VLangCherry completo: primitivo, slice (con su
// elemento, recursivo para multidimensional) o struct (por nombre).
type Tipo struct {
	Base         TipoBase
	Elemento     *Tipo
	NombreStruct string
}

func (t Tipo) String() string {
	switch t.Base {
	case TInt:
		return "int"
	case TFloat64:
		return "float64"
	case TString:
		return "string"
	case TBool:
		return "bool"
	case TRune:
		return "rune"
	case TSlice:
		return "[]" + t.Elemento.String()
	case TStruct:
		return t.NombreStruct
	case TNil:
		return "nil"
	case TVoid:
		return "void"
	}
	return "desconocido"
}

// Igual compara dos tipos por equivalencia estructural (3.6: una variable
// solo admite valores compatibles con su tipo declarado).
func (t Tipo) Igual(o Tipo) bool {
	if t.Base != o.Base {
		return false
	}
	switch t.Base {
	case TSlice:
		if t.Elemento == nil || o.Elemento == nil {
			return t.Elemento == o.Elemento
		}
		return t.Elemento.Igual(*o.Elemento)
	case TStruct:
		return t.NombreStruct == o.NombreStruct
	}
	return true
}

func EsNumerico(t Tipo) bool {
	return t.Base == TInt || t.Base == TFloat64
}

func TipoInt() Tipo    { return Tipo{Base: TInt} }
func TipoFloat() Tipo  { return Tipo{Base: TFloat64} }
func TipoString() Tipo { return Tipo{Base: TString} }
func TipoBool() Tipo   { return Tipo{Base: TBool} }
func TipoRune() Tipo   { return Tipo{Base: TRune} }
func TipoNil() Tipo    { return Tipo{Base: TNil} }
func TipoVoid() Tipo   { return Tipo{Base: TVoid} }
func TipoSliceDe(el Tipo) Tipo {
	e := el
	return Tipo{Base: TSlice, Elemento: &e}
}
func TipoStructDe(nombre string) Tipo {
	return Tipo{Base: TStruct, NombreStruct: nombre}
}

// ValorPorDefecto da el valor inicial de una variable declarada sin
// expresion (3.7: "tomara el valor por defecto del tipo correspondiente").
func ValorPorDefecto(t Tipo) Valor {
	switch t.Base {
	case TInt:
		return Valor{Tipo: t, I: 0}
	case TFloat64:
		return Valor{Tipo: t, F: 0}
	case TString:
		return Valor{Tipo: t, S: ""}
	case TBool:
		return Valor{Tipo: t, B: false}
	case TRune:
		return Valor{Tipo: t, R: 0}
	default:
		// slice, struct y cualquier otro compuesto arrancan "nil" PERO
		// conservando el tipo declarado (3.7 nota final): asi una asignacion
		// posterior (p.ej. xs = []int{...}) valida contra []int y no contra un
		// nil sin tipo. EsNil() ya trata {TSlice,Slice:nil}/{TStruct,Struct:nil}.
		return Valor{Tipo: t}
	}
}
