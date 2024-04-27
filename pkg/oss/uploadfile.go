package oss

import (
	"mime/multipart"
	"strconv"
	"time"
)

func UploadFile(uid string, file *multipart.FileHeader) (url string, err error) {
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	// 生成文件在 OSS 中的存储路径
	objectKey := uid + "/" + strconv.Itoa(int(time.Now().Unix())) + "_" + file.Filename

	// 将文件上传到 OSS 存储桶
	err = bucket.PutObject(objectKey, src)
	if err != nil {
		return
	}
	// 获取上传后的文件 URL
	url = "https://" + bucketName + "." + endpoint + "/" + objectKey
	return
}
