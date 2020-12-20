package dataservices

import (
	"github.com/eltropycal/models/dbmodels"
)

func (pc *PostgresClient) GetAllRestaurantList() ([]dbmodels.Restaurant, error) {
	query := `SELECT id, name, address, lat, lng from restaurants`
	rows, err := pc.DB.Query(query)
	var restaurantList []dbmodels.Restaurant
	if err != nil {
		return restaurantList, err
	}
	defer rows.Close()
	for rows.Next() {
		restaurant := dbmodels.Restaurant{}
		err = rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.AddressText, &restaurant.Lat, &restaurant.Lng)
		if err != nil {
			return restaurantList, err
		}
		restaurantList = append(restaurantList, restaurant)
	}
	return restaurantList, nil
}

func (pc *PostgresClient) GetRestaurantMenuItems(id string) ([]dbmodels.Food, error) {
	query := `SELECT id, name, description, restaurant_id, unit_rate_rupees, unit_time_mins from foods where restaurant_id=$1`
	rows, err := pc.DB.Query(query, id)
	foodItems := make([]dbmodels.Food, 0)
	if err != nil {
		return foodItems, err
	}
	defer rows.Close()
	for rows.Next() {
		food := dbmodels.Food{}
		err = rows.Scan(&food.ID, &food.Name, &food.Description, &food.RestaurantID, &food.UnitRateRupees, &food.UnitTimeMinutes)
		if err != nil {
			return foodItems, err
		}
		foodItems = append(foodItems, food)
	}
	return foodItems, nil
}

func (pc *PostgresClient) GetRestaurantDetailsByID(id string) (dbmodels.Restaurant, error) {
	query := `SELECT id, name, address, lat, lng from restaurants where id=$1`
	rows, err := pc.DB.Query(query, id)
	var restaurantDetails dbmodels.Restaurant
	if err != nil {
		return restaurantDetails, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&restaurantDetails.ID, &restaurantDetails.Name, &restaurantDetails.AddressText, &restaurantDetails.Lat, &restaurantDetails.Lng)
		if err != nil {
			return restaurantDetails, err
		}
		return restaurantDetails, nil
	}
	return restaurantDetails, nil
}
