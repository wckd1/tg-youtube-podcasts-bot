package user

type User struct {
	id              string
	defaultPlaylist string
	playlists       []string
}

func NewUser(id, defaultPlaylist string, playlists []string) User {
	if playlists == nil {
		playlists = make([]string, 0)
	}

	return User{id, defaultPlaylist, playlists}
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
