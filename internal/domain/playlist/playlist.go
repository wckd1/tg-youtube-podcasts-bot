package playlist

const DefaultPlaylistName = "default"

type Playlist struct {
	id            string
	name          string
	episodes      []string
	subscriptions []string
}

func NewPlaylist(id, name string, episodes, subscriptions []string) Playlist {
	if subscriptions == nil {
		subscriptions = make([]string, 0)
	}
	return Playlist{id, name, episodes, subscriptions}
}

func (p Playlist) ID() string   { return p.id }
func (p Playlist) Name() string { return p.name }

func (p Playlist) Episodes() []string { return p.episodes }
func (p *Playlist) AddEpisode(id string) {
	p.episodes = append(p.episodes, id)
}

func (p Playlist) Subscriptions() []string { return p.subscriptions }
func (p *Playlist) AddSubscription(id string) {
	p.subscriptions = append(p.subscriptions, id)
}
