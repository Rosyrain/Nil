package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nil/logic"
	"strconv"
)

func UploadFileHandler(c *gin.Context) {
	//1.参数判断
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("c.FormFile(\"file\") failed", zap.Error(err))
	}
	//2.业务处理
	uid, _, _ := GetCurrentUser(c)
	id := strconv.Itoa(int(uid))
	url, err := logic.UploadFile(id, file)
	if err != nil {
		zap.L().Error("logic.UploadFile(file) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, gin.H{
		"url": url,
	})
}
