package main

type Book struct {
	ID          string `json:"id" bson:"_id"`
	Title       string `json:"title" bson:"title"`
	Author      string `json:"author" bson:"author"`
	Description string `json:"description" bson:"description"`
	Available   bool   `json:"available" bson:"available"`
}
