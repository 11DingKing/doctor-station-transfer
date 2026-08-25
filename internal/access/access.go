package access

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/policy"
)

type Checker struct{}

func (Checker) ProjectWrite(ctx context.Context, u domain.User, p domain.Project) error {
	if err := policy.Require(u.Role, policy.CanCreateProject(u.Role)); err != nil {
		return domain.ErrForbidden
	}
	if p.OwnerID != u.ID && u.Role != domain.RoleAdmin {
		return domain.ErrForbidden
	}
	return nil
}
func (Checker) ReviewWrite(ctx context.Context, u domain.User) error {
	return policy.Require(u.Role, policy.CanReview(u.Role))
}
func (Checker) AuditRead(ctx context.Context, u domain.User) error {
	return policy.Require(u.Role, policy.CanAudit(u.Role))
}
func (Checker) TransferWrite(ctx context.Context, u domain.User) error {
	return policy.Require(u.Role, policy.CanTransfer(u.Role))
}
