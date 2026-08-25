package validation

import (
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"testing"
	"time"
)

func TestEmailValidation(t *testing.T) {
	for _, v := range []string{"a@example.com", "team+1@example.cn"} {
		if !Email(v) {
			t.Errorf("valid %s", v)
		}
	}
	for _, v := range []string{"bad", "@x", "a@"} {
		if Email(v) {
			t.Errorf("invalid %s", v)
		}
	}
}
func TestTitleValidation(t *testing.T) {
	if !Title("药材质量研究") {
		t.Fatal("title")
	}
	if Title("x") {
		t.Fatal("short")
	}
	if Title(" ") {
		t.Fatal("blank")
	}
}
func TestBudgetValidation(t *testing.T) {
	if !Budget(100) || !Budget(100000) {
		t.Fatal("budget")
	}
	if Budget(1) || Budget(2000000000) {
		t.Fatal("range")
	}
}
func TestDueDateValidation(t *testing.T) {
	n := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if !DueDate(n.Add(time.Hour), n) {
		t.Fatal("future")
	}
	if DueDate(n.Add(-time.Hour), n) {
		t.Fatal("past")
	}
}
func TestProjectValidation(t *testing.T) {
	n := time.Now()
	p := domain.Project{Title: "有效课题", BudgetCents: 1000, DueAt: n.Add(time.Hour)}
	if e := Project(p, n); e != nil {
		t.Fatal(e)
	}
	p.BudgetCents = 1
	if Project(p, n) == nil {
		t.Fatal("bad budget")
	}
}
