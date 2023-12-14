package bbolt

import (
	"context"
	"log"

	bolt "go.etcd.io/bbolt"
)

type BBoltStore struct {
	ctx        context.Context
	connection *bolt.DB
}

func NewStore(ctx context.Context) *BBoltStore {
	dbConn, err := bolt.Open("storage/yt_podcasts.db", 0666, nil)
	if err != nil {
		log.Fatal("[ERROR] cannot connect to database:", err)
	}

	return &BBoltStore{
		ctx:        ctx,
		connection: dbConn,
	}
}

func (s BBoltStore) Close() error {
	return s.connection.Close()
}

func (s BBoltStore) Update(handler func(tx *bolt.Tx) error) error {
	return s.connection.Update(handler)
}

func (s BBoltStore) View(handler func(tx *bolt.Tx) error) error {
	return s.connection.View(handler)
}
