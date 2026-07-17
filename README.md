# isura-ledger-ms

The double-entry accounting engine of Isura Bank. **Every financial movement goes through here.**

This service is the financial source of truth — it records all debit and credit entries, maintains real-time account balances, guarantees atomic double-entry bookkeeping, and publishes domain events reliably via the Transactional Outbox pattern.

---

## Responsibilities

- Double-entry bookkeeping — debit + credit = 0, always atomic via a single PostgreSQL transaction
- Real-time balance per account
- Transaction idempotency via `idempotency_key` with `UNIQUE CONSTRAINT`
- Hold/Release for fund reservations (card authorizations, pending Pix)
- Reliable event publishing via Transactional Outbox → Kafka
- Accounting representation of accounts (asset, liability, revenue, expense)

---

## Architecture

This service follows **Hexagonal Architecture** with **tactical Domain-Driven Design**, organized around three bounded contexts internally: `transaction`, `account`, and `outbox`.

### Dependency rule

```
transport → application → domain ← infrastructure
```

The `domain` layer imports nothing external. The `infrastructure` layer implements the interfaces defined in `domain`. The `transport` layer knows only about `application` use cases.

### Project structure

```
isura-ledger-ms/
├── .github/
│   └── workflows/
│       └── golang-build-test.yaml    # CI/CD pipeline de testes e build
├── cmd/
│   └── server/
│       └── main.go                  # Ponto de entrada (inicialização do container/fx)
├── db/
│   ├── accounts.pg.hcl              # Definição Atlas HCL da tabela de contas
│   ├── entries.pg.hcl               # Definição Atlas HCL da tabela de lançamentos
│   ├── outbox.pg.hcl                # Definição Atlas HCL da tabela outbox
│   ├── schema.pg.hcl                # Esquema público do banco de dados
│   └── transactions.pg.hcl          # Definição Atlas HCL da tabela de transações
├── docker/
│   └── tempo/
│       └── tempo.yaml               # Configuração do Grafana Tempo (Traces)
├── internal/
│   ├── application/
│   │   ├── command/
│   │   │   ├── create_account.go    # Use Case: Criação de conta contábil
│   │   │   ├── create_transaction.go# Use Case: Criação de transação double-entry
│   │   │   └── mask.go              # Helpers de mascaramento/logs sensíveis
│   │   ├── event/
│   │   │   └── transaction_create.go# Contrato/Payload do evento de domínio gerado
│   │   ├── logger.go                # Interface/Contrato de logs da aplicação
│   │   ├── metrics.go               # Interface/Contrato de métricas da aplicação
│   │   ├── tracer.go                # Interface/Contrato de traces (OpenTelemetry)
│   │   └── uow.go                   # Interface do Unit of Work
│   ├── domain/
│   │   ├── account/
│   │   │   ├── account.go           # Entidade Account e validações de saldo
│   │   │   └── repository.go        # Interface do repositório de contas
│   │   ├── fault/
│   │   │   ├── fault.go             # Engine customizável de erros estruturados
│   │   │   └── sentinel.go          # Erros sentinelas globais do negócio
│   │   ├── money/
│   │   │   └── money.go             # Value Object Money (Cents + Currency)
│   │   ├── outbox/
│   │   │   ├── outbox.go            # Agregado Outbox e Máquina de Estados
│   │   │   └── repository.go        # Interface do repositório outbox
│   │   └── transaction/
│   │       ├── entry.go             # Entidade Entry (Debito/Credito)
│   │       ├── repository.go        # Interface do repositório de transações
│   │       ├── transaction.go       # Agregado Root Transaction + State machine
│   │       └── types.go             # IDs fortemente tipados (TransactionID, etc.)
│   ├── infra/
│   │   ├── configs/
│   │   │   └── configs.go           # Carregamento de variáveis via Viper
│   │   ├── logger/
│   │   │   └── logger.go            # Implementação slog (JSON/Tint dual handler)
│   │   ├── observability/
│   │   │   ├── otel_tracer.go       # Implementação do Tracer OpenTelemetry
│   │   │   └── prometheus.go        # Registro do subsistema de métricas
│   │   ├── postgres/
│   │   │   ├── database/
│   │   │   │   ├── helper.go        # Tratamentos específicos do driver pgx
│   │   │   │   └── querier.go       # Interface unificada para DB Pool e Tx
│   │   │   ├── model/
│   │   │   │   ├── account.go       # Model do banco para Account
│   │   │   │   ├── entry.go         # Model do banco para Entry
│   │   │   │   ├── outbox.go        # Model do banco para Outbox
│   │   │   │   └── transaction.go   # Model do banco para Transaction
│   │   │   ├── repository/
│   │   │   │   ├── observability/
│   │   │   │   │   ├── account_observability.go     # Decorator para tracing de contas
│   │   │   │   │   ├── outbox_observability.go      # Decorator para tracing de outbox
│   │   │   │   │   └── transaction_observability.go # Decorator para tracing de transações
│   │   │   │   ├── account.go       # Repositório Postgres para contas
│   │   │   │   ├── outbox.go        # Repositório Postgres para outbox
│   │   │   │   ├── resolve_db.go    # Helper para extrair Tx ativa do Context
│   │   │   │   └── transaction.go   # Repositório Postgres para transações
│   │   │   ├── uow/
│   │   │   │   └── uow.go           # Implementação concreta do UoW com pgx.Tx
│   │   │   └── postgres.go          # Inicialização e ping do pool do Postgres
│   │   └── server/
│   │       ├── base_deps.go         # Provider Fx para infra básica (logger, config)
│   │       ├── composer.go          # Wires/Módulos Fx que orquestram a injeção
│   │       ├── graceful_shutdown.go # Controle de encerramento limpo do gRPC/HTTP
│   │       ├── grpc_server.go       # Ciclo de vida do servidor gRPC
│   │       └── http_server.go       # Ciclo de vida do servidor HTTP (Métricas/Health)
│   ├── transport/
│   │   ├── grpc/
│   │   │   ├── handler/
│   │   │   │   └── create_account_handler.go # Traduz Protobuf ↔ Command de Conta
│   │   │   ├── interceptor/
│   │   │   │   ├── logging.go       # Interceptor gRPC de logs
│   │   │   │   ├── metrics.go       # Interceptor gRPC de Prometheus
│   │   │   │   └── tracing.go       # Interceptor gRPC de Spans/Traces
│   │   │   ├── pb/ledger/v1/        # Arquivos `.go` auto-gerados pelo protoc/buf
│   │   │   │   ├── account.pb.go
│   │   │   │   ├── ledger.pb.go
│   │   │   │   ├── ledger_grpc.pb.go
│   │   │   │   └── transaction.pb.go
│   │   │   ├── translator/
│   │   │   │   └── fault_translator.go # Mapeia erros de domínio para gRPC Codes
│   │   │   ├── ledger_module.go     # Módulo gRPC do Ledger
│   │   │   ├── module.go            # Registro de Interceptors no Fx
│   │   │   └── server_registry.go   # Liga o Handler gRPC gerado ao Servidor
│   │   └── rest/
│   │       ├── handler/
│   │       │   └── healthcheck_handler.go
│   │       ├── module/
│   │       │   ├── healthcheck_module.go
│   │       │   └── metrics_module.go
│   │       ├── types/
│   │       │   └── route.go         # Tipagem para acoplamento de rotas HTTP
│   │       ├── register.go          # Registrador do roteador Chi
│   │       └── setup.go             # Criação e configuração do Chi Router
│   └── tests/
│       └── unit/domain/             # Suítes de testes unitários do domínio
│           ├── account/
│           ├── money/
│           ├── outbox/
│           └── transaction/
│               ├── suite_test.go
│               ├── transaction_test.go
│               └── types_test.go
├── proto/ledger/v1/
│   ├── account.proto                # Estrutura de mensagens de conta
│   ├── ledger.proto                 # Definição dos RPCs do serviço
│   └── transaction.proto            # Estrutura de mensagens de transação
├── .air.toml                        # Hot reload para desenvolvimento local
├── .gitignore
├── buf.gen.yaml                     # Configuração do gerador do Buf v2
├── buf.yaml                         # Configuração do módulo Protobuf do Buf v2
├── docker-compose.yml               # Postgres, Prometheus, Tempo, Grafana
├── Dockerfile                       # Build multi-stage para produção
├── Dockerfile.local                 # Build otimizado para o Air local
├── go.mod
├── go.sum
├── Makefile                         # Comandos para rodar migrations, buf, e testes
├── prometheus.yml                   # Configuração de scraping das métricas
└── README.md
```

