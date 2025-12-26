package flowstateroutestests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/luizrgf2/universal-flow/internal/core/usecases"
	"github.com/luizrgf2/universal-flow/internal/presentation"
	"github.com/stretchr/testify/assert"
)

func createNodesToTest() []usecases.CreateFlowUseCaseNodes {

	nodes := []usecases.CreateFlowUseCaseNodes{}

	nodes = append(nodes, usecases.CreateFlowUseCaseNodes{
		ID:         uuid.NewString(),
		Name:       "Test1",
		ScriptPath: "node /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_js_to_test/nodeStartToTest.js",
		OutputNode: []string{"dsfsdf"},
	})
	nodes = append(nodes, usecases.CreateFlowUseCaseNodes{
		ID:         uuid.NewString(),
		Name:       "Test2",
		ScriptPath: "node /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_js_to_test/nodeStartToTest.js",
		OutputNode: []string{"dsfsdf"},
	})
	nodes = append(nodes, usecases.CreateFlowUseCaseNodes{
		ID:         uuid.NewString(),
		Name:       "Test3",
		ScriptPath: "node /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_js_to_test/nodeStartToTest.js",
		OutputNode: []string{"dsfsdf"},
	})

	return nodes

}

func TestToCreateFlowToRun(t *testing.T) {
	route := presentation.StartServer()
	w := httptest.NewRecorder()

	nodesToBody := createNodesToTest()

	body := usecases.CreateFlowUseCaseInput{
		ID:    uuid.NewString(),
		Nodes: nodesToBody,
		Name:  "flowTest",
	}

	bodyJson, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/flow-state/create-flow-to-run", bytes.NewBuffer(bodyJson))
	route.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

}
