package main

import (
	"context"
	"os"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

//This will insert a trip into the database, using testuser.
func InsertTrip(t Trip) error {
	dbUser := os.Getenv("CARPOOL_USER")
	dbPass := os.Getenv("CARPOOL_PASS")
	client, err := mongo.NewClient("mongodb+srv://" +
		dbUser + ":" + dbPass + "@cluster0-80i4c.gcp.mongodb.net/test?retryWrites=true")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	collection := client.Database("carpool-db").Collection("trip-collection")
	b, insertErr := bson.Marshal(t)
	if insertErr != nil {
		return insertErr
	}
	_, e := collection.InsertOne(context.Background(), b)
	if e != nil {
		return e
	}
	//id := res.InsertedID
	return nil
}
