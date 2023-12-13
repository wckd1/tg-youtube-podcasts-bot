package download_file_manager

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/file_manager"
)

type localFile struct {
	path string
	info file_manager.FileInfo
}

// Downloader defines interface to download file to local fs
type Downloader interface {
	Download(ctx context.Context, url string) (file localFile, err error)
	DownloadUpdates(ctx context.Context, url string, date time.Time, filter string) (files []localFile, err error)
}

// Uloader defines interface to upload file from local fs to storage
type Uploader interface {
	Upload(ctx context.Context, file localFile) (url string, err error)
}

type DownloadFileManager struct {
	Downloader Downloader
	Uploader   Uploader
}

func (fm DownloadFileManager) Get(ctx context.Context, url string) (download file_manager.Download, err error) {
	download = file_manager.Download{}

	file, err := fm.Downloader.Download(ctx, url)
	if err != nil {
		return
	}

	download.Info = file.info

	uploadURL, err := fm.Uploader.Upload(ctx, file)
	if err != nil {
		return
	}

	download.URL = uploadURL
	return
}

func (fm DownloadFileManager) CheckUpdate(ctx context.Context, url string, date time.Time, filter string) (downloads []file_manager.Download, err error) {
	files, err := fm.Downloader.DownloadUpdates(ctx, url, date, filter)
	if err != nil {
		return
	}

	for _, f := range files {
		dl := file_manager.Download{Info: f.info}

		uploadURL, err := fm.Uploader.Upload(ctx, f)
		if err != nil {
			continue
		}
		dl.URL = uploadURL

		downloads = append(downloads, dl)
	}

	return
}

func parseInfo(path string) (info file_manager.FileInfo, err error) {
	jsonInfo, err := os.Open(path)
	if err != nil {
		log.Printf("[ERROR] failed to open info json: %v", err)
		return
	}
	defer jsonInfo.Close()

	byteValue, _ := io.ReadAll(jsonInfo)

	json.Unmarshal(byteValue, &info)

	// Remove local info json
	err = os.Remove(path)
	if err != nil {
		log.Printf("[ERROR] failed to delete info json, %v", err)
	}

	return
}