---

## Domain model

### Money — Value Object

Monetary values are stored as `int64` cents — never `float64`. This eliminates floating-point precision errors in financial calculations.

```
Money { amount int64, currency Currency }

BRL → "BRL" | USD → "USD" | EUR → "EUR"
```

Operations: `Add`, `Subtract`, `Equal`, `IsSufficientBalance`, `IsZero`, `IsNegative`, `IsPositive`, `String` (e.g. `"100.50 BRL"`).

### Transaction — Aggregate Root

The central aggregate. Enforces double-entry invariants and owns its state machine.

```
PENDING → COMPLETED  ✓
PENDING → FAILED     ✓
COMPLETED → any      ✗
FAILED → any         ✗
```

A `Transaction` always contains exactly two `Entry` records — one `DEBIT` and one `CREDIT` with equal amounts. These invariants are enforced by `AddEntry()` before any persistence occurs.

### Entry — Entity

Represents one side of a double-entry ledger. Belongs to a `Transaction` aggregate — never created independently.

```
Entry {
  ID             EntryID
  TransactionID  TransactionID
  AccountID      AccountID
  Direction      DEBIT | CREDIT
  Amount         Money
  IdempotencyKey string
  CreatedAt      time.Time
}
```

### Account — Entity

Accounting representation of a bank account within the ledger. Not the same as the customer-facing account in `isura-account-ms` — the ledger holds only what it needs for bookkeeping.

