package db

import (
	"context"
	"path/filepath"
	"testing"
)

func openTest(t *testing.T) *DB {
	t.Helper()
	d, e := Open(context.Background(), "file::memory:?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	if e = Migrate(context.Background(), d); e != nil {
		t.Fatal(e)
	}
	return d
}
func TestMigrationCreatesRelations(t *testing.T) {
	d := openTest(t)
	defer d.Close()
	var n int
	if e := d.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table'`).Scan(&n); e != nil {
		t.Fatal(e)
	}
	if n < 11 {
		t.Fatalf("tables %d", n)
	}
}
func TestMigrationIdempotent(t *testing.T) {
	d := openTest(t)
	defer d.Close()
	if e := Migrate(context.Background(), d); e != nil {
		t.Fatal(e)
	}
	var n int
	d.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&n)
	if n != 11 {
		t.Fatalf("versions %d", n)
	}
}
func TestForeignKeys(t *testing.T) {
	d := openTest(t)
	defer d.Close()
	_, e := d.Exec(`INSERT INTO projects(title,summary,owner_id,budget_cents,state,version,due_at,created_at,updated_at) VALUES('x','x',99,1,'draft',1,datetime('now'),datetime('now'),datetime('now'))`)
	if e == nil {
		t.Fatal("foreign key accepted")
	}
}
func TestReopenPersists(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "reopen.db")
	d, e := Open(ctx, path)
	if e != nil {
		t.Fatal(e)
	}
	Migrate(ctx, d)
	d.Exec(`CREATE TABLE IF NOT EXISTS marker(v TEXT)`)
	d.Exec(`INSERT INTO marker VALUES('ok')`)
	d.Close()
	d2, e := Open(ctx, path)
	if e != nil {
		t.Fatal(e)
	}
	defer d2.Close()
	var v string
	if e = d2.QueryRow(`SELECT v FROM marker`).Scan(&v); e != nil {
		t.Fatal(e)
	}
	if v != "ok" {
		t.Fatal(v)
	}
}
