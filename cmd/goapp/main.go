package main

import (
	"log"

	server "github.com/myrza/bookstore8/internal/database"
)

// main function
func main() {
	//connect to database
	err := server.DatabaseConnect()
	if err != nil {
		log.Fatal(err)
	}

}
