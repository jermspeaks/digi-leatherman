package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ValidateResponse is the JSON response for the validate endpoint.
type ValidateResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// PathRequest is the JSON body for the path query endpoint.
type PathRequest struct {
	Value string `json:"value"`
	Path  string `json:"path"`
}

// DiffRequest is the JSON body for the diff endpoint.
type DiffRequest struct {
	ValueA string `json:"valueA"`
	ValueB string `json:"valueB"`
}

// FormatJSON pretty-prints JSON with 2-space indentation.
func FormatJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req StringRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	var v interface{}
	if err := json.Unmarshal([]byte(req.Value), &v); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("marshal: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, StringResponse{Result: string(out)})
}

// MinifyJSON removes unnecessary whitespace from JSON.
func MinifyJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req StringRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	var v interface{}
	if err := json.Unmarshal([]byte(req.Value), &v); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}
	out, err := json.Marshal(v)
	if err != nil {
		http.Error(w, fmt.Sprintf("marshal: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, StringResponse{Result: string(out)})
}

// ValidateJSON checks whether the input is valid JSON.
func ValidateJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req StringRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	valid := json.Valid([]byte(req.Value))
	resp := ValidateResponse{Valid: valid}
	if !valid {
		var v interface{}
		if err := json.Unmarshal([]byte(req.Value), &v); err != nil {
			resp.Error = err.Error()
		}
	}
	writeJSON(w, resp)
}

// pathGet walks a dot-separated path with optional numeric indices (e.g. "a.b.0.c").
// Returns the value at the path and true, or nil and false if not found.
func pathGet(v interface{}, path string) (interface{}, bool) {
	path = strings.TrimSpace(path)
	if path == "" {
		return v, true
	}
	parts := splitPath(path)
	for _, key := range parts {
		if v == nil {
			return nil, false
		}
		switch m := v.(type) {
		case map[string]interface{}:
			if k, ok := key.(string); ok {
				val, exists := m[k]
				if !exists {
					return nil, false
				}
				v = val
			} else {
				return nil, false
			}
		case []interface{}:
			if i, ok := key.(int); ok && i >= 0 && i < len(m) {
				v = m[i]
			} else {
				return nil, false
			}
		default:
			return nil, false
		}
	}
	return v, true
}

// splitPath splits "a.b.0.c" into ["a","b",0,"c"] (numeric segments as int).
func splitPath(path string) []interface{} {
	var parts []interface{}
	for _, s := range strings.Split(path, ".") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if i, err := strconv.Atoi(s); err == nil {
			parts = append(parts, i)
		} else {
			parts = append(parts, s)
		}
	}
	return parts
}

// PathQueryJSON extracts a value at the given path from JSON.
func PathQueryJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req PathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	var v interface{}
	if err := json.Unmarshal([]byte(req.Value), &v); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}
	got, ok := pathGet(v, req.Path)
	if !ok {
		http.Error(w, "path not found", http.StatusBadRequest)
		return
	}
	out, err := json.Marshal(got)
	if err != nil {
		http.Error(w, fmt.Sprintf("marshal: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, StringResponse{Result: string(out)})
}

// pathJoin returns path + "." + segment, or just segment if path is empty.
func pathJoin(path, segment string) string {
	if path == "" {
		return segment
	}
	return path + "." + segment
}

// diffRecurse builds a structural diff between two values; returns a slice of "path: left -> right" lines.
func diffRecurse(a, b interface{}, path string) []string {
	var lines []string
	pathLabel := path
	if pathLabel == "" {
		pathLabel = "(root)"
	}
	switch av := a.(type) {
	case map[string]interface{}:
		bv, ok := b.(map[string]interface{})
		if !ok {
			lines = append(lines, fmt.Sprintf("%s: (object) -> (%T)", pathLabel, b))
			return lines
		}
		seen := make(map[string]bool)
		for k, v := range av {
			seen[k] = true
			p := pathJoin(path, k)
			if bval, has := bv[k]; has {
				lines = append(lines, diffRecurse(v, bval, p)...)
			} else {
				lines = append(lines, fmt.Sprintf("%s: %v -> (missing)", p, v))
			}
		}
		for k, v := range bv {
			if !seen[k] {
				lines = append(lines, fmt.Sprintf("%s: (missing) -> %v", pathJoin(path, k), v))
			}
		}
	case []interface{}:
		bv, ok := b.([]interface{})
		if !ok {
			lines = append(lines, fmt.Sprintf("%s: (array) -> (%T)", pathLabel, b))
			return lines
		}
		max := len(av)
		if len(bv) > max {
			max = len(bv)
		}
		for i := 0; i < max; i++ {
			p := pathJoin(path, strconv.Itoa(i))
			if i >= len(av) {
				lines = append(lines, fmt.Sprintf("%s: (missing) -> %v", p, bv[i]))
			} else if i >= len(bv) {
				lines = append(lines, fmt.Sprintf("%s: %v -> (missing)", p, av[i]))
			} else {
				lines = append(lines, diffRecurse(av[i], bv[i], p)...)
			}
		}
	default:
		if a != b {
			lines = append(lines, fmt.Sprintf("%s: %v -> %v", pathLabel, a, b))
		}
	}
	return lines
}

// DiffJSON compares two JSON values and returns a structural diff.
func DiffJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req DiffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	var a, b interface{}
	if err := json.Unmarshal([]byte(req.ValueA), &a); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON in valueA: %v", err), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal([]byte(req.ValueB), &b); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON in valueB: %v", err), http.StatusBadRequest)
		return
	}
	lines := diffRecurse(a, b, "")
	result := strings.Join(lines, "\n")
	if result == "" {
		result = "(no differences)"
	}
	writeJSON(w, StringResponse{Result: result})
}
