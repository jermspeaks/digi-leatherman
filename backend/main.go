package main

import (
	"log"
	"net/http"
	"strings"

	"digi-leatherman/backend/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/string/url-encode", cors(handlers.URLEncode))
	mux.HandleFunc("/api/string/url-decode", cors(handlers.URLDecode))
	mux.HandleFunc("/api/string/base64-encode", cors(handlers.Base64Encode))
	mux.HandleFunc("/api/string/base64-decode", cors(handlers.Base64Decode))
	mux.HandleFunc("/api/string/trim", cors(handlers.Trim))
	mux.HandleFunc("/api/string/upper-case", cors(handlers.UpperCase))
	mux.HandleFunc("/api/string/lower-case", cors(handlers.LowerCase))
	mux.HandleFunc("/api/string/capital-case", cors(handlers.CapitalCase))
	mux.HandleFunc("/api/string/snake-case", cors(handlers.SnakeCase))
	mux.HandleFunc("/api/string/kebab-case", cors(handlers.KebabCase))
	mux.HandleFunc("/api/string/camel-case", cors(handlers.CamelCase))
	mux.HandleFunc("/api/string/pascal-case", cors(handlers.PascalCase))
	mux.HandleFunc("/api/string/sentence-case", cors(handlers.SentenceCase))

	addr := ":8100"
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// cors wraps a handler to add CORS headers for the frontend dev server.
// In development, the request Origin is reflected so any localhost port works.
func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// Allow any localhost origin in development (e.g. 5173, 5273, 3000).
		if origin != "" && (strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:")) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}
