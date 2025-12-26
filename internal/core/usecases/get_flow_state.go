package usecases

import (
	"github.com/luizrgf2/universal-flow/internal/core/entities"
	services_interfaces "github.com/luizrgf2/universal-flow/internal/core/services"
)

type GetFlowStateUseCaseInput struct {
	FlowID string
}

type GetFlowStateUseCase struct {
	FlowStateManagerService services_interfaces.FlowStateManagerService
}

func NewGetFlowStateUseCase(flowStateManagerService services_interfaces.FlowStateManagerService) *GetFlowStateUseCase {
	return &GetFlowStateUseCase{
		FlowStateManagerService: flowStateManagerService,
	}
}

func (uc *GetFlowStateUseCase) Execute(input GetFlowStateUseCaseInput) (*entities.Flow, error) {
	flow, err := uc.FlowStateManagerService.GetFlowState(input.FlowID)
	if err != nil {
		return nil, err
	}
	return flow, nil
}