```
Account {
  ID             AccountID
  ExternalID     string      // ID from isura-account-ms (correlation key)
  AccountingType ASSET | LIABILITY | REVENUE | EXPENSE
  Balance        Money
  CreatedAt      time.Time
  UpdatedAt      time.Time
}
```

### OutboxEvent — Aggregate

Ensures reliable event delivery to Kafka without dual writes. Persisted in the same PostgreSQL transaction as the `Transaction`. A background relay reads `PENDING` events with `SELECT FOR UPDATE SKIP LOCKED` and publishes to Kafka.

```
PENDING → FAILED    ✓  (publish attempt failed)
PENDING → SUCCESS   ✓  (published successfully)
FAILED  → PENDING   ✓  (retry — if Attempts < MaxAttempts)
SUCCESS → any       ✗
```

---

## Key design decisions

### Transactional Outbox

All writes within `CreateTransaction` happen in a single PostgreSQL transaction:

```
BEGIN
  INSERT INTO transactions ...
  INSERT INTO entries ...      (debit)
  INSERT INTO entries ...      (credit)
  UPDATE accounts SET balance  (debit account)
  UPDATE accounts SET balance  (credit account)
  INSERT INTO outbox_events ... (event payload)
COMMIT
```

If the commit fails, nothing is persisted — including the outbox event. If the commit succeeds, the relay will eventually publish the event to Kafka. No dual write, no inconsistency.

### Idempotency

Every transaction carries an `idempotency_key`. A `UNIQUE CONSTRAINT` on the `transactions` table ensures that concurrent retries with the same key result in exactly one committed transaction — even if two requests race past the application-level `ExistsByIdempotencyKey` check.

### Unit of Work

The `UnitOfWork` interface wraps the PostgreSQL transaction lifecycle. All repository calls within a use case's `Execute` method receive a `context.Context` carrying the active `pgx.Tx`. Each repository's `resolveDB` method picks the transaction over the connection pool when present.

