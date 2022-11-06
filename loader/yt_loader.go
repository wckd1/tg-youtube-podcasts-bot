package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"wckd1/tg-youtube-podcasts-bot/db"
)

type YTLoader struct {
	Context   context.Context
	Store     db.Store
	Submitter Submitter
}

const (
	cmdArgs  = `-x --audio-format=mp3 --audio-quality=0 -f m4a/bestaudio "https://www.youtube.com/watch?v=%s" --write-info-json --no-progress -o %s.mp3`
	destPath = "./storage/downloads/"
	infoExt  = ".mp3.info.json"
)

func NewLoader(ctx context.Context, db db.Store, submitter Submitter) Interface {
	return &YTLoader{
		Context:   ctx,
		Store:     db,
		Submitter: submitter,
	}
}

func (l YTLoader) Download(id string) {
	// Load audio with metadata
	argsStr := fmt.Sprintf(cmdArgs, id, id)
	cmd := exec.CommandContext(l.Context, "yt-dlp", argsStr)
	cmd.Stdin = os.Stdin
	cmd.Dir = destPath

	log.Printf("[DEBUG] executing command: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] failed to execute command: %v", err)
		l.Submitter.SubmitText(l.Context, "failed to download")
	}

	// Get image, title and description
	jsonInfo, err := os.Open(filepath.Join(destPath, id+infoExt))
	if err != nil {
		log.Printf("[ERROR] failed to open info json: %v", err)
		l.Submitter.SubmitText(l.Context, "failed to download")
	}
	defer jsonInfo.Close()

	byteValue, _ := ioutil.ReadAll(jsonInfo)

	var data struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ImageURL    string `json:"thumbnail"`
	}
	json.Unmarshal(byteValue, &data)

	// TODO: Upload file to Telegram

	// Save uploaded file's info
	params := db.CreateDownloadParams{
		Path:        filepath.Join(destPath, id+".mp3"),
		CoverURL:    data.ImageURL,
		Title:       data.Title,
		Description: data.Description,
	}
	download, err := l.Store.CreateDownload(l.Context, params)
	if err != nil {
		log.Printf("[ERROR] failed to download, %v", err)
		l.Submitter.SubmitText(l.Context, "failed to download")
		return
	}

	l.Submitter.SubmitDowload(l.Context, download)
}
