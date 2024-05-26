package mysql

import (
	"database/sql"
	"nil/models"
)

func SuperuserLogin(user *models.Superuser) (err error) {
	oPassword := user.Password //用户登录的密码
	sqlStr := `select user_id,username,password from superuser where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}

	if err != nil {
		//查询数据库失败
		return err
	}

	//判断密码是否正确
	if oPassword != user.Password {
		return ErrorInvalidPassword
	}
	return
}

func CheckSuperUserExit(uid int64) (err error) {
	sqlStr := `select count(user_id) from superuser where user_id = ?`

	var count int
	if err = db.Get(&count, sqlStr, uid); err != nil {
		return err
	}
	if count > 0 {
		return ErrorSuperuserExist
	}
	return ErrorSuperuserNotExist
}
