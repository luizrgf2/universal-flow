import { FinishInterface, FlowInterface, NodeInterface } from "./type"
import {writeFile} from 'fs/promises'


const baseURL = "http://localhost:8080/api/flow-state"
//process.env.FLOW_ID = "b620fec7-19bc-4c36-957d-248e666410c3"
//process.env.NODE_ID = "f50c3f51-8871-4e77-8132-bd5c8245bbfb"


function getNode(nodes: NodeInterface[], nodeId: string) {
    return nodes.find(node => node.id === nodeId)   
}



async function getStateFromFlow() {
    const flowId = process.env.FLOW_ID
    const nodeId = process.env.NODE_ID as string
    
    const url = baseURL + `/get-flow-state/${flowId}`
    const urlFinish = baseURL + `/finish-node?flowId=${flowId}`

    const resposta = await fetch(url);
    const dados = await resposta.json() as FlowInterface;
    
    const node = getNode(dados.nodes, nodeId)

    const body = {
        flow_id: process.env.FLOW_ID as string,
        node_id: process.env.NODE_ID as string,
        next_node_id: node?.outputNodes[0] as string,
    } as FinishInterface

    const resposta2 = await fetch(urlFinish, {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(body)
    });
    const dados2 = await resposta2.json();
    await writeFile('nodeEndToTestOutput.json', JSON.stringify(dados2, null, 2))
}

(async ()=>{
    await getStateFromFlow()
})()