package audit

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
)

type Service struct{ Repo repository.Audits }

func (s Service) Record(ctx context.Context, a domain.AuditEvent) error {
	if a.EntityType == "" || a.Action == "" {
		return domain.ErrInvalid
	}
	return s.Repo.Add(ctx, a)
}
func (s Service) History(ctx context.Context, typ string, id int64) ([]domain.AuditEvent, error) {
	return s.Repo.ForEntity(ctx, typ, id)
}
