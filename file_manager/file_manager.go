package file_manager

import (
	"context"
	"encoding/json"
	"io"
	"os"
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
	DownloadUpdates(ctx context.Context, url string, date string)  (files []localFile, err error)
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

func (fm FileManager) CheckUpdate(ctx context.Context, url string, date string) (downloads []Download, err error) {
	files, err := fm.Downloader.DownloadUpdates(ctx, url, date)
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


// TODO: Refactor
// func (fm FileManager) GetPending(ctx context.Context, src string) ([]Download, error) {
//     files, err := os.ReadDir(src)
// 	if err != nil {
// 		log.Printf("[ERROR] failed to open source directory: %v", err)
// 		return nil, err
// 	}

// 	var downloads []Download
//     for _, file := range files {
// 		// TODO: Run in gorutines with channel to handle finish
// 		infoPath := filepath.Join(src, file.Name())
// 		info, err := parseInfo(infoPath)
// 		if err != nil {
// 			log.Printf("[ERROR] failed to open info json: %v", err)
// 			return nil, err
// 		}

// 		fpath, err := fm.Downloader.Download(ctx, info.Link)
// 		if err != nil {
// 			log.Printf("[ERROR] failed to download file: %v", err)
// 			return nil, err
// 		}

// 		lfile := localFile {
// 			path: fpath,
// 			info: info,
// 		}
// 		uploadURL, err := fm.Uploader.Upload(ctx, lfile)
// 		if err != nil {
// 			log.Printf("[ERROR] failed to upload file: %v", err)
// 			return nil, err
// 		}
	
// 		downloads = append(downloads, Download{
// 			URL: uploadURL,
// 			Info: info,
// 		})
//     }

// 	return downloads, nil
// }

func parseInfo(path string) (info fileInfo, err error) {
	jsonInfo, err := os.Open(path)
	if err != nil {
		return
	}
	defer jsonInfo.Close()

	byteValue, _ := io.ReadAll(jsonInfo)
	
	json.Unmarshal(byteValue, &info)
	return
}