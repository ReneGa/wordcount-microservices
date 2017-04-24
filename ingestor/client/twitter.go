package client

import (
	"net/url"

	"github.com/chimeracoder/anaconda"
)

// Anaconda is an interface wrapper around the anaconda twitter client
type Anaconda interface {
	SetConsumerKey(key string)
	SetConsumerSecret(keySecret string)
	NewTwitterAPI(token string, tokenSecret string) TwitterAPI
}

type anacondaWrapper struct{}

// NewAnaconda creates a new Anaconda
func NewAnaconda() Anaconda {
	return &anacondaWrapper{}
}

func (a *anacondaWrapper) SetConsumerKey(key string) {
	anaconda.SetConsumerKey(key)
}

func (a *anacondaWrapper) SetConsumerSecret(keySecret string) {
	anaconda.SetConsumerSecret(keySecret)
}

func (a *anacondaWrapper) NewTwitterAPI(token string, tokenSecret string) TwitterAPI {
	twitterAPI := twitterAPI(*anaconda.NewTwitterApi(token, tokenSecret))
	return &twitterAPI
}

// TwitterAPI is an interface wrapper around the anaconda twitter API
type TwitterAPI interface {
	PublicStreamFilter(values url.Values) TwitterStream
}

type twitterAPI anaconda.TwitterApi

func (a *twitterAPI) PublicStreamFilter(values url.Values) TwitterStream {
	return &twitterStream{anaconda.TwitterApi(*a).PublicStreamFilter(values)}
}

type twitterStream struct{ *anaconda.Stream }

func (a *twitterStream) C() chan interface{} {
	return a.Stream.C
}

func (a *twitterStream) Stop() {
	a.Stream.Stop()
}

// TwitterStream is an interface wrapper around the anaconda.Stream struct
type TwitterStream interface {
	C() chan interface{}
	Stop()
}
