package metrics

import (
	"sync"
	"testing"
	"time"
)

func TestCounters(t *testing.T) {
	var c Counters
	c.Observe(time.Millisecond, false)
	c.Observe(2*time.Millisecond, true)
	r, f, avg := c.Snapshot()
	if r != 2 || f != 1 || avg == 0 {
		t.Fatalf("%d %d %v", r, f, avg)
	}
}
func TestCountersConcurrent(t *testing.T) {
	var c Counters
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				c.Observe(time.Microsecond, false)
			}
		}()
	}
	wg.Wait()
	r, _, _ := c.Snapshot()
	if r != 1000 {
		t.Fatal(r)
	}
}
