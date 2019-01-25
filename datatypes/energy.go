package datatypes

type StoneIDsWithBucketsWithDataPoints struct {
	StoneID               string
	BucketsWithDataPoints []BucketWithDataPoints
}

type BucketWithDataPoints struct {
	Bucket           int64
	EnergyDataPoints []Data
}