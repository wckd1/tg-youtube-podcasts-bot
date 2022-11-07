package db

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type Update struct {
	SubscriptionID string
	UpdateInterval string
	LastUpdated    time.Time
}

func (q *Queries) ChangeUpdate(u *Update) error {
	return q.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("updates"))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		uuid := uuid.New().String()
		return b.Put([]byte(uuid), buf)
	})
}
