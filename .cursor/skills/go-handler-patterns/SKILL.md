---
name: go-handler-patterns
description: Go HTTP handler patterns for Digital Leatherman backend. Use when creating API endpoints, handlers, or backend logic. Triggers on "go handler", "backend handler", "api endpoint", "http handler".
---

# Go Handler Patterns

This project uses Go's standard library `net/http` with consistent patterns.

## Handler Inventory

Handlers live in `backend/handlers/`. Routes follow `/api/{category}/{action}` in `backend/main.go`.

| File | Handler functions |
|------|-------------------|
| `string.go` | URLEncode, URLDecode, ParseURLParams, CreateURLWithParams, Base64Encode, Base64Decode, Trim, UpperCase, LowerCase, CapitalCase, SnakeCase, KebabCase, CamelCase, PascalCase, SentenceCase, SpellOut |
| `json.go` | FormatJSON, MinifyJSON, ValidateJSON, PathQueryJSON, DiffJSON |
| `uuid.go` | GenerateUUIDv4, GenerateUUIDv7, BulkGenerateUUID, ValidateUUID, ParseUUID, FormatUUID |
| `lorem.go` | LoremIpsum (single endpoint; dispatches by `tool` in request) |
| `date.go` | DateParse, DateTimezone, DateCountdown |
| `color.go` | ColorConvert, ContrastCheck, ColorBlindness, ExtractPalette, GenerateShades, GenerateHarmonies |

When adding a tool, add to the appropriate file or create a new one for a new category; match the naming style (PascalCase handler names).

## Request/Response Pattern Catalog

Choose a pattern that fits your tool; define types in the same handler file or a shared file.

| Pattern | Use when | Example |
|---------|----------|---------|
| **Simple single-value** | One input, one output | `StringRequest` / `StringResponse` in `string.go` |
| **Multi-field** | One payload, multiple input fields | `PathRequest` (value + path) in `json.go` |
| **Two-input comparison** | Two values to compare | `DiffRequest` (valueA, valueB) in `json.go` |
| **Complex options** | Many optional/configuration fields | `LoremRequest` + `LoremOptions` in `lorem.go` |
| **Multiple output fields** | One input, several result fields | `CountdownResponse` (days, hours, minutes, seconds, isPast, text) in `date.go` |

- **Simple**: use shared `decodeBody(w, r)` when your request is `StringRequest`.
- **Custom types**: decode manually (method check + `json.NewDecoder(r.Body).Decode(&req)`); see Handler Pattern below.

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

Use **decodeBody** when the request is `StringRequest` (single `value`). Use **manual decode** for any other request type (multi-field, two-input, options, or custom structs).

With `StringRequest` (decodeBody):

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

With custom request types (manual decode):

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

## Key Findings

- **POST-only** – Every handler validates `r.Method == http.MethodPost` first; reject others with 405.
- **JSON only** – No binary or file-upload handlers; input is `json.NewDecoder(r.Body).Decode()`, output is `writeJSON()`.
- **Validation** – Invalid input → 400; server/processing errors → 500. Validation results (e.g. “is this valid?”) can be returned as structured JSON in the response body instead of only HTTP errors.
- **Single-endpoint dispatch** – For one route serving multiple tools, dispatch on a request field (e.g. `lorem.go` uses `tool`); keep handler logic in one place.

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

## Adding a New Backend Tool

Pick the right handler file and request/response pattern from this skill (inventory and pattern catalog above), then add your handler and register the route in `main.go`. For the full workflow—handler tests, frontend API client, sidebar config, tool component config, and E2E testing—follow the **add-tool** skill checklist.
