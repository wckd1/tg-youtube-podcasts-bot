package user

import "errors"

var (
	ErrNoUsersStorage = errors.New("no saved users")
	ErrUserNotFound   = errors.New("user not found")
	ErrUserEncoding   = errors.New("can't encode user")
	ErrUserDecoding   = errors.New("can't decode user")
)

type UserRepository interface {
	SaveUser(user *User) error
	GetUser(id string) (User, error)
	DeleteUser(id string) error
}
