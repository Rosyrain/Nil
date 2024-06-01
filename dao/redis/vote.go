package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
)

//本项目使用简化版的投票分数
//用户投一票加423分 (86400是一天的秒数) 86400/200 -> 需要200张赞成票才能可以给帖子续一天 ->出自《redis实战》

/*
direction=1,两种情况：
	1.之前未投票，现在投赞成票		-->更新分数和投票记录	差值：1  +432
	2.之前反对票，现在改投赞成票	-->更新分数和投票记录	差值：2	+432*2

direction=0,两种情况：
	1.之前投反对票，现在取消投票	-->更新分数和投票记录	差值：1	+432
	2.之前投赞成票，现在取消投票	-->更新分数和投票记录	差值：1  -432

direction=-1,两种情况：
	1.之前未投票，现在投反对票		-->更新分数和投票记录	差值：1	-432
	2.之前赞成票，现在改投反对票	-->更新分数和投票记录	差值：2	-432*2

投票的限制：
每个帖子子发表之日起，一个星期内允许投票，超过一个星期不允许投票（减轻后端存储压力）
	1.到期之后讲redis中保存的赞成票与反对票存储到mysql表中

*/
//投票的几种情况

const (
	oneWeekINSeconds = 7 * 24 * 3600
	scorePerVote     = 1 //每一票是多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制

	//暂时取消投票时限
	//取redis取出帖子发布时间
	//postTime := client.ZScore(GetRedisKey(KeyPostTime), postID).Val()
	//if float64(time.Now().Unix())-postTime > oneWeekINSeconds {
	//	return ErrVoteTimeExpire
	//}

	//2和3需要放到一个pipeline事务中操作

	//2.更新分数
	//1.先查当前用户给当前帖子的投票记录
	ov := client.ZScore(GetRedisKey(KeyPostVotedPrefix)+postID, userID).Val()

	//更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}

	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算差值两次投票

	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(GetRedisKey(KeyPostScore), op*diff*scorePerVote, postID)

	//3.记录用户为改帖子投票的数据
	if value == 0 {
		pipeline.ZRem(GetRedisKey(KeyPostVotedPrefix)+postID, userID)
	}
	pipeline.ZAdd(GetRedisKey(KeyPostVotedPrefix)+postID, redis.Z{
		Score:  value, //当前用户投的是赞成票或反对票
		Member: userID,
	})
	_, err := pipeline.Exec()
	return err
}

func VoteForMComment(userID, commentID string, value float64) error {
	//1.判断投票限制

	//暂时取消投票时限
	//取redis取出帖子发布时间
	//postTime := client.ZScore(GetRedisKey(KeyPostTime), postID).Val()
	//if float64(time.Now().Unix())-postTime > oneWeekINSeconds {
	//	return ErrVoteTimeExpire
	//}

	//2和3需要放到一个pipeline事务中操作

	//2.更新分数
	//1.先查当前用户给当前帖子的投票记录
	ov := client.ZScore(GetRedisKey(KeyCommentVotedPrefix)+commentID, userID).Val()

	//更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}

	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算差值两次投票

	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(GetRedisKey(KeyCommentScore), op*diff*scorePerVote, commentID)

	//3.记录用户为改帖子投票的数据
	if value == 0 {
		pipeline.ZRem(GetRedisKey(KeyCommentVotedPrefix)+commentID, userID)
	}
	pipeline.ZAdd(GetRedisKey(KeyCommentVotedPrefix)+commentID, redis.Z{
		Score:  value, //当前用户投的是赞成票或反对票
		Member: userID,
	})
	_, err := pipeline.Exec()
	return err
}

func VoteForSubComment(userID, commentID string, value float64) error {
	//1.判断投票限制

	//暂时取消投票时限
	//取redis取出帖子发布时间
	//postTime := client.ZScore(GetRedisKey(KeyPostTime), postID).Val()
	//if float64(time.Now().Unix())-postTime > oneWeekINSeconds {
	//	return ErrVoteTimeExpire
	//}

	//2和3需要放到一个pipeline事务中操作

	//2.更新分数
	//1.先查当前用户给当前帖子的投票记录
	ov := client.ZScore(GetRedisKey(KeySubCommentVotedPrefix)+commentID, userID).Val()

	//更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}

	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算差值两次投票

	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(GetRedisKey(KeySubCommentScore), op*diff*scorePerVote, commentID)

	//3.记录用户为改帖子投票的数据
	if value == 0 {
		pipeline.ZRem(GetRedisKey(KeySubCommentVotedPrefix)+commentID, userID)
	}
	pipeline.ZAdd(GetRedisKey(KeySubCommentVotedPrefix)+commentID, redis.Z{
		Score:  value, //当前用户投的是赞成票或反对票
		Member: userID,
	})
	_, err := pipeline.Exec()
	return err
}

// GetPostVoteData  根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := GetRedisKey(KeyPostVotedPrefix) + id
		//查找key中分数为1的元素的数量->统计每篇帖子赞成票的数量
		v := client.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}
	return

	//使用pipeline一次发送多条命令，减少RTT
	//pipeline := client.TxPipeline()
	//for _, id := range ids {
	//	key := GetRedisKey(KeyPostVotedPrefix + id)
	//	zap.L().Debug("GetPostVoteData", zap.Any("key", key))
	//	//v := client.ZCount(key, "1", "1").Val()
	//	//zap.L().Debug("GetPostVoteData", zap.Any("v", v))
	//	Zcount := pipeline.ZCount(key, "1", "1")
	//	zap.L().Debug("GetPostVoteData", zap.Any("Zcount", Zcount))
	//}
	//cmders, err := pipeline.Exec()
	//zap.L().Debug("GetPostVoteData", zap.Any("cmders", cmders))
	//if err != nil {
	//	return
	//}
	//data = make([]int64, 0, len(ids))
	//for _, cmder := range cmders {
	//	v := cmder.(*redis.IntCmd).Val()
	//	zap.L().Debug("GetPostVoteData", zap.Any("v", v))
	//	data = append(data, v)
	//}
	//return data, err
}

func GetCommentVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := GetRedisKey(KeyCommentVotedPrefix) + id
		//查找key中分数为1的元素的数量->统计每篇帖子赞成票的数量
		v := client.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}
	return
}

func GetSubCommentVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := GetRedisKey(KeySubCommentVotedPrefix) + id
		//查找key中分数为1的元素的数量->统计每篇帖子赞成票的数量
		v := client.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}
	return
}
