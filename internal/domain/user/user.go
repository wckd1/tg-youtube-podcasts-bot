package user

type User struct {
	id            string
	playlists     []string
	subscriptions []string
}

func NewUser(id string, playlists, subscriptions []string) User {
	if playlists == nil {
		playlists = make([]string, 0)
	}
	if subscriptions == nil {
		subscriptions = make([]string, 0)
	}

	return User{id, playlists, subscriptions}
}

func (u User) ID() string { return u.id }

func (u User) Playlists() []string { return u.playlists }
func (u *User) AddPlaylist(id string) {
	u.playlists = append(u.playlists, id)
}

func (u User) Subscriptions() []string { return u.subscriptions }
