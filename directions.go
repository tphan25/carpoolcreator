package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"googlemaps.github.io/maps"
)

// DistanceMatrixResponse represents a Distance Matrix API response.
type DistanceMatrixResponse struct {

	// OriginAddresses contains an array of addresses as returned by the API from
	// your original request.
	OriginAddresses []string `json:"origin_addresses"`
	// DestinationAddresses contains an array of addresses as returned by the API
	// from your original request.
	DestinationAddresses []string `json:"destination_addresses"`
	// Rows contains an array of elements.
	Rows []DistanceMatrixElementsRow `json:"rows"`
}

// DistanceMatrixElementsRow is a row of distance elements.
type DistanceMatrixElementsRow struct {
	Elements []*DistanceMatrixElement `json:"elements"`
}

// DistanceMatrixElement is the travel distance and time for a pair of origin
// and destination.
type DistanceMatrixElement struct {
	Status string `json:"status"`
	// Duration is the length of time it takes to travel this route.
	Duration time.Duration `json:"duration"`
	// DurationInTraffic is the length of time it takes to travel this route
	// considering traffic.
	DurationInTraffic time.Duration `json:"duration_in_traffic"`
	// Distance is the total distance of this route.
	Distance Distance `json:"distance"`
}

type Distance struct {
	value int    `json:"value"`
	text  string `json:"text"`
}

func testAPI() error {
	c, _ := maps.NewClient(maps.WithAPIKey(os.Getenv("MAPS_CONNECT")))
	r := &maps.DistanceMatrixRequest{
		Origins:       []string{"1.315125,103.76471334"},
		Destinations:  []string{"1.280776,103.8487"},
		DepartureTime: `now`,
		Units:         `UnitsMetric`,
		Mode:          maps.TravelModeDriving,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := c.DistanceMatrix(ctx, r)
	defer cancel()
	if err != nil {
		fmt.Println("r.Get returned non nil error, was")
		return err
	}
	fmt.Println(resp)
	return err
}

func routeByDriver(t TripRead) (Trip, error) {
	var trip Trip
	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("MAPS_CONNECT")))
	if err != nil {
		fmt.Println("Error connecting to Directions API")
		return trip, err
	}
	driverAddresses := make([]string, 10)
	riderAddresses := make([]string, 10)

	drivers := make([]Person, 10)
	riders := make([]Person, 10)

	if t.Host.CanDrive {
		driverAddresses = append(driverAddresses, t.Host.Address)
		drivers = append(drivers, t.Host)
	} else {
		riderAddresses = append(riderAddresses, t.Host.Address)
		riders = append(riders, t.Host)
	}

	for _, guest := range t.GuestList {
		if guest.CanDrive {
			driverAddresses = append(driverAddresses, guest.Address)
			drivers = append(drivers, guest)
		} else {
			riderAddresses = append(riderAddresses, guest.Address)
			riders = append(riders, guest)
		}
	}

	r := &maps.DistanceMatrixRequest{
		Origins:       driverAddresses,
		Destinations:  riderAddresses,
		DepartureTime: `now`,
		Units:         `UnitsMetric`,
		Mode:          maps.TravelModeDriving,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//Returns DistanceMatrixResponse with fields from Dist Matrix API
	resp, err := c.DistanceMatrix(ctx, r)

	/*var riderTaken []bool

	//Populate drivers, O(n * m) right now so consider removing passengers as they are added
	//Current idea: For each driver, add each element into list, and then find lowest x elements (capacity)
	//Add them to driver's car inside trip object & set boolean of array to true for respective rider (meaning no other driver)
	for i, driver := range drivers {
		//Gets corresponding row for driver
		currRow := resp.Rows[i].Elements
		currPassengers := make([]Passenger, len(driver.Riders))
		for j, element := range currRow {

		}
	}*/

	fmt.Println(resp)
	defer cancel()

	return trip, err
}
