package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nil/dao/redis"
	"nil/logic"
	"nil/models"
)

func FocusHandler(c *gin.Context) {
	//1.参数校验
	p := new(models.ParamFocusData)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Debug(" c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create focus with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.业务处理
	c_uid, _, _ := GetCurrentUser(c)
	if err := logic.InsertFocus(p, c_uid); err != nil {
		zap.L().Error("logic.InsertFocus(p) failed", zap.Error(err))
		if errors.Is(err, redis.ErrorFocusRepeated) {
			ResponseError(c, CodeRepeated)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

func UserFocusListHandler(c *gin.Context) {
	//1.参数校验
	p := models.ParamFocusList{
		Page: 1,
		Size: 10,
	}

	uid, _, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.UserID = uid

	//2.业务处理
	data, err := logic.GetUserFocusList(&p)
	if err != nil {
		zap.L().Error("logic.GetUserFocusList(uid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}
