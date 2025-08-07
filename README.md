# 📚 Bookstore App

A REST API for managing a bookstore built with Go, PostgreSQL, and GORM.

## 🚀 Features

- **CRUD Operations**: Create, Read, Update, Delete books
- **RESTful API**: Clean REST endpoints with proper HTTP status codes
- **Database Integration**: PostgreSQL with GORM ORM
- **Comprehensive Testing**: Unit tests, integration tests, and error scenario tests
- **Health Check**: Built-in health check endpoint
- **Docker Support**: PostgreSQL running in Docker container

## 📋 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/books` | Get all books |
| `GET` | `/book` | Get all books (alternative) |
| `GET` | `/book/:id` | Get book by ID |
| `POST` | `/book` | Create new book |
| `PUT` | `/book/:id` | Update book by ID |
| `DELETE` | `/book/:id` | Delete book by ID |

## 🔧 Setup & Installation

### Prerequisites
- Go 1.19 or higher
- Docker and Docker Compose
- Make (for using Makefile commands)

### 1. Clone the repository
```bash
git clone <>
cd bookstore-app
```

### 2. Install dependencies
```bash
go mod download
# or
make deps
```

### 3. Start PostgreSQL database
```bash
docker-compose up -d postgres
# or
make db-up
```

### 4. Run the application
```bash
go run main.go
# or
make run
```

The server will start on `http://localhost:8080`


#### Quick Test Run
```bash
# Run all tests with the test runner script
./run_tests.sh
# or
make test
```

#### Individual Test Categories
```bash
# Unit tests only
go test -v ./pkg/controllers/
# or
make test-unit

# Integration tests only
go test -v .
# or
make test-integration

# Specific test functions
go test -v -run TestBookCRUDFlow .
go test -v -run TestHealthCheck .
go test -v -run TestErrorScenarios .
```

#### Test Coverage
```bash
# Generate coverage report
make test-coverage
# Opens coverage.html in your browser
```

### Test Structure

```
📁 Tests
├── 🧪 pkg/controllers/controllers_test.go  # Unit tests
├── 🌐 integration_test.go                 # Integration tests
└── 📜 run_tests.sh                        # Test runner script
```

### Test Categories Explained

1. **Unit Tests** (`pkg/controllers/controllers_test.go`)
   - Test individual controller functions
   - Mock HTTP requests and validate responses
   - Check error handling for invalid inputs

2. **Integration Tests** (`integration_test.go`)
   - Test complete CRUD workflow
   - Verify database operations
   - Test concurrent operations
   - Validate API endpoint compatibility

3. **Health Check Tests**
   - Ensure service availability endpoint works
   - Validate proper HTTP status codes

4. **Error Scenario Tests**
   - Invalid JSON handling
   - Non-existent resource handling
   - Invalid ID format handling
   - Database error scenarios

## 📊 Test Results Example



## 📚 API Usage Examples

### Create a Book
```bash
curl -X POST http://localhost:8080/book \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language",
    "author": "Alan Donovan",
    "isbn": "9780134190440",
    "publication_year": "2015",
    "genre": "Programming",
    "price": 45.99
  }'
```

### Get All Books
```bash
curl http://localhost:8080/books
```

### Get Book by ID
```bash
curl http://localhost:8080/book/1
```

### Update Book
```bash
curl -X PUT http://localhost:8080/book/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Book Title",
    "author": "Updated Author",
    "isbn": "9780134190440",
    "publication_year": "2024",
    "genre": "Programming",
    "price": 49.99
  }'
```

### Delete Book
```bash
curl -X DELETE http://localhost:8080/book/1
```

### Health Check
```bash
curl http://localhost:8080/health
```

## 🏗️ Project Structure

```
bookstore-app/
├── main.go                     # Application entry point
├── go.mod                      # Go module file
├── go.sum                      # Go dependencies
├── docker-compose.yml          # PostgreSQL container config
├── .env                        # Environment variables
├── init.sql                    # Database initialization
├── Makefile                    # Build and test commands
├── run_tests.sh               # Test runner script
├── integration_test.go         # Integration tests
└── pkg/
    ├── config/
    │   └── config.go          # Database configuration
    ├── models/
    │   └── book.go            # Book model and database operations
    ├── controllers/
    │   ├── controllers.go     # HTTP request handlers
    │   └── controllers_test.go # Unit tests
    └── routes/
        └── routes.go          # Route definitions
```

## 🐳 Docker Commands

```bash
# Start PostgreSQL
docker compose up -d postgres

# View logs
docker compose logs -f postgres

# Stop services
docker compose down

# Reset database (removes all data)
docker compose down -v
```

## 🔄 Development Workflow

```bash
# 1. Start database
make db-up

# 2. Run tests
make test

# 3. Start application
make run

# 4. Make changes and test
make test-unit

# 5. Clean up
make clean
```
