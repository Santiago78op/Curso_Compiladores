// Command parsetest es una herramienta de DESARROLLO (no forma parte del
// pipeline de produccion cli/servidor): corre solo el lexer+parser de ANTLR
// sobre un archivo .vch e imprime el parse tree crudo, para depurar la
// gramatica (grammar/VLangCherry.g4) de forma aislada del traductor y del
// interprete. Uso: parsetest archivo.vch
package main

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"

	"vlangcherry/internal/parser"
)

type erroresListener struct {
	*antlr.DefaultErrorListener
	errores []string
}

func (l *erroresListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.errores = append(l.errores, fmt.Sprintf("linea %d columna %d: %s", line, column+1, msg))
}

func main() {
	datos, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	entrada := antlr.NewInputStream(string(datos))
	lexer := parser.NewVLangCherryLexer(entrada)
	listener := &erroresListener{}
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(listener)

	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewVLangCherryParser(tokens)
	p.RemoveErrorListeners()
	p.AddErrorListener(listener)

	arbol := p.Programa()

	if len(listener.errores) > 0 {
		fmt.Println("ERRORES:")
		for _, e := range listener.errores {
			fmt.Println(" -", e)
		}
		os.Exit(1)
	}

	fmt.Println("OK, parseo exitoso.")
	fmt.Println(arbol.ToStringTree(nil, p))
}
