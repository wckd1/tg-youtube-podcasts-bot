package db

import (
	bolt "go.etcd.io/bbolt"
)

type Querier interface {
	SaveSubsctiption(sub *Subscription) error
	GetSubscriptions() ([]Subscription, error)
	DeleteSubsctiption(sub *Subscription) error
	CreateEpisode(e *Episode) error
	GetEpisodes(limit int) ([]Episode, error)
}

var _ Querier = (*Queries)(nil)

type Queries struct {
	db *bolt.DB
}

func New(db *bolt.DB) *Queries {
	return &Queries{db: db}
}
