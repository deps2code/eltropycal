package dbmodels

import "database/sql"

type Order struct {
	ID                  string         `json:"id"`
	UserID              string         `json:"user_id"`
	UserName            string         `json:"user_name"`
	RestaurantID        string         `json:"restaurant_id"`
	Status              int            `json:"status"`
	DeliveryAddressText string         `json:"delivery_address_text"`
	DeliveryAddressLat  float64        `json:"delivery_address_lat"`
	DeliveryAddressLng  float64        `json:"delivery_address_lng"`
	DriverID            string         `json:"driver_id"`
	Items               string         `json:"items"`
	CreatedAt           sql.NullString `json:"created_at"`
	UpdatedAt           sql.NullString `json:"updated_at"`
}

type ActiveOrder struct {
	ID                    string         `json:"id"`
	UserID                string         `json:"user_id"`
	UserName              string         `json:"user_name"`
	RestaurantID          string         `json:"restaurant_id"`
	Status                int            `json:"status"`
	DeliveryAddressText   string         `json:"delivery_address_text"`
	DeliveryAddressLat    float64        `json:"delivery_address_lat"`
	DeliveryAddressLng    float64        `json:"delivery_address_lng"`
	DriverID              string         `json:"driver_id"`
	Items                 string         `json:"items"`
	CreatedAt             sql.NullString `json:"created_at"`
	UpdatedAt             sql.NullString `json:"updated_at"`
	RestaurantAddressText string         `json:"restaurant_address_text"`
	RestaurantAddressLat  float64        `json:"restaurant_address_lat"`
	RestaurantAddressLng  float64        `json:"restaurant_address_lng"`
}
