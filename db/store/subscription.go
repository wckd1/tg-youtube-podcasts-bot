package db

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

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

func (q *Queries) DeleteSubsctiption(sub *Subscription) error {
    return q.db.Update(func(tx *bolt.Tx) error {
        b, err := tx.CreateBucketIfNotExists([]byte("subscriptions"))
        if err != nil {
            return err
        }

		return b.Delete([]byte(sub.ID))
	})
}
