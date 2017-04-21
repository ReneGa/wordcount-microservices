package gateway

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
)

// Tweets is a gateway to a tweet producing service
type Tweets interface {
	Tweets(query string) domain.Tweets
}

type tweets struct {
	client http.Client
	url    string
}

func streamResponse(res *http.Response, data chan domain.Tweet, stop chan bool) bool {
	var tweet domain.Tweet
	jd := json.NewDecoder(res.Body)
	for {
		jd.Decode(&tweet)
		data <- tweet
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
		reconnect := true
		for reconnect {
			res, err := t.client.Do(req)
			if err == nil {
				reconnect = streamResponse(res, data, stop)
			}
		}
	}()
	return tweets
}
