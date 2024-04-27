package logic

import (
	"errors"
	"nil/dao/mysql"
	"nil/models"
	snowflake "nil/pkg/snowflask"
)

func CreateChunk(p *models.ParamChunk) error {
	//检查板块是否已存在
	if err := mysql.CheckChunkExist(p.ChunkName); err != nil {
		if !errors.Is(err, mysql.ErrorChunkNotExist) {
			//if errors.Is(err,mysql.ErrorChunkExist){
			//	return err
			//}else{
			//	return err
			//}
			return err
		}
	}

	//获取板块id
	chunk_id := snowflake.GenID()

	//整合Chunk需要的信息
	c := models.ChunkDetail{
		ID:           chunk_id,
		Name:         p.ChunkName,
		Introduction: p.Introduction,
	}

	//创建板块信息
	return mysql.InsertChunk(&c)
}

// GetChunkList  查到所有的社区信息并返回
func GetChunkList() ([]*models.Chunk, error) {
	//查找数据库 查到所有的chunk 并返回
	return mysql.GetChunkList()

}

// GetChunkDetail  查询社区详情
func GetChunkDetail(id int64) (*models.ChunkDetail, error) {
	//查找数据库 查到所有的community 并返回
	return mysql.GetChunkDetailByID(id)

}
