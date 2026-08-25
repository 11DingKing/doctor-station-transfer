package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"time"
)

type Audits struct{ DB *sql.DB }

func (r Audits) Add(ctx context.Context, a domain.AuditEvent) error {
	_, e := r.DB.ExecContext(ctx, `INSERT INTO audits(actor_id,entity_type,entity_id,action,result,request_id,created_at) VALUES(?,?,?,?,?,?,?)`, a.ActorID, a.EntityType, a.EntityID, a.Action, a.Result, a.RequestID, a.CreatedAt.Format(time.RFC3339Nano))
	return e
}
func (r Audits) ForEntity(ctx context.Context, t string, id int64) ([]domain.AuditEvent, error) {
	rows, e := r.DB.QueryContext(ctx, `SELECT id,actor_id,entity_type,entity_id,action,result,request_id,created_at FROM audits WHERE entity_type=? AND entity_id=? ORDER BY id`, t, id)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.AuditEvent
	for rows.Next() {
		var a domain.AuditEvent
		var ts string
		if e = rows.Scan(&a.ID, &a.ActorID, &a.EntityType, &a.EntityID, &a.Action, &a.Result, &a.RequestID, &ts); e != nil {
			return nil, e
		}
		a.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, a)
	}
	return out, rows.Err()
}
