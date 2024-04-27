package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	setting "nil/settings"
)

var (
	accessKeyID     string
	accessKeySecret string
	endpoint        string
	bucketName      string
	client          *oss.Client
	bucket          *oss.Bucket
)

func Init(config *setting.OssConfig) (err error) {
	// 配置访问凭证和存储桶信息
	accessKeyID = config.AccessKeyID
	accessKeySecret = config.AccessKeySecret
	endpoint = config.Endpoint
	bucketName = config.BucketName

	// 创建 OSS 客户端实例
	client, err = oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return
	}

	// 获取存储桶对象
	bucket, err = client.Bucket(bucketName)
	if err != nil {
		return
	}
	return
}
