package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/myrza/bookstore8/internal/types"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func GetBookAndAuthor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var b types.AuthorAndBook
		//err := db.QueryRow("SELECT name,surname, biography, birthday, title, ISBN, Year FROM books INNER JOIN authors ON books.authorid = authors.id WHERE books.id = $1", id).Scan(&b.Name, &b.Surname, &b.Biography, &b.Birthday, &b.AuthorID, &b.Title, &b.ISBN, &b.Year, id)
		err := db.QueryRow("SELECT name, surname, biography, birthday, title, authorid, isbn, year FROM books, authors  WHERE books.authorid = authors.id and  books.id = $1", id).Scan(&b.Name, &b.Surname, &b.Biography, &b.Birthday, &b.Title, &b.AuthorID, &b.ISBN, &b.Year)
		if err != nil {
			//w.WriteHeader(http.StatusNotFound)
			//json.NewEncoder(w).Encode("Ошибка при попытке получить книгу по указанному идентификатору")
			json.NewEncoder(w).Encode("Ошибка при попытке получить книгу по указанному идентификатору")
			json.NewEncoder(w).Encode(err)
			//log.Fatal(err)
			//return
		} else {
			json.NewEncoder(w).Encode(b)
		}

	}
}

// обновить книгу и автора в транзакции
func UpdateBookAndAuthor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ab types.AuthorAndBook
		json.NewDecoder(r.Body).Decode(&ab)

		vars := mux.Vars(r)
		book_id := vars["id"]

		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.ExecContext(ctx, "UPDATE books SET title = $1, authorid = $2, isbn=$3, year = $4  WHERE id = $5", ab.Title, ab.AuthorID, ab.ISBN, ab.Year, book_id)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			log.Fatal("\n", (err), "\n ....Transaction rollback!\n")

			//return json(500, "test error")
			json.NewEncoder(w).Encode("Transaction rollback!")
		}

		// Execute the update query
		_, err = db.ExecContext(ctx, "UPDATE authors SET name = $1, surname = $2,biography=$3, birthday = $4  WHERE id = $5", ab.Name, ab.Surname, ab.Biography, ab.Birthday, ab.AuthorID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			log.Fatal("\n", (err), "\n ....Transaction rollback!\n")
			return
		}

		// close the transaction with a Commit() or Rollback() method on the resulting Tx variable.
		// this applies the above changes to our database
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
			json.NewEncoder(w).Encode("Ошибка. Произведен откат транзакции")
		} else {
			json.NewEncoder(w).Encode("Транзакция выполнена")
		}

	}
}
