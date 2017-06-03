package main

import (
	"database/sql"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"flag"
	"log"

	"github.com/ReneGa/tweetcount-microservices/searches/datamapper"
	"github.com/ReneGa/tweetcount-microservices/searches/resource"
	"github.com/julienschmidt/httprouter"
)

var address = flag.String("address", "localhost:8084", "Address to listen on")

func main() {
	db, err := sql.Open("sqlite3", "./searches.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	searchesDataMapper := datamapper.Searches{
		DB: db,
	}

	searchesDataMapper.Init()

	searchesResource := resource.Searches{
		SearchesDataMapper: &searchesDataMapper,
	}

	router := httprouter.New()
	router.GET("/searches", searchesResource.GetAll)
	router.GET("/searches/:searchID", searchesResource.Get)
	router.POST("/searches", searchesResource.Create)

	done := make(chan bool)
	go func() {
		err := http.ListenAndServe(*address, router)
		if err != nil {
			log.Fatal(err)
		}
		done <- true
	}()
	log.Printf("listening on %s", *address)
	<-done
}
