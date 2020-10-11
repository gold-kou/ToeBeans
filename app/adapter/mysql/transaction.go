package mysql

/*
manage database transaction by context.Context of http request
*/

import (
	"context"
	"database/sql"
)

type DBTransaction struct {
	DB *sql.DB
}

type contextKey string

const transactionContextKey contextKey = "transaction"

func SetTransaction(parent context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(parent, transactionContextKey, tx)
}

func GetTransaction(ctx context.Context) *sql.Tx {
	v := ctx.Value(transactionContextKey)
	tx, ok := v.(*sql.Tx)
	if !ok {
		return nil
	}
	return tx
}

func (t DBTransaction) Do(ctx context.Context, f func(ctx context.Context) error) error {
	tx, e := t.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if e != nil {
		return e
	}
	ctx = SetTransaction(ctx, tx)
	e = f(ctx)
	if e != nil {
		_ = tx.Rollback()
		return e
	}
	return tx.Commit()
}

func NewDBTransaction(db *sql.DB) DBTransaction {
	return DBTransaction{DB: db}
}
