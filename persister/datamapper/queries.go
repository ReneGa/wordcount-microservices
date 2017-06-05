package datamapper

import (
	"os"
	"path"
	"time"
)

type Queries struct {
	Directory      string
	BucketDuration time.Duration
	Buckets        map[string]Tweets
}

func (q *Queries) initBuckets() {
	if q.Buckets == nil {
		q.Buckets = map[string]Tweets{}
	}
}

func (q *Queries) createBuckets(query string) Tweets {
	bucketsDirectory := path.Join(q.Directory, query)
	os.MkdirAll(bucketsDirectory, 0777)
	buckets := &TweetBuckets{
		Directory:      bucketsDirectory,
		BucketDuration: q.BucketDuration,
	}
	q.Buckets[query] = buckets
	return buckets
}

func (q *Queries) Get(query string) Tweets {
	q.initBuckets()
	if buckets, ok := q.Buckets[query]; ok {
		return buckets
	}
	return q.createBuckets(query)
}
