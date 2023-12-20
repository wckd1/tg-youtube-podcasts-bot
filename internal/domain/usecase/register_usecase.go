package usecase

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/repository"

	"github.com/google/uuid"
)

var (
	ErrUserExist      = errors.New("user already registered")
	ErrUserCreate     = errors.New("can't create user")
	ErrUserDelete     = errors.New("can't delete user")
	ErrPlaylistCreate = errors.New("can't createa playlist")
	ErrPlaylistDelete = errors.New("can't delete playlist")
)

type RegisterUsecase struct {
	userRepository     repository.UserRepository
	playlistRepository repository.PlaylistRepository
}

func NewRegisterUsecase(
	userRepository repository.UserRepository,
	playlistRepository repository.PlaylistRepository,
) *RegisterUsecase {
	return &RegisterUsecase{userRepository, playlistRepository}
}

func (uc RegisterUsecase) RegisterUser(id string) error {
	// Check if user not registered
	_, err := uc.userRepository.GetUser(id)
	if err == nil {
		return errors.Join(ErrUserExist, err)
	}

	// Create default playlist
	plID := uuid.NewString()
	playlist := entity.NewPlaylist(plID, entity.DefaultPlaylistName, []string{}, []string{})
	if err = uc.playlistRepository.SavePlaylist(&playlist); err != nil {
		return errors.Join(ErrPlaylistCreate, err)
	}

	// Create new user
	user := entity.NewUser(id, playlist.ID(), []string{})

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
