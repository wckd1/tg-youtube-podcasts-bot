package db

import (
	"encoding/json"
	"strings"

	bolt "go.etcd.io/bbolt"
)

type SourceType string

const (
	Channel  SourceType = "channel"
	Playlist SourceType = "playlist"
	Video    SourceType = "video"
)

type Subscription struct {
	YouTubeID string
	Type      SourceType
	Filter    string
}

func (s *Subscription) makeID() string {
	return strings.Join([]string{s.YouTubeID, s.Filter}, "_")
}

func (q *Queries) CreateSubsctiption(sub *Subscription) (string, error) {
	var id string

	err := q.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("subscriptions"))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(sub)
		if err != nil {
			return err
		}

		id = sub.makeID()
		return b.Put([]byte(id), buf)
	})

	return id, err
}

func (q *Queries) DeleteSubsctiption(sub *Subscription) error {
	return q.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("subscriptions"))
		if err != nil {
			return err
		}

		id := sub.makeID()
		return b.Delete([]byte(id))
	})
}
