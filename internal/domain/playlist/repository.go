package playlist

type PlaylistRepository interface {
	CreatePlaylist(*Playlist) error
}
