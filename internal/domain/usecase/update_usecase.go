package usecase

import (
	"errors"
	"log"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/playlist"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
)

var (
	ErrGetSubscriptions  = errors.New("can't get subscriptions")
	ErrNoUpdatesRequired = errors.New("no updates are required")
)

type UpdateUsecase struct {
	userRepository         user.UserRepository
	playlistRepository     playlist.PlaylistRepository
	episodeRepository      episode.EpisodeRepository
	subscriptionRepository subscription.SubscriptionRepository
}

func NewUpdateUsecase(
	userRepository user.UserRepository,
	playlistRepository playlist.PlaylistRepository,
	episodeRepository episode.EpisodeRepository,
	subscriptionRepository subscription.SubscriptionRepository,
) *UpdateUsecase {
	return &UpdateUsecase{userRepository, playlistRepository, episodeRepository, subscriptionRepository}
}

func (uc UpdateUsecase) GetPendingSubscriptions() ([]subscription.Subscription, error) {
	pSubs := make([]subscription.Subscription, 0)

	// Get all subscriptions
	subs, err := uc.subscriptionRepository.GetSubscriptions()
	if err != nil {
		return pSubs, errors.Join(ErrGetSubscriptions, err)
	}

	// Filter only that needs to be update
	now := time.Now()

	for _, sub := range subs {
		// Calculate next update time for subscription
		updt := sub.LastUpdated().Add(time.Hour * 2)

		if updt.Before(now) || updt.Equal(now) {
			pSubs = append(pSubs, sub)
		}
	}

	// Check if any pending
	if len(pSubs) == 0 {
		return pSubs, ErrNoUpdatesRequired
	}

	return pSubs, nil
}

func (uc UpdateUsecase) SaveEpisode(ep *episode.Episode, subID string) error {
	// Save episode
	err := uc.episodeRepository.SaveEpisode(ep)
	if err != nil {
		return err
	}

	// Get playlists with subscription
	pls, err := uc.playlistRepository.GetPlaylistsWithSubscription(subID)
	if err != nil {
		return err
	}

	// Add episode to playlists
	for _, pl := range pls {
		pl.AddEpisode(ep.ID())

		if err := uc.playlistRepository.SavePlaylist(&pl); err != nil {
			log.Printf("[ERROR] can't add episode to playlist, %+v", err)
			continue
		}
	}

	return nil
}

func (uc UpdateUsecase) SaveSubsctiption(sub *subscription.Subscription, lastUpdated time.Time) error {
	sub.SetLastUpdated(time.Now())
	return uc.subscriptionRepository.SaveSubsctiption(sub)
}
