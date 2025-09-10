package server

import (
	"civ/config"
	"civ/internal/routers"
	"log"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()

	config := config.GetConfig()
	routers.SetupRouter(r)
	go func() {
		if err := r.Run(config.System.Host + ":" + string(rune(config.System.Port))); err != nil {
			log.Fatal("Server Run Failed:", err)
		}
	}()
}
