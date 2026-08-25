package worker

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"log/slog"
	"testing"
	"time"
)

func TestRetryAt(t *testing.T) {
	v := RetryAt(time.Unix(0, 0), 2, time.Second)
	if v.Sub(time.Unix(0, 0)) != 4*time.Second {
		t.Fatal(v)
	}
}
func TestWorkerCancel(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	c()
	w := Worker{Jobs: repository.Jobs{}, Milestones: repository.Milestones{}, Log: slog.Default(), Interval: time.Millisecond}
	done := make(chan struct{})
	go func() { w.Run(ctx); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("not stopped")
	}
}
func TestJobStateValues(t *testing.T) {
	for _, s := range []string{"queued", "running", "done"} {
		if s == "" {
			t.Fatal(s)
		}
	}
}
func TestWorkerFixtureMigration(t *testing.T) {
	ctx := context.Background()
	d, e := db.Open(ctx, "file:worker-"+t.Name()+"?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	defer d.Close()
	if e = db.Migrate(ctx, d); e != nil {
		t.Fatal(e)
	}
	j := repository.Jobs{DB: d.DB}
	if e = j.Enqueue(ctx, domain.Job{Kind: "x", Payload: "{}", RunAt: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}); e != nil {
		t.Fatal(e)
	}
}
