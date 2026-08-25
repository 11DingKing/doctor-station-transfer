# Versioned Migrations

The migration runner in `internal/db/migrate.go` applies versions 1 through 11 transactionally. Each version is an immutable SQL batch and is recorded in `schema_migrations`; this directory documents the production migration contract without duplicating SQL sources.
