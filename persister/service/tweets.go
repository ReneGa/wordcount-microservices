package service

import (
	"sync"
	"time"

	"github.com/ReneGa/tweetcount-microservices/persister/datamapper"
	"github.com/ReneGa/tweetcount-microservices/persister/domain"
	"github.com/ReneGa/tweetcount-microservices/persister/gateway"
)

type Tweets struct {
	sync.Mutex
	DataMapper              *datamapper.Queries
	historyWriterRegistered map[string]bool
	Gateway                 gateway.Tweets
}

func copyTweets(from chan domain.Tweet, to chan domain.Tweet) {
	for tweet := range from {
		to <- tweet
	}
}

func (t *Tweets) registerHistoryWriter(query string) bool {
	t.Lock()
	if t.historyWriterRegistered == nil {
		t.historyWriterRegistered = map[string]bool{}
	}
	alreadyRegistered := t.historyWriterRegistered[query]
	t.historyWriterRegistered[query] = true
	t.Unlock()
	return !alreadyRegistered
}

func (t *Tweets) unRegisterHistoryWriter(query string) {
	t.Lock()
	t.historyWriterRegistered[query] = true
	t.Unlock()
}

func (t *Tweets) streamFreshTweets(history datamapper.Tweets, freshTweets domain.Tweets, out chan domain.Tweet, stop chan bool) {
	defer close(out)
	for {
		select {
		case tweet := <-freshTweets.Data:
			if history != nil {
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
		if !t.registerHistoryWriter(query) {
			replayHistory = nil
		} else {
			defer t.unRegisterHistoryWriter(query)
		}
		t.streamFreshTweets(replayHistory, freshTweets, out, stop)
	}()
	return tweets
}
