package milestone

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"time"
)

type Service struct{ DB *sql.DB }

func (s Service) Add(ctx context.Context, m domain.Milestone) error {
	if m.Name == "" || m.Sequence < 1 || m.DueAt.Before(time.Now().UTC()) {
		return domain.ErrInvalid
	}
	_, e := s.DB.ExecContext(ctx, `INSERT INTO milestones(project_id,name,due_at,state,sequence) VALUES(?,?,?,?,?)`, m.ProjectID, m.Name, m.DueAt.Format(time.RFC3339Nano), domain.MilestonePending, m.Sequence)
	return e
}
func (s Service) Complete(ctx context.Context, id int64) error {
	x, e := s.DB.ExecContext(ctx, `UPDATE milestones SET state='accepted',completed_at=? WHERE id=? AND state='pending'`, time.Now().UTC().Format(time.RFC3339Nano), id)
	if e != nil {
		return e
	}
	n, _ := x.RowsAffected()
	if n == 0 {
		return errors.Join(domain.ErrConflict, errors.New("milestone not pending"))
	}
	return nil
}
func (s Service) List(ctx context.Context, pid int64) ([]domain.Milestone, error) {
	rows, e := s.DB.QueryContext(ctx, `SELECT id,project_id,name,due_at,state,sequence,completed_at FROM milestones WHERE project_id=? ORDER BY sequence`, pid)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.Milestone
	for rows.Next() {
		var m domain.Milestone
		var due, st string
		var done sql.NullString
		if e = rows.Scan(&m.ID, &m.ProjectID, &m.Name, &due, &st, &m.Sequence, &done); e != nil {
			return nil, e
		}
		m.DueAt, _ = time.Parse(time.RFC3339Nano, due)
		m.State = domain.MilestoneState(st)
		if done.Valid {
			t, _ := time.Parse(time.RFC3339Nano, done.String)
			m.CompletedAt = &t
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
