package main

/*Person will have name, address, and candrive variable from form submission.*/
type Person struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	CanDrive bool   `json:"canDrive"`
}

/*Group is just a slice of people in each group, pretty much a car.*/
type Group []Person

/*Trip is going to be the majority of response body, as this is created product.*/
type Trip struct {
	Name        string  `json:"name"`
	Location    string  `json:"location"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Group       []Group `json:"group"`
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
