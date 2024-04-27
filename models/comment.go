package models

import "time"

type Comment struct {
	CommentID  int64     `json:"comment_id,string" db:"comment_id"`
	AuthorID   int64     `json:"author_id,string" db:"author_id"`
	PostID     int64     `json:"post_id,string" db:"post_id"`
	Content    string    `json:"content" db:"content"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

type ApiCommentDetail struct {
	AuthorName string        `json:"author_name"`
	VoteNumber int64         `json:"vote_number"`
	*Comment                 //嵌入式帖子结构体
	*Post      `json:"post"` //嵌入帖子信息
}

type CommentList struct {
	*ApiCommentDetail
	SubComments []*ApiCommentDetail
}
