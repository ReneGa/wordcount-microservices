package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/datamapper"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/gateway"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/repository"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/resource"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/service"
	"github.com/julienschmidt/httprouter"
)

var address = flag.String("address", "localhost:8081", "Address to listen on")
var tweetsURL = flag.String("tweetsURL", "http://localhost:8080/tweets", "URL of the tweet producer to connect to")
var stopWordsDirectory = flag.String("stopWordsDirectory", "stopwords/", "Directory to load stopwords files from")

func main() {
	flag.Parse()

	wordSetDataMapper := datamapper.NewStopWordSet(*stopWordsDirectory)
	wordSetRepository := repository.NewStopWordSet(wordSetDataMapper)
	stopWordFilterService := service.NewStopWordFilter(wordSetRepository)
	tweetsGateway := gateway.DefaultTweets{
		Client: http.DefaultClient,
		URL:    *tweetsURL,
	}

	tweetsResource := resource.Tweets{
		Gateway: tweetsGateway,
		Service: stopWordFilterService,
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
