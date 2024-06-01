package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"nil/models"
	"reflect"
	"strconv"
	"time"
)

var (
	searchStatusMap = map[string]interface{}{
		"0": getChunkToBeReviewPostListIDsInOrder,
		"1": getChunkNormalPostListIDsInOrder,
		"3": getChunkToBeDeletePostListIDsInOrder,
	}
)

// getChunkToBeReviewPostListIDsInOrder  根据社区查询ids（待审核的）
func getChunkToBeReviewPostListIDsInOrder(p *models.ParamSearch) (data []string, err error) {
	orderKey := GetRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyPostScore)
	}

	//使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	//针对新的zset 按之前的逻辑取出数据

	//社区的key
	cKey := GetRedisKey(KeyChunkToBeReviewPrefix + strconv.Itoa(int(p.ChunkID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderKey + "super:tobereview:" + strconv.Itoa(int(p.ChunkID))
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

// getChunkToBeDeletePostListIDsInOrder  根据社区查询ids（待确认删除的）
func getChunkToBeDeletePostListIDsInOrder(p *models.ParamSearch) (data []string, err error) {

	orderKey := GetRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyPostScore)
	}

	//使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	//针对新的zset 按之前的逻辑取出数据

	//社区的key
	cKey := GetRedisKey(KeyChunkToBeDeletePrefix + strconv.Itoa(int(p.ChunkID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderKey + "super:delete:" + strconv.Itoa(int(p.ChunkID))
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

// GetChunkNormalPostIDsInOrder  根据社区查询ids（审核通过的）
func getChunkNormalPostListIDsInOrder(p *models.ParamSearch) (data []string, err error) {
	orderKey := GetRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyPostScore)
	}

	//使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	//针对新的zset 按之前的逻辑取出数据

	//社区的key
	cKey := GetRedisKey(KeyChunkNormalPrefix + strconv.Itoa(int(p.ChunkID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderKey + "super:normal:" + strconv.Itoa(int(p.ChunkID))
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

func SuperuserGetPostList(p *models.ParamSearch) (data []string, err error) {
	if _, ok := searchStatusMap[strconv.Itoa(p.Status)]; !ok {
		return nil, errors.New("查询的status超出设定范围")
	}
	funcName := reflect.ValueOf(searchStatusMap[strconv.Itoa(p.Status)])
	results := funcName.Call([]reflect.Value{reflect.ValueOf(p)})
	if len(results) >= 2 {
		if err, ok := results[1].Interface().(error); ok && err != nil {
			// 处理错误
			return nil, err
		} else {
			// 获取data
			data, _ := results[0].Interface().([]string)
			return data, nil
		}
	} else {
		zap.L().Error("function call returned unexpected number of results", zap.Any("len(results):", len(results)))
		return nil, errors.New("redis,SuperuserGetPostList使用reflect发送错误")
	}
}

func SuperuserDeletePost(pid, uid, cid int64) (err error) {
	//需要删除的地方有 postTime,postScore,postComment,userPost;还有各个评论的投票;chunk:
	oKey := GetRedisKey(KeyPostTime)
	sKey := GetRedisKey(KeyPostScore)
	uKey := GetRedisKey(KeyUserPostPrefix + strconv.Itoa(int(uid)))
	chKey := GetRedisKey(KeyChunkPrefix + strconv.Itoa(int(cid)))
	vpKey := GetRedisKey(KeyPostVotedPrefix + strconv.Itoa(int(cid)))
	cpKey := GetRedisKey(KeyPostCommentPrefix + strconv.Itoa(int(pid)))

	//1.先删除time和score中的记录
	_, err = client.ZRem(oKey, pid).Result()
	if err != nil {
		return
	}
	_, err = client.ZRem(sKey, pid).Result()
	if err != nil {
		return
	}

	//2.删除用户下面的帖子
	_, err = client.SRem(uKey, pid).Result()
	if err != nil {
		return
	}

	//4.删除评论对应的投票信息
	_, err = client.Del(vpKey).Result()
	if err != nil {
		return err
	}

	//5.删除chunk下的
	_, err = client.SRem(chKey, cid).Result()
	if err != nil {
		return err
	}

	//6.删除评论
	_, err = client.Del(cpKey).Result()
	if err != nil {
		return err
	}

	return
}

func DeletePostFromTobeDeleteCache(pid, cid int64) error {
	_, err := client.SRem(GetRedisKey(KeyChunkToBeDeletePrefix)+strconv.Itoa(int(cid)), pid).Result()
	return err
}
