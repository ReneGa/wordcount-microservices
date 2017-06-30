package datamapper

import (
	"encoding/json"
	"os"

	"io"

	"github.com/ReneGa/tweetcount-microservices/generic"
	"github.com/ReneGa/tweetcount-microservices/recorder/domain"
)

type TweetBucketWriter interface {
	Open() error
	Close() error
	Append(domain.Tweet) error
}

type JSONFileTweetBucketWriter struct {
	OS          generic.OS
	FileName    string
	FileMode    os.FileMode
	File        io.WriteCloser
	JSONEncoder *json.Encoder
}

func (w *JSONFileTweetBucketWriter) Open() error {
	file, err := w.OS.OpenFile(w.FileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, w.FileMode)
	w.File = file
	w.JSONEncoder = json.NewEncoder(w.File)
	return err
}

func (w *JSONFileTweetBucketWriter) Close() error {
	return w.File.Close()
}

func (w *JSONFileTweetBucketWriter) Append(tweet domain.Tweet) error {
	return w.JSONEncoder.Encode(tweet)
}
