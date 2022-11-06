package loader

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"wckd1/tg-youtube-podcasts-bot/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (l YTLoader) Upload(download db.Download) {
	// Audio can't be mixed with other media types
	// So we will send 2 messages: with meta and with audio itself

	// Meta message
	metaMsg := tgbotapi.NewPhoto(l.chatID, tgbotapi.FileURL(download.CoverURL))
	metaMsg.Caption = strings.Join(
		[]string{"<b>" + download.Title + "</b>", download.Description},
		"\n\n",
	)
	metaMsg.ParseMode = tgbotapi.ModeHTML

	// Audio message
	audioMsg := tgbotapi.NewAudio(l.chatID, tgbotapi.FilePath(download.Path))

	// TODO: Run asynchroniously
	// Send meta
	if _, err := l.BotAPI.Send(metaMsg); err != nil {
		log.Printf("[ERROR] can't send message to telegram: %w", err)
		l.Submitter.SubmitText(l.Context, "failed to upload")
		return
	}

	if _, err := l.BotAPI.Send(audioMsg); err != nil {
		log.Printf("[ERROR] can't send message to telegram: %w", err)
		l.Submitter.SubmitText(l.Context, "failed to upload")
		return
	}

	// TODO: Find out direct link to get uploaded audio
	// {"file_id":"CQACAgIAAxkDAAIBtGNn5s6gqzqkez-aTo9DldWR5GB0AAK_IQACUeZBS_7lDPeEWLtuKwQ","file_unique_id":"AgADvyEAAlHmQUs","duration":0,"file_name":"7cImI9VWESw.mp3","mime_type":"audio/mpeg","file_size":18175055}

	// TODO: Update download.Path or maybe add URL field / rename Path to URL and add Uploaded:bool field
	// Save uploaded file's info
	if err := l.Store.CreateDownload(&download); err != nil {
		log.Printf("[ERROR] failed to download, %v", err)
		l.Submitter.SubmitText(l.Context, "failed to download")
		return
	}

	// Remove local info json
	err := os.Remove(filepath.Join(destPath, string(download.ID)+infoExt))
	log.Printf("[ERROR] failed to delete info json, %v", err)
	// Remove local audio
	err = os.Remove(filepath.Join(destPath, download.Path))
	log.Printf("[ERROR] failed to delete audio jsofilen, %v", err)
}
