package services_interfaces

import "github.com/luizrgf2/universal-flow/internal/core/entities"

type FlowStateManagerService interface {
	CreateFlow(flow *entities.Flow) error
	UpdateFlow(flow *entities.Flow) error
	GetFlowState(flowId string) (*entities.Flow, error)
}
