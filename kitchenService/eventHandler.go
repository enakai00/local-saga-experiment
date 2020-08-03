package kitchenService

import (
	"encoding/json"

	"github.com/enakai00/local-saga-experiment/eventQueue"
	"github.com/enakai00/local-saga-experiment/events"
)

var ticketDatabase = make(map[string]events.TicketStatus) // in-memory DB

func StartHandlers(kitchenEventQueue, orderEventQueue, consumerEventQueue eventQueue.Queue) {
	orderEventSubscription := orderEventQueue.Subscribe("KitchenService:OrderEvent")
	consumerEventSubscription := consumerEventQueue.Subscribe("KitchenService:ConsumerEvent")
	go orderEventHandler(kitchenEventQueue, orderEventSubscription)
	go consumerEventHandler(kitchenEventQueue, consumerEventSubscription)
}

func createTicket(orderStatus *events.OrderStatus) events.TicketStatus {
	orderID := orderStatus.OrderID
	foodType := "vegan" // decided from orderStatus.FoodID
	ticketStatus := events.TicketStatus{
		OrderID:      orderID,
		ConsumerID:   orderStatus.ConsumerID,
		Status:       "pending",
		RestaurantID: orderStatus.RestaurantID,
		FoodID:       orderStatus.FoodID,
		FoodType:     foodType,
	}
	ticketDatabase[orderID] = ticketStatus // store order in DB
	return ticketStatus
}

func updateTicket(orderID string) events.TicketStatus {
	ticketStatus := ticketDatabase[orderID]
	ticketStatus.Status = "approved"
	return ticketStatus
}

func orderEventHandler(kitchenEventQueue eventQueue.Queue,
	orderEventSubscription eventQueue.Subscription) {
	for {
		event := orderEventSubscription.Receive()
		switch event.Type {
		case "OrderStatus":
			orderStatus := new(events.OrderStatus)
			json.Unmarshal(event.Body, orderStatus)
			if orderStatus.Status == "pending" {
				ticketStatus := createTicket(orderStatus)
				jsonBytes, _ := json.Marshal(ticketStatus)
				eventQueue.Send("TicketStatus", jsonBytes, kitchenEventQueue)
			}
		}
	}
}

func consumerEventHandler(kitchenEventQueue eventQueue.Queue,
	consumerEventSubscription eventQueue.Subscription) {
	for {
		event := consumerEventSubscription.Receive()
		switch event.Type {
		case "ConsumerVerification":
			consumerVerification := new(events.ConsumerVerification)
			json.Unmarshal(event.Body, consumerVerification)
			if consumerVerification.Status == "verified" {
				orderID := consumerVerification.OrderID
				ticketStatus := updateTicket(orderID)
				jsonBytes, _ := json.Marshal(ticketStatus)
				eventQueue.Send("TicketStatus", jsonBytes, kitchenEventQueue)
			}
		}
	}
}
