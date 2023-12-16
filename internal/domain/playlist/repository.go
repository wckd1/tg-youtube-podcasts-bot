package playlist

import "errors"

var (
	ErrNoPlaylistsStorage = errors.New("no saved playlists")
	ErrPlaylistNotFound   = errors.New("playlist not found")
	ErrPlaylistEncoding   = errors.New("can't encode playlist")
	ErrPlaylistDecoding   = errors.New("can't decode playlist")
)

type PlaylistRepository interface {
	SavePlaylist(playlist *Playlist) error
	GetPlaylist(id string) (Playlist, error)
	GetPlaylistsWithSubscription(subID string) ([]Playlist, error)
	DeletePlaylist(id string) error
}
