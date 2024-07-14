package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/myrza/bookstore8/internal/types"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM books")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		books := []types.Book{} // array of users
		for rows.Next() {
			var b types.Book
			if err := rows.Scan(&b.ID, &b.Title, &b.AuthorID, &b.ISBN, &b.Year); err != nil {
				//log.Fatal(err)
				json.NewEncoder(w).Encode(err)
			}
			books = append(books, b)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(books)
	}
}

// книга через id
func GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var b types.Book
		err := db.QueryRow("SELECT * FROM books WHERE id = $1", id).Scan(&b.ID, &b.Title, &b.AuthorID, &b.ISBN, &b.Year)
		if err != nil {
			//w.WriteHeader(http.StatusNotFound)
			//json.NewEncoder(w).Encode("Ошибка при попытке получить книгу по указанному идентификатору")
			json.NewEncoder(w).Encode("Ошибка при попытке получить книгу по указанному идентификатору")
			log.Println(err)
		} else {
			json.NewEncoder(w).Encode(b)
		}

	}
}

// создать книгу
func CreateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b types.Book
		json.NewDecoder(r.Body).Decode(&b)

		err := db.QueryRow("INSERT INTO books (title, authorid, isbn, year) VALUES ($1, $2, $3, $4) RETURNING id", b.Title, b.AuthorID, b.ISBN, b.Year).Scan(&b.ID)
		if err != nil {
			//log.Fatal(err)
			json.NewEncoder(w).Encode(err)
			return
		}

		json.NewEncoder(w).Encode(b)
	}
}

// обновить книгу

func UpdateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b types.Book
		json.NewDecoder(r.Body).Decode(&b)

		vars := mux.Vars(r)
		id := vars["id"]

		// Execute the update query
		_, err := db.Exec("UPDATE books SET title = $1, authorid = $2, isbn=$3, year = $4  WHERE id = $5", b.Title, b.AuthorID, b.ISBN, b.Year, id)
		if err != nil {
			json.NewEncoder(w).Encode("Автор не обнаружен")
		}

		// Retrieve the updated user data from the database
		var updatedBook types.Book
		err = db.QueryRow("SELECT * FROM books WHERE id = $1", id).Scan(&updatedBook.ID, &updatedBook.Title, &updatedBook.AuthorID, &updatedBook.ISBN, &updatedBook.Year)
		if err != nil {
			json.NewEncoder(w).Encode(err)
		}

		// Send the updated user data in the response
		json.NewEncoder(w).Encode(updatedBook)
	}
}

// удалить автора
func DeleteBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var b types.Book
		err := db.QueryRow("SELECT * FROM books WHERE id = $1", id).Scan(&b.ID, &b.Title, &b.AuthorID, &b.ISBN, &b.Year)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM books WHERE id = $1", id)
			if err != nil {
				json.NewEncoder(w).Encode(err)
				return
			}

			json.NewEncoder(w).Encode("книга удалена")
		}
	}
}
