package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// PoolInterface defines the interface for database pool operations.
// This allows for mocking in tests.
type PoolInterface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Close()
}

// DatabasePort is the interface used by repositories to execute SQL.
type DatabasePort interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// RowScanner is satisfied by pgx.Row and any type that can scan column values.
type RowScanner interface {
	Scan(dest ...any) error
}

// WorkUnit runs a callback inside a single database transaction boundary.
// Implementations attach transactional state to the context passed to fn so
// DatabasePort calls inside fn use the same transaction.
type WorkUnit interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
