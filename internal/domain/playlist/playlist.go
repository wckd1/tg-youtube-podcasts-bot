package playlist

type Playlist struct {
	id       string
	name     string
	episodes []string
}

func NewPlaylist(id, name string, episodes []string) Playlist {
	return Playlist{id, name, episodes}
}

func (p Playlist) ID() string         { return p.id }
func (p Playlist) Name() string       { return p.name }
func (p Playlist) Episodes() []string { return p.episodes }
