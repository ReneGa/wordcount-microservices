package datamapper

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/ReneGa/tweetcount-microservices/recorder/domain"
)

type Tweets interface {
	Append(tweet domain.Tweet, now time.Time)
	ReplayFrom(now time.Time, out chan domain.Tweet) error
}

type TweetBuckets struct {
	Directory      string
	BucketDuration time.Duration
}

type bucketID string

const bucketFileMode = 0666

func (t *TweetBuckets) bucketForTweet(tweet domain.Tweet) bucketID {
	return t.bucketForTime(tweet.Time)
}

func (t *TweetBuckets) bucketForTime(now time.Time) bucketID {
	return bucketID(fmt.Sprintf("%d", now.Round(t.BucketDuration).Unix()))

}
func (t *TweetBuckets) bucketFileName(bucket bucketID) string {
	return path.Join(t.Directory, string(bucket))
}

func (t *TweetBuckets) listBucketsAfter(startBucket bucketID) []bucketID {
	files, _ := ioutil.ReadDir(t.Directory)
	bucketIDs := make([]bucketID, 0, len(files))
	for _, file := range files {
		bucket := bucketID(file.Name())
		if bucket >= startBucket {
			bucketIDs = append(bucketIDs, bucket)
		}
	}
	return bucketIDs
}

func (t *TweetBuckets) appendToBucket(bucket bucketID, tweet domain.Tweet) error {
	bucketFileName := t.bucketFileName(bucket)
	f, err := os.OpenFile(bucketFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, bucketFileMode)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	je := json.NewEncoder(f)
	return je.Encode(tweet)
}

func (t *TweetBuckets) readStartingFromBucket(startBucket bucketID, startTime time.Time, out chan domain.Tweet) error {
	defer close(out)
	buckets := t.listBucketsAfter(startBucket)
	for _, bucket := range buckets {
		f, err := os.OpenFile(t.bucketFileName(bucket), os.O_RDONLY, bucketFileMode)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		jd := json.NewDecoder(f)
		for {
			var tweet domain.Tweet
			err := jd.Decode(&tweet)
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			if tweet.Time.After(startTime) {
				out <- tweet
			}
		}
	}
	return nil
}

func (t *TweetBuckets) Append(tweet domain.Tweet, now time.Time) {
	bucket := t.bucketForTweet(tweet)
	t.appendToBucket(bucket, tweet)
}

func (t *TweetBuckets) ReplayFrom(startTime time.Time, out chan domain.Tweet) error {
	bucket := t.bucketForTime(startTime)
	return t.readStartingFromBucket(bucket, startTime, out)
}
