package oss

import (
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// UploadFile: upload the file to the OSS
func UploadFile(bucketName, objectKey, localFile string) error {
	// Get the bucket
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return err
	}

	// Upload the file
	return bucket.PutObjectFromFile(objectKey, localFile)
}

// GenerateDownloadURL: generate the download URL for the file in the OSS
func GenerateDownloadURL(bucketName, objectKey string, expiryTime time.Duration) (string, error) {
	// Get the bucket
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	// Generate the download URL
	signedURL, err := bucket.SignURL(objectKey, oss.HTTPGet, int64(expiryTime.Seconds()))
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

// DownloadFile: download the file from the OSS
func DownloadFile(bucketName, objectKey, downloadPath string) error {
	// Get the bucket
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return err
	}

	// Download the file
	return bucket.GetObjectToFile(objectKey, downloadPath)
}

// DeleteFile: delete the file from the OSS
func DeleteFile(bucketName, objectKey string) error {
	// Get the bucket
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return err
	}

	// Delete the file
	return bucket.DeleteObject(objectKey)
}
