package repository

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/converter"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/repository"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"

	bolt "go.etcd.io/bbolt"
)

const usersBucketName = "users"

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	store *bbolt.BBoltStore
}

func NewUserRepository(store *bbolt.BBoltStore) repository.UserRepository {
	return &UserRepository{store}
}

func (r UserRepository) GetUser(id string) (entity.User, error) {
	var u entity.User

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucketName))
		if b == nil {
			return repository.ErrNoUsersStorage
		}

		userData := b.Get([]byte(id))
		if userData == nil {
			return repository.ErrUserNotFound
		}

		decodedUsed, err := converter.BinaryToUser(userData)
		if err != nil {
			return errors.Join(repository.ErrUserDecoding, err)
		}
		u = decodedUsed
		return nil
	})

	return u, err
}

func (r UserRepository) SaveUser(u *entity.User) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(usersBucketName))
		if err != nil {
			return err
		}

		subData, err := converter.UserToBinary(u)
		if err != nil {
			return errors.Join(repository.ErrUserEncoding, err)
		}
		return b.Put([]byte(u.ID()), subData)
	})
}

func (r UserRepository) DeleteUser(id string) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucketName))
		if b == nil {
			return repository.ErrNoUsersStorage
		}

		return b.Delete([]byte(id))
	})
}
