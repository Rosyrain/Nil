package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidToken
	CodeNeedLogin

	CodeNotCaptcha
	CodeCaptchaExpire

	CodeEmail

	CodeChunkExist
	CodeChunkNotExist

	CodeRepeated
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken:  "无效的token",
	CodeNeedLogin:     "需要登录",
	CodeNotCaptcha:    "验证码错误",
	CodeCaptchaExpire: "验证码过期",

	CodeEmail: "激活邮件发送失败",

	CodeChunkExist:    "板块已存在",
	CodeChunkNotExist: "板块不存在",

	CodeRepeated: "禁止重复相同操作",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
