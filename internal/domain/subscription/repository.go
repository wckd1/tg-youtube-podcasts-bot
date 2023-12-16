package subscription

import "errors"

var (
	ErrNoSubscriptionsStorage = errors.New("no saved subscriptions")
	ErrSubscriptionNotFound   = errors.New("subscription not found")
	ErrSubscriptionEncoding   = errors.New("can't encode subscription")
	ErrSubscriptionDecoding   = errors.New("can't decode subscription")
)

type SubscriptionRepository interface {
	SaveSubsctiption(sub *Subscription) error
	GetSubscription(id string) (Subscription, error)
	GetSubscriptions() ([]Subscription, error)
}
