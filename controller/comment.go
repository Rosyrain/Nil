package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nil/logic"
	"nil/models"
)

func CommentToPostHandler(c *gin.Context) {
	//1.参数校验
	p := new(models.ParamComment)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("CommentToPost with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	comment := models.Comment{
		PostID:  p.PostID,
		Content: p.Content,
	}

	comment.AuthorID, _, _ = GetCurrentUser(c)

	//2.业务处理
	if err := logic.CommentToPost(&comment); err != nil {
		zap.L().Error("logic.CommentToPost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

func CommentToCommentHandler(c *gin.Context) {
	//1.参数校验
	p := new(models.ParamSubComment)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("CommentToComment with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	comment := models.Comment{
		PostID:  p.PostID,
		Content: p.Content,
	}
	comment.AuthorID, _, _ = GetCurrentUser(c)

	mcomment_id := p.CommentId

	//2.业务处理
	if err := logic.CommentToComment(&comment, mcomment_id); err != nil {
		zap.L().Error("logic.CommentToPost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

func GetUserCommentListHandler(c *gin.Context) {
	//初始化参数
	p := models.ParamCommentList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	//1.参数校验
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("GetUserCommentList with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.业务处理
	data, err := logic.GetUserCommentList(&p)
	if err != nil {
		zap.L().Error("logic.GetUserCommentList(&p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, data)
}

func GetPostMCommentListHandler(c *gin.Context) {
	//初始化参数
	p := models.ParamCommentList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	//1.参数绑定
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("GetPostMCommentList with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.业务处理
	data, err := logic.GetPostMCommentList(&p)
	if err != nil {
		zap.L().Error("logic.GetPostMCommentList(&p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)

}

func GetSubCommentListHandler(c *gin.Context) {
	//参数初始化
	p := models.ParamCommentList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	//1.参数绑定
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("GetPostMCommentList with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.业务处理
	data, err := logic.GetSubCommentList(&p)
	if err != nil {
		zap.L().Error("logic.GetSubCommentList(&p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)

}
