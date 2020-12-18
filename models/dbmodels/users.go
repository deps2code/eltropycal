package dbmodels

type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"password_hash"`
	Role           int    `json:"role"`
}
