package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func runJSONHandler(t *testing.T, handler http.HandlerFunc, method, body string) (int, string) {
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
	handler(rec, req)
	return rec.Code, strings.TrimSpace(rec.Body.String())
}

func parseJSONResult(t *testing.T, body string) string {
	t.Helper()
	var res StringResponse
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatalf("parseJSONResult: %v", err)
	}
	return res.Result
}

func parseValidateResponse(t *testing.T, body string) (valid bool, errMsg string) {
	t.Helper()
	var res ValidateResponse
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatalf("parseValidateResponse: %v", err)
	}
	return res.Valid, res.Error
}

func TestFormatJSON(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		body        string
		wantStatus  int
		wantResult  string
		checkResult bool
		checkFormat bool // result should be pretty-printed (contains newline)
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false, false},
		{"invalid JSON body", "POST", "{", http.StatusBadRequest, "", false, false},
		{"invalid JSON value", "POST", `{"value":"not json"}`, http.StatusBadRequest, "", false, false},
		{"valid input", "POST", `{"value":"{\"a\":1}"}`, http.StatusOK, "", true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runJSONHandler(t, FormatJSON, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseJSONResult(t, body)
				if tc.checkFormat && !strings.Contains(got, "\n") {
					t.Errorf("expected formatted output with newlines, got %q", got)
				}
				// Round-trip: formatted output should be valid JSON
				var v interface{}
				if err := json.Unmarshal([]byte(got), &v); err != nil {
					t.Errorf("formatted output is not valid JSON: %v", err)
				}
			}
		})
	}
}

func TestMinifyJSON(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		body        string
		wantStatus  int
		wantResult  string
		checkResult bool
		checkMinify bool // result should have no unnecessary spaces
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false, false},
		{"invalid JSON body", "POST", "{", http.StatusBadRequest, "", false, false},
		{"invalid JSON value", "POST", `{"value":"{"}`, http.StatusBadRequest, "", false, false},
		{"valid input", "POST", `{"value":"{\"a\": 1, \"b\": 2}"}`, http.StatusOK, `{"a":1,"b":2}`, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runJSONHandler(t, MinifyJSON, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseJSONResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
				if tc.checkMinify && strings.Contains(got, " ") {
					t.Errorf("minified output should not contain spaces, got %q", got)
				}
			}
		})
	}
}

func TestValidateJSON(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantValid  bool
		checkValid bool
		wantError  bool // error field should be non-empty when valid is false
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, false, false, false},
		{"invalid JSON body", "POST", "{", http.StatusBadRequest, false, false, false},
		{"valid JSON", "POST", `{"value":"{\"a\":1}"}`, http.StatusOK, true, true, false},
		{"invalid JSON", "POST", `{"value":"{invalid}"}`, http.StatusOK, false, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runJSONHandler(t, ValidateJSON, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkValid {
				valid, errMsg := parseValidateResponse(t, body)
				if valid != tc.wantValid {
					t.Errorf("valid = %v, want %v", valid, tc.wantValid)
				}
				if tc.wantError && errMsg == "" {
					t.Errorf("expected non-empty error message for invalid JSON")
				}
			}
		})
	}
}

func TestPathQueryJSON(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		body        string
		wantStatus  int
		wantResult  string
		checkResult bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON body", "POST", "{", http.StatusBadRequest, "", false},
		{"invalid JSON value", "POST", `{"value":"{","path":"a"}`, http.StatusBadRequest, "", false},
		{"path not found", "POST", `{"value":"{\"a\":1}","path":"b"}`, http.StatusBadRequest, "", false},
		{"path found", "POST", `{"value":"{\"a\":{\"b\":42}}","path":"a.b"}`, http.StatusOK, "42", true},
		{"path array index", "POST", `{"value":"{\"x\":[1,2,3]}","path":"x.1"}`, http.StatusOK, "2", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runJSONHandler(t, PathQueryJSON, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseJSONResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestDiffJSON(t *testing.T) {
	cases := []struct {
		name         string
		method       string
		body         string
		wantStatus   int
		wantContains string
		checkResult  bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON body", "POST", "{", http.StatusBadRequest, "", false},
		{"invalid valueA", "POST", `{"valueA":"{","valueB":"{}"}`, http.StatusBadRequest, "", false},
		{"invalid valueB", "POST", `{"valueA":"{}","valueB":"{"}`, http.StatusBadRequest, "", false},
		{"identical", "POST", `{"valueA":"{\"a\":1}","valueB":"{\"a\":1}"}`, http.StatusOK, "no differences", true},
		{"different", "POST", `{"valueA":"{\"a\":1}","valueB":"{\"a\":2}"}`, http.StatusOK, "diff", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runJSONHandler(t, DiffJSON, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult && tc.wantContains != "" {
				// For "diff" check decoded result contains "->" (diff arrow)
				if tc.wantContains == "diff" {
					result := parseJSONResult(t, body)
					if !strings.Contains(result, "->") {
						t.Errorf("result should contain diff arrow, got %s", result)
					}
				} else if !strings.Contains(body, tc.wantContains) {
					t.Errorf("body should contain %q, got %s", tc.wantContains, body)
				}
			}
		})
	}
}

func TestSplitPath(t *testing.T) {
	cases := []struct {
		path string
		want []interface{}
	}{
		{"a.b", []interface{}{"a", "b"}},
		{"a.b.0", []interface{}{"a", "b", 0}},
		{"x.1.y", []interface{}{"x", 1, "y"}},
		{"", nil},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got := splitPath(tc.path)
			if len(got) != len(tc.want) {
				t.Errorf("splitPath(%q) = %v, want %v", tc.path, got, tc.want)
				return
			}
			for i := range got {
				switch g := got[i].(type) {
				case string:
					if w, ok := tc.want[i].(string); !ok || g != w {
						t.Errorf("at %d: got %q (string), want %v", i, g, tc.want[i])
					}
				case int:
					if w, ok := tc.want[i].(int); !ok || g != w {
						t.Errorf("at %d: got %d (int), want %v", i, g, tc.want[i])
					}
				default:
					t.Errorf("at %d: unexpected type %T", i, got[i])
				}
			}
		})
	}
}

func TestPathGet(t *testing.T) {
	v := map[string]interface{}{
		"a": map[string]interface{}{
			"b": 42,
		},
		"x": []interface{}{10, 20, 30},
	}
	cases := []struct {
		path string
		want interface{}
		ok   bool
	}{
		{"", v, true},
		{"a", v["a"], true},
		{"a.b", 42, true},
		{"x.1", 20, true},
		{"missing", nil, false},
		{"a.missing", nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got, ok := pathGet(v, tc.path)
			if ok != tc.ok {
				t.Errorf("pathGet(%q) ok = %v, want %v", tc.path, ok, tc.ok)
				return
			}
			if ok && !reflect.DeepEqual(got, tc.want) {
				t.Errorf("pathGet(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}
