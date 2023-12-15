package repository

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/converter"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"

	bolt "go.etcd.io/bbolt"
)

const usersBucketName = "users"

var _ user.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	store *bbolt.BBoltStore
}

func NewUserRepository(store *bbolt.BBoltStore) user.UserRepository {
	return &UserRepository{store}
}

func (r UserRepository) GetUser(id string) (user.User, error) {
	var u user.User

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucketName))
		if b == nil {
			return user.ErrNoUsersStorage
		}

		userData := b.Get([]byte(id))
		if userData == nil {
			return user.ErrUserNotFound
		}

		decodedUsed, err := converter.BinaryToUser(userData)
		if err != nil {
			return errors.Join(user.ErrUserDecoding, err)
		}
		u = decodedUsed
		return nil
	})

	return u, err
}

func (r UserRepository) SaveUser(u *user.User) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(usersBucketName))
		if err != nil {
			return err
		}

		subData, err := converter.UserToBinary(u)
		if err != nil {
			return errors.Join(user.ErrUserEncoding, err)
		}
		return b.Put([]byte(u.ID()), subData)
	})
}

func (r UserRepository) DeleteUser(id string) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucketName))
		if b == nil {
			return user.ErrNoUsersStorage
		}

		return b.Delete([]byte(id))
	})
}
