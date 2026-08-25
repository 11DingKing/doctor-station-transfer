package notification

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"testing"
	"time"
)

func TestQueueAndPending(t *testing.T) {
	ctx := context.Background()
	d, e := db.Open(ctx, "file:notif-"+t.Name()+"?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	defer d.Close()
	db.Migrate(ctx, d)
	o := Outbox{DB: d.DB}
	if e = o.Queue(ctx, Message{To: "a", Subject: "s", Body: "b", CreatedAt: time.Now()}); e != nil {
		t.Fatal(e)
	}
	n, e := o.Pending(ctx)
	if e != nil || n != 1 {
		t.Fatal(n, e)
	}
}
func TestQueueMultiple(t *testing.T) {
	ctx := context.Background()
	d, _ := db.Open(ctx, "file:notif2-"+t.Name()+"?mode=memory&cache=shared")
	defer d.Close()
	db.Migrate(ctx, d)
	o := Outbox{DB: d.DB}
	for i := 0; i < 5; i++ {
		o.Queue(ctx, Message{To: "a", Subject: "s", Body: "b", CreatedAt: time.Now()})
	}
	n, _ := o.Pending(ctx)
	if n != 5 {
		t.Fatal(n)
	}
}
