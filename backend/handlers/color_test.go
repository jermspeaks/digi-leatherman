package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"testing"
)

// --- Pure function tests ---

func TestParseColor(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantR   uint8
		wantG   uint8
		wantB   uint8
		wantErr bool
	}{
		{"hex 6 digit", "#ff0000", 255, 0, 0, false},
		{"hex 3 digit", "#f00", 255, 0, 0, false},
		{"hex no hash", "ff0000", 255, 0, 0, false},
		{"hex lowercase", "#abcdef", 171, 205, 239, false},
		{"named color red", "red", 255, 0, 0, false},
		{"named color rebeccapurple", "rebeccapurple", 102, 51, 153, false},
		{"rgb format", "rgb(255, 128, 64)", 255, 128, 64, false},
		{"rgb no spaces", "rgb(255,128,64)", 255, 128, 64, false},
		{"rgba format", "rgba(255, 128, 64, 0.5)", 255, 128, 64, false},
		{"hsl format", "hsl(0, 100%, 50%)", 255, 0, 0, false},
		{"hsl blue", "hsl(240, 100%, 50%)", 0, 0, 255, false},
		{"invalid format", "not-a-color", 0, 0, 0, true},
		{"empty", "", 0, 0, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := parseColor(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Errorf("parseColor(%q) expected error, got nil", tc.in)
				}
				return
			}
			if err != nil {
				t.Errorf("parseColor(%q) unexpected error: %v", tc.in, err)
				return
			}
			if c.R != tc.wantR || c.G != tc.wantG || c.B != tc.wantB {
				t.Errorf("parseColor(%q) = RGB(%d,%d,%d), want RGB(%d,%d,%d)",
					tc.in, c.R, c.G, c.B, tc.wantR, tc.wantG, tc.wantB)
			}
		})
	}
}

func TestHSLToRGB(t *testing.T) {
	cases := []struct {
		name  string
		h, s, l float64
		wantR, wantG, wantB uint8
	}{
		{"red", 0, 1, 0.5, 255, 0, 0},
		{"green", 120, 1, 0.5, 0, 255, 0},
		{"blue", 240, 1, 0.5, 0, 0, 255},
		{"white", 0, 0, 1, 255, 255, 255},
		{"black", 0, 0, 0, 0, 0, 0},
		{"gray", 0, 0, 0.5, 128, 128, 128},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, g, b := hslToRGB(tc.h, tc.s, tc.l)
			if r != tc.wantR || g != tc.wantG || b != tc.wantB {
				t.Errorf("hslToRGB(%.0f, %.1f, %.1f) = RGB(%d,%d,%d), want RGB(%d,%d,%d)",
					tc.h, tc.s, tc.l, r, g, b, tc.wantR, tc.wantG, tc.wantB)
			}
		})
	}
}

func TestRGBToHSL(t *testing.T) {
	cases := []struct {
		name  string
		r, g, b uint8
		wantH, wantS, wantL float64
	}{
		{"red", 255, 0, 0, 0, 1, 0.5},
		{"green", 0, 255, 0, 120, 1, 0.5},
		{"blue", 0, 0, 255, 240, 1, 0.5},
		{"white", 255, 255, 255, 0, 0, 1},
		{"black", 0, 0, 0, 0, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, s, l := rgbToHSL(tc.r, tc.g, tc.b)
			if math.Abs(h-tc.wantH) > 1 || math.Abs(s-tc.wantS) > 0.01 || math.Abs(l-tc.wantL) > 0.01 {
				t.Errorf("rgbToHSL(%d,%d,%d) = (%.1f, %.2f, %.2f), want (%.1f, %.2f, %.2f)",
					tc.r, tc.g, tc.b, h, s, l, tc.wantH, tc.wantS, tc.wantL)
			}
		})
	}
}

func TestRelativeLuminance(t *testing.T) {
	cases := []struct {
		name  string
		color Color
		want  float64
	}{
		{"white", Color{255, 255, 255, 1}, 1.0},
		{"black", Color{0, 0, 0, 1}, 0.0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := relativeLuminance(tc.color)
			if math.Abs(got-tc.want) > 0.01 {
				t.Errorf("relativeLuminance(%v) = %.3f, want %.3f", tc.color, got, tc.want)
			}
		})
	}
}

// --- Handler tests ---

