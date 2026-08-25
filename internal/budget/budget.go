package budget

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
)

type Service struct{ DB *sql.DB }

func (s Service) Reserve(ctx context.Context, id, amount, version int64) error {
	if amount <= 0 {
		return domain.ErrInvalid
	}
	x, e := s.DB.ExecContext(ctx, `UPDATE projects SET spent_cents=spent_cents+?,version=version+1 WHERE id=? AND version=? AND spent_cents+?<=budget_cents`, amount, id, version, amount)
	if e != nil {
		return e
	}
	n, _ := x.RowsAffected()
	if n == 0 {
		return errors.Join(domain.ErrConflict, errors.New("budget exhausted or version changed"))
	}
	return nil
}
func (s Service) Release(ctx context.Context, id, amount int64) error {
	x, e := s.DB.ExecContext(ctx, `UPDATE projects SET spent_cents=CASE WHEN spent_cents>=? THEN spent_cents-? ELSE 0 END,version=version+1 WHERE id=?`, amount, amount, id)
	if e != nil {
		return e
	}
	n, _ := x.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
func (s Service) Remaining(ctx context.Context, id int64) (int64, error) {
	var n int64
	e := s.DB.QueryRowContext(ctx, `SELECT budget_cents-spent_cents FROM projects WHERE id=?`, id).Scan(&n)
	return n, e
}
