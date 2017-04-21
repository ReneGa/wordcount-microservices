package gateways

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

// AnacondaTwitter is a Twitter client implemented using the Anaconda library
func AnacondaTwitter(
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
	stop := make(chan struct{})

	tweets := domain.Tweets{
		Tweets: out,
		Stop:   stop,
	}

	go func() {
		defer close(out)
		for {
			select {
			case item := <-anacondaChan:
				if t, ok := item.(anaconda.Tweet); ok {
					tweetTime, _ := t.CreatedAtTime()
					out <- domain.Tweet{
						Text: t.Text,
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
