package main

import (
	"log"
	"net/http"
	//"google.golang.org/api/transport/http"
)

func main() {
	router := NewRouter()
	//If hosting locally, uncomment the following line:
	log.Fatal(http.ListenAndServe(":8080", router))
	http.Handle("/", router)
	//appengine.Main()
	//)
}
