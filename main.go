package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   int    `json:year`
}

var books []Book

func main() {
	// TODO:
	// Define the router.
	router := mux.NewRouter()

	// Create a list of books for astart
	books = append(books,
		Book{ID: 1, Title: "Golang Pointers", Author: "Maunda Alex", Year: 2009},
		Book{ID: 2, Title: "Go Routines", Author: "Maunda Alex", Year: 2009},
		Book{ID: 3, Title: " Golang Routers", Author: "Maunda Alex", Year: 2009},
		Book{ID: 4, Title: "Golang Concurrency", Author: "Maunda Alex", Year: 2009},
		Book{ID: 5, Title: "Golang the Good Parts", Author: "Maunda Alex", Year: 2009})

	// Register the http routes.
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	// Start the server
	log.Fatal(http.ListenAndServe(":2000", router))
}

// define the methods to handle the requests
func getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting all the books")

	json.NewEncoder(w).Encode(books)

}

func getBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting a single book")
	params := mux.Vars(r)
	// Find the book form the books array
	id, _ := strconv.Atoi(params["id"])
	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(&book)
		}
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding a new book")
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	log.Println(book)
	books = append(books, book)
	json.NewEncoder(w).Encode(&books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("updating a book")

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	for i, item := range books {
		if item.ID == book.ID {
			books[i] = book
		}
	}

	json.NewEncoder(w).Encode(&books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting the book")

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	for i, item := range books {
		if item.ID == id {
			books = append(books[:i], books[i+1:]...)
		} else {
			log.Panic("Book not found")
			json.NewEncoder(w).Encode("The book passed was not found")
		}
	}
	json.NewEncoder(w).Encode(books)
}
