package config

import "os"

type Config struct {
	Addr, DBPath, SessionSecret string
	WorkerIntervalSeconds       int
}

func Load() Config {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	db := os.Getenv("DB_PATH")
	if db == "" {
		db = "doctor.db"
	}
	sec := os.Getenv("SESSION_SECRET")
	if sec == "" {
		sec = "development-secret-change-me"
	}
	return Config{Addr: addr, DBPath: db, SessionSecret: sec, WorkerIntervalSeconds: 5}
}
