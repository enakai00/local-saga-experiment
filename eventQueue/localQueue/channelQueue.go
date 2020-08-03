package localQueue

import (
	"log"

	"github.com/enakai00/local-saga-experiment/events"
	"github.com/enakai00/local-saga-experiment/eventQueue"
)

type subscriptionStruct struct { // instance of Subscription
	name string // subscription name for Cloud Pub/Sub
	ch   chan events.Event
}

func (s *subscriptionStruct) Receive() events.Event {
	event := <-s.ch
	log.Printf("Receive event by %s: %s", s.name, event.Body)
	return event
}

type localQueue struct { // instance of Queue
	name          string
	subscriptions []subscriptionStruct // chan events.Event
}

func Create(name string) eventQueue.Queue {
	q := localQueue{
		name:          name,
		subscriptions: []subscriptionStruct{}, // chan events.Event{},
	}
	return &q
}

func (q *localQueue) Subscribe(name string) eventQueue.Subscription {
	c := make(chan events.Event, 10)
	subscription := subscriptionStruct{
		name: name,
		ch:   c,
	}
	q.subscriptions = append(q.subscriptions, subscription)
	return &subscription
}

func (q *localQueue) Send(event events.Event) {
	log.Printf("Send event to %s: %s", q.name, event.Body)
	for _, s := range q.subscriptions {
		s.ch <- event
	}
}
