package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ReneGa/tweetcount-microservices/windower/gateway"
	"github.com/ReneGa/tweetcount-microservices/windower/resource"
	"github.com/ReneGa/tweetcount-microservices/windower/service"
	"github.com/julienschmidt/httprouter"
)

var address = flag.String("address", "localhost:8083", "Address to listen on")
var wordCountsURL = flag.String("wordCountsURL", "http://localhost:8082/wordcounts", "URL of the word counter to connect to")
var searchesURL = flag.String("searchesURL", "http://localhost:8084/searches", "URL of the searches service to connect to")

func main() {
	flag.Parse()

	tweetWordCountsGateway := &gateway.HTTPTweetWordCounts{
		Client: http.DefaultClient,
		URL:    *wordCountsURL,
	}
	searchesGateway := &gateway.HTTPSearches{
		Client: http.DefaultClient,
		URL:    *searchesURL,
	}
	windowService := service.NewWindow(tweetWordCountsGateway, searchesGateway)
	totalsResource := resource.Totals{
		Service: windowService,
	}

	router := httprouter.New()
	router.GET("/totals/:searchID", totalsResource.GET)

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
