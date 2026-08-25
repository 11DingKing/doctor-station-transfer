package query

import "testing"

func TestBuildFilters(t *testing.T) {
	q, a := Build(Filter{State: "active", Owner: "7", Limit: 10, Offset: 2})
	if q == "" || len(a) != 4 {
		t.Fatalf("%s %#v", q, a)
	}
}
func TestBuildDefaults(t *testing.T) {
	q, a := Build(Filter{})
	if q == "" || len(a) != 2 || a[0] != 20 {
		t.Fatalf("%s %#v", q, a)
	}
}
func TestBuildOffsetClamp(t *testing.T) {
	_, a := Build(Filter{Offset: -1})
	if a[1] != 0 {
		t.Fatal(a)
	}
}
