package dataservices

import (
	"time"

	"github.com/eltropycal/models/dbmodels"
)

func (pc *PostgresClient) GetUserDetailsByUsername(username string, role int) (dbmodels.User, error) {
	query := `SELECT id, hash_password from users where username=$1 and role=$2`
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

func (pc *PostgresClient) GetUserSession(token string, userID string) (dbmodels.UserSession, error) {
	query := `SELECT us.user_id, us.expires_at, ur.restaurant_id from user_session us left join user_restaurant_mapping ur on us.user_id = ur.user_id where us.user_id=$1 and us.token=$2`
	rows, err := pc.DB.Query(query, userID, token)
	userDetails := dbmodels.UserSession{}
	if err != nil {
		return userDetails, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&userDetails.UserID, &userDetails.ExpiresAt, &userDetails.RestaurantID)
		if err != nil {
			return userDetails, err
		}
	}
	return userDetails, nil
}

func (pc *PostgresClient) CreateUserSession(userSession dbmodels.UserSession) error {
	query := `INSERT INTO user_session (user_id, token, expires_at, created_at) VALUES($1, $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET token=excluded.token, expires_at=excluded.expires_at`
	rows, err := pc.DB.Query(query, userSession.UserID, userSession.Token, userSession.ExpiresAt, userSession.CreatedAt)
	defer rows.Close()
	if err != nil {
		return err
	}
	return nil
}

func (pc *PostgresClient) InvalidateUserSession(userID string) error {
	query := `UPDATE user_session SET expires_at=$2 where user_id=$1`
	rows, err := pc.DB.Query(query, userID, time.Now())
	defer rows.Close()
	if err != nil {
		return err
	}
	return nil
}
