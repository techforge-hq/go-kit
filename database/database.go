// Package database provides PostgreSQL database connection and query functionality.
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const sqlPreviewMaxLen = 100

// Database wraps a PostgreSQL connection pool with logging functionality.
type Database[L Logger[L]] struct {
	Pool   PoolInterface
	logger L
}

// NewConnection creates a new database connection pool with proper error handling and logging.
func NewConnection[L Logger[L]](ctx context.Context, connString string, log L) (*Database[L], error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create pool: %w", ErrConnectionFailed, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%w: ping after pool creation: %w", ErrConnectionFailed, err)
	}

	db := &Database[L]{
		Pool:   pool,
		logger: log.With("component", "database"),
	}

	poolConfig := pool.Config()
	db.logger.WithContext(ctx).Info("database connection established",
		"max_connections", poolConfig.MaxConns,
		"min_connections", poolConfig.MinConns,
	)

	return db, nil
}

// Close closes the database connection pool.
func (db *Database[L]) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		db.logger.Info("database connection pool closed")
	}
}

// GetPool returns the underlying PostgreSQL connection pool.
func (db *Database[L]) GetPool() PoolInterface {
	return db.Pool
}

// PgxPool returns the concrete *pgxpool.Pool when the database was built with one (see [NewConnection]).
func (db *Database[L]) PgxPool() (*pgxpool.Pool, bool) {
	p, ok := db.Pool.(*pgxpool.Pool)
	return p, ok
}

// Query executes a query that returns multiple rows.
func (db *Database[L]) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	rows, err := db.Pool.Query(ctx, sql, args...)
	if err != nil {
		db.logger.WithContext(ctx).Error("query execution failed",
			"error", err,
			"sql_preview", truncateSQL(sql),
		)
		return nil, fmt.Errorf("database: query failed [%s]: %w", truncateSQL(sql), err)
	}
	return rows, nil
}

// QueryRow executes a query that returns a single row.
func (db *Database[L]) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}

// QueryRowScan executes a query that returns a single row and scans it using the provided function.
func (db *Database[L]) QueryRowScan(ctx context.Context, scanFunc func(row pgx.Row) error, sql string, args ...any) error {
	row := db.Pool.QueryRow(ctx, sql, args...)
	if err := scanFunc(row); err != nil {
		db.logger.WithContext(ctx).Error("query row scan failed",
			"error", err,
			"sql_preview", truncateSQL(sql),
		)
		return fmt.Errorf("database: query row scan failed [%s]: %w", truncateSQL(sql), err)
	}

	return nil
}

// Exec executes a query that doesn't return rows (INSERT, UPDATE, DELETE).
func (db *Database[L]) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	tag, err := db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		db.logger.WithContext(ctx).Error("exec operation failed",
			"error", err,
			"sql_preview", truncateSQL(sql),
		)
		return tag, fmt.Errorf("database: exec failed [%s]: %w", truncateSQL(sql), err)
	}

	db.logger.WithContext(ctx).Debug("exec operation completed",
		"rows_affected", tag.RowsAffected(),
	)

	return tag, nil
}

// HealthCheck performs a health check on the database connection.
func (db *Database[L]) HealthCheck(ctx context.Context) error {
	if db.Pool == nil {
		return fmt.Errorf("%w: %w", ErrHealthCheck, ErrNilPool)
	}
	if err := db.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("%w: %w", ErrHealthCheck, err)
	}
	return nil
}

// Shutdown gracefully closes the database connection pool.
func (db *Database[L]) Shutdown(ctx context.Context) error {
	db.logger.WithContext(ctx).Info("shutting down database connection")
	if db.Pool != nil {
		db.Pool.Close()
	}
	return nil
}

// truncateSQL truncates SQL string for safe logging.
func truncateSQL(sql string) string {
	if len(sql) <= sqlPreviewMaxLen {
		return sql
	}
	return sql[:sqlPreviewMaxLen] + "..."
}
