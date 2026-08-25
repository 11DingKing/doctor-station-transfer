package snapshot

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"testing"
	"time"
)

func TestSnapshotCapture(t *testing.T) {
	ctx := context.Background()
	d, _ := db.Open(ctx, "file:snap-"+t.Name()+"?mode=memory&cache=shared")
	defer d.Close()
	db.Migrate(ctx, d)
	u := repository.Users{DB: d.DB}
	uid, _ := u.Create(ctx, domain.User{Name: "u", Email: t.Name() + "@x", Role: domain.RoleResearcher, PasswordHash: "p", CreatedAt: time.Now()})
	p := repository.Projects{DB: d.DB}
	id, _ := p.Create(ctx, domain.Project{Title: "p", OwnerID: uid, BudgetCents: 1, DueAt: time.Now().Add(time.Hour), CreatedAt: time.Now(), UpdatedAt: time.Now()})
	s := Store{DB: d.DB}
	a, e := s.Capture(ctx, id)
	if e != nil || a.ProjectID != id {
		t.Fatal(a, e)
	}
	if !s.Equal(a, a) {
		t.Fatal("equal")
	}
	b := a
	b.Version++
	if s.Equal(a, b) {
		t.Fatal("different")
	}
	if !s.Changed(a, b) {
		t.Fatal("changed")
	}
	if s.Version(a) != a.Version {
		t.Fatal("version")
	}
}
