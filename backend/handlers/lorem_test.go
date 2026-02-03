package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func runLoremHandler(t *testing.T, method, body string) (int, string) {
	t.Helper()
	var bodyReader *bytes.Reader
	if body != "" {
		bodyReader = bytes.NewReader([]byte(body))
	} else {
		bodyReader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, "http://test", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	LoremIpsum(rec, req)
	return rec.Code, strings.TrimSpace(rec.Body.String())
}

func TestLoremIpsum(t *testing.T) {
	t.Run("method not allowed", func(t *testing.T) {
		status, _ := runLoremHandler(t, "GET", "")
		if status != http.StatusMethodNotAllowed {
			t.Errorf("status = %d, want %d", status, http.StatusMethodNotAllowed)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		status, _ := runLoremHandler(t, "POST", "{")
		if status != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"type":"lines","count":5}`)
		if status != http.StatusBadRequest {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusBadRequest, body)
		}
		if !strings.Contains(body, "invalid type") {
			t.Errorf("body should mention invalid type, got: %s", body)
		}
	})

	t.Run("count too low", func(t *testing.T) {
		status, _ := runLoremHandler(t, "POST", `{"type":"words","count":0}`)
		if status != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("count too high for words", func(t *testing.T) {
		status, _ := runLoremHandler(t, "POST", `{"type":"words","count":1001}`)
		if status != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("count too high for sentences", func(t *testing.T) {
		status, _ := runLoremHandler(t, "POST", `{"type":"sentences","count":101}`)
		if status != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("count too high for paragraphs", func(t *testing.T) {
		status, _ := runLoremHandler(t, "POST", `{"type":"paragraphs","count":51}`)
		if status != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("words returns correct count", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"type":"words","count":10}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		words := strings.Fields(result)
		if len(words) != 10 {
			t.Errorf("result word count = %d, want 10; result: %q", len(words), result)
		}
		// All words should be from lorem set (lowercase)
		for _, w := range words {
			if w == "" {
				continue
			}
			lower := strings.ToLower(w)
			found := false
			for _, lw := range loremWords {
				if lower == lw {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("word %q not in lorem word list", w)
			}
		}
	})

	t.Run("sentences returns correct count", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"type":"sentences","count":3}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		// Sentences are joined with " " and end with "."
		sentences := strings.Split(result, ". ")
		// Last might be "foo." so we count non-empty segments that end with .
		var count int
		for _, s := range sentences {
			s = strings.TrimSpace(s)
			if s != "" {
				count++
			}
		}
		if count != 3 {
			t.Errorf("sentence count = %d, want 3; result: %q", count, result)
		}
	})

	t.Run("paragraphs returns correct count", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"type":"paragraphs","count":2}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		paragraphs := strings.Split(result, "\n\n")
		if len(paragraphs) != 2 {
			t.Errorf("paragraph count = %d, want 2; result: %q", len(paragraphs), result)
		}
	})

	t.Run("type case insensitive", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"type":"WORDS","count":5}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		words := strings.Fields(result)
		if len(words) != 5 {
			t.Errorf("result word count = %d, want 5", len(words))
		}
	})
}
