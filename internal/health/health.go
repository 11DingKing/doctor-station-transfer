package health

import (
	"context"
	"database/sql"
	"time"
)

type Checker struct {
	DB      *sql.DB
	Timeout time.Duration
}

func (c Checker) Ready(ctx context.Context) error {
	if c.Timeout <= 0 {
		c.Timeout = time.Second
	}
	x, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()
	return c.DB.PingContext(x)
}
func Liveness() map[string]string {
	return map[string]string{"status": "ok", "service": "doctor-station-transfer"}
}
