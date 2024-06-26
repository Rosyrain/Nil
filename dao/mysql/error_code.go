package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("密码错误")
	ErrorInvalidID       = errors.New("无效的ID")
	ErrorUnActivate      = errors.New("用户未激活")
	ErrorAnomaly         = errors.New("用户状态异常")

	ErrorChunkExist        = errors.New("板块已存在")
	ErrorChunkNotExist     = errors.New("板块不存在")
	ErrorSuperuserExist    = errors.New("管理员已存在")
	ErrorSuperuserNotExist = errors.New("管理员不存在")

	ErrorNotPower = errors.New("该用户没有权限")
)
