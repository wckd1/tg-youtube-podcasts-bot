package repository

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
)

var (
	ErrNoUsersStorage = errors.New("no saved users")
	ErrUserNotFound   = errors.New("user not found")
	ErrUserEncoding   = errors.New("can't encode user")
	ErrUserDecoding   = errors.New("can't decode user")
)

type UserRepository interface {
	SaveUser(user *entity.User) error
	GetUser(id string) (entity.User, error)
	DeleteUser(id string) error
}
