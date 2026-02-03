package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
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
