.PHONY: proto deploy down_prod up_prod

# Находим все файлы с расширением .proto внутри папки api/proto/
PROTO_FILES := $(shell find api/proto -name "*.proto")

proto:
	@if [ -z "$(PROTO_FILES)" ]; then \
		echo "No .proto files found!"; \
	else \
		protoc --go_out=. --go_opt=paths=source_relative \
		       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		       $(PROTO_FILES); \
		echo "Successfully compiled all proto files."; \
	fi

down_db:
	docker compose down -v

up_db:
	docker compose up -d

run_dev_gateway:
	cd services/tg-bot-gateway && go run cmd/main.go

run_dev_workout:
	cd services/workout-service && go run cmd/main.go

# --- PROD TARGETS ---
deploy:
	docker compose -f docker-compose.prod.yml up -d --build

# Полный останов продакшена
down_prod:
	docker compose -f docker-compose.prod.yml down

# Просмотр логов бэкенда на проде в реальном времени
logs_prod:
	docker compose -f docker-compose.prod.yml logs -f