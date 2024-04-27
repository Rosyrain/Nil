package models

import "time"

type Chunk struct {
	ID   int64  `json:"id,string" db:"chunk_id"`
	Name string `json:"name" db:"chunk_name"`
}

type ChunkDetail struct {
	ID           int64     `json:"id,string" db:"chunk_id"`
	Name         string    `json:"name" db:"chunk_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
