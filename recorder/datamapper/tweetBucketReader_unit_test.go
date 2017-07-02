package datamapper_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"encoding/json"

	"github.com/ReneGa/tweetcount-microservices/generic"
	"github.com/ReneGa/tweetcount-microservices/recorder/datamapper"
	"github.com/ReneGa/tweetcount-microservices/recorder/domain"
	"github.com/stretchr/testify/assert"
)

type CloseableBuffer struct {
	buf *bytes.Buffer
}

func (c CloseableBuffer) Write(p []byte) (n int, err error) {
	return c.buf.Write(p)
}

func (c CloseableBuffer) Read(p []byte) (n int, err error) {
	return c.buf.Read(p)
}

func (c CloseableBuffer) Close() error { return nil }

func TestJSONFileTweetBucketReaderShould(t *testing.T) {
	t.Run("Replay a single tweet", func(t *testing.T) {
		// Given
		tweetTime := time.Now()
		tweet := domain.Tweet{
			ID:   "12345",
			Text: "test text",
			Time: tweetTime,
		}
		tweetBuffer := bytes.NewBuffer([]byte{})
		closeableTweetBuffer := CloseableBuffer{tweetBuffer}

		jsonEncoder := json.NewEncoder(closeableTweetBuffer)
		jsonEncoder.Encode(tweet)

		mockOS := &generic.MockOS{}
		fileName := "bucketFileName"
		fileMode := os.FileMode(0777)
		reader := &datamapper.JSONFileTweetBucketReader{
			OS:       mockOS,
			FileName: fileName,
			FileMode: fileMode,
		}
		out := make(chan domain.Tweet, 1)

		mockOS.
			On("OpenFile", fileName, os.O_RDONLY, fileMode).
			Return(closeableTweetBuffer, nil).
			Once()

		// When
		reader.ReplayFrom(tweetTime, out)

		// Then
		outTweet := <-out

		assert.Equal(t, tweet, outTweet)

	})
}
