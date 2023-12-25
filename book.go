package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Author      string             `json:"author" bson:"author"`
	Description string             `json:"description" bson:"description"`
	Available   bool               `json:"available" bson:"available"`
	BorrowID    string             `json:"borrowId,omitempty" bson:"borrowId,omitempty"`
}
