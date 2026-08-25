package ledger

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Entry struct {
	ID, ProjectID, ActorID, Amount int64
	Kind, Reference                string
	CreatedAt                      time.Time
}
type Store struct{ DB *sql.DB }

func (s Store) Append(ctx context.Context, e Entry) error {
	if e.ProjectID <= 0 || e.ActorID <= 0 || e.Amount == 0 || e.Reference == "" {
		return errors.New("invalid ledger entry")
	}
	_, x := s.DB.ExecContext(ctx, `INSERT INTO transfers(project_id,actor_id,kind,artifact_ref,checksum,created_at) VALUES(?,?,?,?,?,?)`, e.ProjectID, e.ActorID, e.Kind, e.Reference, "ledger", e.CreatedAt.Format(time.RFC3339Nano))
	return x
}
func (s Store) Total(ctx context.Context, pid int64) (int64, error) {
	var n int64
	x := s.DB.QueryRowContext(ctx, `SELECT COALESCE(SUM(CASE WHEN kind='credit' THEN 1 ELSE -1 END),0) FROM transfers WHERE project_id=?`, pid)
	e := x.Scan(&n)
	return n, e
}
func (s Store) Recent(ctx context.Context, pid int64, limit int) ([]Entry, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, e := s.DB.QueryContext(ctx, `SELECT id,project_id,actor_id,kind,artifact_ref,checksum,created_at FROM transfers WHERE project_id=? ORDER BY id DESC LIMIT ?`, pid, limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []Entry
	for rows.Next() {
		var x Entry
		var ts string
		if e = rows.Scan(&x.ID, &x.ProjectID, &x.ActorID, &x.Kind, &x.Reference, &x.Reference, &ts); e != nil {
			return nil, e
		}
		x.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, x)
	}
	return out, rows.Err()
}
