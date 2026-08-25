package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryEventuallySucceeds(t *testing.T) {
	n := 0
	e := Do(context.Background(), Policy{Attempts: 3, Base: time.Millisecond}, func(context.Context) error {
		n++
		if n < 3 {
			return errors.New("wait")
		}
		return nil
	})
	if e != nil || n != 3 {
		t.Fatalf("%v %d", e, n)
	}
}
func TestRetryStopsAfterLimit(t *testing.T) {
	n := 0
	e := Do(context.Background(), Policy{Attempts: 2, Base: time.Millisecond}, func(context.Context) error { n++; return errors.New("bad") })
	if e == nil || n != 2 {
		t.Fatalf("%v %d", e, n)
	}
}
func TestRetryContextCancel(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	c()
	e := Do(ctx, Policy{Attempts: 4, Base: time.Second}, func(context.Context) error { return errors.New("bad") })
	if !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
