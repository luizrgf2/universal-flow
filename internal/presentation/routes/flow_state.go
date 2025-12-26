package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luizrgf2/universal-flow/internal/presentation/controllers"
)

func FlowStateRoutes(routes *gin.RouterGroup) {

	routeGroup := routes.Group("/flow-state")
	routeGroup.POST("/create-flow-to-run", controllers.CreateFlowToRunController)

}
