package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ReneGa/tweetcount-microservices/wordcounter/gateway"
	"github.com/ReneGa/tweetcount-microservices/wordcounter/resource"
	"github.com/julienschmidt/httprouter"
)

var address = flag.String("address", "localhost:8082", "Address to listen on")
var tweetsURL = flag.String("tweetsURL", "http://localhost:8081/tweets", "URL of the tweet producer to connect to")

func main() {
	flag.Parse()

	tweetsGateway := gateway.HTTPTweets{
		Client: http.DefaultClient,
		URL:    *tweetsURL,
	}
	wordCountsResource := resource.WordCounts{
		Gateway: &tweetsGateway,
	}

	router := httprouter.New()
	router.GET("/wordcounts", wordCountsResource.GET)

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
