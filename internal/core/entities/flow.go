package entities

import (
	"fmt"

	"slices"

	"github.com/go-playground/validator/v10"

	"github.com/luizrgf2/universal-flow/internal/core/types"
)

type Flow struct {
	ID                  string           `json:"id" validate:"required,uuid"`
	FlowName            string           `json:"flowName" validate:"required,min=3,max=40"`
	UrlBaseServer       *string          `json:"url_base_server,omitempty" validate:"omitempty,url"`
	Status              types.FlowStatus `json:"status"`
	Nodes               []Node           `json:"nodes"`
	CurrentNode         *string          `json:"currentNode"`
	NextNode            *string          `json:"nextNode"`
	PreviousNode        *string          `json:"previousNode"`
	PreviousNodesRunned []string         `json:"previous_nodes_runned"`
}

var validFlowTransitions = map[types.FlowStatus][]types.FlowStatus{
	"pending":   {"running"},
	"running":   {"completed", "failed"},
	"completed": {},
	"failed":    {},
}

func (f *Flow) ChangeFlowStatus(statusToChange types.FlowStatus) error {
	if f.Status == statusToChange {
		return nil
	}

	allowed, ok := validFlowTransitions[f.Status]
	if !ok {
		return fmt.Errorf("current status '%s' is not a valid known state", f.Status)
	}

	for _, allowedStatus := range allowed {
		if statusToChange == allowedStatus {
			f.Status = statusToChange
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from '%s' to '%s'", f.Status, statusToChange)
}

func (f *Flow) SetCurrentNode(nodeID string) error {
	if f.CurrentNode != nil && *f.CurrentNode == nodeID {
		return fmt.Errorf("node %s is already the current node", nodeID)
	}
	if f.PreviousNode != nil && *f.PreviousNode == nodeID {
		return fmt.Errorf("node %s was the previous node", nodeID)
	}
	if slices.Contains(f.PreviousNodesRunned, nodeID) {
		return fmt.Errorf("node %s has already been run", nodeID)
	}

	f.PreviousNode = f.CurrentNode
	f.CurrentNode = &nodeID

	return nil
}

func (f *Flow) SetNextNode(nodeID string) error {
	if f.CurrentNode != nil && *f.CurrentNode == nodeID {
		return fmt.Errorf("next node %s cannot be the same as the current node", nodeID)
	}
	if f.PreviousNode != nil && *f.PreviousNode == nodeID {
		return fmt.Errorf("next node %s cannot be the same as the previous node", nodeID)
	}
	if slices.Contains(f.PreviousNodesRunned, nodeID) {
		return fmt.Errorf("node %s has already been run", nodeID)
	}

	f.NextNode = &nodeID

	return nil
}

func (f *Flow) SetPreviousNode(nodeID string) error {
	if f.CurrentNode != nil && *f.CurrentNode == nodeID {
		f.PreviousNode = f.CurrentNode
	}
	if f.NextNode != nil && *f.NextNode == nodeID {
		f.NextNode = nil
	}
	if slices.Contains(f.PreviousNodesRunned, nodeID) {
		return fmt.Errorf("node %s has already been run", nodeID)
	}

	f.PreviousNodesRunned = append(f.PreviousNodesRunned, nodeID)
	f.PreviousNode = &nodeID

	return nil
}

func validateFlow(flow *Flow) error {
	validator := validator.New()
	err := validator.Struct(flow)
	if err != nil {
		return err
	}
	return nil
}

func CreateFlow(id string, flowName string, nodes []Node, urlBaseServer *string) (*Flow, error) {
	status, err := types.CreateFlowStatus("pending")
	if err != nil {
		return nil, err
	}

	if len(nodes) < 3 {
		return nil, fmt.Errorf("a flow must have at least 3 nodes")
	}

	flow := &Flow{
		ID:                  id,
		FlowName:            flowName,
		Status:              status,
		Nodes:               nodes,
		UrlBaseServer:       urlBaseServer,
		CurrentNode:         nil,
		NextNode:            nil,
		PreviousNode:        nil,
		PreviousNodesRunned: []string{},
	}

	err = validateFlow(flow)

	if err != nil {
		return nil, err
	}

	return flow, nil
}
