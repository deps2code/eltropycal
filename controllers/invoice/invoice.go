package invoice

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eltropycal/controllers"
	"github.com/eltropycal/models/response"
	"github.com/eltropycal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
)

func GenerateOrderInvoice(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "userContext")
	orderID, err := uuid.Parse(r.URL.Query().Get("order_id"))
	if err != nil {
		log.WithFields(log.Fields{"api": "GenerateOrderInvoice", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid order id")
		return
	}
	orderDetail, err := controllers.DataService.GetOrderDetailsByID(orderID.String())
	if err != nil {
		log.WithFields(log.Fields{"api": "GenerateOrderInvoice", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	if orderDetail.ID == "" {
		log.WithFields(log.Fields{"api": "GenerateOrderInvoice", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error("invalid order id")
		controllers.RespondWithError(w, http.StatusBadRequest, "Invalid order id")
		return
	}

	if orderDetail.Status < 3 || orderDetail.RestaurantID != user.(utils.UserJWTContext).RestaurantID {
		log.WithFields(log.Fields{"api": "GenerateOrderInvoice", "error": "error", "user_id": user.(utils.UserJWTContext).UserID}).Error("status is not yet delivered")
		controllers.RespondWithError(w, http.StatusBadRequest, "Order is not yet delivered, can't generate invoice.")
		return
	}
	restaurantDetail, err := controllers.DataService.GetRestaurantDetailsByID(user.(utils.UserJWTContext).RestaurantID)
	if err != nil {
		log.WithFields(log.Fields{"api": "GenerateOrderInvoice", "error": "db_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	if restaurantDetail.ID == "" {
		log.WithFields(log.Fields{"api": "GenerateOrderInvoice", "error": "invalid_request_body", "user_id": user.(utils.UserJWTContext).UserID}).Error("invalid restaurant id")
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	var orderRespItem response.OrderResponse
	orderRespItem.OrderID = orderDetail.ID
	orderRespItem.Status = orderDetail.Status
	orderRespItem.User = response.User{
		ID:   orderDetail.UserID,
		Name: orderDetail.UserName,
	}
	orderRespItem.DeliveryAddress = response.Address{
		AddressText: orderDetail.DeliveryAddressText,
		Lat:         orderDetail.DeliveryAddressLat,
		Lng:         orderDetail.DeliveryAddressLng,
	}
	orderRespItem.RestaurantAddress = response.Address{
		AddressText: restaurantDetail.AddressText,
		Lat:         restaurantDetail.Lat,
		Lng:         restaurantDetail.Lng,
	}
	var foodItems []response.OrderItem
	var foodResponseItem []response.FoodItem
	json.Unmarshal([]byte(orderDetail.Items), &foodItems)
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
	filename := "static/invoice-" + orderID.String() + "-" + time.Now().Format(time.RFC3339)
	err = utils.GenerateInvoicePdf(orderRespItem, filename)
	if err != nil {
		log.WithFields(log.Fields{"api": "GetOrdersList", "error": "pdf_error", "user_id": user.(utils.UserJWTContext).UserID}).Error(err.Error())
		controllers.RespondWithError(w, http.StatusInternalServerError, "Oops, something went wrong.")
		return
	}
	invoiceResp := response.InvoiceResponse{
		CreatedAt: time.Now().Format(time.RFC3339),
		Link:      "http://ec2-15-207-54-142.ap-south-1.compute.amazonaws.com:9090/" + filename,
	}
	controllers.RespondWithSuccessGeneric(w, http.StatusOK, response.Response{
		Success: true,
		Data:    invoiceResp,
	})
	return

}
