package command

import (
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
)

type add struct {
	subscriptionUsecase *subscription.SubscriptionUseCase
}

func NewAddCommand(subUC *subscription.SubscriptionUseCase) telegram.Command {
	return add{subscriptionUsecase: subUC}
}

// OnMessage return new subscription status
func (a add) OnMessage(msg telegram.Message) telegram.Response {
	if !contains(a.ReactOn(), msg.Command) {
		return telegram.Response{}
	}

	if err := a.subscriptionUsecase.CreateSubscription(); err != nil {
		log.Printf("[ERROR] failed to add subscription. %+v", err)
		return telegram.Response{
			Text: "Failed to add subscription",
			Send: true,
		}
	}

	return telegram.Response{
		Text: "Subscribed",
		Send: true,
	}
}

// ReactOn keys
func (a add) ReactOn() []string {
	return []string{"add"}
}
