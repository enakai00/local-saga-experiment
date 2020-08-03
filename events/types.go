package events

type Event struct {
	ID   string
	Type string // OrderRequest, OrderStatus, ConsumerVerification, etc...
	Body []byte // JSON string
}

// Event.Body
type OrderRequest struct {
	ConsumerID   string `json:"consumerID"`
	RestaurantID string `json:"restaurantID"`
	FoodID       string `json:"foodID"`
}

type OrderStatus struct {
	OrderID      string `json:"orderID"`
	Status       string `json:"orderStatus"` // pending, approved, rejected
	ConsumerID   string `json:"consumerID"`
	RestaurantID string `json:"restaurantID"`
	FoodID       string `json:"foodID"`
}

type ConsumerVerification struct {
	OrderID    string `json:"orderID"`
	ConsumerID string `json:consumerID"`
	Status     string `json:"status"` // verified, failed
}

type TicketStatus struct {
	OrderID      string `json:"orderID"`
	ConsumerID   string `json:"consumerID"`
	Status       string `json:"status"` // peding, approved, rejected
	RestaurantID string `json:"restaurantID"`
	FoodID       string `json:"FoodID"`
	FoodType     string `json:"FoodType"` // vegan, beef, etc...
}
