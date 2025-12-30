package nodes_go_to_test

type FinishInterface struct {
	FlowID       string `json:"flow_id"`
	NodeID       string `json:"node_id"`
	NextNodeID   string `json:"next_node_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	NodeOutput   string `json:"node_output,omitempty"`
}

type NodeInterface struct {
	Error        string            `json:"error,omitempty"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	OutputNodes  []string          `json:"outputNodes"`
	ScriptPath   string            `json:"script_path"`
	SelectedNode string            `json:"selectedNode,omitempty"`
	State        map[string]string `json:"state"`
	Status       string            `json:"status"`
}

type FlowInterface struct {
	CurrentNode         string          `json:"currentNode,omitempty"`
	FlowName            string          `json:"flowName"`
	ID                  string          `json:"id"`
	NextNode            string          `json:"nextNode,omitempty"`
	Nodes               []NodeInterface `json:"nodes"`
	Status              string          `json:"status"`
	PreviousNodesRunned []string        `json:"previous_nodes_runned"`
	PreviousNode        string          `json:"previousNode,omitempty"`
}
