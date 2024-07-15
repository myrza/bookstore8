package server

import (
	"database/sql"
	"fmt"

	//"fmt"
	"log"
	"net/http"

	"github.com/myrza/bookstore8/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func DatabaseConnect() error {
	envFile, _ := godotenv.Read(".env")

	database := envFile["DB_DATABASE"]
	password := envFile["DB_PASSWORD"]
	username := envFile["DB_USERNAME"]
	port := envFile["DB_PORT"]
	host := envFile["DB_HOST"]

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return err
	}
	defer db.Close()

	// create table if not exists
	err = TableCreate(db)
	if err != nil {
		return err
	}
	// create router
	router := mux.NewRouter()

	// операции по авторам книг
	router.HandleFunc("/api/go/authors", handlers.GetAuthors(db)).Methods("GET")
	router.HandleFunc("/api/go/authors", handlers.CreateAuthor(db)).Methods("POST")
	router.HandleFunc("/api/go/authors/{id}", handlers.GetAuthor(db)).Methods("GET")
	router.HandleFunc("/api/go/authors/{id}", handlers.UpdateAuthor(db)).Methods("PUT")
	router.HandleFunc("/api/go/authors/{id}", handlers.DeleteAuthor(db)).Methods("DELETE")

	// операции по книгам
	router.HandleFunc("/api/go/books", handlers.GetBooks(db)).Methods("GET")
	router.HandleFunc("/api/go/books", handlers.CreateBook(db)).Methods("POST")
	router.HandleFunc("/api/go/books/{id}", handlers.GetBook(db)).Methods("GET")
	router.HandleFunc("/api/go/books/{id}", handlers.UpdateBook(db)).Methods("PUT")
	router.HandleFunc("/api/go/books/{id}", handlers.DeleteBook(db)).Methods("DELETE")

	// транзакция
	router.HandleFunc("/api/go/book-author/{id}", handlers.UpdateBookAndAuthor(db)).Methods("PUT")
	router.HandleFunc("/api/go/book-author/{id}", handlers.GetBookAndAuthor(db)).Methods("GET")

	// wrap the router with CORS and JSON content type middlewares
	enhancedRouter := handlers.EnableCORS(handlers.JsonContentTypeMiddleware(router))

	// start server
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
	return nil

}
