package httpapi

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/service"
	"net/http"
	"strings"
)

type ctxKey string

const userKey ctxKey = "user"

func AuthMiddleware(a service.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				writeErr(w, domain.ErrUnauthorized)
				return
			}
			u, e := a.Current(r.Context(), strings.TrimPrefix(h, "Bearer "))
			if e != nil {
				writeErr(w, e)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey, u)))
		})
	}
}
func userFrom(r *http.Request) domain.User {
	u, _ := r.Context().Value(userKey).(domain.User)
	return u
}
