package dataservices

import (
	"github.com/eltropycal/models/dbmodels"
)

func (pc *PostgresClient) GetUserDetailsByUsername(username string, role int) (dbmodels.User, error) {
	query := `SELECT id, hash_password from restaurants where username=$1 and role=$2`
	rows, err := pc.DB.Query(query, username, role)
	userDetails := dbmodels.User{}
	if err != nil {
		return userDetails, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&userDetails.ID, &userDetails.HashedPassword)
		if err != nil {
			return userDetails, err
		}
	}
	return userDetails, nil
}
