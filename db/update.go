package db

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Update struct {
	SubscriptionID string
	UpdateInterval time.Duration
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

		return b.Put([]byte(u.SubscriptionID), buf)
	})
}

func (q *Queries) GetUpdates() ([]Update, error) {
	var result []Update

	err := q.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("updates"))
		if b == nil {
			return fmt.Errorf("no saved updates")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			u := Update{}
			if err := json.Unmarshal(v, &u); err != nil {
				log.Printf("[WARN] failed to unmarshal, %v", err)
				continue
			}
			result = append(result, u)
		}
		return nil
	})

	return result, err
}

func (q *Queries) DeleteUpdate(id string) error {
	return q.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("updates"))
		if b == nil {
			return fmt.Errorf("no saved updates")
		}

		return b.Delete([]byte(id))
	})
}
