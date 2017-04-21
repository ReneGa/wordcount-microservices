package main

import (
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

func main() {

	key := requireEnv("TWITTER_CONSUMER_KEY")
	keySecret := requireEnv("TWITTER_CONSUMER_KEY_SECRET")
	token := requireEnv("TWITTER_ACCESS_TOKEN")
	tokenSecret := requireEnv("TWITTER_ACCESS_TOKEN_SECRET")

	anaconda := client.NewAnaconda()
	twitterGateway := gateway.NewAnacondaTwitter(
		anaconda,
		key,
		keySecret,
		token,
		tokenSecret,
	)

	tweetsResource := resource.Tweets{Gateway: twitterGateway}

	router := httprouter.New()
	router.GET("/tweets", tweetsResource.GET)

	address := "localhost:8080"
	done := make(chan bool)
	go func() {
		err := http.ListenAndServe(address, router)
		if err != nil {
			log.Fatal(err)
		}
		done <- true
	}()
	log.Println("listening...")
	<-done
}
