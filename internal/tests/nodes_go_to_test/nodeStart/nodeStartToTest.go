package main

import (
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

	output := map[string]string{
		"fileName": "password.txt",
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
