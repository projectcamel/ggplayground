package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define a struct to represent the event doc in the database
type Event struct {
	ID   string    `bson:"_id"`
	Date time.Time `bson:"date"`
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
	// Generate the event ID hash (playing with MD5, this should be SHA)
	eventID := generateEventID()

	// Set the event date
	eventDate := time.Now().AddDate(0, 1, 0) // Set the event date to one month from now

	// Store the event in Mongo
	err := storeEvent(eventID, eventDate)
	if err != nil {
		http.Error(w, "Failed to store event", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Event button clicked! Event ID: %s, Date: %s", eventID, eventDate.String())
}

func generateEventID() string {
	// Here's we're using the timestamp to generate the hash
	timestamp := time.Now().String()
	hash := md5.Sum([]byte(timestamp))

	return hex.EncodeToString(hash[:])
}

func storeEvent(eventID string, eventDate time.Time) error {
	// Set up Mongo listener
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	// Access the events collection in the "test" database
	collection := client.Database("test").Collection("events")

	// Create the event doc and add to the collection
	event := Event{
		ID:   eventID,
		Date: eventDate,
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
