package datamapper

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/ReneGa/tweetcount-microservices/persister/domain"
)

type TweetBuckets struct {
	Directory string
}

type bucketID string

const bucketFileMode = 0666

func (t *TweetBuckets) bucketForTweet(tweet domain.Tweet, now time.Time) bucketID {
	return t.bucketForTime(now)
}

func (t *TweetBuckets) bucketForTime(now time.Time) bucketID {
	return bucketID(now.Round(time.Hour).Format(time.RFC3339))
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
		} else {
		}
	}
	return bucketIDs
}

func (t *TweetBuckets) appendToBucket(tweet domain.Tweet, now time.Time, bucket bucketID) error {
	bucketFileName := t.bucketFileName(bucket)
	f, err := os.OpenFile(bucketFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, bucketFileMode)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	je := json.NewEncoder(f)
	return je.Encode(tweet)
}

func (t *TweetBuckets) readStartingFromBucket(startBucket bucketID, bucketTime time.Time, out chan domain.Tweet) error {
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
			if tweet.Time.After(bucketTime) {
				out <- tweet
			}
		}
	}
	return nil
}

func (t *TweetBuckets) Append(tweet domain.Tweet, now time.Time) {
	bucket := t.bucketForTweet(tweet, now)
	t.appendToBucket(tweet, now, bucket)
}

func (t *TweetBuckets) ReplayFrom(bucketTime time.Time, out chan domain.Tweet) error {
	bucket := t.bucketForTime(bucketTime)
	return t.readStartingFromBucket(bucket, bucketTime, out)
}
