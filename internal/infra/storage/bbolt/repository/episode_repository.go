package repository

import (
	"encoding/json"
	"errors"
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

var (
	ErrNoEpisodesBucket = errors.New("no saved episodes")
	ErrEpisodeEncoding  = errors.New("can't encode episode")
	ErrEpisodeDecoding  = errors.New("can't decode episode")
)

const episodesBucketName = "episodes"

var _ episode.EpisodeRepository = (*EpisodeRepository)(nil)

type EpisodeRepository struct {
	store *bbolt.BBoltStore
}

func NewEpisodeRepository(store *bbolt.BBoltStore) episode.EpisodeRepository {
	return &EpisodeRepository{store}
}

func (r *EpisodeRepository) CreateEpisode(e *episode.Episode) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(episodesBucketName))
		if err != nil {
			return err
		}

		uuid := uuid.New().String()
		// e.ID() = uuid

		buf, err := json.Marshal(e)
		if err != nil {
			return err
		}

		return b.Put([]byte(uuid), buf)
	})
}

func (r *EpisodeRepository) GetEpisodes(limit int) ([]episode.Episode, error) {
	var result []episode.Episode

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(episodesBucketName))
		if b == nil {
			return ErrNoEpisodesBucket
		}

		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			e := episode.Episode{}
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
