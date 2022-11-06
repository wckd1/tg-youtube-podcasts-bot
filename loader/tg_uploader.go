package loader

import (
	"log"
	"os"
	"strings"
	"wckd1/tg-youtube-podcasts-bot/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (l YTLoader) Upload(download db.Download) {
	// Config message
	audioMsg := tgbotapi.NewAudio(l.chatID, tgbotapi.FilePath(download.URL))
	audioMsg.Thumb = tgbotapi.FileURL(download.CoverURL)
	audioMsg.Caption = strings.Join(
		[]string{"<b>" + download.Title + "</b>", download.Description},
		"\n\n",
	)
	audioMsg.ParseMode = tgbotapi.ModeHTML

	// Send audio
	audioResp, err := l.BotAPI.Send(audioMsg)
	if err != nil {
		log.Printf("[ERROR] can't send message to telegram: %w", err)
		l.Submitter.SubmitText(l.Context, "failed to upload")
		return
	}

	// Get direct link to file
	fileURL, err := l.BotAPI.GetFileDirectURL(audioResp.Audio.FileID)
	if err != nil {
		log.Printf("[ERROR] failed to upload, %v", err)
		l.Submitter.SubmitText(l.Context, "failed to download")
		return
	}

	// Remove local info json
	err = os.Remove(strings.ReplaceAll(download.URL, ".mp3", infoExt))
	log.Printf("[ERROR] failed to delete info json, %v", err)
	// Remove local audio
	err = os.Remove(download.URL)
	log.Printf("[ERROR] failed to delete audio jsofilen, %v", err)

	// Change URL to uploaded link
	download.URL = fileURL

	// Save uploaded file's info
	if err := l.Store.CreateDownload(&download); err != nil {
		log.Printf("[ERROR] failed to upload, %v", err)
		l.Submitter.SubmitText(l.Context, "failed to download")
		return
	}
}
