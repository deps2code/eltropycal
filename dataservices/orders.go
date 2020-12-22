package dataservices

import (
	"github.com/eltropycal/constants"
	"github.com/eltropycal/models/dbmodels"
)

func (pc *PostgresClient) CreateOrder(order dbmodels.Order) error {
	query := `INSERT INTO orders (id, user_id, restaurant_id, status, delivery_address_text, delivery_address_lat, delivery_address_lng, items) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := pc.DB.Exec(query, order.ID, order.UserID, order.RestaurantID, order.Status, order.DeliveryAddressText, order.DeliveryAddressLat, order.DeliveryAddressLng, order.Items)
	if err != nil {
		return err
	}
	return nil
}

func (pc *PostgresClient) GetOrdersOfRestaurant(id string) ([]dbmodels.Order, error) {
	query := `SELECT o.id, o.user_id, o.status, o.delivery_address_text, o.delivery_address_lat, o.delivery_address_lng, o.created_at, o.updated_at, o.items, u.username from orders o join users u on o.user_id = u.id where o.restaurant_id=$1`
	rows, err := pc.DB.Query(query, id)
	orders := make([]dbmodels.Order, 0)
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		order := dbmodels.Order{}
		err = rows.Scan(&order.ID, &order.UserID, &order.Status, &order.DeliveryAddressText, &order.DeliveryAddressLat, &order.DeliveryAddressLng, &order.CreatedAt, &order.UpdatedAt, &order.Items, &order.UserName)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, err
}

func (pc *PostgresClient) GetOrderDetailsByID(id string) (dbmodels.Order, error) {
	query := `SELECT id, user_id, restaurant_id, driver_id, status, delivery_address_text, delivery_address_lat, delivery_address_lng, created_at, updated_at, items from orders where id=$1`
	rows, err := pc.DB.Query(query, id)
	order := dbmodels.Order{}
	if err != nil {
		return order, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.DriverID, &order.Status, &order.DeliveryAddressText, &order.DeliveryAddressLat, &order.DeliveryAddressLng, &order.CreatedAt, &order.UpdatedAt, &order.Items)
		if err != nil {
			return order, err
		}
		return order, nil
	}
	return order, nil
}

func (pc *PostgresClient) UpdateOrderStatus(status int, orderID string, driverID string) error {
	tx, _ := pc.DB.Begin()
	query := `UPDATE orders SET status=$2, updated_at=now() where id=$1`
	_, err := tx.Exec(query, orderID, status)
	if err != nil {
		tx.Rollback()
		return err
	}
	if status == constants.CANCELLED || status == constants.DELIVERED {
		query := `UPDATE available_drivers SET available=$2 where driver_id=$1`
		_, err := tx.Exec(query, driverID, true)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (pc *PostgresClient) AssignDriver(driverID string, orderID string) error {
	tx, _ := pc.DB.Begin()
	query := `UPDATE orders SET driver_id=$2, updated_at=now() where id=$1`
	_, err := tx.Exec(query, orderID, driverID)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = `INSERT INTO available_drivers (driver_id, available) VALUES ($1, $2) ON CONFLICT (driver_id) DO UPDATE SET available=excluded.available`
	_, err = tx.Exec(query, driverID, false)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (pc *PostgresClient) GetActiveOrderForDriver(driverID string) (dbmodels.ActiveOrder, error) {
	query := `SELECT o.id, o.user_id, o.restaurant_id, o.status, o.delivery_address_text, o.delivery_address_lat, o.delivery_address_lng, o.created_at, o.updated_at, o.items,
	u.username, r.address, r.lat, r.lng from orders o join users u on o.user_id=u.id join restaurants r on r.id=o.restaurant_id where o.driver_id=$1 and o.status < 3 and o.status > 0 order by o.created_at desc limit 1`
	rows, err := pc.DB.Query(query, driverID)
	order := dbmodels.ActiveOrder{}
	if err != nil {
		return order, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.Status, &order.DeliveryAddressText, &order.DeliveryAddressLat, &order.DeliveryAddressLng,
			&order.CreatedAt, &order.UpdatedAt, &order.Items, &order.UserName, &order.RestaurantAddressText, &order.RestaurantAddressLat, &order.RestaurantAddressLng)
		if err != nil {
			return order, err
		}
		return order, nil
	}
	return order, nil
}
