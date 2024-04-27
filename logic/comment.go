package logic

import (
	"go.uber.org/zap"
	"nil/dao/mysql"
	"nil/dao/redis"
	"nil/models"
	snowflake "nil/pkg/snowflask"
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

func CommentToComment(p *models.Comment, mcomment_id int64) (err error) {
	//1.记录次级comment到mysql当中
	p.CommentID = snowflake.GenID()
	err = mysql.CreateComment(p)
	if err != nil {
		return
	}
	//2.记录到redis中主评论的zset中
	err = redis.CreateSubComment(p, mcomment_id)
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
	ids, err := redis.GetSubCommentIDsInOrder(p)
	if err != nil {
		zap.L().Error("GetPostMCommentList failed", zap.Error(err))
		return nil, err
	}
	//fmt.Println("ids:", ids)
	if len(ids) == 0 {
		zap.L().Warn("redis.GetSubCommentIDsInOrder(p) success but return 0 data")
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
