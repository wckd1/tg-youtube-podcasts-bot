package repository

import (
	"errors"
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/converter"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/storage/bbolt"

	bolt "go.etcd.io/bbolt"
)

const subscriptionsBucketName = "subscriptions"

var _ subscription.SubscriptionRepository = (*SubscriptionRepository)(nil)

type SubscriptionRepository struct {
	store *bbolt.BBoltStore
}

func NewSubscriptionRepository(store *bbolt.BBoltStore) subscription.SubscriptionRepository {
	return &SubscriptionRepository{store}
}

func (r SubscriptionRepository) SaveSubsctiption(sub *subscription.Subscription) error {
	return r.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(subscriptionsBucketName))
		if err != nil {
			return err
		}

		subData, err := converter.SubscriptionToBinary(sub)
		if err != nil {
			return errors.Join(subscription.ErrSubscriptionEncoding, err)
		}
		return b.Put([]byte(sub.ID()), subData)
	})
}

func (r SubscriptionRepository) GetSubscription(id string) (subscription.Subscription, error) {
	var sub subscription.Subscription

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(subscriptionsBucketName))
		if b == nil {
			return subscription.ErrNoSubscriptionsStorage
		}

		subData := b.Get([]byte(id))
		if subData == nil {
			return subscription.ErrSubscriptionNotFound
		}

		decodedSub, err := converter.BinaryToSubscription(subData)
		if err != nil {
			return errors.Join(subscription.ErrSubscriptionDecoding, err)
		}
		sub = decodedSub
		return nil
	})

	return sub, err
}

func (r SubscriptionRepository) GetSubscriptions() ([]subscription.Subscription, error) {
	var result []subscription.Subscription

	err := r.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(subscriptionsBucketName))
		if b == nil {
			return subscription.ErrNoSubscriptionsStorage
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			s, err := converter.BinaryToSubscription(v)
			if err != nil {
				log.Printf("[WARN] failed to unmarshal, %+v", errors.Join(subscription.ErrSubscriptionDecoding, err))
				continue
			}
			result = append(result, s)
		}
		return nil
	})

	return result, err
}
