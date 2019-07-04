package main

import (
	"encoding/json"
	"fmt"
	"html"
	"os"

	"io"
	"io/ioutil"
	"net/http"
)

func openCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

/*Index is the current default path at our URL.*/
func Index(w http.ResponseWriter, r *http.Request) {
	openCors(&w)
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

/*TripCreate acts as a handler for POST to /trips, should return location of newly created trip*/
func TripCreate(w http.ResponseWriter, r *http.Request) {
	openCors(&w)
	fmt.Println("Trip create route reached")
	var t TripRead
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &t); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		//Will pass error if json encoding fails
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	//Use data from json in request to map drivers and riders
	drivers, riders := getDriversAndRiders(t)
	d, ri := getAddresses(drivers, riders)
	resp, err := distanceMatrixRequest(d, ri, os.Getenv("MAPS_CONNECT"))
	finalTrip := routeByDriver(t, resp)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Issue with distance matrix API")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	j, err := json.Marshal(finalTrip)
	if err != nil {
		fmt.Println("Failed to make json from trip (testCarpool)")
	}
	w.Write(j)
	//InsertTrip(tempTrip)
}

