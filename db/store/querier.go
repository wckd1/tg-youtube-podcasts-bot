package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateSubsctiption(ctx context.Context, arg CreateSubscriptionParams) error
	DeleteSubsctiption(ctx context.Context, arg DeleteSubscriptionParams) error
	DeleteTitledSubsctiption(ctx context.Context, arg DeleteSubscriptionParams) error
}

var _ Querier = (*Queries)(nil)

type Queries struct {
	db *sql.DB
}

func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}
