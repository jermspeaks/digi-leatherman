package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"digi-leatherman/backend/handlers"
)

// buildMux builds the same mux as main() for integration testing.
func buildMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/color/convert", cors(handlers.ColorConvert))
	mux.HandleFunc("/api/color/contrast", cors(handlers.ContrastCheck))
	mux.HandleFunc("/api/color/blindness", cors(handlers.ColorBlindness))
	mux.HandleFunc("/api/color/palette", cors(handlers.ExtractPalette))
	mux.HandleFunc("/api/color/shades", cors(handlers.GenerateShades))
	mux.HandleFunc("/api/color/harmonies", cors(handlers.GenerateHarmonies))
	return mux
}

func TestColorEndpointsE2E(t *testing.T) {
	mux := buildMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	base := server.URL

	tests := []struct {
		name       string
		path       string
		body       string
		wantStatus int
		wantSubstr string
	}{
		{
			name:       "convert hex",
			path:       "/api/color/convert",
			body:       `{"value":"#ff6b6b"}`,
			wantStatus: http.StatusOK,
			wantSubstr: `"hex":"#ff6b6b"`,
		},
		{
			name:       "convert invalid",
			path:       "/api/color/convert",
			body:       `{"value":"not-a-color"}`,
			wantStatus: http.StatusBadRequest,
			wantSubstr: "",
		},
		{
			name:       "contrast black white",
			path:       "/api/color/contrast",
			body:       `{"color1":"#000000","color2":"#ffffff"}`,
			wantStatus: http.StatusOK,
			wantSubstr: `"ratio":21`,
		},
		{
			name:       "blindness protanopia",
			path:       "/api/color/blindness",
			body:       `{"value":"#ff0000","type":"protanopia"}`,
			wantStatus: http.StatusOK,
			wantSubstr: `"simulated"`,
		},
		{
			name:       "shades",
			path:       "/api/color/shades",
			body:       `{"value":"#3b82f6","count":3}`,
			wantStatus: http.StatusOK,
			wantSubstr: `"shades"`,
		},
		{
			name:       "harmonies",
			path:       "/api/color/harmonies",
			body:       `{"value":"#3b82f6"}`,
			wantStatus: http.StatusOK,
			wantSubstr: `"complementary"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, base+tt.path, bytes.NewReader([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("status = %d, want %d; body: %s", resp.StatusCode, tt.wantStatus, string(body))
				return
			}

			if tt.wantSubstr != "" {
				body, _ := io.ReadAll(resp.Body)
				if !strings.Contains(string(body), tt.wantSubstr) {
					t.Errorf("response body %q does not contain %q", string(body), tt.wantSubstr)
				}
			}
		})
	}
}

func TestColorPaletteEndpointInvalidImage(t *testing.T) {
	mux := buildMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/color/palette",
		bytes.NewReader([]byte(`{"image":"not-valid-base64!!!"}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("status = %d, want 400; body: %s", resp.StatusCode, string(body))
	}
}
