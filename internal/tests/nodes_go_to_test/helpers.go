package nodes_go_to_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const baseURL = "http://localhost:8080/api/flow-state"

func GetNode(nodes []NodeInterface, nodeID string) *NodeInterface {
	for i := range nodes {
		if nodes[i].ID == nodeID {
			return &nodes[i]
		}
	}
	return nil
}

func GetFlowState() (*FlowInterface, error) {
	flowID := os.Getenv("FLOW_ID")
	url := fmt.Sprintf("%s/get-flow-state/%s", baseURL, flowID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var flowState FlowInterface
	err = json.NewDecoder(resp.Body).Decode(&flowState)
	if err != nil {
		return nil, err
	}
	return &flowState, nil
}

func FinishNode(nextNodeID string, output interface{}, errorMsg string) error {
	flowID := os.Getenv("FLOW_ID")
	nodeID := os.Getenv("NODE_ID")

	urlFinish := fmt.Sprintf("%s/finish-node?flowId=%s", baseURL, flowID)

	var outputJSON string
	if output != nil {
		outputBytes, err := json.Marshal(output)
		if err != nil {
			return err
		}
		outputJSON = string(outputBytes)
	}

	body := FinishInterface{
		FlowID:       flowID,
		NodeID:       nodeID,
		NextNodeID:   nextNodeID,
		NodeOutput:   outputJSON,
		ErrorMessage: errorMsg,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, urlFinish, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
