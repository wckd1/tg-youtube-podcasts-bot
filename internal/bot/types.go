package bot

// Response describes bot's answer on particular message
type Response struct {
	Text string
	Send bool
}

// Message is primary record to pass data from/to bots
type Message struct {
	ID        int
	Command   string `json:",omitempty"`
	Arguments string `json:",omitempty"`
}
