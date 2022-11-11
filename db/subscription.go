package db

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Subscription struct {
	ID      string
	URL     string
	IsVideo bool
	Filter  string
	UpdateInterval time.Duration
	LastUpdated    time.Time
}

func (q *Queries) SaveSubsctiption(sub *Subscription) error {
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

func (q *Queries) GetSubscriptions() ([]Subscription, error) {
	var result []Subscription

	err := q.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("subscriptions"))
		if b == nil {
			return fmt.Errorf("no saved subscriptions")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			s:= Subscription{}
			if err := json.Unmarshal(v, &s); err != nil {
				log.Printf("[WARN] failed to unmarshal, %v", err)
				continue
			}
			result = append(result, s)
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
