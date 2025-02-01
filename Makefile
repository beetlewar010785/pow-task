SERVER_IMAGE = pow-server
CLIENT_IMAGE = pow-client
NETWORK = app_network

lint:
	@golangci-lint run --timeout=1m

test:
	@go test -v ./...

build-server:
	docker build -f ./docker/server.Dockerfile -t $(SERVER_IMAGE) .

build-client:
	docker build -f ./docker/client.Dockerfile -t $(CLIENT_IMAGE) .

network:
	docker network create $(NETWORK) || true

run-server: network
	docker run --rm -it -d --name server --network $(NETWORK) -p 8080:8080 $(SERVER_IMAGE)

run-client: network
	docker run --rm -it --name client --network $(NETWORK) -e SERVER_ADDRESS="server:8080" $(CLIENT_IMAGE)

clean:
	@go clean -testcache
	@if [ -n "$$(docker ps -aq --filter ancestor=$(SERVER_IMAGE))" ]; then docker rm -f $$(docker ps -aq --filter ancestor=$(SERVER_IMAGE)); fi
	@if [ -n "$$(docker ps -aq --filter ancestor=$(CLIENT_IMAGE))" ]; then docker rm -f $$(docker ps -aq --filter ancestor=$(CLIENT_IMAGE)); fi
	docker network rm $(NETWORK) || true
