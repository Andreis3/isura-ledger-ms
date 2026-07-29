# TODO List — isura-ledger-ms

## 📋 Phase 1: Consolidar o Use Case `CreateTransaction`
- [ ] **Ajustar mapeamento de contas no `CreateTransaction`**: Corrigir a lógica de atribuição de `debitAccount` e `creditAccount` após a ordenação lexicográfica de IDs para evitar inversão de papéis no cálculo de saldo.
- [ ] **Adicionar validação de moeda (Currency)**: Garantir falha caso `debitAccount.Currency != creditAccount.Currency` ou `input.Currency` divirja das contas.
- [ ] **Garantir Unique Constraint no Postgres**: Validar no schema do Atlas HCL (`db/transactions.pg.hcl`) se a coluna `idempotency_key` possui um índice único para barrar requisições concorrentes idênticas no banco.
- [ ] **Tratamento de Idempotência no Use Case**: Retornar o `TransactionID` existente com status OK caso `ExistsByIdempotencyKey` retorne `true`.

## 📋 Phase 2: Camada de Transporte gRPC (`CreateTransaction`)
- [ ] **Expor a RPC no Proto**: Atualizar `proto/ledger/v1/transaction.proto` com a nova mensagem `CreateTransactionRequest` / `CreateTransactionResponse` e o método no `LedgerService`.
- [ ] **Gerar os stubs gRPC**: Rodar `buf generate` ou o comando do `Makefile` para atualizar `_grpc.pb.go`.
- [ ] **Implementar Handler gRPC**: Criar `CreateTransactionHandler` em `internal/transport/grpc/handler/` orquestrando a conversão de/para Protobuf, chamando `CreateTransaction.Execute` e mapeando erros com o pacote `fault`.
- [ ] **Registrar Handler no Server**: Mapear o novo handler no `internal/transport/grpc/ledger_module.go`.

## 📋 Phase 3: Outbox Relay / Dynamic Dispatcher
- [ ] **Implementar busca segura no `OutboxRepository`**: Adicionar o método `FindPendingForUpdateSkipLocked(ctx, limit)` no repositório do Postgres usando `SELECT ... FOR UPDATE SKIP LOCKED`.
- [ ] **Definir o mecanismo de acionamento do Outbox**:
    - **Opção A (Webhook/HTTP)**: Criar rota `/internal/outbox/process` no `chi` router (`internal/transport/rest/`).
    - **Opção B (gRPC)**: Adicionar RPC `ProcessOutbox` no `.proto`.
- [ ] **Implementar `ProcessOutboxUseCase`**: Lógica de ler os registros pendentes, disparar o payload (HTTP Webhook ou Broker) e marcar como `SUCCESS` ou incrementar `attempts` com cálculo de *backoff*.

## 📋 Phase 4: Testes & Observabilidade
- [ ] **Testes de Unidade do `CreateTransaction`**: Cobrir cenários de saldo insuficiente, moedas incompatíveis e sucesso na pasta `tests/unit/application/command/`.
- [ ] **Testes Concorrentes de Race Condition**: Criar teste executando goroutines paralelas tentando debitar do mesmo saldo simultaneamente.
- [ ] **Validação de Tracing & Spans**: Validar no Grafana Tempo se os spans de `CreateTransaction` e `Outbox` aparecem encadeados na árvore de auditoria distribuída.


