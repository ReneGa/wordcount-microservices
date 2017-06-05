package service

import (
	"time"

	"github.com/ReneGa/tweetcount-microservices/persister/datamapper"
	"github.com/ReneGa/tweetcount-microservices/persister/domain"
	"github.com/ReneGa/tweetcount-microservices/persister/gateway"
)

type Tweets struct {
	DataMapper *datamapper.Queries
	Gateway    gateway.Tweets
}

func copyTweets(from chan domain.Tweet, to chan domain.Tweet) {
	for tweet := range from {
		to <- tweet
	}
}

func (t *Tweets) streamFreshTweets(history datamapper.Tweets, freshTweets domain.Tweets, out chan domain.Tweet, stop chan bool) {
	defer close(out)
	writeHistory := false
	if history.RegisterWriter() {
		defer history.UnregisterWriter()
		writeHistory = true
	}
	for {
		select {
		case tweet := <-freshTweets.Data:
			if writeHistory {
				history.Append(tweet, time.Now())
			}
			out <- tweet
		case <-stop:
			freshTweets.Stop <- true
			return
		}
	}
}

func (t *Tweets) Tweets(query string, startTime time.Time) domain.Tweets {
	out := make(chan domain.Tweet)
	stop := make(chan bool)
	tweets := domain.Tweets{
		Data: out,
		Stop: stop,
	}
	replayTweets := make(chan domain.Tweet)
	freshTweets := t.Gateway.Tweets(query)
	replayHistory := t.DataMapper.Get(query)
	go replayHistory.ReplayFrom(startTime, replayTweets)
	go func() {
		copyTweets(replayTweets, out)
		t.streamFreshTweets(replayHistory, freshTweets, out, stop)
	}()
	return tweets
}
