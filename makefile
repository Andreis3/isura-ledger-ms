run-app:
	@echo "Running app"
	@go run cmd/main.go

run-app-logs:
	@echo "Running app export archive logs"
	@go run cmd/main.go > ~/tmp/app/customers-ms.log 2>&1

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

up:
	@docker compose -f docker-compose.yml up -d --build

down:
	@docker compose -f docker-compose.yml down -v

pgpool-logs:
	@echo "Running pgpool logs"
	@docker logs -f pgpool

tag:
	scripts/bump_version.sh

integration:
	@go test ./tests/integration/... --tags=integration -v -count=1
.PHONY: run-app,
		unit,
		unit-cover,
		unit-report,
		integration,
		docker-dev,
		up,
		down,
		tag,
