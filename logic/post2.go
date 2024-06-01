package logic

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"nil/dao/mysql"
	"nil/dao/redis"
	"nil/models"
)

func ExaminePost(p *models.ParamExamine) error {
	//0.查看管理员用户是否存在
	err := mysql.CheckSuperUserExit(p.SuperuserID)
	if !errors.Is(err, mysql.ErrorSuperuserExist) {
		return err
	}

	//1.查看管理员是否有权限
	if err := mysql.CheckSuperUserPower(p.SuperuserID, p.ChunkID); err != nil {
		return err
	}

	//2.修改状态
	if err := mysql.UpdatePostStatus(p); err != nil {
		return err
	}

	//3.删除在原本的status状态下的缓存
	if err := redis.DeletePostCache(p.Status, p.PostID, p.ChunkID); err != nil {
		return err
	}

	//3.2添加进缓存区
	if p.Direction == 1 { //帖子正常的缓存区
		if err := redis.AddPostToNormalCache(p.PostID, p.ChunkID); err != nil {
			return err
		}
	} else if p.Direction == 2 { //审核失败的没有缓冲区
		return nil
	} else {
		return errors.New("添加缓存异常")
	}

	//3.返回响应
	return nil
}

func SuperuserGetPostList(p *models.ParamSearch) (data []*models.ApiPostDetail, err error) {
	//0.查看管理员用户是否存在
	err = mysql.CheckSuperUserExit(p.SuperuserID)
	if !errors.Is(err, mysql.ErrorSuperuserExist) {
		return nil, err
	}

	//1.查看管理员是否有权限
	if err = mysql.CheckSuperUserPower(p.SuperuserID, p.ChunkID); err != nil {
		return nil, err
	}

	//2.根据chunk_id以及status得到ids
	ids, err := redis.SuperuserGetPostList(p)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.SuperuserGetPostList success but return 0 data")
		return
	}
	zap.L().Debug("SuperuserGetPostList", zap.Any("ids", ids))

	//3.根据id取数据库中查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)

	//将帖子的作者及分区查询出来填充至帖子中
	data = make([]*models.ApiPostDetail, 0, len(posts))
	//提前查询好每篇帖子的投票数据
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("SuperuserGetPostList ", zap.Any("vote data", voteData))

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

func SuperuserPostDelete(pid, uid int64) (err error) {
	//0.查看管理员用户是否存在
	err = mysql.CheckSuperUserExit(uid)
	if !errors.Is(err, mysql.ErrorSuperuserExist) {
		return err
	}

	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed", zap.Error(err))
		return err
	}

	//1.查看管理员是否有权限
	if err = mysql.CheckSuperUserPower(uid, post.ChunkID); err != nil {
		return err
	}

	if post.Status != 3 {
		return mysql.ErrorNotPower
	}

	//2.1 先删除有关post的缓存
	err = redis.SuperuserDeletePost(pid, uid, post.ChunkID)
	if err != nil {
		return err
	}
	fmt.Println(11111111)

	//2.2 移除待删除区的缓存
	err = redis.DeletePostFromTobeDeleteCache(pid, post.ChunkID)
	if err != nil {
		return err
	}
	fmt.Println(22222)

	//2.3去mysql中删除记录
	if err := mysql.SuperuserDeletePost(pid); err != nil {
		return err
	}
	fmt.Println(33333)

	return nil
}
