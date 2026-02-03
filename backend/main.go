package main

import (
	"log"
	"net/http"

	"digi-leatherman/backend/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/string/url-encode", cors(handlers.URLEncode))
	mux.HandleFunc("/api/string/url-decode", cors(handlers.URLDecode))

	addr := ":8100"
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// cors wraps a handler to add CORS headers for the Vite dev server.
func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5273")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}
