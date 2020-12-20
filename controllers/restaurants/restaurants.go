package restaurants

import (
	"net/http"

	"github.com/eltropycal/controllers"
	"github.com/eltropycal/models/response"
	"github.com/eltropycal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
)

//GetAllRestaurants List all resturants for user to choose
func GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "userContext")
	restaurantList, err := controllers.DataService.GetAllRestaurantList()
	if err != nil {
		log.WithFields(log.Fields{"api": "GetAllRestaurants", "error": "db_error", "user_id": user.(utils.JWTInputData).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	controllers.RespondWithSuccessGeneric(w, http.StatusOK, response.Response{
		Success: true,
		Data:    restaurantList,
	})
	return
}

//GetAllRestaurants List all resturants for user to choose
func GetRestaurantMenu(w http.ResponseWriter, r *http.Request) {
	// user := context.Get(r, "userContext")
	restaurantID, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		log.WithFields(log.Fields{"api": "GetRestaurantMenu", "error": "invalid_request_body", "user_id": /*user.(utils.JWTInputData).UserID*/ ""}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid restaurant id")
		return
	}
	restaurantItemsList, err := controllers.DataService.GetRestaurantMenuItems(restaurantID.String())
	if err != nil {
		log.WithFields(log.Fields{"api": "GetRestaurantMenu", "error": "db_error", "user_id": /*user.(utils.JWTInputData).UserID*/ ""}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	controllers.RespondWithSuccessGeneric(w, http.StatusOK, response.Response{
		Success: true,
		Data:    restaurantItemsList,
	})
	return
}
