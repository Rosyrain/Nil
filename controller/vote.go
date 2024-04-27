package controller

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"nil/dao/redis"
	"nil/logic"
	"nil/models"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

//// VoteData  投票数据
//type VoteData struct {
//	//userID 从请求中获取当前的用户
//	PostID    int64 `json:"post_id,string"`   //帖子ID
//	Direction int   `json:"direction,string"` //赞成票1，反对票-1
//}

func PostVoteHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamPostVoteData)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("PostVoteHand failed", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) //翻译并去除错误提示中的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	//校验当前请求的用户ID
	userID, _, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//具体投票函数
	if err := logic.PostForVote(userID, p); err != nil {
		zap.L().Error("logic.PostForVote failed", zap.Error(err))
		if errors.Is(err, redis.ErrVoteRepeated) {
			ResponseError(c, CodeRepeated)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

func MCommentVoteHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamCommentVoteData)
	if err := c.ShouldBindJSON(&p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) //翻译并去除错误提示中的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	//校验当前请求的用户ID
	userID, _, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//具体投票函数
	if err := logic.MCommentForVote(userID, p); err != nil {
		zap.L().Error("logic.PostForVote failed", zap.Error(err))
		if errors.Is(err, redis.ErrVoteRepeated) {
			ResponseError(c, CodeRepeated)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

func SubCommentVoteHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamCommentVoteData)
	if err := c.ShouldBindJSON(&p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) //翻译并去除错误提示中的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	//校验当前请求的用户ID
	userID, _, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//具体投票函数
	if err := logic.SubCommentForVote(userID, p); err != nil {
		zap.L().Error("logic.PostForVote failed", zap.Error(err))
		if errors.Is(err, redis.ErrVoteRepeated) {
			ResponseError(c, CodeRepeated)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}
