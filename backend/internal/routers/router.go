package routers

import (
	"civ/internal/routers/groups"
	"civ/internal/routers/setup"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	Controllers := setup.NewControllers()
	api := router.Group("/api")
	groups.HelloRouters(api, *Controllers)
}
