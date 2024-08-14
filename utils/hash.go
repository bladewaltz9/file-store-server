package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// CalculateSHA1: calculate the SHA1 hash of the file
func CalculateSHA1(file *os.File) (string, error) {
	// move the file pointer to the beginning of the file
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CalculateSHA256: calculate the SHA256 hash of the file
func CalculateSHA256(file *os.File) (string, error) {
	// move the file pointer to the beginning of the file
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CalulateMD5: calculate the MD5 hash of the file
func CalulateMD5(file *os.File) (string, error) {
	// move the file pointer to the beginning of the file
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
