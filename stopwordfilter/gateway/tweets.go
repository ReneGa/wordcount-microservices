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

// DefaultTweets is the gateway to get tweets
type DefaultTweets struct {
	Client *http.Client
	URL    string
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

// Tweets return a stream of tweets for a given search query
func (t *TweetsImpl) Tweets(query string) domain.Tweets {
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
				reconnect = streamResponse(res, data, stop)
			} else {
				log.Println("Error: ", err)
			}
		}
	}()

	return tweets
}
