package worker

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"log/slog"
	"time"
)

type Worker struct {
	Jobs       repository.Jobs
	Milestones repository.Milestones
	Log        *slog.Logger
	Interval   time.Duration
}

func (w Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			w.tick(ctx, now)
		}
	}
}
func (w Worker) tick(ctx context.Context, now time.Time) {
	if _, e := w.Milestones.Overdue(ctx, now); e != nil {
		w.Log.Error("milestone sweep", "error", e)
	}
	j, e := w.Jobs.Claim(ctx, now)
	if e != nil {
		if !errors.Is(e, sql.ErrNoRows) {
			w.Log.Error("job claim", "error", e)
		}
		return
	}
	if e = w.Jobs.Finish(ctx, j.ID, nil); e != nil {
		w.Log.Error("job finish", "error", e)
	}
}
