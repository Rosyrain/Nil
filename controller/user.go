package controller

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"nil/dao/mysql"
	"nil/dao/redis"
	"nil/logic"
	models "nil/models"
	"nil/pkg/email"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	p := new(models.ParamSignUp)

	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}

		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)), //翻译错误
		//})
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))

		return
	}

	//手动对请求参数进行详细的业务规则校验
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.Repassword) == 0 || p.Password != p.Repassword {
	//	//请求参数有误，直接返回响应
	//	zap.L().Error("SignUp with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数错误",
	//	})
	//	return
	//}

	fmt.Println(p)

	//2.业务处理
	if err := logic.SignUp(p); err != nil {
		//方便查看是什么位置出错
		//fmt.Println(err)
		zap.L().Error("logic.SignUp(p) failed", zap.Error(err))

		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		if errors.Is(err, redis.ErrCaptchaTimeExpire) {
			ResponseError(c, CodeCaptchaExpire)
			return
		}
		if errors.Is(err, redis.ErrCaptcha) {
			ResponseError(c, CodeNotCaptcha)
			return
		}

		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

func CaptchaHandler(c *gin.Context) {
	//1.获取参数和参数校验
	p := new(models.ParamActivate)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("Activate with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}

		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)), //翻译错误
		//})
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))

		return

	}

	//2.业务处理
	if err := logic.Captcha(p); err != nil {
		//方便查看是什么位置出错
		//fmt.Println(err)
		zap.L().Error("logic.Captcha(p) failed", zap.Error(err))

		if errors.Is(err, email.ErrorEmail) {
			ResponseError(c, CodeEmail)
			return
		}

		ResponseError(c, CodeServerBusy)
		return
	}

	//返回响应
	ResponseSuccess(c, nil)

}

// VerifyActivateHandler 验证激活
//func VerifyActivateHandler(c *gin.Context) {
//
//	//1.参数校验
//	username := c.Query("username")
//	user_id := c.Query("user_id")
//	if len(username) == 0 || len(user_id) == 0 {
//		zap.L().Error("VerifyActivateHandler failed", zap.Any("VerifyActivateHandler failed,err:", "参数错误"))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	//2.业务处理
//	if err := logic.VerifyActivate(user_id); err != nil {
//		zap.L().Error("logic.VerifyActivateHandler failed", zap.Error(err))
//		return
//	}
//
//	//3.返回响应
//	ResponseSuccess(c, nil)
//}

func LoginHandler(c *gin.Context) {
	//1.获取请求参数以及参数校验
	p := new(models.ParamLogin)

	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))

		return
	}
	fmt.Println(p)

	//2.业务处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   user.UserID, //如果ID值大于 2^53-1  userID最大值是2^63-1
		"user_name": user.Username,
		"token":     user.Token,
	})
}

func UserInfoHandler(c *gin.Context) {
	//1.参数校验
	username := c.Param("username")

	//2.业务处理
	u, err := logic.GetUseInfo(username)
	if err != nil {
		zap.L().Error("logic.GetUseInfo(username)", zap.String("username", username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":  u.UserID,
		"username": u.Username,
		"gender":   u.Gender,
		"email":    u.Email,
		"birthday": u.Birthday,
	})
}

func UpdateUserInfoHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamUpdateUserInfo)

	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("UpdateUserInfoHandler with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//业务处理
	userID, username, err := GetCurrentUser(c)
	if err != nil {
		zap.L().Error("GetCurrentUser(c) failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}

	p.UserId = userID
	p.UserName = username
	if err := logic.UpdateUserInfo(p); err != nil {
		zap.L().Error("logic.UpdateUserInfo(p) failed", zap.Any("username", p.UserName), zap.Error(err))
		if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		//修改password可能出现的报错
		if errors.Is(err, redis.ErrCaptchaTimeExpire) {
			ResponseError(c, CodeCaptchaExpire)
			return
		}
		if errors.Is(err, redis.ErrCaptcha) {
			ResponseError(c, CodeNotCaptcha)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//返回响应
	ResponseSuccess(c, nil)
}
