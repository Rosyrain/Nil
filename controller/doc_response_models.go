package controller

import "nil/models"

//	专门用来接口文档用到的model
//	因为我们的接口文档返回的数据格式是一致的，但是具体的Data类型不一致

type _ResponsePostList struct {
	Code    ResCode                 `json:"code"`    //业务相应状态码
	Message string                  `json:"message"` //提示信息
	Data    []*models.ApiPostDetail `json:"data"`    //数据
}
