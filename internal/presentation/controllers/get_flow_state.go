package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/luizrgf2/universal-flow/internal/core/usecases"
	"github.com/luizrgf2/universal-flow/internal/presentation/factories"
)

func GetFlowStateController(c *gin.Context) {
	flowID := c.Param("id")

	if flowID == "" {
		c.JSON(400, gin.H{"error": "flow id is required"})
		return
	}

	usecase := factories.GetFlowStateFactory()

	flow, err := usecase.Execute(usecases.GetFlowStateUseCaseInput{
		FlowID: flowID,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, flow)
}
