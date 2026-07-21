// Command asttest es una herramienta de DESARROLLO (no forma parte del
// pipeline de produccion cli/servidor): parsea un archivo .vch, lo traduce
// al AST propio (internal/traductor) y construye su grafo (internal/ast),
// imprimiendo conteos y los primeros nodos. Sirve para depurar el
// traductor y ConstruirGrafo sin correr el interprete completo.
// Uso: asttest archivo.vch
package main

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"

	"vlangcherry/internal/ast"
	"vlangcherry/internal/parser"
	"vlangcherry/internal/traductor"
)

func main() {
	datos, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	entrada := antlr.NewInputStream(string(datos))
	lexer := parser.NewVLangCherryLexer(entrada)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewVLangCherryParser(tokens)

	arbolParse := p.Programa()
	programa := traductor.Traducir(arbolParse)

	fmt.Printf("Structs: %d, Funciones: %d, Globales: %d\n", len(programa.Structs), len(programa.Funciones), len(programa.Globales))

	grafo := ast.ConstruirGrafo(programa)
	fmt.Printf("Grafo: %d nodos, %d aristas\n", len(grafo.Nodes), len(grafo.Edges))
	for _, n := range grafo.Nodes[:min(10, len(grafo.Nodes))] {
		fmt.Println(" -", n.ID, ":", n.Label)
	}
}
