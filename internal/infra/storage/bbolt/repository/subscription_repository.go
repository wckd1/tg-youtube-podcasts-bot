package repository

import (
	"errors"
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/converter"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/repository"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"

	bolt "go.etcd.io/bbolt"
)

const subscriptionsBucketName = "subscriptions"

var _ repository.SubscriptionRepository = (*SubscriptionRepository)(nil)

type SubscriptionRepository struct {
	store *bbolt.BBoltStore
}

func NewSubscriptionRepository(store *bbolt.BBoltStore) repository.SubscriptionRepository {
	return &SubscriptionRepository{store}
}

func (r *SubscriptionRepository) CheckExist(id string) error {
	return r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(subscriptionsBucketName))
		if b == nil {
			return repository.ErrNoSubscriptionsStorage
		}

		subData := b.Get([]byte(id))
		if subData == nil {
			return repository.ErrSubscriptionNotFound
		}

		return nil
	})
}

func (r SubscriptionRepository) SaveSubsctiption(sub *entity.Subscription) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(subscriptionsBucketName))
		if err != nil {
			return err
		}

		subData, err := converter.SubscriptionToBinary(sub)
		if err != nil {
			return errors.Join(repository.ErrSubscriptionEncoding, err)
		}
		return b.Put([]byte(sub.ID()), subData)
	})
}

func (r SubscriptionRepository) GetSubscription(id string) (entity.Subscription, error) {
	var sub entity.Subscription

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(subscriptionsBucketName))
		if b == nil {
			return repository.ErrNoSubscriptionsStorage
		}

		subData := b.Get([]byte(id))
		if subData == nil {
			return repository.ErrSubscriptionNotFound
		}

		decodedSub, err := converter.BinaryToSubscription(subData)
		if err != nil {
			return errors.Join(repository.ErrSubscriptionDecoding, err)
		}
		sub = decodedSub
		return nil
	})

	return sub, err
}

func (r SubscriptionRepository) GetSubscriptions() ([]entity.Subscription, error) {
	var result []entity.Subscription

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(subscriptionsBucketName))
		if b == nil {
			return repository.ErrNoSubscriptionsStorage
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			sub, err := converter.BinaryToSubscription(v)
			if err != nil {
				log.Printf("[WARN] failed to unmarshal, %+v", errors.Join(repository.ErrSubscriptionDecoding, err))
				continue
			}
			result = append(result, sub)
		}
		return nil
	})

	return result, err
}
