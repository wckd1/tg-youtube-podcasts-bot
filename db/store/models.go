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
