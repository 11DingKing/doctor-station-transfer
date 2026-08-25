package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/pagination"
	"time"
)

type Projects struct{ DB *sql.DB }

func (r Projects) Create(ctx context.Context, p domain.Project) (int64, error) {
	x, e := r.DB.ExecContext(ctx, `INSERT INTO projects(title,summary,owner_id,budget_cents,spent_cents,state,version,due_at,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, p.Title, p.Summary, p.OwnerID, p.BudgetCents, 0, p.State, 1, p.DueAt.Format(time.RFC3339Nano), p.CreatedAt.Format(time.RFC3339Nano), p.UpdatedAt.Format(time.RFC3339Nano))
	if e != nil {
		return 0, e
	}
	return x.LastInsertId()
}
func (r Projects) Get(ctx context.Context, id int64) (domain.Project, error) {
	var p domain.Project
	var st, due, cr, up string
	e := r.DB.QueryRowContext(ctx, `SELECT id,title,summary,owner_id,budget_cents,spent_cents,state,version,due_at,created_at,updated_at FROM projects WHERE id=?`, id).Scan(&p.ID, &p.Title, &p.Summary, &p.OwnerID, &p.BudgetCents, &p.SpentCents, &st, &p.Version, &due, &cr, &up)
	p.State = domain.ProjectState(st)
	p.DueAt, _ = time.Parse(time.RFC3339Nano, due)
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, cr)
	p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, up)
	return p, e
}
func (r Projects) Transition(ctx context.Context, id int64, from, to domain.ProjectState, version int64) error {
	x, e := r.DB.ExecContext(ctx, `UPDATE projects SET state=?,version=version+1,updated_at=? WHERE id=? AND state=? AND version=?`, to, time.Now().UTC().Format(time.RFC3339Nano), id, from, version)
	if e != nil {
		return e
	}
	n, _ := x.RowsAffected()
	if n == 0 {
		return domain.ErrConflict
	}
	return nil
}
func (r Projects) List(ctx context.Context, q pagination.Request) (pagination.Result[domain.Project], error) {
	like := fmt.Sprintf("%%%s%%", q.Query)
	rows, e := r.DB.QueryContext(ctx, `SELECT id,title,summary,owner_id,budget_cents,spent_cents,state,version,due_at,created_at,updated_at FROM projects WHERE title LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, like, q.Limit, q.Offset)
	if e != nil {
		return pagination.Result[domain.Project]{}, e
	}
	defer rows.Close()
	out := pagination.Result[domain.Project]{Limit: q.Limit, Offset: q.Offset}
	for rows.Next() {
		var p domain.Project
		var st, due, cr, up string
		if e = rows.Scan(&p.ID, &p.Title, &p.Summary, &p.OwnerID, &p.BudgetCents, &p.SpentCents, &st, &p.Version, &due, &cr, &up); e != nil {
			return out, e
		}
		p.State = domain.ProjectState(st)
		p.DueAt, _ = time.Parse(time.RFC3339Nano, due)
		p.CreatedAt, _ = time.Parse(time.RFC3339Nano, cr)
		p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, up)
		out.Items = append(out.Items, p)
	}
	out.Total = len(out.Items)
	return out, rows.Err()
}
