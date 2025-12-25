package flowengine_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/luizrgf2/universal-flow/internal/core/entities"
	flowengine "github.com/luizrgf2/universal-flow/internal/core/flow_engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FlowStateManagerMocked struct {
	mock.Mock
}

func (m *FlowStateManagerMocked) RunNewFlow(flow *entities.Flow) error {
	args := m.Called(flow)
	return args.Error(0)
}

func TestToRunSimpleNodeInNodeJs(t *testing.T) {
	pathScriptJs := "node /home/luizrgf/projetos/meus/universal-flow/internal/tests/nodes_js_to_test/nodeStartToTest.js"

	flowUuid := uuid.New().String()
	nodeUuid := uuid.New().String()

	nodeToTest, err := entities.CreateNode(nodeUuid, "test", pathScriptJs, []string{"sdfsdfsd"})
	assert.ErrorIs(t, err, nil)
	if err != nil {
		return
	}

	flowToTest, err := entities.CreateFlow(flowUuid, "test", []entities.Node{*nodeToTest, *nodeToTest, *nodeToTest})
	assert.ErrorIs(t, err, nil)
	if err != nil {
		return
	}

	flowStateManagerMocked := &FlowStateManagerMocked{}
	flowStateManagerMocked.On("RunNewFlow", flowToTest).Return(nil)

	flowEngine := flowengine.NewFlowEngine(flowStateManagerMocked)

	err = flowEngine.RunFlow(flowToTest)
	assert.ErrorIs(t, err, nil)

}
