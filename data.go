package main

import (
	"context"
	"fmt"
	"os"
	"time"

	//"github.com/mongodb/mongo-go-driver/bson"
	//"github.com/mongodb/mongo-go-driver/mongo"
	//"github.com/mongodb/mongo-go-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//InsertTrip : this will insert a trip into the database, using testuser.
func InsertTrip(t Trip) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_CONNECT")))
	defer cancel()
	if err != nil {
		fmt.Println("Error: mongo.Connect failed")
		return err
	}
	collection := client.Database("carpool-db").Collection("trip-collection")
	b, err := bson.Marshal(t)
	if err != nil {
		fmt.Println("Error, bson.Marshal failed to convert trip to bson")
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, b)
	fmt.Println(res.InsertedID)
	defer cancel()
	if err != nil {
		fmt.Println("Error: Insertion into collection failed")
		return err
	}
	return nil
}

func testDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_CONNECT")))
	defer cancel()
	if err != nil {
		fmt.Println("Error: mongo.Connect failed")
		return err
	}
	collection := client.Database("carpool-db").Collection("trip-collection")
	b := bson.M{"name": "pi", "value": 3.14159}
	if err != nil {
		fmt.Println("Error, bson.Marshal failed to convert trip to bson")
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, b)
	if err != nil {
		fmt.Println("Error, insertion failed bruh")
		defer cancel()
		return err
	}
	fmt.Println(res.InsertedID)
	defer cancel()
	return err
}
