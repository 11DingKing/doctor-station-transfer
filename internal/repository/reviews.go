package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"time"
)

type Reviews struct{ DB *sql.DB }

func (r Reviews) Assign(ctx context.Context, x domain.Review) (int64, error) {
	z, e := r.DB.ExecContext(ctx, `INSERT INTO reviews(project_id,reviewer_id,state,score,comment,due_at,version,created_at) VALUES(?,?,?,?,?,?,1,?)`, x.ProjectID, x.ReviewerID, x.State, x.Score, x.Comment, x.DueAt.Format(time.RFC3339Nano), x.CreatedAt.Format(time.RFC3339Nano))
	if e != nil {
		return 0, e
	}
	return z.LastInsertId()
}
func (r Reviews) Transition(ctx context.Context, id int64, from, to domain.ReviewState, version int64, score int, comment string) error {
	x, e := r.DB.ExecContext(ctx, `UPDATE reviews SET state=?,score=?,comment=?,version=version+1 WHERE id=? AND state=? AND version=?`, to, score, comment, id, from, version)
	if e != nil {
		return e
	}
	n, _ := x.RowsAffected()
	if n == 0 {
		return domain.ErrConflict
	}
	return nil
}
func (r Reviews) ForProject(ctx context.Context, pid int64) ([]domain.Review, error) {
	rows, e := r.DB.QueryContext(ctx, `SELECT id,project_id,reviewer_id,state,score,comment,due_at,version,created_at FROM reviews WHERE project_id=? ORDER BY id`, pid)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.Review
	for rows.Next() {
		var x domain.Review
		var st, due, cr string
		if e = rows.Scan(&x.ID, &x.ProjectID, &x.ReviewerID, &st, &x.Score, &x.Comment, &due, &x.Version, &cr); e != nil {
			return nil, e
		}
		x.State = domain.ReviewState(st)
		x.DueAt, _ = time.Parse(time.RFC3339Nano, due)
		x.CreatedAt, _ = time.Parse(time.RFC3339Nano, cr)
		out = append(out, x)
	}
	return out, rows.Err()
}
