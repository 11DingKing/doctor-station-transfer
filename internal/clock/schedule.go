package clock

import "time"

func UTCDate(t time.Time) time.Time     { return t.UTC().Truncate(24 * time.Hour) }
func IsExpired(now, due time.Time) bool { return !due.After(now) }
func NextBackoff(attempt int, base time.Duration) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 8 {
		attempt = 8
	}
	return base * time.Duration(1<<attempt)
}
func InWindow(now, start, end time.Time) bool { return !now.Before(start) && now.Before(end) }
