package retry

import (
	"context"
	"time"
)

type Policy struct {
	Attempts int
	Base     time.Duration
}

func Do(ctx context.Context, p Policy, fn func(context.Context) error) error {
	if p.Attempts < 1 {
		p.Attempts = 1
	}
	if p.Base <= 0 {
		p.Base = 10 * time.Millisecond
	}
	var e error
	for i := 0; i < p.Attempts; i++ {
		if e = fn(ctx); e == nil {
			return nil
		}
		if i+1 == p.Attempts {
			break
		}
		t := time.NewTimer(p.Base * time.Duration(1<<i))
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
	}
	return e
}
