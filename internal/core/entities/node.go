package entities

import (
	"fmt"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/luizrgf2/universal-flow/internal/core/types"
)

type NodeState struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Node struct {
	ID           string           `json:"id" validate:"required,uuid"`
	Name         string           `json:"name" validate:"required,min=3,max=100"`
	ScriptPath   string           `json:"script_path" validate:"required"`
	Status       types.NodeStatus `json:"status"`
	State        NodeState        `json:"state"`
	Error        *string          `json:"error"`
	OutputNodes  []string         `json:"outputNodes" validate:"required,min=1"`
	SelectedNode *string          `json:"selectedNode"`
}

func validateNode(node *Node) error {
	validator := validator.New()
	err := validator.Struct(node)
	if err != nil {
		return err
	}
	return nil
}

func NewNodeInstance(id string, name string, scriptPath string, outputNodes []string) (*Node, error) {

	status, err := types.CreateNodeStatus("pending")
	if err != nil {
		return nil, err
	}

	node := &Node{
		ID:          id,
		Name:        name,
		ScriptPath:  scriptPath,
		Status:      status,
		OutputNodes: outputNodes,
	}

	err = validateNode(node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

var validTransitions = map[types.NodeStatus][]types.NodeStatus{
	"pending":   {"running"},
	"running":   {"completed", "failed"},
	"completed": {},
	"failed":    {},
}

func (n *Node) ChangeNodeStatus(statusToChange types.NodeStatus) error {
	if n.Status == statusToChange {
		return nil
	}

	allowed, ok := validTransitions[n.Status]
	if !ok {
		return fmt.Errorf("current status '%s' is not a valid known state", n.Status)
	}

	for _, allowedStatus := range allowed {
		if statusToChange == allowedStatus {
			n.Status = statusToChange
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from '%s' to '%s'", n.Status, statusToChange)
}

func (n *Node) ChangeSelectedNode(selectedNode string) error {
	if selectedNode == "" {
		n.SelectedNode = nil
		return nil
	}

	if slices.Contains(n.OutputNodes, selectedNode) {
		n.SelectedNode = &selectedNode
		return nil
	} else {
		return fmt.Errorf("Error not exists this node id: %s in ouputNodes", selectedNode)
	}
}

func (n *Node) ChangeOutput(outputData string) error {
	if outputData == "" {
		n.State.Output = ""
		return nil
	}
	n.State.Output = outputData
	return nil
}

func (n *Node) ChangeInput(inputData string) error {
	if inputData == "" {
		n.State.Input = ""
		return nil
	}
	n.State.Input = inputData
	return nil
}

func (n *Node) ChangeError(errorData string) error {
	if errorData == "" {
		n.Error = nil
		return nil
	}
	n.Error = &errorData
	return nil
}
