package entities

import "errors"

type OrderItem struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Comportamiento de OrderItem
func (i *OrderItem) Validate() error {
	if i.ProductID <= 0 {
		return errors.New("invalid product ID")
	}
	if i.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if i.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return nil
}

func (i *OrderItem) Subtotal() (float64, error) {
	if err := i.Validate(); err != nil {
		return 0, err
	}
	return float64(i.Quantity) * i.Price, nil
}
