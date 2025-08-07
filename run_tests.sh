#!/bin/bash

echo "==================================="
echo "🧪 Bookstore App Test Suite"
echo "==================================="

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2 PASSED${NC}"
    else
        echo -e "${RED}❌ $2 FAILED${NC}"
    fi
}


echo -e "${YELLOW}📋 Checking Prerequisites...${NC}"
if docker compose ps | grep -q "postgres.*Up"; then
    echo -e "${GREEN}✅ PostgreSQL is running${NC}"
else
    echo -e "${YELLOW}⚠️  Starting PostgreSQL...${NC}"
    docker-compose up -d postgres
    echo "⏳ Waiting for PostgreSQL to be ready..."
    sleep 10
fi


echo -e "${YELLOW}🗄️  Setting up test database...${NC}"
docker-compose exec -T postgres psql -U bookstore_user -d bookstore -c "CREATE DATABASE bookstore_test;" 2>/dev/null || true


echo -e "${YELLOW}🔧 Running Unit Tests...${NC}"
go test -v ./pkg/controllers/
unit_exit_code=$?
print_status $unit_exit_code "Unit Tests"

echo -e "${YELLOW}🌐 Running Integration Tests...${NC}"
go test -v -run TestBookCRUDFlow .
integration_exit_code=$?
print_status $integration_exit_code "Integration Tests"


echo -e "${YELLOW}💗 Running Health Check Tests...${NC}"
go test -v -run TestHealthCheck .
health_exit_code=$?
print_status $health_exit_code "Health Check Tests"


echo -e "${YELLOW}🚨 Running Error Scenario Tests...${NC}"
go test -v -run TestErrorScenarios .
error_exit_code=$?
print_status $error_exit_code "Error Scenario Tests"

echo -e "${YELLOW}🔀 Running Endpoint Compatibility Tests...${NC}"
go test -v -run TestMultipleGetEndpoints .
endpoint_exit_code=$?
print_status $endpoint_exit_code "Endpoint Compatibility Tests"

echo -e "${YELLOW}⚡ Running Concurrent Operation Tests...${NC}"
go test -v -run TestConcurrentOperations .
concurrent_exit_code=$?
print_status $concurrent_exit_code "Concurrent Operation Tests"

total_tests=6
failed_tests=0
[ $unit_exit_code -ne 0 ] && ((failed_tests++))
[ $integration_exit_code -ne 0 ] && ((failed_tests++))
[ $health_exit_code -ne 0 ] && ((failed_tests++))
[ $error_exit_code -ne 0 ] && ((failed_tests++))
[ $endpoint_exit_code -ne 0 ] && ((failed_tests++))
[ $concurrent_exit_code -ne 0 ] && ((failed_tests++))

passed_tests=$((total_tests - failed_tests))

echo "==================================="
echo -e "${YELLOW}📊 Test Summary${NC}"
echo "==================================="
echo -e "Total Tests: $total_tests"
echo -e "${GREEN}Passed: $passed_tests${NC}"
echo -e "${RED}Failed: $failed_tests${NC}"

if [ $failed_tests -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}💥 Some tests failed!${NC}"
    exit 1
fi
