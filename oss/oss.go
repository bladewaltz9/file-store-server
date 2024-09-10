package oss

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/bladewaltz9/file-store-server/config"
)

var ossClient *oss.Client

// init: initialize the oss connection
func initOSS(endpoint, accessKeyID, accessKeySecret string) {
	var err error
	ossClient, err = oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the OSS: %v", err.Error()))
	}
}

func init() {
	initOSS(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
}

// GetOSSClient: get the oss client
func GetOSSClient() *oss.Client {
	return ossClient
}
