package usecases

import (
	"github.com/luizrgf2/universal-flow/internal/core/entities"
	flowengine "github.com/luizrgf2/universal-flow/internal/core/flow_engine"
	services_interfaces "github.com/luizrgf2/universal-flow/internal/core/services"
)

type CreateFlowUseCaseNodes struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	ScriptPath string   `json:"script_path"`
	OutputNode []string `json:"output_node"`
}

type CreateFlowUseCaseInput struct {
	ID    string                   `json:"id"`
	Name  string                   `json:"name"`
	Nodes []CreateFlowUseCaseNodes `json:"nodes"`
}

type CreateFlowToRunUseCase struct {
	flowStateManagerService services_interfaces.FlowStateManagerService
}

func MakeCreateFlowToRun(flowStateManagerService services_interfaces.FlowStateManagerService) *CreateFlowToRunUseCase {
	return &CreateFlowToRunUseCase{flowStateManagerService}
}

func (uc *CreateFlowToRunUseCase) createNodesWithInput(nodesInput *[]CreateFlowUseCaseNodes) (*[]entities.Node, error) {
	nodes := []entities.Node{}

	for _, node := range *nodesInput {
		nodeToCreate, err := entities.CreateNode(node.ID, node.Name, node.ScriptPath, node.OutputNode)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, *nodeToCreate)
	}
	return &nodes, nil
}

func (uc *CreateFlowToRunUseCase) Execute(input CreateFlowUseCaseInput) error {

	nodes, err := uc.createNodesWithInput(&input.Nodes)
	if err != nil {
		return err
	}

	flow, err := entities.CreateFlow(input.ID, input.Name, *nodes)
	if err != nil {
		return err
	}

	err = uc.flowStateManagerService.CreateFlow(flow)
	if err != nil {
		return err
	}

	flowengine := flowengine.NewFlowEngine(uc.flowStateManagerService)

	err = flowengine.RunFlow(flow)
	if err != nil {
		return err
	}

	return nil

}
