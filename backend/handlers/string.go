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

// URLEncode encodes the request value using url.QueryEscape and returns JSON.
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

// URLDecode decodes the request value using url.QueryUnescape and returns JSON.
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

// Base64Encode encodes the request value as base64 and returns JSON.
func Base64Encode(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	result := base64.StdEncoding.EncodeToString([]byte(req.Value))
	writeJSON(w, StringResponse{Result: result})
}

// Base64Decode decodes the request value from base64 and returns JSON.
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

// Trim trims leading and trailing whitespace and returns JSON.
func Trim(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: strings.TrimSpace(req.Value)})
}

// UpperCase converts the request value to uppercase and returns JSON.
func UpperCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: strings.ToUpper(req.Value)})
}

// LowerCase converts the request value to lowercase and returns JSON.
func LowerCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: strings.ToLower(req.Value)})
}

// CapitalCase converts the request value to title case and returns JSON.
func CapitalCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	caser := cases.Title(language.English)
	writeJSON(w, StringResponse{Result: caser.String(req.Value)})
}

// SnakeCase converts the request value to snake_case and returns JSON.
func SnakeCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: toSnake(req.Value)})
}

// KebabCase converts the request value to kebab-case and returns JSON.
func KebabCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: toKebab(req.Value)})
}

// CamelCase converts the request value to camelCase and returns JSON.
func CamelCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: toCamel(req.Value)})
}

// PascalCase converts the request value to PascalCase and returns JSON.
func PascalCase(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	writeJSON(w, StringResponse{Result: toPascal(req.Value)})
}

// SentenceCase converts the first character to uppercase and the rest to lowercase and returns JSON.
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
