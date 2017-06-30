package gateway

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/recorder/domain"
)

// Tweets is a gateway to a tweet producing service
type Tweets interface {
	Tweets(query string) domain.Tweets
}

// HTTPTweets is the gateway to get tweets over http
type HTTPTweets struct {
	Client *http.Client
	URL    string
}

type decodeResult int

const (
	decodeError decodeResult = iota
	decodeStopped
)

func decodeResponse(res *http.Response, data chan domain.Tweet, stop chan bool) decodeResult {
	defer res.Body.Close()
	jd := json.NewDecoder(res.Body)
	for {
		select {
		case <-stop:
			return decodeStopped
		default:
			var tweet domain.Tweet
			err := jd.Decode(&tweet)
			if err != nil {
				return decodeError
			}
			select {
			case data <- tweet:
			case <-stop:
				return decodeStopped
			}
		}
	}
}

// Tweets return a stream of tweets for a given search query
func (t *HTTPTweets) Tweets(query string) domain.Tweets {
	url := fmt.Sprintf("%s?q=%s", t.URL, query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	data := make(chan domain.Tweet)
	stop := make(chan bool)

	tweets := domain.Tweets{
		Data: data,
		Stop: stop,
	}

	go func() {
		defer close(data)
		reconnect := true
		for reconnect {
			select {
			case <-stop:
				return
			default:
			}
			res, err := t.Client.Do(req)
			if err == nil {
				decodeResult := decodeResponse(res, data, stop)
				reconnect = decodeResult == decodeError
			} else {
				log.Println("Error: ", err)
			}
		}
	}()

	return tweets
}
