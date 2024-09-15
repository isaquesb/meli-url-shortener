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

type Subscriber struct {
	Name       string
	ParseEvent func() (Event, error)
	Handler    func(*Message) error
}
