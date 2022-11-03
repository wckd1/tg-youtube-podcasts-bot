package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateSubsctiption(ctx context.Context, arg CreateSubscriptionParams) error //(Subscription, error)
}

var _ Querier = (*Queries)(nil)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Queries struct {
	db DBTX
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}
