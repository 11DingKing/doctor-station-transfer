package milestone

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"testing"
	"time"
)

func mileFixture(t *testing.T) (context.Context, *Service, int64) {
	ctx := context.Background()
	d, _ := db.Open(ctx, "file:mile-"+t.Name()+"?mode=memory&cache=shared")
	db.Migrate(ctx, d)
	u := repository.Users{DB: d.DB}
	uid, _ := u.Create(ctx, domain.User{Name: "u", Email: t.Name() + "@x", Role: domain.RoleResearcher, PasswordHash: "p", CreatedAt: time.Now()})
	p := repository.Projects{DB: d.DB}
	id, _ := p.Create(ctx, domain.Project{Title: "p", OwnerID: uid, BudgetCents: 1, DueAt: time.Now().Add(time.Hour), CreatedAt: time.Now(), UpdatedAt: time.Now()})
	return ctx, &Service{DB: d.DB}, id
}
func TestAddAndList(t *testing.T) {
	ctx, s, pid := mileFixture(t)
	if e := s.Add(ctx, domain.Milestone{ProjectID: pid, Name: "m", Sequence: 1, DueAt: time.Now().Add(time.Hour)}); e != nil {
		t.Fatal(e)
	}
	items, e := s.List(ctx, pid)
	if e != nil || len(items) != 1 {
		t.Fatal(items, e)
	}
}
func TestAddRejectsPast(t *testing.T) {
	ctx, s, pid := mileFixture(t)
	if e := s.Add(ctx, domain.Milestone{ProjectID: pid, Name: "m", Sequence: 1, DueAt: time.Now().Add(-time.Hour)}); e == nil {
		t.Fatal("past")
	}
}
func TestCompleteMissing(t *testing.T) {
	ctx, s, _ := mileFixture(t)
	if e := s.Complete(ctx, 9); e == nil {
		t.Fatal("missing")
	}
}
