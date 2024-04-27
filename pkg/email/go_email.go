package email

import (
	"errors"
	"fmt"
	setting "nil/settings"
)

var (
	from       string
	nick_name  string
	secret     string
	host       string
	base_body1 string
	base_body2 string
	port       int
)

var (
	ErrorEmail = errors.New("邮件发送失败")
)

//https://www.cnblogs.com/zhangyafei/p/13918050.html#:~:text=Go%E8%AF%AD%E8%A8%80%E5%8F%91%E9%82%AE%E4%BB%B6%201%201.%20%E7%99%BB%E5%BD%95QQ%E9%82%AE%E7%AE%B1%EF%BC%8C%E9%80%89%E6%8B%A9%E8%B4%A6%E6%88%B7%EF%BC%8C%E5%BC%80%E5%90%AFPOP3%2FSMTP%E6%9C%8D%E5%8A%A1%E5%92%8CIMAP%2FSMTP%E6%9C%8D%E5%8A%A1%2C%E5%B9%B6%E7%94%9F%E6%88%90%E6%8E%88%E6%9D%83%E7%A0%81%202%202.%20%E4%BD%BF%E7%94%A8go%E8%AF%AD%E8%A8%80%E7%9A%84smtp%E5%8C%85%E5%8F%91%E9%80%81%E9%82%AE%E4%BB%B6%20go_email%2Femail.go,is%20the%20type%20used%20for%20email%20messages%20

func InitEmail(config *setting.EmailConfig) error {
	from = config.OfficialEmail
	nick_name = config.NickName
	secret = config.Secret
	host = config.Host
	port = config.Port
	base_body1 = config.BaseBody1
	base_body2 = config.BaseBody2

	return nil
}

func MainSendEmail(username, captcha, email string) error {
	fmt.Println(captcha)
	var to = []string{}
	to = append(to, email)
	subject := "Welcome to Nil,nice to meet you."
	var body string
	if username != "" {
		body = fmt.Sprintf(base_body1, username, captcha)

	} else {
		body = fmt.Sprintf(base_body2, captcha)
	}
	fmt.Println(body)
	if err := SendEmailWithPool(to, from, secret, host, subject, body, nick_name, port); err != nil {
		fmt.Println("发送失败: ", err)
		return ErrorEmail
	} else {
		fmt.Println("发送成功")
		return nil
	}
}
