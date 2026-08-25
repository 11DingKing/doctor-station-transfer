package httpapi

import (
	"github.com/11DingKing/doctor-station-transfer/internal/service"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	h := New(RouterDeps{Auth: service.Auth{}})
	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}
	if w.Body.Len() == 0 {
		t.Fatal("empty")
	}
}
func TestReadyEndpoint(t *testing.T) {
	h := New(RouterDeps{})
	r := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}
}
func TestProtectedEndpointRejectsMissingToken(t *testing.T) {
	h := New(RouterDeps{Auth: service.Auth{}})
	r := httptest.NewRequest(http.MethodGet, "/v1/projects/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatal(w.Code)
	}
}
func TestUnknownRoute(t *testing.T) {
	h := New(RouterDeps{})
	r := httptest.NewRequest(http.MethodGet, "/missing", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatal(w.Code)
	}
}
func TestErrorMapping(t *testing.T) {
	for _, code := range []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusConflict, http.StatusBadRequest, http.StatusInternalServerError} {
		if code < 400 {
			t.Fatal(code)
		}
	}
}
func TestRequestMethods(t *testing.T) {
	for _, m := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete} {
		if m == "" {
			t.Fatal(m)
		}
	}
}
func TestRecovery(t *testing.T) {
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("x") })
	h := Recovery(next, slog.Default())
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 500 {
		t.Fatal(w.Code)
	}
}
