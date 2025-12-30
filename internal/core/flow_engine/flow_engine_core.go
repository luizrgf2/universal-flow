package flowengine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/luizrgf2/universal-flow/internal/core/entities"
	services_interfaces "github.com/luizrgf2/universal-flow/internal/core/services"
	"github.com/luizrgf2/universal-flow/internal/core/types"
)

type FlowEngineCore struct {
	flowStateManagerService services_interfaces.FlowStateManagerService
}

func NewFlowEngine(flowStateManagerService services_interfaces.FlowStateManagerService) *FlowEngineCore {
	return &FlowEngineCore{flowStateManagerService: flowStateManagerService}
}

func (fe *FlowEngineCore) getNodeById(id string, nodes []entities.Node) (*entities.Node, error) {
	for _, node := range nodes {
		if node.ID == id {
			return &node, nil
		}
	}
	return nil, fmt.Errorf("node with id: %s not exists in flow", id)
}

func (fe *FlowEngineCore) selectNextNodeToRun(flow *entities.Flow) (*entities.Node, error) {
	if flow.NextNode == nil {
		return &flow.Nodes[0], nil
	} else {
		node, err := fe.getNodeById(*flow.NextNode, flow.Nodes)
		if err != nil {
			return nil, err
		}
		return node, nil
	}
}

func (fe *FlowEngineCore) updateFlowState(flow *entities.Flow) error {
	err := fe.flowStateManagerService.UpdateFlow(flow)
	if err != nil {
		return err
	}
	return nil
}

func (fe *FlowEngineCore) changeFlowStatus(flow *entities.Flow, status string) error {
	statusToChange, err := types.CreateFlowStatus(status)

	if flow.Status == statusToChange {
		return nil
	}

	err = flow.ChangeFlowStatus(statusToChange)
	if err != nil {
		return err
	}

	err = fe.updateFlowState(flow)
	if err != nil {
		return err
	}

	return nil
}

func (fe *FlowEngineCore) changeNodeStatus(nodeToUpdate *entities.Node, status string) error {
	statusToChange, err := types.CreateNodeStatus(status)

	if nodeToUpdate.Status == statusToChange {
		return nil
	}

	err = nodeToUpdate.ChangeNodeStatus(statusToChange)
	if err != nil {
		return err
	}

	return nil
}

func (fe *FlowEngineCore) execJSNodeOrBun(nodeToRun *entities.Node, flowID string) error {
	comands := strings.Split(nodeToRun.ScriptPath, " ")
	comandMain := comands[0]
	comands = comands[1:]

	if (strings.Contains(nodeToRun.ScriptPath, "node") || strings.Contains(nodeToRun.ScriptPath, "bun")) && (strings.Contains(nodeToRun.ScriptPath, ".js") || strings.Contains(nodeToRun.ScriptPath, ".ts")) {
		execNode := exec.Command(comandMain, comands...)

		execNode.Env = append(os.Environ(),
			"FLOW_ID="+flowID,
			"NODE_ID="+nodeToRun.ID,
		)

		output, errCommand := execNode.CombinedOutput()

		if len(output) > 0 {
			fmt.Print(string(output))
		}

		failedStatus, err := types.CreateNodeStatus("failed")
		if err != nil {
			return err
		}

		if errCommand != nil {
			nodeToRun.ChangeNodeStatus(failedStatus)
			nodeToRun.ChangeError(errCommand.Error())
			return errCommand
		}

		completedStatus, err := types.CreateNodeStatus("completed")
		if err != nil {
			return err
		}

		nodeToRun.ChangeNodeStatus(completedStatus)
		return nil
	}

	failedStatus, err := types.CreateNodeStatus("failed")
	if err != nil {
		return err
	}

	nodeToRun.ChangeNodeStatus(failedStatus)
	nodeToRun.ChangeError("Error not js valid command")

	return fmt.Errorf("Error not js valid command")
}

func (fe *FlowEngineCore) execGoNode(nodeToRun *entities.Node, flowID string) error {
	comands := strings.Split(nodeToRun.ScriptPath, " ")
	comandMain := comands[0]
	comands = comands[1:]

	if strings.Contains(nodeToRun.ScriptPath, "go run") && strings.Contains(nodeToRun.ScriptPath, ".go") {
		execNode := exec.Command(comandMain, comands...)

		execNode.Env = append(os.Environ(),
			"FLOW_ID="+flowID,
			"NODE_ID="+nodeToRun.ID,
		)

		output, errCommand := execNode.CombinedOutput()

		if len(output) > 0 {
			fmt.Print(string(output))
		}

		failedStatus, err := types.CreateNodeStatus("failed")
		if err != nil {
			return err
		}

		if errCommand != nil {
			nodeToRun.ChangeNodeStatus(failedStatus)
			nodeToRun.ChangeError(errCommand.Error())
			return errCommand
		}

		completedStatus, err := types.CreateNodeStatus("completed")
		if err != nil {
			return err
		}

		nodeToRun.ChangeNodeStatus(completedStatus)
		return nil
	}

	failedStatus, err := types.CreateNodeStatus("failed")
	if err != nil {
		return err
	}

	nodeToRun.ChangeNodeStatus(failedStatus)
	nodeToRun.ChangeError("Error not go valid command")

	return fmt.Errorf("Error not go valid command")
}

func (fe *FlowEngineCore) execNode(nodeToRun *entities.Node, flowID string) error {
	isJSOrBun := (strings.Contains(nodeToRun.ScriptPath, "node") || strings.Contains(nodeToRun.ScriptPath, "bun")) && (strings.Contains(nodeToRun.ScriptPath, ".js") || strings.Contains(nodeToRun.ScriptPath, ".ts"))
	isGo := strings.Contains(nodeToRun.ScriptPath, "go run") && strings.Contains(nodeToRun.ScriptPath, ".go")

	if isJSOrBun {
		return fe.execJSNodeOrBun(nodeToRun, flowID)
	}

	if isGo {
		return fe.execGoNode(nodeToRun, flowID)
	}

	failedStatus, err := types.CreateNodeStatus("failed")
	if err != nil {
		return err
	}

	nodeToRun.ChangeNodeStatus(failedStatus)
	errMessage := "Error not valid command, valid commands are node, bun and go"
	nodeToRun.ChangeError(errMessage)

	return errors.New(errMessage)
}

func (fe *FlowEngineCore) updateNodeInFlow(flow *entities.Flow, nodeToUpdate *entities.Node) {
	for i, node := range flow.Nodes {
		if node.ID == nodeToUpdate.ID {
			flow.Nodes[i] = *nodeToUpdate
			return
		}
	}
}

func (fe *FlowEngineCore) RunFlow(flow *entities.Flow) error {
	err := fe.changeFlowStatus(flow, "running")
	if err != nil {
		return err
	}

	nodeToRun, err := fe.selectNextNodeToRun(flow)
	if err != nil {
		return err
	}

	err = fe.changeNodeStatus(nodeToRun, "running")
	if err != nil {
		return err
	}

	err = flow.SetCurrentNode(nodeToRun.ID)
	if err != nil {
		return nil
	}

	fe.updateNodeInFlow(flow, nodeToRun)
	err = fe.updateFlowState(flow)
	if err != nil {
		return err
	}

	err = fe.execNode(nodeToRun, flow.ID)
	if err != nil {
		return err
	}

	return nil
}
