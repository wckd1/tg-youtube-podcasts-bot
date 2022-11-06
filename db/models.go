package db

type SourceType string

const (
	Channel  SourceType = "channel"
	Playlist SourceType = "playlist"
	Video    SourceType = "video"
)

type Subscription struct {
	YouTubeID 	string
	Type 		SourceType
	Filter      string
}

type Download struct {
	ID          int
	URL      string
	CoverURL    string
	Title       string
	Description string
}
