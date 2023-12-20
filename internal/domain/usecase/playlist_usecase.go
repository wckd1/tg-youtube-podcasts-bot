package usecase

import (
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/repository"

	"github.com/google/uuid"
)

type PlaylistUsecase struct {
	userRepository     repository.UserRepository
	playlistRepository repository.PlaylistRepository
}

func NewPlaylistUsecase(
	userRepository repository.UserRepository,
	playlistRepository repository.PlaylistRepository,
) *PlaylistUsecase {
	return &PlaylistUsecase{userRepository, playlistRepository}
}

func (uc PlaylistUsecase) ListPlaylists(userID string) ([]entity.Playlist, error) {
	pls := make([]entity.Playlist, 0)

	// Get user's default playlist
	user, err := uc.userRepository.GetUser(userID)
	if err != nil {
		return pls, err
	}

	defaultPlaylist, err := uc.playlistRepository.GetPlaylist(user.DefaultPlaylist())
	if err != nil {
		return pls, err
	}

	pls = append(pls, defaultPlaylist)

	// Get custom playlists
	for _, plID := range user.Playlists() {
		pl, err := uc.playlistRepository.GetPlaylist(plID)
		if err != nil {
			log.Printf("[ERROR] can't get playlist, %+v", err)
			continue
		}

		pls = append(pls, pl)
	}

	return pls, nil
}

func (uc PlaylistUsecase) CreatePlaylist(userID, name string) (*entity.Playlist, error) {
	// Get user
	user, err := uc.userRepository.GetUser(userID)
	if err != nil {
		return nil, err
	}

	// Create new playlist
	id := uuid.NewString()
	pl := entity.NewPlaylist(id, name, []string{}, []string{})

	if err := uc.playlistRepository.SavePlaylist(&pl); err != nil {
		return &pl, err
	}

	// Update user's playlists
	user.AddPlaylist(id)
	if err := uc.userRepository.SaveUser(&user); err != nil {
		log.Printf("[ERROR] can't update user, %+v", err)
		// Rollback playlist on user update fail
		if err := uc.playlistRepository.DeletePlaylist(id); err != nil {
			log.Printf("[ERROR] can't rollback playlist creation, %+v", err)
			return &pl, err
		}
		return &pl, err
	}

	return &pl, nil
}
