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
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Fatal(err)
	}
	// Optional: Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
	bookDB = client.Database("books")
}
