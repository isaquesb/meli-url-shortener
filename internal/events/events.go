package events

type Event interface {
	GetName() string
	GetKey() []byte
}
