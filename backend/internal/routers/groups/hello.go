package groups

import (
	"civ/internal/routers/setup"

	"github.com/gin-gonic/gin"
)

// HelloRouters registers the GET /hello route on the given router group and binds it to controller.HelloController.HelloGin.
func HelloRouters(router *gin.RouterGroup, controller setup.Controllers) {
	router.GET("/hello", controller.HelloController.HelloGin)
}
