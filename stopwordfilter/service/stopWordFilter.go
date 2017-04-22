package service

import (
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/repository"
	"golang.org/x/text/language"
)

// StopWordFilter is a service that filters stopwords from a stream of tweets
type StopWordFilter interface {
	TweetsWords(tweets domain.Tweets) domain.TweetsWords
}

func NewStopWordFilter(repository repository.WordSet, fallbackLanguage string) StopWordFilter {
	return &stopWordFilter{repository, fallbackLanguage}
}

type stopWordFilter struct {
	repository       repository.WordSet
	fallbackLanguage string
}

func (s *stopWordFilter) TweetsWords(tweets domain.Tweets) domain.TweetsWords {
	stopWordSetIDs := append([]string{s.fallbackLanguage}, s.repository.List()...)
	languages := make([]language.Tag, len(stopWordSetIDs))
	for i, ID := range stopWordSetIDs {
		languages[i] = language.Make(ID)
	}
	matcher := language.NewMatcher(languages)
	data := make(chan domain.TweetWords)
	stop := make(chan bool)

	go func() {
		tweets.Stop <- <-stop
	}()

	go func() {
		defer close(data)
		for tweet := range tweets.Data {
			_, i, _ := matcher.Match(language.Make(tweet.Language))
			stopWords := s.repository.Get(stopWordSetIDs[i])
			data <- domain.FilterStopWords(stopWords, tweet)
		}
	}()

	return domain.TweetsWords{
		Data: data,
		Stop: stop,
	}
}
