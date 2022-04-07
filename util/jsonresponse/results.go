package jsonresponse

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	SuccessCode int = 0
	ErrorCode   int = -1
)

type Response struct {
	// 代码
	Code int `json:"code" example:"200"`
	// 消息
	Message string `json:"message"`
	// 数据集
	Data interface{} `json:"data"`
}

func Success(r *ghttp.Request, data interface{}) {
	Result(r, SuccessCode, "success", data)
}

func SuccessWithMessage(r *ghttp.Request, message string, data interface{}) {
	Result(r, SuccessCode, message, data)
}

func Failed(r *ghttp.Request, message string) {
	Result(r, ErrorCode, message, nil)
}

func FailedWithCode(r *ghttp.Request, code int, message string) {
	Result(r, code, message, nil)
}

func Result(r *ghttp.Request, code int, message string, data interface{}) {
	response := &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	r.SetParam("apiReturnRes", response)
	err := r.Response.WriteJson(response)
	if err != nil {
		g.Log().Error(r.Context(), err.Error())
	}
	r.Exit()
}
