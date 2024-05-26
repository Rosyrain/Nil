package models

type Superuser struct {
	UserID   int64  `db:"user_id" json:"user_id,string"`
	ChunkID  int64  `db:"chunk_id" json:"chunk_id,string"`
	Username string `db:"username"`
	Password string `db:"password"`
	Token    string `json:"token"`
	Birthday string `db:"create_time" json:"birthday"`
}
