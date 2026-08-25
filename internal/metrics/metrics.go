package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

type Counters struct {
	requests  atomic.Uint64
	failures  atomic.Uint64
	mu        sync.Mutex
	latencies []time.Duration
}

func (c *Counters) Observe(d time.Duration, failed bool) {
	c.requests.Add(1)
	if failed {
		c.failures.Add(1)
	}
	c.mu.Lock()
	if len(c.latencies) >= 1000 {
		copy(c.latencies, c.latencies[1:])
		c.latencies = c.latencies[:999]
	}
	c.latencies = append(c.latencies, d)
	c.mu.Unlock()
}
func (c *Counters) Snapshot() (uint64, uint64, time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var sum time.Duration
	for _, d := range c.latencies {
		sum += d
	}
	if len(c.latencies) == 0 {
		return c.requests.Load(), c.failures.Load(), 0
	}
	return c.requests.Load(), c.failures.Load(), sum / time.Duration(len(c.latencies))
}
