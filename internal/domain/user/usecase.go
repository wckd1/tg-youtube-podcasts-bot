package user

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/playlist"

	"github.com/google/uuid"
)

var (
	ErrUserRegistered = errors.New("user already registered")
	ErrUserCreate     = errors.New("can't create user")
	ErrUserDelete     = errors.New("can't delete user")
	ErrPlaylistCreate = errors.New("can't createa playlist")
	ErrPlaylistDelete = errors.New("can't delete playlist")
)

type UserUsecase struct {
	userRepository     UserRepository
	playlistRepository playlist.PlaylistRepository
}

func NewUserUsecase(userRepository UserRepository, playlistRepository playlist.PlaylistRepository) *UserUsecase {
	return &UserUsecase{userRepository, playlistRepository}
}

func (uc UserUsecase) RegisterUser(id string) error {
	// Check if user not registered
	_, err := uc.userRepository.GetUser(id)
	if err == nil {
		return errors.Join(ErrUserRegistered, err)
	}

	// Create default playlist
	plID := uuid.NewString()
	playlist := playlist.NewPlaylist(plID, playlist.DefaultPlaylistName, make([]string, 0))
	if err = uc.playlistRepository.SavePlaylist(&playlist); err != nil {
		return errors.Join(ErrPlaylistCreate, err)
	}

	// Create new user
	user := NewUser(id, playlist.ID(), make([]string, 0), make([]string, 0))

	// Save user
	if err = uc.userRepository.SaveUser(&user); err != nil {
		// Roll back playlist creation
		if err = uc.playlistRepository.DeletePlaylist(id); err != nil {
			return errors.Join(ErrPlaylistDelete, ErrPlaylistCreate, err)
		}
		return errors.Join(ErrUserCreate, err)
	}

	if err = uc.playlistRepository.SavePlaylist(&playlist); err != nil {
		return errors.Join(ErrPlaylistCreate, err)
	}

	return nil
}
