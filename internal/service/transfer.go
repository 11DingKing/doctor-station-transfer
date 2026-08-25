package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/doctor-station-transfer/internal/clock"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
)

type Transfers struct {
	DB       *sql.DB
	Repo     repository.Transfers
	Projects repository.Projects
	Audits   repository.Audits
	Clock    clock.Clock
}

func (s Transfers) Record(ctx context.Context, t domain.Transfer, requestID string) (int64, error) {
	p, e := s.Projects.Get(ctx, t.ProjectID)
	if e != nil {
		return 0, e
	}
	if p.State != domain.ProjectActive {
		return 0, errors.New("project is not active")
	}
	if t.ArtifactRef == "" || t.Checksum == "" {
		return 0, domain.ErrInvalid
	}
	tx, e := s.DB.BeginTx(ctx, nil)
	if e != nil {
		return 0, e
	}
	x, e := tx.ExecContext(ctx, `INSERT INTO transfers(project_id,actor_id,kind,artifact_ref,checksum,version,created_at) VALUES(?,?,?,?,?,1,datetime('now'))`, t.ProjectID, t.ActorID, t.Kind, t.ArtifactRef, t.Checksum)
	if e != nil {
		tx.Rollback()
		return 0, e
	}
	id, _ := x.LastInsertId()
	if _, e = tx.ExecContext(ctx, `INSERT INTO audits(actor_id,entity_type,entity_id,action,result,request_id,created_at) VALUES(?,?,?,?,?,?,datetime('now'))`, t.ActorID, "transfer", id, "record", "ok", requestID); e != nil {
		tx.Rollback()
		return 0, e
	}
	if e = tx.Commit(); e != nil {
		return 0, e
	}
	return id, nil
}

var _ = clock.Real{}
