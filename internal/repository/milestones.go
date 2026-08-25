package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"time"
)

type Milestones struct{ DB *sql.DB }

func (r Milestones) Create(ctx context.Context, m domain.Milestone) error {
	_, e := r.DB.ExecContext(ctx, `INSERT INTO milestones(project_id,name,due_at,state,sequence) VALUES(?,?,?,?,?)`, m.ProjectID, m.Name, m.DueAt.Format(time.RFC3339Nano), m.State, m.Sequence)
	return e
}
func (r Milestones) Overdue(ctx context.Context, now time.Time) (int, error) {
	x, e := r.DB.ExecContext(ctx, `UPDATE milestones SET state='overdue' WHERE state='pending' AND due_at<?`, now.Format(time.RFC3339Nano))
	if e != nil {
		return 0, e
	}
	n, _ := x.RowsAffected()
	return int(n), nil
}
func (r Milestones) Complete(ctx context.Context, id int64) error {
	x, e := r.DB.ExecContext(ctx, `UPDATE milestones SET state='accepted',completed_at=? WHERE id=? AND state='pending'`, time.Now().UTC().Format(time.RFC3339Nano), id)
	if e != nil {
		return e
	}
	n, _ := x.RowsAffected()
	if n == 0 {
		return domain.ErrConflict
	}
	return nil
}
