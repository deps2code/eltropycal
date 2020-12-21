package dbmodels

import "time"

type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"password_hash"`
	Role           int    `json:"role"`
}

type UserSession struct {
	UserID       string    `json:"user_id"`
	Token        string    `json:"token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	RestaurantID string    `json:"restaurant_id"`
}
