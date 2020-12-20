package dataservices

import (
	"github.com/eltropycal/models/dbmodels"
)

func (pc *PostgresClient) CreateOrder(order dbmodels.Order) error {
	query := `INSERT INTO orders (id, user_id, restaurant_id, status, delivery_address_text, delivery_address_lat, delivery_address_lng, items) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	rows, err := pc.DB.Query(query, order.ID, order.UserID, order.RestaurantID, order.Status, order.DeliveryAddressText, order.DeliveryAddressLat, order.DeliveryAddressLng, order.Items)
	defer rows.Close()
	if err != nil {
		return err
	}
	return nil
}

func (pc *PostgresClient) GetOrdersOfRestaurant(id string) ([]dbmodels.Order, error) {
	query := `SELECT o.id, o.user_id, o.status, o.delivery_address_text, o.delivery_address_lat, o.delivery_address_lng, o.created_at, o.delivered_at, o.items, u.username from orders o join users u on o.user_id = u.id where o.restaurant_id=$1`
	rows, err := pc.DB.Query(query, id)
	orders := make([]dbmodels.Order, 0)
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		order := dbmodels.Order{}
		err = rows.Scan(&order.ID, &order.UserID, &order.Status, &order.DeliveryAddressText, &order.DeliveryAddressLat, &order.DeliveryAddressLng, &order.CreatedAt, &order.DeliveredAt, &order.Items, &order.UserName)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, err
}

func (pc *PostgresClient) GetOrderDetailsByID(id string) (dbmodels.Order, error) {
	query := `SELECT id, user_id, restaurant_id, status, delivery_address_text, delivery_address_lat, delivery_address_lng, created_at, delivered_at, items from orders where id=$1`
	rows, err := pc.DB.Query(query, id)
	order := dbmodels.Order{}
	if err != nil {
		return order, err
	}
	defer rows.Close()
	if rows.Next() {
		order := dbmodels.Order{}
		err = rows.Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.Status, &order.DeliveryAddressText, &order.DeliveryAddressLat, &order.DeliveryAddressLng, &order.CreatedAt, &order.DeliveredAt, &order.Items)
		if err != nil {
			return order, err
		}
		return order, nil
	}
	return order, nil
}
