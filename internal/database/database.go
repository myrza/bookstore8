package server

import (
	"database/sql"
	//"fmt"
	"log"
	"net/http"

	//"os"

	"github.com/myrza/bookstore8/internal/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

/*
var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
	schema   = os.Getenv("DB_SCHEMA")
)
*/
// main function
func DatabaseConnect() error {
	//connect to database
	//db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", "postgres://postgres:postgres@db:5432/postgres?sslmode=disable")

	//connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	//log.Fatal(connStr)
	//log.Println("Database^ ", database)
	//db, err := sql.Open("postgres", connStr)

	if err != nil {
		//log.Fatal(err)
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
