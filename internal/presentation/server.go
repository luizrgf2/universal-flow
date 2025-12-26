package presentation

import "github.com/gin-gonic/gin"

func StartServer() *gin.RouterGroup {
	r := gin.Default()

	routeApi := r.Group("/api")

	routeApi.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
	return routeApi
}
