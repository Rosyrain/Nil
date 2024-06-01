package router

import (
	"net/http"
	"nil/controller"
	//_ "nil/docs" // 千万不要忘了导入把你上一步生成的docs
	"nil/logger"
	"nil/middlewares"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/pprof"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//r.LoadHTMLFiles("templates/index.html")
	//r.Static("/static", "./static")
	//r.GET("/", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "index.html", nil)
	//})

	//r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	//注册业务路由
	v1 := r.Group("/api/v1")
	{
		//注册
		v1.POST("/signup", controller.SignUpHandler)

		//发送激活邮件
		v1.POST("/captcha", controller.CaptchaHandler)

		//进行激活
		//v1.GET("/activate_link", controller.VerifyActivateHandler)

		//登录
		v1.POST("/login", controller.LoginHandler)

		//获取用户信息
		v1.GET("/user/:username", controller.UserInfoHandler)

		//获取所有板块信息(简略)(id,name)
		v1.GET("/chunks", controller.ChunkHandler)

		//获取板块具体信息
		v1.GET("/chunk/:id", controller.ChunkDetailHandler)

		//获取帖子具体信息
		v1.GET("/post/:id", controller.GetPostDetailHandler)

		//根据时间或分数获取帖子列表
		v1.GET("/posts", controller.GetPostListHandler2)

		//获取用户发布的帖子列表
		v1.GET("/user/posts", controller.GetUserPostListHandler)

		//获取用户发布的评论列表
		v1.GET("/user/comments", controller.GetUserCommentListHandler)

		//获取帖子主评论列表
		v1.GET("/comments", controller.GetPostMCommentListHandler)

		//获取子评论列表
		v1.GET("/subcomments", controller.GetSubCommentListHandler)

		////
		//	获取帖子不需要登录，下面使用r.GET 而不是v1.GET
		//不分时间或者分数获取帖子列表
		v1.GET("/api/v1/posts2", controller.GetPostListHandler)

		//根据社区获取帖子列表(默认时间，可以改为分数)
		//v1.GET("/posts3", controller.GetChunkPostListHandler)

	}

	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		//更改用户信息
		v1.POST("/user/update_info", controller.UpdateUserInfoHandler)

		//发布帖子
		v1.POST("/post", controller.CreatePostHandler)

		//发布主评论
		v1.POST("/comment_to_post", controller.CommentToPostHandler)

		//发布次级评论
		v1.POST("/comment_to_comment", controller.CommentToCommentHandler)

		//给post投票
		v1.POST("/post/vote", controller.PostVoteHandler)

		//给comment投票
		v1.POST("/comment/vote", controller.MCommentVoteHandler)

		//给subcomment投票
		v1.POST("/subcomment/vote", controller.SubCommentVoteHandler)

		//关注用户
		v1.POST("/focus", controller.FocusHandler)

		//获取用户关注列表
		v1.GET("/user/focus", controller.UserFocusListHandler)

		//插入浏览记录
		v1.POST("/history", controller.HistoryHandler)

		//获取用户浏览记录
		v1.GET("/user/history", controller.GetUserHistoryListHandler)

		//文件上传
		v1.POST("/upload", controller.UploadFileHandler)

		//删除评论
		v1.GET("/comment/:id", controller.CommentDeleteHandler)

		//重新提交帖子
		v1.POST("/post/resubmit", controller.ResubmitPostHandler)

		//删除帖子
		v1.GET("/post/delete/:id", controller.DeletePostHandler)

		////
	}

	v2 := r.Group("/api/v2")
	{
		v2.POST("/login", controller.SuperUserLoginHandler)
	}
	v2.Use(middlewares.JWTAuthMiddleware())
	{
		//创建板块信息
		v2.POST("/create_chunk", controller.CreateChunkHandler)

		//审核帖子
		v2.POST("/examine", controller.ExaminePostHandler)

		//根据选择返回对应帖子列表
		v2.GET("/posts", controller.SuperUserGetPostListHandler)

		//获取帖子详细详细
		v2.GET("/post/:id", controller.GetPostDetailHandler)

		//删除帖子
		v2.GET("/post/delete/:id", controller.SuperuserDeletePostHandler)

		//删除评论
		v2.GET("/comment/:id", controller.SuperUserCommentDeleteHandler)
	}

	pprof.Register(r) //注册pprof相关路由

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
