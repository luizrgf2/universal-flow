package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/luizrgf2/universal-flow/internal/core/usecases"
	"github.com/luizrgf2/universal-flow/internal/presentation/factories"
)

func CreateFlowToRunController(c *gin.Context) {
	var createFlowToRunInput usecases.CreateFlowUseCaseInput

	if err := c.ShouldBindJSON(&createFlowToRunInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	usecase := factories.CreateFlowToRunFactory()
	err := usecase.Execute(createFlowToRunInput)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Flow created successfully"})
}
