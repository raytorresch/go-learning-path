package entities

import (
	"errors"
	"time"
	"user-management/internal/domain/valueobjects"

	"github.com/google/uuid"
)

type Order struct {
	ID          int                      `json:"id"`
	UserID      int                      `json:"user_id"`
	Items       []OrderItem              `json:"items"`
	Total       float64                  `json:"total"`
	Status      valueobjects.OrderStatus `json:"status"`
	CreatedAt   time.Time                `json:"created_at"`
	CompletedAt time.Time                `json:"completed_at,omitempty"`
}

func NewOrder(userID uuid.UUID, items []OrderItem) (*Order, error) {
	order := &Order{
		UserID:    int(userID.ID()),
		Items:     items,
		Status:    valueobjects.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := order.CalculateTotal(); err != nil {
		return nil, err
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	return order, nil
}

func (o *Order) CalculateTotal() error {
	total := 0.0
	for _, item := range o.Items {
		subtotal, err := item.Subtotal()
		if err != nil {
			return err
		}
		total += subtotal
	}
	o.Total = total
	return nil
}

func (o *Order) AddItem(item OrderItem) error {
	if err := item.Validate(); err != nil {
		return err
	}

	// Buscar si ya existe el producto
	for i, existing := range o.Items {
		if existing.ProductID == item.ProductID {
			// Actualizar cantidad
			o.Items[i].Quantity += item.Quantity
			return o.CalculateTotal()
		}
	}

	// Nuevo item
	o.Items = append(o.Items, item)
	return o.CalculateTotal()
}

func (o *Order) Validate() error {
	if o.UserID <= 0 {
		return errors.New("invalid user ID")
	}
	if len(o.Items) == 0 {
		return errors.New("order must have at least one item")
	}
	if o.Total < 0 {
		return errors.New("total cannot be negative")
	}
	if !o.CreatedAt.IsZero() && !o.CompletedAt.IsZero() &&
		o.CompletedAt.Before(o.CreatedAt) {
		return errors.New("completed date cannot be before creation date")
	}
	return nil
}

func (o *Order) Complete() error {
	if o.Status == valueobjects.StatusCompleted {
		return errors.New("order already completed")
	}
	o.Status = valueobjects.StatusCompleted
	o.CompletedAt = time.Now()
	return nil
}
