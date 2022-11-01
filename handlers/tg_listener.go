package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mvdan.cc/xurls"
)

type TelegramListener struct {
	BotAPI *tgbotapi.BotAPI
	// TODO: Add Commands
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
				return fmt.Errorf("telegram update channel closed")
			}

			if update.Message == nil { // ignore any non-Message updates
				continue
			}
			if !update.Message.IsCommand() { // ignore any non-command Messages
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "add":
				args := update.Message.CommandArguments()
				url := xurls.Relaxed.FindString(args)
				title := strings.ReplaceAll(args, url+" ", "")
				msg.Text = fmt.Sprintf("URL: %s\nTitle: %s", url, title)
			default:
				msg.Text = "Unknown command"
			}

			msg.ReplyToMessageID = update.Message.MessageID

			if _, err := l.BotAPI.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
