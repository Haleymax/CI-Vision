package controller

import (
	"civ/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	r "civ/internal/pkg/response"
	"net/http"
)

type Api struct {
	errors.Error
}

// Success 业务成功响应
func (api *Api) Success(c *gin.Context, data ...any) {
	response := r.Resp()
	if data == nil {
		response.WithDataSuccess(c, data[0])
		return
	}
	response.Success(c)
}

// FailCode 业务失败响应
func (api *Api) FailCode(c *gin.Context, code int, data ...any) {
	response := r.Resp()
	if data == nil {
		response.WithData(data[0]).FailCode(c, code)
		return
	}
	response.FailCode(c, code)
}

// Fail 业务失败响应
func (api *Api) Fail(c *gin.Context, code int, message string, data ...any) {
	response := r.Resp()
	if data == nil {
		response.WithData(data[0]).FailCode(c, code, message)
		return
	}
	response.FailCode(c, code, message)
}

func (api *Api) Err(c *gin.Context, e error) {
	businessError, err := api.AsBusinessError(e)
	if err != nil {
		api.FailCode(c, http.StatusInternalServerError, err.Error())
		return
	}

	api.Fail(c, businessError.GetCode(), businessError.GetMessage())
}
