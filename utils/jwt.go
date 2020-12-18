package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//JWTInputData structure
type JWTInputData struct {
	UserID string
	Role   int
}

//EncodeJwt : function to encode data and return a JWT token
func EncodeJwt(data JWTInputData) (tokenString string) {
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": data.UserID,
		"role":   data.Role,
		"exp":    time.Now().Add(time.Hour * time.Duration(24)).Unix(), //should be a unix timestamp of a date greater then the current time
	})
	secretKey := []byte(os.Getenv("ELTROPYCAL_JWT_SECRET"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return tokenString
	}
	return
}
