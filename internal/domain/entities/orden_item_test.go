package entities

import (
	"testing"
)

func TestOrderItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    OrderItem
		wantErr bool
	}{
		{
			name: "valid item",
			item: OrderItem{
				ProductID: 1,
				Name:      "Product A",
				Quantity:  2,
				Price:     10.0,
			},
			wantErr: false,
		},
		{
			name: "invalid product ID",
			item: OrderItem{
				ProductID: 0,
				Name:      "Product B",
				Quantity:  1,
				Price:     5.0,
			},
			wantErr: true,
		},
		{
			name: "negative quantity",
			item: OrderItem{
				ProductID: 2,
				Name:      "Product C",
				Quantity:  -1,
				Price:     15.0,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			item: OrderItem{
				ProductID: 3,
				Name:      "Product D",
				Quantity:  1,
				Price:     -20.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderItem.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderItem_Subtotal(t *testing.T) {
	tests := []struct {
		name      string
		item      OrderItem
		wantTotal float64
		wantErr   bool
	}{
		{
			name: "valid item",
			item: OrderItem{
				ProductID: 1,
				Name:      "Product A",
				Quantity:  2,
				Price:     10.0,
			},
			wantTotal: 20.0,
			wantErr:   false,
		},
		{
			name: "invalid item",
			item: OrderItem{
				ProductID: 0,
				Name:      "Product B",
				Quantity:  1,
				Price:     5.0,
			},
			wantTotal: 0.0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total, err := tt.item.Subtotal()
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderItem.Subtotal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if total != tt.wantTotal {
				t.Errorf("OrderItem.Subtotal() = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}
