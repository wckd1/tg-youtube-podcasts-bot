package db

type SourceType string

const (
	Channel  SourceType = "channel"
	Playlist SourceType = "playlist"
	Video    SourceType = "video"
)

type Subscription struct {
	ID         string
	SourcePath string
	SourceType SourceType
	Title      string
}

type Download struct {
	ID          int
	Path        string
	CoverURL    string
	Title       string
	Description string
}
