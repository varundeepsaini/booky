package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client *mongo.Client

var mongoUrl = "mongodb+srv://varundeepsaini:vlt0824@cluster0.3ay75yy.mongodb.net/books"

func ConnectDB() {
	// Update with your MongoDB URI
	clientOptions := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	bookDB = client.Database("booky")
	log.Println("Connected to MongoDB!")
}
