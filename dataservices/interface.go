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
	GetAllRestuarantList() ([]dbmodels.Restaurant, error)
	GetRestaurantMenuItems(string) ([]dbmodels.Food, error)
}
