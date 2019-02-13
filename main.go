package main

import (
	"log"
	"net/http"

	"google.golang.org/appengine"
)

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
	http.Handle("/", router)
	appengine.Main()
	//)
}
