package mysql

import (
	"database/sql"
	"go.uber.org/zap"
	"nil/models"
)

func CheckChunkExist(chunkname string) (err error) {
	sqlStr := `select count(chunk_id) from chunk where chunk_name = ?`

	var count int
	if err = db.Get(&count, sqlStr, chunkname); err != nil {
		return err
	}
	if count > 0 {
		return ErrorChunkExist
	}
	return ErrorChunkNotExist
}

// 插入新板块信息到Chunk表中
func InsertChunk(c *models.ChunkDetail) (err error) {
	sqlStr := `insert into chunk (chunk_id,chunk_name,introduction) values (?,?,?)`

	_, err = db.Exec(sqlStr, c.ID, c.Name, c.Introduction)
	return err
}

func GetChunkList() (ChunkList []*models.Chunk, err error) {
	sqlStr := `select chunk_id,chunk_name from chunk`
	if err := db.Select(&ChunkList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

// GetChunkDetailByID  依据ID查询社区详情
func GetChunkDetailByID(id int64) (chunk *models.ChunkDetail, err error) {
	chunk = new(models.ChunkDetail)
	sqlStr := `select chunk_id,chunk_name,introduction,create_time from chunk where chunk_id=?`
	if err := db.Get(chunk, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return
}
