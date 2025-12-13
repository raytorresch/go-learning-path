package entities

import "time"

type Order struct {
	ID          int         `json:"id"`
	UserID      int         `json:"user_id"`
	Items       []OrderItem `json:"items"`
	Total       float64     `json:"total"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	CompletedAt time.Time   `json:"completed_at,omitempty"`
}

type OrderItem struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderStatus string

const (
	StatusReceived OrderStatus = "received"
)
