package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"time"
)

type Auth struct {
	Users    repository.Users
	Sessions repository.Sessions
	Now      func() time.Time
	Secret   string
}

func (a Auth) Login(ctx context.Context, email, password string) (domain.User, string, error) {
	u, e := a.Users.ByEmail(ctx, email)
	if e != nil || !u.Active {
		return u, "", domain.ErrUnauthorized
	}
	if hash(password) != u.PasswordHash {
		return u, "", domain.ErrUnauthorized
	}
	b := make([]byte, 32)
	if _, e = rand.Read(b); e != nil {
		return u, "", e
	}
	raw := hex.EncodeToString(b)
	s := domain.Session{UserID: u.ID, TokenHash: a.hashToken(raw), ExpiresAt: a.Now().Add(8 * time.Hour), CreatedAt: a.Now()}
	if e = a.Sessions.Create(ctx, s); e != nil {
		return u, "", e
	}
	return u, raw, nil
}
func (a Auth) Logout(ctx context.Context, token string) error {
	return a.Sessions.Revoke(ctx, a.hashToken(token))
}
func (a Auth) Current(ctx context.Context, token string) (domain.User, error) {
	s, e := a.Sessions.Find(ctx, a.hashToken(token))
	if e != nil {
		return domain.User{}, domain.ErrUnauthorized
	}
	if s.RevokedAt != nil || !s.ExpiresAt.After(a.Now()) {
		return domain.User{}, domain.ErrUnauthorized
	}
	return a.Users.ByID(ctx, s.UserID)
}
func (a Auth) hashToken(t string) string {
	h := sha256.Sum256([]byte(a.Secret + t))
	return hex.EncodeToString(h[:])
}
func hash(t string) string         { h := sha256.Sum256([]byte(t)); return hex.EncodeToString(h[:]) }
func HashPassword(p string) string { return hash(p) }

var _ = fmt.Sprintf
