package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// txContextKey is the private context key used to attach a pgx transaction for DatabasePort routing.
type txContextKey struct{}

func contextWithTx(parent context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(parent, txContextKey{}, tx)
}

func txFromContext(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(txContextKey{}).(pgx.Tx)
	return tx
}

// ContextRouter implements DatabasePort. When ctx carries a transaction (see
// PgWorkUnit), Exec/Query/QueryRow go to that tx; otherwise they delegate to the pool.
type ContextRouter struct {
	pool DatabasePort
}

// NewContextRouter creates a ContextRouter that routes SQL through a context-carried
// transaction or falls back to pool.
func NewContextRouter(pool DatabasePort) *ContextRouter {
	return &ContextRouter{pool: pool}
}

func (r *ContextRouter) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	tag, err := r.exec(ctx, sql, args...)
	return tag, classifyPgErr(err)
}

func (r *ContextRouter) exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if tx := txFromContext(ctx); tx != nil {
		return tx.Exec(ctx, sql, args...)
	}
	return r.pool.Exec(ctx, sql, args...)
}

func (r *ContextRouter) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	rows, err := r.query(ctx, sql, args...)
	if err != nil {
		return nil, classifyPgErr(err)
	}
	return rows, nil
}

func (r *ContextRouter) query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if tx := txFromContext(ctx); tx != nil {
		return tx.Query(ctx, sql, args...)
	}
	return r.pool.Query(ctx, sql, args...)
}

func (r *ContextRouter) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return classifyingRow{inner: r.queryRow(ctx, sql, args...)}
}

func (r *ContextRouter) queryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if tx := txFromContext(ctx); tx != nil {
		return tx.QueryRow(ctx, sql, args...)
	}
	return r.pool.QueryRow(ctx, sql, args...)
}

type classifyingRow struct {
	inner pgx.Row
}

func (c classifyingRow) Scan(dest ...any) error {
	return classifyPgErr(c.inner.Scan(dest...))
}

var _ DatabasePort = (*ContextRouter)(nil)
