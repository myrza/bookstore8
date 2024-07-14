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

func GetAuthors(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM authors")
		if err != nil {
			json.NewEncoder(w).Encode(err)
		}
		defer rows.Close()

		authors := []types.Author{}
		for rows.Next() {
			var a types.Author
			if err := rows.Scan(&a.ID, &a.Name, &a.Surname, &a.Biography, &a.Birthday); err != nil {
				json.NewEncoder(w).Encode(err)
			}
			authors = append(authors, a)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(authors)
	}
}

// Автор через id
func GetAuthor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var a types.Author
		err := db.QueryRow("SELECT * FROM authors WHERE id = $1", id).Scan(&a.ID, &a.Name, &a.Surname, &a.Biography, &a.Birthday)
		if err != nil {
			//w.WriteHeader(http.StatusNotFound)
			//return
			json.NewEncoder(w).Encode("Автор не обнаружен")
		}

		json.NewEncoder(w).Encode(a)
	}
}

// создать автора
func CreateAuthor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var a types.Author
		json.NewDecoder(r.Body).Decode(&a)

		err := db.QueryRow("INSERT INTO authors (name, surname, biography, birthday) VALUES ($1, $2, $3, $4) RETURNING id", a.Name, a.Surname, a.Biography, a.Birthday).Scan(&a.ID)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(a)
	}
}

// обновить автора
func UpdateAuthor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var a types.Author
		json.NewDecoder(r.Body).Decode(&a)

		vars := mux.Vars(r)
		id := vars["id"]

		// Execute the update query
		_, err := db.Exec("UPDATE authors SET name = $1, surname = $2,biography=$3, birthday = $4  WHERE id = $5", a.Name, a.Surname, a.Biography, a.Birthday, id)
		if err != nil {
			json.NewEncoder(w).Encode(err)
		}

		// Retrieve the updated user data from the database
		var updatedAuthor types.Author
		err = db.QueryRow("SELECT * FROM authors WHERE id = $1", id).Scan(&updatedAuthor.ID, &updatedAuthor.Name, &updatedAuthor.Surname, &updatedAuthor.Biography, &updatedAuthor.Birthday)
		if err != nil {
			json.NewEncoder(w).Encode(err)
		}

		// Send the updated user data in the response
		json.NewEncoder(w).Encode(updatedAuthor)
	}
}

// удалить автора
func DeleteAuthor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var a types.Author
		err := db.QueryRow("SELECT * FROM authors WHERE id = $1", id).Scan(&a.ID, &a.Name, &a.Surname, &a.Biography, &a.Birthday)
		if err != nil {
			json.NewEncoder(w).Encode("Автор не обнаружен")
		} else {
			_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
			if err != nil {
				json.NewEncoder(w).Encode(err)
				return
			}

			json.NewEncoder(w).Encode("Автор удален")
		}
	}
}
