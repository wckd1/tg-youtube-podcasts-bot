package telegram

import (
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrInitConnection = errors.New("can't init bot api")
	ErrContextClosed  = errors.New("context closed")
	ErrChannelClosed  = errors.New("telegram update channel closed")
	ErrResponseSend   = errors.New("can't send message to telegram")
)

type TelegramListener struct {
	botAPI   *tgbotapi.BotAPI
	commands commandList
}

func NewTelegramListener(token string, debugMode bool) (*TelegramListener, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Join(ErrInitConnection, err)
	}
	botAPI.Debug = debugMode

	return &TelegramListener{
		botAPI: botAPI,
	}, nil
}

func (l *TelegramListener) RegisterCommands(cs ...Command) {
	var cl commandList

	for _, c := range cs {
		cl = append(cl, c)
	}

	l.commands = cl
}

func (l TelegramListener) Start() {
	log.Println("[INFO] starting telegram listener...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := l.botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := l.transform(update.Message)
		resp := l.commands.OnMessage(*msg)
		if err := l.sendBotResponse(update.Message.Chat.ID, resp); err != nil {
			log.Printf("[WARN] %+v", err)
		}
	}
}

func (l TelegramListener) Shutdown() {
	l.botAPI.StopReceivingUpdates()
	log.Println("[INFO] telegram listener stopped")
}

func (l TelegramListener) transform(msg *tgbotapi.Message) *Message {
	return &Message{
		ID:        msg.MessageID,
		ChatID:    msg.Chat.ID,
		Command:   msg.Command(),
		Arguments: msg.CommandArguments(),
	}
}

func (l TelegramListener) sendBotResponse(chatID int64, resp Response) error {
	if !resp.Send {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, resp.Text)

	if _, err := l.botAPI.Send(msg); err != nil {
		return errors.Join(ErrResponseSend, err)
	}

	return nil
}
