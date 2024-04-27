package logic

import (
	"go.uber.org/zap"
	"nil/dao/mysql"
	"nil/dao/redis"
	"nil/models"
	"strconv"
)

func InsertFocus(p *models.ParamFocusData, c_uid int64) error {
	return redis.InsertFocus(p, c_uid)
}

func GetUserFocusList(p *models.ParamFocusList) (data []*models.User, err error) {
	//1.拿到
	ids, err := redis.GetUserFocusIDs(p)
	if err != nil {
		zap.L().Error("redis.GetUserFocusIDs(uid) failed", zap.Error(err))
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetUserFocusIDs(uid) success but return 0 data")
		return
	}

	data = make([]*models.User, 0, len(ids))

	for _, uid := range ids {
		id, _ := strconv.ParseInt(uid, 10, 64)
		user, err := mysql.GetUserById(id)
		if err != nil {
			continue
		}
		data = append(data, user)
	}
	return
}
