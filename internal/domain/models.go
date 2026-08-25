package domain

import "time"

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleResearcher Role = "researcher"
	RoleReviewer   Role = "reviewer"
	RoleAuditor    Role = "auditor"
)

type User struct {
	ID           int64
	Name, Email  string
	Role         Role
	PasswordHash string
	Active       bool
	CreatedAt    time.Time
}
type Project struct {
	ID                      int64
	Title, Summary          string
	OwnerID                 int64
	BudgetCents, SpentCents int64
	State                   ProjectState
	Version                 int64
	DueAt                   time.Time
	CreatedAt, UpdatedAt    time.Time
}
type Review struct {
	ID, ProjectID, ReviewerID int64
	State                     ReviewState
	Score                     int
	Comment                   string
	DueAt                     time.Time
	Version                   int64
	CreatedAt                 time.Time
}
type Contract struct {
	ID, ProjectID int64
	Number        string
	AmountCents   int64
	State         string
	SignedAt      *time.Time
	CreatedAt     time.Time
}
type Milestone struct {
	ID, ProjectID int64
	Name          string
	DueAt         time.Time
	State         MilestoneState
	Sequence      int
	CompletedAt   *time.Time
}
type Transfer struct {
	ID, ProjectID, ActorID      int64
	Kind, ArtifactRef, Checksum string
	Version                     int64
	CreatedAt                   time.Time
}
type AuditEvent struct {
	ID, ActorID               int64
	EntityType                string
	EntityID                  int64
	Action, Result, RequestID string
	CreatedAt                 time.Time
}
type Session struct {
	ID, UserID int64
	TokenHash  string
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	CreatedAt  time.Time
}
type Job struct {
	ID                   int64
	Kind, Payload, State string
	Attempts             int
	RunAt                time.Time
	LastError            string
	CreatedAt, UpdatedAt time.Time
}
type IdempotencyRecord struct {
	Key, Scope string
	Response   string
	CreatedAt  time.Time
}
