package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"wckd1/tg-youtube-podcasts-bot/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mvdan.cc/xurls/v2"
)

type YTLoader struct {
	Context   context.Context
	BotAPI    *tgbotapi.BotAPI
	chatID    int64
	Store     db.Store
	Submitter Submitter
}

const (
	ytdlpCmd = "yt-dlp -x --audio-format=mp3 --audio-quality=0 -f m4a/bestaudio --write-info-json --no-progress -o %s.tmp https://www.youtube.com/watch?v=%s"
	destPath = "./storage/downloads/"
	infoExt  = ".tmp.info.json"
)

func NewLoader(ctx context.Context, botAPI *tgbotapi.BotAPI, chatID int64, db db.Store, submitter Submitter) Interface {
	return &YTLoader{
		Context:   ctx,
		BotAPI:    botAPI,
		chatID:    chatID,
		Store:     db,
		Submitter: submitter,
	}
}

func (l YTLoader) Download(id string) {
	// Load audio with metadata
	cmdStr := fmt.Sprintf(ytdlpCmd, id, id)
	cmd := exec.CommandContext(l.Context, "sh", "-c", cmdStr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Dir = destPath

	log.Printf("[DEBUG] executing command: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] failed to execute command: %v", err)
		l.Submitter.SubmitText(l.Context, "failed to download")
		return
	}

	// Get image, title and description
	jsonInfo, err := os.Open(filepath.Join(destPath, id+infoExt))
	if err != nil {
		log.Printf("[ERROR] failed to open info json: %v", err)
		l.Submitter.SubmitText(l.Context, "failed to download")
		return
	}
	defer jsonInfo.Close()

	byteValue, _ := io.ReadAll(jsonInfo)

	var data struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ImageURL    string `json:"thumbnail"`
	}
	json.Unmarshal(byteValue, &data)

	// Sanitize description
	desc := data.Description
	furls := xurls.Relaxed().FindAllString(desc, -1)
	for _, u := range furls {
		desc = strings.ReplaceAll(desc, u, "")
	}
	data.Description = desc

	download := db.Download{
		URL:         filepath.Join(destPath, id+".mp3"),
		CoverURL:    data.ImageURL,
		Title:       data.Title,
		Description: data.Description,
	}

	l.Upload(download)
}
