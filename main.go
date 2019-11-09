package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
	"os"
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book

var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	// TODO:
	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANT_SQL_URL"))
	logFatalError(err)

	db, err = sql.Open("postgres", pgUrl)
	logFatalError(err)
	err = db.Ping()
	logFatalError(err)

	log.Println(pgUrl)

	// Define the router.
	router := mux.NewRouter()

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
	var book Book
	books = []Book{}

	rows, err := db.Query("SELECT * from books")
	logFatalError(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatalError(err)

		books = append(books, book)
	}

	json.NewEncoder(w).Encode(&books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.Vars(r)

	// get it from the database.
	rows := db.QueryRow("select * from books where id =" + params["id"])

	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatalError(err)
	json.NewEncoder(w).Encode(&book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	// Adding the book model to the database.
	var book Book
	var bookId int
	json.NewDecoder(r.Body).Decode(&book)

	// Insert to the database.
	err := db.QueryRow("INSERT INTO books (title, author, year) VALUES ($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookId)

	logFatalError(err)

	json.NewEncoder(w).Encode(&bookId)

}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	json.NewDecoder(r.Body).Decode(&book)

	// database query

	result, err := db.Exec("UPDATE books set title=$1, author=$2, year=$3 WHERE id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, &book.ID)

	logFatalError(err)

	rowsAffected, err := result.RowsAffected()

	logFatalError(err)

	json.NewEncoder(w).Encode(rowsAffected)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {

	// get the id from the request
	params := mux.Vars(r)


	result, err := db.Exec("DELETE FROM books WHERE id=$1", params["id"])
	logFatalError(err)


	rowsAffected, err := result.RowsAffected()


	logFatalError(err)

	json.NewEncoder(w).Encode(rowsAffected)

}
