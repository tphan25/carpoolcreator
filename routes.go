package main

import "net/http"

/*Route will have its own name (like an identifier), method (GET, POST, etc) and the
pattern represents the route in URI. It also has a handler function to return values. */
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

/*Routes is just a slice of arrays.*/
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TripCreate",
		"POST",
		"/trips",
		TripCreate,
	},
	Route{
		"TestJson",
		"POST",
		"/testjson",
		TestJson,
	},
}
