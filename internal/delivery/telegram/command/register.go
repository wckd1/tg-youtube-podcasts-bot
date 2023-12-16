package command

import (
	"errors"
	"log"
	"strconv"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"
	"wckd1/tg-youtube-podcasts-bot/utils"
)

type register struct {
	registerUsecase *usecase.RegisterUsecase
}

func NewRegisterCommand(registerUsecase *usecase.RegisterUsecase) telegram.Command {
	return register{registerUsecase}
}

// OnMessage return new subscription status
func (r register) OnMessage(msg telegram.Message) telegram.Response {
	if !utils.Contains(r.ReactOn(), msg.Command) {
		return telegram.Response{}
	}

	id := strconv.Itoa(int(msg.ChatID))
	if err := r.registerUsecase.RegisterUser(id); err != nil {
		log.Printf("[ERROR] failed to register user. %+v", err)

		if errors.Is(err, usecase.ErrUserExist) {
			return telegram.Response{
				Text: "You are already registered",
				Send: true,
			}
		}
	}

	return telegram.Response{
		Text: "Succesfully registered",
		Send: true,
	}
}

// ReactOn keys
func (r register) ReactOn() []string {
	return []string{"reg"}
}
