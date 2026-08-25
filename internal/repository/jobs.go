package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"time"
)

type Jobs struct{ DB *sql.DB }

func (r Jobs) Enqueue(ctx context.Context, j domain.Job) error {
	_, e := r.DB.ExecContext(ctx, `INSERT INTO jobs(kind,payload,state,attempts,run_at,last_error,created_at,updated_at) VALUES(?,?,?,0,?,?,?,?)`, j.Kind, j.Payload, "queued", j.RunAt.Format(time.RFC3339Nano), "", j.CreatedAt.Format(time.RFC3339Nano), j.UpdatedAt.Format(time.RFC3339Nano))
	return e
}
func (r Jobs) Claim(ctx context.Context, now time.Time) (domain.Job, error) {
	tx, e := r.DB.BeginTx(ctx, nil)
	if e != nil {
		return domain.Job{}, e
	}
	var j domain.Job
	var run, cr, up string
	e = tx.QueryRowContext(ctx, `SELECT id,kind,payload,state,attempts,run_at,last_error,created_at,updated_at FROM jobs WHERE state='queued' AND run_at<=? ORDER BY id LIMIT 1`, now.Format(time.RFC3339Nano)).Scan(&j.ID, &j.Kind, &j.Payload, &j.State, &j.Attempts, &run, &j.LastError, &cr, &up)
	if e != nil {
		tx.Rollback()
		return j, e
	}
	if _, e = tx.ExecContext(ctx, `UPDATE jobs SET state='running',attempts=attempts+1,updated_at=? WHERE id=?`, now.Format(time.RFC3339Nano), j.ID); e != nil {
		tx.Rollback()
		return j, e
	}
	if e = tx.Commit(); e != nil {
		return j, e
	}
	j.State = "running"
	j.Attempts++
	j.RunAt, _ = time.Parse(time.RFC3339Nano, run)
	return j, nil
}
func (r Jobs) Finish(ctx context.Context, id int64, err error) error {
	state := "done"
	msg := ""
	if err != nil {
		state = "queued"
		msg = err.Error()
	}
	_, e := r.DB.ExecContext(ctx, `UPDATE jobs SET state=?,last_error=?,updated_at=? WHERE id=?`, state, msg, time.Now().UTC().Format(time.RFC3339Nano), id)
	return e
}
