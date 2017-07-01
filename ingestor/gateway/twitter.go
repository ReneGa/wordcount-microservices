package gateway

import (
	"net/url"
	"time"

	"github.com/ReneGa/tweetcount-microservices/ingestor/client"
	"github.com/ReneGa/tweetcount-microservices/ingestor/domain"
	"github.com/chimeracoder/anaconda"
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

func convertTweet(anacondaTweet anaconda.Tweet) domain.Tweet {
	Text := anacondaTweet.Text
	ID := anacondaTweet.IdStr
	Language := anacondaTweet.Lang
	timeString := anacondaTweet.CreatedAt
	Time, _ := time.Parse(time.RubyDate, timeString)
	return domain.Tweet{
		Text:     Text,
		ID:       ID,
		Time:     Time,
		Language: Language,
	}
}

func processStream(stream client.TwitterStream, out chan domain.Tweet, stop chan bool) {
	items := stream.C()
	defer close(out)
	defer stream.Stop()
	for {
		select {
		case item := <-items:
			if item != nil {
				switch typedItem := item.(type) {
				case anaconda.Tweet:
					tweet := convertTweet(typedItem)
					select {
					case out <- tweet:
					case <-stop:
						return
					}
				default:
				}
			}
		case <-stop:
			return
		}
	}
}

// Tweets returns a stream of public Tweets for a given search query.
func (a twitter) Tweets(query string) domain.Tweets {
	api := a.newTwitterAPI()
	stream := api.PublicStreamFilter(url.Values{
		"track": {query},
	})

	out := make(chan domain.Tweet)
	stop := make(chan bool)

	tweets := domain.Tweets{
		Data: out,
		Stop: stop,
	}

	go processStream(stream, out, stop)

	return tweets
}
