package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	collection := bookDB.Collection("books")
	result, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the generated ID to the book object
	book.ID = result.InsertedID.(primitive.ObjectID)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func BrowseBooks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []Book

	collection := bookDB.Collection("books")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println(err)
		}
	}(cursor, context.Background())

	for cursor.Next(context.Background()) {
		var book Book
		if err = cursor.Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	if err = cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		return
	}
}

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookID := params["book_id"]
	borrowID := params["borrow_id"]

	collection := bookDB.Collection("books")
	objID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		http.Error(w, "Invalid book ID format", http.StatusBadRequest)
		return
	}

	var book Book
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "No book found with the given ID", http.StatusNotFound)
			return
		}
		http.Error(w, "An error occurred while retrieving the book", http.StatusInternalServerError)
		return
	}

	if book.BorrowID != borrowID {
		http.Error(w, "Incorrect borrow ID", http.StatusBadRequest)
		return
	}

	update := bson.M{"$set": bson.M{"available": true, "borrowId": ""}}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID, "borrowId": borrowID}, update)
	if err != nil {
		http.Error(w, "Failed to update the book as returned", http.StatusInternalServerError)
		return
	}

	message := "Book returned successfully"
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"message": message})
	if err != nil {
		return
	}
}

func BorrowBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookID := params["book_id"]

	collection := bookDB.Collection("books")
	objID, _ := primitive.ObjectIDFromHex(bookID)

	var book Book
	err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "No book found with the given ID", http.StatusNotFound)
			return
		}
		http.Error(w, "An error occurred while retrieving the book", http.StatusInternalServerError)
		return
	}

	if !book.Available {
		http.Error(w, "Book already borrowed", http.StatusBadRequest)
		return
	}

	borrowId := uuid.New().String()
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"available": false,
			"borrowId":  borrowId,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update the book as borrowed", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":  "Book borrowed successfully",
		"borrowId": borrowId,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
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

	log.Fatal(http.ListenAndServe(":8000", r))
}
