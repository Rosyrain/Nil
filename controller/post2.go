package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"nil/dao/mysql"
	"nil/logic"
	"nil/models"
	"strconv"
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

	if p.Status == p.Direction {
		zap.L().Error("ExaminePost failed", zap.Error(errors.New("禁止重复操作")))
		ResponseError(c, CodeRepeated)
		return
	}

	uid, _, _ := GetCurrentUser(c)

	p.SuperuserID = uid
	//2.业务处理
	if err := logic.ExaminePost(p); err != nil {
		zap.L().Error("logic.ExaminePost(p)", zap.Error(err))
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

func SuperUserGetPostListHandler(c *gin.Context) {
	//1.初始化参数
	p := models.ParamSearch{
		Page:    1,
		Size:    10,
		ChunkID: 1,
		Status:  0,
		Order:   "time",
	}

	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("SuperUserGetPostListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	uid, _, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.SuperuserID = uid

	//2.业务处理
	data, err := logic.SuperuserGetPostList(&p)
	if err != nil {
		zap.L().Error("logic.SuperuserGetPostList(&p) failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorNotPower) {
			ResponseErrorWithMsg(c, CodeServerBusy, "没有权限")
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, data)
}

func SuperuserDeletePostHandler(c *gin.Context) {
	//1.参数校验
	id := c.Param("id")
	pid, _ := strconv.ParseInt(id, 10, 64)

	uid, _, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	//2.业务处理
	if err := logic.SuperuserPostDelete(pid, uid); err != nil {
		zap.L().Error("logic.SuperuserPostDelete(pid,uid) failed", zap.Error(err))
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
