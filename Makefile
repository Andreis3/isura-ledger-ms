# ── Variáveis ──────────────────────────────────────────────

DOCKER_COMPOSE = docker compose
SERVICE_NAME = ledger
DB_URL  = postgres://admin:admin@localhost:5432/isura_ledger_main?sslmode=disable
SCHEMA_DIR = db
URL ?= http://localhost:8080/accounts
DURATION ?= 5h

help:
	@echo "======================================================================"
	@echo " 🚀 ISURA LEDGER MS - COMANDOS DISPONÍVEIS"
	@echo "======================================================================"
	@echo ""
	@echo " [ Aplicação & Execução ]"
	@echo "   make run-app          - Executa a aplicação localmente"
	@echo "   make run-race         - Executa a aplicação com detecção de race conditions"
	@echo "   make run-app-logs     - Executa a aplicação exportando logs para arquivo"
	@echo "   make air              - Executa com hot-reload (via Air)"
	@echo ""
	@echo " [ Testes Unitários & Cobertura ]"
	@echo "   make unit             - Roda os testes unitários básicos"
	@echo "   make unit-verbose     - Roda os testes unitários via Ginkgo com race detector"
	@echo "   make unit-cover       - Roda testes unitários medindo cobertura"
	@echo "   make unit-report      - Gera relatório HTML e de funções da cobertura"
	@echo ""
	@echo " [ Testes de Carga (Vegeta) ]"
	@echo "   make test-load        - Roda o teste estável (5.000 req/s, 10.000 conexões)"
	@echo "   make test-extreme     - Roda o teste de alto rendimento (100.000 req/s)"
	@echo ""
	@echo " [ Docker & Infraestrutura ]"
	@echo "   make build            - Faz build da imagem local sem cache"
	@echo "   make up               - Sobe o ambiente completo (DB + App)"
	@echo "   make down             - Para e remove os containers"
	@echo "   make logs             - Exibe os logs do container em tempo real"
	@echo "   make bash             - Entra no terminal do container da aplicação"
	@echo "   make restart          - Para, faz o build e sobe o ambiente novamente"
	@echo ""
	@echo " [ Protobuf & Migrations ]"
	@echo "   make proto-lint       - Valida a sintaxe dos arquivos .proto com Buf"
	@echo "   make proto-gen        - Gera o código Go a partir dos protos"
	@echo "   make migrate          - Aplica as migrations do banco via Atlas"
	@echo "======================================================================"

run-app:
	@echo "Running app"
	@go run cmd/server/main.go
run-race:
	@echo "Running race active"
	go run -race cmd/server/main.go


run-app-logs:
	@echo "Running app export archive logs"
	@go run cmd/main.go > ~/tmp/app/customers-ms.log 2>&1

air:
	@echo "Running with reload"
	@air -c .air.toml

unit:
	@go test ./tests/unit/... --tags=unit -v

unit-verbose:
	ginkgo -r --race --tags=unit --randomize-all --randomize-suites --fail-on-pending

unit-cover:
	@go test ./tests/unit/... -coverpkg ./internal/... --tags=unit

unit-report:
	mkdir -p "coverage" \
	&& go test ./tests/unit/... -coverprofile=coverage/cover.out -coverpkg ./internal/... --tags=unit \
	&& go tool cover -html=coverage/cover.out -o coverage/cover.html \
	&& go tool cover -func=coverage/cover.out -o coverage/cover.functions.html

test-load:
	@echo "🚀 Iniciando teste de carga estável (5k req/s)..."
	go run ./vegeta/account/create_account.go  -rate=1000 -connections=5000 -workers=500 -duration=$(DURATION) -url=$(URL)

test-extreme:
	@echo "🔥 Iniciando teste de carga extremo (100k req/s)..."
	go run ./vegeta/account/create_account.go  -rate=100000 -connections=100000 -workers=500 -duration=$(DURATION) -url=$(URL)

# 1. Build da imagem local usando o Dockerfile.local
build:
	$(DOCKER_COMPOSE) build --no-cache $(SERVICE_NAME)

# 2. Sobe o ambiente completo (DB + App com Air)
up:
	$(DOCKER_COMPOSE) up -d

# 3. Para tudo e remove containers
down:
	$(DOCKER_COMPOSE) down

# 4. Ver logs do Air/App em tempo real
logs:
	$(DOCKER_COMPOSE) logs -f $(SERVICE_NAME)

# 5. Atalho para entrar no container (útil para rodar migrations manuais)
bash:
	docker exec -it $$(docker ps -q -f name=$(SERVICE_NAME)) bash

# 6. Build + Up combinado
restart: down build up logs

# Geração de stubs Protobuf/gRPC usando Buf v2
proto-lint:
	@echo "Checking proto syntax with buf..."
	@buf lint

proto-gen:
	@echo "Generating Go code from protos..."
	@buf generate

migrate:
	atlas schema apply \
	  -u "$(DB_URL)" \
	  --to "file://$(SCHEMA_DIR)"



.PHONY: build,
		up,
		down,
		logs,
		bash,
		restart,
		unit,
		unit-verbose,
		unit-cover,
		unit-report,
		integration,
		proto-lint,
		proto-gen,
		migrate,
		air,
		run-race,
		test-load,
		test-extreme,
		help