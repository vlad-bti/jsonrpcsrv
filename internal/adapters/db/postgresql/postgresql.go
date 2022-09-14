package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/vlad-bti/jsonrpcsrv/pkg/logger"
	"github.com/vlad-bti/jsonrpcsrv/pkg/postgres"
)

type baseStorage struct {
	db *postgres.Postgres
}

type transactor struct {
	baseStorage
	log *logger.Logger
}

func NewTransactor(log *logger.Logger, pg *postgres.Postgres) *transactor {
	return &transactor{
		baseStorage{pg},
		log,
	}
}

type txKey struct{}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*pgx.Tx); ok {
		return tx
	}
	return nil
}

func (r *baseStorage) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	tx := extractTx(ctx)
	if tx != nil {
		return (*tx).Exec(ctx, sql, args...)
	}
	return r.db.Pool.Exec(ctx, sql, args...)
}

func (r *baseStorage) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	tx := extractTx(ctx)
	if tx != nil {
		return (*tx).Query(ctx, sql, args...)
	}
	return r.db.Pool.Query(ctx, sql, args...)
}

// WithinTransaction runs function within transaction
//
// The transaction commits when function were finished without error
func (r *transactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) (err error) {
	// begin transaction
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		// finalize transaction on panic, etc.
		errTx := tx.Rollback(ctx)
		if errTx != nil && !errors.Is(errTx, pgx.ErrTxClosed) {
			r.log.Info("close transaction: %v", errTx)
			err = errTx
		}
	}()

	// run callback
	err = tFunc(injectTx(ctx, &tx))
	if err != nil {
		// if error, rollback
		_ = tx.Rollback(ctx)
		r.log.Info("rollback transaction: %v", err)
		return err
	}

	// if no error, commit
	return tx.Commit(ctx)
}
