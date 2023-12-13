package handlers

import (
	"context"
	"fmt"
	"log"
	"wckd1/tg-youtube-podcasts-bot/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	BotAPI      *tgbotapi.BotAPI
	Commands    bot.Command
	AdminChatID int64
}

// Process events
func (l *Telegram) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
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

			msg := l.transform(update.Message)
			resp := l.Commands.OnMessage(*msg)
			if err := l.sendBotResponse(resp, l.ChatID); err != nil {
				log.Printf("[WARN] failed to respond on update, %v", err)
			}
		}
	}
}

func (l *Telegram) sendBotResponse(resp bot.Response, chatID int64) error {
	if !resp.Send {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, resp.Text)

	if _, err := l.BotAPI.Send(msg); err != nil {
		return fmt.Errorf("can't send message to telegram %q: %w", resp.Text, err)
	}

	return nil
}

func (l *Telegram) transform(msg *tgbotapi.Message) *bot.Message {
	message := bot.Message{
		ID:        msg.MessageID,
		Command:   msg.Command(),
		Arguments: msg.CommandArguments(),
	}

	return &message
}
