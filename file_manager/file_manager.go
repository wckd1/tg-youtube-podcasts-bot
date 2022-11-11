package file_manager

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"mvdan.cc/xurls/v2"
	"os"
	"strings"
	"time"
)

const (
	maxDescriptionLenght = 500
)

type localFile struct {
	path string
	info fileInfo
}

type fileInfo struct {
	Link        string `json:"webpage_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"thumbnail"`
	Author      string `json:"uploader"`
	Length      int    `json:"filesize"`
	Duration    int    `json:"duration"`
	Date        string `json:"upload_date"`
}

type Download struct {
	URL  string
	Info fileInfo
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

type FileManager struct {
	Downloader Downloader
	Uploader   Uploader
}

func (fm FileManager) Get(ctx context.Context, url string) (download Download, err error) {
	download = Download{}

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

func (fm FileManager) CheckUpdate(ctx context.Context, url string, date time.Time, filter string) (downloads []Download, err error) {
	files, err := fm.Downloader.DownloadUpdates(ctx, url, date, filter)
	if err != nil {
		return
	}

	for _, f := range files {
		dl := Download{Info: f.info}

		uploadURL, err := fm.Uploader.Upload(ctx, f)
		if err != nil {
			continue
		}
		dl.URL = uploadURL

		downloads = append(downloads, dl)
	}

	return
}

func parseInfo(path string) (info fileInfo, err error) {
	jsonInfo, err := os.Open(path)
	if err != nil {
		log.Printf("[ERROR] failed to open info json: %v", err)
		return
	}
	defer jsonInfo.Close()

	byteValue, _ := io.ReadAll(jsonInfo)

	json.Unmarshal(byteValue, &info)

	// Sanitize description
	desc := info.Description
	furls := xurls.Relaxed().FindAllString(desc, -1)
	for _, u := range furls {
		desc = strings.ReplaceAll(desc, u, "")
	}
	if len(desc) > maxDescriptionLenght {
		info.Description = desc[:maxDescriptionLenght - 3] + "..."
	} else {
		info.Description = desc
	}

	// Remove local info json
	err = os.Remove(path)
	if err != nil {
		log.Printf("[ERROR] failed to delete info json, %v", err)
	}

	return
}
