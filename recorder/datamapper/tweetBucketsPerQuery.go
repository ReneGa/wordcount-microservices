package datamapper

import (
	"os"
	"path"
	"time"

	"net/url"

	"github.com/ReneGa/tweetcount-microservices/generic"
	"github.com/ReneGa/tweetcount-microservices/recorder/domain"
)

type TweetBucketsPerQuery interface {
	BucketWriterCreator(query string) (TweetBucketWriterCreator, error)
	ReplayFrom(query string, startTime time.Time, out chan domain.Tweet) error
}

type JSONFileTweetBucketsPerQuery struct {
	IOUtil    generic.IOUtil
	OS        generic.OS
	Directory string
	FileMode  os.FileMode
	Duration  time.Duration
	FileNamer BucketFileNamer
}

func (q *JSONFileTweetBucketsPerQuery) bucketsDirectoryForQuery(query string) string {
	return path.Join(q.Directory, url.QueryEscape(query))
}

func (q *JSONFileTweetBucketsPerQuery) listBucketsForQuery(query string) ([]string, error) {
	fileInfos, err := q.IOUtil.ReadDir(q.bucketsDirectoryForQuery(query))
	if err != nil {
		return nil, err
	}
	fileNames := make([]string, len(fileInfos))
	for i, fileInfo := range fileInfos {
		fileNames[i] = fileInfo.Name()
	}
	return fileNames, nil
}

func (q *JSONFileTweetBucketsPerQuery) ensureDirectoryExists(directory string) error {
	err := q.OS.MkdirAll(directory, os.ModeDir|q.FileMode)
	if err != nil {
		return err
	}
	return nil
}

func (q *JSONFileTweetBucketsPerQuery) BucketWriterCreator(query string) (TweetBucketWriterCreator, error) {
	bucketsDirectory := q.bucketsDirectoryForQuery(query)
	err := q.ensureDirectoryExists(bucketsDirectory)
	if err != nil {
		return nil, err
	}

	return &JSONFileTweetBucketWriterCreator{
		OS:        q.OS,
		Directory: bucketsDirectory,
		Duration:  q.Duration,
		FileMode:  q.FileMode,
		FileNamer: q.FileNamer,
	}, nil
}

func copyTweets(in chan domain.Tweet, out chan domain.Tweet) {
	for tweet := range in {
		out <- tweet
	}
}

func (q *JSONFileTweetBucketsPerQuery) ReplayFrom(query string, startTime time.Time, out chan domain.Tweet) error {
	defer close(out)
	bucketsDirectory := q.bucketsDirectoryForQuery(query)
	err := q.ensureDirectoryExists(bucketsDirectory)
	if err != nil {
		return err
	}
	bucketFileNames, err := q.listBucketsForQuery(query)
	startBucketFileName := q.FileNamer.Name(startTime)
	if err != nil {
		return err
	}
	for _, bucketFileName := range bucketFileNames {
		if bucketFileName >= startBucketFileName {
			bucketReader := &JSONFileTweetBucketReader{
				OS:       q.OS,
				FileName: path.Join(bucketsDirectory, bucketFileName),
				FileMode: q.FileMode,
			}
			bucketReader.ReplayFrom(startTime, out)
		}
	}
	return nil
}
