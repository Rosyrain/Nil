package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"nil/models"
)

//把每一步数据库操作封装成函数
//待logic层根据业务需求调用

const secret = "liweizhou.com"

// CheUserExist  检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`

	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return ErrorUserNotExist
}

func CheckUserExistByUID(user_id int64) (err error) {
	sqlStr := `select count(username) from user where user_id = ?`

	var count int
	if err = db.Get(&count, sqlStr, user_id); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return ErrorUserNotExist
}

// InserUser  向数据库中插入用户一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 执行sql语句入库
	sqlStr := `insert into user(user_id,username,password,email,gender,status) values(?,?,?,?,?,?)`

	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password, user.Email, user.Gender, user.Status)
	return err
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password //用户登录的密码
	sqlStr := `select user_id,username,password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}

	if err != nil {
		//查询数据库失败
		return err
	}

	//判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

func GetUserById(uID int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username from user where user_id=?`
	err = db.Get(user, sqlStr, uID)
	return
}

func GetUserByUsername(username string) (u *models.User, err error) {
	u = new(models.User)
	sqlStr := `select user_id,username,gender,email,create_time from user where username=?`
	err = db.Get(u, sqlStr, username)
	return
}

func CheckUserStatus(username string) error {
	sqlStr := `select status from user where user_name=?`
	var status int
	err := db.Get(&status, sqlStr, username)
	if err != nil {
		return err
	}
	if status == 0 {
		return ErrorUnActivate
	} else {
		return ErrorAnomaly
	}
	return nil
}

func VerifyActivate(status int, uID int64) error {
	sqlStr := `UPDATE user SET status = ? WHERE user_id = ?`
	_, err := db.Exec(sqlStr, status, uID)
	return err
}

func UpdateUserinfo(order, o, n string, uid int64) (err error) {
	fmt.Println(order, o, n, uid)
	switch order {
	case "password":
		opassword := encryptPassword(o)
		sqlStr1 := `select password from user where user_id=?`
		var password string
		err = db.Get(&password, sqlStr1, uid)
		if err == sql.ErrNoRows {
			return ErrorUserNotExist
		}
		if password != opassword {
			return ErrorInvalidPassword
		}
		npassword := encryptPassword(n)
		sqlStr2 := `UPDATE user SET password = :value WHERE user_id = :uid`
		_, err = db.NamedExec(sqlStr2, map[string]interface{}{"value": npassword, "uid": uid})

	case "email", "username", "gender":
		sqlStr := "UPDATE user SET " + order + " = :value WHERE user_id = :uid"
		_, err = db.NamedExec(sqlStr, map[string]interface{}{"value": n, "uid": uid})

	default:
		return fmt.Errorf("Invalid order: %s", order)
	}
	return err
}
