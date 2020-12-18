package dbmodels

type Restaurant struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	AddressText string  `json:"address_text"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

type Food struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	RestaurantID    string  `json:"restaurant_id, omitempty"`
	UnitRateRupees  float64 `json:"unit_rate_rupees"`
	UnitTimeMinutes float64 `json:"unit_time_mins"`
}
