package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func parseDateResponse(t *testing.T, body string) DateParseResponse {
	t.Helper()
	var res DateParseResponse
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatalf("parseDateResponse: %v", err)
	}
	return res
}

func parseTimezoneResponse(t *testing.T, body string) TimezoneResponse {
	t.Helper()
	var res TimezoneResponse
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatalf("parseTimezoneResponse: %v", err)
	}
	return res
}

func parseCountdownResponse(t *testing.T, body string) CountdownResponse {
	t.Helper()
	var res CountdownResponse
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatalf("parseCountdownResponse: %v", err)
	}
	return res
}

func TestDateParse(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		body        string
		wantStatus  int
		checkResult bool
		wantISO     string
		wantUnix    int64
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, false, "", 0},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, false, "", 0},
		{"empty value", "POST", `{"value":""}`, http.StatusBadRequest, false, "", 0},
		{"invalid date", "POST", `{"value":"not a date"}`, http.StatusBadRequest, false, "", 0},
		{
			"RFC3339",
			"POST",
			`{"value":"2024-03-15T10:30:00Z"}`,
			http.StatusOK,
			true,
			"2024-03-15T10:30:00Z",
			1710498600,
		},
		{
			"Unix timestamp seconds",
			"POST",
			`{"value":"1710498600"}`,
			http.StatusOK,
			true,
			"2024-03-15T10:30:00Z",
			1710498600,
		},
		{
			"Unix timestamp milliseconds",
			"POST",
			`{"value":"1710498600000"}`,
			http.StatusOK,
			true,
			"2024-03-15T10:30:00Z",
			1710498600,
		},
		{
			"ISO date only",
			"POST",
			`{"value":"2024-03-15"}`,
			http.StatusOK,
			true,
			"2024-03-15T00:00:00Z",
			1710460800,
		},
		{
			"US format",
			"POST",
			`{"value":"03/15/2024"}`,
			http.StatusOK,
			true,
			"2024-03-15T00:00:00Z",
			1710460800,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, DateParse, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkResult {
				res := parseDateResponse(t, body)
				if res.ISO != tc.wantISO {
					t.Errorf("ISO = %q, want %q", res.ISO, tc.wantISO)
				}
				if res.Unix != tc.wantUnix {
					t.Errorf("Unix = %d, want %d", res.Unix, tc.wantUnix)
				}
				if res.RFC2822 == "" {
					t.Error("RFC2822 should not be empty")
				}
				if res.Relative == "" {
					t.Error("Relative should not be empty")
				}
			}
		})
	}
}

func TestDateTimezone(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		checkFn    func(t *testing.T, res TimezoneResponse)
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, nil},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, nil},
		{"invalid date", "POST", `{"value":"bad","fromZone":"UTC","toZone":"UTC"}`, http.StatusBadRequest, nil},
		{"invalid fromZone", "POST", `{"value":"2024-03-15T12:00:00","fromZone":"BadZone","toZone":"UTC"}`, http.StatusBadRequest, nil},
		{"invalid toZone", "POST", `{"value":"2024-03-15T12:00:00","fromZone":"UTC","toZone":"BadZone"}`, http.StatusBadRequest, nil},
		{
			"UTC to America/New_York",
			"POST",
			`{"value":"2024-03-15T12:00:00","fromZone":"UTC","toZone":"America/New_York"}`,
			http.StatusOK,
			func(t *testing.T, res TimezoneResponse) {
				if res.Result == "" {
					t.Error("Result should not be empty")
				}
				if res.ToTime == "" {
					t.Error("ToTime should not be empty")
				}
				if res.FromTime == "" {
					t.Error("FromTime should not be empty")
				}
			},
		},
		{
			"America/Los_Angeles to Europe/London",
			"POST",
			`{"value":"2024-06-15T09:00:00","fromZone":"America/Los_Angeles","toZone":"Europe/London"}`,
			http.StatusOK,
			func(t *testing.T, res TimezoneResponse) {
				if res.Result == "" {
					t.Error("Result should not be empty")
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, DateTimezone, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkFn != nil && status == http.StatusOK {
				res := parseTimezoneResponse(t, body)
				tc.checkFn(t, res)
			}
		})
	}
}

func TestDateCountdown(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		checkFn    func(t *testing.T, res CountdownResponse)
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, nil},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, nil},
		{"invalid date", "POST", `{"value":"bad date"}`, http.StatusBadRequest, nil},
		{
			"past date",
			"POST",
			`{"value":"2020-01-01T00:00:00Z"}`,
			http.StatusOK,
			func(t *testing.T, res CountdownResponse) {
				if !res.IsPast {
					t.Error("IsPast should be true for past date")
				}
				if res.Days <= 0 {
					t.Error("Days should be > 0 for past date")
				}
				if res.Text == "" {
					t.Error("Text should not be empty")
				}
			},
		},
		{
			"future date",
			"POST",
			`{"value":"2099-12-31T23:59:59Z"}`,
			http.StatusOK,
			func(t *testing.T, res CountdownResponse) {
				if res.IsPast {
					t.Error("IsPast should be false for future date")
				}
				if res.Days <= 0 {
					t.Error("Days should be > 0 for future date")
				}
				if res.Text == "" {
					t.Error("Text should not be empty")
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, DateCountdown, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkFn != nil && status == http.StatusOK {
				res := parseCountdownResponse(t, body)
				tc.checkFn(t, res)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty", "", true},
		{"RFC3339", "2024-03-15T10:30:00Z", false},
		{"RFC3339 with offset", "2024-03-15T10:30:00+05:00", false},
		{"Unix seconds", "1710499800", false},
		{"Unix milliseconds", "1710499800000", false},
		{"Unix float", "1710499800.5", false},
		{"ISO date", "2024-03-15", false},
		{"US format", "03/15/2024", false},
		{"US format with time", "03/15/2024 10:30:00", false},
		{"Natural month", "Jan 2, 2006", false},
		{"Natural full", "January 2, 2006", false},
		{"Invalid", "not a date", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseDate(tc.input)
			if tc.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestFormatRelative(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		t           time.Time
		wantContain string
	}{
		{"seconds ago", now.Add(-30 * time.Second), "seconds ago"},
		{"minutes ago", now.Add(-5 * time.Minute), "minutes ago"},
		{"hours ago", now.Add(-2 * time.Hour), "hours ago"},
		{"days ago", now.Add(-3 * 24 * time.Hour), "days ago"},
		{"in seconds", now.Add(30 * time.Second), "in"},
		{"in minutes", now.Add(5 * time.Minute), "in"},
		{"in hours", now.Add(2 * time.Hour), "in"},
		{"in days", now.Add(3 * 24 * time.Hour), "in"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := formatRelative(tc.t)
			if got == "" {
				t.Error("formatRelative returned empty string")
			}
			if !containsSubstring(got, tc.wantContain) {
				t.Errorf("formatRelative = %q, want to contain %q", got, tc.wantContain)
			}
		})
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
