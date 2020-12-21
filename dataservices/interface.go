package dataservices

import (
	"database/sql"

	"github.com/eltropycal/models/dbmodels"
)

type PostgresClient struct {
	DB *sql.DB
}

// IPostgresClient - DB interface
type IPostgresClient interface {
	Connect()
	GetUserDetailsByUsername(string, int) (dbmodels.User, error)
	GetAllRestaurantList() ([]dbmodels.Restaurant, error)
	GetRestaurantMenuItems(string) ([]dbmodels.Food, error)
	CreateOrder(dbmodels.Order) error
	GetOrdersOfRestaurant(string) ([]dbmodels.Order, error)
	GetOrderDetailsByID(string) (dbmodels.Order, error)
	GetFoodDetailsByID(string) (dbmodels.Food, error)
	GetRestaurantDetailsByID(string) (dbmodels.Restaurant, error)
	GetUserSession(string, string) (dbmodels.UserSession, error)
	CreateUserSession(dbmodels.UserSession) error
	InvalidateUserSession(string) error
}
