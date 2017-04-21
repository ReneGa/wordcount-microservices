package gateway

import (
	"net/url"

	"github.com/ReneGa/tweetcount-microservices/ingestor/client"
	"github.com/ReneGa/tweetcount-microservices/ingestor/domain"
	"github.com/chimeracoder/anaconda"
)

// Twitter is a gateway to the Twitter streaming API
type Twitter interface {
	Tweets(query string) domain.Tweets
}

// NewAnacondaTwitter creates a new Twitter client
func NewAnacondaTwitter(
	anaconda client.Anaconda,
	key string,
	keySecret string,
	token string,
	tokenSecret string,
) Twitter {
	anaconda.SetConsumerKey(key)
	anaconda.SetConsumerSecret(keySecret)
	api := anaconda.NewTwitterAPI(token, tokenSecret)
	return &anacondaTwitter{api}
}

type anacondaTwitter struct{ client.AnacondaAPI }

// Tweets returns a stream of public Tweets for a given search query.
func (a anacondaTwitter) Tweets(query string) domain.Tweets {
	stream := a.PublicStreamFilter(url.Values{
		"track": {query},
	})
	anacondaChan := stream.C()
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
			case item := <-anacondaChan:
				if t, ok := item.(anaconda.Tweet); ok {
					tweetTime, _ := t.CreatedAtTime()
					out <- domain.Tweet{
						Text: t.Text,
						ID:   t.IdStr,
						Time: tweetTime,
					}
				}
			case <-stop:
				return
			}
		}
	}()

	return tweets
}
