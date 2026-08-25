package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"time"
)

type Users struct{ DB *sql.DB }

func (r Users) Create(ctx context.Context, u domain.User) (int64, error) {
	q := `INSERT INTO users(name,email,role,password_hash,active,created_at) VALUES(?,?,?,?,1,?)`
	x, e := r.DB.ExecContext(ctx, q, u.Name, u.Email, u.Role, u.PasswordHash, u.CreatedAt.Format(time.RFC3339Nano))
	if e != nil {
		return 0, e
	}
	return x.LastInsertId()
}
func (r Users) ByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	var role string
	var active int
	var ts string
	e := r.DB.QueryRowContext(ctx, `SELECT id,name,email,role,password_hash,active,created_at FROM users WHERE email=?`, email).Scan(&u.ID, &u.Name, &u.Email, &role, &u.PasswordHash, &active, &ts)
	u.Role = domain.Role(role)
	u.Active = active == 1
	u.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
	return u, e
}
func (r Users) ByID(ctx context.Context, id int64) (domain.User, error) {
	var u domain.User
	var role string
	var active int
	var ts string
	e := r.DB.QueryRowContext(ctx, `SELECT id,name,email,role,password_hash,active,created_at FROM users WHERE id=?`, id).Scan(&u.ID, &u.Name, &u.Email, &role, &u.PasswordHash, &active, &ts)
	u.Role = domain.Role(role)
	u.Active = active == 1
	u.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
	return u, e
}
