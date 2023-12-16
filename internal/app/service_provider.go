package app

import (
	"context"
	"log"
	"wckd1/tg-youtube-podcasts-bot/configs"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/playlist"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/rss"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/content/youtube"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt/repository"
)

type serviceProvider struct {
	ctx    context.Context
	config configs.Config

	store          *bbolt.BBoltStore
	contentManager episode.ContentManager

	// User
	userRepository user.UserRepository
	userUsecase    *user.UserUsecase

	// Playlist
	playlistRepository playlist.PlaylistRepository

	// Subscription
	subscriptionRepository subscription.SubscriptionRepository
	subscriptionUsecase    *subscription.SubscriptionUsecase

	// Episode
	episodeRepository episode.EpisodeRepository
	episodeUsecase    *episode.EpisodeUsecase

	// RSS
	rssUseCase *rss.RSSUseCase
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

func (s *serviceProvider) ContentManager() episode.ContentManager {
	if s.contentManager == nil {
		cm, err := youtube.NewYouTubeContentManager()
		if err != nil {
			log.Fatalf("[ERROR] can't init YouTubeContentManager, %+v", err)
		}
		s.contentManager = cm
	}

	return s.contentManager
}

// User
func (s *serviceProvider) UserRepository() user.UserRepository {
	if s.userRepository == nil {
		s.userRepository = repository.NewUserRepository(s.Store())
	}

	return s.userRepository
}
func (s *serviceProvider) UserUsecase() *user.UserUsecase {
	if s.userUsecase == nil {
		s.userUsecase = user.NewUserUsecase(s.UserRepository(), s.PlaylistRepository())
	}

	return s.userUsecase
}

// Playlist
func (s *serviceProvider) PlaylistRepository() playlist.PlaylistRepository {
	if s.playlistRepository == nil {
		s.playlistRepository = repository.NewPlaylistRepository(s.Store())
	}

	return s.playlistRepository
}

// Subscription
func (s *serviceProvider) SubscriptionRepository() subscription.SubscriptionRepository {
	if s.subscriptionRepository == nil {
		s.subscriptionRepository = repository.NewSubscriptionRepository(s.Store())
	}

	return s.subscriptionRepository
}
func (s *serviceProvider) SubscriptionUsecase() *subscription.SubscriptionUsecase {
	if s.subscriptionUsecase == nil {
		s.subscriptionUsecase = subscription.NewSubscriptionUsecase(s.SubscriptionRepository(), s.UserRepository())
	}

	return s.subscriptionUsecase
}

// Episode
func (s *serviceProvider) EpisodeRepository() episode.EpisodeRepository {
	if s.episodeRepository == nil {
		s.episodeRepository = repository.NewEpisodeRepository(s.Store())
	}

	return s.episodeRepository
}
func (s *serviceProvider) EpisodeUsecase() *episode.EpisodeUsecase {
	if s.episodeUsecase == nil {
		s.episodeUsecase = episode.NewEpisodeUsecase(
			s.EpisodeRepository(),
			s.UserRepository(),
			s.PlaylistRepository(),
			s.ContentManager(),
		)
	}

	return s.episodeUsecase
}

// RSS
func (s *serviceProvider) RSSUseCase() *rss.RSSUseCase {
	if s.rssUseCase == nil {
		s.rssUseCase = rss.NewRSSUseCase()
	}

	return s.rssUseCase
}
