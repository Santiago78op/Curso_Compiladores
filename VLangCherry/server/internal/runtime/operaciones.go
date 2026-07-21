package runtime

import (
	"fmt"
	"strings"
)

// Operaciones aritmeticas/comparacion/logicas segun las tablas del
// enunciado (4.3-4.5). Decision de diseno documentada: las tablas del PDF
// llegaron con celdas rotas por la conversion en algunos casos (resta,
// multiplicacion, modulo); se completaron por simetria con el patron
// consistente del resto de la seccion (int/float64 se mezclan siempre
// hacia float64) y se restringio modulo a int%int, tal como dice el texto
// descriptivo ("el modulo produce el residuo... entre tipos int").

func Suma(a, b Valor, linea, col int) (Valor, error) {
	switch {
	case a.Tipo.Base == TString && b.Tipo.Base == TString:
		return Valor{Tipo: TipoString(), S: a.S + b.S}, nil
	case a.Tipo.Base == TInt && b.Tipo.Base == TInt:
		return Valor{Tipo: TipoInt(), I: a.I + b.I}, nil
	case EsNumerico(a.Tipo) && EsNumerico(b.Tipo):
		return Valor{Tipo: TipoFloat(), F: aFloat(a) + aFloat(b)}, nil
	}
	return Valor{}, errorOperacion("+", a, b)
}

func Resta(a, b Valor, linea, col int) (Valor, error) {
	switch {
	case a.Tipo.Base == TInt && b.Tipo.Base == TInt:
		return Valor{Tipo: TipoInt(), I: a.I - b.I}, nil
	case EsNumerico(a.Tipo) && EsNumerico(b.Tipo):
		return Valor{Tipo: TipoFloat(), F: aFloat(a) - aFloat(b)}, nil
	}
	return Valor{}, errorOperacion("-", a, b)
}

func Multiplicacion(a, b Valor, linea, col int) (Valor, error) {
	switch {
	case a.Tipo.Base == TInt && b.Tipo.Base == TInt:
		return Valor{Tipo: TipoInt(), I: a.I * b.I}, nil
	case EsNumerico(a.Tipo) && EsNumerico(b.Tipo):
		return Valor{Tipo: TipoFloat(), F: aFloat(a) * aFloat(b)}, nil
	}
	return Valor{}, errorOperacion("*", a, b)
}

func Division(a, b Valor, linea, col int) (Valor, error) {
	switch {
	case a.Tipo.Base == TInt && b.Tipo.Base == TInt:
		if b.I == 0 {
			return Valor{}, fmt.Errorf("división entre cero")
		}
		return Valor{Tipo: TipoInt(), I: a.I / b.I}, nil
	case EsNumerico(a.Tipo) && EsNumerico(b.Tipo):
		if aFloat(b) == 0 {
			return Valor{}, fmt.Errorf("división entre cero")
		}
		return Valor{Tipo: TipoFloat(), F: aFloat(a) / aFloat(b)}, nil
	}
	return Valor{}, errorOperacion("/", a, b)
}

func Modulo(a, b Valor, linea, col int) (Valor, error) {
	if a.Tipo.Base == TInt && b.Tipo.Base == TInt {
		if b.I == 0 {
			return Valor{}, fmt.Errorf("división entre cero")
		}
		return Valor{Tipo: TipoInt(), I: a.I % b.I}, nil
	}
	return Valor{}, errorOperacion("%", a, b)
}

func NegacionUnaria(a Valor, linea, col int) (Valor, error) {
	switch a.Tipo.Base {
	case TInt:
		return Valor{Tipo: TipoInt(), I: -a.I}, nil
	case TFloat64:
		return Valor{Tipo: TipoFloat(), F: -a.F}, nil
	}
	return Valor{}, fmt.Errorf("no se puede negar un valor de tipo %s", a.Tipo)
}

