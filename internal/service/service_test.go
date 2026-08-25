package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/doctor-station-transfer/internal/clock"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/pagination"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"testing"
	"time"
)

func fixture(t *testing.T) (context.Context, *sql.DB, clock.Fixed, repository.Users) {
	t.Helper()
	ctx := context.Background()
	d, e := db.Open(ctx, "file:svc-"+t.Name()+"?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	if e = db.Migrate(ctx, d); e != nil {
		t.Fatal(e)
	}
	c := clock.Fixed{T: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)}
	u := repository.Users{DB: d.DB}
	if _, e = u.Create(ctx, domain.User{Name: "Researcher", Email: "r-" + t.Name() + "@x", Role: domain.RoleResearcher, PasswordHash: "p", CreatedAt: c.Now()}); e != nil {
		t.Fatal(e)
	}
	return ctx, d.DB, c, u
}
func TestCreateProjectWritesAudit(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	p := Projects{DB: raw, Repo: repository.Projects{DB: raw}, Contracts: repository.Contracts{DB: raw}, Milestones: repository.Milestones{DB: raw}, Audits: repository.Audits{DB: raw}, Clock: c}
	v, e := p.Create(ctx, domain.Project{Title: "药材", Summary: "样品", OwnerID: x.ID, BudgetCents: 1000, DueAt: c.Now().Add(time.Hour)}, "req")
	if e != nil {
		t.Fatal(e)
	}
	if v.ID == 0 {
		t.Fatal("id")
	}
	var n int
	raw.QueryRow(`SELECT COUNT(*) FROM audits WHERE entity_id=?`, v.ID).Scan(&n)
	if n != 1 {
		t.Fatalf("audit %d", n)
	}
}
func TestCreateProjectRejectsBudget(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	p := Projects{DB: raw, Repo: repository.Projects{DB: raw}, Clock: c}
	if _, e := p.Create(ctx, domain.Project{Title: "bad", OwnerID: x.ID, BudgetCents: 0}, ""); !errors.Is(e, domain.ErrInvalid) {
		t.Fatalf("%v", e)
	}
}
func TestTransitionUsesVersion(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	p := Projects{DB: raw, Repo: repository.Projects{DB: raw}, Audits: repository.Audits{DB: raw}, Clock: c}
	v, _ := p.Create(ctx, domain.Project{Title: "a", OwnerID: x.ID, BudgetCents: 1, DueAt: c.Now()}, "")
	if e := p.Transition(ctx, v.ID, domain.ProjectSubmitted, 1, x.ID, ""); e != nil {
		t.Fatal(e)
	}
	if e := p.Transition(ctx, v.ID, domain.ProjectReviewing, 1, x.ID, ""); !errors.Is(e, domain.ErrConflict) {
		t.Fatalf("expected conflict %v", e)
	}
}
func TestTransitionRejectsIllegal(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	p := Projects{DB: raw, Repo: repository.Projects{DB: raw}, Audits: repository.Audits{DB: raw}, Clock: c}
	v, _ := p.Create(ctx, domain.Project{Title: "a", OwnerID: x.ID, BudgetCents: 1, DueAt: c.Now()}, "")
	if e := p.Transition(ctx, v.ID, domain.ProjectCompleted, 1, x.ID, ""); e == nil {
		t.Fatal("illegal transition accepted")
	}
}
func TestListFiltersTitle(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	p := Projects{DB: raw, Repo: repository.Projects{DB: raw}, Clock: c}
	for _, title := range []string{"黄芪", "党参", "黄精"} {
		p.Create(ctx, domain.Project{Title: title, OwnerID: x.ID, BudgetCents: 1, DueAt: c.Now()}, "")
	}
	r, e := p.List(ctx, pagination.Request{Limit: 2, Query: "黄"})
	if e != nil || len(r.Items) != 2 {
		t.Fatalf("%+v %v", r, e)
	}
}
func TestContractRequiresApproval(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	p := Projects{DB: raw, Repo: repository.Projects{DB: raw}, Contracts: repository.Contracts{DB: raw}, Clock: c}
	v, _ := p.Create(ctx, domain.Project{Title: "a", OwnerID: x.ID, BudgetCents: 2, DueAt: c.Now()}, "")
	if e := p.CreateContract(ctx, v.ID, domain.Contract{Number: "C1", AmountCents: 2}, x.ID, ""); e == nil {
		t.Fatal("contract before approval")
	}
}
func TestReviewSelfAssignmentForbidden(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	s := Reviews{DB: raw, Repo: repository.Reviews{DB: raw}, Audits: repository.Audits{DB: raw}, Clock: c}
	if _, e := s.Assign(ctx, domain.Review{ProjectID: 1, ReviewerID: x.ID}, x.ID, ""); !errors.Is(e, domain.ErrForbidden) {
		t.Fatal(e)
	}
}
func TestReviewScoreBounds(t *testing.T) {
	ctx, raw, c, _ := fixture(t)
	defer raw.Close()
	s := Reviews{DB: raw, Repo: repository.Reviews{DB: raw}, Audits: repository.Audits{DB: raw}, Clock: c}
	for _, score := range []int{-1, 0, 101} {
		if e := s.Submit(ctx, 1, 1, score, "", 1, ""); !errors.Is(e, domain.ErrInvalid) {
			t.Fatalf("score %d: %v", score, e)
		}
	}
}
func TestTransferRequiresActive(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	p := repository.Projects{DB: raw}
	v, _ := p.Create(ctx, domain.Project{Title: "a", OwnerID: x.ID, BudgetCents: 1, State: domain.ProjectDraft, DueAt: c.Now(), CreatedAt: c.Now(), UpdatedAt: c.Now()})
	s := Transfers{DB: raw, Repo: repository.Transfers{DB: raw}, Projects: p, Audits: repository.Audits{DB: raw}, Clock: c}
	if _, e := s.Record(ctx, domain.Transfer{ProjectID: v, ActorID: x.ID, ArtifactRef: "a", Checksum: "c"}, ""); e == nil {
		t.Fatal("inactive transfer")
	}
}
func TestTransferValidatesArtifacts(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	p := repository.Projects{DB: raw}
	v, _ := p.Create(ctx, domain.Project{Title: "a", OwnerID: x.ID, BudgetCents: 1, State: domain.ProjectActive, DueAt: c.Now(), CreatedAt: c.Now(), UpdatedAt: c.Now()})
	s := Transfers{DB: raw, Repo: repository.Transfers{DB: raw}, Projects: p, Audits: repository.Audits{DB: raw}, Clock: c}
	if _, e := s.Record(ctx, domain.Transfer{ProjectID: v, ActorID: x.ID}, ""); !errors.Is(e, domain.ErrInvalid) {
		t.Fatal(e)
	}
}
func TestHashPasswordDeterministic(t *testing.T) {
	if HashPassword("abc") != HashPassword("abc") {
		t.Fatal("hash")
	}
	if HashPassword("abc") == HashPassword("def") {
		t.Fatal("collision")
	}
}
func TestAuthUnknownUser(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	a := Auth{Users: u, Sessions: repository.Sessions{DB: raw}, Now: c.Now, Secret: "s"}
	if _, _, e := a.Login(ctx, "none", "x"); !errors.Is(e, domain.ErrUnauthorized) {
		t.Fatal(e)
	}
}
func TestAuthSessionLifecycle(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	raw.Exec(`UPDATE users SET password_hash=? WHERE id=?`, HashPassword("pw"), x.ID)
	a := Auth{Users: u, Sessions: repository.Sessions{DB: raw}, Now: c.Now, Secret: "s"}
	_, tok, e := a.Login(ctx, x.Email, "pw")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = a.Current(ctx, tok); e != nil {
		t.Fatal(e)
	}
	a.Logout(ctx, tok)
	if _, e = a.Current(ctx, tok); !errors.Is(e, domain.ErrUnauthorized) {
		t.Fatal(e)
	}
}
func TestAuthExpiry(t *testing.T) {
	ctx, raw, c, u := fixture(t)
	defer raw.Close()
	x, _ := u.ByEmail(ctx, "r-"+t.Name()+"@x")
	raw.Exec(`UPDATE users SET password_hash=? WHERE id=?`, HashPassword("pw"), x.ID)
	a := Auth{Users: u, Sessions: repository.Sessions{DB: raw}, Now: c.Now, Secret: "s"}
	_, tok, _ := a.Login(ctx, x.Email, "pw")
	raw.Exec(`UPDATE sessions SET expires_at=?`, c.Now().Add(-time.Minute).Format(time.RFC3339Nano))
	if _, e := a.Current(ctx, tok); !errors.Is(e, domain.ErrUnauthorized) {
		t.Fatal(e)
	}
}
