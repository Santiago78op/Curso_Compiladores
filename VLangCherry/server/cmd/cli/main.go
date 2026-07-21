// Command cli corre el pipeline completo de VLangCherry (internal/analizar)
// sobre un archivo .vch y vuelca el resultado a stdout en texto plano; con
// el flag --json al final imprime ademas el Resultado completo como JSON
// (mismo formato que responde cmd/servidor). Uso: cli archivo.vch [--json]
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"vlangcherry/internal/analizar"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "uso: cli <archivo.vch> [--json]")
		os.Exit(1)
	}
	datos, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "no se pudo leer %q: %v\n", os.Args[1], err)
		os.Exit(1)
	}
	r := analizar.Analizar(string(datos))

	fmt.Println("=== ERRORES ===")
	for _, e := range r.Errores {
		fmt.Printf("[%s] línea %d col %d: %s\n", e.Tipo, e.Linea, e.Columna, e.Descripcion)
	}
	fmt.Println("=== CONSOLA ===")
	fmt.Println(r.Consola)
	fmt.Println("=== SIMBOLOS ===")
	for _, s := range r.Simbolos {
		fmt.Printf("%s | %s | %s | %s | %s | %d:%d\n", s.ID, s.TipoSimbolo, s.TipoDato, s.Ambito, s.Valor, s.Linea, s.Columna)
	}
	fmt.Println("=== AST ===")
	fmt.Printf("%d nodos, %d aristas\n", len(r.AST.Nodes), len(r.AST.Edges))

	if len(os.Args) > 2 && os.Args[2] == "--json" {
		b, _ := json.MarshalIndent(r, "", "  ")
		fmt.Println(string(b))
	}
}
