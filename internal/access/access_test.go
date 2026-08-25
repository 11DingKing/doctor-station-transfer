package access

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"testing"
)

func TestCheckerProject(t *testing.T) {
	c := Checker{}
	p := domain.Project{OwnerID: 2}
	if c.ProjectWrite(context.Background(), domain.User{ID: 2, Role: domain.RoleResearcher}, p) != nil {
		t.Fatal("owner")
	}
	if c.ProjectWrite(context.Background(), domain.User{ID: 3, Role: domain.RoleResearcher}, p) == nil {
		t.Fatal("other")
	}
	if c.ProjectWrite(context.Background(), domain.User{ID: 3, Role: domain.RoleAdmin}, p) != nil {
		t.Fatal("admin")
	}
}
func TestCheckerReview(t *testing.T) {
	c := Checker{}
	if c.ReviewWrite(context.Background(), domain.User{Role: domain.RoleReviewer}) != nil {
		t.Fatal("review")
	}
	if c.ReviewWrite(context.Background(), domain.User{Role: domain.RoleAuditor}) == nil {
		t.Fatal("audit")
	}
}
func TestCheckerAudit(t *testing.T) {
	c := Checker{}
	if c.AuditRead(context.Background(), domain.User{Role: domain.RoleAuditor}) != nil {
		t.Fatal("audit")
	}
	if c.AuditRead(context.Background(), domain.User{Role: domain.RoleResearcher}) == nil {
		t.Fatal("research")
	}
}
func TestCheckerTransfer(t *testing.T) {
	c := Checker{}
	if c.TransferWrite(context.Background(), domain.User{Role: domain.RoleResearcher}) != nil {
		t.Fatal("transfer")
	}
	if c.TransferWrite(context.Background(), domain.User{Role: domain.RoleReviewer}) == nil {
		t.Fatal("review")
	}
}
