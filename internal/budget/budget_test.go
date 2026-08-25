package budget

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/clock"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"testing"
	"time"
)

func setupBudget(t *testing.T) (context.Context, *Service, int64) {
	ctx := context.Background()
	d, e := db.Open(ctx, "file:budget-"+t.Name()+"?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	db.Migrate(ctx, d)
	u := repository.Users{DB: d.DB}
	uid, _ := u.Create(ctx, domain.User{Name: "u", Email: t.Name() + "@x", Role: domain.RoleResearcher, PasswordHash: "p", CreatedAt: time.Now()})
	p := repository.Projects{DB: d.DB}
	id, e := p.Create(ctx, domain.Project{Title: "b", OwnerID: uid, BudgetCents: 1000, DueAt: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now(), State: domain.ProjectActive})
	if e != nil {
		t.Fatal(e)
	}
	return ctx, &Service{DB: d.DB}, id
}
func TestReserveAndRemaining(t *testing.T) {
	ctx, s, id := setupBudget(t)
	if e := s.Reserve(ctx, id, 300, 1); e != nil {
		t.Fatal(e)
	}
	left, e := s.Remaining(ctx, id)
	if e != nil || left != 700 {
		t.Fatalf("%d %v", left, e)
	}
}
func TestReserveVersionConflict(t *testing.T) {
	ctx, s, id := setupBudget(t)
	if e := s.Reserve(ctx, id, 300, 99); e == nil {
		t.Fatal("stale version")
	}
}
func TestReserveExhausted(t *testing.T) {
	ctx, s, id := setupBudget(t)
	if e := s.Reserve(ctx, id, 2000, 1); e == nil {
		t.Fatal("overspend")
	}
}
func TestRelease(t *testing.T) {
	ctx, s, id := setupBudget(t)
	s.Reserve(ctx, id, 500, 1)
	if e := s.Release(ctx, id, 200); e != nil {
		t.Fatal(e)
	}
	left, _ := s.Remaining(ctx, id)
	if left != 700 {
		t.Fatal(left)
	}
}
func TestReserveInvalid(t *testing.T) {
	ctx, s, id := setupBudget(t)
	if e := s.Reserve(ctx, id, 0, 1); e == nil {
		t.Fatal("zero")
	}
}

var _ = clock.Real{}
