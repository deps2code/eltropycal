package utils

import (
	"os"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//JWTInputData structure
type JWTInputData struct {
	UserID string
	Role   int
}

type UserJWTContext struct {
	UserID       string
	Role         int
	RestaurantID string
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

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
