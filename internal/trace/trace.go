package trace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type key struct{}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key{}, id)
}
func RequestID(ctx context.Context) string {
	if v, ok := ctx.Value(key{}).(string); ok && v != "" {
		return v
	}
	return NewID()
}
func NewID() string {
	b := make([]byte, 12)
	if _, e := rand.Read(b); e != nil {
		return "fallback"
	}
	return hex.EncodeToString(b)
}
