// Servidor Express-like (net/http) para VLangCherry.
// Endpoint principal: POST /interpretar { "codigo": "<fuente .vch>" }
// Responde JSON: { errores, consola, consolaLineas, simbolos, ast, dot }.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"vlangcherry/internal/analizar"
)

type peticion struct {
	Codigo string `json:"codigo"`
}

// codificarJSON escribe v como JSON en la respuesta; si la codificacion
// falla (p.ej. el cliente corto la conexion a mitad de escritura) se
// registra en el log en vez de ignorarse en silencio, aunque ya no haya
// forma de corregir la respuesta parcial que el cliente pudo haber recibido.
func codificarJSON(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("error codificando respuesta JSON:", err)
	}
}

func conCORS(siguiente http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		siguiente(w, r)
	}
}

func main() {
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "4100"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", conCORS(func(w http.ResponseWriter, r *http.Request) {
		codificarJSON(w, map[string]string{
			"nombre": "VLangCherry", "estado": "ok", "endpoint": "POST /interpretar { codigo }",
		})
	}))

	mux.HandleFunc("/salud", conCORS(func(w http.ResponseWriter, r *http.Request) {
		codificarJSON(w, map[string]string{"estado": "ok"})
	}))

	mux.HandleFunc("/interpretar", conCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var p peticion
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			codificarJSON(w, map[string]string{"error": "cuerpo inválido"})
			return
		}
		resultado := analizar.Analizar(p.Codigo)
		w.Header().Set("Content-Type", "application/json")
		codificarJSON(w, resultado)
	}))

	log.Println("VLangCherry server escuchando en http://localhost:" + puerto)
	log.Fatal(http.ListenAndServe("localhost:"+puerto, mux))
}
