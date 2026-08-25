package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"testing"
	"time"
)

func repoDB(t *testing.T) (context.Context, *db.DB, int64) {
	ctx := context.Background()
	d, e := db.Open(ctx, "file:repo-"+t.Name()+"?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	db.Migrate(ctx, d)
	u := Users{DB: d.DB}
	id, e := u.Create(ctx, domain.User{Name: "owner", Email: t.Name() + "@x", Role: domain.RoleResearcher, PasswordHash: "p", CreatedAt: time.Now()})
	if e != nil {
		t.Fatal(e)
	}
	return ctx, d, id
}
func makeProject(t *testing.T, ctx context.Context, d *db.DB, uid int64) int64 {
	p := Projects{DB: d.DB}
	id, e := p.Create(ctx, domain.Project{Title: "研究", Summary: "样品", OwnerID: uid, BudgetCents: 5000, State: domain.ProjectDraft, Version: 1, DueAt: time.Now().Add(time.Hour), CreatedAt: time.Now(), UpdatedAt: time.Now()})
	if e != nil {
		t.Fatal(e)
	}
	return id
}
func TestUsersCRUD(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	u := Users{DB: d.DB}
	got, e := u.ByID(ctx, id)
	if e != nil || got.Email == "" {
		t.Fatal(got, e)
	}
	_, e = u.ByEmail(ctx, got.Email)
	if e != nil {
		t.Fatal(e)
	}
}
func TestSessionRepository(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	s := Sessions{DB: d.DB}
	x := domain.Session{UserID: id, TokenHash: "hash", ExpiresAt: time.Now().Add(time.Hour), CreatedAt: time.Now()}
	if e := s.Create(ctx, x); e != nil {
		t.Fatal(e)
	}
	got, e := s.Find(ctx, "hash")
	if e != nil || got.UserID != id {
		t.Fatal(got, e)
	}
	if e = s.Revoke(ctx, "hash"); e != nil {
		t.Fatal(e)
	}
	got, e = s.Find(ctx, "hash")
	if e != nil || got.RevokedAt == nil {
		t.Fatal(got, e)
	}
}
func TestProjectsGetAndTransition(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	pid := makeProject(t, ctx, d, id)
	p := Projects{DB: d.DB}
	got, e := p.Get(ctx, pid)
	if e != nil || got.State != domain.ProjectDraft {
		t.Fatal(got, e)
	}
	if e = p.Transition(ctx, pid, domain.ProjectDraft, domain.ProjectSubmitted, 1); e != nil {
		t.Fatal(e)
	}
	if _, e = p.Get(ctx, pid); e != nil {
		t.Fatal(e)
	}
}
func TestProjectConflict(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	pid := makeProject(t, ctx, d, id)
	p := Projects{DB: d.DB}
	if e := p.Transition(ctx, pid, domain.ProjectDraft, domain.ProjectSubmitted, 2); e == nil {
		t.Fatal("stale")
	}
}
func TestProjectListPagination(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	p := Projects{DB: d.DB}
	for i := 0; i < 7; i++ {
		p.Create(ctx, domain.Project{Title: "课题", OwnerID: id, BudgetCents: 100, DueAt: time.Now().Add(time.Hour), CreatedAt: time.Now(), UpdatedAt: time.Now()})
	}
	r, e := p.List(ctx, struct {
		Limit, Offset int
		Query, Sort   string
	}{Limit: 3, Offset: 2})
	if e != nil || len(r.Items) != 3 {
		t.Fatalf("%+v %v", r, e)
	}
}
func TestReviewRepository(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	pid := makeProject(t, ctx, d, id)
	r := Reviews{DB: d.DB}
	rid, e := r.Assign(ctx, domain.Review{ProjectID: pid, ReviewerID: id, State: domain.ReviewPending, DueAt: time.Now().Add(time.Hour), CreatedAt: time.Now()})
	if e != nil {
		t.Fatal(e)
	}
	if e = r.Transition(ctx, rid, domain.ReviewPending, domain.ReviewAccepted, 1, 80, "ok"); e != nil {
		t.Fatal(e)
	}
	items, e := r.ForProject(ctx, pid)
	if e != nil || len(items) != 1 {
		t.Fatal(items, e)
	}
}
func TestReviewUnique(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	pid := makeProject(t, ctx, d, id)
	r := Reviews{DB: d.DB}
	x := domain.Review{ProjectID: pid, ReviewerID: id, State: domain.ReviewPending, DueAt: time.Now().Add(time.Hour), CreatedAt: time.Now()}
	r.Assign(ctx, x)
	if _, e := r.Assign(ctx, x); e == nil {
		t.Fatal("duplicate")
	}
}
func TestAuditRepository(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	a := Audits{DB: d.DB}
	if e := a.Add(ctx, domain.AuditEvent{ActorID: id, EntityType: "project", EntityID: 1, Action: "create", Result: "ok", RequestID: "r", CreatedAt: time.Now()}); e != nil {
		t.Fatal(e)
	}
	items, e := a.ForEntity(ctx, "project", 1)
	if e != nil || len(items) != 1 {
		t.Fatal(items, e)
	}
}
func TestContractRepository(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	pid := makeProject(t, ctx, d, id)
	c := Contracts{DB: d.DB}
	if e := c.Create(ctx, domain.Contract{ProjectID: pid, Number: "N1", AmountCents: 100, State: "draft", CreatedAt: time.Now()}); e != nil {
		t.Fatal(e)
	}
	if e := c.Sign(ctx, 1); e != nil {
		t.Fatal(e)
	}
}
func TestMilestoneRepository(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	pid := makeProject(t, ctx, d, id)
	m := Milestones{DB: d.DB}
	if e := m.Create(ctx, domain.Milestone{ProjectID: pid, Name: "阶段一", DueAt: time.Now().Add(-time.Hour), State: domain.MilestonePending, Sequence: 1}); e != nil {
		t.Fatal(e)
	}
	n, e := m.Overdue(ctx, time.Now())
	if e != nil || n != 1 {
		t.Fatal(n, e)
	}
}
func TestTransferRepository(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	pid := makeProject(t, ctx, d, id)
	r := Transfers{DB: d.DB}
	n, e := r.Create(ctx, domain.Transfer{ProjectID: pid, ActorID: id, Kind: "sample", ArtifactRef: "s", Checksum: "c", CreatedAt: time.Now()})
	if e != nil || n == 0 {
		t.Fatal(n, e)
	}
}
func TestJobRepository(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	_ = id
	j := Jobs{DB: d.DB}
	if e := j.Enqueue(ctx, domain.Job{Kind: "sweep", Payload: "{}", RunAt: time.Now().Add(-time.Second), CreatedAt: time.Now(), UpdatedAt: time.Now()}); e != nil {
		t.Fatal(e)
	}
	x, e := j.Claim(ctx, time.Now())
	if e != nil || x.ID == 0 {
		t.Fatal(x, e)
	}
	if e = j.Finish(ctx, x.ID, nil); e != nil {
		t.Fatal(e)
	}
}
func TestIdempotencyRepository(t *testing.T) {
	ctx, d, _ := repoDB(t)
	defer d.Close()
	r := Idempotency{DB: d.DB}
	if e := r.Put(ctx, "k", "create", "{}"); e != nil {
		t.Fatal(e)
	}
	v, e := r.Get(ctx, "k", "create")
	if e != nil || v != "{}" {
		t.Fatal(v, e)
	}
}
func TestExecTxCommit(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	e := ExecTx(ctx, d.DB, func(tx *sql.Tx) error {
		_, e := tx.ExecContext(ctx, `UPDATE users SET active=0 WHERE id=?`, id)
		return e
	})
	if e != nil {
		t.Fatal(e)
	}
}
func TestExecTxRollback(t *testing.T) {
	ctx, d, id := repoDB(t)
	defer d.Close()
	e := ExecTx(ctx, d.DB, func(tx *sql.Tx) error {
		tx.ExecContext(ctx, `UPDATE users SET active=0 WHERE id=?`, id)
		return errors.New("fail")
	})
	if e == nil {
		t.Fatal("rollback")
	}
	u := Users{DB: d.DB}
	x, _ := u.ByID(ctx, id)
	if !x.Active {
		t.Fatal("state leaked")
	}
}
