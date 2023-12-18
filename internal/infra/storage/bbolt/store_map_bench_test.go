package bbolt

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	bbolt "go.etcd.io/bbolt"
)

func BenchmarkSearchWithMap(b *testing.B) {
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
				p := binaryMapToTestPlaylist(v, b)

				if p.ID == targetUUID {
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

func testPlaylistMapToBinary(p *testPlaylist) ([]byte, error) {
	sData := map[string]interface{}{
		"id":            p.ID,
		"name":          p.Name,
		"episodes":      p.Episodes,
		"subscriptions": p.Subscriptions,
	}

	return json.Marshal(sData)
}

func binaryMapToTestPlaylist(d []byte, b *testing.B) testPlaylist {
	var plData map[string]interface{}
	if err := json.Unmarshal(d, &plData); err != nil {
		b.Fatal(err)
	}

	id, ok := plData["ID"].(string)
	if !ok {
		b.Fatal(fmt.Errorf("missing or invalid ID field"))
	}

	name, ok := plData["Name"].(string)
	if !ok {
		b.Fatal(fmt.Errorf("missing or invalid Name field"))
	}

	epInterfase, ok := plData["Episodes"]
	if !ok {
		b.Fatal(fmt.Errorf("missing Episodes field"))
	}
	var episodes []string
	if epSlice, ok := epInterfase.([]interface{}); ok {
		for _, item := range epSlice {
			if str, isString := item.(string); isString {
				episodes = append(episodes, str)
			} else {
				b.Fatal(fmt.Errorf("invalid Episodes field"))
			}
		}
	}

	subInterface, ok := plData["Subscriptions"]
	if !ok {
		b.Fatal(fmt.Errorf("missing Subscriptions field"))
	}
	var subscriptions []string
	if subSlice, ok := subInterface.([]interface{}); ok {
		for _, item := range subSlice {
			if str, isString := item.(string); isString {
				subscriptions = append(subscriptions, str)
			} else {
				b.Fatal(fmt.Errorf("invalid Subscriptions field"))
			}
		}
	}

	return testPlaylist{ID: id, Name: name, Episodes: episodes, Subscriptions: subscriptions}
}
