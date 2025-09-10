package groups

import (
	"civ/internal/routers/setup"

	"github.com/gin-gonic/gin"
)

func HelloRouters(router *gin.RouterGroup, controller setup.Controllers) {
	router.GET("/hello", controller.HelloController.HelloGin)
}
