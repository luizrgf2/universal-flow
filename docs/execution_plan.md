## 8. Plano de Execução e Orquestração

Esta seção detalha o ciclo de vida de uma instância de fluxo, desde sua criação até a sua finalização.

### 8.1. Fase 1: Iniciação do Fluxo

Um cliente (um sistema externo ou uma interface de usuário) inicia a execução de um novo fluxo de trabalho através de uma chamada de API REST.

-   **Endpoint:** `POST /api/v1/flows`
-   **Payload de Exemplo:**
    ```json
    {
      "flow_name": "customer-onboarding-v2",
      "initial_data": {
        "customer_id": "abc-123",
        "tier": "premium"
      }
    }
    ```
-   **Ação do Sistema:**
    1.  O `FlowManager` recebe a requisição.
    2.  Ele cria uma nova instância de fluxo com base em um "template" de fluxo pré-definido (identificado pelo `flow_name`).
    3.  Gera um `id` único para a instância.
    4.  Define o `status` inicial da instância como `pending`.
    5.  Armazena os `initial_data` no estado do nó inicial.
    6.  Salva a nova instância no banco de dados.
-   **Resposta da API:** O `FlowManager` retorna a representação completa da nova instância de fluxo, incluindo seu `id` e o `status: "pending"`.

### 8.2. Fase 2: Orquestração pelo FlowEngine

O `FlowEngine` é o componente responsável por conduzir a instância do fluxo através de seus estágios.

1.  **Seleção do Fluxo:** O motor monitora o banco de dados por instâncias com status `pending` ou aguarda uma notificação (via pub/sub) sobre uma nova instância.

2.  **Início da Execução:**
    -   O `FlowEngine` carrega a instância do fluxo.
    -   Altera o `status` do fluxo para `running`.
    -   Identifica o primeiro nó a ser executado (o "start node").
    -   Define o campo `current_node` do fluxo com o ID deste nó.

3.  **Ciclo de Execução do Nó:** Para cada nó no `current_node`:
    -   **a. Invocação do Nó:** O motor invoca a lógica do nó. Esta pode ser uma chamada de função Go, um processo externo, ou uma mensagem para um serviço de worker.
    -   **b. Execução Autônoma do Nó:** A lógica do nó é executada:
        -   Ele consulta a API de estado (`GET /api/v1/flows/{flow_id}/states`) para buscar os dados de que precisa.
        -   Executa sua tarefa de negócio (enviar um e-mail, processar dados, etc.).
        -   Ao terminar, o nó **obrigatoriamente** chama a API para registrar seu resultado.
            -   **Endpoint:** `POST /api/v1/flows/{flow_id}/nodes/{node_id}/finish`
            -   **Payload:** `{ "status": "completed", "output": { ... }, "selected_node": "next-node-id" }`
            -   Em caso de falha, o payload seria: `{ "status": "failed", "error": "details here" }`
    -   **c. Atualização de Estado pelo FlowEngine:** O `FlowEngine` é notificado da finalização do nó (através de um evento ou ao receber a resposta da API).
        -   Ele atualiza o `status`, `output` e `error` do nó que acabou de rodar.
        -   Adiciona o ID do nó finalizado ao array `previous_nodes_runned`.
        -   Define `previous_node` com o valor de `current_node`.
        -   Lê o campo `selected_node` retornado pelo nó para determinar o próximo passo.
        -   Atualiza `current_node` com o valor de `selected_node`.
        -   Define `next_node` com base nas `output_nodes` do novo `current_node`.

4.  **Finalização do Fluxo:**
    -   O ciclo continua até que o `current_node` seja um nó final (sem `selected_node` para um próximo passo).
    -   O `FlowEngine` então altera o `status` geral do fluxo para `completed`.

5.  **Tratamento de Falhas:**
    -   Se qualquer nó reportar um `status` de `failed`, o `FlowEngine` interrompe imediatamente o ciclo.
    -   Ele define o `status` geral do fluxo para `failed`.
    -   Nenhum nó subsequente é executado.

### 8.3. Visualização do Ciclo de Vida

O estado de uma instância de fluxo transita de forma previsível:

`pending` → `running` → (`completed` | `failed`)

