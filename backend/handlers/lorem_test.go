package handlers

import (
	"bytes"
	"encoding/json"
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

	t.Run("generator with startWithClassic", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"generator","count":2,"options":{"type":"paragraphs","startWithClassic":true}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		if !strings.HasPrefix(result, "Lorem ipsum dolor sit amet") {
			t.Errorf("result should start with classic Latin, got: %q", result[:min(60, len(result))])
		}
		paragraphs := strings.Split(result, "\n\n")
		if len(paragraphs) != 2 {
			t.Errorf("paragraph count = %d, want 2", len(paragraphs))
		}
	})

	t.Run("generator with vocabulary bacon", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"generator","count":5,"options":{"type":"words","vocabulary":"bacon"}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		words := strings.Fields(result)
		if len(words) != 5 {
			t.Errorf("word count = %d, want 5", len(words))
		}
		// At least one word should be from bacon set
		found := false
		for _, w := range words {
			for _, bw := range baconWords {
				if strings.ToLower(w) == bw {
					found = true
					break
				}
			}
		}
		if !found {
			t.Errorf("expected at least one bacon word in %q", result)
		}
	})

	t.Run("invalid tool", func(t *testing.T) {
		status, _ := runLoremHandler(t, "POST", `{"tool":"unknown","count":5}`)
		if status != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("characters returns correct length", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"characters","count":50}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		if len(result) != 50 {
			t.Errorf("result length = %d, want 50", len(result))
		}
	})

	t.Run("characters wholeWordsOnly", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"characters","count":20,"options":{"wholeWordsOnly":true}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		if len(result) > 20 {
			t.Errorf("wholeWordsOnly result length = %d, should be <= 20", len(result))
		}
		// Should not end mid-word (no trailing partial word)
		words := strings.Fields(result)
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
				for _, bw := range baconWords {
					if lower == bw {
						found = true
						break
					}
				}
			}
			if !found {
				for _, hw := range hipsterWords {
					if lower == hw {
						found = true
						break
					}
				}
			}
			if !found {
				t.Errorf("word %q not in any vocabulary", w)
			}
		}
	})

	t.Run("bytes returns correct byte count", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"bytes","count":30}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		if len([]byte(result)) != 30 {
			t.Errorf("result byte length = %d, want 30", len([]byte(result)))
		}
	})

	t.Run("title returns word count", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"title","count":5}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		words := strings.Fields(result)
		if len(words) != 5 {
			t.Errorf("title word count = %d, want 5", len(words))
		}
		if result != "" && result[0] < 'A' || result[0] > 'Z' {
			t.Errorf("title should be capitalized, got %q", result[:1])
		}
	})

	t.Run("slug format", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"slug","count":4}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		parts := strings.Split(result, "-")
		if len(parts) != 4 {
			t.Errorf("slug parts = %d, want 4; got %q", len(parts), result)
		}
		if result != strings.ToLower(result) {
			t.Errorf("slug should be lowercase, got %q", result)
		}
	})

	t.Run("camelCase format", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"camelcase","count":3}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		if result == "" {
			t.Errorf("camelCase result empty")
		}
		if result[0] < 'a' || result[0] > 'z' {
			t.Errorf("camelCase should start with lower, got %q", result[:1])
		}
		if strings.Contains(result, " ") || strings.Contains(result, "-") {
			t.Errorf("camelCase should have no spaces or hyphens, got %q", result)
		}
	})

	t.Run("list bullet", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"list","count":3,"options":{"listStyle":"bullet"}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		lines := strings.Split(result, "\n")
		if len(lines) != 3 {
			t.Errorf("list lines = %d, want 3", len(lines))
		}
		for i, line := range lines {
			if !strings.HasPrefix(line, "- ") {
				t.Errorf("line %d should start with '- ', got %q", i, line)
			}
		}
	})

	t.Run("list numbered", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"list","count":2,"options":{"listStyle":"numbered"}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		lines := strings.Split(result, "\n")
		if len(lines) != 2 {
			t.Errorf("list lines = %d, want 2", len(lines))
		}
		if !strings.HasPrefix(lines[0], "1. ") {
			t.Errorf("first line should start with '1. ', got %q", lines[0])
		}
		if !strings.HasPrefix(lines[1], "2. ") {
			t.Errorf("second line should start with '2. ', got %q", lines[1])
		}
	})

	t.Run("headings count", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"headings","count":4,"options":{"headingStyle":"plain"}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		lines := strings.Split(strings.TrimSuffix(result, "\n"), "\n")
		if len(lines) != 4 {
			t.Errorf("headings lines = %d, want 4", len(lines))
		}
	})

	t.Run("headings markdown style", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"headings","count":2,"options":{"headingStyle":"markdown"}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		if !strings.HasPrefix(result, "# ") {
			t.Errorf("markdown headings should start with '# ', got %q", result[:20])
		}
	})

	t.Run("html paragraphs", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"html","count":2,"options":{"format":"paragraphs"}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		if !strings.Contains(result, "<p>") || !strings.Contains(result, "</p>") {
			t.Errorf("html should contain <p>...</p>, got %q", result[:min(80, len(result))])
		}
	})

	t.Run("html list", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"html","count":2,"options":{"format":"list"}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		if !strings.Contains(result, "<ul>") || !strings.Contains(result, "<li>") {
			t.Errorf("html list should contain <ul> and <li>, got %q", result[:min(80, len(result))])
		}
	})

	t.Run("markdown paragraphs", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"markdown","count":2,"options":{"format":"paragraphs"}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		paragraphs := strings.Split(result, "\n\n")
		if len(paragraphs) < 2 {
			t.Errorf("markdown paragraphs count = %d, want at least 2", len(paragraphs))
		}
	})

	t.Run("markdown list", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"markdown","count":2,"options":{"format":"list"}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		if !strings.Contains(result, "- ") {
			t.Errorf("markdown list should contain '- ', got %q", result[:min(80, len(result))])
		}
	})

	t.Run("json valid and keys", func(t *testing.T) {
		status, body := runLoremHandler(t, "POST", `{"tool":"json","count":0,"options":{"keys":["title","body"]}}`)
		if status != http.StatusOK {
			t.Errorf("status = %d, want %d; body: %s", status, http.StatusOK, body)
			return
		}
		result := parseResult(t, body)
		var m map[string]string
		if err := json.Unmarshal([]byte(result), &m); err != nil {
			t.Errorf("json result should be valid JSON: %v", err)
		}
		if m["title"] == "" || m["body"] == "" {
			t.Errorf("json should have title and body keys with values, got %v", m)
		}
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
