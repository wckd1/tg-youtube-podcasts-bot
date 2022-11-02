package handlers

import (
	"context"
	"fmt"
	"log"
	"wckd1/tg-youtube-podcasts-bot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramListener struct {
	BotAPI   *tgbotapi.BotAPI
	Commands bot.Command
}

// Process events
func (l *TelegramListener) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := l.BotAPI.GetUpdatesChan(u)

	for {
		select {

		case <-ctx.Done():
			return ctx.Err()

		case update, ok := <-updates:
			if !ok {
				return fmt.Errorf("[INFO] telegram update channel closed")
			}

			if update.Message == nil { // ignore any non-Message updates
				continue
			}
			if !update.Message.IsCommand() { // ignore any non-command Messages
				continue
			}

			resp := l.Commands.OnMessage(*update.Message)
			if err := l.sendBotResponse(resp, update.Message.Chat.ID); err != nil {
				log.Printf("[WARN] failed to respond on update, %v", err)
			}
		}
	}
}

func (l *TelegramListener) sendBotResponse(resp bot.Response, chatID int64) error {
	if !resp.Send {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, resp.Text)
	// msg.ReplyToMessageID = update.Message.MessageID

	if _, err := l.BotAPI.Send(msg); err != nil {
		return fmt.Errorf("can't send message to telegram %q: %w", resp.Text, err)
	}

	return nil
}
