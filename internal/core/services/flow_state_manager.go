package services_interfaces

import "github.com/luizrgf2/universal-flow/internal/core/entities"

type FlowStateManagerService interface {
	RunNewFlow(flow entities.Flow) error
}
