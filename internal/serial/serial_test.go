package serial

import (
	"testing"
	"time"
)

func TestEnvelopeRoundTrip(t *testing.T) {
	in := map[string]any{"title": "课题", "n": 3}
	b, e := Encode("project", in)
	if e != nil {
		t.Fatal(e)
	}
	var out map[string]any
	if e = Decode(b, &out); e != nil {
		t.Fatal(e)
	}
	if out["title"] != "课题" {
		t.Fatal(out)
	}
}
func TestEnvelopeVersion(t *testing.T) {
	if e := Decode([]byte(`{"version":2,"type":"x","at":"2026-01-01T00:00:00Z","data":{}}`), &map[string]any{}); e == nil {
		t.Fatal("version accepted")
	}
}
func TestEnvelopeTime(t *testing.T) {
	b, e := Encode("x", time.Now())
	if e != nil || len(b) == 0 {
		t.Fatal(e)
	}
}
