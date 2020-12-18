package router

import (
	"net/http"

	"github.com/eltropycal/controllers"
)

// Route structure defining the api routes
type Route struct {
	AuthMiddleware int // 0 = no auth, 1 = Authenticate web session, 2 = authenticate else redirect, 3 = validate app request
	Name           string
	Method         string
	Pattern        string
	HandlerFunc    http.HandlerFunc
}

// Routes defines the type Routes which is just an array (slice) of Route structs.
type Routes []Route

// Initialize our routes
var internalRoutes = Routes{
	Route{
		0,
		"Backend server healthceck API ",
		"GET",
		"/healthcheck",
		controllers.WebServerHealthCheck,
	},
}
