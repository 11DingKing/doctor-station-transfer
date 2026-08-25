package policy

import (
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"testing"
)

func TestRoleMatrix(t *testing.T) {
	cases := []struct {
		r                               domain.Role
		create, review, audit, transfer bool
	}{{domain.RoleAdmin, true, true, true, true}, {domain.RoleResearcher, true, false, false, true}, {domain.RoleReviewer, false, true, false, false}, {domain.RoleAuditor, false, false, true, false}, {domain.Role("guest"), false, false, false, false}}
	for _, c := range cases {
		if CanCreateProject(c.r) != c.create || CanReview(c.r) != c.review || CanAudit(c.r) != c.audit || CanTransfer(c.r) != c.transfer {
			t.Errorf("matrix %+v", c)
		}
	}
}
func TestRequire(t *testing.T) {
	if Require(domain.RoleAdmin, true) != nil {
		t.Fatal("true rejected")
	}
	if Require(domain.RoleReviewer, false) == nil {
		t.Fatal("false accepted")
	}
}
func TestRoleStrings(t *testing.T) {
	for _, r := range []domain.Role{domain.RoleAdmin, domain.RoleResearcher, domain.RoleReviewer, domain.RoleAuditor} {
		if r == "" {
			t.Fatal("empty")
		}
	}
}
