package health

import (
	"context"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"testing"
)

func TestLiveness(t *testing.T) {
	if Liveness()["status"] != "ok" {
		t.Fatal("status")
	}
}
func TestReady(t *testing.T) {
	d, e := db.Open(context.Background(), "file:health?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	defer d.Close()
	if e = (&Checker{DB: d.DB}).Ready(context.Background()); e != nil {
		t.Fatal(e)
	}
}