-   **pending:** A instância foi criada e está na fila para execução.
-   **running:** O `FlowEngine` está ativamente orquestrando a execução de seus nós.
-   **completed:** Todos os nós foram executados com sucesso até um nó final.
-   **failed:** A execução foi interrompida devido a um erro em um dos nós.

### 8.4. Estruturas de Dados Detalhadas

#### 8.4.1. FlowInstance
Representa a instância completa de um fluxo de trabalho em execução.

```json
{
  "id": "flow-instance-uuid",
  "flow_name": "Customer Onboarding",
  "status": "running",
  "initial_data": { "customer_id": "abc-123" },
  "current_node": "send-welcome-email-uuid",
  "previous_node": "start-node-uuid",
  "next_node": "wait-for-activation-uuid",
  "previous_nodes_runned": ["start-node-uuid"],
  "nodes": [
    {
      "id": "start-node-uuid",
      "name": "StartNode",
      "script_path": "node start.js",
      "status": "completed",
      "state": {
        "input": { "customer_id": "abc-123" }
      },
      "error": null,
      "output_nodes": ["send-welcome-email-uuid"],
      "selected_node": "send-welcome-email-uuid"
    },
    {
      "id": "send-welcome-email-uuid",
      "name": "SendWelcomeEmail",
      "script_path": "node sendEmail.js",
      "status": "running",
      "state": null,
      "error": null,
      "output_nodes": ["wait-for-activation-uuid"],
      "selected_node": null
    }
  ]
}
```

#### 8.4.2. NodeInstance
Representa o estado de um nó específico dentro de uma `FlowInstance`.

-   `id`: UUID único do nó na instância.
-   `name`: Nome do "template" do nó.
-   `script_path`: O caminho para o script que o nó irá executar (ex: `node sendEmail.js` ou `go run main.go`).
-   `status`: Estado de execução do nó (`pending`, `running`, `completed`, `failed`).
-   `state`: Dados de entrada ou intermediários usados pelo nó.
-   `error`: Mensagem de erro, caso a execução falhe.
-   `output_nodes`: Lista de IDs dos nós para os quais este nó pode transicionar.
-   `selected_node`: O ID do nó escolhido como próximo passo, preenchido ao finalizar a execução.

### 8.5. API de Controle e Estado

A interação com o motor do fluxo ocorre através de uma API REST bem definida.

#### `POST /api/v1/flows`
-   **Propósito:** Iniciar uma nova instância de fluxo.
-   **Payload:**
    ```json
    {
      "flow_name": "nome-do-fluxo",
      "initial_data": { "chave": "valor" }
    }
    ```
-   **Resposta (201 Created):** O objeto `FlowInstance` recém-criado com `status: "pending"`.

#### `GET /api/v1/flows/{flow_id}`
-   **Propósito:** Obter o estado completo e atual de uma instância de fluxo.
-   **Resposta (200 OK):** O objeto `FlowInstance` completo.

#### `GET /api/v1/flows/{flow_id}/states`
-   **Propósito:** Usado por um nó em execução para buscar os dados de que necessita.
-   **Ação do Sistema:** Agrega os campos `initial_data` e os `output` de todos os nós já concluídos (`status: "completed"`) para fornecer uma visão de estado consolidada.
-   **Resposta (200 OK):**
    ```json
    {
      "consolidated_state": {
        "customer_id": "abc-123",
        "tier": "premium",
        "start_node_output": { "validation": "ok" }
      }
    }
    ```

#### `POST /api/v1/flows/{flow_id}/nodes/{node_id}/finish`
-   **Propósito:** Endpoint obrigatório que um nó chama para sinalizar sua finalização.
-   **Payload:**
    ```json
    {
      "status": "completed",
      "output": { "email_sent_at": "2025-12-24T10:00:00Z" },
      "selected_node": "proximo-node-id",
      "error": null
    }
    ```
-   **Ação do Sistema:** O `FlowEngine` recebe esta chamada, atualiza o estado do nó e do fluxo, e continua a orquestração para o `selected_node`.
-   **Resposta (202 Accepted):** Confirmação de que a atualização foi recebida e está sendo processada.