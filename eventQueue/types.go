package eventQueue

import (
	"github.com/google/uuid"

	"github.com/enakai00/local-saga-experiment/events"
)

type Subscription interface {
	Receive() events.Event
}

type Queue interface {
	Subscribe(string) Subscription
	Send(events.Event)
}

func Send(eventType string, jsonBytes []byte, queue Queue) {
	id, _ := uuid.NewRandom()
	event := events.Event{
		ID:   id.String(),
		Type: eventType,
		Body: jsonBytes,
	}
	queue.Send(event)
}
