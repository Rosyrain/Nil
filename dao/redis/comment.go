package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"nil/models"
	"strconv"
	"time"
)

// CreateComment 创建主评论
func CreateComment(p *models.Comment) error {
	pipeline := client.TxPipeline()
	//评论时间
	pipeline.ZAdd(GetRedisKey(KeyCommentTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: p.CommentID,
	})

	//评论分数
	pipeline.ZAdd(GetRedisKey(KeyCommentScore), redis.Z{
		Score:  float64(0),
		Member: p.CommentID,
	})

	//把评论id加到帖子的set
	cKey := GetRedisKey(KeyPostCommentPrefix + strconv.Itoa(int(p.PostID)))
	pipeline.SAdd(cKey, p.CommentID)

	//把评论id加到user的set下
	ukey := GetRedisKey(KeyUserCommentPrefix + strconv.Itoa(int(p.AuthorID)))
	pipeline.SAdd(ukey, p.CommentID)

	_, err := pipeline.Exec() //事务同时成功或者同时失败

	return err
}

func CreateSubComment(p *models.Comment, mcomment_id int64) error {
	pipeline := client.TxPipeline()
	//评论时间
	pipeline.ZAdd(GetRedisKey(KeySubCommentTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: p.CommentID,
	})

	//评论分数
	pipeline.ZAdd(GetRedisKey(KeySubCommentScore), redis.Z{
		Score:  float64(0),
		Member: p.CommentID,
	})

	//把次级评论id加到主评论的set
	cKey := GetRedisKey(KeyCommentPrefix + strconv.Itoa(int(mcomment_id)))
	pipeline.SAdd(cKey, p.CommentID)

	//把评论id加到user的set下
	ukey := GetRedisKey(KeyUserCommentPrefix + strconv.Itoa(int(p.AuthorID)))
	pipeline.SAdd(ukey, p.CommentID)

	_, err := pipeline.Exec() //事务同时成功或者同时失败

	return err
}

func GetUserCommentIDsInOrder(p *models.ParamCommentList) (data []string, err error) {
	orderKey := GetRedisKey(KeyCommentTime)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyCommentScore)
	}

	//用户的key
	uKey := GetRedisKey(KeyUserCommentPrefix + strconv.Itoa(int(p.UserID)))

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

	//fmt.Println("key:", key)
	//存在的话直接根据key查询ids
	return GetIDsFormKey(key, p.Page, p.Size)
}

func GetPostCommentIDsInOrder(p *models.ParamCommentList) (data []string, err error) {
	orderKey := GetRedisKey(KeyCommentTime)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyCommentScore)
	}

	//帖子的key
	pKey := GetRedisKey(KeyPostCommentPrefix + strconv.Itoa(int(p.PostID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.UserID))
	if client.Exists(key).Val() < 1 {
		//不存在，需要计算
		//pipeline := client.TxPipeline()

		//组合一个临时Zset集合存储查询结果(60s),作为缓存
		//这个ordKey形如 post:time:12121 下面存id+time/score
		client.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, pKey, orderKey) //zinterstore 计算
		client.Expire(key, 60*time.Second)
		//_, err = pipeline.Exec()
		//if err != redisNil {
		//	return
		//}
	}

	//fmt.Println("key:", key)
	//存在的话直接根据key查询ids
	return GetIDsFormKey(key, p.Page, p.Size)
}

func GetSubCommentIDsInOrder(p *models.ParamCommentList) (data []string, err error) {
	orderKey := GetRedisKey(KeySubCommentTime)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeySubCommentScore)
	}

	//帖子的key
	cKey := GetRedisKey(KeyCommentPrefix + strconv.Itoa(int(p.CommentID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommentID))
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
	fmt.Println("key", key)

	//fmt.Println("key:", key)
	//存在的话直接根据key查询ids
	return GetIDsFormKey(key, p.Page, p.Size)
}
