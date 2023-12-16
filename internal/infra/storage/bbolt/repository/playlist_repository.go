package repository

import (
	"errors"
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/converter"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/playlist"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"
	"wckd1/tg-youtube-podcasts-bot/utils"

	bolt "go.etcd.io/bbolt"
)

const playlistsBucketName = "playlists"

var _ playlist.PlaylistRepository = (*PlaylistRepository)(nil)

type PlaylistRepository struct {
	store *bbolt.BBoltStore
}

func NewPlaylistRepository(store *bbolt.BBoltStore) playlist.PlaylistRepository {
	return &PlaylistRepository{store}
}

func (r PlaylistRepository) SavePlaylist(pl *playlist.Playlist) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(playlistsBucketName))
		if err != nil {
			return err
		}

		subData, err := converter.PlaylistToBinary(pl)
		if err != nil {
			return errors.Join(playlist.ErrPlaylistEncoding, err)
		}
		return b.Put([]byte(pl.ID()), subData)
	})
}

func (r PlaylistRepository) GetPlaylist(id string) (playlist.Playlist, error) {
	var pl playlist.Playlist

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(playlistsBucketName))
		if b == nil {
			return playlist.ErrNoPlaylistsStorage
		}

		plData := b.Get([]byte(id))

		if plData == nil {
			return playlist.ErrPlaylistNotFound
		}

		decodedPl, err := converter.BinaryToPlaylist(plData)
		if err != nil {
			return errors.Join(playlist.ErrPlaylistDecoding, err)
		}
		pl = decodedPl
		return nil
	})

	return pl, err
}

func (r PlaylistRepository) GetPlaylistsWithSubscription(subID string) ([]playlist.Playlist, error) {
	var result []playlist.Playlist

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(playlistsBucketName))
		if b == nil {
			return playlist.ErrNoPlaylistsStorage
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			pl, err := converter.BinaryToPlaylist(v)
			if err != nil {
				log.Printf("[WARN] failed to unmarshal, %+v", errors.Join(playlist.ErrPlaylistDecoding, err))
				continue
			}

			if utils.Contains(pl.Subscriptions(), subID) {
				result = append(result, pl)
			}
		}
		return nil
	})

	return result, err
}

func (r PlaylistRepository) DeletePlaylist(id string) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(playlistsBucketName))
		if b == nil {
			return playlist.ErrNoPlaylistsStorage
		}

		return b.Delete([]byte(id))
	})
}
