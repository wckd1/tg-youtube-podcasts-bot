package usecase

import (
	"errors"
	"log"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/repository"
)

var (
	ErrGetSubscriptions  = errors.New("can't get subscriptions")
	ErrNoUpdatesRequired = errors.New("no updates are required")
)

type UpdateUsecase struct {
	userRepository         repository.UserRepository
	playlistRepository     repository.PlaylistRepository
	episodeRepository      repository.EpisodeRepository
	subscriptionRepository repository.SubscriptionRepository
}

func NewUpdateUsecase(
	userRepository repository.UserRepository,
	playlistRepository repository.PlaylistRepository,
	episodeRepository repository.EpisodeRepository,
	subscriptionRepository repository.SubscriptionRepository,
) *UpdateUsecase {
	return &UpdateUsecase{userRepository, playlistRepository, episodeRepository, subscriptionRepository}
}

func (uc UpdateUsecase) GetPendingSubscriptions() ([]entity.Subscription, error) {
	pSubs := make([]entity.Subscription, 0)

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

func (uc UpdateUsecase) SaveEpisode(ep *entity.Episode, subID string) error {
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

func (uc UpdateUsecase) SaveSubsctiption(sub *entity.Subscription, lastUpdated time.Time) error {
	sub.SetLastUpdated(time.Now())
	return uc.subscriptionRepository.SaveSubsctiption(sub)
}
