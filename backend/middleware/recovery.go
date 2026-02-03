package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery wraps next with panic recovery: on panic, log stack trace and return 500.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered",
					"err", err,
					"path", r.URL.Path,
					"method", r.Method,
					"stack", string(debug.Stack()),
				)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
