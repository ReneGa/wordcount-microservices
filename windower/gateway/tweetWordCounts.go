package gateway

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/windower/domain"
)

// TweetWordCounts is a gateway to a TweetWordCount producing service
type TweetWordCounts interface {
	TweetWordCounts(query string) domain.TweetWordCounts
}

// HTTPTweetWordCounts is the gateway to get tweets over http
type HTTPTweetWordCounts struct {
	Client *http.Client
	URL    string
}

type decodeResult int

const (
	decodeError decodeResult = iota
	decodeStopped
)

func decodeResponse(res *http.Response, data chan domain.TweetWordCount, stop chan bool) decodeResult {
	defer res.Body.Close()
	jd := json.NewDecoder(res.Body)
	for {
		select {
		case <-stop:
			return decodeStopped
		default:
			var tweetWordCount domain.TweetWordCount
			err := jd.Decode(&tweetWordCount)
			if err != nil {
				return decodeError
			}
			select {
			case data <- tweetWordCount:
			case <-stop:
				return decodeStopped
			}
		}
	}
}

// TweetWordCounts returns a stream of TweetWordCounts for a given search query
func (t *HTTPTweetWordCounts) TweetWordCounts(query string) domain.TweetWordCounts {
	url := fmt.Sprintf("%s?q=%s", t.URL, query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	data := make(chan domain.TweetWordCount)
	stop := make(chan bool)

	tweets := domain.TweetWordCounts{
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
