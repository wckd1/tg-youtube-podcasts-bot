package episode

import (
	"context"
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/playlist"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
)

var (
	ErrEpisodeCreate = errors.New("can't add episode")
)

type EpisodeUsecase struct {
	episodeRepository  EpisodeRepository
	userRepository     user.UserRepository
	playlistRepository playlist.PlaylistRepository
	contentManager     ContentManager
}

func NewEpisodeUsecase(
	epR EpisodeRepository,
	uR user.UserRepository,
	plR playlist.PlaylistRepository,
	cm ContentManager,
) *EpisodeUsecase {
	return &EpisodeUsecase{episodeRepository: epR, userRepository: uR, playlistRepository: plR, contentManager: cm}
}

func (uc EpisodeUsecase) AddEpisode(userID, id, url string) error {
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

func (uc EpisodeUsecase) SaveEpisode(ep *Episode) error {
	return uc.episodeRepository.SaveEpisode(ep)
}

func (uc EpisodeUsecase) fetchEpisodeIfNeeded(id, url string) error {
	// Check if episode already exists
	err := uc.episodeRepository.CheckExist(id)

	if err != nil {
		// Fetch and save episode if not exists
		if errors.Is(err, ErrEpisodeNotFound) || errors.Is(err, ErrNoEpisodesStorage) {
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
