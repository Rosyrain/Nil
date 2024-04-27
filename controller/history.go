package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nil/logic"
	"nil/models"
)

func HistoryHandler(c *gin.Context) {
	//1.参数判断
	p := new(models.ParamHistory)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Debug(" c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create history with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	uid, _, _ := GetCurrentUser(c)
	//2.业务处理
	if err := logic.InsertHistory(p.PostID, uid); err != nil {
		zap.L().Error("logic.InsertHistory(post_id) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

func GetUserHistoryListHandler(c *gin.Context) {
	//1.参数判断
	p := models.ParamHistoryList{
		Page: 1,
		Size: 10,
	}
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Debug(" c.ShouldBindQuery(&p) error", zap.Any("err", err))
		zap.L().Error("GetUserHistoryList with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	uid, _, _ := GetCurrentUser(c)
	p.UserID = uid

	//2.业务处理
	data, err := logic.GetUserHistoryList(&p)
	if err != nil {
		zap.L().Error("logic.GetUserHistoryList(&p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, data)
}
