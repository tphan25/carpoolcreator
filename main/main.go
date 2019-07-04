package main

//"google.golang.org/api/transport/http"
import (
	"CarpoolCreatorAPI/authentication"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := authentication.Init()
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := NewRouter()
	//If hosting locally, uncomment the following line:
	log.Fatal(http.ListenAndServe(":8080", router))
	http.Handle("/", router)
	// t := tripReadJsonFromFile("TripRead.json")
	// //Costs an API request
	// writeDistanceMatrixToFile(t)
	// testCarpool(t, "DistMatrix.json")
	//testHeap()
}
