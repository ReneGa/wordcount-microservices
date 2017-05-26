package main

import (
	"database/sql"
	"log"

	"github.com/ReneGa/tweetcount-microservices/searches/repository"
)

func main() {

	db, err := sql.Open("sqlite3", "./searches.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	searchesRepo := repository.Searches{
		DB: db,
	}

	searchesResource := resource.Searches{
		SearchesRepository: searchesRepo,
	}

}
