package telegram

// Message is primary record to pass data from/to bots
type Message struct {
	ID        int
	ChatID    int64
	Command   string `json:",omitempty"`
	Arguments string `json:",omitempty"`
}

// Response describes bot's answer on particular message
type Response struct {
	Text string
	Send bool
}
