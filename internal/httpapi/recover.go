package httpapi

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				log.Error("panic", "value", v, "stack", string(debug.Stack()))
				writeJSON(w, http.StatusInternalServerError, map[string]string{"code": "internal_error", "message": "internal server error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
