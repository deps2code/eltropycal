package dataservices

import "database/sql"

type PostgresClient struct {
	DB *sql.DB
}

// IPostgresClient - DB interface
type IPostgresClient interface {
	Connect()
}
