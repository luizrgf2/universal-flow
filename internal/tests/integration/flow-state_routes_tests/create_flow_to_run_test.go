package flowstateroutestests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/luizrgf2/universal-flow/internal/core/usecases"
	"github.com/luizrgf2/universal-flow/internal/presentation"
	"github.com/stretchr/testify/assert"
)

var baseUrl = "http://localhost:8080/api/flow-state"

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

func createNodesInGoToTeste() []usecases.CreateFlowUseCaseNodes {
	startID := uuid.NewString()
	midID := uuid.NewString()
	endID := uuid.NewString()

	return []usecases.CreateFlowUseCaseNodes{
		{
			ID:         startID,
			Name:       "start",
			ScriptPath: "go run /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_go_to_test/nodeStart/nodeStartToTest.go",
			OutputNode: []string{midID},
		},
		{
			ID:         midID,
			Name:       "mid",
			ScriptPath: "go run /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_go_to_test/nodeMid/nodeMidToTest.go",
			OutputNode: []string{endID},
		},
		{
			ID:         endID,
			Name:       "end",
			ScriptPath: "go run /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_go_to_test/nodeEnd/nodeEndToTest.go",
			OutputNode: []string{},
		},
	}
}

func createNodesInBunToTeste() []usecases.CreateFlowUseCaseNodes {
	startID := uuid.NewString()
	midID := uuid.NewString()
	endID := uuid.NewString()

	return []usecases.CreateFlowUseCaseNodes{
		{
			ID:         startID,
			Name:       "start",
			ScriptPath: "bun /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_js_to_test/nodeStartToTest.ts",
			OutputNode: []string{midID},
		},
		{
			ID:         midID,
			Name:       "mid",
			ScriptPath: "bun /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_js_to_test/nodeMidToTest.ts",
			OutputNode: []string{endID},
		},
		{
			ID:         endID,
			Name:       "end",
			ScriptPath: "bun /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_js_to_test/nodeEndToTest.ts",
			OutputNode: []string{},
		},
	}
}

func waitForFileAndCleanUp(t *testing.T, path string) {
	t.Helper()
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatalf("timed out after 5 seconds waiting for file %s", path)
		case <-ticker.C:
			_, err := os.Stat(path)
			if err == nil {
				os.Remove(path) // Clean up the file
				return          // File found, exit successfully
			} else if !os.IsNotExist(err) {
				t.Fatalf("error checking for file %s: %v", path, err)
			}
			// If file does not exist, continue loop
		}
	}
}

func TestToCreateFlowToRun(t *testing.T) {
	route := presentation.StartServer()
	server := httptest.NewServer(route)
	defer server.Close()

	nodesToBody := createNodesInGoToTeste()

	body := usecases.CreateFlowUseCaseInput{
		ID:    uuid.NewString(),
		Nodes: nodesToBody,
		Name:  "flowTestGo",
	}

	bodyJson, _ := json.Marshal(body)

	resp, err := http.Post(baseUrl+"/create-flow-to-run", "application/json", bytes.NewBuffer(bodyJson))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	dir, err := os.Getwd()
	assert.NoError(t, err)
	pathOfFile := filepath.Join(dir, "..", "..", "..", "..", "nodeEndToTestOutput.json")

	waitForFileAndCleanUp(t, pathOfFile)
}

func TestToCreateFlowWithBunNodes(t *testing.T) {
	route := presentation.StartServer()
	server := httptest.NewServer(route)
	defer server.Close()

	nodesToBody := createNodesInBunToTeste()

	body := usecases.CreateFlowUseCaseInput{
		ID:    uuid.NewString(),
		Nodes: nodesToBody,
		Name:  "flowTestBun",
	}

	bodyJson, _ := json.Marshal(body)

	resp, err := http.Post(baseUrl+"/create-flow-to-run", "application/json", bytes.NewBuffer(bodyJson))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	dir, err := os.Getwd()
	assert.NoError(t, err)
	pathOfFile := filepath.Join(dir, "..", "..", "..", "..", "nodeEndToTestOutput.json")

	waitForFileAndCleanUp(t, pathOfFile)
}
