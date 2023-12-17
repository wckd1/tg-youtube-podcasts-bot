package command

import (
	"fmt"
	"log"
	"strconv"
	commandparser "wckd1/tg-youtube-podcasts-bot/internal/delivery/command_parser"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"
	"wckd1/tg-youtube-podcasts-bot/utils"
)

var _ telegram.Command = (*subscribe)(nil)

type subscribe struct {
	addUsecase usecase.AddUsecase
}

func NewSubscribeCommand(addUsecase usecase.AddUsecase) telegram.Command {
	return subscribe{addUsecase}
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

	id := args[commandparser.SubIDKey]
	url := args[commandparser.SubURLKey]
	playlist, ok := args[commandparser.SubPlaylistKey]
	if !ok {
		playlist = ""
	}
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
