package bbolt

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/google/uuid"
	bbolt "go.etcd.io/bbolt"
)

func BenchmarkDirectUnmarshal(b *testing.B) {
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

			id := uuid.NewString()
			bData := []byte(pl.ID)

			nameBytes := []byte(pl.Name)
			nameLen := make([]byte, 4)
			binary.BigEndian.PutUint32(nameLen, uint32(len(nameBytes)))
			bData = append(bData, nameLen...)
			bData = append(bData, nameBytes...)

			epBytes, err := marshalUUIDSliceToBinary(pl.Episodes)
			if err != nil {
				return err
			}
			epBytesLen := make([]byte, 4)
			binary.BigEndian.PutUint32(epBytesLen, uint32(len(pl.Episodes)))
			bData = append(bData, epBytesLen...)
			bData = append(bData, epBytes...)

			subSlice, err := marshalUUIDSliceToBinary(pl.Subscriptions)
			if err != nil {
				return err
			}
			subSliceLen := make([]byte, 4)
			binary.BigEndian.PutUint32(subSliceLen, uint32(len(pl.Subscriptions)))
			bData = append(bData, subSliceLen...)
			bData = append(bData, subSlice...)

			if i == b.N-2 {
				targetUUID = id
			}

			// Store binary data in the bucket
			err = bucket.Put([]byte(pl.ID), bData)
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
				data := v
				p := testPlaylist{}

				p.ID = string(data[:36])
				data = data[36:]

				nameLen := binary.BigEndian.Uint32(data[:4])
				data = data[4:]
				p.Name = string(data[:nameLen])
				data = data[nameLen:]

				// Decode slices of UUIDs
				epLen := binary.BigEndian.Uint32(data[:4])
				data = data[4:]
				p.Episodes = unmarshalBinaryToUUIDSlice(data[:epLen*36])
				data = data[epLen*36:]

				subLen := binary.BigEndian.Uint32(data[:4])
				data = data[4:]
				p.Subscriptions = unmarshalBinaryToUUIDSlice(data[:subLen*36])

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

func marshalUUIDSliceToBinary(slice []string) ([]byte, error) {
	var result []byte
	for _, s := range slice {
		result = append(result, []byte(s)...)
	}
	return result, nil
}

func unmarshalBinaryToUUIDSlice(data []byte) []string {
	var result []string
	for i := 0; i < len(data); i += 36 {
		result = append(result, string(data[i:i+36]))
	}
	return result
}
