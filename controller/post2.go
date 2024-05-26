package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"nil/logic"
	"nil/models"
)

func ExaminePostHandler(c *gin.Context) {
	//1.参数校验
	p := new(models.ParamExamine)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("ExaminePost with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))

		return
	}

	uid, _, _ := GetCurrentUser(c)
	p.UserID = uid
	//2.业务处理
	if err := logic.ExaminePost(p); err != nil {
		zap.L().Error("logic.ExaminePost(p)", zap.Error(err))

		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}
