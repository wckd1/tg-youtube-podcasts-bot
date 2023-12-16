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
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/content"

	"github.com/google/uuid"
)

var ErrNotAvailable = errors.New("yt-dlp not available")

var _ episode.ContentManager = (*youTubeContentManager)(nil)

const (
	baseCmd    = "yt-dlp --skip-download --write-info-json --no-progress %s"
	loadArgs   = "-o %s.tmp"
	updateArgs = "--no-write-playlist-metafiles --playlist-end 10 --dateafter %s -P \"%s\""
	filterArgs = "--match-filters title~='%s'"

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

func (cm youTubeContentManager) Get(ctx context.Context, url string) (episode.Episode, error) {
	info, err := getInfo(ctx, url)
	if err != nil {
		log.Printf("[ERROR] failed to get info: %v", err)
		return episode.Episode{}, err
	}

	return infoToEpisode(info)
}

func (cm youTubeContentManager) CheckUpdate(ctx context.Context, url string, date time.Time, filter string) (eps []episode.Episode, err error) {
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

	resps := make(chan episode.Episode)
	wg := sync.WaitGroup{}
	wg.Add(len(dlFiles))

	for _, fp := range dlFiles {
		go func(path string) {
			defer wg.Done()

			// Get metadata
			info, err := parseInfo(path)
			if err != nil {
				return
			}

			ep, err := infoToEpisode(info)
			if err != nil {
				log.Printf("[ERROR] failed to decode episode: %v", err)
				return
			}

			resps <- ep
		}(fp)
	}

	go func() {
		wg.Wait()
		close(resps)
	}()

	for r := range resps {
		eps = append(eps, r)
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

func infoToEpisode(info content.FileInfo) (episode.Episode, error) {
	var format content.Format
	for _, f := range info.Formats {
		if f.ID == "140" {
			format = f
			break
		}
	}

	pubDate, err := time.Parse("20060102", info.Date)
	if err != nil {
		return episode.Episode{}, err
	}

	return episode.NewEpisode(
		info.ID,
		"audio/"+format.Extension, // TODO: Fix audio type
		format.URL,
		info.Link,
		info.ImageURL,
		info.Title,
		info.Description, // TODO: "<![CDATA[" + dl.Info.Description + "]]>",
		info.Author,
		pubDate.Format("Mon, 2 Jan 2006 15:04:05 GMT"), // TODO: Save as converter.DateFormat, parse only on RSS build
		format.Length,
		info.Duration,
	), nil
}
