package repository

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/playlist"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"
)

const playlistsBucketName = "playlists"

var (
	ErrNoPlaylistsBucket = errors.New("no saved playlists")
	ErrPlaylistNotFound  = errors.New("playlist not found")
	ErrPlaylistEncoding  = errors.New("can't encode playlist")
	ErrPlaylistDecoding  = errors.New("can't decode playlist")
)

var _ playlist.PlaylistRepository = (*PlaylistRepository)(nil)

type PlaylistRepository struct {
	store *bbolt.BBoltStore
}

func NewPlaylistRepository(store *bbolt.BBoltStore) playlist.PlaylistRepository {
	return &PlaylistRepository{store}
}

func (r PlaylistRepository) CreatePlaylist(*playlist.Playlist) error {
	return nil
}
