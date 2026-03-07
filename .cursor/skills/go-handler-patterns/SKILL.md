---
name: go-handler-patterns
description: Go HTTP handler patterns for Digital Leatherman backend. Use when creating API endpoints, handlers, or backend logic. Triggers on "go handler", "backend handler", "api endpoint", "http handler".
---

# Go Handler Patterns

This project uses Go's standard library `net/http` with consistent patterns.

## Request/Response Types

Defined in `backend/handlers/string.go`:

```go
type StringRequest struct {
	Value string `json:"value"`
}

type StringResponse struct {
	Result string `json:"result"`
}
```

For custom endpoints, define new types following this pattern.

## Helper Functions

### decodeBody

Validates POST method and decodes JSON body:

```go
func decodeBody(w http.ResponseWriter, r *http.Request) (StringRequest, bool) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return StringRequest{}, false
	}
	var req StringRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return StringRequest{}, false
	}
	return req, true
}
```

### writeJSON

Writes JSON response with proper content-type:

```go
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
```

## Handler Pattern

Standard handler structure:

```go
func MyHandler(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	
	// Process the request
	result := processValue(req.Value)
	
	// Return response
	writeJSON(w, StringResponse{Result: result})
}
```

For custom request types, decode manually:

```go
func CustomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CustomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Process and respond
	result := process(req)
	writeJSON(w, CustomResponse{Result: result})
}
```

## CORS Wrapper

All handlers are wrapped with CORS in `backend/main.go`:

```go
func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && (strings.HasPrefix(origin, "http://localhost:") || 
		                    strings.HasPrefix(origin, "http://127.0.0.1:")) {
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
```

## Route Registration

In `backend/main.go`:

```go
mux.HandleFunc("/api/string/my-handler", cors(handlers.MyHandler))
```

## Error Handling

Use `http.Error` with appropriate status codes:

```go
http.Error(w, "method not allowed", http.StatusMethodNotAllowed)  // 405
http.Error(w, "invalid JSON", http.StatusBadRequest)              // 400
http.Error(w, "invalid value", http.StatusBadRequest)             // 400
http.Error(w, "internal error", http.StatusInternalServerError)   // 500
```

## Pure Functions

Extract pure logic into testable functions:

```go
// Pure function - easy to test
func toSnake(s string) string {
	words := wordsFrom(s)
	return strings.Join(words, "_")
}

// Handler uses pure function
func SnakeCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: toSnake(req.Value)})
}
```

## Middleware

Middleware in `backend/middleware/`:

- `logging.go` - Request/response logging with slog
- `recovery.go` - Panic recovery

Applied in `main.go`:

```go
handler := middleware.Recovery(middleware.Logging(mux))
```
