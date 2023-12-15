package episode

type EpisodeUsecase struct {
	episodeRepository EpisodeRepository
}

func NewEpisodeUsecase(episodeRepository EpisodeRepository) *EpisodeUsecase {
	return &EpisodeUsecase{episodeRepository}
}
