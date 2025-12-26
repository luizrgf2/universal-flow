package usecases

import (
	flowengine "github.com/luizrgf2/universal-flow/internal/core/flow_engine"
	services_interfaces "github.com/luizrgf2/universal-flow/internal/core/services"
)

type FinishNodeUseCaseInput struct {
	FlowID       string  `json:"flow_id"`
	NodeID       string  `json:"node_id"`
	NextNodeID   *string `json:"next_node_id"`
	ErrorMessage *string `json:"error_message"`
	NodeOutput   *string `json:"node_output"`
}

type FinishNodeUseCase struct {
	FlowStateManagerService services_interfaces.FlowStateManagerService
}

func NewFinishNodeUseCase(flowStateManagerService services_interfaces.FlowStateManagerService) *FinishNodeUseCase {
	return &FinishNodeUseCase{
		FlowStateManagerService: flowStateManagerService,
	}
}

func (uc *FinishNodeUseCase) Execute(input FinishNodeUseCaseInput) error {
	flowEngineFinishService := flowengine.NewFlowEngineFinish(uc.FlowStateManagerService)

	inputData := flowengine.FinishNodeInput{
		FlowID:       input.FlowID,
		NodeID:       input.NodeID,
		NextNodeID:   input.NextNodeID,
		ErrorMessage: input.ErrorMessage,
		NodeOutput:   input.NodeOutput,
	}

	return flowEngineFinishService.FinishNode(inputData)

}
