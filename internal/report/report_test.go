package report

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"testing"
	"time"
)

func reportFixture(t *testing.T) (context.Context, *Service, int64) {
	ctx := context.Background()
	d, e := db.Open(ctx, "file:report-"+t.Name()+"?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	db.Migrate(ctx, d)
	u := repository.Users{DB: d.DB}
	uid, _ := u.Create(ctx, domain.User{Name: "u", Email: t.Name() + "@x", Role: domain.RoleResearcher, PasswordHash: "p", CreatedAt: time.Now()})
	p := repository.Projects{DB: d.DB}
	id, _ := p.Create(ctx, domain.Project{Title: "课题", OwnerID: uid, BudgetCents: 1000, DueAt: time.Now().Add(time.Hour), CreatedAt: time.Now(), UpdatedAt: time.Now()})
	return ctx, &Service{DB: d.DB}, id
}
func TestBuildEmpty(t *testing.T) {
	ctx, s, id := reportFixture(t)
	x, e := s.Build(ctx, id)
	if e != nil || x.ProjectID != id || x.Reviews != 0 {
		t.Fatal(x, e)
	}
}
func TestExportContainsTotals(t *testing.T) {
	ctx, s, id := reportFixture(t)
	v, e := s.Export(ctx, id)
	if e != nil || v == "" {
		t.Fatal(v, e)
	}
}
func TestMissingProject(t *testing.T) {
	ctx, s, _ := reportFixture(t)
	if _, e := s.Build(ctx, 999); e == nil {
		t.Fatal("missing")
	}
}
