package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func runUUIDHandler(t *testing.T, handler http.HandlerFunc, method, body string) (int, string) {
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

func parseUUIDResult(t *testing.T, body string) string {
	t.Helper()
	var res UUIDGenerateResponse
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatalf("parseUUIDResult: %v", err)
	}
	return res.Result
}

func TestGenerateUUIDv4(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{"method not allowed", "GET", http.StatusMethodNotAllowed},
		{"valid generation", "POST", http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runUUIDHandler(t, GenerateUUIDv4, tc.method, "")
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.wantStatus == http.StatusOK {
				result := parseUUIDResult(t, body)
				parsed, err := uuid.Parse(result)
				if err != nil {
					t.Errorf("invalid UUID generated: %v", err)
				}
				if parsed.Version() != 4 {
					t.Errorf("UUID version = %d, want 4", parsed.Version())
				}
			}
		})
	}
}

func TestGenerateUUIDv7(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{"method not allowed", "GET", http.StatusMethodNotAllowed},
		{"valid generation", "POST", http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runUUIDHandler(t, GenerateUUIDv7, tc.method, "")
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.wantStatus == http.StatusOK {
				result := parseUUIDResult(t, body)
				parsed, err := uuid.Parse(result)
				if err != nil {
					t.Errorf("invalid UUID generated: %v", err)
				}
				if parsed.Version() != 7 {
					t.Errorf("UUID version = %d, want 7", parsed.Version())
				}
			}
		})
	}
}

func TestBulkGenerateUUID(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		body        string
		wantStatus  int
		wantCount   int
		wantVersion int
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, 0, 0},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, 0, 0},
		{"generate 5 v4", "POST", `{"version":4,"count":5}`, http.StatusOK, 5, 4},
		{"generate 3 v7", "POST", `{"version":7,"count":3}`, http.StatusOK, 3, 7},
		{"default version when invalid", "POST", `{"version":99,"count":2}`, http.StatusOK, 2, 4},
		{"default count when zero", "POST", `{"version":4,"count":0}`, http.StatusOK, 1, 4},
		{"cap count at 1000", "POST", `{"version":4,"count":2000}`, http.StatusOK, 1000, 4},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runUUIDHandler(t, BulkGenerateUUID, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.wantStatus == http.StatusOK {
				var res BulkGenerateResponse
				if err := json.Unmarshal([]byte(body), &res); err != nil {
					t.Fatalf("parse response: %v", err)
				}
				if len(res.Results) != tc.wantCount {
					t.Errorf("result count = %d, want %d", len(res.Results), tc.wantCount)
				}
				for _, r := range res.Results {
					parsed, err := uuid.Parse(r)
					if err != nil {
						t.Errorf("invalid UUID: %v", err)
					}
					if int(parsed.Version()) != tc.wantVersion {
						t.Errorf("UUID version = %d, want %d", parsed.Version(), tc.wantVersion)
					}
				}
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		body        string
		wantStatus  int
		wantValid   bool
		wantVersion int
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, false, 0},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, false, 0},
		{"valid v4 UUID", "POST", `{"value":"550e8400-e29b-41d4-a716-446655440000"}`, http.StatusOK, true, 4},
		{"valid v4 UUID uppercase", "POST", `{"value":"550E8400-E29B-41D4-A716-446655440000"}`, http.StatusOK, true, 4},
		{"invalid UUID", "POST", `{"value":"not-a-uuid"}`, http.StatusOK, false, 0},
		{"empty value", "POST", `{"value":""}`, http.StatusOK, false, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runUUIDHandler(t, ValidateUUID, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.wantStatus == http.StatusOK {
				var res UUIDValidateResponse
				if err := json.Unmarshal([]byte(body), &res); err != nil {
					t.Fatalf("parse response: %v", err)
				}
				if res.Valid != tc.wantValid {
					t.Errorf("valid = %v, want %v", res.Valid, tc.wantValid)
				}
				if tc.wantValid && res.Version != tc.wantVersion {
					t.Errorf("version = %d, want %d", res.Version, tc.wantVersion)
				}
			}
		})
	}
}

func TestParseUUID(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		body        string
		wantStatus  int
		wantValid   bool
		wantVersion int
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, false, 0},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, false, 0},
		{"valid v4 UUID", "POST", `{"value":"550e8400-e29b-41d4-a716-446655440000"}`, http.StatusOK, true, 4},
		{"invalid UUID", "POST", `{"value":"invalid"}`, http.StatusOK, false, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runUUIDHandler(t, ParseUUID, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.wantStatus == http.StatusOK {
				var res UUIDParseResponse
				if err := json.Unmarshal([]byte(body), &res); err != nil {
					t.Fatalf("parse response: %v", err)
				}
				if res.Valid != tc.wantValid {
					t.Errorf("valid = %v, want %v", res.Valid, tc.wantValid)
				}
				if tc.wantValid && res.Version != tc.wantVersion {
					t.Errorf("version = %d, want %d", res.Version, tc.wantVersion)
				}
			}
		})
	}
}

func TestFormatUUID(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		body        string
		wantStatus  int
		wantResult  string
		checkResult bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"invalid UUID", "POST", `{"value":"invalid","format":"standard"}`, http.StatusBadRequest, "", false},
		{"standard format", "POST", `{"value":"550e8400-e29b-41d4-a716-446655440000","format":"standard"}`, http.StatusOK, "550e8400-e29b-41d4-a716-446655440000", true},
		{"no-hyphens format", "POST", `{"value":"550e8400-e29b-41d4-a716-446655440000","format":"no-hyphens"}`, http.StatusOK, "550e8400e29b41d4a716446655440000", true},
		{"uppercase format", "POST", `{"value":"550e8400-e29b-41d4-a716-446655440000","format":"uppercase"}`, http.StatusOK, "550E8400-E29B-41D4-A716-446655440000", true},
		{"braces format", "POST", `{"value":"550e8400-e29b-41d4-a716-446655440000","format":"braces"}`, http.StatusOK, "{550e8400-e29b-41d4-a716-446655440000}", true},
		{"urn format", "POST", `{"value":"550e8400-e29b-41d4-a716-446655440000","format":"urn"}`, http.StatusOK, "urn:uuid:550e8400-e29b-41d4-a716-446655440000", true},
		{"unknown format", "POST", `{"value":"550e8400-e29b-41d4-a716-446655440000","format":"unknown"}`, http.StatusBadRequest, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runUUIDHandler(t, FormatUUID, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				result := parseUUIDResult(t, body)
				if result != tc.wantResult {
					t.Errorf("result = %q, want %q", result, tc.wantResult)
				}
			}
		})
	}
}
