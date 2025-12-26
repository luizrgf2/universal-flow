package presentation

import (
	"github.com/gin-gonic/gin"
	"github.com/luizrgf2/universal-flow/internal/presentation/routes"
)

func StartServer() *gin.Engine {
	r := gin.Default()

	routeApi := r.Group("/api")

	routeApi.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.FlowStateRoutes(routeApi)
	return r
}
