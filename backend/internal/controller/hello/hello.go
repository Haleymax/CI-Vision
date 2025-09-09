package hello

import (
	"civ/internal/controller"
	"civ/internal/service"
	"github.com/gin-gonic/gin"
)

type HelloController struct {
	controller.Api
}

func NewHelloController() *HelloController {
	return &HelloController{}
}

func (api HelloController) HelloGin(c *gin.Context) {
	result, err := service.NewHelloService().Hello()
	if err != nil {
		api.Err(c, err)
		return
	}
	api.Success(c, result)
	return
}