func NotLogico(a Valor) (Valor, error) {
	if a.Tipo.Base != TBool {
		return Valor{}, fmt.Errorf("el operador ! requiere BOOL, se recibió %s", a.Tipo)
	}
	return Valor{Tipo: TipoBool(), B: !a.B}, nil
}

func Y(a, b Valor) (Valor, error) {
	if a.Tipo.Base != TBool || b.Tipo.Base != TBool {
		return Valor{}, fmt.Errorf("&& requiere BOOL en ambos lados, se recibió %s y %s", a.Tipo, b.Tipo)
	}
	return Valor{Tipo: TipoBool(), B: a.B && b.B}, nil
}

func O(a, b Valor) (Valor, error) {
	if a.Tipo.Base != TBool || b.Tipo.Base != TBool {
		return Valor{}, fmt.Errorf("|| requiere BOOL en ambos lados, se recibió %s y %s", a.Tipo, b.Tipo)
	}
	return Valor{Tipo: TipoBool(), B: a.B || b.B}, nil
}

func Igualdad(a, b Valor, negar bool) (Valor, error) {
	igual, err := valoresIguales(a, b)
	if err != nil {
		return Valor{}, err
	}
	if negar {
		igual = !igual
	}
	return Valor{Tipo: TipoBool(), B: igual}, nil
}

func valoresIguales(a, b Valor) (bool, error) {
	switch {
	case EsNumerico(a.Tipo) && EsNumerico(b.Tipo):
		return aFloat(a) == aFloat(b), nil
	case a.Tipo.Base == TString && b.Tipo.Base == TString:
		return a.S == b.S, nil
	case a.Tipo.Base == TBool && b.Tipo.Base == TBool:
		return a.B == b.B, nil
	case a.Tipo.Base == TRune && b.Tipo.Base == TRune:
		return a.R == b.R, nil
	case a.Tipo.Base == TNil || b.Tipo.Base == TNil:
		return EsNil(a) && EsNil(b), nil
	}
	return false, fmt.Errorf("no se puede comparar (==) %s y %s", a.Tipo, b.Tipo)
}

// Relacional evalua <, <=, >, >= entre operandos numericos (con la mezcla
// int/float64 habitual), entre dos runes comparando su valor ASCII (4.4.2,
// consideracion final: "la comparacion de valores tipo rune se realiza
// comparando su valor ASCII"), o entre dos strings de forma lexicografica
// caracter por caracter (4.4.1, consideracion: "las comparaciones entre
// cadenas se hacen lexicograficamente" - esa nota aparece bajo el titulo de
// igualdad pero "lexicografico" es un criterio de ORDEN, asi que se aplica
// aqui, a los relacionales, y no solo a ==/!=).
func Relacional(op string, a, b Valor) (Valor, error) {
	var cmp int
	switch {
	case EsNumerico(a.Tipo) && EsNumerico(b.Tipo):
		fa, fb := aFloat(a), aFloat(b)
		cmp = compararFloats(fa, fb)
	case a.Tipo.Base == TRune && b.Tipo.Base == TRune:
		cmp = int(a.R) - int(b.R)
	case a.Tipo.Base == TString && b.Tipo.Base == TString:
		cmp = strings.Compare(a.S, b.S)
	default:
		return Valor{}, fmt.Errorf("no se puede comparar (%s) entre %s y %s", op, a.Tipo, b.Tipo)
	}
	var r bool
	switch op {
	case "<":
		r = cmp < 0
	case "<=":
		r = cmp <= 0
	case ">":
		r = cmp > 0
	case ">=":
		r = cmp >= 0
	}
	return Valor{Tipo: TipoBool(), B: r}, nil
}

func compararFloats(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func aFloat(v Valor) float64 {
	if v.Tipo.Base == TInt {
		return float64(v.I)
	}
	return v.F
}

func errorOperacion(op string, a, b Valor) error {
	return fmt.Errorf("no se puede realizar la operación %s entre %s y %s", op, a.Tipo, b.Tipo)
}
