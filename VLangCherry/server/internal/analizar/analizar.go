// Package analizar orquesta el pipeline completo: fuente .vch ->
// { errores, consola, simbolos, ast }. Entorno FRESCO por llamada.
package analizar

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"vlangcherry/internal/ast"
	"vlangcherry/internal/parser"
	"vlangcherry/internal/runtime"
	"vlangcherry/internal/traductor"
)

type Resultado struct {
	Errores       []runtime.ErrorReporte `json:"errores"`
	Consola       string                 `json:"consola"`
	ConsolaLineas []string               `json:"consolaLineas"`
	Simbolos      []runtime.FilaSimbolo  `json:"simbolos"`
	AST           ast.Grafo              `json:"ast"`
	Dot           string                 `json:"dot"`
}

type oyenteErrores struct {
	*antlr.DefaultErrorListener
	errores *runtime.ListaErrores
	esLexer bool
}

func (o *oyenteErrores) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	if o.esLexer {
		o.errores.Lexico(msg, line, column+1)
	} else {
		o.errores.Sintactico(msg, line, column+1)
	}
}

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

	// Best-effort: aunque haya errores lexicos/sintacticos ya reportados,
	// se intenta traducir y ejecutar lo que se pudo parsear (igual criterio
	// que los demas proyectos del curso: recolectar TODOS los errores).
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

	consolaLineas := interprete.Consola()
	consolaTexto := ""
	for i, l := range consolaLineas {
		if i > 0 {
			consolaTexto += "\n"
		}
		consolaTexto += l
	}

	listaErrores := errores.Errores
	return Resultado{
		Errores:       listaErrores,
		Consola:       consolaTexto,
		ConsolaLineas: consolaLineas,
		Simbolos:      interprete.Simbolos(),
		AST:           grafo,
		Dot:           ast.ADot(grafo),
	}
}
