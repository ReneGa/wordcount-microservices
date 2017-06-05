package datamapper

import "path"
import "os"

type Queries struct {
	Directory string
	Buckets   map[string]TweetBuckets
}

func (q *Queries) initBuckets() {
	if q.Buckets == nil {
		q.Buckets = map[string]TweetBuckets{}
	}
}

func (q *Queries) createBuckets(query string) TweetBuckets {
	bucketsDirectory := path.Join(q.Directory, query)
	os.MkdirAll(bucketsDirectory, 0777)
	buckets := TweetBuckets{
		Directory: bucketsDirectory,
	}
	q.Buckets[query] = buckets
	return buckets
}

func (q *Queries) Get(query string) TweetBuckets {
	q.initBuckets()
	if buckets, ok := q.Buckets[query]; ok {
		return buckets
	}
	return q.createBuckets(query)
}
