package user

type User struct {
	id              string
	defaultPlaylist string
	playlists       []string
	subscriptions   []string
}

func NewUser(id, defaultPlaylist string, playlists, subscriptions []string) User {
	if playlists == nil {
		playlists = make([]string, 0)
	}
	if subscriptions == nil {
		subscriptions = make([]string, 0)
	}

	return User{id, defaultPlaylist, playlists, subscriptions}
}

func (u User) ID() string { return u.id }

func (u User) DefaultPlaylist() string { return u.defaultPlaylist }
func (u *User) SetDefaultPlaylist(id string) {
	u.defaultPlaylist = id
}

func (u User) Playlists() []string { return u.playlists }
func (u *User) AddPlaylist(id string) {
	u.playlists = append(u.playlists, id)
}

func (u User) Subscriptions() []string { return u.subscriptions }
