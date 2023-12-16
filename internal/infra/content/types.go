package content

type FileInfo struct {
	ID          string   `json:"id"`
	Formats     []Format `json:"formats"`
	Link        string   `json:"webpage_url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ImageURL    string   `json:"thumbnail"`
	Author      string   `json:"uploader"`
	Duration    int      `json:"duration"`
	Date        string   `json:"upload_date"`
}

type Format struct {
	ID        string `json:"format_id"`
	Length    int    `json:"filesize"`
	Extension string `json:"ext"`
	URL       string `json:"url"`
}
