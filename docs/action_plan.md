# Plano de Ação - Universal Flow

Este documento descreve o plano de ação para o desenvolvimento do sistema **Universal Flow**, com base nos diagramas de arquitetura e nos documentos de requisitos e ciclo de vida.

## 1. Configuração Inicial do Projeto

-   [ ] Inicializar o módulo Go (`go mod init`).
-   [ ] Definir a estrutura de diretórios principal (`/cmd` para os executáveis, `/internal` para o código do projeto).
-   [ ] Configurar o `docker-compose.yml` para levantar um banco de dados **Postgres**, que servirá como o repositório de estado.

## 2. Estruturas de Dados e Modelos

-   [ ] Implementar os `structs` em Go para `FlowInstance` e `NodeInstance`.
-   [ ] As estruturas devem espelhar o JSON definido nos documentos de requisitos.
-   [ ] Adicionar tags de mapeamento objeto-relacional (ex: GORM) para a persistência no banco de dados.

## 3. Camada de Gerenciamento de Estado (Database)

-   [ ] Configurar a conexão com o banco de dados Postgres utilizando um driver e, opcionalmente, um ORM como GORM.
-   [ ] Criar um Data Access Object (DAO) para abstrair as operações de banco de dados.
-   [ ] Implementar as seguintes funções no DAO:
    -   `CreateFlowInstance(flow *FlowInstance)`
    -   `GetFlowInstance(id string) (*FlowInstance, error)`
    -   `UpdateFlowState(...)`
    -   `GetConsolidatedState(flowID string) (map[string]interface{}, error)`

## 4. Implementação da API REST

-   [ ] Configurar um servidor HTTP utilizando um framework (ex: Gin) para agilidade.
-   [ ] Implementar os endpoints REST críticos definidos na documentação:
    -   `POST /api/v1/flows`: Cria uma nova instância de fluxo com status `pending`.
    -   `GET /api/v1/flows/:flow_id`: Retorna o estado completo de uma instância.
    -   `GET /api/v1/flows/:flow_id/states`: Endpoint para os nós obterem seu estado de entrada consolidado.
    -   `POST /api/v1/flows/:flow_id/nodes/:node_id/finish`: Endpoint para um nó sinalizar sua finalização. Esta chamada deve acionar o `FlowEngine` para continuar a orquestração.

## 5. Lógica de Orquestração (`FlowEngine`)

-   [ ] Desenvolver o `FlowEngine` como um componente que roda em background.
-   [ ] Implementar um mecanismo para que o `FlowEngine` inicie a execução de fluxos com status `pending`.
-   [ ] Implementar a lógica principal que é acionada pela chamada à API `/finish`:
    1.  Atualizar o estado do fluxo com o resultado do nó que terminou.
    2.  Identificar o próximo nó a ser executado a partir do campo `selected_node`.
    3.  Invocar o próximo nó. A estratégia inicial será fazer uma chamada HTTP para um serviço externo pré-configurado.

## 6. Implementação do Notificador (`FlowNotificator`)

-   [ ] Criar um componente que possa escutar eventos de mudança de estado do fluxo (ex: `flow_completed`, `flow_failed`).
-   [ ] Implementar a lógica para, ao capturar um evento, enviar uma notificação para um `Backend` externo via um webhook HTTP.

## 7. Criação de Nós Exemplo (Serviços Externos)

-   [ ] Desenvolver 2 ou 3 serviços web simples e independentes (em Go, Python, ou Node.js) que simularão a execução de nós.
-   [ ] Cada serviço de nó deve ter:
    1.  Um endpoint para ser "invocado" pelo `FlowEngine`.
    2.  A lógica para, ao ser invocado, fazer as chamadas de volta para a API do Universal Flow (`/states` para obter dados e `/finish` para se declarar concluído).

## 8. Testes e Integração

-   [ ] Escrever testes unitários para a camada de acesso ao banco de dados, os handlers da API e a lógica do `FlowEngine`.
-   [ ] Desenvolver um teste de integração ponta a ponta que:
    1.  Inicia a aplicação Universal Flow e os serviços de nós de exemplo.
    2.  Faz uma chamada `POST /api/v1/flows` para iniciar um fluxo.
    3.  Verifica se o fluxo transita pelos status corretamente até ser `completed` ou `failed`.
    4.  Valida se o estado final no banco de dados está correto.
