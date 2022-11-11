package file_manager

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"mvdan.cc/xurls/v2"
)

type YTDLPLoader struct{}

const (
	loadCmd = "yt-dlp -x --audio-format=mp3 --audio-quality=0 -f m4a/bestaudio --write-info-json --no-progress -o %s.tmp %s"
	destPath = "./storage/downloads/"
	infoExt  = ".tmp.info.json"
	checkCmd = "yt-dlp %s --skip-download --write-info-json --no-write-playlist-metafiles --dateafter %s"
	titleFilter = "--match-filters title~='%s'"
)

func (l YTDLPLoader) DownloadWithInfo(ctx context.Context, url string) (file localFile, err error) {
	file = localFile{}
	id := uuid.New().String()

	// Load audio with metadata
	cmdStr := fmt.Sprintf(loadCmd, id, url)
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Dir = destPath

	log.Printf("[DEBUG] executing command: %s", cmd.String())
	if err = cmd.Run(); err != nil {
		log.Printf("[ERROR] failed to execute command: %v", err)
		return
	}

	file.path = filepath.Join(destPath, id+".mp3")

	// Parse image, title and description
	infoPath := filepath.Join(destPath, id+infoExt)
	info, err := parseInfo(infoPath)
	if err != nil {
		log.Printf("[ERROR] failed to open info json: %v", err)
		return
	}

	// Sanitize description
	desc := info.Description
	furls := xurls.Relaxed().FindAllString(desc, -1)
	for _, u := range furls {
		desc = strings.ReplaceAll(desc, u, "")
	}
	info.Description = desc

	file.info = info

	// Remove local info json
	err = os.Remove(infoPath)
	if err != nil {
		log.Printf("[ERROR] failed to delete info json, %v", err)
	}

	return
}

func (l YTDLPLoader) DownloadUpdates(ctx context.Context, url string, date string)  (files []localFile, err error) {

	return
}

// All filters, separate json files
// yt-dlp --dateafter 20200520 --match-filters title~='MOUNTAIN BIKE' https://www.youtube.com/playlist\?list\=PLWx61XgoQmqdkfWC58_sYKAZvdQt9eBxQ --write-info-json --skip-download --no-write-playlist-metafiles

// For filter by title
// --match-filters title~='{title}'

// For filter by date
// --dateafter "YYYYMMDD"