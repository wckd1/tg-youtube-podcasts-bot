package db

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type Episode struct {
	UUID        string    `xml:"guid"`
	Enclosure   Enclosure `xml:"enclosure"`
	Link        string    `xml:"link"`
	Image       string    `xml:"image"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Author      string    `xml:"author,omitempty"`
	Duration    int       `xml:"duration,omitempty"`
	PubDate     string    `xml:"pubDate,omitempty"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

func (q *Queries) CreateEpisode(e *Episode) error {
	return q.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("episodes"))
		if err != nil {
			return err
		}

		uuid := uuid.New().String()
		e.UUID = uuid

		buf, err := json.Marshal(e)
		if err != nil {
			return err
		}

		return b.Put([]byte(uuid), buf)
	})
}

func (q *Queries) GetEpisodes(limit int) ([]Episode, error) {
	var result []Episode

	err := q.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("episodes"))
		if b == nil {
			return fmt.Errorf("no saved episodes")
		}

		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			e := Episode{}
			if err := json.Unmarshal(v, &e); err != nil {
				log.Printf("[WARN] failed to unmarshal, %v", err)
				continue
			}
			if len(result) >= limit {
				break
			}
			result = append(result, e)
		}
		return nil
	})

	return result, err
}
