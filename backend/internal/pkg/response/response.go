package response

import (
	"civ/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"civ/config"
	"time"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Cost string      `json:"cost"`
}

type Response struct {
	httpCode int
	result   *Result
}

func Resp() *Response {
	return &Response{
		httpCode: http.StatusOK,
		result: &Result{
			Code: 0,
			Msg:  "",
			Data: nil,
			Cost: "",
		},
	}
}

func (r *Response) Fail(c *gin.Context, code int, msg string, data ...any) {
	r.SetCode(code)
	r.SetMessage(msg)
	if data != nil {
		r.WithData(data[0])
	}
	r.json(c)
}

// SetCode 设置返回code码
func (r *Response) SetCode(code int) *Response {
	r.result.Code = code
	return r
}

// SetHttpCode 设置http状态码
func (r *Response) SetHttpCode(code int) *Response {
	r.httpCode = code
	return r
}

type defaultRes struct {
	Result any `json:"result"`
}

// WithData 设置返回data数据
func (r *Response) WithData(data any) *Response {
	switch data.(type) {
	case string, int, bool:
		r.result.Data = &defaultRes{Result: data}
	default:
		r.result.Data = data
	}
	return r
}

// SetMessage 设置返回自定义错误消息
func (r *Response) SetMessage(message string) *Response {
	r.result.Msg = message
	return r
}

var ErrorText = errors.NewErrorText(config.GetConfig().System.Language)

// json 返回 gin 框架的 HandlerFunc
func (r *Response) json(c *gin.Context) {
	if r.result.Msg == "" {
		r.result.Msg = ErrorText.Text(r.result.Code)
	}
	r.result.Cost = time.Since(c.GetTime("requestStartTime")).String()
	c.AbortWithStatusJSON(r.httpCode, r.result)
}
