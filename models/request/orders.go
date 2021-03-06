package request

import (
	"errors"
	"strings"

	"github.com/eltropycal/utils"
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

type OrderStatusUpdateRequest struct {
	OrderID string `json:"order_id"`
	Status  int    `json:"status"`
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

func (b *OrderStatusUpdateRequest) Validate(user interface{}) error {
	b.OrderID = strings.TrimSpace(b.OrderID)
	if b.OrderID == "" {
		return errors.New("order id is a required field")
	}
	// checking that only user with role 2 i.e the customer can cancel the order
	if b.Status < -1 || b.Status == 0 || b.Status > 3 {
		return errors.New("invalid status")
	}
	if b.Status == -1 && user.(utils.UserJWTContext).Role != 2 {
		return errors.New("you are not allowed to cancel the order")
	}
	if b.Status != -1 && user.(utils.UserJWTContext).Role != 3 {
		return errors.New("you are not allowed to update the order status")
	}
	return nil
}
