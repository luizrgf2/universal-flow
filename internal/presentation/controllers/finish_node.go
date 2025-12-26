package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/luizrgf2/universal-flow/internal/core/usecases"
	"github.com/luizrgf2/universal-flow/internal/presentation/factories"
)

func FinishNodeController(c *gin.Context) {
	var finishNodeInput usecases.FinishNodeUseCaseInput
	if err := c.ShouldBindJSON(&finishNodeInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	flowID := c.Query("flowId")
	if flowID == "" {
		c.JSON(400, gin.H{"error": "flowId query param is required"})
		return
	}

	finishNodeInput.FlowID = flowID

	usecase := factories.FinishNodeFactory()
	err := usecase.Execute(finishNodeInput)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "node finished successfully"})
}
