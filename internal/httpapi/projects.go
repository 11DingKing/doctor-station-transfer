package httpapi

import (
	"encoding/json"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/pagination"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

func (d RouterDeps) createProject(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	var in struct {
		Title, Summary string
		BudgetCents    int64
		DueAt          string
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, domain.ErrInvalid)
		return
	}
	due, _ := time.Parse(time.RFC3339, in.DueAt)
	p, e := d.Projects.Create(r.Context(), domain.Project{Title: in.Title, Summary: in.Summary, BudgetCents: in.BudgetCents, DueAt: due, OwnerID: u.ID}, r.Header.Get("X-Request-ID"))
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}
func (d RouterDeps) listProjects(w http.ResponseWriter, r *http.Request) {
	q := pagination.Parse(r.URL.Query().Get("limit"), r.URL.Query().Get("offset"))
	q.Query = r.URL.Query().Get("q")
	v, e := d.Projects.List(r.Context(), q)
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, http.StatusOK, v)
}
func (d RouterDeps) transitionProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var in struct {
		To      string
		Version int64
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, domain.ErrInvalid)
		return
	}
	e := d.Projects.Transition(r.Context(), id, domain.ProjectState(in.To), in.Version, userFrom(r).ID, r.Header.Get("X-Request-ID"))
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
func (d RouterDeps) contractProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var in struct {
		Number      string
		AmountCents int64
	}
	json.NewDecoder(r.Body).Decode(&in)
	e := d.Projects.CreateContract(r.Context(), id, domain.Contract{ProjectID: id, Number: in.Number, AmountCents: in.AmountCents}, userFrom(r).ID, r.Header.Get("X-Request-ID"))
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "contracted"})
}
func (d RouterDeps) createTransfer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var in struct{ Kind, ArtifactRef, Checksum string }
	json.NewDecoder(r.Body).Decode(&in)
	tid, e := d.Transfers.Record(r.Context(), domain.Transfer{ProjectID: id, ActorID: userFrom(r).ID, Kind: in.Kind, ArtifactRef: in.ArtifactRef, Checksum: in.Checksum}, r.Header.Get("X-Request-ID"))
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": tid})
}
