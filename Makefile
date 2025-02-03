SERVER_IMAGE = pow-server
CLIENT_IMAGE = pow-client
NETWORK = app_network

lint:
	@golangci-lint run --timeout=1m

test:
	@go test -coverprofile=coverage.out ./internal/...

build-server:
	docker build -f ./docker/server.Dockerfile -t $(SERVER_IMAGE) .

build-client:
	docker build -f ./docker/client.Dockerfile -t $(CLIENT_IMAGE) .

network:
	docker network ls | grep -wq $(NETWORK) || docker network create $(NETWORK)

run-server: network
	docker run --rm -d --name server --network $(NETWORK) -p 8080:8080 -e LOG_LEVEL=4 $(SERVER_IMAGE)

run-client: network
	docker run --rm --network $(NETWORK) -e SERVER_ADDRESS="server:8080" -e LOG_LEVEL=4 $(CLIENT_IMAGE)

clean-docker:
	docker ps -q --filter "name=server" | xargs -r docker stop
	docker ps -q --filter "name=client" | xargs -r docker stop
	docker network rm $(NETWORK) || true
