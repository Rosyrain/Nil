package models

import "time"

//内存对齐概念

type Post struct {
	ID         int64     `json:"id,string" db:"post_id"`
	AuthorID   int64     `json:"author_id,string" db:"author_id"`
	ChunkID    int64     `json:"chunk_id,string" db:"chunk_id" binding:"required"`
	Status     int32     `json:"status" db:"status"` //是否审核通过
	Title      string    `json:"title" db:"title" binding:"required"`
	Content    string    `json:"content" db:"content" binding:"required"`
	VoteNum    int64     `json:"vote_number" db:"vote_number"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

type ApiPostDetail struct {
	AuthorName   string                `json:"author_name"`
	VoteNumber   int64                 `json:"vote_number"`
	*Post                              //嵌入式帖子结构体
	*ChunkDetail `json:"chunk_detail"` //嵌入帖子信息
}
