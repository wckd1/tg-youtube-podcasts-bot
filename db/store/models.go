package db

type Subscription struct {
	ID      int64  `json:"id"`
	Channel string `json:"channel"`
	Title   string `json:"title"`
}
