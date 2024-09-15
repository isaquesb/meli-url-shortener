package urls

import (
	"fmt"
	"github.com/isaquesb/meli-url-shortener/internal/events"
	"time"
)

const Created = "urls.created"
const Deleted = "urls.deleted"
const Visited = "urls.visited"

func MakeEvent(name string) (events.Event, error) {
	if name == Created {
		return &CreateEvent{}, nil
	}

	if name == Deleted {
		return &DeleteEvent{}, nil
	}

	if name == Visited {
		return &VisitEvent{}, nil
	}

	return nil, fmt.Errorf("event not found: %s", name)
}

func EventParserFor(evtName string) func() (events.Event, error) {
	return func() (events.Event, error) { return MakeEvent(evtName) }
}

type CreateEvent struct {
	ShortCode []byte    `json:"short"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}

func (ce *CreateEvent) GetName() string { return Created }
func (ce *CreateEvent) GetKey() []byte  { return ce.ShortCode }

func NewCreateEvent(short []byte, url string) *CreateEvent {
	return &CreateEvent{
		ShortCode: short,
		Url:       url,
		CreatedAt: time.Now(),
	}
}

type DeleteEvent struct {
	ShortCode []byte `json:"short"`
}

func (de *DeleteEvent) GetName() string { return Deleted }
func (de *DeleteEvent) GetKey() []byte  { return de.ShortCode }

func NewDeleteEvent(short []byte) *DeleteEvent {
	return &DeleteEvent{
		ShortCode: short,
	}
}

type VisitEvent struct {
	ShortCode []byte    `json:"short"`
	Date      time.Time `json:"date"`
}

func (ve *VisitEvent) GetName() string { return Visited }
func (ve *VisitEvent) GetKey() []byte  { return ve.ShortCode }

func NewVisitEvent(short []byte) *VisitEvent {
	return &VisitEvent{
		ShortCode: short,
		Date:      time.Now(),
	}
}
