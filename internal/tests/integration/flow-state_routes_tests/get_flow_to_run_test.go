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

func TestToGetFlowState(t *testing.T) {
	route := presentation.StartServer()
	w := httptest.NewRecorder()

	nodesToBody := createNodesToTest()

	onlyCreateFlow := true

	body := usecases.CreateFlowUseCaseInput{
		ID:             uuid.NewString(),
		Nodes:          nodesToBody,
		Name:           "flowTest",
		OnlyCreateFlow: &onlyCreateFlow,
	}

	bodyJson, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/flow-state/create-flow-to-run", bytes.NewBuffer(bodyJson))
	route.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()

	getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/flow-state/get-flow-state/%s", body.ID), nil)

	route.ServeHTTP(w, getReq)

	var flowResponse entities.Flow

	json.Unmarshal(w.Body.Bytes(), &flowResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, body.ID, flowResponse.ID)
	assert.Equal(t, body.Name, flowResponse.FlowName)
	assert.Equal(t, len(body.Nodes), len(flowResponse.Nodes))
}
