package consistency

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Checker struct{ DB *sql.DB }

func (c Checker) Project(ctx context.Context, id int64) error {
	var budget, spent int64
	var state string
	if e := c.DB.QueryRowContext(ctx, `SELECT budget_cents,spent_cents,state FROM projects WHERE id=?`, id).Scan(&budget, &spent, &state); e != nil {
		return e
	}
	if budget < 0 || spent < 0 || spent > budget {
		return fmt.Errorf("budget invariant violated")
	}
	if state == "completed" && spent == 0 {
		return errors.New("completed project has no spending")
	}
	return nil
}
func (c Checker) Reviews(ctx context.Context, id int64) error {
	var total, done int
	row := c.DB.QueryRowContext(ctx, `SELECT COUNT(*),COALESCE(SUM(CASE WHEN state='submitted' THEN 1 ELSE 0 END),0) FROM reviews WHERE project_id=?`, id)
	if e := row.Scan(&total, &done); e != nil {
		return e
	}
	if done > total {
		return errors.New("review count invariant")
	}
	return nil
}
func (c Checker) ForeignKeys(ctx context.Context) error {
	var v int
	if e := c.DB.QueryRowContext(ctx, `PRAGMA foreign_keys`).Scan(&v); e != nil {
		return e
	}
	if v != 1 {
		return errors.New("foreign keys disabled")
	}
	return nil
}
