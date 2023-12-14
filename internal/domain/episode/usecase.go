package episode

type EpisodeUseCase struct {
	episodeRepository EpisodeRepository
}

func NewEpisodeUseCase(episodeRepository EpisodeRepository) *EpisodeUseCase {
	return &EpisodeUseCase{episodeRepository}
}
