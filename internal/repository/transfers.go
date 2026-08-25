package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"time"
)

type Transfers struct{ DB *sql.DB }

func (r Transfers) Create(ctx context.Context, t domain.Transfer) (int64, error) {
	x, e := r.DB.ExecContext(ctx, `INSERT INTO transfers(project_id,actor_id,kind,artifact_ref,checksum,version,created_at) VALUES(?,?,?,?,?,1,?)`, t.ProjectID, t.ActorID, t.Kind, t.ArtifactRef, t.Checksum, t.CreatedAt.Format(time.RFC3339Nano))
	if e != nil {
		return 0, e
	}
	id, e := x.LastInsertId()
	return id, e
}
