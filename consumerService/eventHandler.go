package consumerService

import (
	"encoding/json"

	"github.com/enakai00/local-saga-experiment/eventQueue"
	"github.com/enakai00/local-saga-experiment/events"
)

func StartHandlers(consumerEventQueue, kitchenEventQueue eventQueue.Queue) {
	kitchenEventSubscription := kitchenEventQueue.Subscribe("ConsumerService:KitchenEvent")
	go kitchenEventHandler(consumerEventQueue, kitchenEventSubscription)
}

func kitchenEventHandler(consumerEventQueue eventQueue.Queue,
	kitchenEventSubscription eventQueue.Subscription) {
	for {
		event := kitchenEventSubscription.Receive()
		switch event.Type {
		case "TicketStatus":
			ticketStatus := new(events.TicketStatus)
			json.Unmarshal(event.Body, ticketStatus)
			if ticketStatus.Status == "pending" {
				status := "verified" // deceided from ticketStatus.consumerID and FoodTYpe
				consumerVerification := events.ConsumerVerification{
					OrderID:    ticketStatus.OrderID,
					ConsumerID: ticketStatus.ConsumerID,
					Status:     status,
				}
				jsonBytes, _ := json.Marshal(consumerVerification)
				eventQueue.Send("ConsumerVerification", jsonBytes, consumerEventQueue)
			}
		default:
			eventQueue.Send("Error", []byte{}, consumerEventQueue)
		}
	}
}
