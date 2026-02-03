package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseRecorder wraps http.ResponseWriter to capture status code and bytes written.
type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(code int) {
	if r.status == 0 {
		r.status = code
	}
	r.ResponseWriter.WriteHeader(code)
}

// Logging wraps next with request logging: method, path, remote, status, duration.
// Uses slog: Info for 2xx/3xx, Warn for 4xx, Error for 5xx.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		if rec.status == 0 {
			rec.status = http.StatusOK
		}
		duration := time.Since(start)
		attrs := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
			"status", rec.status,
			"duration_ms", duration.Milliseconds(),
		}
		switch {
		case rec.status >= 500:
			slog.Error("request", attrs...)
		case rec.status >= 400:
			slog.Warn("request", attrs...)
		default:
			slog.Info("request", attrs...)
		}
	})
}
