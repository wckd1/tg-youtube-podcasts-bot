package app

import (
	"context"
	"log"
	"wckd1/tg-youtube-podcasts-bot/configs"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/playlist"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/service"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/content/youtube"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt/repository"
)

type serviceProvider struct {
	ctx    context.Context
	config configs.Config
	store  *bbolt.BBoltStore

	// Repositories
	userRepository         user.UserRepository
	playlistRepository     playlist.PlaylistRepository
	subscriptionRepository subscription.SubscriptionRepository
	episodeRepository      episode.EpisodeRepository

	// Services
	contentManager service.ContentManager

	// Usecases
	registerUsecase *usecase.RegisterUsecase
	updateUsecase   *usecase.UpdateUsecase
	addUsecase      *usecase.AddUsecase
	rssUseCase      *usecase.RSSUseCase
}

func newServiceProvider(ctx context.Context, config configs.Config) *serviceProvider {
	return &serviceProvider{
		ctx:    ctx,
		config: config,
	}
}

func (s *serviceProvider) Store() *bbolt.BBoltStore {
	if s.store == nil {
		s.store = bbolt.NewStore(s.ctx)
	}
	return s.store
}

// Repositories
func (s *serviceProvider) UserRepository() user.UserRepository {
	if s.userRepository == nil {
		s.userRepository = repository.NewUserRepository(s.Store())
	}

	return s.userRepository
}

func (s *serviceProvider) PlaylistRepository() playlist.PlaylistRepository {
	if s.playlistRepository == nil {
		s.playlistRepository = repository.NewPlaylistRepository(s.Store())
	}

	return s.playlistRepository
}

func (s *serviceProvider) SubscriptionRepository() subscription.SubscriptionRepository {
	if s.subscriptionRepository == nil {
		s.subscriptionRepository = repository.NewSubscriptionRepository(s.Store())
	}

	return s.subscriptionRepository
}

func (s *serviceProvider) EpisodeRepository() episode.EpisodeRepository {
	if s.episodeRepository == nil {
		s.episodeRepository = repository.NewEpisodeRepository(s.Store())
	}

	return s.episodeRepository
}

// Services
func (s *serviceProvider) ContentManager() service.ContentManager {
	if s.contentManager == nil {
		cm, err := youtube.NewYouTubeContentManager()
		if err != nil {
			log.Fatalf("[ERROR] can't init YouTubeContentManager, %+v", err)
		}
		s.contentManager = cm
	}

	return s.contentManager
}

// Usecases
func (s *serviceProvider) RegisterUsecase() *usecase.RegisterUsecase {
	if s.registerUsecase == nil {
		s.registerUsecase = usecase.NewRegisterUsecase(s.UserRepository(), s.PlaylistRepository())
	}

	return s.registerUsecase
}

func (s *serviceProvider) UpdateUsecase() *usecase.UpdateUsecase {
	if s.updateUsecase == nil {
		s.updateUsecase = usecase.NewUpdateUsecase(
			s.UserRepository(),
			s.PlaylistRepository(),
			s.EpisodeRepository(),
			s.SubscriptionRepository(),
		)
	}

	return s.updateUsecase
}

func (s *serviceProvider) AddUsecase() *usecase.AddUsecase {
	if s.addUsecase == nil {
		s.addUsecase = usecase.NewAddUsecase(
			s.UserRepository(),
			s.PlaylistRepository(),
			s.EpisodeRepository(),
			s.SubscriptionRepository(),
			s.ContentManager(),
		)
	}

	return s.addUsecase
}

func (s *serviceProvider) RSSUseCase() *usecase.RSSUseCase {
	if s.rssUseCase == nil {
		s.rssUseCase = usecase.NewRSSUseCase(
			s.PlaylistRepository(),
			s.EpisodeRepository(),
		)
	}

	return s.rssUseCase
}
