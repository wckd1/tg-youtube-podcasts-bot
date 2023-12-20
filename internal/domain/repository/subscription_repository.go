package repository

import (
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
)

var (
	ErrNoSubscriptionsStorage = errors.New("no saved subscriptions")
	ErrSubscriptionNotFound   = errors.New("subscription not found")
	ErrSubscriptionEncoding   = errors.New("can't encode subscription")
	ErrSubscriptionDecoding   = errors.New("can't decode subscription")
)

type SubscriptionRepository interface {
	CheckExist(id string) error
	SaveSubsctiption(sub *entity.Subscription) error
	GetSubscription(id string) (entity.Subscription, error)
	GetSubscriptions() ([]entity.Subscription, error)
}
