package db

import (
	"encoding/binary"
	"strings"
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

func subID(s *Subscription) []byte {
	id := strings.Join([]string{s.YouTubeID, s.Filter}, "_")
    return []byte(id)
}

func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}
