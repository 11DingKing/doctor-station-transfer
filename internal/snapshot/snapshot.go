package snapshot

import (
	"context"
	"database/sql"
	"time"
)

type Record struct {
	ProjectID  int64
	State      string
	Version    int64
	CapturedAt time.Time
}
type Store struct{ DB *sql.DB }

func (s Store) Capture(ctx context.Context, id int64) (Record, error) {
	var r Record
	r.ProjectID = id
	r.CapturedAt = time.Now().UTC()
	e := s.DB.QueryRowContext(ctx, `SELECT state,version FROM projects WHERE id=?`, id).Scan(&r.State, &r.Version)
	return r, e
}
func (s Store) Equal(a, b Record) bool {
	return a.ProjectID == b.ProjectID && a.State == b.State && a.Version == b.Version
}
func (s Store) Changed(a, b Record) bool { return !s.Equal(a, b) }
func (s Store) Version(a Record) int64   { return a.Version }
