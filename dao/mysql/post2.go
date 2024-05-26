package mysql

import "nil/models"

func CheckSuperUserPower(p *models.ParamExamine) error {
	sqlStr := `select chunk_id from superuser where user_id=?`
	var chunkID int64
	err := db.Get(&chunkID, sqlStr, p.UserID)
	if err != nil {
		return err
	}

	if chunkID != p.ChunkID && chunkID != 0 {
		return ErrorNotPower
	}
	return nil
}

func UpdatePostStatus(p *models.ParamExamine) (err error) {
	sqlStr := `update post set status=? where post_id=?`
	_, err = db.Exec(sqlStr, p.Direction, p.PostID)
	return
}
