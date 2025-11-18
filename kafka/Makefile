.PHONY: up down logs kafka-ui proto test

up: ## Запустить Kafka инфраструктуру
	docker-compose up -d

down: ## Остановить Kafka инфраструктуру
	docker-compose down

logs: ## Просмотр логов Kafka
	docker-compose logs -f kafka

kafka-ui: ## Открыть Kafka UI (http://localhost:8090)
	@echo "Kafka UI доступен на http://localhost:8090"

test: ## Запустить тесты
	go test ./...

tidy: ## Обновить зависимости
	go mod tidy

help: ## Показать доступные команды
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

