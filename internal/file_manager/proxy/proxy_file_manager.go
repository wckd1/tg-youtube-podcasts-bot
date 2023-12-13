package proxy_file_manager

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
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/file_manager"

	"github.com/google/uuid"
)

var _ file_manager.FileManager = (*ProxyFileManager)(nil)

const (
	baseCmd  = "yt-dlp --skip-download --write-info-json --no-progress %s"
	loadArgs = "-o %s.tmp"
	dlPath   = "./storage/downloads/"
	infoExt  = ".info.json"

	audioCmd = "yt-dlp --get-url -f 140 %s"
)

type ProxyFileManager struct{}

func (fm ProxyFileManager) Get(ctx context.Context, url string) (file_manager.Download, error) {
	dl := file_manager.Download{
		URL:  url,
		Info: file_manager.FileInfo{},
	}

	// Get link
	audio, err := getAudioLink(ctx, url)
	if err != nil {
		log.Printf("[ERROR] failed to get audio link: %v", err)
		return file_manager.Download{}, err
	}
	dl.URL = audio

	// Get metadata
	info, err := getInfo(ctx, url)
	if err != nil {
		log.Printf("[ERROR] failed to get info: %v", err)
		return file_manager.Download{}, err
	}
	dl.Info = info

	return dl, nil
}

// TODO: Add
func (fm ProxyFileManager) CheckUpdate(ctx context.Context, url string, date time.Time, filter string) (downloads []file_manager.Download, err error) {
	return
}

func getInfo(ctx context.Context, url string) (file_manager.FileInfo, error) {
	id := uuid.New().String()
	// Prepare yt-dlp command
	fcmd := fmt.Sprintf(baseCmd, url)
	args := fmt.Sprintf(loadArgs, id)
	cmdStr := strings.Join([]string{fcmd, args}, " ")

	// Load metadata
	if err := executeCommand(ctx, cmdStr); err != nil {
		return file_manager.FileInfo{}, fmt.Errorf("failed to execute command: %v", err)
	}

	// Parse image, title and description
	infoPath := filepath.Join(dlPath, id+".tmp"+infoExt)

	jsonInfo, err := os.Open(infoPath)
	if err != nil {
		return file_manager.FileInfo{}, fmt.Errorf("failed to open info json: %v", err)
	}
	defer jsonInfo.Close()

	var info file_manager.FileInfo
	byteValue, _ := io.ReadAll(jsonInfo)
	json.Unmarshal(byteValue, &info)

	// Remove local info json
	err = os.Remove(infoPath)
	if err != nil {
		log.Printf("[ERROR] failed to delete info json, %v", err)
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

func executeCommand(ctx context.Context, cmdStr string) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Dir = dlPath

	log.Printf("[DEBUG] executing command: %s", cmd.String())
	return cmd.Run()
}
