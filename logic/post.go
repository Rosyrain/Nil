package logic

import (
	"nil/dao/mysql"
	"nil/dao/redis"
	"nil/models"
	snowflake "nil/pkg/snowflask"

	"go.uber.org/zap"
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
	ids, err := redis.GetPostListIDsInOrder(p)
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
	ids, err := redis.GetChunkCheckPostListIDsInOrder(p)
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
