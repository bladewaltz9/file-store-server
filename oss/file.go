package oss

// UploadFile: upload the file to the OSS
func UploadFile(bucketName, objectName, localFile string) error {
	// Get the bucket
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return err
	}

	// Upload the file
	return bucket.PutObjectFromFile(objectName, localFile)
}

// DownloadFile: download the file from the OSS
func DownloadFile(bucketName, objectName, downloadPath string) error {
	// Get the bucket
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return err
	}

	// Download the file
	return bucket.GetObjectToFile(objectName, downloadPath)
}

// DeleteFile: delete the file from the OSS
func DeleteFile(bucketName, objectName string) error {
	// Get the bucket
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return err
	}

	// Delete the file
	return bucket.DeleteObject(objectName)
}