func TestColorConvert(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		checkHex   string
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, ""},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, ""},
		{"invalid color", "POST", `{"value":"not-a-color"}`, http.StatusBadRequest, ""},
		{"hex input", "POST", `{"value":"#ff0000"}`, http.StatusOK, "#ff0000"},
		{"named input", "POST", `{"value":"red"}`, http.StatusOK, "#ff0000"},
		{"rgb input", "POST", `{"value":"rgb(255, 0, 0)"}`, http.StatusOK, "#ff0000"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, ColorConvert, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkHex != "" {
				var res ColorConvertResponse
				if err := json.Unmarshal([]byte(body), &res); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if res.Hex != tc.checkHex {
					t.Errorf("hex = %q, want %q", res.Hex, tc.checkHex)
				}
			}
		})
	}
}

func TestContrastCheck(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		checkAA    bool
		wantAA     bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, false, false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, false, false},
		{"invalid color1", "POST", `{"color1":"not-a-color","color2":"#fff"}`, http.StatusBadRequest, false, false},
		{"invalid color2", "POST", `{"color1":"#000","color2":"not-a-color"}`, http.StatusBadRequest, false, false},
		{"black white high contrast", "POST", `{"color1":"#000000","color2":"#ffffff"}`, http.StatusOK, true, true},
		{"similar colors low contrast", "POST", `{"color1":"#777777","color2":"#888888"}`, http.StatusOK, true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, ContrastCheck, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkAA {
				var res ContrastResponse
				if err := json.Unmarshal([]byte(body), &res); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if res.AANormal != tc.wantAA {
					t.Errorf("AANormal = %v, want %v (ratio: %.2f)", res.AANormal, tc.wantAA, res.Ratio)
				}
			}
		})
	}
}

func TestColorBlindness(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed},
		{"invalid JSON", "POST", "{", http.StatusBadRequest},
		{"invalid color", "POST", `{"value":"not-a-color","type":"protanopia"}`, http.StatusBadRequest},
		{"unknown type", "POST", `{"value":"#ff0000","type":"unknown"}`, http.StatusBadRequest},
		{"protanopia", "POST", `{"value":"#ff0000","type":"protanopia"}`, http.StatusOK},
		{"deuteranopia", "POST", `{"value":"#00ff00","type":"deuteranopia"}`, http.StatusOK},
		{"tritanopia", "POST", `{"value":"#0000ff","type":"tritanopia"}`, http.StatusOK},
		{"achromatopsia", "POST", `{"value":"#ff0000","type":"achromatopsia"}`, http.StatusOK},
		{"default type", "POST", `{"value":"#ff0000"}`, http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, ColorBlindness, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
		})
	}
}

func TestGenerateShades(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		checkCount int
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, 0},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, 0},
		{"invalid color", "POST", `{"value":"not-a-color"}`, http.StatusBadRequest, 0},
		{"default count", "POST", `{"value":"#ff0000"}`, http.StatusOK, 5},
		{"custom count", "POST", `{"value":"#ff0000","count":3}`, http.StatusOK, 3},
		{"max count", "POST", `{"value":"#ff0000","count":100}`, http.StatusOK, 10},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, GenerateShades, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.checkCount > 0 {
				var res ShadesResponse
				if err := json.Unmarshal([]byte(body), &res); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if len(res.Shades) != tc.checkCount {
					t.Errorf("shades count = %d, want %d", len(res.Shades), tc.checkCount)
				}
				if len(res.Tints) != tc.checkCount {
					t.Errorf("tints count = %d, want %d", len(res.Tints), tc.checkCount)
				}
			}
		})
	}
}

func TestGenerateHarmonies(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed},
		{"invalid JSON", "POST", "{", http.StatusBadRequest},
		{"invalid color", "POST", `{"value":"not-a-color"}`, http.StatusBadRequest},
		{"red color", "POST", `{"value":"#ff0000"}`, http.StatusOK},
		{"named color", "POST", `{"value":"blue"}`, http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, GenerateHarmonies, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
			if tc.wantStatus == http.StatusOK {
				var res HarmoniesResponse
				if err := json.Unmarshal([]byte(body), &res); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if res.Complementary == "" {
					t.Error("complementary should not be empty")
				}
				if len(res.Triadic) != 2 {
					t.Errorf("triadic count = %d, want 2", len(res.Triadic))
				}
				if len(res.Analogous) != 2 {
					t.Errorf("analogous count = %d, want 2", len(res.Analogous))
				}
			}
		})
	}
}

func TestExtractPalette(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed},
		{"invalid JSON", "POST", "{", http.StatusBadRequest},
		{"invalid base64", "POST", `{"image":"not-base64!!!"}`, http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, ExtractPalette, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", status, tc.wantStatus, body)
			}
		})
	}
}
