package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// StringRequest is the JSON body for string tool endpoints.
type StringRequest struct {
	Value string `json:"value"`
}

// StringResponse is the JSON response for string tool endpoints.
type StringResponse struct {
	Result string `json:"result"`
}

// URLEncode encodes the request value for safe use in URL query strings (spaces become +, special chars percent-encoded).
// Example: "hello world" -> "hello+world".
func URLEncode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req StringRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	result := url.QueryEscape(req.Value)
	writeJSON(w, StringResponse{Result: result})
}

// URLDecode decodes URL-encoded text back to plain text.
// Example: "hello+world" -> "hello world".
func URLDecode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req StringRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	result, err := url.QueryUnescape(req.Value)
	if err != nil {
		http.Error(w, "invalid encoded value", http.StatusBadRequest)
		return
	}
	writeJSON(w, StringResponse{Result: result})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

var nonWordRe = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// wordsFrom splits s on non-alphanumeric runs and returns non-empty tokens (lowercased for consistency).
func wordsFrom(s string) []string {
	parts := nonWordRe.Split(strings.TrimSpace(s), -1)
	var words []string
	for _, p := range parts {
		if p != "" {
			words = append(words, strings.ToLower(p))
		}
	}
	return words
}

func toSnake(s string) string {
	words := wordsFrom(s)
	return strings.Join(words, "_")
}

func toKebab(s string) string {
	words := wordsFrom(s)
	return strings.Join(words, "-")
}

func toCamel(s string) string {
	words := wordsFrom(s)
	if len(words) == 0 {
		return ""
	}
	caser := cases.Title(language.English)
	for i := 1; i < len(words); i++ {
		words[i] = caser.String(words[i])
	}
	return strings.Join(words, "")
}

func toPascal(s string) string {
	words := wordsFrom(s)
	if len(words) == 0 {
		return ""
	}
	caser := cases.Title(language.English)
	for i := range words {
		words[i] = caser.String(words[i])
	}
	return strings.Join(words, "")
}

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

// Base64Encode encodes the request value as Base64.
// Example: "Hi" -> "SGk=".
func Base64Encode(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	result := base64.StdEncoding.EncodeToString([]byte(req.Value))
	writeJSON(w, StringResponse{Result: result})
}

// Base64Decode decodes Base64 back to plain text.
// Example: "SGk=" -> "Hi".
func Base64Decode(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	result, err := base64.StdEncoding.DecodeString(req.Value)
	if err != nil {
		http.Error(w, "invalid base64 value", http.StatusBadRequest)
		return
	}
	writeJSON(w, StringResponse{Result: string(result)})
}

// Trim removes leading and trailing whitespace from the request value.
// Example: "  hello world  " -> "hello world".
func Trim(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: strings.TrimSpace(req.Value)})
}

// UpperCase converts all characters in the request value to uppercase.
// Example: "Hello World" -> "HELLO WORLD".
func UpperCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: strings.ToUpper(req.Value)})
}

// LowerCase converts all characters in the request value to lowercase.
// Example: "Hello World" -> "hello world".
func LowerCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: strings.ToLower(req.Value)})
}

// CapitalCase converts the request value to title case (first letter of each word uppercase).
// Example: "hello world" -> "Hello World".
func CapitalCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	caser := cases.Title(language.English)
	writeJSON(w, StringResponse{Result: caser.String(req.Value)})
}

// SnakeCase converts words in the request value to snake_case (lowercase with underscores).
// Example: "hello world" -> "hello_world".
func SnakeCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: toSnake(req.Value)})
}

// KebabCase converts words in the request value to kebab-case (lowercase with hyphens).
// Example: "hello world" -> "hello-world".
func KebabCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: toKebab(req.Value)})
}

// CamelCase converts words in the request value to camelCase (first word lowercase, rest capitalized).
// Example: "hello world" -> "helloWorld".
func CamelCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: toCamel(req.Value)})
}

// PascalCase converts words in the request value to PascalCase (each word capitalized).
// Example: "hello world" -> "HelloWorld".
func PascalCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: toPascal(req.Value)})
}

// SentenceCase converts the request value to sentence case (first character uppercase, rest lowercase).
// Example: "HELLO WORLD" -> "Hello world".
func SentenceCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	s := req.Value
	if s == "" {
		writeJSON(w, StringResponse{Result: ""})
		return
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	writeJSON(w, StringResponse{Result: string(runes)})
}

// ParseURLParams parses query parameters from a full URL or raw query string and returns them as JSON.
// Input can be e.g. "https://example.com?foo=bar&baz=qux" or "foo=bar&baz=qux".
func ParseURLParams(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	raw := strings.TrimSpace(req.Value)
	var query string
	if strings.Contains(raw, "?") {
		u, err := url.Parse(raw)
		if err != nil {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}
		query = u.RawQuery
	} else {
		u, err := url.Parse(raw)
		if err == nil && (u.Scheme != "" || u.Host != "") {
			query = u.RawQuery // URL with no query string
		} else {
			query = raw // raw query string
		}
	}
	if query == "" {
		writeJSON(w, StringResponse{Result: "{}"})
		return
	}
	vals, err := url.ParseQuery(query)
	if err != nil {
		http.Error(w, "invalid query string", http.StatusBadRequest)
		return
	}
	// url.Values is map[string][]string; marshal as JSON for readable output
	out, err := json.Marshal(vals)
	if err != nil {
		http.Error(w, "failed to format result", http.StatusInternalServerError)
		return
	}
	writeJSON(w, StringResponse{Result: string(out)})
}

// CreateURLWithParams builds a URL from a base URL (first line) and key=value params (remaining lines).
// Empty lines are ignored. Supports key=value or key: value.
func CreateURLWithParams(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	lines := strings.Split(req.Value, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		http.Error(w, "base URL required (first line)", http.StatusBadRequest)
		return
	}
	baseRaw := strings.TrimSpace(lines[0])
	u, err := url.Parse(baseRaw)
	if err != nil {
		http.Error(w, "invalid base URL", http.StatusBadRequest)
		return
	}
	u.RawQuery = "" // strip existing query
	vals := make(url.Values)
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		key, val := "", ""
		if idx := strings.Index(line, "="); idx >= 0 {
			key = strings.TrimSpace(line[:idx])
			val = strings.TrimSpace(line[idx+1:])
		} else if idx := strings.Index(line, ":"); idx >= 0 {
			key = strings.TrimSpace(line[:idx])
			val = strings.TrimSpace(line[idx+1:])
		}
		if key != "" {
			vals.Set(key, val)
		}
	}
	u.RawQuery = vals.Encode()
	writeJSON(w, StringResponse{Result: u.String()})
}
