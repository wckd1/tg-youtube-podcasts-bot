package repository

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/converter"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/repository"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"

	bolt "go.etcd.io/bbolt"
)

const episodesBucketName = "episodes"

var _ repository.EpisodeRepository = (*EpisodeRepository)(nil)

type EpisodeRepository struct {
	store *bbolt.BBoltStore
}

func NewEpisodeRepository(store *bbolt.BBoltStore) repository.EpisodeRepository {
	return &EpisodeRepository{store}
}

func (r *EpisodeRepository) CheckExist(id string) error {
	return r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(episodesBucketName))
		if b == nil {
			return repository.ErrNoEpisodesStorage
		}

		epData := b.Get([]byte(id))
		if epData == nil {
			return repository.ErrEpisodeNotFound
		}

		return nil
	})
}

func (r *EpisodeRepository) SaveEpisode(ep *entity.Episode) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(episodesBucketName))
		if err != nil {
			return err
		}

		epData, err := converter.EpisodeToBinary(ep)
		if err != nil {
			return errors.Join(repository.ErrEpisodeEncoding, err)
		}

		return b.Put([]byte(ep.ID()), epData)
	})
}

func (r *EpisodeRepository) GetEpisode(id string) (entity.Episode, error) {
	var ep entity.Episode

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(episodesBucketName))
		if b == nil {
			return repository.ErrNoEpisodesStorage
		}

		epData := b.Get([]byte(id))
		if epData == nil {
			return repository.ErrEpisodeNotFound
		}

		decodedEpisode, err := converter.BinaryToEpisode(epData)
		if err != nil {
			return errors.Join(repository.ErrEpisodeDecoding, err)
		}
		ep = decodedEpisode
		return nil
	})

	return ep, err
}
