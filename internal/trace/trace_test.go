package trace

import (
	"context"
	"testing"
)

func TestRequestIDRoundTrip(t *testing.T) {
	ctx := WithRequestID(context.Background(), "abc")
	if RequestID(ctx) != "abc" {
		t.Fatal("id")
	}
}
func TestRequestIDGenerated(t *testing.T) {
	id := RequestID(context.Background())
	if len(id) < 8 {
		t.Fatal(id)
	}
}
func TestIDsUnique(t *testing.T) {
	a, b := NewID(), NewID()
	if a == b {
		t.Fatal("collision")
	}
}
