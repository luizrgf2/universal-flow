package usecases

import (
	"fmt"
	"slices"

	"github.com/luizrgf2/universal-flow/internal/core/entities"
	flowengine "github.com/luizrgf2/universal-flow/internal/core/flow_engine"
	services_interfaces "github.com/luizrgf2/universal-flow/internal/core/services"
)

type FinishNodeUseCaseInput struct {
	FlowID       string  `json:"flow_id"`
	NodeID       string  `json:"node_id"`
	NextNodeID   *string `json:"next_node_id"`
	ErrorMessage *string `json:"error_message"`
	NodeOutput   *string `json:"node_output"`
	NodeInput    *string `json:"node_input"`
}

type FinishNodeUseCase struct {
	FlowStateManagerService services_interfaces.FlowStateManagerService
}

func NewFinishNodeUseCase(flowStateManagerService services_interfaces.FlowStateManagerService) *FinishNodeUseCase {
	return &FinishNodeUseCase{
		FlowStateManagerService: flowStateManagerService,
	}
}

func (uc *FinishNodeUseCase) findNode(flow *entities.Flow, nodeID string) (*entities.Node, int, error) {
	for i, node := range flow.Nodes {
		if node.ID == nodeID {
			return &node, i, nil
		}
	}
	return nil, -1, fmt.Errorf("node with id %s not found in flow", nodeID)
}

func (uc *FinishNodeUseCase) setError(node *entities.Node, flow *entities.Flow, errorMessage string) error {
	err := node.ChangeError(errorMessage)
	if err != nil {
		return err
	}
	err = node.ChangeNodeStatus("failed")
	if err != nil {
		return err
	}
	err = flow.ChangeFlowStatus("failed")
	if err != nil {
		return err
	}
	return nil
}

func (uc *FinishNodeUseCase) setNextNode(node *entities.Node, flow *entities.Flow, nextNodeID string, nodeInput *string) error {
	if !slices.Contains(node.OutputNodes, nextNodeID) {
		return fmt.Errorf("next node is not in the list of output nodes")
	}

	err := node.ChangeSelectedNode(nextNodeID)
	if err != nil {
		return err
	}

	err = node.ChangeNodeStatus("completed")
	if err != nil {
		return err
	}

	err = flow.SetNextNode(nextNodeID)
	if err != nil {
		return err
	}

	nextNode, nextNodeIndex, err := uc.findNode(flow, nextNodeID)
	if err != nil {
		return err
	}

	if nodeInput != nil {
		err = nextNode.ChangeInput(*nodeInput)
		if err != nil {
			return err
		}
	}
	flow.Nodes[nextNodeIndex] = *nextNode
	return nil
}

func (uc *FinishNodeUseCase) Execute(input FinishNodeUseCaseInput) error {
	flow, err := uc.FlowStateManagerService.GetFlowState(input.FlowID)
	if err != nil {
		return err
	}

	node, nodeIndex, err := uc.findNode(flow, input.NodeID)
	if err != nil {
		return err
	}

	if input.ErrorMessage != nil {
		uc.setError(node, flow, *input.ErrorMessage)
	} else if input.NextNodeID != nil {
		uc.setNextNode(node, flow, *input.NextNodeID, input.NodeInput)
	} else {
		return fmt.Errorf("you must provide a next node or an error")
	}

	if input.NodeOutput != nil {
		err = node.ChangeOutput(*input.NodeOutput)
		if err != nil {
			return err
		}
	}

	flow.Nodes[nodeIndex] = *node

	err = uc.FlowStateManagerService.UpdateFlow(flow)
	if err != nil {
		return err
	}

	flowengine := flowengine.NewFlowEngine(uc.FlowStateManagerService)

	if flow.NextNode != nil {
		err = flowengine.RunFlow(flow)
		if err != nil {
			return err
		}
	}

	return nil
}
