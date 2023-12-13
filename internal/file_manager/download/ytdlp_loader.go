package download_file_manager

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-pkgz/syncs"
	"github.com/google/uuid"
)

type YTDLPLoader struct{}

const (
	baseCmd    = "yt-dlp -x --audio-format=mp3 --audio-quality=0 -f m4a/bestaudio --write-info-json --no-progress %s"
	loadArgs   = "-o %s.tmp"
	updateArgs = "--no-write-playlist-metafiles --playlist-end 10 --dateafter %s -P \"%s\""
	filterArgs = "--match-filters title~='%s'"

	dlPath  = "./storage/downloads/"
	infoExt = ".info.json"
)

func (l YTDLPLoader) Download(ctx context.Context, url string) (file localFile, err error) {
	file = localFile{}
	id := uuid.New().String()

	// Prepare yt-dlp command
	fcmd := fmt.Sprintf(baseCmd, url)
	args := fmt.Sprintf(loadArgs, id)
	cmdStr := strings.Join([]string{fcmd, args}, " ")

	// Load audio with metadata
	if err = executeCommand(ctx, cmdStr); err != nil {
		log.Printf("[ERROR] failed to execute command: %v", err)
		return
	}

	file.path = filepath.Join(dlPath, id+".mp3")

	// Parse image, title and description
	infoPath := filepath.Join(dlPath, id+".tmp"+infoExt)
	file.info, err = parseInfo(infoPath)
	if err != nil {
		return
	}

	return
}

func (l YTDLPLoader) DownloadUpdates(ctx context.Context, url string, date time.Time, filter string) (files []localFile, err error) {
	id := uuid.New().String()

	// Prepare yt-dlp command
	fcmd := fmt.Sprintf(baseCmd, url)
	args := fmt.Sprintf(
		updateArgs,
		date.Format("20060102"), // For filter by date --dateafter "YYYYMMDD" // Maybe -1
		"./"+id,                 // Output directory
	)
	if len(filter) > 0 {
		args = args + " " + fmt.Sprintf(filterArgs, filter)
	}
	cmdStr := strings.Join([]string{fcmd, args}, " ")

	// Load audio with metadata
	if err = executeCommand(ctx, cmdStr); err != nil {
		log.Printf("[ERROR] failed to execute command: %v", err)
		return
	}

	// Get downloaded files
	dlFiles, err := filepath.Glob(filepath.Join(dlPath, id, "*.mp3"))
	if err != nil {
		log.Printf("[ERROR] failed to open source directory: %v", err)
		return
	}

	resps := make(chan localFile)
	wg := syncs.NewSizedGroup(len(dlFiles))

	for _, f := range dlFiles {
		wg.Go(func(ctx context.Context) {
			file := localFile{
				path: f,
			}

			infoPath := f[0:len(f)-4] + infoExt
			file.info, err = parseInfo(infoPath)
			if err != nil {
				return
			}

			resps <- file
		})
	}

	go func() {
		wg.Wait()
		close(resps)
	}()

	for r := range resps {
		files = append(files, r)
	}

	return
}

func executeCommand(ctx context.Context, cmdStr string) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Dir = dlPath

	log.Printf("[DEBUG] executing command: %s", cmd.String())
	return cmd.Run()
}
