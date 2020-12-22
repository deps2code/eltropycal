package router

import (
	"net/http"

	"github.com/eltropycal/controllers"
	"github.com/eltropycal/controllers/invoice"
	"github.com/eltropycal/controllers/orders"
	"github.com/eltropycal/controllers/restaurants"
)

// Route structure defining the api routes
type Route struct {
	Authenticate bool
	Roles        []int // 1 = restaurant admin, 2 = customer, 3 = driver
	Name         string
	Method       string
	Pattern      string
	HandlerFunc  http.HandlerFunc
}

// Routes defines the type Routes which is just an array (slice) of Route structs.
type Routes []Route

// Initialize our routes
var internalRoutes = Routes{
	Route{
		false,
		[]int{},
		"Backend server healthceck API",
		"GET",
		"/healthcheck",
		controllers.WebServerHealthCheck,
	},
	Route{
		false,
		[]int{},
		"Login API",
		"POST",
		"/login",
		controllers.UserLogin,
	},
	Route{
		true,
		[]int{1, 2, 3},
		"Logout API",
		"POST",
		"/logout",
		controllers.UserLogout,
	},
	Route{
		true,
		[]int{2},
		"List all restaurants API",
		"GET",
		"/restaurants",
		restaurants.GetAllRestaurants,
	},
	Route{
		true,
		[]int{2},
		"Get a restaurant menu API",
		"GET",
		"/restaurant-menu",
		restaurants.GetRestaurantMenu,
	},
	Route{
		true,
		[]int{2},
		"Create order API",
		"POST",
		"/order",
		orders.CreateOrder,
	},
	Route{
		true,
		[]int{1},
		"Get orders of Restaurant API",
		"GET",
		"/restaurant-orders",
		orders.GetOrdersList,
	},
	Route{
		true,
		[]int{2},
		"Get order estimation API",
		"GET",
		"/order-estimation",
		orders.GetOrderEstimation,
	},
	Route{
		true,
		[]int{2, 3},
		"Update order status API",
		"POST",
		"/order-status",
		orders.UpdateOrderStatus,
	},
	Route{
		true,
		[]int{3},
		"Get driver active order API",
		"GET",
		"/active-order",
		orders.GetOrderToDeliver,
	},
	Route{
		true,
		[]int{1},
		"Get Order invoice link API",
		"GET",
		"/order-invoice",
		invoice.GenerateOrderInvoice,
	},
}
