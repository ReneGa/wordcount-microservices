package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/datamapper"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/gateway"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/repository"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/resource"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/service"
	"github.com/julienschmidt/httprouter"
)

func requireEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Sprintf("Required environment variable '%s' not set.", name))
	}
	return value
}

var address = flag.String("address", "localhost:8081", "Address to listen on")
var tweetsURL = flag.String("tweetsURL", "http://localhost:8080/tweets", "URL of the tweet producer to connect to")
var stopWordsDirectory = flag.String("stopWordsDirectory", "stopwords/", "Directory to load stopwords files from")

func main() {
	flag.Parse()

	wordSetDataMapper := datamapper.NewStopWordSet(*stopWordsDirectory)
	wordSetRepository := repository.NewStopWordSet(wordSetDataMapper)
	stopWordFilterService := service.NewStopWordFilter(wordSetRepository)
	tweetsGateway := gateway.NewTweets(http.DefaultClient, *tweetsURL)
	tweetsWordsResource := resource.TweetsWords{
		Gateway: tweetsGateway,
		Service: stopWordFilterService,
	}

	router := httprouter.New()
	router.GET("/tweetsWords", tweetsWordsResource.GET)

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
