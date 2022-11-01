package handlers

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateChecker struct {
	BotAPI *tgbotapi.BotAPI
}

// TODO: Remove hardcoded value
const chatID = -826459712

func (uc *UpdateChecker) Start(ctx context.Context, delay time.Duration) error {
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			msg := tgbotapi.NewMessage(chatID, "Checked for update")
			if _, err := uc.BotAPI.Send(msg); err != nil {
				log.Panic(err)
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
