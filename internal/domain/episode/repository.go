package episode

type EpisodeRepository interface {
	CreateEpisode(e *Episode) error
	GetEpisodes(limit int) ([]Episode, error)
}
