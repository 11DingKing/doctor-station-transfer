package notification

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Message struct {
	To, Subject, Body string
	CreatedAt         time.Time
}
type Sender interface {
	Send(context.Context, Message) error
}
type Outbox struct{ DB *sql.DB }

func (o Outbox) Queue(ctx context.Context, m Message) error {
	_, e := o.DB.ExecContext(ctx, `INSERT INTO jobs(kind,payload,state,run_at,last_error,created_at,updated_at) VALUES('notification',?,?,?, '',?,?)`, fmt.Sprintf("%s|%s|%s", m.To, m.Subject, m.Body), "queued", m.CreatedAt.Format(time.RFC3339Nano), m.CreatedAt.Format(time.RFC3339Nano), m.CreatedAt.Format(time.RFC3339Nano))
	return e
}
func (o Outbox) Pending(ctx context.Context) (int, error) {
	var n int
	e := o.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE kind='notification' AND state='queued'`).Scan(&n)
	return n, e
}
