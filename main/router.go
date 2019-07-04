package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

/* NewRouter creates a new mux.Router, and also wraps each route in our list with a Logger function
so that we can keep logs everytime a route is accessed. */
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
