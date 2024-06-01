package logic

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"nil/dao/mysql"
	"nil/dao/redis"
	"nil/models"
	snowflake "nil/pkg/snowflask"
	"strconv"
)

func CommentToPost(p *models.Comment) (err error) {
	//1.记录comment到mysql当中
	p.CommentID = snowflake.GenID()
	err = mysql.CreateComment(p)
	if err != nil {
		return
	}
	//2.记录到redis中post的主评论zset中
	err = redis.CreateComment(p)
	if err != nil {
		return
	}
	//3.返回
	return
}

func CommentToComment(p *models.Comment, commentid int64) (err error) {
	//1.记录次级comment到mysql当中
	p.CommentID = snowflake.GenID()
	err = mysql.CreateComment(p)
	if err != nil {
		return
	}
	//2.记录到redis中主评论的zset中
	err = redis.CreateSubComment(p, commentid)
	//3.返回
	return
}

func GetUserCommentList(p *models.ParamCommentList) (data []*models.ApiCommentDetail, err error) {
	//1.先拿到用户评论的id列表
	ids, err := redis.GetUserCommentIDsInOrder(p)
	if err != nil {
		zap.L().Error("GetUserCommentList failed", zap.Error(err))
		return nil, err
	}
	//fmt.Println("ids:", ids)
	if len(ids) == 0 {
		zap.L().Warn("redis.GetUserCommentIDsInOrder success but return 0 data")
		return
	}

	zap.L().Debug("GetUserCommentList", zap.Any("ids", ids))

	//
	data = make([]*models.ApiCommentDetail, 0, len(ids))

	//3.去mysql中拿到comment的基本信息列表
	comments, err := mysql.GetCommentListByIDs(ids)
	if err != nil {
		return
	}

	//4.提取查询得票数
	voteData, err := redis.GetCommentVoteData(ids)

	for idx, comment := range comments {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(comment.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("authorID", comment.AuthorID),
				zap.Error(err))
			continue
		}

		//根据post_id查询post信息
		post, err := mysql.GetPostByID(comment.PostID)
		if err != nil {
			zap.L().Error("mysql.GetPostByID(comment.PostID) failed",
				zap.Int64("postId", comment.PostID),
				zap.Error(err))
			continue
		}

		commentdetail := &models.ApiCommentDetail{
			AuthorName: user.Username,
			VoteNumber: voteData[idx],
			Comment:    comment,
			Post:       post,
		}
		data = append(data, commentdetail)
	}

	return
}

func GetPostMCommentList(p *models.ParamCommentList) (data []*models.ApiCommentDetail, err error) {
	//1.先拿到帖子评论的id列表
	ids, err := redis.GetPostCommentIDsInOrder(p)
	if err != nil {
		zap.L().Error("GetPostMCommentList failed", zap.Error(err))
		return nil, err
	}
	//fmt.Println("ids:", ids)
	if len(ids) == 0 {
		zap.L().Warn("redis.GetSubCommentIDsInOrder(p) success but return 0 data")
		return nil, nil
	}

	zap.L().Debug("GetPostMCommentList", zap.Any("ids", ids))

	//
	data = make([]*models.ApiCommentDetail, 0, len(ids))

	//3.去mysql中拿到comment的基本信息列表
	comments, err := mysql.GetCommentListByIDs(ids)
	if err != nil {
		return
	}

	//4.提取查询得票数
	voteData, err := redis.GetCommentVoteData(ids)

	for idx, comment := range comments {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(comment.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("authorID", comment.AuthorID),
				zap.Error(err))
			continue
		}

		//根据post_id查询post信息
		post, err := mysql.GetPostByID(comment.PostID)
		if err != nil {
			zap.L().Error("mysql.GetPostByID(comment.PostID) failed",
				zap.Int64("postId", comment.PostID),
				zap.Error(err))
			continue
		}

		commentdetail := &models.ApiCommentDetail{
			AuthorName: user.Username,
			VoteNumber: voteData[idx],
			Comment:    comment,
			Post:       post,
		}
		data = append(data, commentdetail)
	}

	return
}

func GetSubCommentList(p *models.ParamCommentList) (data []*models.ApiCommentDetail, err error) {
	//1.先拿到帖子评论的id列表
	ids, err := redis.GetSubCommentIDsInOrder(p)
	if err != nil {
		zap.L().Error("GetPostMCommentList failed", zap.Error(err))
		return nil, err
	}
	//fmt.Println("ids:", ids)
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostCommentIDsInOrder(p) success but return 0 data")
		return
	}

	zap.L().Debug("GetPostMCommentList", zap.Any("ids", ids))

	//
	data = make([]*models.ApiCommentDetail, 0, len(ids))

	//3.去mysql中拿到comment的基本信息列表
	comments, err := mysql.GetCommentListByIDs(ids)
	if err != nil {
		return
	}

	//4.提取查询得票数
	voteData, err := redis.GetSubCommentVoteData(ids)

	for idx, comment := range comments {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(comment.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("authorID", comment.AuthorID),
				zap.Error(err))
			continue
		}

		//根据post_id查询post信息
		post, err := mysql.GetPostByID(comment.PostID)
		if err != nil {
			zap.L().Error("mysql.GetPostByID(comment.PostID) failed",
				zap.Int64("postId", comment.PostID),
				zap.Error(err))
			continue
		}

		commentdetail := &models.ApiCommentDetail{
			AuthorName: user.Username,
			VoteNumber: voteData[idx],
			Comment:    comment,
			Post:       post,
		}
		data = append(data, commentdetail)
	}

	return
}

func CommentDelete(cid, uid int64) error {
	//1.检查用户是否存在
	err := mysql.CheckUserExistByUID(uid)
	if !errors.Is(err, mysql.ErrorUserExist) {
		zap.L().Error("mysql.CheckUserExistByUID(id) failed", zap.Error(err))
		return err
	}

	//2.检查用户是否有删除权限(是否为评论的发布者)
	c, err := mysql.GetCommentByID(cid)
	if err != nil {
		zap.L().Error("mysql.GetCommentByID(cid) failed", zap.Error(err))
		return err
	}

	if c.AuthorID != uid {
		return mysql.ErrorNotPower
	}

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

func CommentDeleteForPost(cid, uid, pid int64) error {

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
	err = redis.DeleteComment(ids, cid, uid, pid)
	if err != nil {
		return err
	}

	//返回响应
	return nil
}
