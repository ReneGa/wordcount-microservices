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
func (a *MockAnaconda) NewTwitterAPI(token string, tokenSecret string) AnacondaAPI {
	returns := a.Called(token, tokenSecret)
	return returns.Get(0).(AnacondaAPI)
}

// MockAnacondaAPI is a mock MockAnacondaAPI
type MockAnacondaAPI struct{ mock.Mock }

// PublicStreamFilter is a mockable PublicStreamFilter
func (a *MockAnacondaAPI) PublicStreamFilter(values url.Values) AnacondaStream {
	returns := a.Called(values)
	return returns.Get(0).(AnacondaStream)
}

// MockAnacondaStream is a mock MockAnacondaStream
type MockAnacondaStream struct{ mock.Mock }

// C is a mockable C
func (a *MockAnacondaStream) C() chan interface{} {
	returns := a.Called()
	return returns.Get(0).(chan interface{})
}

// Stop is a mockable Stop
func (a *MockAnacondaStream) Stop() {
	a.Called()
}
