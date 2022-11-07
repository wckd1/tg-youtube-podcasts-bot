package db

import (
	"encoding/json"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type Episode struct {
	URL         string
	CoverURL    string
	Title       string
	Description string
}

func (q *Queries) CreateEpisode(e *Episode) error {
	return q.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("episodes"))
		if err != nil {
			return err
		}

		uuid := uuid.New().String()

		buf, err := json.Marshal(e)
		if err != nil {
			return err
		}

		return b.Put([]byte(uuid), buf)
	})
}
