package db

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

func (q *Queries) CreateDownload(d *Download) error {
    return q.db.Update(func(tx *bolt.Tx) error {
        b, err := tx.CreateBucketIfNotExists([]byte("downloads"))
        if err != nil {
            return err
        }

        id, _ := b.NextSequence()
        d.ID = int(id) 

        buf, err := json.Marshal(d)
        if err != nil {
            return err
        }

        return b.Put(itob(d.ID), buf)
	})
}
