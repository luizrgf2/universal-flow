# Universal Flow

O Universal Flow é um motor de fluxo leve e orientado a estado, projetado para integração universal. Seu princípio fundamental é orquestrar fluxos de trabalho onde cada passo (um "nó" ou *node*) é um processo independente e autônomo que se comunica exclusivamente através de um gerenciador de estado centralizado.

## 🎯 Sobre o Projeto

O objetivo principal do Universal Flow é fornecer um motor simples e poderoso para gerenciar fluxos de trabalho complexos em qualquer aplicação. Em vez de passar dados diretamente entre funções ou serviços, os nós são completamente desacoplados. Eles não recebem entradas diretas nem retornam saídas diretas. Em vez disso, eles:

1.  **Consultam** uma API central para obter os dados de que precisam.
2.  **Executam** sua lógica de negócio.
3.  **Chamam** a API central novamente para reportar seu resultado e determinar o próximo passo no fluxo.

Isso torna o sistema altamente modular, escalável e fácil de depurar, pois todo o estado de cada fluxo de trabalho é persistido e auditável em todos os momentos.

## 🚀 Começando

Para executar o motor do Universal Flow, simplesmente execute o seguinte comando no diretório raiz do projeto:

```bash
go run main.go
```

O servidor será iniciado na porta `8080`.

## ⚙️ Como Funciona: A API

O motor é controlado através de uma API REST simples.

### 1. Criar e Executar um Fluxo

Para iniciar um novo fluxo de trabalho, você envia uma requisição `POST` com a estrutura do fluxo. O motor irá salvá-lo e começar imediatamente a executar o primeiro nó.

**Endpoint:** `POST /api/flow-state/create-flow-to-run`

**Corpo (Body):**

```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "name": "Meu Primeiro Fluxo",
  "nodes": [
    {
      "id": "a1b2c3d4-e5f6-a7b8-c9d0-e1f2a3b4c5d6",
      "name": "Nó Inicial",
      "script_path": "node /path/to/your/start-script.js",
      "output_node": ["b2c3d4e5-f6a7-b8c9-d0e1-f2a3b4b5c6d7"]
    },
    {
      "id": "b2c3d4e5-f6a7-b8c9-d0e1-f2a3b4b5c6d7",
      "name": "Nó Intermediário",
      "script_path": "node /path/to/your/middle-script.js",
      "output_node": ["c3d4e5f6-a7b8-c9d0-e1f2a3b4b5c6d7", "d4e5f6a7-b8c9-d0e1-f2a3b4b5c6d8"]
    },
    {
      "id": "c3d4e5f6-a7b8-c9d0-e1f2a3b4b5c6d7",
      "name": "Nó Final (Caminho A)",
      "script_path": "node /path/to/your/end-script-A.js",
      "output_node": []
    },
    {
      "id": "d4e5f6a7-b8c9-d0e1-f2a3b4b5c6d8",
      "name": "Nó Final (Caminho B)",
      "script_path": "node /path/to/your/end-script-B.js",
      "output_node": []
    }
  ]
}
```

#### O Campo `output_node`

O campo `output_node` é fundamental para definir a estrutura do seu fluxo. Ele é um array de strings que contém os IDs de todos os nós que são **destinos possíveis** a partir do nó atual.

-   **Fluxos Lineares:** Se um nó tem apenas um caminho possível, o `output_node` conterá um único ID.
-   **Bifurcações e Condicionais:** Se um nó pode levar a múltiplos caminhos diferentes (ex: "Aprovar" vs. "Rejeitar"), o `output_node` conterá os IDs de todos os nós de destino possíveis.
-   **Nós Finais:** Um nó que finaliza um fluxo (ou um caminho do fluxo) terá um array vazio `[]`.

A responsabilidade de **escolher** qual caminho seguir, dentre as opções listadas no `output_node`, é da lógica interna do script do nó. Ao chamar o endpoint `finish-node`, o script deve passar o ID do nó escolhido no campo `next_node_id`.

### 2. Obter Estado do Fluxo

Você pode recuperar o estado completo e em tempo real de qualquer fluxo de trabalho a qualquer momento. Isso é útil para monitoramento e depuração.

**Endpoint:** `GET /api/flow-state/get-flow-state/:id`

-   `:id` é o ID do fluxo que você deseja inspecionar.

### 3. Finalizar um Nó

Este é o endpoint mais crítico para desenvolvedores de nós. Quando o script de um nó termina sua tarefa, ele **deve** chamar este endpoint para informar ao motor que terminou e o que deve acontecer a seguir.

**Endpoint:** `PATCH /api/flow-state/finish-node?flowId=<FLOW_ID>`

#### Cenário de Sucesso

Para marcar o nó como `completed` (concluído) e informar ao motor qual nó executar em seguida.

