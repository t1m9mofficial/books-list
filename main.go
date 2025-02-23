package main

// import
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	// "os"
	"fmt"

	"books-list/controllers"
	"books-list/driver"
	"books-list/models"
	"github.com/gorilla/mux"
	// "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

// Book model
// type Book struct {
// 	ID     int    `json:"id"`
// 	Title  string `json:"title"`
// 	Author string `json:"author"`
// 	Year   string `json:"year"`
// }

var books []models.Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	// logFatal(err)

	// db, err = sql.Open("postgres", pgURL)
	// logFatal(err)

	// err = db.Ping()
	// logFatal(err)

	db = driver.ConnectDB()
	controller := controllers.Controller{}

	router := mux.NewRouter()

	router.HandleFunc("/books", controller.GetBooks(db)).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	fmt.Println("Server is running at port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// func getBooks(w http.ResponseWriter, r *http.Request) {
// 	var book models.Book
// 	books = []models.Book{}

// 	rows, err := db.Query("select * from books")
// 	logFatal(err)

// 	defer rows.Close()

// 	for rows.Next() {
// 		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
// 		logFatal(err)

// 		books = append(books, book)
// 	}

// 	json.NewEncoder(w).Encode(books)
// }

func getBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	params := mux.Vars(r)

	rows := db.QueryRow("select * from books where id=$1", params["id"])
	rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)

	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)
	logFatal(err)

	json.NewEncoder(w).Encode(bookID)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id;", &book.Title, &book.Author, &book.Year, &book.ID)
	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	result, err := db.Exec("delete from books where id=$1", params["id"])
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}
