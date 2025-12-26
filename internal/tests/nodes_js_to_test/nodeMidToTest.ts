import { FinishInterface, FlowInterface, NodeInterface } from "./type"
import {writeFile} from 'fs/promises'


const baseURL = "http://localhost:8080/api/flow-state"
//process.env.FLOW_ID = "b620fec7-19bc-4c36-957d-248e666410c3"
//process.env.NODE_ID = "f50c3f51-8871-4e77-8132-bd5c8245bbfb"


function getNode(nodes: NodeInterface[], nodeId: string) {
    return nodes.find(node => node.id === nodeId)   
}

async function getFlowState() {
    const flowId = process.env.FLOW_ID
    const url = baseURL + `/get-flow-state/${flowId}`
    const resposta = await fetch(url);
    const dados = await resposta.json() as FlowInterface;
    return dados
}

async function finishNode(node?: NodeInterface, output?: object, error?: string) {
    const flowId = process.env.FLOW_ID as string
    const nodeId = process.env.NODE_ID as string
    
    const urlFinish = baseURL + `/finish-node?flowId=${flowId}`
        const body = {
        flow_id: flowId,
        node_id: nodeId,
        next_node_id: node ? node?.outputNodes[0] as string : undefined,
        node_output: output ? JSON.stringify(output) : undefined,
        error_message: error ? error : undefined
    } as FinishInterface

    const resposta2 = await fetch(urlFinish, {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(body)
    });
    const dados2 = await resposta2.json();

}




async function getStateFromFlow() {
    const flowId = process.env.FLOW_ID
    const nodeId = process.env.NODE_ID as string
    
    const FLOW_STATE = await getFlowState()
    const node = getNode(FLOW_STATE.nodes, nodeId)
    const nodeInput = node?.state?.input ? JSON.parse(node.state.input) as any : undefined 

    //get filename from flow state
    const fileName = nodeInput.fileName as string




    if(!node) 
        return await finishNode(node, undefined, "Node not found")
    
    await writeFile(fileName, "Hello Password")

    await finishNode(node, {password: "Hello Password"})

}

(async ()=>{
    await getStateFromFlow()
})()