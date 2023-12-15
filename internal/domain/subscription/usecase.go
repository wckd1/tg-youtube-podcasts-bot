package subscription

import (
	"errors"
	"time"
)

var (
	ErrGetSubscriptions  = errors.New("can't get subscriptions")
	ErrNoUpdatesRequired = errors.New("no updates are required")
)

type SubscriptionUsecase struct {
	subscriptionRepository SubscriptionRepository
}

func NewSubscriptionUsecase(subscriptionRepository SubscriptionRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{subscriptionRepository}
}

func (uc SubscriptionUsecase) CreateSubscription() error {
	return nil
}

func (uc SubscriptionUsecase) RemoveSubscription() error {
	return nil
}

func (uc SubscriptionUsecase) GetPendingSubscriptions() ([]Subscription, error) {
	pSubs := make([]Subscription, 0)

	// Get all subscriptions
	subs, err := uc.subscriptionRepository.GetSubscriptions()
	if err != nil {
		return pSubs, errors.Join(ErrGetSubscriptions, err)
	}

	// Filter only that needs to be update
	now := time.Now()

	for _, sub := range subs {
		// Calculate next update time for subscription
		updt := sub.LastUpdated().Add(time.Hour * 2)

		if updt.Before(now) || updt.Equal(now) {
			pSubs = append(pSubs, sub)
		}
	}

	// Check if any pending
	if len(pSubs) == 0 {
		return pSubs, ErrNoUpdatesRequired
	}

	return pSubs, nil
}
