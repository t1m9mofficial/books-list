package controllers

import (
	"books-list/models"
	"log"
	"database/sql"
	"encoding/json"
	"net/http"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Controller type
type Controller struct{}

var books []models.Book

// GetBooks export
func (c Controller) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		books = []models.Book{}

		rows, err := db.Query("select * from books")
		logFatal(err)

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
			logFatal(err)

			books = append(books, book)
		}

		json.NewEncoder(w).Encode(books)
	}
}
