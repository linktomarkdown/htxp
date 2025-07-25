package htxp

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Body 定义了API响应的标准格式
type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 处理成功响应
func Success(w http.ResponseWriter, data interface{}) {
	Response(w, data, nil, 200)
}

// Error 处理错误响应（默认 code -1）
func Error(w http.ResponseWriter, err error) {
	Response(w, nil, err, -1)
}

// ErrorWithCode 支持自定义 code
func ErrorWithCode(w http.ResponseWriter, err error, code int) {
	Response(w, nil, err, code)
}

// Response 统一处理HTTP响应
func Response(w http.ResponseWriter, resp interface{}, err error, code int) {
	var body Body
	if err != nil {
		body.Code = code
		body.Msg = err.Error()
	} else {
		body.Code = code
		body.Msg = "success"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}
