package router

import (
	"github.com/gorilla/mux"
)

// NewRouter returns a router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range internalRoutes {
		switch route.AuthMiddleware {
		case 0:
			router.Methods(route.Method).
				Path("/api" + route.Pattern).
				Name(route.Name).
				Handler(route.HandlerFunc)
			break

		}
	}
	return router
}
