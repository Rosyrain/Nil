package logic

import (
	"nil/dao/mysql"
	"nil/models"
	"nil/pkg/jwt"
)

func SuperUserLogin(p *models.ParamLogin) (user *models.Superuser, err error) {
	user = &models.Superuser{
		Username: p.Username,
		Password: p.Password,
	}

	//传递的是一个指针，就能拿到user.UserID
	if err := mysql.SuperuserLogin(user); err != nil {
		return nil, err
	}

	//生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
