.PHONY: help start stop test benchmark clean logs deps

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install Go dependencies
	go mod tidy
	go mod download

start: ## Start all database services
	docker-compose up -d
	@echo "Waiting for databases to be ready..."
	@echo "Waiting for TimescaleDB..."
	@until docker-compose exec timescaledb pg_isready -U postgres; do echo "TimescaleDB not ready, waiting..."; sleep 2; done
	@echo "Waiting for ClickHouse..."
	@until docker-compose exec clickhouse clickhouse-client --query "SELECT 1" > /dev/null 2>&1; do echo "ClickHouse not ready, waiting..."; sleep 2; done
	@echo "Creating ClickHouse database..."
	@docker-compose exec clickhouse clickhouse-client --query "CREATE DATABASE IF NOT EXISTS benchmark"
	@echo "Services are ready!"

stop: ## Stop all database services
	docker-compose down

logs: ## Show logs from all services
	docker-compose logs -f

test: ## Run connection tests
	go test -v -run TestDatabaseConnections

benchmark: ## Run all benchmark tests
	go test -v -bench=. -benchmem -benchtime=10s

benchmark-timescale: ## Run only TimescaleDB benchmarks
	go test -v -bench=BenchmarkTimescaleDB -benchmem -benchtime=10s

benchmark-clickhouse: ## Run only ClickHouse benchmarks
	go test -v -bench=BenchmarkClickHouse -benchmem -benchtime=10s

benchmark-async: ## Run only async insert benchmarks
	go test -v -bench=Async -benchmem -benchtime=10s

clean: ## Clean up containers and volumes
	docker-compose down -v
	docker system prune -f

restart: stop start ## Restart all services

status: ## Show status of all containers
	docker-compose ps