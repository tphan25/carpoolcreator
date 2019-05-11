package main

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"googlemaps.github.io/maps"
)

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

func distanceMatrixRequest(originAddresses []string, destinationAddresses []string, apiKey string) (maps.DistanceMatrixResponse, error) {
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	r := &maps.DistanceMatrixRequest{
		Origins:       originAddresses,
		Destinations:  destinationAddresses,
		DepartureTime: `now`,
		Units:         `UnitsMetric`,
		Mode:          maps.TravelModeDriving,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//Returns DistanceMatrixResponse with fields from Dist Matrix API
	resp, err := c.DistanceMatrix(ctx, r)
	if err != nil {
		log.Fatal("Error making request")
	}
	return *resp, err
}

/*getDriversAndRiders will take a TripRead object (from our API request) and return Person objects corresponding to drivers and riders*/
func getDriversAndRiders(t TripRead) (driverList []Person, riderList []Person) {
	var drivers []Person
	var riders []Person
	//Check host first
	if t.Host.CanDrive {
		drivers = append(drivers, t.Host)
	} else {
		riders = append(riders, t.Host)
	}
	//Check guests
	for _, guest := range t.GuestList {
		if guest.CanDrive {
			drivers = append(drivers, guest)
		} else {
			riders = append(riders, guest)
		}
	}
	return drivers, riders
}

func getAddresses(drivers []Person, riders []Person) (driverAddresses []string, riderAddresses []string) {
	var d []string
	var r []string
	for _, driver := range drivers {
		d = append(d, driver.Address)
	}
	for _, rider := range riders {
		r = append(r, rider.Address)
	}
	return d, r
}

/*Calls routeByDriver using the API*/
func routeFromAPI(t TripRead, apiKey string) Trip {
	drivers, riders := getDriversAndRiders(t)
	driverAddresses, riderAddresses := getAddresses(drivers, riders)
	distResponse, err := distanceMatrixRequest(driverAddresses, riderAddresses, apiKey)
	if err != nil {
		fmt.Println("Distance matrix request failed")
		log.Fatal()
	}
	return routeByDriver(t, distResponse)
}

func routeFromFile(t TripRead, fileName string) Trip {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Could not open file (routeFromFile)")
	}
	byteArray, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Failed to read from file (testDistMatrix)")
	}

	var distResponse maps.DistanceMatrixResponse
	err = json.Unmarshal(byteArray, &distResponse)
	if err != nil {
		fmt.Println("Failed to decode as json (testDistMatrix)")
	}
	return routeByDriver(t, distResponse)
}

/*Takes TripRead object and converts to a Trip object (returned in HTTP response)*/
func routeByDriver(t TripRead, resp maps.DistanceMatrixResponse) Trip {
	var trip Trip
	drivers, riders := getDriversAndRiders(t)
	//driverAddresses, riderAddresses := getAddresses(drivers, riders)
	//directionsMatrixResponse, err := distanceMatrixRequest(driverAddresses, riderAddresses, os.Getenv("MAPS_CONNECT"))
	//We'll be using addresses for most work, so having a mapping of addresses back to names will be good
	mapAddressToName := make(map[string]string)
	for _, d := range drivers {
		mapAddressToName[d.Address] = d.Name
	}
	for _, r := range riders {
		mapAddressToName[r.Address] = r.Name
	}
	mapNameToIndex := make(map[string]int)
	for x, d := range drivers {
		mapNameToIndex[d.Name] = x
	}
	for y, r := range riders {
		mapNameToIndex[r.Name] = y
	}
	//BIG NOTE: If two people share an address we will have issues
	//To check if already taken, use a set, we will be using a map[string]bool
	ridersTaken := make(map[string]bool) //Use the address as a key, if true, rider is taken

	//For each driver, look through set of distances, and correspond those to a passenger to be added to their personal heap
	for i, driver := range drivers {
		passengerHeap := &PassengerHeap{}
		heap.Init(passengerHeap)
		for j, element := range resp.Rows[i].Elements {
			currRider := riders[j] //This is the current rider as a Person object
			//If rider is not taken, push it onto the current heap
			if !ridersTaken[currRider.Address] {
				heap.Push(passengerHeap, Passenger{
					address:  currRider.Address,
					distance: element.Distance.Meters,
				})

			}
		}
		fmt.Println(passengerHeap)
		//Fill up the car using min heap of passengers
		driver.Riders = make([]Person, driver.Capacity)
		for k := range driver.Riders {
			currPassenger := heap.Pop(passengerHeap).(Passenger)
			fmt.Println(currPassenger.distance)
			//Make a person object from Passsenger, or rather, find the passenger in riders
			//Find passenger in riders:
			passengerName := mapAddressToName[currPassenger.address]
			passengerIndex := mapNameToIndex[passengerName]
			driver.Riders[k] = riders[passengerIndex]
			ridersTaken[currPassenger.address] = true
		}
		//Honestly use pointers instead at some point though
		drivers[i] = driver
	}
	//At this point, we have an array of drivers
	//Find leftover passengers
	var extraGuests []Person
	for i, currRider := range riders {
		if !ridersTaken[currRider.Address] {
			extraGuests = append(extraGuests, riders[i])
		}
	}
	trip = Trip{
		Host:        t.Host,
		Name:        t.TripName,
		Location:    t.TripLocation,
		Date:        t.TripDate,
		Description: t.Description,
		Cars:        drivers,
		ExtraGuests: extraGuests,
	}
	return trip

}
