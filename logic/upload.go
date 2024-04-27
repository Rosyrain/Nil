package logic

import (
	"mime/multipart"
	"nil/pkg/oss"
)

func UploadFile(uid string, file *multipart.FileHeader) (string, error) {
	return oss.UploadFile(uid, file)
}
