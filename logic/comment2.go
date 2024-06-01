package logic

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"nil/dao/mysql"
	"nil/dao/redis"
	"strconv"
)

func SuperuserCommentDelete(cid, uid int64) error {
	//1.检查用户是否存在
	err := mysql.CheckSuperUserExit(uid)
	if !errors.Is(err, mysql.ErrorSuperuserExist) {
		zap.L().Error("mysql.CheckUserExistByUID(id) failed", zap.Error(err))
		return err
	}

	//Superuser删除评论不需要检查权限问题,但需要拿到pid做相关缓存删除
	//2.检查用户是否有删除权限(是否为评论的发布者)
	c, err := mysql.GetCommentByID(cid)
	if err != nil {
		zap.L().Error("mysql.GetCommentByID(cid) failed", zap.Error(err))
		return err
	}

	//if c.AuthorID != uid {
	//	return mysql.ErrorNotPower
	//}

	//3.删除mysql与redis中的记录
	//3.1 先查询该评论的子评论的ids
	ids, err := redis.GetAllSubCommentIDs(cid)
	ids = append(ids, strconv.Itoa(int(cid)))
	fmt.Println("ids:", ids)

	//ps:这里考虑将其看作一个事务，同时发生或都不发生
	//3.2 去mysql中删除所以的评论记录
	err = mysql.DeleteCommentByIDs(ids)
	if err != nil {
		return err
	}

	//3.2 删除redis中的记录
	err = redis.DeleteComment(ids, cid, uid, c.PostID)
	if err != nil {
		return err
	}

	//返回响应
	return nil
}
