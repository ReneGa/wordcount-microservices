package service

import (
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/repository"
)

// StopWordFilter is a service that filters stopwords from a stream of tweets
type StopWordFilter interface {
	TweetsWords(tweets domain.Tweets) domain.TweetsWords
}

func NewStopWordFilter(repository repository.WordSet) StopWordFilter {
	return &stopWordFilter{repository}
}

type stopWordFilter struct {
	repository repository.WordSet
}

func (s *stopWordFilter) TweetsWords(tweets domain.Tweets) domain.TweetsWords {
	stopWords := s.repository.Get()
	data := make(chan domain.TweetWords)
	stop := make(chan bool)

	go func() {
		tweets.Stop <- <-stop
	}()

	go func() {
		defer close(data)
		for tweet := range tweets.Data {
			data <- domain.FilterStopWords(stopWords, tweet)
		}
	}()

	return domain.TweetsWords{
		Data: data,
		Stop: stop,
	}
}
