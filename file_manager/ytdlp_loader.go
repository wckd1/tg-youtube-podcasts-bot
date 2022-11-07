package file_manager

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

	"mvdan.cc/xurls/v2"
)

type YTDLPLoader struct{}

const (
	ytdlpCmd = "yt-dlp -x --audio-format=mp3 --audio-quality=0 -f m4a/bestaudio --write-info-json --no-progress -o %s.tmp https://www.youtube.com/watch?v=%s"
	destPath = "./storage/downloads/"
	infoExt  = ".tmp.info.json"
)

func (l YTDLPLoader) Download(ctx context.Context, id string) (file localFile, err error) {
	file = localFile{}

	// Load audio with metadata
	cmdStr := fmt.Sprintf(ytdlpCmd, id, id)
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Dir = destPath

	log.Printf("[DEBUG] executing command: %s", cmd.String())
	if err = cmd.Run(); err != nil {
		log.Printf("[ERROR] failed to execute command: %v", err)
		return
	}

	file.path = filepath.Join(destPath, id+".mp3")

	// Parse image, title and description
	jsonInfoPath := filepath.Join(destPath, id+infoExt)
	jsonInfo, err := os.Open(jsonInfoPath)
	if err != nil {
		log.Printf("[ERROR] failed to open info json: %v", err)
		return
	}
	defer jsonInfo.Close()

	byteValue, _ := io.ReadAll(jsonInfo)

	var info fileInfo
	json.Unmarshal(byteValue, &info)

	// Sanitize description
	desc := info.Description
	furls := xurls.Relaxed().FindAllString(desc, -1)
	for _, u := range furls {
		desc = strings.ReplaceAll(desc, u, "")
	}
	info.Description = desc

	file.info = info

	// Remove local info json
	err = os.Remove(jsonInfoPath)
	if err != nil {
		log.Printf("[ERROR] failed to delete info json, %v", err)
	}

	return
}
