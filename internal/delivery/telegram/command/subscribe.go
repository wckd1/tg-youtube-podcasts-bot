package command

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"wckd1/tg-youtube-podcasts-bot/internal/converter"
	commandparser "wckd1/tg-youtube-podcasts-bot/internal/delivery/command_parser"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"
	"wckd1/tg-youtube-podcasts-bot/utils"
)

var _ telegram.Command = (*subscribe)(nil)

type subscribe struct {
	addUsecase          *usecase.AddUsecase
	subscriptionUsecase *usecase.SubscriptionUsecase
}

func NewSubscribeCommand(addUsecase *usecase.AddUsecase, subscriptionUsecase *usecase.SubscriptionUsecase) telegram.Command {
	return subscribe{addUsecase, subscriptionUsecase}
}

// OnMessage return new subscription status
func (s subscribe) OnMessage(msg telegram.Message) telegram.Response {
	if !utils.Contains(s.ReactOn(), msg.Command) {
		return telegram.Response{}
	}

	userID := strconv.Itoa(int(msg.ChatID))

	// Parse command arguments
	args, err := commandparser.ParseSubscribeArguments(msg.Arguments)
	if err != nil {
		log.Printf("[ERROR] can't parse subscription arguments, %+v", err)
		return telegram.Response{
			Text: fmt.Sprintf("Can't execute command: %s", err.Error()),
			Send: true,
		}
	}

	// List all subscription
	playlist, ok := args[commandparser.SubPlaylistKey]
	if !ok {
		playlist = ""
	}

	if len(args) == 0 || len(args) == 1 && playlist != "" {
		subs, err := s.subscriptionUsecase.ListSubscriptions(userID, playlist)
		if err != nil {
			log.Printf("[ERROR] can't list subscriptions, %+v", err)
			return telegram.Response{
				Text: "Something went wrong",
				Send: true,
			}
		}

		if len(subs) == 0 {
			return telegram.Response{
				Text: "No subscriptions found",
				Send: true,
			}
		}

		subsStrings := make([]string, 0)
		for _, sub := range subs {
			subsStrings = append(subsStrings, converter.SubscriptionToString(&sub))
		}

		return telegram.Response{
			Text: strings.Join(subsStrings, "\n"),
			Send: true,
		}
	}

	// Create new subscription
	id := args[commandparser.SubIDKey]
	url := args[commandparser.SubURLKey]
	filter, ok := args[commandparser.SubFilterKey]
	if !ok {
		filter = ""
	}

	if err := s.addUsecase.AddSubscription(userID, id, url, playlist, filter); err != nil {
		log.Printf("[ERROR] failed to add subscription. %+v", err)
		return telegram.Response{
			Text: "Failed to add subscription",
			Send: true,
		}
	}

	return telegram.Response{
		Text: "Subscription added",
		Send: true,
	}
}

// ReactOn keys
func (s subscribe) ReactOn() []string {
	return []string{"sub"}
}
