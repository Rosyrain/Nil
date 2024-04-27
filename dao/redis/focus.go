package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"nil/models"
	"strconv"
	"time"
)

var (
	ErrorFocusRepeated = errors.New("禁止重复相同操作")
)

func InsertFocus(p *models.ParamFocusData, c_uid int64) (err error) {
	value := float64(p.Direction)
	ukey := GetRedisKey(KeyUserFocusPrefix + strconv.Itoa(int(c_uid)))
	ov := client.ZScore(ukey, p.UserID).Val()

	if ov == value {
		return ErrorFocusRepeated
	}

	op := value - ov
	if op > 0 {
		client.ZAdd(ukey, redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: p.UserID,
		})
	} else {
		//删除关注用户列表
		client.ZRem(ukey, p.UserID)
	}

	return
}

func GetUserFocusIDs(p *models.ParamFocusList) (data []string, err error) {
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1

	ukey := GetRedisKey(KeyUserFocusPrefix + strconv.Itoa(int(p.UserID)))

	//3.ZREVRANGE  按分数从大到小查询指定数量的元素
	return client.ZRevRange(ukey, start, end).Result()
}
