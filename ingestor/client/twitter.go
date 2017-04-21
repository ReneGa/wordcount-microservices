package client

import (
	"net/url"

	"github.com/chimeracoder/anaconda"
)

// Anaconda is an interface wrapper around the anaconda twitter client
type Anaconda interface {
	SetConsumerKey(key string)
	SetConsumerSecret(keySecret string)
	NewTwitterAPI(token string, tokenSecret string) AnacondaAPI
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

func (a *anacondaWrapper) NewTwitterAPI(token string, tokenSecret string) AnacondaAPI {
	anacondaAPI := anacondaAPI(*anaconda.NewTwitterApi(token, tokenSecret))
	return &anacondaAPI
}

// AnacondaAPI is an interface wrapper around the anaconda twitter API
type AnacondaAPI interface {
	PublicStreamFilter(values url.Values) AnacondaStream
}

type anacondaAPI anaconda.TwitterApi

func (a *anacondaAPI) PublicStreamFilter(values url.Values) AnacondaStream {
	return &anacondaStream{anaconda.TwitterApi(*a).PublicStreamFilter(values)}
}

type anacondaStream struct{ *anaconda.Stream }

func (a *anacondaStream) C() chan interface{} {
	return a.Stream.C
}

func (a *anacondaStream) Stop() {
	a.Stream.Stop()
}

// AnacondaStream is an interface wrapper around the anaocnda.Stream struct
type AnacondaStream interface {
	C() chan interface{}
	Stop()
}
