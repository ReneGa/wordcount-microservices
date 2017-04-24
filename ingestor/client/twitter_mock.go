package client

import (
	"net/url"

	"github.com/stretchr/testify/mock"
)

// MockAnaconda is a mock MockAnaconda Twitter client
type MockAnaconda struct{ mock.Mock }

// SetConsumerKey is a mockable SetConsumerKey
func (a *MockAnaconda) SetConsumerKey(key string) {
	a.Called(key)
}

// SetConsumerSecret is a mockable SetConsumerSecret
func (a *MockAnaconda) SetConsumerSecret(keySecret string) {
	a.Called(keySecret)
}

// NewTwitterAPI is a mockable NewTwitterAPI
func (a *MockAnaconda) NewTwitterAPI(token string, tokenSecret string) TwitterAPI {
	returns := a.Called(token, tokenSecret)
	return returns.Get(0).(TwitterAPI)
}

// MockTwitterAPI is a mock MockTwitterAPI
type MockTwitterAPI struct{ mock.Mock }

// PublicStreamFilter is a mockable PublicStreamFilter
func (a *MockTwitterAPI) PublicStreamFilter(values url.Values) TwitterStream {
	returns := a.Called(values)
	return returns.Get(0).(TwitterStream)
}

// MockTwitterStream is a mock MockTwitterStream
type MockTwitterStream struct{ mock.Mock }

// C is a mockable C
func (a *MockTwitterStream) C() chan interface{} {
	returns := a.Called()
	return returns.Get(0).(chan interface{})
}

// Stop is a mockable Stop
func (a *MockTwitterStream) Stop() {
	a.Called()
}
