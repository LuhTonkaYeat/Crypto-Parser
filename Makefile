.PHONY: build server cli test clean stop status

build:
	docker build -t crypto-parser .

server:
	docker run -d --name crypto-api -p 8080:8080 crypto-parser server
	@echo "Server started on http://localhost:8080"
	@echo "Stop with: make stop"

stop:
	docker stop crypto-api || true
	docker rm crypto-api || true

cli:
	docker run --rm --network="host" crypto-parser cli $(COIN)

status:
	@echo "Checking server status..."
	@curl -s http://localhost:8080/health || echo "Server not running"

test: build
	@echo "Starting server..."
	docker run -d --name crypto-api -p 8080:8080 crypto-parser server
	sleep 2
	@echo "\n=== Testing btc ==="
	docker run --rm --network="host" crypto-parser cli btc
	@echo "\n=== Testing health ==="
	curl -s http://localhost:8080/health
	@echo ""
	docker stop crypto-api
	docker rm crypto-api
	@echo "\n✅ Tests passed"

clean:
	docker stop crypto-api 2>/dev/null || true
	docker rm crypto-api 2>/dev/null || true
	docker rmi crypto-parser 2>/dev/null || true