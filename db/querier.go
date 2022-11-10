package db

import (
	bolt "go.etcd.io/bbolt"
)

type Querier interface {
	CreateSubsctiption(sub *Subscription) error
	GetSubsctiption(id string) (Subscription, error)
	DeleteSubsctiption(sub *Subscription) error
	CreateEpisode(e *Episode) error
	GetEpisodes(limit int) ([]Episode, error)
	ChangeUpdate(u *Update) error
	GetUpdates() ([]Update, error)
	DeleteUpdate(id string) error
}

var _ Querier = (*Queries)(nil)

type Queries struct {
	db *bolt.DB
}

func New(db *bolt.DB) *Queries {
	return &Queries{db: db}
}
