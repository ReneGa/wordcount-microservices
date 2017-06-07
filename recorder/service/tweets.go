package service

import (
	"sync"
	"time"

	"github.com/ReneGa/tweetcount-microservices/recorder/datamapper"
	"github.com/ReneGa/tweetcount-microservices/recorder/domain"
	"github.com/ReneGa/tweetcount-microservices/recorder/gateway"
)

type Tweets struct {
	sync.Mutex
	DataMapper       *datamapper.Queries
	Gateway          gateway.Tweets
	writerRegistered map[string]bool
}

func copyTweets(from chan domain.Tweet, to chan domain.Tweet) {
	for tweet := range from {
		to <- tweet
	}
}

func (t *Tweets) registerWriter(query string) bool {
	t.Lock()
	if t.writerRegistered == nil {
		t.writerRegistered = map[string]bool{}
	}
	alreadyRegistered := t.writerRegistered[query]
	t.writerRegistered[query] = true
	t.Unlock()
	return !alreadyRegistered
}

func (t *Tweets) unregisterWriter(query string) {
	t.Lock()
	t.writerRegistered[query] = false
	t.Unlock()
}

func (t *Tweets) streamFreshTweets(freshTweets domain.Tweets, out chan domain.Tweet, stop chan bool) {
	for {
		select {
		case tweet := <-freshTweets.Data:
			out <- tweet
		case <-stop:
			freshTweets.Stop <- true
			return
		}
	}
}

func (t *Tweets) streamAndRecordFreshTweets(history datamapper.Tweets, freshTweets domain.Tweets, out chan domain.Tweet, stop chan bool) {
	for {
		select {
		case tweet := <-freshTweets.Data:
			history.Append(tweet)
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
	freshTweets := t.Gateway.Tweets(query)
	replayTweets := make(chan domain.Tweet)
	history := t.DataMapper.Get(query)
	go history.ReplayFrom(startTime, replayTweets)
	go func() {
		defer close(out)
		copyTweets(replayTweets, out)
		if !t.registerWriter(query) {
			t.streamFreshTweets(freshTweets, out, stop)
		} else {
			t.streamAndRecordFreshTweets(history, freshTweets, out, stop)
			t.unregisterWriter(query)
		}
	}()
	return tweets
}
