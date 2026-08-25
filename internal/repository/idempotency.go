package repository

import (
	"context"
	"database/sql"
	"time"
)

type Idempotency struct{ DB *sql.DB }

func (r Idempotency) Put(ctx context.Context, key, scope, response string) error {
	_, e := r.DB.ExecContext(ctx, `INSERT INTO idempotency(key,scope,response,created_at) VALUES(?,?,?,?)`, key, scope, response, time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
func (r Idempotency) Get(ctx context.Context, key, scope string) (string, error) {
	var v string
	e := r.DB.QueryRowContext(ctx, `SELECT response FROM idempotency WHERE key=? AND scope=?`, key, scope).Scan(&v)
	return v, e
}
