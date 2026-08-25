package repository

import (
	"context"
	"database/sql"
	"errors"
)

func ExecTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, e := db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	if e = fn(tx); e != nil {
		if r := tx.Rollback(); r != nil {
			return errors.Join(e, r)
		}
		return e
	}
	return tx.Commit()
}
func RowsAffected(res sql.Result) (int64, error) {
	if res == nil {
		return 0, errors.New("nil result")
	}
	return res.RowsAffected()
}
