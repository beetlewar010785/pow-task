SERVER_IMAGE = pow-server
CLIENT_IMAGE = pow-client
NETWORK = app_network

lint:
	@golangci-lint run --timeout=1m

test:
	@go test -coverprofile=coverage.out ./...

build-server:
	docker build -f ./docker/server.Dockerfile -t $(SERVER_IMAGE) .

build-client:
	docker build -f ./docker/client.Dockerfile -t $(CLIENT_IMAGE) .

network:
	docker network create $(NETWORK) || true

run-server: network
	docker run --rm -i -d --name server --network $(NETWORK) -p 8080:8080 $(SERVER_IMAGE)

run-client: network
	docker run --rm -i --name client --network $(NETWORK) -e SERVER_ADDRESS="server:8080" $(CLIENT_IMAGE)

clean-docker:
	docker ps -q --filter "name=server" | xargs -r docker stop
	docker ps -q --filter "name=client" | xargs -r docker stop
	docker network rm $(NETWORK) || true
