package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/11DingKing/doctor-station-transfer/internal/clock"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/pagination"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"time"
)

type Projects struct {
	DB         *sql.DB
	Repo       repository.Projects
	Contracts  repository.Contracts
	Milestones repository.Milestones
	Audits     repository.Audits
	Clock      clock.Clock
}

func (s Projects) Create(ctx context.Context, p domain.Project, requestID string) (domain.Project, error) {
	if p.BudgetCents <= 0 || p.Title == "" {
		return p, domain.ErrInvalid
	}
	now := s.Clock.Now()
	p.State = domain.ProjectDraft
	p.CreatedAt = now
	p.UpdatedAt = now
	tx, e := s.DB.BeginTx(ctx, nil)
	if e != nil {
		return p, e
	}
	x, e := tx.ExecContext(ctx, `INSERT INTO projects(title,summary,owner_id,budget_cents,spent_cents,state,version,due_at,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, p.Title, p.Summary, p.OwnerID, p.BudgetCents, 0, p.State, 1, p.DueAt.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if e != nil {
		tx.Rollback()
		return p, e
	}
	p.ID, _ = x.LastInsertId()
	if _, e = tx.ExecContext(ctx, `INSERT INTO audits(actor_id,entity_type,entity_id,action,result,request_id,created_at) VALUES(?,?,?,?,?,?,?)`, p.OwnerID, "project", p.ID, "create", "ok", requestID, now.Format(time.RFC3339Nano)); e != nil {
		tx.Rollback()
		return p, e
	}
	if e = tx.Commit(); e != nil {
		return p, e
	}
	return p, nil
}
func (s Projects) Transition(ctx context.Context, id int64, to domain.ProjectState, version int64, actor int64, requestID string) error {
	p, e := s.Repo.Get(ctx, id)
	if e != nil {
		return e
	}
	if e = p.State.ValidateMove(to); e != nil {
		return e
	}
	if e = s.Repo.Transition(ctx, id, p.State, to, version); e != nil {
		return e
	}
	return s.Audits.Add(ctx, domain.AuditEvent{ActorID: actor, EntityType: "project", EntityID: id, Action: "transition", Result: string(to), RequestID: requestID, CreatedAt: s.Clock.Now()})
}
func (s Projects) List(ctx context.Context, q pagination.Request) (pagination.Result[domain.Project], error) {
	return s.Repo.List(ctx, q)
}
func (s Projects) CreateContract(ctx context.Context, pid int64, c domain.Contract, actor int64, requestID string) error {
	tx, e := s.DB.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	p, e := s.Repo.Get(ctx, pid)
	if e != nil {
		tx.Rollback()
		return e
	}
	if p.State != domain.ProjectApproved {
		tx.Rollback()
		return errors.New("project must be approved")
	}
	if _, e = tx.ExecContext(ctx, `INSERT INTO contracts(project_id,number,amount_cents,state,created_at) VALUES(?,?,?,?,?)`, pid, c.Number, c.AmountCents, "draft", s.Clock.Now().Format(time.RFC3339Nano)); e != nil {
		tx.Rollback()
		return e
	}
	if _, e = tx.ExecContext(ctx, `UPDATE projects SET state='contracted',version=version+1 WHERE id=? AND state='approved'`, pid); e != nil {
		tx.Rollback()
		return e
	}
	if _, e = tx.ExecContext(ctx, `INSERT INTO audits(actor_id,entity_type,entity_id,action,result,request_id,created_at) VALUES(?,?,?,?,?,?,?)`, actor, "project", pid, "contract", "ok", requestID, s.Clock.Now().Format(time.RFC3339Nano)); e != nil {
		tx.Rollback()
		return e
	}
	return tx.Commit()
}
func Encode(v any) string { b, _ := json.Marshal(v); return string(b) }

var _ = fmt.Sprintf
