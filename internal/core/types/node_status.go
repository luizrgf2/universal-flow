package types

import "errors"

type NodeStatus string

func validateStatus(status string) error {
	if status != "pending" && status != "completed" && status != "running" && status != "failed" {
		return errors.New("invalid node status: must be one of 'pending', 'completed', 'running', or 'failed'")
	}
	return nil
}

func CreateNodeStatus(status string) (NodeStatus, error) {
	err := validateStatus(status)
	if err != nil {
		return "", err
	}
	return NodeStatus(status), nil
}
