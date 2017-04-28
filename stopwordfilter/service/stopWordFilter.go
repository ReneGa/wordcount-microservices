package service

import (
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/domain"
	"github.com/ReneGa/tweetcount-microservices/stopwordfilter/repository"
	"golang.org/x/text/language"
)

// StopWordFilter is a service that filters stopwords from a stream of tweets
type StopWordFilter interface {
	Filter(tweets domain.Tweets) domain.Tweets
}

// NewRepositoryBackedStopWordFilter allocates a new stopword filter that loads stopwords from a repository
func NewRepositoryBackedStopWordFilter(repository repository.StopWordSet) *RepositoryBackedStopWordFilter {
	stopWordFilter := &RepositoryBackedStopWordFilter{repository, makeLanguageMatcher(repository)}
	return stopWordFilter
}

func makeLanguageMatcher(repository repository.StopWordSet) language.Matcher {
	stopWordSetIDs := repository.List()
	languages := make([]language.Tag, len(stopWordSetIDs))
	for i, ID := range stopWordSetIDs {
		languages[i] = language.Make(ID)
	}
	return language.NewMatcher(languages)
}

// RepositoryBackedStopWordFilter is a stopword filter that loads stopwords from a repository
type RepositoryBackedStopWordFilter struct {
	Repository      repository.StopWordSet
	languageMatcher language.Matcher
}

func (s *RepositoryBackedStopWordFilter) filterTweet(tweet domain.Tweet) domain.Tweet {
	languageTag, _, _ := s.languageMatcher.Match(language.Make(tweet.Language))
	stopWords := s.Repository.Get(languageTag.String())
	return domain.FilterStopWords(stopWords, tweet)
}

func (s *RepositoryBackedStopWordFilter) filterTweets(in <-chan domain.Tweet, out chan<- domain.Tweet) {
	for tweet := range in {
		out <- s.filterTweet(tweet)
	}
}

// Filter filters the given stream of tweets using
func (s *RepositoryBackedStopWordFilter) Filter(tweets domain.Tweets) domain.Tweets {
	data := make(chan domain.Tweet)
	stop := make(chan bool)

	go func() {
		tweets.Stop <- <-stop
	}()

	go func() {
		s.filterTweets(tweets.Data, data)
		close(data)
	}()

	return domain.Tweets{
		Data: data,
		Stop: stop,
	}
}
