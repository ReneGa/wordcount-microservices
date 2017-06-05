package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ReneGa/tweetcount-microservices/persister/datamapper"
	"github.com/ReneGa/tweetcount-microservices/persister/gateway"
	"github.com/ReneGa/tweetcount-microservices/persister/resource"
	"github.com/ReneGa/tweetcount-microservices/persister/service"
	"github.com/julienschmidt/httprouter"
)

var address = flag.String("address", "localhost:8085", "Address to listen on")
var bucketsDirectory = flag.String("bucketsDirectory", "./buckets", "Directory to write tweet buckets to.")
var tweetsURL = flag.String("tweetsURL", "http://localhost:8080/tweets", "URL of the tweet producer to connect to")

func main() {
	flag.Parse()

	queriesDataMapper := &datamapper.Queries{
		Directory: *bucketsDirectory,
	}
	tweetsGateway := &gateway.HTTPTweets{
		Client: http.DefaultClient,
		URL:    *tweetsURL,
	}
	tweetsService := service.Tweets{
		DataMapper: queriesDataMapper,
		Gateway:    tweetsGateway,
	}
	tweetsResource := resource.Tweets{
		Service: tweetsService,
	}

	router := httprouter.New()
	router.GET("/tweets", tweetsResource.GET)

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
