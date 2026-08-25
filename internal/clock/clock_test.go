package clock

import (
	"testing"
	"time"
)

func TestFixedClock(t *testing.T) {
	v := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	if (Fixed{T: v}).Now() != v {
		t.Fatal("fixed")
	}
}
func TestUTCDate(t *testing.T) {
	v := time.Date(2026, 2, 3, 4, 5, 6, 0, time.FixedZone("x", 8*3600))
	if UTCDate(v).Hour() != 0 {
		t.Fatal("date")
	}
}
func TestExpired(t *testing.T) {
	n := time.Now()
	if !IsExpired(n, n.Add(-time.Second)) {
		t.Fatal("expired")
	}
	if IsExpired(n, n.Add(time.Second)) {
		t.Fatal("future")
	}
}
func TestWindow(t *testing.T) {
	n := time.Now()
	if !InWindow(n, n.Add(-time.Second), n.Add(time.Second)) {
		t.Fatal("window")
	}
	if InWindow(n.Add(time.Second), n, n.Add(time.Second)) {
		t.Fatal("end")
	}
}
func TestBackoff(t *testing.T) {
	if NextBackoff(0, time.Second) != time.Second {
		t.Fatal("base")
	}
	if NextBackoff(3, time.Second) != 8*time.Second {
		t.Fatal("shift")
	}
	if NextBackoff(99, time.Second) != 256*time.Second {
		t.Fatal("cap")
	}
}
func TestBackoffNegative(t *testing.T) {
	if NextBackoff(-1, time.Second) != time.Second {
		t.Fatal("negative")
	}
}
