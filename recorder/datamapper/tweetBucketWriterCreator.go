package datamapper

import (
	"os"
	"path"
	"time"

	"github.com/ReneGa/tweetcount-microservices/generic"
)

type TweetBucketWriterCreator interface {
	CreateForTime(time.Time) TweetBucketWriter
	ShouldCreateNew(time.Time) bool
}

type JSONFileTweetBucketWriterCreator struct {
	OS                     generic.OS
	Directory              string
	FileMode               os.FileMode
	Duration               time.Duration
	FileNamer              BucketFileNamer
	CurrentBucketStartTime time.Time
}

func (w JSONFileTweetBucketWriterCreator) CreateForTime(tweetTime time.Time) TweetBucketWriter {
	bucketStartTime := tweetTime.Truncate(w.Duration)
	bucketFileName := path.Join(w.Directory, w.FileNamer.Name(bucketStartTime))
	bucket := &JSONFileTweetBucketWriter{
		OS:       w.OS,
		FileMode: w.FileMode,
		FileName: bucketFileName,
	}
	w.CurrentBucketStartTime = bucketStartTime
	return bucket
}

func (w JSONFileTweetBucketWriterCreator) ShouldCreateNew(tweetTime time.Time) bool {
	return w.CurrentBucketStartTime.Add(w.Duration).Before(tweetTime)
}
