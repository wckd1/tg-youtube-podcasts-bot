package subscription

type SubscriptionRepository interface {
	SaveSubsctiption(sub *Subscription) error
	GetSubscriptions() ([]Subscription, error)
	DeleteSubsctiption(sub *Subscription) error
}
