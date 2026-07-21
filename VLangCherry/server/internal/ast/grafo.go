package ast

import (
	"fmt"
	"reflect"
	"strings"
)

// NodoGrafo/AristaGrafo/Grafo son el formato de salida para el reporte
// AST (8.3), pensado para viajar como JSON y graficarse con vis-network
// en el cliente (mismo formato que usa CompInterpreter).
type NodoGrafo struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type AristaGrafo struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Label string `json:"label"`
}

type Grafo struct {
	Nodes []NodoGrafo   `json:"nodes"`
	Edges []AristaGrafo `json:"edges"`
}

// ConstruirGrafo recorre CUALQUIER Nodo por reflection: los campos que son
// Nodo (directos, en slices, o anidados dentro de structs contenedores como
// RamaIf/CasoSwitch) se vuelven aristas; los campos escalares se pegan a la
// etiqueta del nodo. Un tipo de nodo nuevo aparece en el grafo sin tocar
// esta funcion (misma idea que reportes/ast-grafo.js en CompInterpreter).
func ConstruirGrafo(raiz Nodo) Grafo {
	g := &Grafo{}
	contador := 0

	var visitarNodo func(n Nodo) string
	var procesarCampo func(nombreCampo string, fv reflect.Value, padreID string, etiqueta *string)

	// esNodoValor intenta ver un reflect.Value como Nodo (puntero no nulo
	// que implementa la interfaz); lo usan procesarCampo (para decidir si
	// un campo es arista) y visitarNodo (para castear el hijo antes de
	// recursar).
	esNodoValor := func(v reflect.Value) (Nodo, bool) {
		if !v.IsValid() {
			return nil, false
		}
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return nil, false
		}
		if !v.CanInterface() {
			return nil, false
		}
		iv := v.Interface()
		if iv == nil {
			return nil, false
		}
		if n, ok := iv.(Nodo); ok {
			return n, true
		}
		return nil, false
	}

	// visitarNodo registra un nodo nuevo en g.Nodes (ID autoincremental) y
	// recorre sus campos exportados via procesarCampo para descubrir
	// aristas hacia sus hijos; retorna el ID asignado para que el llamador
	// pueda conectar la arista padre->hijo.
	visitarNodo = func(n Nodo) string {
		miID := fmt.Sprintf("n%d", contador)
		contador++
		etiqueta := n.Tipo()

		rv := reflect.ValueOf(n)
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				g.Nodes = append(g.Nodes, NodoGrafo{ID: miID, Label: etiqueta})
				return miID
			}
			rv = rv.Elem()
		}

		if rv.Kind() == reflect.Struct {
			t := rv.Type()
			for i := 0; i < rv.NumField(); i++ {
				campo := t.Field(i)
				if campo.Name == "base" || !campo.IsExported() {
					continue
				}
				procesarCampo(campo.Name, rv.Field(i), miID, &etiqueta)
			}
		}

		g.Nodes = append(g.Nodes, NodoGrafo{ID: miID, Label: etiqueta})
		return miID
	}

	// procesarCampo clasifica un campo por su Kind: si contiene (o
	// contiene una coleccion de) Nodo, genera aristas hacia esos hijos;
	// si es un TipoAST o un escalar, lo agrega como texto a la etiqueta
	// del nodo padre; si es un struct "contenedor transparente" (RamaIf,
	// CasoSwitch, Parametro, ...), recursa sobre sus propios campos
	// propagando el nombre compuesto (p.ej. "Ramas[0].Condicion").
	procesarCampo = func(nombreCampo string, fv reflect.Value, padreID string, etiqueta *string) {
		if !fv.IsValid() {
			return
		}

		switch fv.Kind() {
		case reflect.Ptr, reflect.Interface:
			if fv.IsNil() {
				return
			}
			if nodo, ok := esNodoValor(fv); ok {
				hijoID := visitarNodo(nodo)
				g.Edges = append(g.Edges, AristaGrafo{From: padreID, To: hijoID, Label: nombreCampo})
				return
			}
			procesarCampo(nombreCampo, fv.Elem(), padreID, etiqueta)

		case reflect.Slice, reflect.Array:
			if fv.Len() == 0 {
				return
			}
			for i := 0; i < fv.Len(); i++ {
				procesarCampo(fmt.Sprintf("%s[%d]", nombreCampo, i), fv.Index(i), padreID, etiqueta)
			}

		case reflect.Struct:
			if nodo, ok := esNodoValor(fv); ok {
				hijoID := visitarNodo(nodo)
				g.Edges = append(g.Edges, AristaGrafo{From: padreID, To: hijoID, Label: nombreCampo})
				return
			}
			if ta, ok := fv.Interface().(TipoAST); ok {
				*etiqueta += "\n" + nombreCampo + ": " + ta.String()
				return
			}
			// contenedor transparente (RamaIf, CasoSwitch, Parametro,
			// CampoStruct, CampoValorLiteral, ...): recorrer sus campos
			// propagando el nombre compuesto.
			t := fv.Type()
			for i := 0; i < fv.NumField(); i++ {
				campo := t.Field(i)
				if !campo.IsExported() {
					continue
				}
				procesarCampo(nombreCampo+"."+campo.Name, fv.Field(i), padreID, etiqueta)
			}

		default:
			if nodo, ok := esNodoValor(fv); ok {
				hijoID := visitarNodo(nodo)
				g.Edges = append(g.Edges, AristaGrafo{From: padreID, To: hijoID, Label: nombreCampo})
				return
			}
			valorTexto := fmt.Sprintf("%v", fv.Interface())
			if valorTexto == "" {
				return
			}
			*etiqueta += "\n" + nombreCampo + ": " + valorTexto
		}
	}

	if raiz != nil {
		visitarNodo(raiz)
	}
	return *g
}

func (t TipoAST) String() string {
	if t.EsSlice {
		if t.Elemento == nil {
			return "[]?"
		}
		return "[]" + t.Elemento.String()
	}
	return t.NombreBase
}

// ADot genera el equivalente en Graphviz DOT, por si se quiere graficar
// del lado servidor (sugerencia del enunciado, seccion 8.3).
func ADot(g Grafo) string {
	var sb strings.Builder
	sb.WriteString("digraph AST {\n  node [shape=box, style=rounded, fontname=\"Consolas\"];\n")
	for _, n := range g.Nodes {
		lbl := strings.ReplaceAll(n.Label, `"`, `\"`)
		lbl = strings.ReplaceAll(lbl, "\n", `\n`)
		sb.WriteString(fmt.Sprintf("  %s [label=\"%s\"];\n", n.ID, lbl))
	}
	for _, e := range g.Edges {
		sb.WriteString(fmt.Sprintf("  %s -> %s;\n", e.From, e.To))
	}
	sb.WriteString("}\n")
	return sb.String()
}
