package routers

import (
	"civ/internal/routers/groups"
	"civ/internal/routers/setup"

	"github.com/gin-gonic/gin"
)

// SetupRouter registers API routes on the provided gin.Engine.
// It creates controller instances via setup.NewControllers(), mounts the
// "/api" route group on the given router, and registers application routes
// (currently groups.HelloRouters) onto that group.
func SetupRouter(router *gin.Engine) {
	Controllers := setup.NewControllers()
	api := router.Group("/api")
	groups.HelloRouters(api, *Controllers)
}
