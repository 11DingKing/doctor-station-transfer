package contract

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"time"
)

type Service struct{ DB *sql.DB }

func (s Service) Draft(ctx context.Context, pid int64, number string, amount int64) error {
	if number == "" || amount <= 0 {
		return domain.ErrInvalid
	}
	_, e := s.DB.ExecContext(ctx, `INSERT INTO contracts(project_id,number,amount_cents,state,created_at) VALUES(?,?,?,?,?)`, pid, number, amount, "draft", time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
func (s Service) Sign(ctx context.Context, id int64) error {
	x, e := s.DB.ExecContext(ctx, `UPDATE contracts SET state='signed',signed_at=? WHERE id=? AND state='draft'`, time.Now().UTC().Format(time.RFC3339Nano), id)
	if e != nil {
		return e
	}
	n, _ := x.RowsAffected()
	if n == 0 {
		return errors.Join(domain.ErrConflict, errors.New("contract not draft"))
	}
	return nil
}
func (s Service) ByProject(ctx context.Context, pid int64) (domain.Contract, error) {
	var c domain.Contract
	var ts sql.NullString
	var created string
	e := s.DB.QueryRowContext(ctx, `SELECT id,project_id,number,amount_cents,state,signed_at,created_at FROM contracts WHERE project_id=?`, pid).Scan(&c.ID, &c.ProjectID, &c.Number, &c.AmountCents, &c.State, &ts, &created)
	if ts.Valid {
		t, _ := time.Parse(time.RFC3339Nano, ts.String)
		c.SignedAt = &t
	}
	c.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return c, e
}
