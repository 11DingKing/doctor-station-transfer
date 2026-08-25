package filter

import (
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"testing"
	"time"
)

func TestApplyStateOwnerText(t *testing.T) {
	base := time.Now()
	items := []domain.Project{{ID: 1, Title: "黄芪研究", OwnerID: 2, State: domain.ProjectActive, CreatedAt: base}, {ID: 2, Title: "党参", OwnerID: 2, State: domain.ProjectDraft, CreatedAt: base.Add(time.Minute)}, {ID: 3, Title: "黄精", OwnerID: 3, State: domain.ProjectActive, CreatedAt: base.Add(2 * time.Minute)}}
	x := Apply(items, Projects{State: domain.ProjectActive, OwnerID: 2, Text: "黄"})
	if len(x) != 1 || x[0].ID != 1 {
		t.Fatal(x)
	}
}
func TestApplyDoesNotMutate(t *testing.T) {
	items := []domain.Project{{ID: 1, Title: "a", CreatedAt: time.Now()}}
	Apply(items, Projects{})
	if len(items) != 1 {
		t.Fatal("mutated")
	}
}
func TestGroupByState(t *testing.T) {
	m := GroupByState([]domain.Project{{State: domain.ProjectDraft}, {State: domain.ProjectDraft}, {State: domain.ProjectActive}})
	if m[domain.ProjectDraft] != 2 || m[domain.ProjectActive] != 1 {
		t.Fatal(m)
	}
}
