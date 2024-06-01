package logic

import (
	"errors"
	"go.uber.org/zap"
	"nil/dao/mysql"
	"nil/dao/redis"
	"nil/models"
	snowflake "nil/pkg/snowflask"
)

func CreatePost(p *models.Post) (err error) {
	//1.生成post id
	p.ID = snowflake.GenID()

	//2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.ChunkID, p.AuthorID)
	//3.返回
	return
}

// GetPostDetail 根据帖子id查询帖子详情数据
func GetPostDetail(pid int64) (data *models.ApiPostDetail, err error) {
	//查询数据并拼接组合接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}

	data = new(models.ApiPostDetail)

	//根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("authorID", post.AuthorID),
			zap.Error(err))
		return
	}
	//根据社区id查询社区详情信息
	chunk, err := mysql.GetChunkDetailByID(post.ChunkID)
	if err != nil {
		zap.L().Error("mysql.GetChunkDetailByID(post.ChunkID) failed",
			zap.Int64("authorID", post.AuthorID),
			zap.Error(err))
		return
	}

	data = &models.ApiPostDetail{
		AuthorName:  user.Username,
		Post:        post,
		ChunkDetail: chunk,
	}
	return
}

// GetPostList  获取帖子列表	直接从mysql中查询
func GetPostList(offset, limit int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(offset, limit)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详情信息
		chunk, err := mysql.GetChunkDetailByID(post.ChunkID)
		if err != nil {
			zap.L().Error("mysql.GetChunkDetailByID(post.ChunkID) failed",
				zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}

		postdetail := &models.ApiPostDetail{
			AuthorName:  user.Username,
			Post:        post,
			ChunkDetail: chunk,
		}
		data = append(data, postdetail)
	}

	return
}

// GetPostList2 从获取帖子列表（不分社区，只按照时间和分数）
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//2.取redis查询id列表
	ids, err := redis.GetNormalPostListIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostListIDsInOrder success but return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))

	//3.根据id取数据库中查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)

	//将帖子的作者及分区查询出来填充至帖子中
	data = make([]*models.ApiPostDetail, 0, len(posts))
	//提前查询好每篇帖子的投票数据
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("GetPostList2 ", zap.Any("vote data", voteData))

	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详情信息
		chunk, err := mysql.GetChunkDetailByID(post.ChunkID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}

		postdetail := &models.ApiPostDetail{
			AuthorName:  user.Username,
			VoteNumber:  voteData[idx],
			Post:        post,
			ChunkDetail: chunk,
		}
		data = append(data, postdetail)
	}

	return
}

// GetChunkPostList 根据板块，时间/分数获取帖子列表
func GetChunkPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//2.取redis查询id列表
	ids, err := redis.GetChunkNormalPostListIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetChunkCheckPostListIDsInOrder success but return 0 data")
		return
	}
	zap.L().Debug("GetCommunityPostList", zap.Any("ids", ids))

	//3.根据id取数据库中查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)

	//将帖子的作者及分区查询出来填充至帖子中
	data = make([]*models.ApiPostDetail, 0, len(posts))
	//提前查询好每篇帖子的投票数据
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("GetChunkPostList ", zap.Any("vote data", voteData))

	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详情信息
		chunk, err := mysql.GetChunkDetailByID(post.ChunkID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}

		postdetail := &models.ApiPostDetail{
			AuthorName:  user.Username,
			VoteNumber:  voteData[idx],
			Post:        post,
			ChunkDetail: chunk,
		}
		data = append(data, postdetail)
	}

	return
}

// GetPostListNew  将两个查询逻辑合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.ChunkID == 0 {
		//	查所有
		data, err = GetPostList2(p)
	} else {
		// 根据社区查询
		data, err = GetChunkPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}

