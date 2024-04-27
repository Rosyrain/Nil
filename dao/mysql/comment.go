package mysql

import (
	"github.com/jmoiron/sqlx"
	"nil/models"
	"strings"
)

func CreateComment(p *models.Comment) (err error) {
	sqlStr := `insert into comment(comment_id,author_id,post_id,content) values (?,?,?,?)`
	_, err = db.Exec(sqlStr, p.CommentID, p.AuthorID, p.PostID, p.Content)
	return err
}

// GetCommentByID  根据ID查询一条评论的详情信息
func GetCommentByID(cid int64) (comment *models.Comment, err error) {
	comment = new(models.Comment)
	sqlStr := `select 
    comment_id,post_id,author_i,contentd,create_time 
	from comment where comment_id=?`
	err = db.Get(comment, sqlStr, cid)
	return
}

// GetCommentListByIDs  根据给点的id列表查询帖子数据
func GetCommentListByIDs(ids []string) (data []*models.Comment, err error) {
	sqlStr := `select comment_id,post_id,author_id,content,create_time
			from comment
			where comment_id in (?)
			order by FIND_IN_SET(comment_id,?)
			    `

	//https://www.liwenzhou.com/posts/Go/sqlx/
	//此处的query是新的查询语句列表，args是参数列表
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	//将query重绑至db，然后查询
	query = db.Rebind(query)
	err = db.Select(&data, query, args...) //!!!!别忘记最后的 ...
	return
}
