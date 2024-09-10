package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

// CreateBucket: create a bucket
func CreateBucket(bucketName string) error {
	return ossClient.CreateBucket(bucketName)
}

// DeleteBucket: delete a bucket
func DeleteBucket(bucketName string) error {
	return ossClient.DeleteBucket(bucketName)
}

// ListBuckets: list all the buckets
func ListBuckets() ([]oss.BucketProperties, error) {
	result, err := ossClient.ListBuckets()
	if err != nil {
		return nil, err
	}
	return result.Buckets, nil
}
