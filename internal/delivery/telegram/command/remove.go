package command

import (
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"
)

type remove struct {
	subscriptionUsecase *subscription.SubscriptionUsecase
}

func NewRemoveCommand(subUC *subscription.SubscriptionUsecase) telegram.Command {
	return remove{subscriptionUsecase: subUC}
}

// OnMessage return deleted subscription status
func (r remove) OnMessage(msg telegram.Message) telegram.Response {
	if !contains(r.ReactOn(), msg.Command) {
		return telegram.Response{}
	}

	if err := r.subscriptionUsecase.RemoveSubscription(); err != nil {
		log.Printf("[ERROR] failed to remove subscription. %+v", err)
		return telegram.Response{
			Text: "Failed to remove subscription",
			Send: true,
		}
	}

	return telegram.Response{
		Text: "Unsubscribed",
		Send: true,
	}
}

// ReactOn keys
func (r remove) ReactOn() []string {
	return []string{"remove"}
}
