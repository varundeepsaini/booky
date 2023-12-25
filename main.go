package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

var bookDB *mongo.Database

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := bookDB.Collection("booky")
	_, err = collection.InsertOne(context.TODO(), book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		return
	}
}

func BrowseBooks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []Book

	collection := bookDB.Collection("booky")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.TODO())

	for cursor.Next(context.TODO()) {
		var book Book
		err := cursor.Decode(&book)
		if err != nil {
			return
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	errorr := json.NewEncoder(w).Encode(books)
	if errorr != nil {
		return
	}
}

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookID := params["book_id"]

	collection := bookDB.Collection("booky")
	filter := bson.M{"_id": bookID}
	update := bson.M{"$set": bson.M{"available": true}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err2 := w.Write([]byte(`{"message":"Book returned successfully"}`))
	if err2 != nil {
		return
	}
}

func BorrowBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookID := params["book_id"]

	collection := bookDB.Collection("booky")
	filter := bson.M{"_id": bookID}
	update := bson.M{"$set": bson.M{"available": false}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err2 := w.Write([]byte(`{"message":"Book borrowed successfully"}`))
	if err2 != nil {
		return
	}
}

func main() {
	ConnectDB()
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/booky/", AddBook).Methods("PUT")
	r.HandleFunc("/api/v1/booky/", BrowseBooks).Methods("GET")
	r.HandleFunc("/api/v1/booky/{book_id}/borrow", BorrowBook).Methods("PUT")
	r.HandleFunc("/api/v1/booky/{book_id}/borrow/{borrow_id}", ReturnBook).Methods("POST")

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		return
	}
}
