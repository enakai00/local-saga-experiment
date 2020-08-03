package orderService

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/enakai00/local-saga-experiment/eventQueue"
	"github.com/enakai00/local-saga-experiment/events"
)

var orderDatabase = make(map[string]events.OrderStatus) // in-memory DB

func StartHandlers(orderEventQueue, clientEventQueue, kitchenEventQueue eventQueue.Queue) {
	clientEventSubscription := clientEventQueue.Subscribe("OrderService:ClientEvent")
	kitchenEventSubscription := kitchenEventQueue.Subscribe("OrderService:KitchenEvent")
	go clientEventHandler(orderEventQueue, clientEventSubscription)
	go kitchenEventHandler(orderEventQueue, kitchenEventSubscription)
}

func createOrder(orderRequest *events.OrderRequest) events.OrderStatus {
	id, _ := uuid.NewRandom()
	orderID := id.String()
	orderStatus := events.OrderStatus{
		OrderID:      orderID,
		Status:       "pending",
		ConsumerID:   orderRequest.ConsumerID,
		RestaurantID: orderRequest.RestaurantID,
		FoodID:       orderRequest.FoodID,
	}
	orderDatabase[orderID] = orderStatus // store order in DB
	return orderStatus
}

func updateOrder(orderID string, status string) events.OrderStatus {
	orderStatus := orderDatabase[orderID]
	orderStatus.Status = "approved"
	return orderStatus
}

func clientEventHandler(orderEventQueue eventQueue.Queue,
	clientEventSubscription eventQueue.Subscription) {
	for {
		event := clientEventSubscription.Receive()
		switch event.Type {
		case "OrderRequest":
			orderRequest := new(events.OrderRequest)
			json.Unmarshal(event.Body, orderRequest)
			orderStatus := createOrder(orderRequest)
			jsonBytes, _ := json.Marshal(orderStatus)
			eventQueue.Send("OrderStatus", jsonBytes, orderEventQueue)
		}
	}
}

func kitchenEventHandler(orderEventQueue eventQueue.Queue,
	kitchenEventSubscription eventQueue.Subscription) {
	for {
		event := kitchenEventSubscription.Receive()
		switch event.Type {
		case "TicketStatus":
			ticketStatus := new(events.TicketStatus)
			json.Unmarshal(event.Body, ticketStatus)
			orderID := ticketStatus.OrderID
			status := ticketStatus.Status
			if status == "approved" {
				orderStatus := updateOrder(orderID, status)
				jsonBytes, _ := json.Marshal(orderStatus)
				eventQueue.Send("OrderStatus", jsonBytes, orderEventQueue)
			}
		}
	}
}
