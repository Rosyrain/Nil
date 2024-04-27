package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
			"chunk_id":"1",
			"title":"test",
			"content":"just a test"
		}`

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	//判断响应内容是不是按预期返回了需要登陆的错误
	//assert.Equal(t, "pong", w.Body.String())

	//方法一：判断响应的内容是不是包含指定字符串
	//assert.Contains(t, w.Body.String(), "需要登录")

	//方法二: 将响应的内容反序列化到ResponseData，然后判断字段与预期是否一致
	res := new(ResponseData)
	if err := sonic.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("sonic.Unmarshsal w.Body failed,err:%v\n", err)
	}
	assert.Equal(t, res.Code, CodeNeedLogin)

}
