package orders

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/eltropycal/constants"
	"github.com/eltropycal/controllers"
	"github.com/eltropycal/models/dbmodels"
	"github.com/eltropycal/models/request"
	"github.com/eltropycal/models/response"
	"github.com/eltropycal/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	// user := context.Get(r, "userContext")
	var apiReq request.OrderRequest
	err := json.NewDecoder(r.Body).Decode(&apiReq)
	if err != nil {
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "invalid_api_request"}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err = apiReq.Validate(); err != nil {
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "invalid_request_body"}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	foodItems, err := json.Marshal(apiReq.Items)
	if err != nil {
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "invalid_request_body"}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "food items are not correct")
		return
	}

	order := dbmodels.Order{
		ID:                  uuid.New().String(),
		UserID:              "60bc850b-5048-4d55-8e59-552a85ee1a47", /*user.(utils.JWTInputData).UserID*/
		RestaurantID:        apiReq.RestaurantID,
		Status:              constants.PREPARING,
		DeliveryAddressText: apiReq.DeliveryAddress.AddressText,
		DeliveryAddressLat:  apiReq.DeliveryAddress.Lat,
		DeliveryAddressLng:  apiReq.DeliveryAddress.Lng,
		DriverID:            "",
		Items:               string(foodItems),
	}
	err = controllers.DataService.CreateOrder(order)
	if err != nil {
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "db_error"}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong")
		return
	}

	//start go routine for assigniing driver
	controllers.RespondWithSuccess(w, http.StatusOK, "We have received your order, it will be delivered to you soon")
	return

}

//GetAllRestaurants List all resturants for user to choose
func GetOrdersList(w http.ResponseWriter, r *http.Request) {
	// user := context.Get(r, "userContext")
	restaurantID, err := uuid.Parse(r.URL.Query().Get("restaurant_id"))
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrdersList", "error": "invalid_request_body", "user_id": /*user.(utils.JWTInputData).UserID*/ ""}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid restaurant id")
		return
	}
	ordersList, err := controllers.DataService.GetOrdersOfRestaurant(restaurantID.String())
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrdersList", "error": "db_error", "user_id": /*user.(utils.JWTInputData).UserID*/ ""}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	ordersResponse := make([]response.OrderResponse, 0)
	for _, order := range ordersList {
		var orderRespItem response.OrderResponse
		orderRespItem.OrderID = order.ID
		orderRespItem.User = response.User{
			ID:   order.UserID,
			Name: order.UserName,
		}
		orderRespItem.DeliveryAddress = response.Address{
			AddressText: order.DeliveryAddressText,
			Lat:         order.DeliveryAddressLat,
			Lng:         order.DeliveryAddressLng,
		}
		var foodItems []response.OrderItem
		json.Unmarshal([]byte(order.Items), &foodItems)
		orderRespItem.Items = foodItems

		ordersResponse = append(ordersResponse, orderRespItem)
	}
	controllers.RespondWithSuccessGeneric(w, http.StatusOK, response.Response{
		Success: true,
		Data:    ordersResponse,
	})
	return
}

//GetAllRestaurants List all resturants for user to choose
func GetOrderEstimation(w http.ResponseWriter, r *http.Request) {
	// user := context.Get(r, "userContext")
	orderID, err := uuid.Parse(r.URL.Query().Get("order_id"))
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "invalid_request_body", "user_id": /*user.(utils.JWTInputData).UserID*/ ""}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid order id")
		return
	}
	orderDetail, err := controllers.DataService.GetOrderDetailsByID(orderID.String())
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "db_error", "user_id": /*user.(utils.JWTInputData).UserID*/ ""}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	if orderDetail.ID == "" {
		log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "invalid_request_body", "user_id": /*user.(utils.JWTInputData).UserID*/ ""}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid order id")
		return
	}

	restaurantDetail, err := controllers.DataService.GetRestaurantDetailsByID(orderDetail.RestaurantID)
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "db_error", "user_id": /*user.(utils.JWTInputData).UserID*/ ""}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	if restaurantDetail.ID == "" {
		log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "invalid_request_body", "user_id": /*user.(utils.JWTInputData).UserID*/ ""}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}

	var foodItems []response.OrderItem
	json.Unmarshal([]byte(orderDetail.Items), &foodItems)

	var estimatedPrice float64
	for _, foodItem := range foodItems {
		foodItemDetails, err := controllers.DataService.GetFoodDetailsByID(foodItem.FoodID)
		if err != nil {
			log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "db_error", "user_id": /*user.(utils.JWTInputData).UserID*/ ""}).Error(err.Error())
			controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
			return
		}
		estimatedPrice += float64(foodItem.Quantity) * foodItemDetails.UnitRateRupees
	}
	estimatedPriceWithTax := estimatedPrice + 0.05*estimatedPrice
	deliveryDistance := utils.Distance(restaurantDetail.Lat, restaurantDetail.Lng, orderDetail.DeliveryAddressLat, orderDetail.DeliveryAddressLng)
	totalPrice := estimatedPriceWithTax + deliveryDistance/1000

	controllers.RespondWithSuccessGeneric(w, http.StatusOK, response.Response{
		Success: true,
		Data: response.OrderEstimationResponse{
			EstimatedPrice: math.Round(totalPrice*100) / 100,
			EstimatedTime:  math.Round(deliveryDistance / (40 * 16.66)),
		},
	})
	return
}
