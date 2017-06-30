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
	DataMapper             datamapper.TweetBucketsPerQuery
	Gateway                gateway.Tweets
	tweetConsumersForQuery map[string][]domain.Tweets
}

const BucketsWriteBufferSize = 16

func copyTweets(from chan domain.Tweet, to chan domain.Tweet) {
	for tweet := range from {
		to <- tweet
	}
}

func copyStoppableTweets(from domain.Tweets, to chan domain.Tweet, stop chan bool) {
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
	consumers[i] = consumers[len(consumers)-1]
	consumers = consumers[:len(consumers)-1]
	t.tweetConsumersForQuery[query] = consumers
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *Tweets) writeToBuckets(query string, tweetsToWriteToBuckets chan domain.Tweet) {

	var bucketWriter datamapper.TweetBucketWriter
	bucketWriterCreator, err := t.DataMapper.BucketWriterCreator(query)
	panicOnErr(err)

	for tweet := range tweetsToWriteToBuckets {

		if bucketWriter == nil {
			bucketWriter = bucketWriterCreator.CreateForTime(tweet.Time)
			panicOnErr(bucketWriter.Open())

		} else if bucketWriterCreator.ShouldCreateNew(tweet.Time) {
			panicOnErr(bucketWriter.Close())
			bucketWriter = bucketWriterCreator.CreateForTime(tweet.Time)
			panicOnErr(bucketWriter.Open())
		}

		panicOnErr(bucketWriter.Append(tweet))
	}
	panicOnErr(bucketWriter.Close())
}

func (t *Tweets) broadcastTweets(query string, tweetsForQuery domain.Tweets, tweetsToWriteToBuckets chan domain.Tweet) {
	defer close(tweetsToWriteToBuckets)
	for tweet := range tweetsForQuery.Data {
		tweetsToWriteToBuckets <- tweet
		t.Lock()
		for i, consumer := range t.tweetConsumersForQuery[query] {
			select {
			case consumer.Data <- tweet:
			case <-consumer.Stop:
				t.removeConsumerForQuery(query, i)
				if len(t.tweetConsumersForQuery[query]) == 0 {
					delete(t.tweetConsumersForQuery, query)
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
		tweetsToWriteToBuckets := make(chan domain.Tweet, BucketsWriteBufferSize)

		go t.writeToBuckets(query, tweetsToWriteToBuckets)
		go t.broadcastTweets(query, tweetsForQuery, tweetsToWriteToBuckets)
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
	go func() {
		panicOnErr(t.DataMapper.ReplayFrom(query, startTime, replayTweets))
	}()
	go func() {
		defer close(out)
		copyTweets(replayTweets, out)
		copyStoppableTweets(freshTweets, out, stop)
	}()
	return tweets
}
