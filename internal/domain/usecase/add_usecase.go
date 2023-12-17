package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/playlist"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/service"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"

	"github.com/google/uuid"
)

var (
	ErrEpisodeCreate      = errors.New("can't add episode")
	ErrSubscriptionCreate = errors.New("can't add subscription")
	ErrPlaylistNotFound   = errors.New("playlist not found")
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

func (uc AddUsecase) AddEpisode(userID, id, url, pl string) error {
	// Get target playlist or default
	playlist, err := uc.getTargetPlaylist(userID, pl)
	if err != nil {
		return err
	}

	// Check if episode is already in playlist
	for _, epID := range playlist.Episodes() {
		if epID == id {
			return nil
		}
	}

	// Ensure episode saved
	if err := uc.fetchEpisodeIfNeeded(id, url); err != nil {
		return err
	}

	// Add episode to playlist
	playlist.AddEpisode(id)
	err = uc.playlistRepository.SavePlaylist(&playlist)
	if err != nil {
		return err
	}

	return nil
}

func (uc AddUsecase) AddSubscription(userID, id, url, pl string) error {
	// Get target playlist or default
	playlist, err := uc.getTargetPlaylist(userID, pl)
	if err != nil {
		return err
	}

	// Check if subscription is already in playlist
	for _, subID := range playlist.Subscriptions() {
		if subID == id {
			return nil
		}
	}

	// Get subscription if already exist
	err = uc.subscriptionRepository.CheckExist(id)

	if err != nil {
		// Create subscribtion if not exist
		if errors.Is(err, subscription.ErrSubscriptionNotFound) || errors.Is(err, subscription.ErrNoSubscriptionsStorage) {
			sub := subscription.NewSubscription(id, url, time.Now())
			// Save episode to database
			err = uc.subscriptionRepository.SaveSubsctiption(&sub)
			if err != nil {
				return errors.Join(ErrSubscriptionCreate, err)
			}
		} else {
			return errors.Join(ErrSubscriptionCreate, err)
		}
	}

	// Add subscription to playlist
	playlist.AddSubscription(id)
	err = uc.playlistRepository.SavePlaylist(&playlist)
	if err != nil {
		return errors.Join(ErrSubscriptionCreate, err)
	}

	return nil
}

func (uc AddUsecase) getTargetPlaylist(userID, pl string) (playlist.Playlist, error) {
	var playlist playlist.Playlist

	if pl != "" {
		plID, err := uuid.Parse(pl)
		// Playlist name is passed
		if err != nil {
			return playlist, fmt.Errorf("add to playlist by name is not supported yet") // TODO
		} else {
			// Playlist id is passed
			p, err := uc.playlistRepository.GetPlaylist(plID.String())
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
