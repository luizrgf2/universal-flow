# Requisitos do Sistema - Universal Flow

## 1. Filosofia Principal: Arquitetura Orientada a Estado

O princípio fundamental do Universal Flow é que os componentes de um fluxo de trabalho (Nós) são completamente desacoplados e não se comunicam diretamente através de parâmetros de função ou valores de retorno. Toda a comunicação e passagem de dados **deve** ser mediada por um serviço central de gerenciamento de estado (`FlowStateManager`).

Isso garante que cada nó seja uma unidade de trabalho universal e autônoma, e que todo o ciclo de vida de um fluxo seja auditável, rastreável e persistido.

## 2. Requisitos para a Interface e Comportamento de um Nó

-   **R2.1: Sem Parâmetros de Entrada Diretos:** Uma função ou processo que representa um nó **não deve** aceitar parâmetros de negócio em sua assinatura.
    -   *Exemplo de violação:* `function sendEmail(to, body)` está **proibido**.

-   **R2.2: Sem Valores de Retorno Diretos:** A função de um nó **não deve** usar uma instrução `return` para passar seu resultado para o orquestrador.
    -   *Exemplo de violação:* `return "email sent"` está **proibido**.

-   **R2.3: Aquisição de Estado de Entrada via API:** Um nó **deve** obter todos os seus dados de entrada (informações que seriam tradicionalmente passadas como parâmetros) fazendo uma chamada a uma API do `FlowStateManager`.
    -   *Exemplo:* Para enviar um e-mail, o `SendEmailNode` deve consultar a API para obter o endereço do destinatário e o corpo da mensagem, que foram possivelmente definidos por um nó anterior no mesmo fluxo.

-   **R2.4: Persistência do Estado de Saída via API:** Um nó **deve** reportar sua conclusão, resultado ou erro fazendo uma chamada para a API do `FlowStateManager`. Esta ação substitui o `return` tradicional e sinaliza o fim de sua execução.
    -   *Exemplo:* Após enviar o e-mail, o `SendEmailNode` deve chamar a API para salvar seu estado de conclusão, como `{"status": "completed", "output": {"message_id": "xyz-123"}}`.

## 3. Requisitos para a API do FlowStateManager

-   **R3.1: Acesso ao Estado do Fluxo:** A API **deve** fornecer um endpoint para que um nó possa consultar o estado atual do fluxo, incluindo os resultados de todos os nós executados anteriormente.
    -   Ex: `GET /api/v1/flows/{flow_id}/states`

-   **R3.2: Acesso ao Estado de Nós Específicos:** A API **deve** permitir a consulta de resultados de um ou mais nós específicos para facilitar o acesso aos dados de entrada.
    -   Ex: `GET /api/v1/flows/{flow_id}/states?node_id={node_id}`

-   **R3.3: Atualização de Estado do Nó:** A API **deve** fornecer um endpoint para que o nó em execução possa salvar seu próprio estado (output, status, erros). Esta chamada é fundamental, pois sinaliza ao `FlowEngine` que o nó terminou seu trabalho.
    -   Ex: `POST /api/v1/flows/{flow_id}/states`
    -   Payload de exemplo: `{"node_id": "send-email-node-instance-1", "status": "completed", "output": {"timestamp": "...", "details": "Email sent successfully"}}`

## 4. Requisitos para a Lógica de Orquestração (FlowEngine)

-   **R4.1: Execução Sequencial:** O `FlowEngine` **deve** executar os nós de um fluxo na ordem definida na sua estrutura.

-   **R4.2: Orquestração Assíncrona Baseada em Estado:** O `FlowEngine` invoca um nó e **deve** aguardar que este nó chame a API de atualização de estado (R3.3) para considerar sua execução concluída e então decidir o próximo passo (executar o próximo nó ou finalizar o fluxo).

-   **R4.3: Tratamento de Erros Centralizado:** Se um nó reportar um estado de `error` através da API, o `FlowEngine` **deve** interromper a execução dos nós subsequentes e marcar o fluxo inteiro como `failed`, registrando a origem do erro.

## 5. Requisitos de Universalidade e Desacoplamento

-   **R5.1: Independência de Linguagem:** A arquitetura **deve** permitir que nós sejam escritos em qualquer linguagem de programação ou executados em qualquer ambiente, contanto que possam se comunicar com a API do `FlowStateManager` (via HTTP, gRPC, ou outro protocolo definido).

