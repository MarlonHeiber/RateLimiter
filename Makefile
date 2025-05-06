
.PHONY: up down check_bins


up: 
	@echo "Running redis in Docker..."
	@docker-compose up -d
	@echo "Waiting for Redis to start..."
	@echo "You can send requests to the app by typing:"
	@echo 'curl http://localhost:8080/test'
	@echo "or"	
	@echo 'curl -H "API_KEY: Example123" http://localhost:8080/test'
	@echo "Starting the app..."
	@echo 
	@go run ./cmd/main.go


down: check_bins
	@echo "Deleting environment..."
	@docker-compose down