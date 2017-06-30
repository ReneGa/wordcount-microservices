package datamapper

import (
	"encoding/json"
	"os"
	"time"

	"io"

	"github.com/ReneGa/tweetcount-microservices/generic"
	"github.com/ReneGa/tweetcount-microservices/recorder/domain"
)

type TweetBucketReader interface {
	ReplayFrom(time.Time, chan domain.Tweet) error
}

type JSONFileTweetBucketReader struct {
	OS          generic.OS
	FileName    string
	FileMode    os.FileMode
	File        io.ReadCloser
	JSONDecoder *json.Decoder
}

func (r *JSONFileTweetBucketReader) ReplayFrom(startTime time.Time, out chan domain.Tweet) error {
	file, err := r.OS.OpenFile(r.FileName, os.O_RDONLY, r.FileMode)
	if err != nil {
		return err
	}
	defer file.Close()
	jsonDecoder := json.NewDecoder(file)
	for {
		var tweet domain.Tweet
		err := jsonDecoder.Decode(&tweet)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if tweet.Time.After(startTime) {
			out <- tweet
		}
	}
	return nil
}
