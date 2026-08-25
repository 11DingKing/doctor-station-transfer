package ledger

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"testing"
	"time"
)

func ledgerFixture(t *testing.T) (context.Context, *Store, int64, int64) {
	ctx := context.Background()
	d, _ := db.Open(ctx, "file:ledger-"+t.Name()+"?mode=memory&cache=shared")
	db.Migrate(ctx, d)
	u := repository.Users{DB: d.DB}
	uid, _ := u.Create(ctx, domain.User{Name: "u", Email: t.Name() + "@x", Role: domain.RoleResearcher, PasswordHash: "p", CreatedAt: time.Now()})
	p := repository.Projects{DB: d.DB}
	pid, _ := p.Create(ctx, domain.Project{Title: "p", OwnerID: uid, BudgetCents: 1, DueAt: time.Now().Add(time.Hour), CreatedAt: time.Now(), UpdatedAt: time.Now()})
	return ctx, &Store{DB: d.DB}, pid, uid
}
func TestLedgerAppendValidation(t *testing.T) {
	ctx, s, pid, uid := ledgerFixture(t)
	if e := s.Append(ctx, Entry{ProjectID: pid, ActorID: uid, Amount: 1, Kind: "credit", Reference: "r", CreatedAt: time.Now()}); e != nil {
		t.Fatal(e)
	}
	if e := s.Append(ctx, Entry{ProjectID: 0, ActorID: uid, Amount: 1, Reference: "r"}); e == nil {
		t.Fatal("invalid")
	}
}
func TestLedgerRecent(t *testing.T) {
	ctx, s, pid, uid := ledgerFixture(t)
	for i := 0; i < 3; i++ {
		s.Append(ctx, Entry{ProjectID: pid, ActorID: uid, Amount: 1, Kind: "credit", Reference: "r", CreatedAt: time.Now()})
	}
	items, e := s.Recent(ctx, pid, 2)
	if e != nil || len(items) != 2 {
		t.Fatal(items, e)
	}
}
func TestLedgerTotal(t *testing.T) {
	ctx, s, pid, uid := ledgerFixture(t)
	s.Append(ctx, Entry{ProjectID: pid, ActorID: uid, Amount: 1, Kind: "credit", Reference: "r", CreatedAt: time.Now()})
	n, e := s.Total(ctx, pid)
	if e != nil || n != 1 {
		t.Fatal(n, e)
	}
}
