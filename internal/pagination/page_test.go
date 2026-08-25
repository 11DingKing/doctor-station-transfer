package pagination

import "testing"

func TestParseBounds(t *testing.T) {
	for _, c := range []struct {
		l, o   string
		el, eo int
	}{{"", "", 20, 0}, {"5", "2", 5, 2}, {"1000", "-1", 20, 0}, {"10", "9", 10, 9}} {
		r := Parse(c.l, c.o)
		if r.Limit != c.el || r.Offset != c.eo {
			t.Errorf("%+v", r)
		}
	}
}
func TestResultShape(t *testing.T) {
	r := Result[string]{Items: []string{"a"}, Total: 1, Limit: 10, Offset: 0}
	if len(r.Items) != 1 || r.Total != 1 {
		t.Fatal("shape")
	}
}
