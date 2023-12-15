package command

import (
	"log"
	"strconv"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
)

type register struct {
	userUsecase *user.UserUsecase
}

func NewRegisterCommand(uUC *user.UserUsecase) telegram.Command {
	return register{userUsecase: uUC}
}

// OnMessage return new subscription status
func (r register) OnMessage(msg telegram.Message) telegram.Response {
	if !contains(r.ReactOn(), msg.Command) {
		return telegram.Response{}
	}

	id := strconv.Itoa(int(msg.ChatID))
	if err := r.userUsecase.RegisterUser(id); err != nil {
		log.Printf("[ERROR] failed to register user. %+v", err)
		text := "Failed to register user"
		if err.Error() == user.ErrUserRegistered.Error() {
			text = text + ": already registered"
		}

		return telegram.Response{
			Text: text,
			Send: true,
		}
	}

	return telegram.Response{
		Text: "You are registered",
		Send: true,
	}
}

// ReactOn keys
func (r register) ReactOn() []string {
	return []string{"reg"}
}
