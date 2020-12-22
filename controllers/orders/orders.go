package orders

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/eltropycal/constants"
	"github.com/eltropycal/controllers"
	"github.com/eltropycal/models/dbmodels"
	"github.com/eltropycal/models/request"
	"github.com/eltropycal/models/response"
	"github.com/eltropycal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "userContext")
	var apiReq request.OrderRequest
	err := json.NewDecoder(r.Body).Decode(&apiReq)
	if err != nil {
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "invalid_api_request", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err = apiReq.Validate(); err != nil {
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	foodItems, err := json.Marshal(apiReq.Items)
	if err != nil {
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "food items are not correct")
		return
	}

	orderID := uuid.New().String()
	order := dbmodels.Order{
		ID:                  orderID,
		UserID:              user.(utils.UserJWTContext).UserID,
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
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong")
		return
	}

	//start go routine for assigniing driver
	go assignDriverToOrder(orderID, user)
	controllers.RespondWithSuccess(w, http.StatusOK, "We have received your order, it will be delivered to you soon")
	return

}

//GetOrdersList List all orders for a restaurant
func GetOrdersList(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "userContext")
	restaurantID := user.(utils.UserJWTContext).RestaurantID
	ordersList, err := controllers.DataService.GetOrdersOfRestaurant(restaurantID)
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrdersList", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	ordersResponse := make([]response.OrderResponse, 0)
	for _, order := range ordersList {
		var orderRespItem response.OrderResponse
		orderRespItem.OrderID = order.ID
		orderRespItem.Status = order.Status
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
		var foodResponseItem []response.FoodItem
		json.Unmarshal([]byte(order.Items), &foodItems)
		for _, foodItem := range foodItems {
			foodItemDetails, err := controllers.DataService.GetFoodDetailsByID(foodItem.FoodID)
			if err != nil {
				log.WithFields(log.Fields{"api": "GetOrdersList", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
				controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
				return
			}
			if foodItem.Quantity > 0 {
				foodResponseItem = append(foodResponseItem, response.FoodItem{
					FoodName:    foodItemDetails.Name,
					Description: foodItemDetails.Description,
					Quantity:    foodItem.Quantity,
					Price:       foodItemDetails.UnitRateRupees,
				})
			}
		}
		orderRespItem.Items = foodResponseItem

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
	user := context.Get(r, "userContext")
	orderID, err := uuid.Parse(r.URL.Query().Get("order_id"))
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid order id")
		return
	}
	orderDetail, err := controllers.DataService.GetOrderDetailsByID(orderID.String())
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	if orderDetail.ID == "" {
		log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid order id")
		return
	}

	restaurantDetail, err := controllers.DataService.GetRestaurantDetailsByID(orderDetail.RestaurantID)
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	if restaurantDetail.ID == "" {
		log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}

	var foodItems []response.OrderItem
	json.Unmarshal([]byte(orderDetail.Items), &foodItems)

	var estimatedPrice float64
	for _, foodItem := range foodItems {
		foodItemDetails, err := controllers.DataService.GetFoodDetailsByID(foodItem.FoodID)
		if err != nil {
			log.WithFields(log.Fields{"api": "GetOrderEstimation", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
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

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "userContext")
	var apiReq request.OrderStatusUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&apiReq)
	if err != nil {
		log.WithFields(log.Fields{"api": "UpdateOrderStatus", "error": "invalid_api_request", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err = apiReq.Validate(user); err != nil {
		log.WithFields(log.Fields{"api": "UpdateOrderStatus", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	orderDetail, err := controllers.DataService.GetOrderDetailsByID(apiReq.OrderID)
	if err != nil {
		log.WithFields(log.Fields{"api": "UpdateOrderStatus", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	if orderDetail.ID == "" {
		log.WithFields(log.Fields{"api": "UpdateOrderStatus", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid order id")
		return
	}

	if apiReq.Status == -1 && orderDetail.UserID != user.(utils.UserJWTContext).UserID {
		log.WithFields(log.Fields{"api": "UpdateOrderStatus", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Looks like its not your order")
		return
	}

	if apiReq.Status != -1 && orderDetail.DriverID != user.(utils.UserJWTContext).UserID {
		log.WithFields(log.Fields{"api": "UpdateOrderStatus", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Looks like its not your order")
		return
	}

	if apiReq.Status == -1 {
		err := controllers.DataService.UpdateOrderStatus(constants.CANCELLED, apiReq.OrderID, orderDetail.DriverID)
		if err != nil {
			log.WithFields(log.Fields{"api": "UpdateOrderStatus", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
			controllers.RespondWithError(w, http.StatusBadRequest, "Could not cancel order at this time")
			return
		}
		controllers.RespondWithSuccess(w, http.StatusOK, "Order status updated")
		return
	}
	if apiReq.Status < orderDetail.Status {
		controllers.RespondWithSuccess(w, http.StatusOK, "You cannot change the status back.")
		return
	}
	err = controllers.DataService.UpdateOrderStatus(apiReq.Status, apiReq.OrderID, orderDetail.DriverID)
	if err != nil {
		log.WithFields(log.Fields{"api": "UpdateOrderStatus", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Could not update order status at this time")
		return
	}
	controllers.RespondWithSuccess(w, http.StatusOK, "Order status updated")
	return
}

func GetOrderToDeliver(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "userContext")
	order, err := controllers.DataService.GetActiveOrderForDriver(user.(utils.UserJWTContext).UserID)
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrderToDeliver", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	if order.ID == "" {
		log.WithFields(log.Fields{"api": "GetOrderToDeliver", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Info("No orders")
		controllers.RespondWithSuccess(w, http.StatusOK, "No orders yet")
		return
	}
	var orderRespItem response.OrderResponse
	orderRespItem.OrderID = order.ID
	orderRespItem.Status = order.Status
	orderRespItem.User = response.User{
		ID:   order.UserID,
		Name: order.UserName,
	}
	orderRespItem.DeliveryAddress = response.Address{
		AddressText: order.DeliveryAddressText,
		Lat:         order.DeliveryAddressLat,
		Lng:         order.DeliveryAddressLng,
	}
	orderRespItem.RestaurantAddress = response.Address{
		AddressText: order.RestaurantAddressText,
		Lat:         order.RestaurantAddressLat,
		Lng:         order.RestaurantAddressLng,
	}
	var foodItems []response.OrderItem
	var foodResponseItem []response.FoodItem
	json.Unmarshal([]byte(order.Items), &foodItems)
	for _, foodItem := range foodItems {
		foodItemDetails, err := controllers.DataService.GetFoodDetailsByID(foodItem.FoodID)
		if err != nil {
			log.WithFields(log.Fields{"api": "GetOrdersList", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
			controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
			return
		}
		if foodItem.Quantity > 0 {
			foodResponseItem = append(foodResponseItem, response.FoodItem{
				FoodName:    foodItemDetails.Name,
				Description: foodItemDetails.Description,
				Quantity:    foodItem.Quantity,
				Price:       foodItemDetails.UnitRateRupees,
			})
		}
	}
	orderRespItem.Items = foodResponseItem
	controllers.RespondWithSuccessGeneric(w, http.StatusOK, response.Response{
		Success: true,
		Data:    orderRespItem,
	})
	return
}

func assignDriverToOrder(orderID string, user interface{}) {
	driverID, err := controllers.DataService.GetAvailableDriver()
	if err != nil {
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		time.AfterFunc(1*time.Minute, func() {
			assignDriverToOrder(orderID, user)
		})
	}
	if driverID == "" {
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "no free driver", "user_id": user.(utils.UserJWTContext).UserID}).Error("Couldn't find driver id")
		time.AfterFunc(1*time.Minute, func() {
			assignDriverToOrder(orderID, user)
		})
	}
	err = controllers.DataService.AssignDriver(driverID, orderID)
	if err != nil {
		log.WithFields(log.Fields{"api": "CreateOrder", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		time.AfterFunc(1*time.Minute, func() {
			assignDriverToOrder(orderID, user)
		})
	}
}
