package httpapi

import (
	"encoding/json"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"net/http"
)

func (d RouterDeps) assignReview(w http.ResponseWriter, r *http.Request) {
	var in domain.Review
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, domain.ErrInvalid)
		return
	}
	id, e := d.Reviews.Assign(r.Context(), in, userFrom(r).ID, r.Header.Get("X-Request-ID"))
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}
