package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	//"google.golang.org/api/transport/http"
)

func main() {
	dbUser := os.Getenv("CARPOOL_USER")
	dbPass := os.Getenv("CARPOOL_PASS")
	fmt.Print(dbUser + ":" + dbPass)
	router := NewRouter()
	//If hosting locally, uncomment the following line:
	log.Fatal(http.ListenAndServe(":8080", router))
	http.Handle("/", router)
	//appengine.Main()
	//)
}
