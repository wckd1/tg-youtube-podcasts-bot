package app

import (
	"context"
	"wckd1/tg-youtube-podcasts-bot/configs"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/rss"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt/repository"
)

type serviceProvider struct {
	ctx    context.Context
	config configs.Config
	store  *bbolt.BBoltStore

	// Subscription
	subscriptionRepository subscription.SubscriptionRepository
	subscriptionUseCase    *subscription.SubscriptionUseCase

	// Episode
	episodeRepository episode.EpisodeRepository
	episodeUseCase    *episode.EpisodeUseCase

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

// Subscription
func (s *serviceProvider) SubscriptionRepository() subscription.SubscriptionRepository {
	if s.subscriptionRepository == nil {
		s.subscriptionRepository = repository.NewSubscriptionRepository(s.Store())
	}

	return s.subscriptionRepository
}
func (s *serviceProvider) SubscriptionUseCase() *subscription.SubscriptionUseCase {
	if s.subscriptionUseCase == nil {
		s.subscriptionUseCase = subscription.NewSubscriptionUseCase(s.SubscriptionRepository())
	}

	return s.subscriptionUseCase
}

// Episode
func (s *serviceProvider) EpisodeRepository() episode.EpisodeRepository {
	if s.episodeRepository == nil {
		s.episodeRepository = repository.NewEpisodeRepository(s.Store())
	}

	return s.episodeRepository
}
func (s *serviceProvider) EpisodeUseCase() *episode.EpisodeUseCase {
	if s.episodeUseCase == nil {
		s.episodeUseCase = episode.NewEpisodeUseCase(s.EpisodeRepository())
	}

	return s.episodeUseCase
}

// RSS
func (s *serviceProvider) RSSUseCase() *rss.RSSUseCase {
	if s.rssUseCase == nil {
		s.rssUseCase = rss.NewRSSUseCase()
	}

	return s.rssUseCase
}
