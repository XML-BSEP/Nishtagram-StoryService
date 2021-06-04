package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	g := gin.Default()
	g.GET("ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})
	g.Run("localhost:8084")
}

