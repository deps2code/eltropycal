package request

import (
	"errors"
	"strings"
)

//LoginRequest -
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Validate for validating Request for Login
func (c *LoginRequest) Validate() error {
	if strings.TrimSpace(c.Username) == "" {
		return errors.New("the username field is required")
	}
	if strings.TrimSpace(c.Password) == "" {
		return errors.New("the password field is required")
	}
	return nil
}
