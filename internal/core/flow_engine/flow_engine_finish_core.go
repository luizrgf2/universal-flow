package flowengine

import (
	"fmt"
	"slices"

	"github.com/luizrgf2/universal-flow/internal/core/entities"
	services_interfaces "github.com/luizrgf2/universal-flow/internal/core/services"
)

type FinishNodeInput struct {
	FlowID       string  `json:"flow_id"`
	NodeID       string  `json:"node_id"`
	NextNodeID   *string `json:"next_node_id"`
	ErrorMessage *string `json:"error_message"`
	NodeOutput   *string `json:"node_output"`
}

type FlowEngineFinishCore struct {
	flowStateManagerService services_interfaces.FlowStateManagerService
}

func NewFlowEngineFinish(flowStateManagerService services_interfaces.FlowStateManagerService) *FlowEngineFinishCore {
	return &FlowEngineFinishCore{
		flowStateManagerService: flowStateManagerService,
	}
}

func (fe *FlowEngineFinishCore) findNode(flow *entities.Flow, nodeID string) (*entities.Node, int, error) {
	for i, node := range flow.Nodes {
		if node.ID == nodeID {
			return &node, i, nil
		}
	}
	return nil, -1, fmt.Errorf("node with id %s not found in flow", nodeID)
}

func (fe *FlowEngineFinishCore) setError(node *entities.Node, nodeIndex int, flow *entities.Flow, errorMessage string) error {
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

	flow.Nodes[nodeIndex] = *node

	return nil
}

func (fe *FlowEngineFinishCore) setNextNode(node *entities.Node, nodeIndex int, flow *entities.Flow, nextNodeID string) error {
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

	flow.Nodes[nodeIndex] = *node

	return nil
}

func (fe *FlowEngineFinishCore) changeOuputOfCurrentNode(input FinishNodeInput, flow *entities.Flow) error {
	if input.NodeOutput == nil {
		return nil
	}
	node, nodeIndex, err := fe.findNode(flow, input.NodeID)
	if err != nil {
		return err
	}
	if input.NodeOutput != nil {
		err = node.ChangeOutput(*input.NodeOutput)
		if err != nil {
			return err
		}
		flow.Nodes[nodeIndex] = *node
	}
	return nil

}

func (fe *FlowEngineFinishCore) changeInputOfNextNode(input FinishNodeInput, flow *entities.Flow) error {
	if input.NextNodeID == nil {
		return nil
	}
	nextNode, nextNodeIndex, err := fe.findNode(flow, *input.NextNodeID)
	if err != nil {
		return err
	}
	if input.NodeOutput != nil {
		err = nextNode.ChangeInput(*input.NodeOutput)
		if err != nil {
			return err
		}
		flow.Nodes[nextNodeIndex] = *nextNode
	}
	return nil
}

func (fe *FlowEngineFinishCore) setPreviousNode(flow *entities.Flow, nodeId string) error {
	flow.SetPreviousNode(nodeId)
	err := fe.flowStateManagerService.UpdateFlow(flow)
	if err != nil {
		return err
	}
	return nil
}

func (fe *FlowEngineFinishCore) runFlow(flow *entities.Flow) error {
	flowEngineService := NewFlowEngine(fe.flowStateManagerService)
	err := flowEngineService.RunFlow(flow)
	if err != nil {
		return err
	}
	if flow.NextNode != nil {
		err := flowEngineService.RunFlow(flow)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fe *FlowEngineFinishCore) FinishNode(input FinishNodeInput) error {
	flow, err := fe.flowStateManagerService.GetFlowState(input.FlowID)
	if err != nil {
		return err
	}

	node, nodeIndex, err := fe.findNode(flow, input.NodeID)
	if err != nil {
		return err
	}

	fe.setPreviousNode(flow, node.ID)

	if input.ErrorMessage != nil {
		fe.setError(node, nodeIndex, flow, *input.ErrorMessage)
	} else if input.NextNodeID != nil {
		fe.setNextNode(node, nodeIndex, flow, *input.NextNodeID)
	} else {
		return nil
	}

	err = fe.changeOuputOfCurrentNode(input, flow)
	if err != nil {
		return err
	}

	err = fe.changeInputOfNextNode(input, flow)
	if err != nil {
		return err
	}

	err = fe.flowStateManagerService.UpdateFlow(flow)
	if err != nil {
		return err
	}

	err = fe.runFlow(flow)
	if err != nil {
		return err
	}

	return nil
}
