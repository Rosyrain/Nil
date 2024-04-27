package logic

import (
	"go.uber.org/zap"
	"nil/dao/mysql"
	"nil/dao/redis"
	"nil/models"
	"strconv"
)

func InsertHistory(post_id string, uid int64) error {
	return redis.InsertHistory(post_id, uid)
}

func GetUserHistoryList(p *models.ParamHistoryList) (data []*models.Post, err error) {
	ids, err := redis.GetUserHistoryList(p)
	if err != nil {
		zap.L().Error("redis.GetUserHistoryList(p) failed", zap.Error(err))
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetUserHistoryList(p) success but return 0 data")
		return
	}

	data = make([]*models.Post, 0, len(ids))

	for _, pid := range ids {
		id, _ := strconv.ParseInt(pid, 10, 64)
		post, err := mysql.GetPostByID(id)
		if err != nil {
			continue
		}
		data = append(data, post)
	}
	return
}
