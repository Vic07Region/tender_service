package database

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func NewService(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

//func (q *Queries) WithTx(tx *sql.Tx) *Queries {
//	return &Queries{
//		db: tx,
//	}
//}
