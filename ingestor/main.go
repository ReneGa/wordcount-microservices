package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ReneGa/tweetcount-microservices/ingestor/client"
	"github.com/ReneGa/tweetcount-microservices/ingestor/gateway"
	"github.com/ReneGa/tweetcount-microservices/ingestor/resource"
	"github.com/julienschmidt/httprouter"
)

func requireEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Sprintf("Required environment variable '%s' not set.", name))
	}
	return value
}

var address = flag.String("address", "localhost:8080", "Address to listen on")

func main() {
	flag.Parse()

	key := requireEnv("TWITTER_CONSUMER_KEY")
	keySecret := requireEnv("TWITTER_CONSUMER_KEY_SECRET")
	token := requireEnv("TWITTER_ACCESS_TOKEN")
	tokenSecret := requireEnv("TWITTER_ACCESS_TOKEN_SECRET")

	anaconda := client.NewAnaconda()
	anaconda.SetConsumerKey(key)
	anaconda.SetConsumerSecret(keySecret)
	newTwitterAPI := func() client.TwitterAPI {
		return anaconda.NewTwitterAPI(token, tokenSecret)
	}

	twitterGateway := gateway.NewTwitter(newTwitterAPI)

	tweetsResource := resource.Tweets{Gateway: twitterGateway}

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
