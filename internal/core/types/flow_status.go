package types

import "errors"

type FlowStatus string

func validateFlowStatus(status string) error {
	if status != "pending" && status != "completed" && status != "running" && status != "failed" {
		return errors.New("invalid flow status: must be one of 'pending', 'completed', 'running', or 'failed'")
	}
	return nil
}

func CreateFlowStatus(status string) (FlowStatus, error) {
	err := validateFlowStatus(status)
	if err != nil {
		return "", err
	}
	return FlowStatus(status), nil
}
