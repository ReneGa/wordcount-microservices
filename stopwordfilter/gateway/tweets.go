package gateway

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
)

// Tweets is a gateway to a tweet producing service
type Tweets interface {
	Tweets(query string) domain.Tweets
}

// NewTweets creates a new gateway to a tweets producing service
func NewTweets(client *http.Client, url string) Tweets {
	return &tweets{client, url}
}

type tweets struct {
	client *http.Client
	url    string
}

func streamResponse(res *http.Response, data chan domain.Tweet, stop chan bool) bool {
	defer res.Body.Close()
	var tweet domain.Tweet
	jd := json.NewDecoder(res.Body)
	for {
		select {
		case <-stop:
			return false
		default:
			err := jd.Decode(&tweet)
			if err != nil {
				return true
			}
			data <- tweet
		}
	}
}

func (t *tweets) Tweets(query string) domain.Tweets {
	url := fmt.Sprintf("%s?q=%s", t.url, query)
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
			res, err := t.client.Do(req)
			if err == nil {
				reconnect = streamResponse(res, data, stop)
			} else {
				log.Println("Error: ", err)
			}
		}
	}()

	return tweets
}
