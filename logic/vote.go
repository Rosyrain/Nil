package logic

import (
	"nil/dao/redis"
	"nil/models"
	"strconv"

	"go.uber.org/zap"
)

//投票功能
//基于用户投票的相关算法 http://www.ruanyifeng .com/blog/algorithm/

//本项目使用简化版的投票分数
//用户投一票加423分 (86400是一天的秒数) 86400/200 -> 需要200张赞成票才能可以给帖子续一天 ->出自《redis实战》

//投票的几种情况
/*
direction=1,两种情况：
	1.之前未投票，现在投赞成票
	2.之前反对票，现在改投赞成票

direction=0,两种情况：
	1.之前投赞成票，现在取消投票
	2.之前投反对票，现在取消投票

direction=-1,两种情况：
	1.之前未投票，现在投反对票
	2.之前赞成票，现在改投反对票

投票的限制：
每个帖子子发表之日起，一个星期内允许投票，超过一个星期不允许投票（减轻后端存储压力）
	1.到期之后将redis中保存的赞成票与反对票存储到mysql表中

*/

// PostForVote  投票功能实现
func PostForVote(userID int64, p *models.ParamPostVoteData) error {
	zap.L().Debug("PostForVote",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))

}

func MCommentForVote(userID int64, p *models.ParamCommentVoteData) error {
	zap.L().Debug("PostForVote",
		zap.Int64("userID", userID),
		zap.String("postID", p.CommentID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForMComment(strconv.Itoa(int(userID)), p.CommentID, float64(p.Direction))

}

func SubCommentForVote(userID int64, p *models.ParamCommentVoteData) error {
	zap.L().Debug("PostForVote",
		zap.Int64("userID", userID),
		zap.String("postID", p.CommentID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForSubComment(strconv.Itoa(int(userID)), p.CommentID, float64(p.Direction))

}
