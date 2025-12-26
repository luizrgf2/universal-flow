export interface FinishInterface {
    flow_id: string
    node_id: string
    next_node_id: string
    error_message?: string,
    node_output?: string
}

export interface NodeInterface {
    error?: string
    id: string
    name: string
    outputNodes: string[]
    script_path: string
    selectedNode?: string
    state: {input: '', output: ''}
    status: string
}

export interface FlowInterface {
    currentNode?: string
    flowName: string
    id: string
    nextNode?: string
    nodes: NodeInterface[]
    status: string
    previous_nodes_runned: string[]
    previousNode?: string
}