package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// define a struct to represent the event doc in the database
type Event struct {
	ID string `bson:"_id"`
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/event", eventHandler)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Here's our HTML button - modify UI as desired
	fmt.Fprintf(w, `<html>
	<head>
		<title>Event Button</title>
	</head>
	<body>
		<h1>Event Button Demo</h1>
		<button onclick="location.href='/event'">Click me!</button>
	</body>
	</html>`)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	// Below we generate an arbitary EventID. I want to use hashgen going forward, but need to test the package with staging DBM to make sure it acts as intended
	eventID := "abcd1234"

	// store the event ID in Mongo
	err := storeEventID(eventID)
	if err != nil {
		http.Error(w, "Failed to store event ID", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Event created! Event ID: %s", eventID)
}

func storeEventID(eventID string) error {
	// Set up the Mongo listener
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	// Access the events collection in the "test" database
	collection := client.Database("test").Collection("events")

	// Generate event doc and add to collection
	event := Event{
		ID: eventID,
	}

	_, err = collection.InsertOne(context.Background(), event)
	if err != nil {
		return err
	}

	return nil
}

// couple hanging threads here I want to comment on:
// for the local test I also have to initialize the package using:
// go get go.mongodb.org/mongo-driver/mongo
// this is a CLI function so will need to test the equivalent on the domain servr
// we can repeat this template ad infinitum for other generations
// next hurdle will be connecting generations to the ID refs (eg. attach userID to eventID)
// not sure if it makes more sense to focus on filling out the struct for all generations first
// or to lock in the distributed hashing first so we don't have to refactor when ID storage changes

