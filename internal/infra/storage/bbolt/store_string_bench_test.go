package bbolt

import (
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	bbolt "go.etcd.io/bbolt"
)

func BenchmarkSearchWithString(b *testing.B) {
	db, err := bbolt.Open("test.db", 0666, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		db.Close()
		_ = os.Remove("test.db")
	}()

	testEps := make([]string, 0)
	for i := 0; i < 10; i++ {
		testEps = append(testEps, uuid.NewString())
	}

	testSubs := make([]string, 0)
	for i := 0; i < 10; i++ {
		testSubs = append(testSubs, uuid.NewString())
	}

	var targetUUID string

	err = db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("myBucket"))
		if err != nil {
			return err
		}

		for i := 0; i < b.N; i++ {
			// Serialize testPlaylist struct to binary
			pl := testPlaylist{
				ID:            uuid.NewString(),
				Name:          "Sample Name",
				Episodes:      testEps,
				Subscriptions: testSubs,
			}

			if i == b.N-2 {
				targetUUID = pl.ID
			}

			plData, err := testPlaylistMapToBinary(&pl)
			if err != nil {
				b.Fatal(err)
			}

			// Store binary data in the bucket
			err = bucket.Put([]byte(pl.ID), plData)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		b.Fatal(err)
	}

	// Run the benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := db.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte("myBucket"))
			if bucket == nil {
				b.Fatal("No bucket")
			}

			c := bucket.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				p := string(v[:])
				if strings.Contains(p, "\"ID\":\""+targetUUID+"\"") {
					break
				}
			}
			return nil
		})

		if err != nil {
			b.Fatal(err)
		}
	}
}
