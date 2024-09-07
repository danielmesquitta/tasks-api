package transactioner

import (
	"context"
	"database/sql"
)

type SQLTransactioner struct {
	db *sql.DB
}

func NewSQLTransactioner(db *sql.DB) *SQLTransactioner {
	return &SQLTransactioner{
		db: db,
	}
}

func (tm *SQLTransactioner) Do(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	ctx = context.WithValue(ctx, CtxTxKey, tx)

	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
