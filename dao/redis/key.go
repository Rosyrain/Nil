package redis

//redis key

//redis key尽量使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix                = "nil:"
	KeyPostTime              = "post:time:"        //zset;帖子及发帖时间
	KeyPostScore             = "post:score:"       //zset;帖子及投票分数
	KeyPostVotedPrefix       = "post:voted:"       //zset;记录用户投票类型；参数是post_id
	KeyActivateCaptcha       = "user:captcha:"     //zset:记录用户的验证码
	KeyChunkPrefix           = "chunk:"            //set:保存每个分区下的id
	KeyPostCommentPrefix     = "post:comment:"     //set:保存每个帖子下的主评论id
	KeyCommentTime           = "comment:time:"     //zset;评论及发帖时间
	KeyCommentScore          = "comment:score:"    //zset;评论及评论分数
	KeyCommentPrefix         = "comment:"          //set;保存每个主评论下的子评论
	KeySubCommentTime        = "subcomment:time:"  //zset;子评论及发帖时间
	KeySubCommentScore       = "subcomment:score:" //zset;主评论及评论分数
	KeyCommentVotedPrefix    = "comment:voted:"    //zset;记录用户投票类型；参数是post_id
	KeyUserPostPrefix        = "user:post:"        //set;记录用户发布的帖子
	KeyUserCommentPrefix     = "user:comment:"     //set;记录用户发布的评论(不分主次)
	KeySubCommentVotedPrefix = "subcomment:voted:" //zset;记录用户投票类型；参数是post_id
	KeyUserFocusPrefix       = "user:focus:"       //zset;记录用户的关注用户
	KeyUserHistoryPrefix     = "user:history:"     //zset;记录用户浏览记录

	KeyChunkNormalPrefix = "chunk:normal:" //set:保存每个分区下的正常的帖子id
	KeyNormalPost        = "normal:"       //set:用于保存全部的正常的帖子id

	KeyChunkToBeReviewPrefix = "chunk:tobereview:" //set:保存每个分区下的待审查的帖子id

	KeyChunkToBeDeletePrefix = "chunk:tobedelete:" //set:保存每个分区下的待删除的帖子id

)

// GetRedisKey  返回redis key加上前缀
func GetRedisKey(key string) string {
	return KeyPrefix + key
}
