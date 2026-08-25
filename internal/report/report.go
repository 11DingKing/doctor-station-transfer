package report

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Summary struct {
	ProjectID                  int64
	Reviews, Accepted, Pending int
	Spent, Budget              int64
	GeneratedAt                time.Time
}
type Service struct{ DB *sql.DB }

func (s Service) Build(ctx context.Context, id int64) (Summary, error) {
	var x Summary
	x.ProjectID = id
	x.GeneratedAt = time.Now().UTC()
	e := s.DB.QueryRowContext(ctx, `SELECT budget_cents,spent_cents FROM projects WHERE id=?`, id).Scan(&x.Budget, &x.Spent)
	if e != nil {
		return x, e
	}
	if e = s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM reviews WHERE project_id=?`, id).Scan(&x.Reviews); e != nil {
		return x, e
	}
	if e = s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM reviews WHERE project_id=? AND state='submitted'`, id).Scan(&x.Accepted); e != nil {
		return x, e
	}
	x.Pending = x.Reviews - x.Accepted
	return x, nil
}
func (s Service) Export(ctx context.Context, id int64) (string, error) {
	x, e := s.Build(ctx, id)
	if e != nil {
		return "", e
	}
	return fmt.Sprintf("project=%d reviews=%d accepted=%d pending=%d spent=%d budget=%d", x.ProjectID, x.Reviews, x.Accepted, x.Pending, x.Spent, x.Budget), nil
}
func (s Service) WithinBudget(x Summary) bool { return x.Budget >= 0 && x.Spent <= x.Budget }
func (s Service) CompletionRatio(x Summary) float64 {
	if x.Reviews == 0 {
		return 0
	}
	return float64(x.Accepted) / float64(x.Reviews)
}
