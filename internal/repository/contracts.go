package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"time"
)

type Contracts struct{ DB *sql.DB }

func (r Contracts) Create(ctx context.Context, c domain.Contract) error {
	_, e := r.DB.ExecContext(ctx, `INSERT INTO contracts(project_id,number,amount_cents,state,created_at) VALUES(?,?,?,?,?)`, c.ProjectID, c.Number, c.AmountCents, c.State, c.CreatedAt.Format(time.RFC3339Nano))
	return e
}
func (r Contracts) Sign(ctx context.Context, id int64) error {
	_, e := r.DB.ExecContext(ctx, `UPDATE contracts SET state='signed',signed_at=? WHERE id=?`, time.Now().UTC().Format(time.RFC3339Nano), id)
	return e
}
