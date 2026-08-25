package worker

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/clock"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"log/slog"
	"time"
)

type RetryWorker struct {
	Jobs        repository.Jobs
	Log         *slog.Logger
	MaxAttempts int
	Base        time.Duration
}

func (w RetryWorker) RunOnce(ctx context.Context) error {
	j, e := w.Jobs.Claim(ctx, time.Now().UTC())
	if e != nil {
		return e
	}
	if j.Attempts >= w.MaxAttempts {
		w.Jobs.Finish(ctx, j.ID, nil)
		return nil
	}
	e = w.Jobs.Finish(ctx, j.ID, nil)
	if e != nil {
		w.Log.Error("finish", "error", e)
	}
	return e
}
func RetryAt(now time.Time, attempt int, base time.Duration) time.Time {
	return now.Add(clock.NextBackoff(attempt, base))
}
