

async function getStateFromFlow() {
    const flowId = process.env.FLOW_ID
    const nodeId = process.env.NODE_ID
    
    const url = `http://localhost:8080/api/flow-state/get-flow-state/${flowId}`

    const resposta = await fetch(url);
    const dados = await resposta.json();
    console.log(dados);
}

(async ()=>{
    getStateFromFlow()
})()