package db

import "context"

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL);`,
	`CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT NOT NULL,email TEXT NOT NULL UNIQUE,role TEXT NOT NULL,password_hash TEXT NOT NULL,active INTEGER NOT NULL DEFAULT 1,created_at TEXT NOT NULL); CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);`,
	`CREATE TABLE IF NOT EXISTS projects(id INTEGER PRIMARY KEY AUTOINCREMENT,title TEXT NOT NULL,summary TEXT NOT NULL,owner_id INTEGER NOT NULL REFERENCES users(id),budget_cents INTEGER NOT NULL,spent_cents INTEGER NOT NULL DEFAULT 0,state TEXT NOT NULL,version INTEGER NOT NULL DEFAULT 1,due_at TEXT NOT NULL,created_at TEXT NOT NULL,updated_at TEXT NOT NULL); CREATE INDEX IF NOT EXISTS idx_projects_state_due ON projects(state,due_at);`,
	`CREATE TABLE IF NOT EXISTS reviews(id INTEGER PRIMARY KEY AUTOINCREMENT,project_id INTEGER NOT NULL REFERENCES projects(id),reviewer_id INTEGER NOT NULL REFERENCES users(id),state TEXT NOT NULL,score INTEGER NOT NULL DEFAULT 0,comment TEXT NOT NULL,due_at TEXT NOT NULL,version INTEGER NOT NULL DEFAULT 1,created_at TEXT NOT NULL,UNIQUE(project_id,reviewer_id)); CREATE INDEX IF NOT EXISTS idx_reviews_due ON reviews(state,due_at);`,
	`CREATE TABLE IF NOT EXISTS contracts(id INTEGER PRIMARY KEY AUTOINCREMENT,project_id INTEGER NOT NULL UNIQUE REFERENCES projects(id),number TEXT NOT NULL UNIQUE,amount_cents INTEGER NOT NULL,state TEXT NOT NULL,signed_at TEXT,created_at TEXT NOT NULL);`,
	`CREATE TABLE IF NOT EXISTS milestones(id INTEGER PRIMARY KEY AUTOINCREMENT,project_id INTEGER NOT NULL REFERENCES projects(id),name TEXT NOT NULL,due_at TEXT NOT NULL,state TEXT NOT NULL,sequence INTEGER NOT NULL,completed_at TEXT,UNIQUE(project_id,sequence)); CREATE INDEX IF NOT EXISTS idx_milestones_due ON milestones(state,due_at);`,
	`CREATE TABLE IF NOT EXISTS transfers(id INTEGER PRIMARY KEY AUTOINCREMENT,project_id INTEGER NOT NULL REFERENCES projects(id),actor_id INTEGER NOT NULL REFERENCES users(id),kind TEXT NOT NULL,artifact_ref TEXT NOT NULL,checksum TEXT NOT NULL,version INTEGER NOT NULL DEFAULT 1,created_at TEXT NOT NULL);`,
	`CREATE TABLE IF NOT EXISTS audits(id INTEGER PRIMARY KEY AUTOINCREMENT,actor_id INTEGER REFERENCES users(id),entity_type TEXT NOT NULL,entity_id INTEGER NOT NULL,action TEXT NOT NULL,result TEXT NOT NULL,request_id TEXT NOT NULL,created_at TEXT NOT NULL); CREATE INDEX IF NOT EXISTS idx_audits_entity ON audits(entity_type,entity_id);`,
	`CREATE TABLE IF NOT EXISTS sessions(id INTEGER PRIMARY KEY AUTOINCREMENT,user_id INTEGER NOT NULL REFERENCES users(id),token_hash TEXT NOT NULL UNIQUE,expires_at TEXT NOT NULL,revoked_at TEXT,created_at TEXT NOT NULL); CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token_hash);`,
	`CREATE TABLE IF NOT EXISTS idempotency(key TEXT NOT NULL,scope TEXT NOT NULL,response TEXT NOT NULL,created_at TEXT NOT NULL,PRIMARY KEY(key,scope));`,
	`CREATE TABLE IF NOT EXISTS jobs(id INTEGER PRIMARY KEY AUTOINCREMENT,kind TEXT NOT NULL,payload TEXT NOT NULL,state TEXT NOT NULL,attempts INTEGER NOT NULL DEFAULT 0,run_at TEXT NOT NULL,last_error TEXT NOT NULL DEFAULT '',created_at TEXT NOT NULL,updated_at TEXT NOT NULL); CREATE INDEX IF NOT EXISTS idx_jobs_ready ON jobs(state,run_at);`,
}

func Migrate(ctx context.Context, d *DB) error {
	for i, m := range migrations {
		var exists int
		e := d.QueryRowContext(ctx, "SELECT COUNT(1) FROM schema_migrations WHERE version=?", i+1).Scan(&exists)
		if e != nil {
			if i != 0 {
				return e
			}
		}
		if exists > 0 {
			continue
		}
		tx, e := d.BeginTx(ctx, nil)
		if e != nil {
			return e
		}
		if _, e = tx.ExecContext(ctx, m); e != nil {
			tx.Rollback()
			return e
		}
		if _, e = tx.ExecContext(ctx, "INSERT INTO schema_migrations(version,applied_at) VALUES(?,datetime('now'))", i+1); e != nil {
			tx.Rollback()
			return e
		}
		if e = tx.Commit(); e != nil {
			return e
		}
	}
	return nil
}
