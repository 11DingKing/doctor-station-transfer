package main

import (
	"context"
	"database/sql"
	"github.com/11DingKing/doctor-station-transfer/internal/clock"
	"github.com/11DingKing/doctor-station-transfer/internal/config"
	"github.com/11DingKing/doctor-station-transfer/internal/db"
	"github.com/11DingKing/doctor-station-transfer/internal/domain"
	"github.com/11DingKing/doctor-station-transfer/internal/httpapi"
	"github.com/11DingKing/doctor-station-transfer/internal/repository"
	"github.com/11DingKing/doctor-station-transfer/internal/service"
	"github.com/11DingKing/doctor-station-transfer/internal/worker"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	d, e := db.Open(ctx, cfg.DBPath)
	if e != nil {
		panic(e)
	}
	defer d.Close()
	if e = db.Migrate(ctx, d); e != nil {
		panic(e)
	}
	now := clock.Real{}
	users := repository.Users{DB: d.DB}
	seed(ctx, users, now)
	auth := service.Auth{Users: users, Sessions: repository.Sessions{DB: d.DB}, Now: now.Now, Secret: cfg.SessionSecret}
	ps := repository.Projects{DB: d.DB}
	aud := repository.Audits{DB: d.DB}
	projects := service.Projects{DB: d.DB, Repo: ps, Contracts: repository.Contracts{DB: d.DB}, Milestones: repository.Milestones{DB: d.DB}, Audits: aud, Clock: now}
	reviews := service.Reviews{DB: d.DB, Repo: repository.Reviews{DB: d.DB}, Projects: ps, Audits: aud, Clock: now}
	transfers := service.Transfers{DB: d.DB, Repo: repository.Transfers{DB: d.DB}, Projects: ps, Audits: aud, Clock: now}
	w := worker.Worker{Jobs: repository.Jobs{DB: d.DB}, Milestones: repository.Milestones{DB: d.DB}, Log: slog.Default(), Interval: time.Duration(cfg.WorkerIntervalSeconds) * time.Second}
	go w.Run(ctx)
	srv := &http.Server{Addr: cfg.Addr, Handler: httpapi.New(httpapi.RouterDeps{Auth: auth, Projects: projects, Reviews: reviews, Transfers: transfers}), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if e := srv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
			slog.Error("server", "error", e)
		}
	}()
	<-ctx.Done()
	shut, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	srv.Shutdown(shut)
}
func seed(ctx context.Context, u repository.Users, c clock.Clock) {
	if _, e := u.ByEmail(ctx, "admin@doctor.local"); e == sql.ErrNoRows {
		u.Create(ctx, domain.User{Name: "Admin", Email: "admin@doctor.local", Role: domain.RoleAdmin, PasswordHash: service.HashPassword("admin"), CreatedAt: c.Now()})
	}
}
