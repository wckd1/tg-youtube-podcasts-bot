package db

import (
	bolt "go.etcd.io/bbolt"
)

type Store interface {
	Querier
}

type BoltStore struct {
	db *bolt.DB
	*Queries
}

func NewStore(db *bolt.DB) Store {
	return &BoltStore{
		db:      db,
		Queries: New(db),
	}
}
