package events

type Message struct {
	Uuid    string
	Payload []byte
}

type Handler struct {
	Name   string
	Handle func(*Message) error
}
