package redis

import (
	"errors"
	"nil/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	FiveMinutesINSeconds = 5 * 60 //五分钟的时间
)

var (
	ErrCaptchaTimeExpire = errors.New("验证码已过期")
	ErrCaptcha           = errors.New("验证码错误")
)

func CreatePost(postID, ChunkID, uid int64) error {
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(GetRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//帖子分数
	pipeline.ZAdd(GetRedisKey(KeyPostScore), redis.Z{
		Score:  float64(0),
		Member: postID,
	})

	//把帖子id加到社区的set
	cKey := GetRedisKey(KeyChunkPrefix + strconv.Itoa(int(ChunkID)))
	pipeline.SAdd(cKey, postID)

	//把帖子id加到user的set下
	ukey := GetRedisKey(KeyUserPostPrefix + strconv.Itoa(int(uid)))
	pipeline.SAdd(ukey, postID)

	_, err := pipeline.Exec() //事务同时成功或者同时失败

	return err
}

func GetIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1

	//3.ZREVRANGE  按分数从大到小查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()

}

func GetPostListIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//根据用户请求中携带的order参数确定要查询的redis key
	key := GetRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		key = GetRedisKey(KeyPostScore)
	}

	//2.确定查询的索引起始
	return GetIDsFormKey(key, p.Page, p.Size)
}

// GetChunkPostIDsInOrder  根据社区查询ids
func GetChunkPostListIDsInOrder(p *models.ParamPostList) (data []string, err error) {

	orderKey := GetRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyPostScore)
	}

	//使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	//针对新的zset 按之前的逻辑取出数据

	//社区的key
	cKey := GetRedisKey(KeyChunkPrefix + strconv.Itoa(int(p.ChunkID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.ChunkID))
	if client.Exists(key).Val() < 1 {
		//不存在，需要计算
		//pipeline := client.TxPipeline()

		//组合一个临时Zset集合存储查询结果(60s),作为缓存
		//这个ordKey形如 post:time:12121 下面存id+time/score
		client.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey) //zinterstore 计算
		client.Expire(key, 60*time.Second)
		//_, err = pipeline.Exec()
		//if err != redisNil {
		//	return
		//}

	}

	//存在的话直接根据key查询ids
	return GetIDsFormKey(key, p.Page, p.Size)
}

// GetChunkPostIDsInOrder  根据社区查询ids
func GetChunkCheckPostListIDsInOrder(p *models.ParamPostList) (data []string, err error) {

	orderKey := GetRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyPostScore)
	}

	//使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	//针对新的zset 按之前的逻辑取出数据

	//社区的key
	cKey := GetRedisKey(KeyChunkCheckPrefix + strconv.Itoa(int(p.ChunkID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.ChunkID))
	if client.Exists(key).Val() < 1 {
		//不存在，需要计算
		//pipeline := client.TxPipeline()

		//组合一个临时Zset集合存储查询结果(60s),作为缓存
		//这个ordKey形如 post:time:12121 下面存id+time/score
		client.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey) //zinterstore 计算
		client.Expire(key, 60*time.Second)
		//_, err = pipeline.Exec()
		//if err != redisNil {
		//	return
		//}

	}

	//存在的话直接根据key查询ids
	return GetIDsFormKey(key, p.Page, p.Size)
}

func GetUserPostIDsInOrder(p *models.ParamPostList) (data []string, err error) {
	orderKey := GetRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyPostScore)
	}

	//用户的key
	uKey := GetRedisKey(KeyUserPostPrefix + strconv.Itoa(int(p.UserID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.UserID))
	if client.Exists(key).Val() < 1 {
		//不存在，需要计算
		//pipeline := client.TxPipeline()

		//组合一个临时Zset集合存储查询结果(60s),作为缓存
		//这个ordKey形如 post:time:12121 下面存id+time/score
		client.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, uKey, orderKey) //zinterstore 计算
		client.Expire(key, 60*time.Second)
		//_, err = pipeline.Exec()
		//if err != redisNil {
		//	return
		//}
	}

	//存在的话直接根据key查询ids
	return GetIDsFormKey(key, p.Page, p.Size)
}
