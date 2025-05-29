.PHONY: build test run docker-build docker-run clean

# Development
run:
	docker compose up

build:
	docker compose build

# Testing
test-backend:
	cd backend && go test -v ./...

test-frontend:
	cd frontend && npm test

test: test-backend test-frontend

# Cleanup
clean:
	docker compose down
	rm -rf data/*
	git checkout -- data/.gitkeep

# Production
prod-build:
	docker compose -f docker-compose.yml build

prod-run:
	docker compose -f docker-compose.yml up -d

prod-logs:
	docker compose -f docker-compose.yml logs -f

prod-stop:
	docker compose -f docker-compose.yml down 