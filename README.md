
# Booky - Book Sharing Microservice 

This project is a backend microservice for sharing books. It allows users to add, browse, borrow, and return books. The service is built using Go and MongoDB.

## Features

- **Add Books**: Users can add new books to the system.
- **Browse Books**: Users can view all books available for borrowing.
- **Borrow Books**: Users can borrow available books.
- **Return Books**: Users can return books they have borrowed.

## Getting Started

These instructions will help you get the project up and running on your local machine for development and testing purposes.

### Prerequisites

- [Go](https://golang.org/dl/) (Version 1.x or later)
- [MongoDB](https://www.mongodb.com/try/download/community)
- [Postman](https://www.postman.com/downloads/) (for API testing)

### Installation

1. **Clone the repository:**
   ```
   git clone https://github.com/yourusername/book-sharing-microservice.git
   ```
2. **Navigate to the project directory:**
   ```
   cd book-sharing-microservice
   ```
3. **Run the service:**
   ```
   go run .
   ```

### Environment Variables

Set up the following environment variable:
- `MONGODB_URI`: The URI for your MongoDB instance.

## API Endpoints

The service exposes several endpoints for interacting with books:

### Add a Book

- **URL**: `/api/v1/booky/`
- **Method**: `PUT`
- **Body** (JSON):
  ```json
  {
    "title": "String",
    "author": "String",
    "description": "String",
    "available": Boolean
  }
  ```

### Browse Books

- **URL**: `/api/v1/booky/`
- **Method**: `GET`

### Borrow a Book

- **URL**: `/api/v1/booky/{book_id}/borrow`
- **Method**: `PUT`
- **URL Parameters**: `book_id=[string]`

### Return a Book

- **URL**: `/api/v1/booky/{book_id}/borrow/{borrow_id}`
- **Method**: `POST`
- **URL Parameters**:
    - `book_id=[string]`
    - `borrow_id=[string]`

## Testing

Use Postman or a similar API testing tool to test the above endpoints.

## Built With

- [Go](https://golang.org/) - The programming language used.
- [MongoDB](https://www.mongodb.com/) - Database used.
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP routing and URL matcher.
