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
	DataMapper             *datamapper.Queries
	Gateway                gateway.Tweets
	tweetConsumersForQuery map[string][]domain.Tweets
}

func copyTweets(from chan domain.Tweet, to chan domain.Tweet) {
	for tweet := range from {
		to <- tweet
	}
}

func copyTweetsWithStop(from domain.Tweets, to chan domain.Tweet, stop chan bool) {
	for {
		select {
		case tweet := <-from.Data:
			select {
			case to <- tweet:
			case <-stop:
				from.Stop <- true
				return

			}
		case <-stop:
			from.Stop <- true
			return
		}
	}
}

func (t *Tweets) removeConsumerForQuery(query string, i int) {
	consumers := t.tweetConsumersForQuery[query]
	t.tweetConsumersForQuery[query] = append(
		consumers[0:i],
		consumers[i:len(consumers)]...,
	)
}

func (t *Tweets) broadcastTweets(query string, tweetsForQuery domain.Tweets) {
	dataMapper := t.DataMapper.Get(query)
	for tweet := range tweetsForQuery.Data {
		dataMapper.Append(tweet)
		t.Lock()
		for i, consumer := range t.tweetConsumersForQuery[query] {
			select {
			case consumer.Data <- tweet:
			case <-consumer.Stop:
				t.removeConsumerForQuery(query, i)
				if len(t.tweetConsumersForQuery[query]) == 0 {
					t.Unlock()
					tweetsForQuery.Stop <- true
					return
				}
			}
		}
		t.Unlock()
	}
}

func (t *Tweets) freshTweets(query string) domain.Tweets {
	t.Lock()
	defer t.Unlock()
	if t.tweetConsumersForQuery == nil {
		t.tweetConsumersForQuery = map[string][]domain.Tweets{}
	}
	if _, ok := t.tweetConsumersForQuery[query]; !ok {
		t.tweetConsumersForQuery[query] = []domain.Tweets{}
		tweetsForQuery := t.Gateway.Tweets(query)
		go t.broadcastTweets(query, tweetsForQuery)
	}

	data := make(chan domain.Tweet)
	stop := make(chan bool)
	tweets := domain.Tweets{
		Data: data,
		Stop: stop,
	}
	t.tweetConsumersForQuery[query] = append(t.tweetConsumersForQuery[query], tweets)

	return tweets
}

func (t *Tweets) Tweets(query string, startTime time.Time) domain.Tweets {
	out := make(chan domain.Tweet)
	stop := make(chan bool)
	tweets := domain.Tweets{
		Data: out,
		Stop: stop,
	}
	freshTweets := t.freshTweets(query)
	replayTweets := make(chan domain.Tweet)
	history := t.DataMapper.Get(query)
	go history.ReplayFrom(startTime, replayTweets)
	go func() {
		defer close(out)
		copyTweets(replayTweets, out)
		copyTweetsWithStop(freshTweets, out, stop)
	}()
	return tweets
}
