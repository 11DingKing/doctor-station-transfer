package httpapi

import (
	"encoding/json"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"net/http"
)

func (d RouterDeps) login(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, domain.ErrInvalid)
		return
	}
	u, t, e := d.Auth.Login(r.Context(), in.Email, in.Password)
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": t, "user": u})
}
func (d RouterDeps) logout(w http.ResponseWriter, r *http.Request) {
	h := r.Header.Get("Authorization")
	if e := d.Auth.Logout(r.Context(), h[len("Bearer "):]); e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}
