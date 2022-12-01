package file_manager

import (
	"context"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramUploader struct {
	BotAPI *tgbotapi.BotAPI
	ChatID int64
}

func (u TelegramUploader) Upload(ctx context.Context, file localFile) (url string, err error) {
	// Config message
	msg := tgbotapi.NewAudio(u.ChatID, tgbotapi.FilePath(file.path))
	msg.Thumb = tgbotapi.FileURL(file.info.ImageURL)
	msg.Title = file.info.Title
	msg.Caption = strings.Join(
		[]string{"<b>" + file.info.Title + "</b>", file.info.Description},
		"\n\n",
	)
	msg.ParseMode = tgbotapi.ModeHTML

	// Send audio
	log.Printf("[DEBUG] uploading file: %s", file.path)
	resp, err := u.BotAPI.Send(msg)
	if err != nil {
		log.Printf("[ERROR] can't send message to telegram: %w", err)
		return
	}

	// Get direct link to file
	url, err = u.BotAPI.GetFileDirectURL(resp.Audio.FileID)
	if err != nil {
		log.Printf("[ERROR] failed to upload, %v", err)
		return
	}

	// Remove local audio
	err = os.Remove(file.path)
	log.Printf("[ERROR] failed to delete audio file, %v", err)

	return
}
