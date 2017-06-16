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

// attemptToRegisterWriter attempts to register a history writer
// and returns a boolean indicating whether the registration was
// successful.
func (t *Tweets) attemptToRegisterWriter(query string) bool {
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

func (t *Tweets) copyAndRecordTweets(query string, history datamapper.Tweets, freshTweets domain.Tweets, out chan domain.Tweet, stop chan bool) {
	writeHistory := false
	defer func() {
		if writeHistory {
			t.unregisterWriter(query)
		}
	}()
	for {
		select {
		case tweet := <-freshTweets.Data:
			writeHistory := t.attemptToRegisterWriter(query)
			if writeHistory {
				history.Append(tweet)
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
	freshTweets := t.Gateway.Tweets(query)
	replayTweets := make(chan domain.Tweet)
	history := t.DataMapper.Get(query)
	go history.ReplayFrom(startTime, replayTweets)
	go func() {
		defer close(out)
		copyTweets(replayTweets, out)
		t.copyAndRecordTweets(query, history, freshTweets, out, stop)
	}()
	return tweets
}
