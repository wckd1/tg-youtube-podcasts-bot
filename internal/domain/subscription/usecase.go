package subscription

import (
	"errors"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
)

var (
	ErrGetSubscriptions   = errors.New("can't get subscriptions")
	ErrNoUpdatesRequired  = errors.New("no updates are required")
	ErrSubscriptionCreate = errors.New("can't add subscription")
)

type SubscriptionUsecase struct {
	subscriptionRepository SubscriptionRepository
	userRepository         user.UserRepository
}

func NewSubscriptionUsecase(subR SubscriptionRepository, uR user.UserRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{subscriptionRepository: subR, userRepository: uR}
}

func (uc SubscriptionUsecase) CreateSubscription(userID, id, url, filter string) error {
	// Get subscription if already exist
	sub, err := uc.subscriptionRepository.GetSubscription(id)

	if err != nil {
		// Create subscribtion if not exist
		if errors.Is(err, ErrSubscriptionNotFound) || errors.Is(err, ErrNoSubscriptionsStorage) {
			sub = NewSubscription(id, url, filter, time.Now())
			// Save episode to database
			err = uc.subscriptionRepository.SaveSubsctiption(&sub)
			if err != nil {
				return errors.Join(ErrSubscriptionCreate, err)
			}
		} else {
			return errors.Join(ErrSubscriptionCreate, err)
		}
	}

	// Get user
	user, err := uc.userRepository.GetUser(userID)
	if err != nil {
		return err
	}
	user.AddSubscription(sub.ID())

	// Save updated user
	err = uc.userRepository.SaveUser(&user)
	if err != nil {
		return errors.Join(ErrSubscriptionCreate, err)
	}

	return nil
}

func (uc SubscriptionUsecase) SaveSubsctiption(sub *Subscription) error {
	return uc.subscriptionRepository.SaveSubsctiption(sub)
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