**Corpo (Body):**

```json
{
  "node_id": "b2c3d4e5-f6a7-b8c9-d0e1-f2a3b4b5c6d7",
  "next_node_id": "c3d4e5f6-a7b8-c9d0-e1f2a3b4b5c6d7",
  "node_output": "{\"resultado\":\"ok\"}"
}
```

#### Cenário de Falha

Para marcar o nó como `failed` (falhou), o que também irá parar todo o fluxo e marcá-lo como `failed`.

**Corpo (Body):**

```json
{
  "node_id": "b2c3d4e5-f6a7-b8c9-d0e1-f2a3b4b5c6d7",
  "error_message": "Falha ao conectar com o serviço externo."
}
```

## 💻 Desenvolvimento de Nós (Nodes)

Um nó é simplesmente um comando executável (ex: `node meu-script.js`, `python process.py`). O motor em si não conhece nem se importa com a lógica interna do nó. Ele apenas se preocupa com a comunicação do nó com a API.

### Injeção de Contexto

Quando o motor do Universal Flow executa um nó, ele injeta o contexto de execução como **variáveis de ambiente**:

-   `FLOW_ID`: O ID do fluxo atualmente em execução.
-   `NODE_ID`: O ID da instância específica do nó que está sendo executada.

Seu script deve ler essas variáveis para interagir com a API corretamente.

### Exemplo de Ciclo de Vida de um Nó

1.  O motor executa seu script (ex: `node meu-script.js`).
2.  Dentro do seu script, você lê as variáveis de ambiente: `process.env.FLOW_ID` e `process.env.NODE_ID`.
3.  (Opcional) Seu script pode fazer uma requisição `GET` para `/api/flow-state/get-flow-state/:id` para obter dados produzidos por nós anteriores.
4.  Seu script executa sua lógica de negócio.
5.  Seu script faz uma requisição `PATCH` para `/api/flow-state/finish-node` para sinalizar sua conclusão, escolhendo o próximo nó a ser executado ou reportando um erro.

### Exemplo de Nó (Node.js)

Aqui está um exemplo completo de um script de nó. Ele lê seu contexto, executa uma tarefa e reporta sua conclusão ao motor do Universal Flow.

```javascript
// Um script de nó simples para o Universal Flow

// Esta é a função que fará o trabalho principal.
async function main() {
    // 1. Ler o contexto das variáveis de ambiente.
    const flowId = process.env.FLOW_ID;
    const nodeId = process.env.NODE_ID;

    if (!flowId || !nodeId) {
        console.error("Erro: FLOW_ID e NODE_ID devem ser definidos.");
        // Se o contexto estiver faltando, não podemos prosseguir.
        // Em um cenário real, você poderia querer reportar isso como uma falha.
        return;
    }

    console.log(`Executando o nó ${nodeId} para o fluxo ${flowId}`);

    // 2. (Opcional) Buscar o estado atual do fluxo para obter dados.
    // Neste exemplo, assumimos que este nó precisa da saída do nó anterior.
    // const flowStateResponse = await fetch(`http://localhost:8080/api/flow-state/get-flow-state/${flowId}`);
    // const flowState = await flowStateResponse.json();
    // console.log("Estado atual do fluxo:", flowState);

    // 3. Executar a lógica de negócio.
    // Para este exemplo, vamos apenas esperar um segundo para simular trabalho.
    await new Promise(resolve => setTimeout(resolve, 1000));
    const resultado = { message: "Tarefa concluída com sucesso!", timestamp: new Date().toISOString() };

    // 4. Finalizar o nó chamando a API.
    // Este nó decidirá mover para o próximo nó, que vamos definir aqui para simplificar.
    // Em um cenário real, o `nextNodeId` pode ser determinado pela lógica de negócio.
    const nextNodeId = "c3d4e5f6-a7b8-c9d0-e1f2a3b4b5c6d7"; // ID de exemplo

    const finishUrl = `http://localhost:8080/api/flow-state/finish-node?flowId=${flowId}`;

    const finishBody = {
        node_id: nodeId,
        next_node_id: nextNodeId,
        node_output: JSON.stringify(resultado)
    };

    try {
        const response = await fetch(finishUrl, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(finishBody)
        });

        if (response.ok) {
            console.log(`Nó ${nodeId} finalizado com sucesso.`);
        } else {
            const errorData = await response.json();
            console.error(`Falha ao finalizar o nó ${nodeId}:`, errorData);
        }
    } catch (error) {
        console.error("Erro ao chamar a API de finish-node:", error);
    }
}

// Executa a função principal.
main();
```

Este exemplo demonstra o padrão principal: ler o contexto, fazer o trabalho e reportar de volta ao motor. Isso mantém seus nós simples, sem estado (stateless) e universalmente compatíveis com o motor de fluxo.
