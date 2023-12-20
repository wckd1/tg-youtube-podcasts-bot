package repository

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
)

var (
	ErrNoPlaylistsStorage = errors.New("no saved playlists")
	ErrPlaylistNotFound   = errors.New("playlist not found")
	ErrPlaylistEncoding   = errors.New("can't encode playlist")
	ErrPlaylistDecoding   = errors.New("can't decode playlist")
)

type PlaylistRepository interface {
	SavePlaylist(playlist *entity.Playlist) error
	GetPlaylist(id string) (entity.Playlist, error)
	GetPlaylistByName(name string) (entity.Playlist, error)
	GetPlaylistsWithSubscription(subID string) ([]entity.Playlist, error)
	DeletePlaylist(id string) error
}
