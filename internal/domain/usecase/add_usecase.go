package usecase

import (
	"context"
	"errors"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/playlist"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/service"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
)

var (
	ErrEpisodeCreate      = errors.New("can't add episode")
	ErrSubscriptionCreate = errors.New("can't add subscription")
)

type AddUsecase struct {
	userRepository         user.UserRepository
	playlistRepository     playlist.PlaylistRepository
	episodeRepository      episode.EpisodeRepository
	subscriptionRepository subscription.SubscriptionRepository
	contentManager         service.ContentManager
}

func NewAddUsecase(
	userRepository user.UserRepository,
	playlistRepository playlist.PlaylistRepository,
	episodeRepository episode.EpisodeRepository,
	subscriptionRepository subscription.SubscriptionRepository,
	contentManager service.ContentManager,
) *AddUsecase {
	return &AddUsecase{userRepository, playlistRepository, episodeRepository, subscriptionRepository, contentManager}
}

func (uc AddUsecase) AddEpisode(userID, id, url string) error {
	// Get user's default playlist
	user, err := uc.userRepository.GetUser(userID)
	if err != nil {
		return err
	}

	defaultPlaylist, err := uc.playlistRepository.GetPlaylist(user.DefaultPlaylist())
	if err != nil {
		return err
	}

	// Check if episode is already in playlist
	for _, epID := range defaultPlaylist.Episodes() {
		if epID == id {
			return nil
		}
	}

	// Ensure episode saved
	if err := uc.fetchEpisodeIfNeeded(id, url); err != nil {
		return err
	}

	// Add episode to default playlist
	defaultPlaylist.AddEpisode(id)
	err = uc.playlistRepository.SavePlaylist(&defaultPlaylist)
	if err != nil {
		return err
	}

	return nil
}

func (uc AddUsecase) CreateSubscription(userID, id, url, filter string) error {
	// Get subscription if already exist
	sub, err := uc.subscriptionRepository.GetSubscription(id)

	if err != nil {
		// Create subscribtion if not exist
		if errors.Is(err, subscription.ErrSubscriptionNotFound) || errors.Is(err, subscription.ErrNoSubscriptionsStorage) {
			sub = subscription.NewSubscription(id, url, filter, time.Now())
			// Save episode to database
			err = uc.subscriptionRepository.SaveSubsctiption(&sub)
			if err != nil {
				return errors.Join(ErrSubscriptionCreate, err)
			}
		} else {
			return errors.Join(ErrSubscriptionCreate, err)
		}
	}

	// Get user
	user, err := uc.userRepository.GetUser(userID)
	if err != nil {
		return err
	}

	defaultPlaylist, err := uc.playlistRepository.GetPlaylist(user.DefaultPlaylist())
	if err != nil {
		return err
	}

	// Check if subscription is already in playlist
	for _, subID := range defaultPlaylist.Subscriptions() {
		if subID == id {
			return nil
		}
	}

	defaultPlaylist.AddSubscription(sub.ID())

	// Save updated user
	err = uc.playlistRepository.SavePlaylist(&defaultPlaylist)
	if err != nil {
		return errors.Join(ErrSubscriptionCreate, err)
	}

	return nil
}

func (uc AddUsecase) fetchEpisodeIfNeeded(id, url string) error {
	// Check if episode already exists
	err := uc.episodeRepository.CheckExist(id)

	if err != nil {
		// Fetch and save episode if not exists
		if errors.Is(err, episode.ErrEpisodeNotFound) || errors.Is(err, episode.ErrNoEpisodesStorage) {
			// Content manager get episode
			ctx := context.Background()
			ep, err := uc.contentManager.Get(ctx, url)
			if err != nil {
				return errors.Join(ErrEpisodeCreate, err)
			}

			// Save episode to database
			err = uc.episodeRepository.SaveEpisode(&ep)
			if err != nil {
				return errors.Join(ErrEpisodeCreate, err)
			}
		} else {
			return errors.Join(ErrEpisodeCreate, err)
		}
	}

	return nil
}
