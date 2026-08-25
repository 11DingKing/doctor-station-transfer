package httpapi

import (
	"encoding/json"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type RouterDeps struct {
	Auth      service.Auth
	Projects  service.Projects
	Reviews   service.Reviews
	Transfers service.Transfers
}

func New(d RouterDeps) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	r.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	})
	r.Post("/v1/auth/login", d.login)
	r.With(AuthMiddleware(d.Auth)).Post("/v1/auth/logout", d.logout)
	r.With(AuthMiddleware(d.Auth)).Route("/v1/projects", func(sr chi.Router) {
		sr.Get("/", d.listProjects)
		sr.Post("/", d.createProject)
		sr.Post("/{id}/transition", d.transitionProject)
		sr.Post("/{id}/contract", d.contractProject)
		sr.Post("/{id}/transfers", d.createTransfer)
	})
	r.With(AuthMiddleware(d.Auth)).Post("/v1/reviews", d.assignReview)
	return r
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, e error) {
	status := http.StatusInternalServerError
	code := "internal_error"
	if e == domain.ErrUnauthorized {
		status = http.StatusUnauthorized
		code = "unauthorized"
	}
	if e == domain.ErrForbidden {
		status = http.StatusForbidden
		code = "forbidden"
	}
	if e == domain.ErrConflict {
		status = http.StatusConflict
		code = "conflict"
	}
	if e == domain.ErrInvalid {
		status = http.StatusBadRequest
		code = "invalid"
	}
	writeJSON(w, status, map[string]string{"code": code, "message": e.Error()})
}
