# Ciclo de Vida de um Node no Universal Flow

Este documento detalha o ciclo de vida e o modelo de execução de um nó (node) dentro do Universal Flow. A arquitetura foi desenhada para garantir que os nós sejam independentes, sem estado (stateless) e agnósticos em relação ao fluxo de dados, promovendo reusabilidade e simplicidade.

## Princípios Fundamentais

1.  **Sem Entradas (Inputs) Diretas:** As funções que representam os nós **não recebem parâmetros** em sua assinatura. Toda a informação necessária para a execução de um nó deve ser obtida a partir do serviço de estado central do Universal Flow.
2.  **Sem Saídas (Outputs) Diretas:** As funções dos nós **não retornam valores** (como `return`). Qualquer resultado, estado intermediário ou dado gerado deve ser persistido de volta no serviço de estado central.
3.  **Estado Centralizado:** O Universal Flow mantém um objeto de estado global para cada instância de fluxo em execução. Este objeto é a única fonte da verdade e o principal meio de comunicação e transferência de dados entre os nós.

## Etapas do Ciclo de Vida de um Nó

O ciclo de vida de um nó pode ser dividido nas seguintes etapas:

### 1. Início da Execução

Quando o motor do Universal Flow determina que um nó é o próximo a ser executado (baseado no campo `next_node` do estado do fluxo), ele invoca a função correspondente àquele nó.

### 2. Consulta ao Estado

A primeira ação de qualquer nó é fazer uma chamada de API para o serviço de estado do Universal Flow para recuperar o estado atual do fluxo. Isso permite que o nó acesse dados que foram produzidos por nós anteriores.

**Exemplo:**
Um nó de "Enviar Email de Boas-Vindas" consultaria o estado para obter o email do cliente, que foi coletado em um passo anterior como o `start-node`.

### 3. Execução da Lógica de Negócio

Com os dados em mãos, o nó executa sua lógica de negócio principal. Isso pode ser qualquer coisa, desde uma simples manipulação de dados até a integração com um serviço externo (como enviar um email, processar um pagamento, etc.).

- **Tratamento de Erro:** Se ocorrer um erro durante a execução, o nó deve capturá-lo e atualizar seu próprio estado no fluxo com a mensagem de erro no campo `error`. Isso interrompe ou desvia o fluxo, conforme a lógica definida, e fornece visibilidade sobre a falha.

### 4. Atualização do Estado

Após a conclusão de sua lógica, o nó faz uma nova chamada à API do serviço de estado para salvar qualquer novo dado ou alteração. Isso pode incluir:

-   O resultado de sua operação.
-   Um status de "concluído".
-   A decisão de qual será o próximo nó.

O nó deve selecionar o próximo passo e atualizar o campo `selected_node` em seu próprio objeto de estado para registrar a decisão.

Imediatamente após a conclusão bem-sucedida de um nó, o motor do fluxo (FlowEngine) deve adicionar o ID desse nó ao array `previous_nodes_runned` no estado global do fluxo. Isso cria um histórico cronológico e imutável do caminho de execução.

### 5. Indicação do Próximo Nó

Finalmente, o motor do fluxo utiliza o valor de `selected_node` para atualizar o campo `next_node` no nível raiz do estado do fluxo, preparando o terreno para a próxima iteração do ciclo. O fluxo então prossegue para o nó indicado.

Este modelo garante que cada passo do fluxo seja auditável e que o estado seja consistentemente rastreado, tornando o sistema robusto e fácil de depurar.
