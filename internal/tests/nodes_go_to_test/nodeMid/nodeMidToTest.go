package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/luizrgf2/universal-flow/internal/tests/nodes_go_to_test"
)

func main() {
	flowState, err := nodes_go_to_test.GetFlowState()
	if err != nil {
		log.Fatal(err)
	}

	nodeID := os.Getenv("NODE_ID")
	node := nodes_go_to_test.GetNode(flowState.Nodes, nodeID)

	if node == nil {
		nodes_go_to_test.FinishNode("", nil, "Node not found")
		log.Fatal("Node not found")
	}

	var inputState map[string]interface{}
	if input, ok := node.State["input"]; ok && input != "" {
		err := json.Unmarshal([]byte(input), &inputState)
		if err != nil {
			nodes_go_to_test.FinishNode(node.OutputNodes[0], nil, "Error parsing input state")
			log.Fatal(err)
		}
	}

	fileName, ok := inputState["fileName"].(string)
	if !ok {
		nodes_go_to_test.FinishNode(node.OutputNodes[0], nil, "fileName not found in input state")
		log.Fatal("fileName not found in input state")
	}

	err = os.WriteFile(fileName, []byte("Hello Password"), 0644)
	if err != nil {
		nodes_go_to_test.FinishNode(node.OutputNodes[0], nil, "Error writing file")
		log.Fatal(err)
	}

	output := map[string]string{
		"password": "Hello Password",
	}

	var nextNodeID string
	if len(node.OutputNodes) > 0 {
		nextNodeID = node.OutputNodes[0]
	}

	err = nodes_go_to_test.FinishNode(nextNodeID, output, "")
	if err != nil {
		log.Fatal(err)
	}
}
