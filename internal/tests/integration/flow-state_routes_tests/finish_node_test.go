package flowstateroutestests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/luizrgf2/universal-flow/internal/core/entities"
	"github.com/luizrgf2/universal-flow/internal/core/usecases"
	"github.com/luizrgf2/universal-flow/internal/presentation"
	"github.com/stretchr/testify/assert"
)

func TestFinishNode(t *testing.T) {
	route := presentation.StartServer()
	w := httptest.NewRecorder()

	// 1. Create a flow
	nodesToBody := []usecases.CreateFlowUseCaseNodes{
		{ID: uuid.NewString(), Name: "Start", ScriptPath: "bun /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_js_to_test/nodeStartToTest.ts", OutputNode: []string{uuid.NewString()}},
		{ID: uuid.NewString(), Name: "Middle", ScriptPath: "bun /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_js_to_test/nodeMidToTest.ts", OutputNode: []string{uuid.NewString()}},
		{ID: uuid.NewString(), Name: "End", ScriptPath: "bun /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_js_to_test/nodeEndToTest.ts", OutputNode: []string{}},
	}
	nodesToBody[0].OutputNode[0] = nodesToBody[1].ID
	nodesToBody[1].OutputNode[0] = nodesToBody[2].ID

	onlyCreateFlow := true

	createFlowBody := usecases.CreateFlowUseCaseInput{
		ID:             uuid.NewString(),
		Nodes:          nodesToBody,
		Name:           "flowTest",
		OnlyCreateFlow: &onlyCreateFlow,
	}

	bodyJson, _ := json.Marshal(createFlowBody)
	req, _ := http.NewRequest("POST", "/api/flow-state/create-flow-to-run", bytes.NewBuffer(bodyJson))
	route.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 2. Get the flow to know the current node
	w = httptest.NewRecorder()
	getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/flow-state/get-flow-state/%s", createFlowBody.ID), nil)
	route.ServeHTTP(w, getReq)
	assert.Equal(t, 200, w.Code)

	var flowResponse entities.Flow
	json.Unmarshal(w.Body.Bytes(), &flowResponse)

	currentNodeID := *flowResponse.CurrentNode
	var nextNodeID string
	for _, node := range flowResponse.Nodes {
		if node.ID == currentNodeID {
			nextNodeID = node.OutputNodes[0]
			break
		}
	}

	// 3. Finish the current node
	w = httptest.NewRecorder()
	finishNodeBody := usecases.FinishNodeUseCaseInput{
		NodeID:     currentNodeID,
		NextNodeID: &nextNodeID,
	}
	bodyJson, _ = json.Marshal(finishNodeBody)
	finishReq, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/flow-state/finish-node?flowId=%s", createFlowBody.ID), bytes.NewBuffer(bodyJson))
	route.ServeHTTP(w, finishReq)
	assert.Equal(t, 200, w.Code)

	// 4. Get the flow again to check the new state
	w = httptest.NewRecorder()
	getReq, _ = http.NewRequest("GET", fmt.Sprintf("/api/flow-state/get-flow-state/%s", createFlowBody.ID), nil)
	route.ServeHTTP(w, getReq)
	assert.Equal(t, 200, w.Code)

	var finalFlowResponse entities.Flow
	json.Unmarshal(w.Body.Bytes(), &finalFlowResponse)

	// Assertions
	assert.Equal(t, nextNodeID, *finalFlowResponse.CurrentNode)

	for _, node := range finalFlowResponse.Nodes {
		if node.ID == currentNodeID {
			assert.Equal(t, "completed", string(node.Status))
			assert.Equal(t, nextNodeID, *node.SelectedNode)
		}
		if node.ID == nextNodeID {
			assert.Equal(t, "running", string(node.Status))
		}
	}
}
