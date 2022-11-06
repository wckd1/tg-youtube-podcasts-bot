package db

import (
	"encoding/binary"
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

func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}
