package database

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const workUnitTestTable = `public.work_unit_tx_test`

func setupWorkUnitIntegration(t *testing.T) (ctx context.Context, db *Database[noopLogger], router *ContextRouter, wu *PgWorkUnit) {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set: skipping Postgres-backed work unit tests")
	}
	ctx = context.Background()
	log := NewNoopLogger()
	var err error
	db, err = NewConnection(ctx, dsn, log)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	pgPool, ok := db.PgxPool()
	require.True(t, ok, "pool must be *pgxpool.Pool")

	_, err = db.Exec(ctx, `CREATE TABLE IF NOT EXISTS `+workUnitTestTable+` (k INT PRIMARY KEY)`)
	require.NoError(t, err)

	_, err = db.Exec(ctx, `TRUNCATE `+workUnitTestTable)
	require.NoError(t, err)

	router = NewContextRouter(db)
	wu = NewPgWorkUnit(pgPool)
	return ctx, db, router, wu
}

func countWorkUnitTestRows(t *testing.T, ctx context.Context, db *Database[noopLogger]) int {
	t.Helper()
	rows, err := db.Query(ctx, `SELECT COUNT(*) FROM `+workUnitTestTable)
	require.NoError(t, err)
	defer rows.Close()
	require.True(t, rows.Next())
	var n int
	require.NoError(t, rows.Scan(&n))
	require.NoError(t, rows.Err())
	return n
}

func TestPgWorkUnit_Run_CommitsOnSuccess(t *testing.T) {
	ctx, db, router, wu := setupWorkUnitIntegration(t)

	err := wu.Run(ctx, func(txCtx context.Context) error {
		_, err := router.Exec(txCtx, `INSERT INTO `+workUnitTestTable+` (k) VALUES ($1)`, 4242)
		return err
	})
	require.NoError(t, err)

	assert.Equal(t, 1, countWorkUnitTestRows(t, ctx, db))
}

func TestPgWorkUnit_Run_RollsBackOnError(t *testing.T) {
	ctx, db, router, wu := setupWorkUnitIntegration(t)

	err := wu.Run(ctx, func(txCtx context.Context) error {
		if _, err := router.Exec(txCtx, `INSERT INTO `+workUnitTestTable+` (k) VALUES ($1)`, 1); err != nil {
			return err
		}
		return errors.New("force rollback")
	})
	require.Error(t, err)

	assert.Equal(t, 0, countWorkUnitTestRows(t, ctx, db))
}

func TestPgWorkUnit_Run_RollsBackOnPanic(t *testing.T) {
	ctx, db, router, wu := setupWorkUnitIntegration(t)

	assert.Panics(t, func() {
		_ = wu.Run(ctx, func(txCtx context.Context) error {
			if _, err := router.Exec(txCtx, `INSERT INTO `+workUnitTestTable+` (k) VALUES ($1)`, 99); err != nil {
				return err
			}
			panic("boom")
		})
	})

	assert.Equal(t, 0, countWorkUnitTestRows(t, ctx, db))
}

func TestPgWorkUnit_Run_MultiStepRollback(t *testing.T) {
	ctx, db, router, wu := setupWorkUnitIntegration(t)

	err := wu.Run(ctx, func(txCtx context.Context) error {
		if _, err := router.Exec(txCtx, `INSERT INTO `+workUnitTestTable+` (k) VALUES ($1)`, 10); err != nil {
			return err
		}
		if _, err := router.Exec(txCtx, `INSERT INTO `+workUnitTestTable+` (k) VALUES ($1)`, 20); err != nil {
			return err
		}
		return errors.New("abort after two writes")
	})
	require.Error(t, err)

	assert.Equal(t, 0, countWorkUnitTestRows(t, ctx, db))
}

func TestPgWorkUnit_Run_SecondExecFailureRollsBackFirst(t *testing.T) {
	ctx, db, router, wu := setupWorkUnitIntegration(t)

	err := wu.Run(ctx, func(txCtx context.Context) error {
		if _, err := router.Exec(txCtx, `INSERT INTO `+workUnitTestTable+` (k) VALUES ($1)`, 7); err != nil {
			return err
		}
		_, err := router.Exec(txCtx, `INSERT INTO `+workUnitTestTable+` (k) VALUES ($1)`, 7)
		return err
	})
	require.Error(t, err)

	assert.Equal(t, 0, countWorkUnitTestRows(t, ctx, db))
}
