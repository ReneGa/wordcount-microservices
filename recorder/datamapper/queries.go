package datamapper

import (
	"net/url"
	"os"
	"path"
	"sync"
	"time"
)

type Queries struct {
	sync.Mutex
	Directory      string
	BucketDuration time.Duration
	Buckets        map[string]Tweets
}

const bucketsDirectoryMode = 0777

func (q *Queries) initBuckets() {
	if q.Buckets == nil {
		q.Buckets = map[string]Tweets{}
	}
}

func (q *Queries) createBuckets(query string) Tweets {
	bucketsDirectory := path.Join(q.Directory, url.QueryEscape(query))
	os.MkdirAll(bucketsDirectory, bucketsDirectoryMode)

	buckets := &TweetBuckets{
		Directory:      bucketsDirectory,
		BucketDuration: q.BucketDuration,
	}
	q.Buckets[query] = buckets
	return buckets
}

func (q *Queries) Get(query string) Tweets {
	q.Lock()
	defer q.Unlock()

	q.initBuckets()
	if buckets, ok := q.Buckets[query]; ok {
		return buckets
	}
	return q.createBuckets(query)
}
