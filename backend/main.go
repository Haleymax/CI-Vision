package main

import (
	"civ/data"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db, err := data.IniDB()
	if db != nil {
		log.Println("数据库连接成功")
	}
	defer data.CloseDB()

	if err != nil {
		log.Println(err)
		panic(err)
	}
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8085")
}
