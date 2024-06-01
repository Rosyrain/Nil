package controller

import (
	"fmt"
	"nil/logic"
	"nil/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateHandler  创建帖子
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数的校验

	//c.ShouldBindJSON()  //validator --> binding tag
	p := new(models.Post)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Debug(" c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	//从token中取到当前用户的ID
	userID, _, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	p.Status = 0 //待审核
	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情处理函数
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数（从url获取帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.根据id获取帖子数据（查询数据库）
	data, err := logic.GetPostDetail(pid)
	if err != nil {
		zap.L().Error("logic.GetPostDetail(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//返回响应
	ResponseSuccess(c, data)
}

// GetPostListHanlder  获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	offset, limit := GetPageInfo(c)
	//1.获取数据
	data, err := logic.GetPostList(offset, limit)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//2.返回响应
	ResponseSuccess(c, data)
}

// GetPostList2  获升级版帖子裂口接口
// 根据前端传来的参数动态获取帖子列表
//按照创建时间或分数进行排序

//1.获取参数
//2.取redis查询id列表
//3.根据id取数据库中查询帖子详细信息

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	//Get请求参数(query string)：/api/v1/posts2?page=1&size=10&order=time
	//获取参数
	//c.ShouldBindQuery()		//根据请求的参数类型选择相应的方法取获取参数
	//c.ShouldBindJSON()	//如果请求中携带的是json格式的参数，才能用这个方法获取参数

	//初始化结构体参数
	p := models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	//1.获取参数
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.去redis查询帖子列表
	data, err := logic.GetPostListNew(&p)
	if err != nil {
		zap.L().Error("logic.GetPostListNew(&p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, data)
}

// GetChunkPostListHandler  根据社区去查询帖子列表
//func GetChunkPostListHandler(c *gin.Context) {
//	//初始化结构体参数
//	p := models.ParamChunkPostList{
//		ParamPostList: &models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime,
//		},
//	}
//
//	//1.获取参数
//	if err := c.ShouldBindQuery(&p); err != nil {
//		zap.L().Error("GetChunkPostListHandler with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	//2.取redis查询帖子列表
//	data, err := logic.GetChunkPostList(&p)
//	if err != nil {
//		zap.L().Error("logic.GetChunkPostList() failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//
//	//3.返回响应
//	ResponseSuccess(c, data)
//}

func GetUserPostListHandler(c *gin.Context) {

	//初始化结构体参数
	p := models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	//参数校验
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("GetUserPostList with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//业务处理
	//1.去redis中取出帖子的id列表
	data, err := logic.GetUserPostList(&p)
	if err != nil {
		zap.L().Error("logic.GetUserPostList(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//返回响应
	ResponseSuccess(c, data)
}

func ResubmitPostHandler(c *gin.Context) {
	//1.参数校验
	p := new(models.ParamResubmitPost)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Debug(" c.ShouldBindJSON(&p) error", zap.Any("err", err))
		zap.L().Error("Resubmit post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	if p.PostID != p.Post.ID {
		fmt.Println(p.PostID, p.Post.ID)
		zap.L().Error("p.PostID != p.Post.ID", zap.Any("修改id不相等和提交的id", p.Post))
		ResponseError(c, CodeServerBusy)
		return
	}

	uid, _, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	if uid != p.Post.AuthorID {
		zap.L().Error("no power to resubmit post", zap.Any("用户没有权限", p.Post.AuthorID))
		ResponseError(c, CodeServerBusy)
		return
	}

	//2.业务处理
	if err := logic.ResubmitPost(p); err != nil {
		zap.L().Error("logic.ResubmitPost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)

}

func DeletePostHandler(c *gin.Context) {
	//1.参数校验
	id := c.Param("id")
	pid, _ := strconv.ParseInt(id, 10, 64)

	uid, _, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	//2.业务处理
	if err := logic.DeletePost(pid, uid); err != nil {
		zap.L().Error("logic.DeletePost(pid,cid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}
