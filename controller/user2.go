package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"nil/dao/mysql"
	"nil/logic"
	"nil/models"
)

func SuperUserLoginHandler(c *gin.Context) {
	//1.参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SuperUserLogin with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))

		return
	}

	//2.业务处理
	user, err := logic.SuperUserLogin(p)
	if err != nil {
		zap.L().Error("logic.SuperUserLogin failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   user.UserID,
		"user_name": user.Username,
		"token":     user.Token,
	})
}
