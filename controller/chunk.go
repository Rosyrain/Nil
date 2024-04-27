package controller

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"nil/dao/mysql"
	"nil/logic"
	"nil/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/// ---- 跟社区相关 ----

// CreateChunkHandler	创建板块信息
func CreateChunkHandler(c *gin.Context) {
	//1.参数判断
	p := new(models.ParamChunk)

	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("CreateChunk with invalid param", zap.Error(err))
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
	if err := logic.CreateChunk(p); err != nil {
		zap.L().Error("logic.CreateChunk(p) failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorChunkExist) {
			ResponseError(c, CodeChunkExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

// ChunkHandler  返回所有社区的id和name
func ChunkHandler(c *gin.Context) {
	//希望查询到所有的社区（chunk_id,chunk_name）以列表的形式返回
	data, err := logic.GetChunkList()

	if err != nil {
		zap.L().Error("logic.GetChunkList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端报错暴露给外面
	}

	ResponseSuccess(c, data)
}

// ChunkDetailHandler  社区分类详情
func ChunkDetailHandler(c *gin.Context) {
	//1.获取社区id
	ChunkID := c.Param("id") //获取参数
	id, err := strconv.ParseInt(ChunkID, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	//希望查询到对应的社区（chunk_id,chunk_name）以列表的形式返回
	data, err := logic.GetChunkDetail(id)

	if err != nil {
		zap.L().Error("logic.GetChunkDetailList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端报错暴露给外面
		return
	}

	ResponseSuccess(c, data)
}
