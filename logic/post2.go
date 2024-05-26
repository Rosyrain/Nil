package logic

import (
	"errors"
	"nil/dao/mysql"
	"nil/models"
)

func ExaminePost(p *models.ParamExamine) error {
	//0.查看管理员用户是否存在
	err := mysql.CheckSuperUserExit(p.UserID)
	if !errors.Is(err, mysql.ErrorSuperuserExist) {
		return err
	}

	//1.查看管理员是否有权限
	if err := mysql.CheckSuperUserPower(p); err != nil {
		return err
	}

	//2.修改状态
	return mysql.UpdatePostStatus(p)

	//3.返回响应

}
