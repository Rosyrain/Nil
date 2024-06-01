package models

// ParamExamine 帖子审核信息
type ParamExamine struct {
	SuperuserID int64 `json:"superuser_id,string"`
	ChunkID     int64 `json:"chunk_id,string" binding:"required"`
	PostID      int64 `json:"post_id,string" binding:"required"`
	Status      int   `json:"status" binding:"oneof=0 1"`
	Direction   int   `json:"direction" binding:"oneof=1 2,required"`
}

// ParamChunk  板块创建信息
type ParamChunk struct {
	ChunkName    string `json:"chunk_name" binding:"required"`
	ChunkId      int64  `json:"chunk_id,string" binding:"required"`
	Introduction string `json:"introduction" binding:"required"`
}

type ParamSearch struct {
	SuperuserID int64  `json:"superuser_id,string" form:"superuser_id"`
	Page        int64  `json:"page" form:"page" example:"1"`       // 页码
	Size        int64  `json:"size" form:"size" example:"10"`      // 每页数量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
	ChunkID     int64  `json:"chunk_id" form:"chunk_id" example:"1"`
	Status      int    `json:"status" form:"status" example:"0"`
}
