package repository

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
)

var (
	ErrNoEpisodesStorage = errors.New("no saved episodes")
	ErrEpisodeNotFound   = errors.New("episode not found")
	ErrEpisodeEncoding   = errors.New("can't encode episode")
	ErrEpisodeDecoding   = errors.New("can't decode episode")
)

type EpisodeRepository interface {
	CheckExist(id string) error
	SaveEpisode(e *entity.Episode) error
	GetEpisode(id string) (entity.Episode, error)
}
