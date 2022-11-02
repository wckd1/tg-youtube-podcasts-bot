package handlers

import (
	"context"
	"fmt"
	"log"
	"sync"
	"wckd1/tg-youtube-podcasts-bot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramListener struct {
	BotAPI   *tgbotapi.BotAPI
	Commands bot.Command
	ChatID   int64

	msgs struct {
		once sync.Once
		ch   chan bot.Response
	}
}

// Process events
func (l *TelegramListener) Start(ctx context.Context) error {
	l.msgs.once.Do(func() {
		l.msgs.ch = make(chan bot.Response, 100)
	})

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

		case resp := <-l.msgs.ch: // publish messages from outside clients
			if err := l.sendBotResponse(resp, l.ChatID); err != nil {
				log.Printf("[WARN] failed to respond on submitted request, %v", err)
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

func (l *TelegramListener) Submit(ctx context.Context, text string) error {
	l.msgs.once.Do(func() { l.msgs.ch = make(chan bot.Response, 100) })

	select {
	case <-ctx.Done():
		return ctx.Err()
	case l.msgs.ch <- bot.Response{Text: text, Send: true}:
	}

	return nil
}

func (l *TelegramListener) transform(msg *tgbotapi.Message) *bot.Message {
	message := bot.Message{
		ID:        msg.MessageID,
		Command:   msg.Command(),
		Arguments: msg.CommandArguments(),
	}

	return &message
}