```go
return c.uow.WithTransaction(ctx, func(ctxTx context.Context) error {
    // all writes here share the same pgx.Tx
    c.transactionRepo.Save(ctxTx, tx)
    c.accountRepo.UpdateBalance(ctxTx, ...)
    c.outboxRepo.Save(ctxTx, event)
    return nil
})
```

### CQRS

Write operations live in `application/command/`, read operations in `application/query/`. Commands return only `error`. Queries return `(Result, error)`.

---

## Tech stack

| Layer | Technology                          |
|---|-------------------------------------|
| Transport | gRPC + Protobuf                     |
| Language | Go 1.26.4                           |
| Persistence | PostgreSQL 16 + pgx/v5              |
| Migrations | golang-migrate                      |
| Events | Apache Kafka (Transactional Outbox) |
| Observability | OpenTelemetry                       |
| Container | Docker + Kubernetes                 |

### Key dependencies

```
github.com/jackc/pgx/v5          # PostgreSQL driver
github.com/google/uuid           # UUID generation
github.com/golang-migrate/migrate # Database migrations
google.golang.org/grpc           # gRPC server
google.golang.org/protobuf       # Protobuf serialization
go.opentelemetry.io/otel         # Observability
```

---

## Running locally

### Prerequisites

- Go 1.26.4+
- Docker and Docker Compose
- `golang-migrate` CLI

### Setup

```bash
# clone
git clone https://github.com/andreis3/isura-ledger-ms
cd isura-ledger-ms

# start PostgreSQL and Kafka
docker compose up -d

# run migrations
make migrate-up

# start the service
go run ./cmd/server/main.go
```

### Makefile commands

```bash
make migrate-up      # apply all pending migrations
make migrate-down    # rollback last migration
make test            # run unit tests
make test-int        # run integration tests
make proto           # regenerate protobuf files
make lint            # run golangci-lint
make build           # build binary
```

---

## Environment variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=isura_ledger
DB_USER=postgres
DB_PASSWORD=postgres
DB_MAX_CONNS=20
DB_MIN_CONNS=5

# gRPC
GRPC_PORT=50051

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_LEDGER_EVENTS=ledger.events

# Outbox relay
OUTBOX_RELAY_INTERVAL_MS=500
OUTBOX_RELAY_BATCH_SIZE=50
OUTBOX_MAX_ATTEMPTS=3

# Observability
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
OTEL_SERVICE_NAME=isura-ledger-ms
```

---

## Database schema

### tables

| Table | Description |
|---|---|
| `transactions` | Aggregate root — one record per transaction |
| `entries` | Double-entry records — always two per transaction |
| `accounts` | Accounting representation of accounts |
| `outbox_events` | Transactional outbox — pending Kafka events |

### Key constraints

```sql
-- idempotency guarantee
UNIQUE (idempotency_key) ON transactions

-- double-entry integrity
FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON entries

-- outbox ordering
INDEX (status, created_at) ON outbox_events
```

---

## gRPC API

```protobuf
service LedgerService {
  rpc CreateTransaction (CreateTransactionRequest) returns (CreateTransactionResponse);
  rpc GetBalance        (GetBalanceRequest)         returns (GetBalanceResponse);
  rpc GetTransaction    (GetTransactionRequest)      returns (GetTransactionResponse);
  rpc CreateAccount     (CreateAccountRequest)       returns (CreateAccountResponse);
}
```

---

## Testing

```bash
# unit tests (domain logic — no external dependencies)
go test ./internal/domain/...

# integration tests (requires Docker)
go test ./internal/infrastructure/... -tags=integration

# all tests with race detector
go test -race ./...
```

---

## References

- *Implementing Domain-Driven Design* — Vaughn Vernon
- *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi
- *Concurrency in Go* — Katherine Cox-Buday
- *The Go Programming Language* — Donovan & Kernighan
- [pgx/v5 documentation](https://github.com/jackc/pgx)
- [gRPC Go documentation](https://grpc.io/docs/languages/go/)

---

## License

MIT — see [LICENSE](./LICENSE) for details.
