package command

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	commandparser "wckd1/tg-youtube-podcasts-bot/internal/delivery/command_parser"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"
	"wckd1/tg-youtube-podcasts-bot/utils"
)

var (
	ErrNoURL           = errors.New("no source url provided")
	ErrParseURL        = errors.New("can't parse URL")
	ErrNotYoutubeURL   = errors.New("only youtube links are supported")
	ErrNotSupportedURL = errors.New("unrecognized link type")
)

var _ telegram.Command = (*add)(nil)

type add struct {
	addUsecase *usecase.AddUsecase
}

func NewAddCommand(addUsecase *usecase.AddUsecase) telegram.Command {
	return add{addUsecase}
}

// OnMessage return new subscription status
func (a add) OnMessage(msg telegram.Message) telegram.Response {
	if !utils.Contains(a.ReactOn(), msg.Command) {
		return telegram.Response{}
	}

	userID := strconv.Itoa(int(msg.ChatID))

	// Parse command arguments
	args, err := commandparser.ParseAddArguments(msg.Arguments)
	if err != nil {
		log.Printf("[ERROR] can't parse playlist arguments, %+v", err)
		return telegram.Response{
			Text: fmt.Sprintf("Can't execute command: %s", err.Error()),
			Send: true,
		}
	}

	id := args[commandparser.AddIDKey]
	url := args[commandparser.AddURLKey]
	playlist, ok := args[commandparser.AddPlaylistKey]
	if !ok {
		playlist = ""
	}

	if err := a.addUsecase.AddEpisode(userID, id, url, playlist); err != nil {
		log.Printf("[ERROR] failed to add episode. %+v", err)
		return telegram.Response{
			Text: "Failed to add episode",
			Send: true,
		}
	}

	return telegram.Response{
		Text: "Episode added",
		Send: true,
	}
}

// ReactOn keys
func (a add) ReactOn() []string {
	return []string{"add"}
}
