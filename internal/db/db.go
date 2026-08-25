package db

import (
	"context"
	"database/sql"
	_ "modernc.org/sqlite"
)

type DB struct{ *sql.DB }

func Open(ctx context.Context, path string) (*DB, error) {
	d, e := sql.Open("sqlite", path)
	if e != nil {
		return nil, e
	}
	d.SetMaxOpenConns(8)
	d.SetMaxIdleConns(8)
	if _, e = d.ExecContext(ctx, "PRAGMA foreign_keys = ON; PRAGMA busy_timeout = 5000"); e != nil {
		d.Close()
		return nil, e
	}
	return &DB{d}, nil
}
func (d *DB) Close() error { return d.DB.Close() }
