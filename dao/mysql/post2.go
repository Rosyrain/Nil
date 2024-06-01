package mysql

import (
	"nil/models"
	"strconv"
)

func CheckSuperUserPower(uid, cid int64) error {
	sqlStr := `select chunk_id from superuser where user_id=?`
	var chunkID int64
	err := db.Get(&chunkID, sqlStr, uid)
	if err != nil {
		return err
	}

	if chunkID != cid && chunkID != 0 {
		return ErrorNotPower
	}
	return nil
}

func UpdatePostStatus(p *models.ParamExamine) (err error) {
	sqlStr := `update post set status=? where post_id=?`
	_, err = db.Exec(sqlStr, p.Direction, p.PostID)
	return
}

func SuperuserDeletePost(pid int64) error {
	id := strconv.Itoa(int(pid))
	sqlStr := `delete from post where post_id=?`
	_, err := db.Exec(sqlStr, id)
	return err
}
