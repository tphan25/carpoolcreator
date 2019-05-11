package main

/*Person will have name, address, and candrive variable from form submission.
  If the user is a driver, they will have a slice field "Riders" representing people in their vehicle.*/
type Person struct {
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	CanDrive bool     `json:"canDrive"`
	Capacity int      `json:"capacity"`
	Riders   []Person `json:"riders"`
}

type Rider struct {
	person   Person
	distance int
}

/*Group is just a slice of people in each group, pretty much a car.*/
type Group []Person

/*Trip is going to be the majority of response body, as this is created product.*/
type Trip struct {
	Host        Person   `json:"host"`
	Name        string   `json:"name"`
	Location    string   `json:"location"`
	Date        string   `json:"date"`
	Description string   `json:"description"`
	Cars        []Person `json:"cars"`
	ExtraGuests []Person `json:"extraGuests"`
}

/*TripRead represents the request body's json of a trip, before processing. */
type TripRead struct {
	Host         Person   `json:"host"`
	GuestList    []Person `json:"guestList"`
	TripName     string   `json:"tripName"`
	TripLocation string   `json:"tripLocation"`
	TripDate     string   `json:"tripDate"`
	Description  string   `json:"description"`
}

/*ProcessTrip will take a TripRead (i.e. from HTTP POST) and convert it into a Trip, to be stored
To be deprecated once directionsmatrix done*/
/*
func ProcessTrip(t TripRead) Trip {
	var trip Trip
	trip.Host = t.Host
	trip.Name = t.TripName
	trip.Location = t.TripLocation
	trip.Date = t.TripDate
	trip.Description = t.Description

	people := append(t.GuestList, t.Host)
	var drivers []Person
	for i, p := range people {
		//Separate drivers and nondrivers
		if p.CanDrive {
			drivers = append(drivers, p)
			people = append(people[:i], people[i+1:]...)
		}
	}
	groups := make([]Group, len(drivers))
	//for now, just group people into 4
	for i, d := range drivers {
		groups[i] = append(groups[i], d)
		for j := 0; j < 4; j++ {
			//This will append 4 people from the car, and remove them from the full persons list
			if len(people) > 0 {
				groups[i] = append(groups[i], people[0])
				people = append(people[:0], people[1:]...)
			}

		}
	}

	trip.Cars = groups
	trip.ExtraGuests = people
	return trip
}
*/
