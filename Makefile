.PHONY: help build run test test-unit test-integration clean db-up db-down db-reset

help:
	@echo "📚 Bookstore App Commands:"
	@echo "  make build           - Build the application"
	@echo "  make run             - Run the application"
	@echo "  make test            - Run all tests"
	@echo "  make test-unit       - Run unit tests only"
	@echo "  make test-integration- Run integration tests only"
	@echo "  make db-up           - Start PostgreSQL database"
	@echo "  make db-down         - Stop PostgreSQL database"
	@echo "  make db-reset        - Reset database (stop, remove, start)"
	@echo "  make clean           - Clean build artifacts"


build:
	@echo "🏗️  Building application..."
	go build -o bin/bookstore-app .

# Run the application
run: build
	@echo "🚀 Starting bookstore application..."
	./bin/bookstore-app

# Run all tests
test:
	@echo "🧪 Running all tests..."
	./run_tests.sh


test-unit:
	@echo "🔧 Running unit tests..."
	go test -v ./pkg/controllers/

# Run integration tests only  
test-integration:
	@echo "🌐 Running integration tests..."
	go test -v .

db-up:
	@echo "🗄️  Starting PostgreSQL database..."
	docker-compose up -d postgres
	@echo "⏳ Waiting for database to be ready..."
	sleep 5


db-down:
	@echo "🛑 Stopping PostgreSQL database..."
	docker-compose down


db-reset: db-down
	@echo "🔄 Resetting database..."
	docker-compose down -v
	docker-compose up -d postgres
	@echo "⏳ Waiting for database to be ready..."
	sleep 10


clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html


deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy


# Run a quick health check
health-check:
	@echo "💗 Checking application health..."
	curl -f http://localhost:8080/health || echo "❌ Application is not running"


dev: db-up build run

prod-build:
	@echo "🏭 Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/bookstore-app .

logs:
	docker-compose logs -f

test-db-setup:
	@echo "🗄️  Setting up test database..."
	docker-compose exec postgres psql -U bookstore_user -d bookstore -c "CREATE DATABASE IF NOT EXISTS bookstore_test;"
