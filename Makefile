.PHONY: build run test lint migrate-up migrate-down migrate-create docker-up docker-down

# Variables
APP_NAME=api
BUILD_DIR=bin
DB_URL?=postgres://postgres:postgres_password@localhost:5432/shop_db?sslmode=disable
MIGRATIONS_DIR=migrations

build:
	@echo "Building binary..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api/main.go

run: build
	@echo "Running application..."
	./$(BUILD_DIR)/$(APP_NAME)

test:
	@echo "Running unit and integration tests..."
	go test -v -race ./...

lint:
	@echo "Running go vet and gofmt check..."
	go vet ./...
	@test -z "$$(gofmt -l .)" || (echo "Unformatted files found. Run gofmt -w ." && exit 1)

migrate-up:
	@echo "Applying up migrations..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	@echo "Applying down migrations..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-create:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down -v
