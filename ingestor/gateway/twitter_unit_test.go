package gateway_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/ReneGa/tweetcount-microservices/ingestor/client"
	"github.com/ReneGa/tweetcount-microservices/ingestor/domain"
	"github.com/ReneGa/tweetcount-microservices/ingestor/gateway"
	"github.com/chimeracoder/anaconda"
	"github.com/stretchr/testify/assert"
)

func TestTwitterGatewayShould(t *testing.T) {
	t.Run("return tweets produced by the API", func(t *testing.T) {
		// Given
		mockTweets := make(chan interface{}, 1)
		anacondaTweet1 := anaconda.Tweet{
			IdStr:     "123456",
			Text:      "dog",
			CreatedAt: time.Now().Format(time.RubyDate),
		}
		tweetTime, _ := time.Parse(time.RubyDate, anacondaTweet1.CreatedAt)
		mockTweets <- anacondaTweet1

		mockTwitterStream := &client.MockTwitterStream{}
		mockTwitterStream.On("C").Return(mockTweets)

		mockTwitterAPI := &client.MockTwitterAPI{}
		mockTwitterAPI.
			On("PublicStreamFilter", url.Values{"track": []string{"dog"}}).
			Return(mockTwitterStream)

		twitter := gateway.NewTwitter(func() client.TwitterAPI {
			return mockTwitterAPI
		})

		// When
		dogTweets := twitter.Tweets("dog")
		actualTweet := <-dogTweets.Data

		// Then
		mockTwitterAPI.AssertExpectations(t)
		assert.Equal(t, domain.Tweet{
			ID:   "123456",
			Text: "dog",
			Time: tweetTime,
		}, actualTweet)

	})
}
