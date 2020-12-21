package router

import (
	"github.com/eltropycal/middlewares"
	"github.com/gorilla/mux"
)

// NewRouter returns a router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range internalRoutes {
		if route.Authenticate {
			router.Methods(route.Method).
				Path("/api" + route.Pattern).
				Name(route.Name).
				Handler(middlewares.Validate(route.HandlerFunc, route.Roles))
		} else {
			router.Methods(route.Method).
				Path("/api" + route.Pattern).
				Name(route.Name).
				Handler(route.HandlerFunc)
		}
	}
	return router
}
