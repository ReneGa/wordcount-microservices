package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ReneGa/tweetcount-microservices/generator/domain"
	"github.com/ReneGa/tweetcount-microservices/generator/resource"
	"github.com/julienschmidt/httprouter"
)

var address = flag.String("address", ":8080", "Address to listen on")
var delay = flag.Duration("delay", 0, "Delay between tweets")
var language = flag.String("language", "en", "Language of generated tweets")

func main() {
	flag.Parse()

	tweetsGenerator := &domain.TweetsGenerator{
		Delay:    *delay,
		Language: *language,
	}

	tweetsResource := resource.Tweets{Generator: tweetsGenerator}

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
