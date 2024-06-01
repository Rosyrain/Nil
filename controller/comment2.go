package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nil/dao/mysql"
	"nil/logic"
	"strconv"
)

func SuperUserCommentDeleteHandler(c *gin.Context) {
	//1. 参数校验
	id := c.Param("id")
	cid, _ := strconv.ParseInt(id, 10, 64)

	uid, _, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	//2.业务处理
	if err := logic.SuperuserCommentDelete(cid, uid); err != nil {

		zap.L().Error("logic.CommentDelete failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorSuperuserNotExist) {
			ResponseError(c, CodeNeedLogin)
			return
		}
		if errors.Is(err, mysql.ErrorNotPower) {
			ResponseErrorWithMsg(c, CodeServerBusy, "没有权限")
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}
