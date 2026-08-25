package policy

import (
	"fmt"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
)

func CanCreateProject(r domain.Role) bool { return r == domain.RoleAdmin || r == domain.RoleResearcher }
func CanReview(r domain.Role) bool        { return r == domain.RoleReviewer || r == domain.RoleAdmin }
func CanAudit(r domain.Role) bool         { return r == domain.RoleAuditor || r == domain.RoleAdmin }
func CanTransfer(r domain.Role) bool      { return r == domain.RoleResearcher || r == domain.RoleAdmin }
func Require(r domain.Role, ok bool) error {
	if !ok {
		return fmt.Errorf("role %s is not permitted", r)
	}
	return nil
}