-   **R5.2: Reusabilidade e Modularidade:** Os nós **devem** ser inteiramente reutilizáveis em diferentes fluxos, pois não possuem dependências diretas uns dos outros, apenas com a API de estado. Um `SendEmailNode`, por exemplo, pode ser usado em qualquer fluxo que necessite enviar um e-mail.

## 6. Modelo de Dados de Estado do Fluxo

A estrutura de dados principal que o `FlowStateManager` gerencia para representar o estado de um fluxo de trabalho em execução deve seguir o seguinte formato.

```json
{
  "id": "flow-instance-uuid",
  "flow_name": "Customer Onboarding",
  "status": "running",
  "current_node": "send-welcome-email-uuid",
  "next_node": "end-node-uuid",
  "previous_node": "start-node-uuid",
  "nodes": [
    {
      "id": "start-node-uuid",
      "name": "StartNode",
      "status": "completed",
      "state": {
        "input": { "customer_id": "abc-123" },
        "output": { "customer_id": "abc-123", "onboarding_kit": "full" }
      },
      "error": null,
      "output_nodes": ["send-welcome-email-uuid", "another-possible-path-uuid"],
      "selected_node": "send-welcome-email-uuid"
    },
    {
      "id": "send-welcome-email-uuid",
      "name": "SendEmailNode",
      "status": "running",
      "state": {
        "input": { "email_address": "customer@email.com", "kit": "full" },
        "output": null
      },
      "error": "Failed to connect to SMTP server",
      "output_nodes": ["end-node-uuid"],
      "selected_node": null
    },
    {
      "id": "end-node-uuid",
      "name": "EndNode",
      "status": "pending",
      "state": {},
      "error": null,
      "output_nodes": [],
      "selected_node": null
    }
  ]
}
```

### Descrição dos Campos

*   **`id`**: (string) O ID único da **instância** do fluxo em execução.
*   **`flow_name`**: (string) Um nome descritivo para o modelo do fluxo.
*   **`status`**: (string) O estado geral do fluxo: `pending`, `running`, `completed`, `failed`.
*   **`current_node`**: (string) O ID do nó que está atualmente em execução ou pronto para ser executado. Essencial para o `FlowEngine` saber onde retomar o trabalho.
*   **`next_node`**: (string | null) O ID do próximo nó que será executado assim que o `current_node` for concluído. Ele é determinado pelo valor de `selected_node` do nó atual.
*   **`previous_node`**: (string) O ID do último nó que foi concluído. Útil para rastreabilidade e lógica de compensação.
*   **`nodes`**: (array) A lista de todos os nós que compõem o fluxo.
    *   **`id`**: (string) O ID único da **instância** de um nó dentro do fluxo.
    *   **`name`**: (string) O nome do *tipo* de nó (ex: `SendEmailNode`).
    *   **`status`**: (string) O estado de execução específico daquele nó.
    *   **`state`**: (json) Um objeto que contém os dados que o nó produziu como resultado (`output`). Este campo é preenchido pelo próprio nó ao chamar a API de estado para finalizar sua execução.
    *   **`error`**: (string | json | null) Descreve qualquer erro que tenha ocorrido durante a execução do nó. Pode ser uma simples string ou um objeto JSON estruturado. É `null` se nenhuma falha ocorrer. Essencial para depuração e visibilidade.
    *   **`output_nodes`**: (array de strings) Define as possíveis saídas do nó, contendo a lista de IDs dos próximos nós que *podem* ser executados.
    *   **`selected_node`**: (string | null) O ID do nó que foi efetivamente escolhido como o próximo passo a partir das opções em `output_nodes`. Este campo é preenchido pelo nó ao finalizar sua execução, indicando qual caminho o fluxo deve seguir. É fundamental para a rastreabilidade.

## 7. Requisitos Não-Funcionais

*   **R7.1: Tecnologia de Implementação:** O núcleo do `FlowEngine` e `FlowStateManager` **deve** ser construído utilizando a linguagem de programação **Go (Golang)**.
*   **R7.2: Desempenho e Eficiência:** A aplicação **deve** ser otimizada para alto desempenho e baixo consumo de recursos (CPU e memória), garantindo que seja "extremamente leve".
*   **R7.3: Facilidade de Integração:** O motor do fluxo **deve** ser projetado como um componente de fácil integração, permitindo que outros sistemas o incorporem para gerenciar seus próprios fluxos de trabalho com o mínimo de sobrecarga.
