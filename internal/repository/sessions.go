package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"time"
)

type Sessions struct{ DB *sql.DB }

func (r Sessions) Create(ctx context.Context, s domain.Session) error {
	_, e := r.DB.ExecContext(ctx, `INSERT INTO sessions(user_id,token_hash,expires_at,created_at) VALUES(?,?,?,?)`, s.UserID, s.TokenHash, s.ExpiresAt.Format(time.RFC3339Nano), s.CreatedAt.Format(time.RFC3339Nano))
	return e
}
func (r Sessions) Revoke(ctx context.Context, hash string) error {
	_, e := r.DB.ExecContext(ctx, `UPDATE sessions SET revoked_at=? WHERE token_hash=? AND revoked_at IS NULL`, time.Now().UTC().Format(time.RFC3339Nano), hash)
	return e
}
func (r Sessions) Find(ctx context.Context, hash string) (domain.Session, error) {
	var s domain.Session
	var exp, created string
	var revoked sql.NullString
	e := r.DB.QueryRowContext(ctx, `SELECT id,user_id,token_hash,expires_at,revoked_at,created_at FROM sessions WHERE token_hash=?`, hash).Scan(&s.ID, &s.UserID, &s.TokenHash, &exp, &revoked, &created)
	s.ExpiresAt, _ = time.Parse(time.RFC3339Nano, exp)
	s.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	if revoked.Valid {
		t, _ := time.Parse(time.RFC3339Nano, revoked.String)
		s.RevokedAt = &t
	}
	return s, e
}
