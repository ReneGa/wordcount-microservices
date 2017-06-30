package datamapper

import "time"

type BucketFileNamer interface {
	Name(time.Time) string
}

type RFC3339BucketFileNamer struct{}

func (_ RFC3339BucketFileNamer) Name(bucketTime time.Time) string {
	return bucketTime.In(time.UTC).Format(time.RFC3339)
}
