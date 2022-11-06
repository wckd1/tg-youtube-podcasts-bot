package db

import (
	bolt "go.etcd.io/bbolt"
)

type Querier interface {
	CreateSubsctiption(sub *Subscription) error
	DeleteSubsctiption(sub *Subscription) error
}

var _ Querier = (*Queries)(nil)

type Queries struct {
	db *bolt.DB
}

func New(db *bolt.DB) *Queries {
	return &Queries{db: db}
}
