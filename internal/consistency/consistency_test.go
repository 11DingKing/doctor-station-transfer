package consistency

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"testing"
	"time"
)

func consistencyFixture(t *testing.T) (context.Context, *Checker, int64) {
	ctx := context.Background()
	d, _ := db.Open(ctx, "file:cons-"+t.Name()+"?mode=memory&cache=shared")
	db.Migrate(ctx, d)
	u := repository.Users{DB: d.DB}
	uid, _ := u.Create(ctx, domain.User{Name: "u", Email: t.Name() + "@x", Role: domain.RoleResearcher, PasswordHash: "p", CreatedAt: time.Now()})
	p := repository.Projects{DB: d.DB}
	id, _ := p.Create(ctx, domain.Project{Title: "p", OwnerID: uid, BudgetCents: 100, DueAt: time.Now().Add(time.Hour), CreatedAt: time.Now(), UpdatedAt: time.Now()})
	return ctx, &Checker{DB: d.DB}, id
}
func TestProjectConsistency(t *testing.T) {
	ctx, c, id := consistencyFixture(t)
	if e := c.Project(ctx, id); e != nil {
		t.Fatal(e)
	}
}
func TestReviewConsistency(t *testing.T) {
	ctx, c, id := consistencyFixture(t)
	if e := c.Reviews(ctx, id); e != nil {
		t.Fatal(e)
	}
}
func TestForeignKeyConsistency(t *testing.T) {
	ctx, c, _ := consistencyFixture(t)
	if e := c.ForeignKeys(ctx); e != nil {
		t.Fatal(e)
	}
}
