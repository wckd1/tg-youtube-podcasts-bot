package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/content"

	"github.com/google/uuid"
)

var ErrNotAvailable = errors.New("yt-dlp not available")

var _ content.ContentManager = (*youTubeContentManager)(nil)

const (
	baseCmd    = "yt-dlp --skip-download --write-info-json --no-progress %s"
	loadArgs   = "-o %s.tmp"
	updateArgs = "--no-write-playlist-metafiles --playlist-end 10 --dateafter %s -P \"%s\""
	filterArgs = "--match-filters title~='%s'"

	audioCmd = "yt-dlp --get-url -f 140 %s"

	dlPath  = "./storage/downloads/"
	infoExt = ".info.json"
)

type youTubeContentManager struct{}

func NewYouTubeContentManager() (*youTubeContentManager, error) {
	cmd := exec.Command("which", "yt-dlp")
	err := cmd.Run()
	if err != nil {
		return nil, errors.Join(ErrNotAvailable, err)
	}

	return &youTubeContentManager{}, nil
}

func (cm youTubeContentManager) Get(ctx context.Context, url string) (content.Download, error) {
	dl := content.Download{
		URL:  url,
		Info: content.FileInfo{},
	}

	// Get link
	audio, err := getAudioLink(ctx, url)
	if err != nil {
		log.Printf("[ERROR] failed to get audio link: %v", err)
		return content.Download{}, err
	}
	dl.URL = audio

	// Get metadata
	info, err := getInfo(ctx, url)
	if err != nil {
		log.Printf("[ERROR] failed to get info: %v", err)
		return content.Download{}, err
	}
	dl.Info = info

	return dl, nil
}

func (cm youTubeContentManager) CheckUpdate(ctx context.Context, url string, date time.Time, filter string) (dls []content.Download, err error) {
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

	// Load metadata files
	if err = executeCommand(ctx, cmdStr); err != nil {
		log.Printf("[ERROR] failed to execute command: %v", err)
		return
	}

	// Get downloaded files
	dlFiles, err := filepath.Glob(filepath.Join(dlPath, id, "*"+infoExt))
	if err != nil {
		log.Printf("[ERROR] failed to open source directory: %v", err)
		return
	}

	resps := make(chan content.Download)
	wg := sync.WaitGroup{}
	wg.Add(len(dlFiles))

	for _, fp := range dlFiles {
		go func(path string) {
			defer wg.Done()

			dl := content.Download{
				URL:  url,
				Info: content.FileInfo{},
			}

			// Get metadata
			info, err := parseInfo(path)
			if err != nil {
				return
			}
			dl.Info = info

			// Get link
			audio, err := getAudioLink(ctx, info.Link)
			if err != nil {
				log.Printf("[ERROR] failed to get audio link: %v", err)
				return
			}
			dl.URL = audio

			resps <- dl
		}(fp)
	}

	go func() {
		wg.Wait()
		close(resps)
	}()

	for r := range resps {
		dls = append(dls, r)
	}

	return
}

func getInfo(ctx context.Context, url string) (content.FileInfo, error) {
	id := uuid.New().String()
	// Prepare yt-dlp command
	fcmd := fmt.Sprintf(baseCmd, url)
	args := fmt.Sprintf(loadArgs, id)
	cmdStr := strings.Join([]string{fcmd, args}, " ")

	// Load metadata
	if err := executeCommand(ctx, cmdStr); err != nil {
		return content.FileInfo{}, fmt.Errorf("failed to execute command: %v", err)
	}

	// Parse image, title and description
	infoPath := filepath.Join(dlPath, id+".tmp"+infoExt)
	info, err := parseInfo(infoPath)
	if err != nil {
		return content.FileInfo{}, fmt.Errorf("failed to parse metadata: %v", err)
	}

	return info, nil
}

func getAudioLink(ctx context.Context, url string) (string, error) {
	fcmd := fmt.Sprintf(audioCmd, url)
	cmd := exec.CommandContext(ctx, "sh", "-c", fcmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get output: %v", err)
	}

	return string(output), nil
}

func parseInfo(path string) (info content.FileInfo, err error) {
	jsonInfo, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("[ERROR] failed to open info json: %v", err)
		return
	}
	defer jsonInfo.Close()

	byteValue, _ := io.ReadAll(jsonInfo)
	json.Unmarshal(byteValue, &info)

	// Remove local info json
	err = os.Remove(path)
	if err != nil {
		err = fmt.Errorf("[ERROR] failed to delete info json, %v", err)
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
