package repository

import (
	"encoding/json"
	"errors"
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/converter"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"

	bolt "go.etcd.io/bbolt"
)

const episodesBucketName = "episodes"

var _ episode.EpisodeRepository = (*EpisodeRepository)(nil)

type EpisodeRepository struct {
	store *bbolt.BBoltStore
}

func NewEpisodeRepository(store *bbolt.BBoltStore) episode.EpisodeRepository {
	return &EpisodeRepository{store}
}

func (r *EpisodeRepository) SaveEpisode(ep *episode.Episode) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(episodesBucketName))
		if err != nil {
			return err
		}

		epData, err := converter.EpisodeToBinary(ep)
		if err != nil {
			return errors.Join(episode.ErrEpisodeEncoding, err)
		}

		return b.Put([]byte(ep.ID()), epData)
	})
}

func (r *EpisodeRepository) GetEpisode(id string) (episode.Episode, error) {
	var ep episode.Episode

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(episodesBucketName))
		if b == nil {
			return episode.ErrNoEpisodesStorage
		}

		epData := b.Get([]byte(id))
		if epData == nil {
			return episode.ErrEpisodeNotFound
		}

		decodedEpisode, err := converter.BinaryToEpisode(epData)
		if err != nil {
			return errors.Join(episode.ErrEpisodeDecoding, err)
		}
		ep = decodedEpisode
		return nil
	})

	return ep, err
}

func (r *EpisodeRepository) GetEpisodes(limit int) ([]episode.Episode, error) {
	var result []episode.Episode

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(episodesBucketName))
		if b == nil {
			return episode.ErrNoEpisodesStorage
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
