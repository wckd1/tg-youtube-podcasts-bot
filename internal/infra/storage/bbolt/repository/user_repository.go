package repository

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/converter"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"

	bolt "go.etcd.io/bbolt"
)

const usersBucketName = "users"

var (
	ErrNoUsersBucket = errors.New("no saved users")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserEncoding  = errors.New("can't encode user")
	ErrUserDecoding  = errors.New("can't decode user")
)

var _ user.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	store *bbolt.BBoltStore
}

func NewUserRepository(store *bbolt.BBoltStore) user.UserRepository {
	return &UserRepository{store}
}

func (r UserRepository) GetUser(id string) (user.User, error) {
	var user user.User

	err := r.store.View(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte(usersBucketName))
		if usersBucket == nil {
			return ErrNoUsersBucket
		}

		userData := usersBucket.Get([]byte(id))
		if userData == nil {
			return ErrUserNotFound
		}

		u, err := converter.BinaryToUser(userData)
		if err != nil {
			return errors.Join(ErrUserDecoding, err)
		}
		user = u
		return nil
	})

	return user, err
}

func (r UserRepository) CreateUser(user *user.User) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(usersBucketName))
		if err != nil {
			return err
		}

		subData, err := converter.UserToBinary(user)
		if err != nil {
			return errors.Join(ErrUserEncoding, err)
		}
		return b.Put([]byte(user.ID()), subData)
	})
}

func (r UserRepository) DeleteUser(id string) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucketName))
		if b == nil {
			return ErrNoUsersBucket
		}

		return b.Delete([]byte(id))
	})
}
