package db

import (
	"encoding/json"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

type Subscription struct {
	ID      string
	URL     string
	IsVideo bool
	Filter  string
}

func (q *Queries) CreateSubsctiption(sub *Subscription) error {
	return q.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("subscriptions"))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(sub)
		if err != nil {
			return err
		}
		return b.Put([]byte(sub.ID), buf)
	})
}

func (q *Queries) GetSubsctiption(id string) (Subscription, error) {
	result := Subscription{}

	err := q.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("subscriptions"))
		if b == nil {
			return fmt.Errorf("no saved subscriptions")
		}

		v := b.Get([]byte(id))
		if err := json.Unmarshal(v, &result); err != nil {
			log.Printf("[WARN] failed to unmarshal, %v", err)
			return err
		}

		return nil
	})

	return result, err
}

func (q *Queries) DeleteSubsctiption(sub *Subscription) error {
	return q.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("subscriptions"))
		if b == nil {
			return fmt.Errorf("no saved subscriptions")
		}

		return b.Delete([]byte(sub.ID))
	})
}
