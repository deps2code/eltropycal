package request

import (
	"errors"
	"strings"
)

type OrderRequest struct {
	RestaurantID    string      `json:"restaurant_id"`
	Items           []OrderItem `json:"items"`
	DeliveryAddress Address     `json:"address"`
}

type OrderItem struct {
	FoodID   string `json:"id"`
	Quantity int    `json:"quantity"`
}

type Address struct {
	AddressText string  `json:"address_text"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

func (b *OrderRequest) Validate() error {
	b.RestaurantID = strings.TrimSpace(b.RestaurantID)
	if b.RestaurantID == "" {
		return errors.New("restaurant id is required to order")
	}
	if len(b.Items) < 1 {
		return errors.New("min 1 item is required to order")
	}
	if b.DeliveryAddress.Lat == 0.0 || b.DeliveryAddress.Lng == 0.0 {
		return errors.New("delivery lat and lng are required")
	}
	return nil
}
