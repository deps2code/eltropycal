package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/eltropycal/controllers"
	"github.com/eltropycal/utils"
	"github.com/gorilla/context"
)

func Validate(next http.HandlerFunc, roles []int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					secret := os.Getenv("ELTROPYCAL_JWT_SECRET")
					return []byte(secret), nil
				})
				if err != nil {
					controllers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
					return
				}
				if token.Valid {
					tokenMapClaims := token.Claims.(jwt.MapClaims)
					userSession, err := controllers.DataService.GetUserSession(bearerToken[1], tokenMapClaims["userID"].(string))
					if err != nil || userSession.UserID == "" {
						controllers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
						return
					}
					if userSession.ExpiresAt.Before(time.Now()) {
						controllers.RespondWithError(w, http.StatusUnauthorized, "Session timed out. Please login again.")
						return
					}
					isRoleAllowed, _ := utils.InArray(int(tokenMapClaims["role"].(float64)), roles)
					if !isRoleAllowed {
						controllers.RespondWithError(w, http.StatusConflict, "You do not have correct access to call this api")
						return
					}
					userContext := utils.UserJWTContext{
						UserID:       tokenMapClaims["userID"].(string),
						Role:         int(tokenMapClaims["role"].(float64)),
						RestaurantID: userSession.RestaurantID,
					}
					context.Set(req, "userContext", userContext)
					next(w, req)
				} else {
					controllers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
					return
				}
			}
		} else {
			controllers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
	})
}
