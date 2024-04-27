package mysql

import (
	"nil/models"
	setting "nil/settings"
	"testing"
)

func init() {
	dbcfg := setting.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "jx20031002",
		DB:           "nil",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbcfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:       10,
		AuthorID: 123,
		ChunkID:  1,
		Title:    "test",
		Content:  "just a test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed,err%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success....")
}
