package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Foobar struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

func TestJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	foobar := Foobar{
		Foo: "foo1",
		Bar: "bar1",
	}
	//s, _ := json.Marshal(foobar)

	json.NewEncoder(w).Encode(foobar)
}

func tripReadJsonFromFile(s string) TripRead {
	var t TripRead
	file, err := os.Open(s)
	if err != nil {
		fmt.Printf("Could not open file (testDistMatrix): %s\n", s)
	}
	byteArray, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Failed to read from file (testDistMatrix)")
	}
	err = json.Unmarshal(byteArray, &t)
	if err != nil {
		fmt.Println("Failed to decode as json (testDistMatrix)")
	}
	return t
}

func writeDistanceMatrixToFile(t TripRead) {
	drivers, riders := getDriversAndRiders(t)
	d, r := getAddresses(drivers, riders)
	// fmt.Println(drivers)
	// fmt.Println(riders)
	// fmt.Println(d)
	// fmt.Println(r)
	resp, err := distanceMatrixRequest(d, r, os.Getenv("MAPS_CONNECT"))
	if err != nil {
		fmt.Println("Failed to make request")
	}
	f, err := os.Create("DistMatrix.json")
	if err != nil {
		fmt.Println("Failed to create file")
	}
	j, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Failed to marshal JSON")
	}
	f.Write(j)
}

func testCarpool(t TripRead, filename string) {
	finalTrip := routeFromFile(t, filename)
	f, err := os.Create("FinalTrip.json")
	if err != nil {
		fmt.Println("Failed to create file (testCarpool)")
	}
	j, err := json.Marshal(finalTrip)
	if err != nil {
		fmt.Println("Failed to make json from trip (testCarpool)")
	}
	f.Write(j)
}
