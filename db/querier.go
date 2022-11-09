package db

import (
	bolt "go.etcd.io/bbolt"
)

type Querier interface {
	CreateSubsctiption(sub *Subscription) (string, error)
	DeleteSubsctiption(sub *Subscription) error
	CreateEpisode(e *Episode) error
	GetEpisodes(limit int) ([]Episode, error)
	ChangeUpdate(u *Update) error
}

var _ Querier = (*Queries)(nil)

type Queries struct {
	db *bolt.DB
}

func New(db *bolt.DB) *Queries {
	return &Queries{db: db}
}
