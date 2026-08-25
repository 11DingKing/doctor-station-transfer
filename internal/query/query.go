package query

import (
	"fmt"
	"strings"
)

type Filter struct {
	State, Owner  string
	Limit, Offset int
}

func Build(f Filter) (string, []any) {
	parts := []string{"1=1"}
	args := []any{}
	if f.State != "" {
		parts = append(parts, "state=?")
		args = append(args, f.State)
	}
	if f.Owner != "" {
		parts = append(parts, "owner_id=?")
		args = append(args, f.Owner)
	}
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	return fmt.Sprintf("WHERE %s ORDER BY created_at DESC LIMIT ? OFFSET ?", strings.Join(parts, " AND ")), append(args, f.Limit, f.Offset)
}
