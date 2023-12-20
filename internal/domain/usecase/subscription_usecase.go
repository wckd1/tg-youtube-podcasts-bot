package usecase

import (
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/repository"

	"github.com/google/uuid"
)

type SubscriptionUsecase struct {
	userRepository         repository.UserRepository
	playlistRepository     repository.PlaylistRepository
	subscriptionRepository repository.SubscriptionRepository
}

func NewSubscriptionUsecase(
	userRepository repository.UserRepository,
	playlistRepository repository.PlaylistRepository,
	subscriptionRepository repository.SubscriptionRepository,
) *SubscriptionUsecase {
	return &SubscriptionUsecase{userRepository, playlistRepository, subscriptionRepository}
}

func (uc SubscriptionUsecase) ListSubscriptions(userID, pl string) ([]entity.Subscription, error) {
	subs := make([]entity.Subscription, 0)

	// Get target playlist or default
	playlist, err := uc.getTargetPlaylist(userID, pl)
	if err != nil {
		return subs, err
	}

	// Get playlist's subscriptions
	for _, subID := range playlist.Subscriptions() {
		sub, err := uc.subscriptionRepository.GetSubscription(subID)
		if err != nil {
			log.Printf("[ERROR] can't get subscription, %+v", err)
			continue
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func (uc SubscriptionUsecase) getTargetPlaylist(userID, pl string) (entity.Playlist, error) {
	var playlist entity.Playlist

	if pl != "" {
		_, err := uuid.Parse(pl)
		// Playlist name is passed
		if err != nil {
			p, err := uc.playlistRepository.GetPlaylistByName(pl)
			if err != nil {
				return playlist, ErrPlaylistNotFound
			}
			playlist = p
		} else {
			// Playlist id is passed
			p, err := uc.playlistRepository.GetPlaylist(pl)
			if err != nil {
				return playlist, ErrPlaylistNotFound
			}
			playlist = p
		}
	} else {
		// Get user's default playlist
		user, err := uc.userRepository.GetUser(userID)
		if err != nil {
			return playlist, err
		}

		defaultPl, err := uc.playlistRepository.GetPlaylist(user.DefaultPlaylist())
		if err != nil {
			return playlist, err
		}
		playlist = defaultPl
	}

	return playlist, nil
}
