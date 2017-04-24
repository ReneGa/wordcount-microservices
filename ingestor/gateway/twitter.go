package gateway

import (
	"net/url"
	"time"

	"reflect"

	"github.com/ReneGa/tweetcount-microservices/ingestor/client"
	"github.com/ReneGa/tweetcount-microservices/ingestor/domain"
)

// Twitter is a gateway to the Twitter streaming API
type Twitter interface {
	Tweets(query string) domain.Tweets
}

// NewTwitter creates a new Twitter client
func NewTwitter(
	newTwitterAPI func() client.TwitterAPI,
) Twitter {
	return &twitter{newTwitterAPI}
}

type twitter struct {
	newTwitterAPI func() client.TwitterAPI
}

// Tweets returns a stream of public Tweets for a given search query.
func (a twitter) Tweets(query string) domain.Tweets {
	api := a.newTwitterAPI()
	stream := api.PublicStreamFilter(url.Values{
		"track": {query},
	})
	items := stream.C()
	out := make(chan domain.Tweet)
	stop := make(chan bool)

	tweets := domain.Tweets{
		Data: out,
		Stop: stop,
	}

	go func() {
		defer close(out)
		defer stream.Stop()
		for {
			select {
			case item := <-items:
				itemValue := reflect.ValueOf(item)
				Text := itemValue.FieldByName("Text").String()
				ID := itemValue.FieldByName("IdStr").String()
				Language := itemValue.FieldByName("Lang").String()
				timeString := itemValue.FieldByName("CreatedAt").String()
				Time, _ := time.Parse(time.RubyDate, timeString)
				out <- domain.Tweet{
					Text:     Text,
					ID:       ID,
					Time:     Time,
					Language: Language,
				}
			case <-stop:
				return
			}
		}
	}()

	return tweets
}
