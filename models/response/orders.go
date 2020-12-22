package response

type OrderResponse struct {
	OrderID           string      `json:"order_id"`
	Items             []FoodItem  `json:"items"`
	DeliveryAddress   Address     `json:"address"`
	User              User        `json:"user"`
	Status            int         `json:"status"`
	RestaurantAddress interface{} `json:"restaurant_address,omitempty"`
}

type OrderItem struct {
	FoodID   string `json:"id"`
	Quantity int    `json:"quantity"`
}

type FoodItem struct {
	FoodName    string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

type Address struct {
	AddressText string  `json:"address_text"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type OrderEstimationResponse struct {
	EstimatedPrice float64 `json:"price_in_rupees"`
	EstimatedTime  float64 `json:"time_in_minutes"`
}
