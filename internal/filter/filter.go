package filter

import (
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"sort"
	"strings"
)

type Projects struct {
	State   domain.ProjectState
	OwnerID int64
	Text    string
}

func Apply(items []domain.Project, f Projects) []domain.Project {
	out := make([]domain.Project, 0, len(items))
	for _, p := range items {
		if f.State != "" && p.State != f.State {
			continue
		}
		if f.OwnerID != 0 && p.OwnerID != f.OwnerID {
			continue
		}
		if f.Text != "" && !strings.Contains(strings.ToLower(p.Title), strings.ToLower(f.Text)) {
			continue
		}
		out = append(out, p)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out
}
func GroupByState(items []domain.Project) map[domain.ProjectState]int {
	m := make(map[domain.ProjectState]int)
	for _, p := range items {
		m[p.State]++
	}
	return m
}
func Clone(items []domain.Project) []domain.Project {
	out := make([]domain.Project, len(items))
	copy(out, items)
	return out
}
