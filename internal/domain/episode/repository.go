package episode

import "errors"

var (
	ErrNoEpisodesStorage = errors.New("no saved episodes")
	ErrEpisodeNotFound   = errors.New("episode not found")
	ErrEpisodeEncoding   = errors.New("can't encode episode")
	ErrEpisodeDecoding   = errors.New("can't decode episode")
)

type EpisodeRepository interface {
	CheckExist(id string) error
	SaveEpisode(e *Episode) error
	GetEpisode(id string) (Episode, error)
}
