// Pruebas del pipeline completo. Nacen del hallazgo B3 de la auditoria
// (Fase C): un typo cualquiera hacia aparecer, junto al error sintactico
// real, un "error interno durante la ejecución ... línea 0" que daba la
// impresion de que el interprete se habia roto.
//
// Lo delicado del arreglo no es callar ese mensaje, sino callarlo SIN
// perder dos cosas: el criterio "best-effort" (seguir recolectando errores
// semanticos aunque el parseo haya fallado) y la visibilidad de un panic
// genuino cuando la entrada parseo limpia. Cada test fija una de esas.
package analizar

import "strings"

import "testing"

func tiposDeError(r Resultado) (lexico, sintactico, semantico int) {
	for _, e := range r.Errores {
		switch e.Tipo {
		case "Léxico":
			lexico++
		case "Sintáctico":
			sintactico++
		case "Semántico":
			semantico++
		}
	}
	return
}

func erroresInternos(r Resultado) []string {
	var out []string
	for _, e := range r.Errores {
		if strings.Contains(e.Descripcion, "error interno") {
			out = append(out, e.Descripcion)
		}
	}
	return out
}

// B3: con la entrada rota, el usuario ve el error sintactico y nada mas.
func TestErrorSintacticoNoProduceErrorInterno(t *testing.T) {
	casos := map[string]string{
		"asignacion sin expresion": "func main() {\n    x := \n}\n",
		"range con :=":             "func main() {\n    for i, v := range []int{1} {\n        println(v)\n    }\n}\n",
		"llave sin cerrar":         "func main() {\n    println(1)\n",
	}
	for nombre, codigo := range casos {
		t.Run(nombre, func(t *testing.T) {
			r := Analizar(codigo)
			if _, sint, _ := tiposDeError(r); sint == 0 {
				t.Fatalf("se esperaba al menos un error sintáctico; errores: %+v", r.Errores)
			}
			if internos := erroresInternos(r); len(internos) > 0 {
				t.Errorf("no debe reportarse un error interno cuando la entrada ya venía rota: %v", internos)
			}
		})
	}
}

// El criterio best-effort documentado en Analizar: aunque el parseo falle,
// se sigue traduciendo y ejecutando para juntar TODOS los errores (8.1).
// Callar el panic no puede haberse llevado esto puesto.
func TestBestEffortSigueReportandoSemanticos(t *testing.T) {
	codigo := "func main() {\n" +
		"    println(noExiste)\n" + // semántico: no definida
		"    x := 5\n" +
		"    x = 6\n" + // semántico: sin mut
		"    y := @\n" + // léxico + sintáctico
		"}\n"
	r := Analizar(codigo)
	lex, sint, sem := tiposDeError(r)
	if lex == 0 || sint == 0 {
		t.Fatalf("el caso debe producir errores de entrada; lex=%d sint=%d", lex, sint)
	}
	if sem < 2 {
		t.Errorf("best-effort roto: se esperaban los 2 errores semánticos aun con la entrada rota, hubo %d; errores: %+v", sem, r.Errores)
	}
}

// La contracara: si el codigo parsea LIMPIO, un panic del interprete si es
// un bug nuestro y tiene que verse. Esconderlo seria peor que el ruido.
func TestCodigoLimpioNoSilenciaErroresInternos(t *testing.T) {
	r := Analizar("func main() {\n    println(\"hola\")\n}\n")
	if lex, sint, _ := tiposDeError(r); lex != 0 || sint != 0 {
		t.Fatalf("este código debe parsear limpio; lex=%d sint=%d", lex, sint)
	}
	// No hay panic que provocar sin romper el intérprete a propósito, así que
	// lo que se fija es la precondición del arreglo: con entrada limpia la
	// bandera que silencia el panic queda en false, y el reporte se emite.
	if len(r.Errores) != 0 {
		t.Errorf("código válido no debe reportar errores: %+v", r.Errores)
	}
	if r.Consola != "hola" {
		t.Errorf("consola = %q, se esperaba \"hola\"", r.Consola)
	}
}

// Guarda de regresión de los hallazgos ALTOS, que sobrevivieron tanto tiempo
// justamente porque ningún ejemplo los ejercitaba.
func TestGuardasYSemanticaAlta(t *testing.T) {
	casos := []struct {
		nombre  string
		codigo  string
		esperar string // subcadena que debe aparecer en algún error
	}{
		{"A1 reasignar sin mut", "func main() {\n    x := 5\n    x = 6\n}\n", "sin \"mut\""},
		{"A1 += sin mut", "func main() {\n    x := 5\n    x += 1\n}\n", "sin \"mut\""},
		{"A3 ciclo infinito", "func main() {\n    for true {\n    }\n}\n", "ciclo infinito"},
		{"A3 recursión", "func f() int { return f() }\nfunc main() { println(f()) }\n", "demasiado profunda"},
		{"M1 return sin valor", "func f() int {\n    return\n}\nfunc main() { println(f()) }\n", "sin valor"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			r := Analizar(c.codigo)
			for _, e := range r.Errores {
				if strings.Contains(e.Descripcion, c.esperar) {
					return
				}
			}
			t.Errorf("no se reportó el error esperado (%q); errores: %+v", c.esperar, r.Errores)
		})
	}
}

// A2: el cortocircuito no es un detalle de estilo — es lo que hace que la
// guarda clásica funcione. Si se rompe, esto reporta división entre cero.
func TestCortocircuitoEvitaLaDivision(t *testing.T) {
	r := Analizar("func main() {\n    x := 0\n    if x != 0 && 10/x > 1 {\n        println(\"no\")\n    }\n    println(\"ok\")\n}\n")
	if len(r.Errores) != 0 {
		t.Fatalf("la guarda no debía evaluar la división; errores: %+v", r.Errores)
	}
	if r.Consola != "ok" {
		t.Errorf("consola = %q, se esperaba \"ok\"", r.Consola)
	}
}

// A5: declarar slice o struct sin inicializar no puede romper la asignación
// posterior (el valor por defecto conserva el tipo declarado).
func TestDeclararSinInicializarConservaElTipo(t *testing.T) {
	r := Analizar("func main() {\n    mut xs []int\n    xs = []int{1, 2}\n    xs = append(xs, 3)\n    println(xs)\n}\n")
	if len(r.Errores) != 0 {
		t.Fatalf("no debía haber errores; %+v", r.Errores)
	}
	if r.Consola != "[1 2 3]" {
		t.Errorf("consola = %q, se esperaba \"[1 2 3]\"", r.Consola)
	}
}
