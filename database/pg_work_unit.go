package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PgWorkUnit implements WorkUnit using the primary pgx pool.
type PgWorkUnit struct {
	pool *pgxpool.Pool
}

// NewPgWorkUnit creates a new PgWorkUnit backed by the given connection pool.
func NewPgWorkUnit(pool *pgxpool.Pool) *PgWorkUnit {
	return &PgWorkUnit{pool: pool}
}

// Run begins a transaction, calls fn with a ctx that carries the tx (for ContextRouter), commits
// if fn returns nil, and rolls back if fn returns a non-nil error or if fn panics (panic is re-raised).
func (w *PgWorkUnit) Run(ctx context.Context, fn func(context.Context) error) (err error) {
	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	err = fn(contextWithTx(ctx, tx))
	return
}

var _ WorkUnit = (*PgWorkUnit)(nil)
