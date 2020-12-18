package request

import (
	"errors"
	"strings"
)

//LoginRequest -
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

//Validate for validating Request for Login
func (b *LoginRequest) Validate() error {
	b.Username = strings.TrimSpace(b.Username)
	if b.Username == "" {
		return errors.New("sername field is required")
	}
	b.Password = strings.TrimSpace(b.Password)
	if b.Password == "" {
		return errors.New("password field is required")
	}
	if b.Role == 0 {
		return errors.New("role field is required")
	}
	return nil
}
