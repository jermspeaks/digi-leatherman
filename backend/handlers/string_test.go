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

// --- Pure function tests ---

func TestWordsFrom(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"empty", "", nil},
		{"single word", "hello", []string{"hello"}},
		{"multiple words", "hello world", []string{"hello", "world"}},
		{"punctuation", "hello, world! foo-bar", []string{"hello", "world", "foo", "bar"}},
		{"spaces only", "   ", nil},
		{"mixed case", "Hello World", []string{"hello", "world"}},
		{"leading trailing space", "  foo bar  ", []string{"foo", "bar"}},
		{"numbers", "a1b2c3", []string{"a1b2c3"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := wordsFrom(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("wordsFrom(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestToSnake(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"single word", "hello", "hello"},
		{"two words", "hello world", "hello_world"},
		{"with punctuation", "hello, world!", "hello_world"},
		{"already snake", "hello_world", "hello_world"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := toSnake(tc.in)
			if got != tc.want {
				t.Errorf("toSnake(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestToKebab(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"single word", "hello", "hello"},
		{"two words", "hello world", "hello-world"},
		{"with punctuation", "hello, world!", "hello-world"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := toKebab(tc.in)
			if got != tc.want {
				t.Errorf("toKebab(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestToCamel(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"single word", "hello", "hello"},
		{"two words", "hello world", "helloWorld"},
		{"three words", "hello world foo", "helloWorldFoo"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := toCamel(tc.in)
			if got != tc.want {
				t.Errorf("toCamel(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestToPascal(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"single word", "hello", "Hello"},
		{"two words", "hello world", "HelloWorld"},
		{"three words", "hello world foo", "HelloWorldFoo"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := toPascal(tc.in)
			if got != tc.want {
				t.Errorf("toPascal(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// --- Handler test helper ---

func runHandler(t *testing.T, handler http.HandlerFunc, method, body string) (int, string) {
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

func parseResult(t *testing.T, body string) string {
	t.Helper()
	var res StringResponse
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatalf("parseResult: %v", err)
	}
	return res.Result
}

// --- HTTP handler tests ---

func TestURLEncode(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"hello world"}`, http.StatusOK, "hello+world", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, URLEncode, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestURLDecode(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"invalid encoded value", "POST", `{"value":"%"}`, http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"hello+world"}`, http.StatusOK, "hello world", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, URLDecode, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestBase64Encode(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"hello"}`, http.StatusOK, "aGVsbG8=", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, Base64Encode, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestBase64Decode(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"invalid base64 value", "POST", `{"value":"!!"}`, http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"aGVsbG8="}`, http.StatusOK, "hello", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, Base64Decode, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestTrim(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"  hello  "}`, http.StatusOK, "hello", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, Trim, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestUpperCase(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"hello"}`, http.StatusOK, "HELLO", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, UpperCase, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestLowerCase(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"HELLO"}`, http.StatusOK, "hello", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, LowerCase, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestCapitalCase(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"hello world"}`, http.StatusOK, "Hello World", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, CapitalCase, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestSnakeCase(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"hello world"}`, http.StatusOK, "hello_world", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, SnakeCase, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestKebabCase(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"hello world"}`, http.StatusOK, "hello-world", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, KebabCase, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestCamelCase(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"hello world"}`, http.StatusOK, "helloWorld", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, CamelCase, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestPascalCase(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"hello world"}`, http.StatusOK, "HelloWorld", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, PascalCase, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}

func TestSentenceCase(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantResult     string
		checkResult    bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"empty value", "POST", `{"value":""}`, http.StatusOK, "", true},
		{"valid input", "POST", `{"value":"HELLO WORLD"}`, http.StatusOK, "Hello world", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, SentenceCase, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}
