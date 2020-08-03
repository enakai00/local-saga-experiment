package main

import (
	"encoding/json"

	"github.com/enakai00/local-saga-experiment/consumerService"
	"github.com/enakai00/local-saga-experiment/eventQueue"
	"github.com/enakai00/local-saga-experiment/eventQueue/localQueue"
	"github.com/enakai00/local-saga-experiment/events"
	"github.com/enakai00/local-saga-experiment/kitchenService"
	"github.com/enakai00/local-saga-experiment/orderService"
)

func main() {
	// Create event queues
	clientEventQueue := localQueue.Create("ClientEvents") // interface eventQueue.Queue
	orderEventQueue := localQueue.Create("OrderEvents")
	consumerEventQueue := localQueue.Create("ConsumerEvents")
	kitchenEventQueue := localQueue.Create("KitchenEvents")

	// Start handlers
	orderService.StartHandlers(orderEventQueue, clientEventQueue, kitchenEventQueue)
	kitchenService.StartHandlers(kitchenEventQueue, orderEventQueue, consumerEventQueue)
	consumerService.StartHandlers(consumerEventQueue, kitchenEventQueue)

	// Start saga with an OrderRequest event
	orderRequest := events.OrderRequest{
		ConsumerID:   "consumer001",
		RestaurantID: "restaurant001",
		FoodID:       "food001",
	}
	jsonBytes, _ := json.Marshal(orderRequest)
	eventQueue.Send("OrderRequest", jsonBytes, clientEventQueue)

	// Oberve OrderEvent queue
	orderEventSubscription := orderEventQueue.Subscribe("Client:OrderEvent")
	for {
		orderEventSubscription.Receive()
	}
}
