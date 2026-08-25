package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/doctor-station-transfer/internal/clock"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"time"
)

type Reviews struct {
	DB       *sql.DB
	Repo     repository.Reviews
	Projects repository.Projects
	Audits   repository.Audits
	Clock    clock.Clock
}

func (s Reviews) Assign(ctx context.Context, r domain.Review, actor int64, requestID string) (int64, error) {
	if r.ReviewerID == actor {
		return 0, domain.ErrForbidden
	}
	r.State = domain.ReviewPending
	r.CreatedAt = s.Clock.Now()
	id, e := s.Repo.Assign(ctx, r)
	if e != nil {
		return 0, e
	}
	e = s.Audits.Add(ctx, domain.AuditEvent{ActorID: actor, EntityType: "review", EntityID: id, Action: "assign", Result: "ok", RequestID: requestID, CreatedAt: s.Clock.Now()})
	return id, e
}
func (s Reviews) Submit(ctx context.Context, id int64, version int64, score int, comment string, actor int64, requestID string) error {
	if score < 1 || score > 100 {
		return domain.ErrInvalid
	}
	if e := s.Repo.Transition(ctx, id, domain.ReviewAccepted, domain.ReviewSubmitted, version, score, comment); e != nil {
		return e
	}
	return s.Audits.Add(ctx, domain.AuditEvent{ActorID: actor, EntityType: "review", EntityID: id, Action: "submit", Result: "ok", RequestID: requestID, CreatedAt: s.Clock.Now()})
}
func (s Reviews) Expire(ctx context.Context, now time.Time) error {
	_, e := s.DB.ExecContext(ctx, `UPDATE reviews SET state='recused' WHERE state='pending' AND due_at<?`, now.Format(time.RFC3339Nano))
	return e
}
