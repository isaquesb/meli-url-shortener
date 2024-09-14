package events

type Envelop struct {
	Name  string `json:"name"`
	Event Event  `json:"event"`
}

type Message struct {
	Uuid  string
	Name  string
	Event Event
}

type Handler struct {
	Id         string
	ParseEvent func() (Event, error)
	Handle     func(*Message) error
}
