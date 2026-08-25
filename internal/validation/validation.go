package validation

import (
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"net/mail"
	"strings"
	"time"
)

func Email(v string) bool { _, e := mail.ParseAddress(v); return e == nil }
func Title(v string) bool {
	v = strings.TrimSpace(v)
	return len([]rune(v)) >= 2 && len([]rune(v)) <= 120
}
func Budget(v int64) bool                     { return v >= 100 && v <= 1000000000 }
func DueDate(v time.Time, now time.Time) bool { return v.After(now) && v.Before(now.AddDate(5, 0, 0)) }
func Project(p domain.Project, now time.Time) error {
	if !Title(p.Title) || !Budget(p.BudgetCents) || !DueDate(p.DueAt, now) {
		return domain.ErrInvalid
	}
	return nil
}