// GetUserPostList	获取用户帖子列表 日期/分数
func GetUserPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//1.根据要求获取帖子id列表
	ids, err := redis.GetUserPostIDsInOrder(p)
	if err != nil {
		zap.L().Error("GetUserPostList failed", zap.Error(err))
		return nil, err
	}
	//2.参数判断
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostListIDsInOrder success but return 0 data")
		return
	}
	zap.L().Debug("GetUserPostList", zap.Any("ids", ids))

	data = make([]*models.ApiPostDetail, 0, len(ids))

	//3.根据id取数据库中查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)

	//4.提取查询得票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("GetChunkPostList ", zap.Any("vote data", voteData))

	//idx->index
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详情信息
		chunk, err := mysql.GetChunkDetailByID(post.ChunkID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("authorID", post.AuthorID),
				zap.Error(err))
			continue
		}

		postdetail := &models.ApiPostDetail{
			AuthorName:  user.Username,
			VoteNumber:  voteData[idx],
			Post:        post,
			ChunkDetail: chunk,
		}
		data = append(data, postdetail)
	}

	return

}

func ResubmitPost(p *models.ParamResubmitPost) error {
	//2.2判断chunk_id是否发生改变，发生改变也需要修改缓存
	//这点提前了，我需要知道之前的chunk_id才能完成下面的 1.
	oldPost, err := mysql.GetPostByID(p.PostID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID failed", zap.Error(err))
		return err
	}

	//1.先将Post从之前状态的缓存中进行删除
	if err := redis.DeletePostCache(p.Status, p.PostID, oldPost.ChunkID); err != nil {
		zap.L().Error("redis.DeletePostCache failed", zap.Error(err))
		return err
	}

	//2.将帖子放入审核缓存
	if err := redis.AddPostTobeReviewCache(p.PostID, p.Post.ChunkID); err != nil {
		zap.L().Error("redis.AddPostTobeReviewCache failed", zap.Error(err))
		return err
	}

	if oldPost.ChunkID != p.Post.ChunkID {
		//2.2.1删除原来的chunk下的pid
		err := redis.DeletePostFromChunkCache(p.PostID, oldPost.ChunkID)
		if err != nil {
			zap.L().Error("redis.DeleteChunkCache failed", zap.Error(err))
			return err
		}

		//2.2.1将pid加入新的chunk缓存
		err = redis.AddPostToChunkCache(p.PostID, p.Post.ChunkID)
		if err != nil {
			zap.L().Error("redis.AddPostToChunkCache failed", zap.Error(err))
			return err
		}

	}

	//3.将帖子内容更新至数据库当中
	p.Post.Status = 0
	err = mysql.ResubmitPost(p.Post)

	return err
}

func DeletePost(pid, uid int64) (err error) {
	//1.拿到post的信息
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		return err
	}

	//1.检查用户是否有权限
	if post.AuthorID != uid {
		zap.L().Error("DeletePost no power", zap.Any("没有权限-uid:", uid))
		return mysql.ErrorNotPower
	}

	if post.Status == 3 {
		return errors.New("以及进入待删除状态")
	}

	//2.更改post的状态
	if err = mysql.DeletePost(pid); err != nil {
		return err
	}

	//3.删除post原来的缓存
	err = redis.DeletePostCache(int(post.Status), pid, post.ChunkID)
	if err != nil {
		zap.L().Error("redis.DeletePostCache failed", zap.Error(err))
		return err
	}

	//4.加入待删除缓存
	err = redis.AddPostToDeleteCache(pid, post.ChunkID)
	if err != nil {
		zap.L().Error("redis.AddPostToDeleteCache failed", zap.Error(err))
		return err
	}

	//5.将有关的评论全部删除
	//5.1得到将帖子主评论ids
	p := models.ParamCommentList{
		PostID: pid,
		Page:   1,
		Size:   -1,
		Order:  "time",
	}

	comments, err := GetPostMCommentList(&p)
	if err != nil {
		return
	}

	for _, comment := range comments {
		if err := CommentDeleteForPost(comment.CommentID, comment.Post.AuthorID, comment.PostID); err != nil {
			continue
		}
	}

	//有关帖子的缓存，在管理端进行删除

	return nil
}
